package ast

import (
	"fmt"
	"os"
)

type TopNode struct {
	Node interface{}
}

type ConvertTops2Package struct {
	Blocks    []*Block
	Functions []*Function
	Classes   []*Class
	Enums     []*Enum
	Constants []*Constant
	Imports   []*Import
	TypeAlias []*TypeAlias
}

func (this *ConvertTops2Package) ConvertTops2Package(nodes []*TopNode) (
	redeclareErrors []*RedeclareError, errs []error) {
	//
	if len(nodes) == 0 {
		errs = make([]error, 1)
		errs[0] = fmt.Errorf("nothing to compile")
		return
	}
	if err := PackageBeenCompile.loadCorePackage(); err != nil {
		fmt.Printf("load lucy buildin package failed,err:%v\n", err)
		os.Exit(1)
	}
	errs = make([]error, 0)
	PackageBeenCompile.files = make(map[string]*SourceFile)
	this.Blocks = []*Block{}
	this.Functions = make([]*Function, 0)
	this.Classes = make([]*Class, 0)
	this.Enums = make([]*Enum, 0)
	this.Constants = make([]*Constant, 0)
	expressions := []*Expression{}
	for _, v := range nodes {
		switch v.Node.(type) {
		case *Block:
			t := v.Node.(*Block)
			this.Blocks = append(this.Blocks, t)
		case *Function:
			t := v.Node.(*Function)
			this.Functions = append(this.Functions, t)
		case *Enum:
			t := v.Node.(*Enum)
			this.Enums = append(this.Enums, t)
		case *Class:
			t := v.Node.(*Class)
			this.Classes = append(this.Classes, t)
		case *Constant:
			t := v.Node.(*Constant)
			this.Constants = append(this.Constants, t)
		case *Import:
			i := v.Node.(*Import)
			if i.Alias != UnderScore {
				err := PackageBeenCompile.insertImport(i)
				if err != nil {
					errs = append(errs, err)
				}
			} else {
				if PackageBeenCompile.unUsedPackage == nil {
					PackageBeenCompile.unUsedPackage = make(map[string]*Import)
				}
				PackageBeenCompile.unUsedPackage[i.Import] = i
			}
		case *Expression: // a,b = f();
			t := v.Node.(*Expression)
			if t.Type == ExpressionTypeVar || t.Type == ExpressionTypeVarAssign {
				expressions = append(expressions, t)
			} else {
				errs = append(errs, fmt.Errorf("%s cannot have '%s' in top",
					t.Pos.ErrMsgPrefix(), t.Op))
			}
		case *TypeAlias:
			t := v.Node.(*TypeAlias)
			this.TypeAlias = append(this.TypeAlias, t)
		default:
			panic("tops have unKnow  type")
		}
	}
	redeclareErrors = this.redeclareErrors()
	PackageBeenCompile.Block.Constants = make(map[string]*Constant)
	for _, v := range this.Constants {
		v.IsGlobal = true
		err := PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.Variables = make(map[string]*Variable)
	PackageBeenCompile.Block.Functions = make(map[string]*Function)
	for _, v := range this.Functions {
		v.IsGlobal = true
		if err := PackageBeenCompile.Block.nameIsValid(v.Name, v.Pos); err != nil {
			PackageBeenCompile.errors = append(PackageBeenCompile.errors, err)
		}
		PackageBeenCompile.Block.Functions[v.Name] = v
	}
	PackageBeenCompile.Block.Classes = make(map[string]*Class)
	for _, v := range this.Classes {
		v.IsGlobal = true
		if err := PackageBeenCompile.Block.nameIsValid(v.Name, v.Pos); err != nil {
			PackageBeenCompile.errors = append(PackageBeenCompile.errors, err)
		}
		err := PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.Enums = make(map[string]*Enum)
	PackageBeenCompile.Block.EnumNames = make(map[string]*EnumName)
	for _, v := range this.Enums {
		v.IsGlobal = true
		if err := PackageBeenCompile.Block.nameIsValid(v.Name, v.Pos); err != nil {
			PackageBeenCompile.errors = append(PackageBeenCompile.errors, err)
		}
		err := PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	PackageBeenCompile.Block.TypeAliases = make(map[string]*Type)
	for _, v := range this.TypeAlias {
		if err := PackageBeenCompile.Block.nameIsValid(v.Name, v.Pos); err != nil {
			PackageBeenCompile.errors = append(PackageBeenCompile.errors, err)
		}
		PackageBeenCompile.Block.TypeAliases[v.Name] = v.Type
		v.Type.Alias = v
	}
	if len(expressions) > 0 {
		s := make([]*Statement, len(expressions))
		for k, v := range expressions {
			s[k] = &Statement{
				Type:       StatementTypeExpression,
				Expression: v,
				Pos:        v.Pos,
			}
		}
		b := &Block{}
		b.Pos = expressions[0].Pos
		b.Statements = s
		this.Blocks = append([]*Block{b}, this.Blocks...)
	}
	PackageBeenCompile.mkInitFunctions(this.Blocks)
	return
}

func (this *ConvertTops2Package) redeclareErrors() []*RedeclareError {
	ret := []*RedeclareError{}
	m := make(map[string][]interface{})
	//enums
	for _, v := range this.Enums {
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
	for _, v := range this.Constants {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//functions
	for _, v := range this.Functions {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	//classes
	for _, v := range this.Classes {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}
	// type alias
	for _, v := range this.TypeAlias {
		if _, ok := m[v.Name]; ok {
			m[v.Name] = append(m[v.Name], v)
		} else {
			m[v.Name] = []interface{}{v}
		}
	}

	for name, v := range m {
		if len(v) == 1 || len(v) == 0 { //very good  , 0 looks is impossible
			continue
		}
		r := &RedeclareError{}
		r.Name = name
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
			case *TypeAlias:
				t := vv.(*TypeAlias)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "type alias"
			case *EnumName:
				t := vv.(*EnumName)
				r.Positions[kk] = t.Pos
				r.Types[kk] = "enum name"
			}
		}
		ret = append(ret, r)
	}
	return ret
}
