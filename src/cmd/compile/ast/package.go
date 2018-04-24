package ast

import (
	"errors"
	"fmt"
	"strings"
)

//const (
//	_ = iota
//	PACKAGE_KIND_LUCY
//	PACKAGE_KIND_JAVA
//)

type Package struct {
	//Kind                         int
	TriggerPackageInitMethodName string
	Name                         string
	Main                         *Function
	DestPath                     string
	LoadedPackages               map[string]*Package
	loadedClasses                map[string]*Class
	Block                        Block // package always have a default block
	Files                        map[string]*File
	InitFunctions                []*Function
	NErros2Stop                  int // number of errors should stop compile
	Errors                       []error
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

type LoadedResouces struct {
	T   interface{}
	Err error
}

func (p *Package) mkInitFunctions(bs []*Block) {
	p.InitFunctions = make([]*Function, len(bs))
	for k, b := range bs {
		f := &Function{}
		f.Block = b
		f.isGlobalVariableDefinition = b.isGlobalVariableDefinition
		f.Typ = &FunctionType{}
		p.InitFunctions[k] = f
		f.Used = true
		f.isPackageBlockFunction = true
	}
}

func (p *Package) shouldStop(errs []error) bool {
	return (len(p.Errors) + len(errs)) >= p.NErros2Stop
}

func (p *Package) addBuildFunctions() {
	if p.Block.Funcs == nil {
		p.Block.Funcs = make(map[string]*Function)
	}
	for k, f := range buildinFunctionsMap {
		ff := mkBuildinFunction(k, f.args, f.returnList, f.checker)
		p.Block.Funcs[k] = ff
	}
}

func (p *Package) TypeCheck() []error {
	p.addBuildFunctions()
	if p.NErros2Stop <= 2 {
		p.NErros2Stop = 10
	}
	p.Errors = []error{}
	p.Errors = append(p.Errors, p.Block.checkConst()...)
	//
	for _, v := range p.Block.Funcs {
		if v.IsBuildin {
			continue
		}
		v.Block.inherite(&p.Block)
		v.Block.InheritedAttribute.Function = v
		v.checkParaMeterAndRetuns(&p.Errors)
		if p.shouldStop(nil) {
			return p.Errors
		}
	}
	for _, v := range p.Block.Enums {
		v.Name = p.Name + "/" + v.Name
	}
	for _, v := range p.Block.Classes {
		v.Name = p.Name + "/" + v.Name
	}

	for _, v := range p.Block.Classes {
		es := v.checkPhase1(&p.Block)
		if errsNotEmpty(es) {
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
		es := v.checkPhase2(&p.Block)
		if errsNotEmpty(es) {
			p.Errors = append(p.Errors, es...)
		}
		if p.shouldStop(nil) {
			return p.Errors
		}
	}

	for _, v := range p.Block.Funcs {
		if v.IsBuildin {
			continue
		}
		v.checkBlock(&p.Errors)
		if PackageBeenCompile.shouldStop(nil) {
			return p.Errors
		}
	}
	return p.Errors
}

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
	var pp *Package
	pp, t, err := NameLoader.LoadName(resource)
	if pp != nil {
		PackageBeenCompile.LoadedPackages[resource] = pp
		p.mkClassCache(pp)
	}
	if pp, ok := t.(*Package); ok && pp != nil {
		PackageBeenCompile.LoadedPackages[resource] = pp
		p.mkClassCache(pp)
	}
	if c, ok := t.(*Class); ok && c != nil {
		PackageBeenCompile.loadedClasses[resource] = c
	}
	return t, err
}
func (p *Package) mkClassCache(load *Package) {
	for _, v := range load.Block.Classes {
		p.loadedClasses[v.Name] = v // binary name
	}
}

//different for other file
type File struct {
	Imports map[string]*Import // n
}

type Import struct {
	AccessName string
	Resource   string // full name
	Pos        *Pos
	Used       bool
}

/*
	import "github.com/lucy" should access by lucy.Println
	import "github.com/std" as std should access by std.Println
*/
func (i *Import) GetAccessName() (string, error) {
	if i.AccessName == "_" { //special case _ is a identifer
		return "", fmt.Errorf("'_' is not legal package name")
	}
	if i.AccessName != "" {
		return i.AccessName, nil
	}
	name := i.Resource
	if strings.Contains(i.Resource, "/") {
		name = name[strings.LastIndex(name, "/")+1:]
		if name == "" {
			return "", fmt.Errorf("no last element after/")
		}
	}
	//check if legal
	if !packageAliasReg.Match([]byte(name)) {
		return "", fmt.Errorf("%s is not legal package name", name)
	}
	i.AccessName = name
	return name, nil
}

type RedeclareError struct {
	Name  string
	Poses []*Pos
	Types []string //varialbe or function
}

func (r *RedeclareError) Error() error {
	var firstPos *Pos
	for _, v := range r.Poses {
		if firstPos == nil {
			firstPos = v
			continue
		}
		if v.StartLine < firstPos.StartLine {
			firstPos = v
		}
	}
	s := fmt.Sprintf("%s name named '%s' defined multi times,which are:\n",
		errMsgPrefix(firstPos), r.Name)
	for k, v := range r.Poses {
		if v == firstPos {
			continue
		}
		s += fmt.Sprintf("\t%s '%s' named '%s'\n", errMsgPrefix(v), r.Types[k], r.Name)
	}
	return errors.New(s)
}
