package ast

import (
	"fmt"
	"os"
)

//代表语法数的一个节点
type Node struct {
	Data interface{}
}

type ConvertTops2Package struct {
	Name      []string //package name
	Blocks    []*Block
	Funcs     []*Function
	Classes   []*Class
	Enums     []*Enum
	Vars      []*VariableDefinition
	Consts    []*Const
	Import    []*Import
	TypeAlias []*ExpressionTypeAlias
}

func (conversion *ConvertTops2Package) ConvertTops2Package(t []*Node) (redeclareErrors []*RedeclareError, errs []error) {
	//
	if err := PackageBeenCompile.loadBuildinPackage(); err != nil {
		fmt.Printf("load lucy buildin package failed,err:%v\n", err)
		os.Exit(1)
	}
	errs = make([]error, 0)
	PackageBeenCompile.Files = make(map[string]*File)
	conversion.Name = []string{}
	conversion.Blocks = []*Block{}
	conversion.Funcs = make([]*Function, 0)
	conversion.Classes = make([]*Class, 0)
	conversion.Enums = make([]*Enum, 0)
	conversion.Vars = make([]*VariableDefinition, 0)
	conversion.Consts = make([]*Const, 0)
	expressions := []*Expression{}
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			t := v.Data.(*Block)
			conversion.Blocks = append(conversion.Blocks, t)
		case *Function:
			t := v.Data.(*Function)
			conversion.Funcs = append(conversion.Funcs, t)
		case *Enum:
			t := v.Data.(*Enum)
			conversion.Enums = append(conversion.Enums, t)
		case *Class:
			t := v.Data.(*Class)
			conversion.Classes = append(conversion.Classes, t)
		case *Const:
			t := v.Data.(*Const)
			conversion.Consts = append(conversion.Consts, t)
		case *Import:
			i := v.Data.(*Import)
			if PackageBeenCompile.Files[i.Pos.Filename] == nil {
				PackageBeenCompile.Files[i.Pos.Filename] = &File{Imports: make(map[string]*Import)}
			}
			PackageBeenCompile.Files[i.Pos.Filename].Imports[i.AccessName] = i
		case *Expression: // a,b = f();
			t := v.Data.(*Expression)
			expressions = append(expressions, t)
		case *ExpressionTypeAlias:
			t := v.Data.(*ExpressionTypeAlias)
			conversion.TypeAlias = append(conversion.TypeAlias, t)
		default:
			panic("tops have unkown type")
		}
	}
	errs = append(errs, checkEnum(conversion.Enums)...)
	redeclareErrors = conversion.redeclareErrors()
	PackageBeenCompile.Block.Consts = make(map[string]*Const)
	for _, v := range conversion.Consts {
		PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
	}
	PackageBeenCompile.Block.Vars = make(map[string]*VariableDefinition)
	PackageBeenCompile.Block.Funcs = make(map[string]*Function)
	for _, v := range conversion.Funcs {
		v.IsGlobal = true
		err := PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.Classes = make(map[string]*Class)
	for _, v := range conversion.Classes {
		err := PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.Enums = make(map[string]*Enum)
	PackageBeenCompile.Block.EnumNames = make(map[string]*EnumName)
	for _, v := range conversion.Enums {
		err := PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	//after class inserted,then resolve type
	for _, v := range conversion.TypeAlias {
		err := v.Typ.resolve(&PackageBeenCompile.Block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if PackageBeenCompile.Block.Types == nil {
			PackageBeenCompile.Block.Types = make(map[string]*VariableType)
		}
		v.Typ.Alias = v.Name
		PackageBeenCompile.Block.Types[v.Name] = v.Typ
	}
	if len(expressions) > 0 {
		s := make([]*Statement, len(expressions))
		for k, v := range expressions {
			s[k] = &Statement{
				Typ:        STATEMENT_TYPE_EXPRESSION,
				Expression: v,
			}
		}
		b := &Block{}
		b.Statements = s
		b.isGlobalVariableDefinition = true
		conversion.Blocks = append([]*Block{b}, conversion.Blocks...)
	}
	PackageBeenCompile.mkInitFunctions(conversion.Blocks)
	return
}

func (conversion *ConvertTops2Package) redeclareErrors() []*RedeclareError {
	ret := []*RedeclareError{}
	m := make(map[string][]interface{})
	//eums
	for _, v := range conversion.Enums {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
		for _, vv := range v.Enums {
			if _, ok := m[vv.Name]; ok {
				m[vv.Name] = append(m[vv.Name], vv)
			} else {
				m[vv.Name] = []interface{}{vv}
			}
		}
	}
	//const
	for _, v := range conversion.Consts {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//vars
	for _, v := range conversion.Vars {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//funcs
	for _, v := range conversion.Funcs {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//classes
	for _, v := range conversion.Classes {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	// type alias
	for _, v := range conversion.TypeAlias {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}

	for k, v := range m {
		if len(v) == 1 || len(v) == 0 { //very good  , 0 looks is impossible
			continue
		}
		r := &RedeclareError{}
		r.Name = k
		r.Positions = make([]*Pos, len(v))
		r.Types = make([]string, len(v))
		for kk, vv := range v {
			switch vv.(type) {
			case *Const:
				t := vv.(*Const)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "const"
			case *Enum:
				t := vv.(*Enum)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "enum"
			case *Function:
				t := vv.(*Function)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "function"
			case *Class:
				t := vv.(*Class)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "class"
			case *ExpressionTypeAlias:
				t := vv.(*ExpressionTypeAlias)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "type alias"
			case *EnumName:
				t := vv.(*EnumName)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "enum name"
			default:
				panic("make error")
			}
		}
		ret = append(ret, r)
	}
	return ret
}
