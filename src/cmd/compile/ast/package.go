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
	Kind           int
	Name           string //if error,should be multi names ,taken first is ok
	FullName       string
	Main           *Function
	DestPath       string
	loadedPackages map[string]*Package
	loaded         map[string]*LoadedName
	Block          Block // package always have a default block
	Files          map[string]*File
	InitFunctions  []*Function
	NErros2Stop    int // number of errors should stop compile
	Errors         []error
}

type LoadedName struct {
	T   interface{}
	Err error
}

func (p *Package) mkShortName() {
	if strings.Contains(p.FullName, "/") {
		t := strings.Split(p.FullName, "/")
		p.Name = t[len(t)-1]
		if p.Name == "" {
			panic("last element is null string")
		}
	} else {
		p.Name = p.FullName
	}
}

func (p *Package) mkInitFunctions(bs []*Block) {
	p.InitFunctions = make([]*Function, len(bs))
	for k, b := range bs {
		f := &Function{}
		f.Block = b
		f.isGlobalVariableDefinition = b.isGlobalVariableDefinition
		f.Typ = &FunctionType{}
		f.MkVariableType()
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
	p.mkShortName()
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

func (p *Package) load(pname, name string) (interface{}, error) {
	if p.loadedPackages == nil {
		p.loadedPackages = make(map[string]*Package)
	}
	if p.loaded == nil {
		p.loaded = make(map[string]*LoadedName)
	}
	var err error
	fullname := pname + "/" + name
	if t, ok := p.loaded[fullname]; ok { // look up in cache
		return t.T, t.Err
	}

	if t := p.loadedPackages[pname]; t != nil {
		if t.Kind == PACKAGE_KIND_LUCY {
			tt := t.Block.SearchByName(name)
			if tt == nil {
				err = fmt.Errorf("%s is not found", name)
			}
			return tt, nil
		} else { //java package

		}
	}
	if _, ok := p.loadedPackages[pname]; ok == false {
		p.loadedPackages[pname] = &Package{}
	}
	p.loaded[fullname] = &LoadedName{}

	t, err := NameLoader.LoadName(p.loadedPackages[pname], pname, name)
	if err != nil {
		p.loaded[fullname].Err = err
		return nil, err
	}
	p.loaded[fullname].T = t
	return t, nil

}

//different for other file
type File struct {
	Imports map[string]*Imports // n
}

type Imports struct {
	AccessName string
	Name       string // full name
	Pos        *Pos
	Used       bool
}

/*
	import "github.com/lucy" should access by lucy.Println
	import "github.com/std" as std should access by std.Println
*/
func (i *Imports) GetAccessName() (string, error) {
	if i.AccessName == "_" { //special case _ is a identifer
		return "", fmt.Errorf("'_' is not legal package name")
	}
	if i.AccessName != "" {
		return i.AccessName, nil
	}
	name := i.Name
	if strings.Contains(i.Name, "/") {
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
