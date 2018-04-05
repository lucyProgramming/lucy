package ast

import (
	"errors"
	"fmt"
	"strings"
)

const (
	_ = iota
	PACKAGE_KIND_LUCY
	PACKAGE_KIND_JAVA
)

type Package struct {
	TriggerPackageInitMethodName string
	Kind                         int
	Name                         string
	Main                         *Function
	DestPath                     string
	LoadedPackages               map[string]*Package
	loaded                       map[string]*LoadedResouces
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

func (p *Package) addBuildFunctions() {
	if p.Block.Funcs == nil {
		p.Block.Funcs = make(map[string]*Function)
	}
	for k, f := range buildinFunctionsMap {
		ff := mkBuildinFunction(k, f.args, f.returns, f.checker)
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
		if p.Block.shouldStop(nil) {
			return p.Errors
		}
	}
	for _, v := range p.InitFunctions {
		p.Errors = append(p.Errors, v.check(&p.Block)...)
		if p.Block.shouldStop(nil) {
			return p.Errors
		}
	}
	for _, v := range p.Block.Classes {
		p.Errors = append(p.Errors, v.check(&p.Block)...)
	}
	for _, v := range p.Block.Funcs {
		if v.IsBuildin {
			continue
		}
		v.checkBlock(&p.Errors)
		if p.Block.shouldStop(nil) {
			return p.Errors
		}
	}
	return p.Errors
}

func (p *Package) load(resource string) (interface{}, error) {
	if p.loaded != nil {
		t, ok := p.loaded[resource]
		if ok {
			return t.T, t.Err
		}
	}
	if p.loaded == nil {
		p.loaded = make(map[string]*LoadedResouces)
	}
	p.loaded[resource] = &LoadedResouces{}
	var pp *Package

	pp, t, err := NameLoader.LoadName(resource)
	if pp != nil {
		if PackageBeenCompile.LoadedPackages == nil {
			PackageBeenCompile.LoadedPackages = make(map[string]*Package)
		}
		PackageBeenCompile.LoadedPackages[pp.Name] = pp
	}
	p.loaded[resource] = &LoadedResouces{}
	p.loaded[resource].T = t
	p.loaded[resource].Err = err
	if p.loaded[resource].T != nil {
		if tp, ok := p.loaded[resource].T.(*Package); ok && tp != nil && tp != pp {
			if PackageBeenCompile.LoadedPackages == nil {
				PackageBeenCompile.LoadedPackages = make(map[string]*Package)
			}
			PackageBeenCompile.LoadedPackages[tp.Name] = tp
		}
	}
	return p.loaded[resource].T, p.loaded[resource].Err
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
	Name string
	Pos  []*Pos
	Type string //varialbe or function
}

func (r *RedeclareError) Error() error {
	s := fmt.Sprintf("%s:%s redeclare")
	for _, v := range r.Pos {
		s += fmt.Sprintf("\t%s", errMsgPrefix(v))
	}
	return errors.New(s)
}
