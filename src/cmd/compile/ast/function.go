package ast

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
)

type CallChecker func(errs *[]error, args []*VariableType, pos *Pos)

type Function struct {
	isPackageBlockFunction     bool
	callchecker                CallChecker // used in build function
	ClassMethod                *cg.MethodHighLevel
	IsGlobal                   bool
	Isbuildin                  bool
	Used                       bool
	AccessFlags                uint16 // public private or protected
	Typ                        *FunctionType
	ClosureVars                ClosureVars
	Name                       string // if name is nil string,means no name function
	Block                      *Block
	Pos                        *Pos
	Descriptor                 string
	Signature                  *class_json.MethodSignature
	VariableType               VariableType
	Varoffset                  uint16
	ArrayListVarForMultiReturn ArrayListVarForMultiReturn
}
type ArrayListVarForMultiReturn struct {
	Offset uint16
}

func (f *Function) IsClosureFunction() bool {
	return f.ClosureVars.NotEmpty()
}

func (f *Function) ClosureVarExist(name string, v *VariableDefinition) (uint8, bool) {
	return f.ClosureVars.ClosureVarsExist(name, v)
}
func (f *Function) readableMsg() string {
	s := "fn" + f.Name + "("
	for k, v := range f.Typ.ParameterList {
		if k != len(f.Typ.ParameterList)-1 {
			s += v.Typ.TypeString() + ","
		} else {
			s += v.Typ.TypeString()
		}
	}
	s += ")"
	if len(f.Typ.ReturnList) > 0 {
		s += "->"
		for k, v := range f.Typ.ReturnList {
			s += "("
			if k != len(f.Typ.ReturnList)-1 {
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
	for _, v := range f.Typ.ParameterList {
		s += v.NameWithType.Typ.Descriptor()
	}
	s += ")"
	f.Descriptor = s
}

func (f *Function) checkBlocks(errs *[]error) {
	f.Block.InheritedAttribute.function = f
	*errs = append(*errs, f.Block.check(nil)...)
}

func (f *Function) checkParaMeterAndRetuns(errs *[]error) {
	f.Typ.checkParaMeterAndRetuns(f.Block, errs)
}

func (f *Function) check(b *Block) []error {
	errs := make([]error, 0)
	f.Block.inherite(b)
	f.Block.InheritedAttribute.function = f
	f.checkParaMeterAndRetuns(&errs)
	//
	{
		s := &StatementReturn{}
		es := []*Expression{}
		for _, v := range f.Typ.ReturnList {
			identifer := &ExpressionIdentifer{}
			identifer.Name = v.Name
			es = append(es, &Expression{
				Typ:  EXPRESSION_TYPE_IDENTIFIER,
				Data: identifer,
			})
		}
		f.Block.Statements = append(f.Block.Statements, &Statement{Typ: STATEMENT_TYPE_RETURN, StatementReturn: s})
	}

	f.checkBlocks(&errs)
	return errs
}

func (f *FunctionType) checkParaMeterAndRetuns(functionblock *Block, errs *[]error) {
	var err error
	var es []error
	for _, v := range f.ParameterList {
		v.IsFunctionParameter = true
		es = functionblock.checkVar(v)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			continue
		}
		err = v.Typ.resolve(functionblock)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(v.Pos), err.Error()))
		}
		err = functionblock.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
			continue
		}
	}
	//handler return
	for _, v := range f.ReturnList {
		es = functionblock.checkVar(v)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			continue
		}
		err = v.Typ.resolve(functionblock)
		if err != nil {
			*errs = append(*errs, err)
		}
		err = functionblock.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s err:%v", errMsgPrefix(v.Pos), err))
		}
	}
}

type FunctionType struct {
	ParameterList ParameterList
	ReturnList    ReturnList
}

type ParameterList []*VariableDefinition // actually local variables
type ReturnList []*VariableDefinition    // actually local variables

func (r ReturnList) retTypes(pos *Pos) []*VariableType {
	if r == nil || len(r) == 0 {
		t := &VariableType{}
		t.Typ = VARIABLE_TYPE_FUNCTION
		return []*VariableType{t}
	}
	ret := make([]*VariableType, len(r))
	for k, v := range r {
		ret[k] = v.Typ.Clone()
		ret[k].Pos = pos
	}
	return ret
}
