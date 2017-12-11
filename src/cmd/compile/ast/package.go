package ast

import (
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

func (f *File) fullPackageName(accessname, name string) (fullname string, accessResouce string, err error) {
	if f.Imports[accessname] == nil {
		err = fmt.Errorf("package %s not imported", accessname)
		return
	}
	t := strings.Split(name, ".")
	fullname = f.Imports[accessname].Name
	for i := 0; i < len(t)-1; i++ {
		fullname += "/" + t[i]
	}
	accessResouce = t[len(t)-1] // last element
	return
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
