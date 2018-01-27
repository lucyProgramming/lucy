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
	Name           string //if error,should be multi names ,taken first is ok
	InitFunctions  []*Function
	NErros2Stop    int // number of errors should stop compile
	Errors         []error
}

func (p *Package) mkInitFunctions(bs ...*Block) {
	p.InitFunctions = make([]*Function, len(bs))
	for k, b := range bs {
		f := &Function{}
		f.Name = fmt.Sprintf("<init%d>", k)
		f.Block = b
		f.Typ = &FunctionType{}
		f.mkVariableType()
		p.InitFunctions[k] = f
		f.Used = true
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
		if v.Isbuildin {
			continue
		}
		v.Block.inherite(&p.Block)
		p.Errors = append(p.Errors, v.check(&p.Block)...)
	}
	for _, v := range p.Block.Classes {
		p.Errors = append(p.Errors, v.check(&p.Block)...)
	}
	for _, v := range p.InitFunctions {
		p.Errors = append(p.Errors, v.check(&p.Block)...)
	}
	return p.Errors
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
	import "github.com/std" as std should access by std.Println
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
