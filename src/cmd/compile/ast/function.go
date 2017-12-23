package ast

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
)

type FunctionBuildProperty struct {
	IsAnyNumberParameter bool
	CallChcker           func(errs *[]error, args []*VariableType, pos *Pos)
}

type Function struct {
	FunctionBuildProperty
	Isbuildin    bool
	Used         bool
	AccessFlags  uint16 // public private or protected
	Typ          *FunctionType
	Name         string // if name is nil string,means no name function
	Block        *Block
	Pos          *Pos
	Descriptor   string
	Signature    *class_json.MethodSignature
	VariableType VariableType
}

func (f *Function) readableMsg() string {
	s := "fn" + f.Name + "("
	for k, v := range f.Typ.Parameters {
		if k != len(f.Typ.Parameters)-1 {
			s += v.Typ.TypeString() + ","
		} else {
			s += v.Typ.TypeString()
		}
	}
	s += ")"
	if len(f.Typ.Returns) > 0 {
		for k, v := range f.Typ.Returns {
			s += "("
			if k != len(f.Typ.Returns)-1 {
				s += v.Typ.TypeString() + ","
			} else {
				s += v.Typ.TypeString()
			}
			s += ")"
		}
	}
	return s
}

func (f *Function) mkVariableType() {
	f.VariableType.Typ = VARIABLE_TYPE_FUNCTION
	f.VariableType.Function = f
}
func (f *Function) MkVariableType() {
	f.mkVariableType()
}

func (f *Function) MkDescriptor() {
	s := "("
	for _, v := range f.Typ.Parameters {
		s += v.NameWithType.Typ.Descriptor()
	}
	s += ")"
	f.Descriptor = s
}

func (f *Function) check(b *Block) []error {
	f.Block.inherite(b)
	errs := make([]error, 0)
	f.Typ.checkParaMeterAndRetuns(f.Block, errs)
	f.Block.InheritedAttribute.function = f
	errs = append(errs, f.Block.check(b)...)
	return errs
}

func (f *FunctionType) checkParaMeterAndRetuns(block *Block, errs []error) {
	//handler parameter first
	var err error
	errs = []error{}
	var es []error
	for _, v := range f.Parameters {
		v.isFunctionParameter = true
		es = block.checkVar(v)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
			continue
		}
		err = v.Typ.resolve(block)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s", errMsgPrefix(v.Pos), err.Error()))
		}
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	//handler return
	for _, v := range f.Returns {
		es = block.checkVar(v)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
			continue
		}
		err = v.Typ.resolve(block)
		if err != nil {
			errs = append(errs, err)
		}
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s err:%v", errMsgPrefix(v.Pos), err))
		}
	}
}

type FunctionType struct {
	ClosureVars map[string]*VariableDefinition
	Parameters  ParameterList
	Returns     ReturnList
}

type ParameterList []*VariableDefinition // actually local variables
type ReturnList []*VariableDefinition    // actually local variables

func (r ReturnList) retTypes(pos *Pos) []*VariableType {
	if r == nil || len(r) == 0 {
		return mkVoidVariableTypes(pos)
	}
	ret := make([]*VariableType, len(r))
	for k, v := range r {
		ret[k] = v.Typ.Clone()
		ret[k].Pos = pos
	}
	return ret
}
