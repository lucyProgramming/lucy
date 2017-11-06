package ast

import (
	"errors"
	"fmt"
)

type Package struct {
	Files   map[string]*File
	Name    string //if error,should be multi names
	Inits   []*Block
	Funcs   map[string]*Function
	Classes map[string]*Class
	Enums   []*Enum
	Vars    map[string]*GlobalVariable
	Consts  map[string]*Const
}

//different for other file
type File struct {
	Imports map[string]*Imports // n
	Package *PackageNameDeclare
}
type Imports struct {
	Name  string
	Alias string
	Pos   Pos
}

type PackageNameDeclare struct {
	Name string
	Pos  Pos
}

type RedeclareError struct {
	Name string
	Pos  []*Pos
	Type string //varialbe or function
}

func (r *RedeclareError) Error() string {
	s := fmt.Sprintf("%s:%s redeclare")
	for _, v := range r.Pos {
		s += fmt.Sprintf("\t%s %d:%d\n", v.Filename, v.StartLine, v.StartColumn)
	}
	return s
}

type PackageNameNotConsistentError struct {
	Names []*PackageNameDeclare
}

func (p *PackageNameNotConsistentError) Error() string {
	if len(p.Names) == 0 {
		panic("zero length")
	}
	s := fmt.Sprintf("package named not consistently")
	for _, v := range p.Names {
		s += fmt.Sprintf("\tnamed by %s %s %d:%d\n", v.Name, v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn)
	}
	return s
}

func (p *Package) TypeCheck() []error {
	//name conflict,such as function name and class names
	errs := []error{}
	errs = append(errs, p.checkConst()...)
	if len(errs) > 10 {
		return errs
	}
	return errs
}

func (p *Package) checkConst() []error {
	errs := make([]error, 0)
	for _, v := range p.Consts {
		is, t, value, err := v.Init.getConstValue()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if is == false {
			errs = append(errs, errors.New("is not a const value"))
			continue
		}
		//rewrite
		v.Init = &Expression{}
		v.Init.Typ = t
		v.Init.Data = value
	}
	return errs
}
func (p *Package) checkGlobalVariables() []error {
	errs := make([]error, 0)
	for _, v := range p.Vars {

	}
	return errs
}
