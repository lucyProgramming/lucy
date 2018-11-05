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

func (this *Package) isSame(compare *Package) bool {
	return this.Name == compare.Name
}

func (this *Package) markBuildIn() {
	for _, v := range this.Block.Variables {
		v.IsBuildIn = true
	}
	for _, v := range this.Block.Constants {
		v.IsBuildIn = true
	}
	for _, v := range this.Block.Enums {
		v.IsBuildIn = true
	}
	for _, v := range this.Block.Classes {
		v.IsBuildIn = true
	}
	for _, v := range this.Block.Functions {
		v.IsBuildIn = true
		v.LoadedFromCorePackage = true
	}
	for _, v := range this.Block.TypeAliases {
		v.IsBuildIn = true
	}
}
func (this *Package) loadCorePackage() error {
	if this.Name == common.CorePackage {
		return nil
	}
	pp, err := this.load(common.CorePackage)
	if err != nil {
		return err
	}
	lucyBuildInPackage = pp.(*Package)
	lucyBuildInPackage.markBuildIn()
	this.Block.Outer = &lucyBuildInPackage.Block
	return nil
}

func (this *Package) getImport(file string, accessName string) *Import {
	if this.files == nil {
		return nil
	}
	if file, ok := this.files[file]; ok == false {
		return nil
	} else {
		return file.Imports[accessName]
	}
}

func (this *Package) mkInitFunctions(bs []*Block) {
	this.InitFunctions = make([]*Function, len(bs))
	for k, b := range bs {
		b.IsFunctionBlock = true
		f := &Function{}
		b.Fn = f
		f.Pos = b.Pos
		f.Block = *b
		this.InitFunctions[k] = f
		f.isPackageInitBlockFunction = true
		f.Used = true
	}
}

func (this *Package) shouldStop(errs []error) bool {
	return len(this.errors)+len(errs) >= this.NErrors2Stop
}

func (this *Package) TypeCheck() []error {
	if this.NErrors2Stop <= 2 {
		this.NErrors2Stop = 10
	}
	this.errors = []error{}
	this.errors = append(this.errors, this.Block.checkConstants()...)
	for _, v := range this.Block.Enums {
		v.Name = this.Name + "/" + v.Name
		this.errors = append(this.errors, v.check()...)
	}
	for _, v := range this.Block.TypeAliases {
		err := v.resolve(&PackageBeenCompile.Block)
		if err != nil {
			this.errors = append(this.errors, err)
		}
	}
	for _, v := range this.Block.Functions {
		if v.IsBuildIn {
			continue
		}
		v.Block.inherit(&this.Block)
		v.Block.InheritedAttribute.Function = v
		v.checkParametersAndReturns(&this.errors, false, false)
		if v.IsGlobalMain() {
			defineMainOK := false
			if len(v.Type.ParameterList) == 1 {
				defineMainOK = v.Type.ParameterList[0].Type.Type == VariableTypeArray &&
					v.Type.ParameterList[0].Type.Array.Type == VariableTypeString
			}
			if defineMainOK == false {
				this.errors = append(this.errors,
					fmt.Errorf("%s function '%s' expect declared as 'main(args []string)'",
						errMsgPrefix(v.Pos), MainFunctionName))
			}
		}
		if this.shouldStop(nil) {
			return this.errors
		}
	}
	for _, v := range this.Block.Classes {
		v.Name = this.Name + "/" + v.Name
		this.errors = append(this.errors, v.Block.checkConstants()...)
		v.mkDefaultConstruction()
		v.Block.inherit(&PackageBeenCompile.Block)
		v.Block.InheritedAttribute.Class = v
	}

	for _, v := range this.Block.Classes {
		err := v.resolveFather()
		if err != nil {
			this.errors = append(this.errors, err)
		}
		this.errors = append(this.errors, v.resolveInterfaces()...)
		this.errors = append(this.errors, v.resolveFieldsAndMethodsType()...)
	}

	for _, v := range this.Block.Classes {
		es := v.checkPhase1()
		this.errors = append(this.errors, es...)
		if this.shouldStop(nil) {
			return this.errors
		}
	}
	for _, v := range this.Block.Functions {
		if v.TemplateFunction != nil {
			continue
		}
		this.errors = append(this.errors, v.checkReturnVarExpression()...)
	}
	for _, v := range this.InitFunctions {
		this.errors = append(this.errors, v.check(&this.Block)...)
		if this.shouldStop(nil) {
			return this.errors
		}
	}
	for _, v := range this.Block.Classes {
		this.errors = append(this.errors, v.checkPhase2()...)
		if this.shouldStop(nil) {
			return this.errors
		}
	}
	for _, v := range this.Block.Functions {
		if v.IsBuildIn {
			continue
		}
		if v.TemplateFunction != nil {
			continue
		}
		v.checkBlock(&this.errors)
		if PackageBeenCompile.shouldStop(nil) {
			return this.errors
		}
	}
	for _, v := range this.statementLevelFunctions {
		v.IsClosureFunction = v.Closure.CaptureCount(v) > 0
	}
	for _, v := range this.statementLevelClass {
		for f, meta := range v.closure.Functions {
			if f.IsClosureFunction == false {
				continue
			}
			this.errors = append(this.errors,
				fmt.Errorf("%s trying to access capture function '%s' from outside",
					meta.pos.ErrMsgPrefix(), f.Name))
		}
	}
	if this.shouldStop(nil) {
		return this.errors
	}
	this.errors = append(this.errors, this.checkUnUsedPackage()...)
	return this.errors
}

/*
	load package or class
*/
func (this *Package) load(resource string) (interface{}, error) {
	if resource == "" {
		panic("null string")
	}
	if this.loadedClasses == nil {
		this.loadedClasses = make(map[string]*Class)
	}
	if t, ok := this.loadedClasses[resource]; ok {
		return t, nil
	}
	if this.LoadedPackages == nil {
		this.LoadedPackages = make(map[string]*Package)
	}
	if t, ok := this.LoadedPackages[resource]; ok {
		return t, nil
	}
	t, err := ImportsLoader.LoadImport(resource)
	if pp, ok := t.(*Package); ok && pp != nil {
		PackageBeenCompile.LoadedPackages[resource] = pp
		this.mkClassCache(pp)
		if pp.Name == common.CorePackage {
			pp.markBuildIn()
		}
	}
	if c, ok := t.(*Class); ok && c != nil {
		if c.IsJava == false {
			return nil, fmt.Errorf("load lucy class not allow")
		}
		PackageBeenCompile.loadedClasses[resource] = c
	}
	return t, err
}

func (this *Package) checkUnUsedPackage() []error {
	if compileCommon.CompileFlags.DisableCheckUnUse {
		return nil
	}
	errs := []error{}
	for _, v := range this.files {
		for _, i := range v.Imports {
			if i.Used == false {
				errs = append(errs, fmt.Errorf("%s '%s' imported not used",
					errMsgPrefix(i.Pos), i.Import))
			}
		}
	}
	for _, i := range this.unUsedPackage {
		pp, err := this.load(i.Import)
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

func (this *Package) loadClass(className string) (*Class, error) {
	if this.loadedClasses == nil {
		this.loadedClasses = make(map[string]*Class)
	}
	if c, ok := this.loadedClasses[className]; ok && c != nil {
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
	this.loadedClasses[className] = cc
	return cc, nil
}

func (this *Package) mkClassCache(loadedPackage *Package) {
	for _, v := range loadedPackage.Block.Classes {
		this.loadedClasses[v.Name] = v // binary name
	}
}

func (this *Package) insertImport(i *Import) error {
	if this.files == nil {
		this.files = make(map[string]*SourceFile)
	}
	x, ok := this.files[i.Pos.Filename]
	if ok == false {
		x = &SourceFile{}
		this.files[i.Pos.Filename] = x
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
