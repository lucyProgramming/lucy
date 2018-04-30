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

func (convertor *ConvertTops2Package) ConvertTops2Package(t []*Node) (redeclareErrors []*RedeclareError, errs []error) {
	//
	if err := PackageBeenCompile.loadBuildinPackage(); err != nil {
		fmt.Printf("load lucy buildin package failed,err:%v\n", err)
		os.Exit(1)
	}

	errs = make([]error, 0)
	PackageBeenCompile.Files = make(map[string]*File)
	convertor.Name = []string{}
	convertor.Blocks = []*Block{}
	convertor.Funcs = make([]*Function, 0)
	convertor.Classes = make([]*Class, 0)
	convertor.Enums = make([]*Enum, 0)
	convertor.Vars = make([]*VariableDefinition, 0)
	convertor.Consts = make([]*Const, 0)
	expressions := []*Expression{}
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			t := v.Data.(*Block)
			convertor.Blocks = append(convertor.Blocks, t)
		case *Function:
			t := v.Data.(*Function)
			convertor.Funcs = append(convertor.Funcs, t)
		case *Enum:
			t := v.Data.(*Enum)
			convertor.Enums = append(convertor.Enums, t)
		case *Class:
			t := v.Data.(*Class)
			convertor.Classes = append(convertor.Classes, t)
		case *Const:
			t := v.Data.(*Const)
			convertor.Consts = append(convertor.Consts, t)
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
			convertor.TypeAlias = append(convertor.TypeAlias, t)
		default:
			panic("tops have unkown type")
		}
	}
	errs = append(errs, checkEnum(convertor.Enums)...)
	redeclareErrors = convertor.redeclareErrors()
	PackageBeenCompile.Block.Consts = make(map[string]*Const)
	for _, v := range convertor.Consts {
		PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
	}
	PackageBeenCompile.Block.Vars = make(map[string]*VariableDefinition)
	PackageBeenCompile.Block.Funcs = make(map[string]*Function)
	for _, v := range convertor.Funcs {
		err := PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
		v.IsGlobal = true
	}
	PackageBeenCompile.Block.Classes = make(map[string]*Class)
	for _, v := range convertor.Classes {
		err := PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.Enums = make(map[string]*Enum)
	PackageBeenCompile.Block.EnumNames = make(map[string]*EnumName)
	for _, v := range convertor.Enums {
		err := PackageBeenCompile.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	//after class inserted,then resolve type
	for _, v := range convertor.TypeAlias {
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
		convertor.Blocks = append([]*Block{b}, convertor.Blocks...)
	}
	PackageBeenCompile.mkInitFunctions(convertor.Blocks)
	return
}

func (convertor *ConvertTops2Package) redeclareErrors() []*RedeclareError {
	ret := []*RedeclareError{}
	m := make(map[string][]interface{})
	//eums
	for _, v := range convertor.Enums {
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
	for _, v := range convertor.Consts {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//vars
	for _, v := range convertor.Vars {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//funcs
	for _, v := range convertor.Funcs {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//classes
	for _, v := range convertor.Classes {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	// type alias
	for _, v := range convertor.TypeAlias {
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
		r.Poses = make([]*Pos, len(v))
		r.Types = make([]string, len(v))
		for kk, vv := range v {
			switch vv.(type) {
			case *Const:
				t := vv.(*Const)
				r.Poses[kk] = t.Pos
				r.Types[kk] = "const"
			case *Enum:
				t := vv.(*EnumName)
				r.Poses[kk] = t.Pos
				r.Types[kk] = "enum"
			case *Function:
				t := vv.(*Function)
				r.Poses[kk] = t.Pos
				r.Types[kk] = "function"
			case *Class:
				t := vv.(*Class)
				r.Poses[kk] = t.Pos
				r.Types[kk] = "class"
			case *ExpressionTypeAlias:
				t := vv.(*ExpressionTypeAlias)
				r.Poses[kk] = t.Pos
				r.Types[kk] = "type alias"
			default:
				panic("make error")
			}
		}
		ret = append(ret, r)
	}
	return ret
}
