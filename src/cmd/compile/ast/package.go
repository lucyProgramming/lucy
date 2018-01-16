package ast

import (
	"errors"
	"fmt"
	"strings"
)

type Package struct {
	loadedPackages map[string]*Package
	Block          Block // package always have a default block
	Files          map[string]*File
	Name           string //if error,should be multi names
	Blocks         []*Block
	NErros         int // number of errors should stop compile
}

func (p *Package) addBuildFunctions() {
	if p.Block.Funcs == nil {
		p.Block.Funcs = make(map[string]*Function)
	}
	{
		name := "print"
		f := mkBuildFunction(name, true, nil, nil)
		f.CallChecker = func(errs *[]error, args []*VariableType, pos *Pos) {}
		p.Block.Funcs[name] = f
	}
	{
		name := "panic"
		f := mkBuildFunction(name, true, nil, nil)
		f.CallChecker = oneAnyTypeParameterChecker
		p.Block.Funcs[name] = f
	}
	{
		name := "recover"
		f := mkBuildFunction(name, false, nil, nil)
		f.CallChecker = oneAnyTypeParameterChecker
		p.Block.Funcs[name] = f
	}
}
func (p *Package) TypeCheck() []error {
	p.addBuildFunctions()
	if p.NErros <= 2 {
		p.NErros = 10
	}
	errs := []error{}
	errs = append(errs, p.Block.checkConst()...)
	//
	for _, v := range p.Block.Funcs {
		if v.Isbuildin {
			continue
		}
		v.Block.inherite(&p.Block)
		errs = append(errs, v.check(&p.Block)...)
	}
	for _, v := range p.Block.Classes {
		errs = append(errs, v.check(&p.Block)...)
	}
	for _, v := range p.Blocks {
		errs = append(errs, v.check(&p.Block)...)
	}
	return errs
}

func (p *Package) loadPackage(name string) (*Package, error) {
	if p.loadedPackages == nil {
		p.loadedPackages = make(map[string]*Package)
	}
	if t := p.loadedPackages[name]; t != nil {
		return t, nil
	}
	pp, err := PackageLoad.LoadPackage(name)
	if err != nil {
		return nil, err
	}
	p.loadedPackages[name] = pp
	return pp, nil

}

//different for other file
type File struct {
	Imports map[string]*Imports // n
	Package *PackageNameDeclare
}

type Imports struct {
	AccessName string
	Name       string // full name
	Pos        *Pos
	Used       bool
}

/*
	import "github.com/lucy" should access by lucy.Println
	import "github.com/lucy" as std should access by std.Println
*/
func (i *Imports) GetAccessName() (string, error) {
	if i.AccessName == "_" { //special case _ is a identifer
		return "", fmt.Errorf("_ is not legal package name")
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

type PackageNameDeclare struct {
	Name string
	Pos  *Pos
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

type PackageNameNotConsistentError struct {
	Names []*PackageNameDeclare
}

func (p *PackageNameNotConsistentError) Error() string {
	if len(p.Names) == 0 {
		panic("zero length")
	}
	s := fmt.Sprintf("package named not consistently\n")
	for _, v := range p.Names {
		s += fmt.Sprintf("%s named by %s\n", errMsgPrefix(v.Pos), v.Name)
	}
	return s
}
