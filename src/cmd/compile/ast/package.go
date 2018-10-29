package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	compileCommon "gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"path/filepath"
)

type Package struct {
	Name                         string
	LoadedPackages               map[string]*Package
	loadedClasses                map[string]*Class
	Block                        Block // package always have a default block
	files                        map[string]*SourceFile
	InitFunctions                []*Function
	NErrors2Stop                 int // number of errors should stop compile
	errors                       []error
	TriggerPackageInitMethodName string //
	unUsedPackage                map[string]*Import
	statementLevelFunctions      []*Function
	statementLevelClass          []*Class
}

func (p *Package) isSame(compare *Package) bool {
	return p.Name == compare.Name
}

func (p *Package) loadCorePackage() error {
	if p.Name == common.CorePackage {
		return nil
	}
	pp, err := p.load(common.CorePackage)
	if err != nil {
		return err
	}
	lucyBuildInPackage = pp.(*Package)
	for _, v := range lucyBuildInPackage.Block.Variables {
		v.IsBuildIn = true
	}
	for _, v := range lucyBuildInPackage.Block.Constants {
		v.IsBuildIn = true
	}
	for _, v := range lucyBuildInPackage.Block.Enums {
		v.IsBuildIn = true
	}
	for _, v := range lucyBuildInPackage.Block.Classes {
		v.IsBuildIn = true
	}
	for _, v := range lucyBuildInPackage.Block.Functions {
		v.IsBuildIn = true
		v.LoadedFromCorePackage = true
	}
	for _, v := range lucyBuildInPackage.Block.TypeAliases {
		v.IsBuildIn = true
	}
	p.Block.Outer = &lucyBuildInPackage.Block
	return nil
}

func (p *Package) getImport(file string, accessName string) *Import {
	if p.files == nil {
		return nil
	}
	if file, ok := p.files[file]; ok == false {
		return nil
	} else {
		return file.Imports[accessName]
	}
}

func (p *Package) mkInitFunctions(bs []*Block) {
	p.InitFunctions = make([]*Function, len(bs))
	for k, b := range bs {
		b.IsFunctionBlock = true
		f := &Function{}
		b.Fn = f
		f.Pos = b.Pos
		f.Block = *b
		p.InitFunctions[k] = f
		f.isPackageInitBlockFunction = true
		f.Used = true
	}
}

func (p *Package) shouldStop(errs []error) bool {
	return len(p.errors)+len(errs) >= p.NErrors2Stop
}

func (p *Package) TypeCheck() []error {
	if p.NErrors2Stop <= 2 {
		p.NErrors2Stop = 10
	}
	p.errors = []error{}
	p.errors = append(p.errors, p.Block.checkConstants()...)
	for _, v := range p.Block.Enums {
		v.Name = p.Name + "/" + v.Name
		p.errors = append(p.errors, v.check()...)
	}
	for _, v := range p.Block.TypeAliases {
		err := v.resolve(&PackageBeenCompile.Block)
		if err != nil {
			p.errors = append(p.errors, err)
		}
	}
	for _, v := range p.Block.Functions {
		if v.IsBuildIn {
			continue
		}
		v.Block.inherit(&p.Block)
		v.Block.InheritedAttribute.Function = v
		v.checkParametersAndReturns(&p.errors, false, false)
		if v.IsGlobalMain() {
			defineMainOK := false
			if len(v.Type.ParameterList) == 1 {
				defineMainOK = v.Type.ParameterList[0].Type.Type == VariableTypeArray &&
					v.Type.ParameterList[0].Type.Array.Type == VariableTypeString
			}
			if defineMainOK == false {
				p.errors = append(p.errors,
					fmt.Errorf("%s function '%s' expect declared as 'main(args []string)'",
						errMsgPrefix(v.Pos), MainFunctionName))
			}
		}
		if p.shouldStop(nil) {
			return p.errors
		}
	}
	for _, v := range p.Block.Classes {
		v.Name = p.Name + "/" + v.Name
		p.errors = append(p.errors, v.Block.checkConstants()...)
		v.mkDefaultConstruction()
		v.Block.inherit(&PackageBeenCompile.Block)
		v.Block.InheritedAttribute.Class = v
	}

	for _, v := range p.Block.Classes {
		err := v.resolveFather()
		if err != nil {
			p.errors = append(p.errors, err)
		}
		p.errors = append(p.errors, v.resolveInterfaces()...)
		p.errors = append(p.errors, v.resolveFieldsAndMethodsType()...)
	}

	for _, v := range p.Block.Classes {
		es := v.checkPhase1()
		p.errors = append(p.errors, es...)
		if p.shouldStop(nil) {
			return p.errors
		}
	}
	for _, v := range p.Block.Functions {
		if v.TemplateFunction != nil {
			continue
		}
		p.errors = append(p.errors, v.checkReturnVarExpression()...)
	}
	for _, v := range p.InitFunctions {
		p.errors = append(p.errors, v.check(&p.Block)...)
		if p.shouldStop(nil) {
			return p.errors
		}
	}
	for _, v := range p.Block.Classes {
		p.errors = append(p.errors, v.checkPhase2()...)
		if p.shouldStop(nil) {
			return p.errors
		}
	}
	for _, v := range p.Block.Functions {
		if v.IsBuildIn {
			continue
		}
		if v.TemplateFunction != nil {
			continue
		}
		v.checkBlock(&p.errors)
		if PackageBeenCompile.shouldStop(nil) {
			return p.errors
		}
	}
	for _, v := range p.statementLevelFunctions {
		v.IsClosureFunction = v.Closure.CaptureCount(v) > 0
	}
	for _, v := range p.statementLevelClass {
		for f, meta := range v.closure.Functions {
			if f.IsClosureFunction == false {
				continue
			}
			p.errors = append(p.errors,
				fmt.Errorf("%s trying to access capture function '%s' from outside",
					meta.pos.ErrMsgPrefix(), f.Name))
		}
	}
	p.errors = append(p.errors, p.checkUnUsedPackage()...)
	return p.errors
}

/*
	load package or class
*/
func (p *Package) load(resource string) (interface{}, error) {
	if resource == "" {
		panic("null string")
	}
	if p.loadedClasses == nil {
		p.loadedClasses = make(map[string]*Class)
	}
	if t, ok := p.loadedClasses[resource]; ok {
		return t, nil
	}
	if p.LoadedPackages == nil {
		p.LoadedPackages = make(map[string]*Package)
	}
	if t, ok := p.LoadedPackages[resource]; ok {
		return t, nil
	}
	t, err := ImportsLoader.LoadImport(resource)
	if pp, ok := t.(*Package); ok && pp != nil {
		PackageBeenCompile.LoadedPackages[resource] = pp
		p.mkClassCache(pp)
	}
	if c, ok := t.(*Class); ok && c != nil {
		if c.IsJava == false {
			return nil, fmt.Errorf("load lucy class not allow")
		}
		PackageBeenCompile.loadedClasses[resource] = c
	}
	return t, err
}

func (p *Package) checkUnUsedPackage() []error {
	if compileCommon.CompileFlags.DisableCheckUnUse {
		return nil
	}
	errs := []error{}
	for _, v := range p.files {
		for _, i := range v.Imports {
			if i.Used == false {
				errs = append(errs, fmt.Errorf("%s '%s' imported not used",
					errMsgPrefix(i.Pos), i.Import))
			}
		}
	}
	for _, i := range p.unUsedPackage {
		pp, err := p.load(i.Import)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %v",
				errMsgPrefix(i.Pos), err))
			continue
		}
		if ppp, ok := pp.(*Package); ok == false {
			errs = append(errs, fmt.Errorf("%s '%s' not a package",
				errMsgPrefix(i.Pos), i.Import))
		} else {
			if ppp.TriggerPackageInitMethodName == "" {
				errs = append(errs,
					fmt.Errorf("%s  package named '%s' have no global vars and package "+
						"init blocks, no need to trigger package init method",
						errMsgPrefix(i.Pos), i.Import))
			}
		}
	}
	return errs
}

func (p *Package) loadClass(className string) (*Class, error) {
	if p.loadedClasses == nil {
		p.loadedClasses = make(map[string]*Class)
	}
	if c, ok := p.loadedClasses[className]; ok && c != nil {
		return c, nil
	}
	c, err := ImportsLoader.LoadImport(className)
	if err != nil {
		return nil, err
	}
	if t, ok := c.(*Class); ok == false || t == nil {
		return nil, fmt.Errorf("'%s' is not class", className)
	}
	cc := c.(*Class)
	p.loadedClasses[className] = cc
	return cc, nil
}

func (p *Package) mkClassCache(loadedPackage *Package) {
	for _, v := range loadedPackage.Block.Classes {
		p.loadedClasses[v.Name] = v // binary name
	}
}

func (p *Package) insertImport(i *Import) error {
	if p.files == nil {
		p.files = make(map[string]*SourceFile)
	}
	x, ok := p.files[i.Pos.Filename]
	if ok == false {
		x = &SourceFile{}
		p.files[i.Pos.Filename] = x
	}
	return x.insertImport(i)
}

//different from different source file
type SourceFile struct {
	Imports            map[string]*Import // accessName -> *Import
	ImportsByResources map[string]*Import // resourceName -> *Import
}

func (s *SourceFile) insertImport(i *Import) error {
	if s.Imports == nil {
		s.Imports = make(map[string]*Import)
	}
	if s.ImportsByResources == nil {
		s.ImportsByResources = make(map[string]*Import)
	}
	if err := i.MkAccessName(); err != nil {
		return err
	}
	if _, ok := s.Imports[i.Import]; ok {
		return fmt.Errorf("%s '%s' reimported",
			i.Pos.ErrMsgPrefix(), i.Import)
	}
	if _, ok := s.ImportsByResources[i.Alias]; ok {
		return fmt.Errorf("%s '%s' reimported",
			i.Pos.ErrMsgPrefix(), i.Alias)
	}
	s.Imports[i.Import] = i
	s.Imports[i.Alias] = i
	return nil
}

type Import struct {
	Alias  string
	Import string // full name
	Pos    *Pos
	Used   bool
}

/*
	import "github.com/lucy" should access by lucy.doSomething()
	import "github.com/std" as std2 should access by std2.doSomething()
*/
func (i *Import) MkAccessName() error {
	if i.Alias != "" {
		return nil
	}
	if false == PackageNameIsValid(i.Import) {
		return fmt.Errorf("%s '%s' is not a valid name",
			i.Pos.ErrMsgPrefix(), i.Import)
	}
	i.Alias = filepath.Base(i.Import)
	return nil
}

type RedeclareError struct {
	Name      string
	Positions []*Pos
	Types     []string
}

func (r *RedeclareError) Error() error {
	s := fmt.Sprintf("name '%s' defined  multi times,which are:\n", r.Name)
	for k, v := range r.Positions {
		s += fmt.Sprintf("\t%s '%s' named '%s'\n", errMsgPrefix(v), r.Types[k], r.Name)
	}
	return errors.New(s)
}
