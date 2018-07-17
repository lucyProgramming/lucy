package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"strings"
)

type Package struct {
	Name                         string
	LoadedPackages               map[string]*Package
	loadedClasses                map[string]*Class
	Block                        Block // package always have a default block
	Files                        map[string]*SourceFile
	InitFunctions                []*Function
	NErrors2Stop                 int // number of errors should stop compile
	Errors                       []error
	TriggerPackageInitMethodName string
}

func (p *Package) loadBuildInPackage() error {
	if p.Name == common.CorePackage {
		return nil
	}
	pp, err := p.load(common.CorePackage)
	if err != nil {
		return err
	}
	lucyBuildInPackage = pp.(*Package)
	p.Block.Outer = &lucyBuildInPackage.Block
	return nil
}

func (p *Package) getImport(file string, accessName string) *Import {
	if p.Files == nil {
		return nil
	}
	if _, ok := p.Files[file]; ok == false {
		return nil
	}
	return p.Files[file].Imports[accessName]
}

func (p *Package) mkInitFunctions(bs []*Block) {
	p.InitFunctions = make([]*Function, len(bs))
	for k, b := range bs {
		b.IsFunctionBlock = true
		f := &Function{}
		f.Pos = b.Pos
		f.Block = *b
		f.isGlobalVariableDefinition = b.IsGlobalVariableDefinitionBlock
		p.InitFunctions[k] = f
		f.Used = true
		f.isPackageBlockFunction = true
	}
}

func (p *Package) shouldStop(errs []error) bool {
	return (len(p.Errors) + len(errs)) >= p.NErrors2Stop
}

func (p *Package) TypeCheck() []error {
	if p.NErrors2Stop <= 2 {
		p.NErrors2Stop = 10
	}
	p.Errors = []error{}
	p.Errors = append(p.Errors, p.Block.checkConstants()...)
	//
	for _, v := range p.Block.Functions {
		if v.IsBuildIn {
			continue
		}
		v.Block.inherit(&p.Block)
		v.Block.InheritedAttribute.Function = v
		v.checkParametersAndReturns(&p.Errors)
		if p.shouldStop(nil) {
			return p.Errors
		}
	}
	for _, v := range p.Block.Enums {
		v.Name = p.Name + "/" + v.Name
		err := v.check()
		if err != nil {
			p.Errors = append(p.Errors, err)
		}
	}
	for _, v := range p.Block.Classes {
		v.Name = p.Name + "/" + v.Name
		v.mkDefaultConstruction()
		v.Block.inherit(&PackageBeenCompile.Block)
		v.Block.InheritedAttribute.Class = v
	}
	for _, v := range p.Block.Classes {
		err := v.resolveFather()
		if err != nil {
			p.Errors = append(p.Errors, err)
		}
		es := v.resolveInterfaces()
		if esNotEmpty(es) {
			p.Errors = append(p.Errors, es...)
		}
		es = v.resolveFieldsAndMethodsType()
		if esNotEmpty(es) {
			p.Errors = append(p.Errors, es...)
		}
	}
	for _, v := range p.Block.Classes {
		es := v.checkPhase1()
		if esNotEmpty(es) {
			p.Errors = append(p.Errors, es...)
		}
		if p.shouldStop(nil) {
			return p.Errors
		}
	}
	for _, v := range p.InitFunctions {
		p.Errors = append(p.Errors, v.check(&p.Block)...)
		if p.shouldStop(nil) {
			return p.Errors
		}
	}
	for _, v := range p.Block.Classes {
		es := v.checkPhase2()
		if esNotEmpty(es) {
			p.Errors = append(p.Errors, es...)
		}
		if p.shouldStop(nil) {
			return p.Errors
		}
	}
	for _, v := range p.Block.Functions {
		if v.IsBuildIn {
			continue
		}
		if v.TemplateFunction != nil {
			continue
		}
		v.checkBlock(&p.Errors)
		if PackageBeenCompile.shouldStop(nil) {
			return p.Errors
		}
	}
	return p.Errors
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
	if p.LoadedPackages == nil {
		p.LoadedPackages = make(map[string]*Package)
	}
	if t, ok := p.loadedClasses[resource]; ok {
		return t, nil
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

func (p *Package) mkClassCache(load *Package) {
	for _, v := range load.Block.Classes {
		p.loadedClasses[v.Name] = v // binary name
	}
}

//different from different source file
type SourceFile struct {
	Imports map[string]*Import // n
}

type Import struct {
	AccessName string
	Import     string // full name
	Pos        *Position
	Used       bool
}

/*
	import "github.com/lucy" should access by lucy.Println
	import "github.com/std" as std should access by std.Println
*/
func (i *Import) GetAccessName() (string, error) {
	if i.AccessName == "_" { //special case _ is a identifier
		return "", fmt.Errorf("'_' is not legal package access name")
	}
	if i.AccessName != "" {
		return i.AccessName, nil
	}
	name := i.Import
	if strings.Contains(i.Import, "/") {
		name = name[strings.LastIndex(name, "/")+1:]
		if name == "" {
			return "", fmt.Errorf("no last element after/")
		}
	}
	//check if legal
	if false == packageAccessNameReg.Match([]byte(name)) {
		return "", fmt.Errorf("%s is not legal package name", name)
	}
	i.AccessName = name
	return name, nil
}

type RedeclareError struct {
	Name      string
	Positions []*Position
	Types     []string
}

func (r *RedeclareError) Error() error {
	s := fmt.Sprintf("name '%s' defined  multi times,which are:\n", r.Name)
	for k, v := range r.Positions {
		s += fmt.Sprintf("\t%s '%s' named '%s'\n", errMsgPrefix(v), r.Types[k], r.Name)
	}
	return errors.New(s)
}
