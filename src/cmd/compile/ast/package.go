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
	Enums   []*Enum
	Vars    map[string]*GlobalVariable
	Consts  map[string]*Const
}

func (p *ConvertTops2Package) redeclareErrors() []*RedeclareError {
	ret := []*RedeclareError{}
	m := make(map[string][]interface{})
	//eums
	for _, v := range p.Enums {
		for _, vv := range v.Names {
			if _, ok := m[vv.Name]; ok {
				m[vv.Name] = append(m[vv.Name], vv)
			} else {
				m[vv.Name] = []interface{}{vv}
			}
		}
	}
	//const
	for _, v := range p.Consts {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}

	//vars
	for _, v := range p.Vars {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}

	//funcs
	for _, v := range p.Funcs {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//classes
	for _, v := range p.Classes {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	for k, v := range m {
		if len(v) == 1 { //very good
			continue
		}
		r := &RedeclareError{}
		r.Name = k
		r.Pos = make([]*Pos, 0)
		for _, vv := range v {
			switch vv.(type) {
			case *Const:
				t := vv.(*Const)
				r.Pos = append(r.Pos, &t.Pos)
				r.Type = "const"
			case *Enum:
				t := vv.(*EnumNames)
				r.Pos = append(r.Pos, &t.Pos)
				r.Type = "enum"
			case *GlobalVariable:
				t := vv.(*GlobalVariable)
				r.Name = t.Name
				r.Pos = append(r.Pos, &t.Pos)
				r.Type = "global varialbe"
			case *Function:
				t := vv.(*Function)
				r.Pos = append(r.Pos, &t.Pos)
				r.Type = "function"
			case *Class:
				t := vv.(*Class)
				r.Pos = append(r.Pos, &t.Pos)
				r.Type = "class"
			default:
				panic("make error")
			}
		}
		ret = append(ret, r)
	}
	return ret
}

func (p *ConvertTops2Package) checkEnum() []error {
	ret := make([]error, 0)
	for _, v := range p.Enums {
		if len(v.Names) == 0 {
			continue
		}
		is, typ, value, err := v.Init.getConstValue()
		if err != nil || is == false || typ != EXPRESSION_TYPE_INT {
			ret = append(ret, fmt.Errorf("enum type must inited by integer"))
			continue
		}
		v.Value = value.(int64)
		for k, vv := range v.Names {
			vv.Value = int64(k) + v.Value
		}
	}
	return ret
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
	Funcs   []*Function
	Classes []*Class
	Enums   []*Enum
	Vars    []*GlobalVariable
	Consts  []*Const
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
	if len(p.Names) == 0 {
		panic("zero length")
	}
	s := fmt.Sprintf("package named not consistently")
	for _, v := range p.Names {
		s += fmt.Sprintf("\tnamed by %s %s %d:%d\n", v.Name, v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn)
	}
	return s
}

func (c *ConvertTops2Package) ConvertTops2Package(t []*Node) (p *Package, redeclareErrors []*RedeclareError, errs []error) {
	errs = make([]error, 0)
	p = &Package{}
	p.Files = make(map[string]*File)
	c.Name = []string{}
	c.Inits = []*Block{}
	c.Funcs = make([]*Function, 0)
	c.Classes = make([]*Class, 0)
	c.Enums = make([]*Enum, 0)
	c.Vars = make([]*GlobalVariable, 0)
	c.Consts = make([]*Const, 0)
	//主要是检查重复申明
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			c.Inits = append(c.Inits, v.Data.(*Block))
		case *Function:
			t := v.Data.(*Function)
			c.Funcs = append(c.Funcs, t)
		case *Enum:
			t := v.Data.(*Enum)
			c.Enums = append(c.Enums, t)
		case *Class:
			t := v.Data.(*Class)
			c.Classes = append(c.Classes, t)
		case *GlobalVariable:
			t := v.Data.(*GlobalVariable)
			c.Vars = append(c.Vars, t)
		case *Const:
			t := v.Data.(*Const)
			c.Consts = append(c.Consts, t)
		case *Imports:
			i := v.Data.(*Imports)
			if p.Files[i.Pos.Filename] == nil {
				p.Files[i.Pos.Filename] = &File{Imports: make(map[string]*Imports)}
			}

			if p.Files[i.Pos.Filename] == nil {
				p.Files[i.Pos.Filename] = &File{}
			}
			if p.Files[i.Pos.Filename].Imports == nil {
				p.Files[i.Pos.Filename].Imports = make(map[string]*Imports)
			}
			p.Files[i.Pos.Filename].Imports[i.Name] = i
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
				m[v.Package.Name] = []*PackageNameDeclare{v.Package}
			}
		}
		if len(m) > 0 {
			t := []*PackageNameDeclare{}
			for _, v := range m {
				t = append(t, v...)
			}
			errs = append(errs, &PackageNameNotConsistentError{t})
		}
	}
	errs = append(errs, c.checkEnum()...)
	redeclareErrors = c.redeclareErrors()
	return
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
	for _, v := range p.Consts {
		v.Init.getConstValue()
	}
	return nil
}
