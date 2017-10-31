package ast

import (
	"fmt"
)

type Package struct {
	Files   map[string]*File
	Name    string //if error,should be multi names
	Inits   []*Block
	Funcs   map[string]*Function
	Classes map[string]*Class
	Enums   map[string]*Enum
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

type ConvertTops2Package struct {
	Name    []string //package name
	Inits   []*Block
	Funcs   map[string][]*Function
	Classes map[string][]*Class
	Enums   map[string][]*Enum
	Vars    map[string][]*GlobalVariable
	Consts  map[string][]*Const
	Import  []*Imports
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
	return ""
}

func (c *ConvertTops2Package) ConvertTops2Package(t Tops) (p *Package, errs []error) {
	p := &Package{}
	p.Files = make(map[string]*File)
	c.Name = []string{}
	c.Inits = []*Block{}
	c.Funcs = make(map[string][]*Function)
	c.Classes = make(map[string][]*Class)
	c.Enums = make(map[string][]*Enum)
	c.Vars = make(map[string][]*GlobalVariable)
	c.Consts = make(map[string][]*Const)
	//主要是检查重复申明
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			c.Inits = append(c.Inits, v.Data.(*Block))
		case *Function:
			t := v.Data.(*Function)
			if c.Funcs[t.Name] == nil {
				c.Funcs[t.Name] = []*Function{t}
			} else {
				c.Funcs[t.Name] = append(c.Funcs[t.Name], t)
			}
		case *Class:
			t := v.Data.(*Class)
			if c.Classes[t.Name] == nil {
				c.Classes[t.Name] = []*Class{t}
			} else {
				c.Classes[t.Name] = append(c.Classes[t.Name], t)
			}
		case *Enum:
			t := v.Data.(*Enum)
			if c.Enums[t.Name] == nil {
				c.Enums[t.Name] = []*Enum{t}
			} else {
				c.Enums[t.Name] = append(c.Enums[t.Name], t)
			}
		case *GlobalVariable:
			t := v.Data.(*GlobalVariable)
			if c.Enums[t.Name] == nil {
				c.Vars[t.Name] = []*GlobalVariable{t}
			} else {
				c.Vars[t.Name] = append(c.Vars[t.Name], t)
			}
		case *Const:
			t := v.Data.(*Const)
			if c.Enums[t.Name] == nil {
				c.Consts[t.Name] = []*GlobalVariable{t}
			} else {
				c.Consts[t.Name] = append(c.Vars[t.Name], t)
			}
		case *Imports:
			i := v.Data.(*Imports)
			if p.Files[i.Pos.Filename] == nil {
				p.Files[i.Pos.Filename] = &File{Imports: make(map[string]*Imports)}
			}
			p.Files[i.Pos.Filename][i.Name] = i
		case *PackageNameDeclare:
			t := v.Data.(*PackageNameDeclare)
			if p.Files[t.Pos.Filename] == nil {
				p.Files[t.Pos.Filename] = &File{Imports: make(map[string]*Imports)}
			}
			p.Files[t.Pos.Filename].Package = t
		default:
			panic("tops have unkown type")
		}
	}
	//package name no be the same
	{
		m := make(map[string][]*PackageNameDeclare)
		for _, v := range p.Files {
			if m[v.Package.Name] == nil {
				m[v.Package.Name] = []*PackageNameDeclare{v}
			}
		}
		if len(m) > 0 {
			t := []*PackageNameDeclare{}
			for _, v := range m {
				t = append(t, v)
			}
			errs = append(errs, &PackageNameNotConsistentError{t})
		}
	}
	//check redeclare error
	for name, v := range c.Funcs {
		if len(v) > 1 {
			t := []*Pos{}
			for _, vv := range v {
				t = append(t, &vv.Pos)
			}
			errs = append(errs, &RedeclareError{
				Name: name,
				Type: "function",
				Pos:  t,
			})
		}
	}
	//class redeclare
	for name, v := range c.Classes {
		if len(v) > 1 {
			t := []*Pos{}
			for _, vv := range v {
				t = append(t, &vv.Pos)
			}
			errs = append(errs, &RedeclareError{
				Name: name,
				Type: "class",
				Pos:  t,
			})
		}
	}
	for name, v := range c.Enums {
		if len(v) > 1 {
			t := []*Pos{}
			for _, vv := range v {
				t = append(t, &vv.Pos)
			}
			errs = append(errs, &RedeclareError{
				Name: name,
				Type: "enum",
				Pos:  t,
			})
		}
	}
	for name, v := range c.Vars {
		if len(v) > 1 {
			t := []*Pos{}
			for _, vv := range v {
				t = append(t, &vv.Pos)
			}
			errs = append(errs, &RedeclareError{
				Name: name,
				Type: "variable",
				Pos:  t,
			})
		}
	}
	for name, v := range c.Consts {
		if len(v) > 1 {
			t := []*Pos{}
			for _, vv := range v {
				t = append(t, &vv.Pos)
			}
			errs = append(errs, &RedeclareError{
				Name: name,
				Type: "const",
				Pos:  t,
			})
		}
	}
	return p, nil
}

func (p *Package) TypeCheck() []error {
	errs := []error{}
	errs = append(errs, p.checkConst()...)
	if len(errs) > 10 {
		return errs
	}

}

func (p *Package) checkConst() []error {
	for _, v := range p.Consts {
		v.Init.Typ
	}
}
