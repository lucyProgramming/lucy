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
	errs = []error{}
	var es []error
	for _, v := range f.Parameters {
		if v.Name != "" {
			es = block.checkVar(v)
			if errsNotEmpty(es) {
				errs = append(errs, es...)
				continue
			}
			err = block.insert(v.Name, v.Pos, v)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s err:%v", errMsgPrefix(v.Pos), err))
				continue
			}
		}
	}
	//handler return
	for _, v := range f.Returns {
		es = block.checkVar(v)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
			continue
		}
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %d:%d err:%v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
		}
	}
}

type FunctionType struct {
	Parameters ParameterList
	Returns    ReturnList
}

type ParameterList []*VariableDefinition // actually local variables
type ReturnList []*VariableDefinition    // actually local variables
