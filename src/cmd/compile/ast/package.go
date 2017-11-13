package ast

import (
	"fmt"
)

type Package struct {
	Files     map[string]*File
	Name      string //if error,should be multi names
	Blocks    []*Block
	Funcs     map[string]*Function
	Classes   map[string]*Class
	Enums     map[string]*Enum
	EnumNames map[string]*EnumName
	Vars      map[string]*VariableDefinition
	Consts    map[string]*Const
	NErros    int // number of errors should stop compile
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
