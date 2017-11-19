package ast

import (
	"fmt"
)

type Function struct {
	Access int // public private or protected
	Typ    FunctionType
	Name   string // if name is nil string,means no name function
	Block  *Block
	Pos    *Pos
}

func (f *Function) check(b *Block) []error {
	if b != nil {
		f.Block.inherite(b)
	}
	errs := make([]error, 0)
	f.Typ.checkParaMeterAndRetuns(f.Block, errs)
	f.Block.InheritedAttribute.infunction = true
	f.Block.InheritedAttribute.returns = f.Typ.Returns
	errs = append(errs, f.Block.check()...)
	return errs
}

func (f *FunctionType) checkParaMeterAndRetuns(block *Block, errs []error) {
	//handler parameter first
	var err error
	for _, v := range f.Parameters {
		if v.Name != "" {
			vd := &VariableDefinition{}
			vd.Name = v.Name
			vd.Typ = v.Typ
			err = block.SymbolicTable.Insert(v.Name, nil, vd)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %d:%d err:%v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
				continue
			}
			if v.Expression != nil {
				is, typ, value, err := v.Expression.getConstValue()
				if err != nil {
					errs = append(errs, fmt.Errorf("%s %d:%d default value is wrong because of %v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
					continue
				}
				if !is {
					errs = append(errs, fmt.Errorf("%s %d:%d default value is not a const value %v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
					continue
				}
				t, _ := block.getTypeFromExpression(&Expression{Typ: typ, Data: value})
				if !v.Typ.typeCompatible(t) {
					errs = append(errs, fmt.Errorf("%s %d:%d default value can not assign to variale type %v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn))
					continue
				}
			}
		}
	}
	//handler return
	for _, v := range f.Returns {
		if v.Name != "" {
			t := VariableDefinition{}
			t.Name = v.Name
			t.Typ = v.Typ
			err = block.SymbolicTable.Insert(v.Name, nil, t)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %d:%d err:%v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
				continue
			}
		}
	}
}

type FunctionType struct {
	Parameters ParameterList
	Returns    ReturnList
}

//type Parameter struct {
//	VariableDefinition
//	Default *Expression //f(a int = 1) default parameter
//}

type ParameterList []*VariableDefinition // actually local variables
type ReturnList []*VariableDefinition    // actually local variables
