package ast

import (
	"fmt"
)

type FunctionType struct {
	ParameterList ParameterList
	ReturnList    ReturnList
	VArgs         *Variable
}
type ParameterList []*Variable
type ReturnList []*Variable

func (ft *FunctionType) equal(compare *FunctionType) bool {
	if len(ft.ParameterList) != len(compare.ParameterList) ||
		len(ft.ReturnList) != len(compare.ReturnList) {
		return false
	}
	for k, v := range ft.ParameterList {
		if false == v.Type.StrictEqual(compare.ParameterList[k].Type) {
			return false
		}
	}
	for k, v := range ft.ReturnList {
		if false == v.Type.StrictEqual(compare.ReturnList[k].Type) {
			return false
		}
	}
	return true
}

func (ft *FunctionType) NoReturnValue() bool {
	return len(ft.ReturnList) == 0 ||
		ft.ReturnList[0].Type.Type == VariableTypeVoid
}

func (ft FunctionType) getReturnTypes(pos *Pos) []*Type {
	if ft.ReturnList == nil || len(ft.ReturnList) == 0 {
		t := &Type{}
		t.Type = VariableTypeVoid // means no return ;
		t.Pos = pos
		return []*Type{t}
	}
	ret := make([]*Type, len(ft.ReturnList))
	for k, v := range ft.ReturnList {
		ret[k] = v.Type.Clone()
		ret[k].Pos = pos
	}
	return ret
}

func (ft FunctionType) getParameterTypes() []*Type {
	ret := make([]*Type, len(ft.ParameterList))
	for k, v := range ft.ParameterList {
		ret[k] = v.Type
	}
	return ret
}

func (ft *FunctionType) fitCallArgs(from *Pos, args *CallArgs,
	callArgsTypes []*Type, f *Function) (vArgs *CallVArgs, err error) {
	//trying to convert literal
	convertExpressionsToNeeds(*args, ft.getParameterTypes(), callArgsTypes)

	if ft.VArgs != nil {
		vArgs = &CallVArgs{}
		vArgs.NoArgs = true
		vArgs.Type = ft.VArgs.Type
	}
	msg := fmt.Sprintf("\thave %s\n", callHave(callArgsTypes))
	msg += fmt.Sprintf("\twant %s\n", callWant(ft))
	errs := []error{}
	if len(callArgsTypes) > len(ft.ParameterList) {
		if ft.VArgs == nil {
			errMsg := fmt.Sprintf("%s too many paramaters to call\n", errMsgPrefix(from))
			errMsg += msg
			err = fmt.Errorf(errMsg)
			return
		}
		v := ft.VArgs
		for _, t := range callArgsTypes[len(ft.ParameterList):] {
			if t == nil { // some error before
				return
			}
			if t.IsVArgs {
				if len(callArgsTypes[len(ft.ParameterList):]) > 1 {
					errMsg := fmt.Sprintf("%s too many argument to call\n",
						errMsgPrefix(t.Pos))
					errMsg += msg
					err = fmt.Errorf(errMsg)
					return
				}
				if false == v.Type.Equal(&errs, t) {
					err = fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(t.Pos),
						t.TypeString())
					return
				}
				vArgs.PackageJavaArray2VArgs = true
				continue
			}
			if false == v.Type.Array.Equal(&errs, t) {
				err = fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix(t.Pos),
					t.TypeString(), v.Type.TypeString())
				return
			}
		}
		vArgs.NoArgs = false
		k := len(ft.ParameterList)
		vArgs.Length = len(callArgsTypes) - k
		vArgs.Expressions = (*args)[k:]
		*args = (*args)[:k]
		vArgs.Length = len(callArgsTypes) - k
	}

	if len(callArgsTypes) < len(ft.ParameterList) {
		if f != nil && f.HaveDefaultValue && len(callArgsTypes) >= f.DefaultValueStartAt {
			for i := len(callArgsTypes); i < len(f.Type.ParameterList); i++ {
				*args = append(*args, f.Type.ParameterList[i].Expression)
			}
		} else { // no default value
			errMsg := fmt.Sprintf("%s too few paramaters to call\n", errMsgPrefix(from))
			errMsg += msg
			err = fmt.Errorf(errMsg)
			return
		}
	}
	for k, v := range ft.ParameterList {
		if k < len(callArgsTypes) && callArgsTypes[k] != nil {
			if false == v.Type.Equal(&errs, callArgsTypes[k]) {
				errMsg := fmt.Sprintf("%s cannot use '%s' as '%s'",
					errMsgPrefix((callArgsTypes)[k].Pos),
					callArgsTypes[k].TypeString(), v.Type.TypeString())
				errMsg += msg
				err = fmt.Errorf(errMsg)
				return
			}
		}
	}
	return
}

type CallVArgs struct {
	Expressions []*Expression
	Length      int
	/*
		a := new int[](10)
		print(a...)
	*/
	PackageJavaArray2VArgs bool
	NoArgs                 bool
	Type                   *Type
}
