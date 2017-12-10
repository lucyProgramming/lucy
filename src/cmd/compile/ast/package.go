package ast

import (
	"fmt"
	"strings"
)

type Package struct {
	Block     Block // package always have a default block
	Files     map[string]*File
	Name      string //if error,should be multi names
	Blocks    []*Block
	Funcs     map[string][]*Function
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
	Alias string
	Name  string // full name
	Pos   *Pos
	Used  bool
}

/*
	import "github.com/lucy" should access by lucy.Println
	import "github.com/lucy" as std should access by std.Println
*/
func (i *Imports) AccessName() (string, error) {
	if i.Alias == "_" { //special case _ is a identifer
		return "", fmt.Errorf("_ is not legal package name")
	}
	if i.Alias != "" {
		return i.Alias, nil
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
	i.Alias = name
	return name, nil
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
