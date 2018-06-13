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
	Functions []*Function
	Classes   []*Class
	Enums     []*Enum
	Variables []*VariableDefinition
	Constants []*Constant
	Import    []*Import
	TypeAlias []*ExpressionTypeAlias
}

func (conversion *ConvertTops2Package) ConvertTops2Package(t []*Node) (redeclareErrors []*RedeclareError, errs []error) {
	//
	if err := PackageBeenCompile.loadBuildInPackage(); err != nil {
		fmt.Printf("load lucy buildin package failed,err:%v\n", err)
		os.Exit(1)
	}
	errs = make([]error, 0)
	PackageBeenCompile.Files = make(map[string]*LucyFile)
	conversion.Name = []string{}
	conversion.Blocks = []*Block{}
	conversion.Functions = make([]*Function, 0)
	conversion.Classes = make([]*Class, 0)
	conversion.Enums = make([]*Enum, 0)
	conversion.Variables = make([]*VariableDefinition, 0)
	conversion.Constants = make([]*Constant, 0)
	expressions := []*Expression{}
	for _, v := range t {
		switch v.Data.(type) {
		case *Block:
			t := v.Data.(*Block)
			conversion.Blocks = append(conversion.Blocks, t)
		case *Function:
			t := v.Data.(*Function)
			conversion.Functions = append(conversion.Functions, t)
		case *Enum:
			t := v.Data.(*Enum)
			conversion.Enums = append(conversion.Enums, t)
		case *Class:
			t := v.Data.(*Class)
			conversion.Classes = append(conversion.Classes, t)
		case *Constant:
			t := v.Data.(*Constant)
			conversion.Constants = append(conversion.Constants, t)
		case *Import:
			i := v.Data.(*Import)
			if PackageBeenCompile.Files[i.Pos.Filename] == nil {
				PackageBeenCompile.Files[i.Pos.Filename] = &LucyFile{Imports: make(map[string]*Import)}
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
	PackageBeenCompile.Block.Constants = make(map[string]*Constant)
	for _, v := range conversion.Constants {
		PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
	}
	PackageBeenCompile.Block.Variables = make(map[string]*VariableDefinition)
	PackageBeenCompile.Block.Functions = make(map[string]*Function)
	for _, v := range conversion.Functions {
		v.IsGlobal = true
		err := PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.Classes = make(map[string]*Class)
	for _, v := range conversion.Classes {
		err := PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.Enums = make(map[string]*Enum)
	PackageBeenCompile.Block.EnumNames = make(map[string]*EnumName)
	for _, v := range conversion.Enums {
		err := PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	//after class inserted,then resolve type
	for _, v := range conversion.TypeAlias {
		err := v.Type.resolve(&PackageBeenCompile.Block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if PackageBeenCompile.Block.TypeAlias == nil {
			PackageBeenCompile.Block.TypeAlias = make(map[string]*VariableType)
		}
		v.Type.Alias = v.Name
		PackageBeenCompile.Block.TypeAlias[v.Name] = v.Type
	}
	if len(expressions) > 0 {
		s := make([]*Statement, len(expressions))
		for k, v := range expressions {
			s[k] = &Statement{
				Type:       STATEMENT_TYPE_EXPRESSION,
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
	for _, v := range conversion.Constants {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//vars
	for _, v := range conversion.Variables {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//funcs
	for _, v := range conversion.Functions {
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
			case *Constant:
				t := vv.(*Constant)
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
