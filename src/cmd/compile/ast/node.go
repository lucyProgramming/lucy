package ast

import (
//"fmt"
)

//import "fmt"

//代表语法数的一个节点
type Node struct {
	//Pos  Pos
	Data interface{} //class defination or varialbe Defination
}

//type Tops []*Node //语法树顶层结构

type ConvertTops2Package struct {
	Name    []string //package name
	Blocks  []*Block
	Funcs   []*Function
	Classes []*Class
	Enums   []*Enum
	Vars    []*VariableDefinition
	Consts  []*Const
	Import  []*Imports
}

func (p *ConvertTops2Package) redeclareErrors() []*RedeclareError {
	ret := []*RedeclareError{}
	m := make(map[string][]interface{})
	//eums
	for _, v := range p.Enums {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
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
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "const"
			case *Enum:
				t := vv.(*EnumName)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "enum"
			case *VariableDefinition:
				t := vv.(*VariableDefinition)
				r.Name = t.Name
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "global varialbe"
			case *Function:
				t := vv.(*Function)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "function"
			case *Class:
				t := vv.(*Class)
				r.Pos = append(r.Pos, t.Pos)
				r.Type = "class"
			default:
				panic("make error")
			}
		}
		ret = append(ret, r)
	}
	return ret
}

func (c *ConvertTops2Package) ConvertTops2Package(t []*Node) (p *Package, redeclareErrors []*RedeclareError, errs []error) {
	errs = make([]error, 0)
	p = &Package{}
	p.Files = make(map[string]*File)
	c.Name = []string{}
	c.Blocks = []*Block{}
	c.Funcs = make([]*Function, 0)
	c.Classes = make([]*Class, 0)
	c.Enums = make([]*Enum, 0)
	c.Vars = make([]*VariableDefinition, 0)
	c.Consts = make([]*Const, 0)
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			t := v.Data.(*Block)
			t.p = p
			c.Blocks = append(c.Blocks, t)
		case *Function:
			t := v.Data.(*Function)
			c.Funcs = append(c.Funcs, t)
		case *Enum:
			t := v.Data.(*Enum)
			c.Enums = append(c.Enums, t)
		case *Class:
			t := v.Data.(*Class)
			c.Classes = append(c.Classes, t)
		case *VariableDefinition:
			t := v.Data.(*VariableDefinition)
			c.Vars = append(c.Vars, t)
		case *Const:
			t := v.Data.(*Const)
			c.Consts = append(c.Consts, t)
		case *Imports:
			i := v.Data.(*Imports)
			if p.Files[i.Pos.Filename] == nil {
				p.Files[i.Pos.Filename] = &File{Imports: make(map[string]*Imports)}
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
	//package name not be the same one
	{
		m := make(map[string][]*PackageNameDeclare)
		for _, v := range p.Files {
			if m[v.Package.Name] == nil {
				m[v.Package.Name] = []*PackageNameDeclare{v.Package}
			}
		}
		if len(m) > 1 {
			t := []*PackageNameDeclare{}
			for _, v := range m {
				t = append(t, v...)
			}
			errs = append(errs, &PackageNameNotConsistentError{t})
		}
	}
	errs = append(errs, checkEnum(c.Enums)...)
	redeclareErrors = c.redeclareErrors()
	p.Block.Consts = make(map[string]*Const)
	for _, v := range c.Consts {
		p.Block.insert(v.Name, v.Pos, v)

	}
	p.Block.Vars = make(map[string]*VariableDefinition)
	for _, v := range c.Vars {
		p.Block.Vars[v.Name] = v
	}
	p.Block.Funcs = make(map[string]*Function)
	for _, v := range c.Funcs {
		v.mkVariableType()
		p.Block.insert(v.Name, v.Pos, v)
	}
	p.Block.Classes = make(map[string]*Class)
	for _, v := range c.Classes {
		v.mkVariableType()
		p.Block.Classes[v.Name] = v
	}
	p.Block.Enums = make(map[string]*Enum)
	p.Block.EnumNames = make(map[string]*EnumName)
	for _, v := range c.Enums {
		v.mkVariableType()
		p.Block.Enums[v.Name] = v
		for _, vv := range v.Names {
			p.Block.EnumNames[vv.Name] = vv
		}
	}
	p.Blocks = c.Blocks
	return
}
