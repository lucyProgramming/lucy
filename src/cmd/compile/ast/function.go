package ast

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
)

type Function struct {
	AccessFlags uint16 // public private or protected
	Typ         *FunctionType
	Name        string // if name is nil string,means no name function
	Block       *Block
	Pos         *Pos
	Descriptor  string
	Signature   *class_json.MethodSignature
}

func (f *Function) MkDescriptor() {
	s := "("
	for _, v := range f.Typ.Parameters {
		s += v.NameWithType.Typ.Descriptor() + ";"
	}
	s += ")"
	f.Descriptor = s
}

func (f *Function) check(b *Block) []error {
	f.Block.inherite(b)
	errs := make([]error, 0)
	f.Typ.checkParaMeterAndRetuns(f.Block, errs)
	f.Block.InheritedAttribute.infunction = true
	f.Block.InheritedAttribute.returns = f.Typ.Returns
	errs = append(errs, f.Block.check(nil)...)
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
			err = block.insert(v.Name, nil, vd)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %d:%d err:%v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
				continue
			}
			if v.Expression != nil {
				is, typ, value, err := v.Expression.getConstValue()
				if err != nil {
					errs = append(errs, fmt.Errorf("%s default value is wrong because of %v", errMsgPrefix(v.Pos), err))
					continue
				}
				if !is {
					errs = append(errs, fmt.Errorf("%s default value is not a const value %v", errMsgPrefix(v.Pos), err))
					continue
				}
				t, _ := block.getTypeFromExpression(&Expression{Typ: typ, Data: value})
				if !v.Typ.typeCompatible(t) {
					errs = append(errs, fmt.Errorf("%s default value can not assign to variale type %v", errMsgPrefix(v.Pos)))
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
			err = block.insert(v.Name, nil, t)
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

type ParameterList []*VariableDefinition // actually local variables
type ReturnList []*VariableDefinition    // actually local variables
