package ast

import (
	"errors"
	"fmt"
)

type FunctionType struct {
	ParameterList ParameterList
	ReturnList    ReturnList
	VArgs         *Variable
}
type ParameterList []*Variable
type ReturnList []*Variable

func (functionType *FunctionType) NoReturnValue() bool {
	return len(functionType.ReturnList) == 0 ||
		functionType.ReturnList[0].Type.Type == VariableTypeVoid
}

func (functionType FunctionType) getReturnTypes(pos *Pos) []*Type {
	if functionType.ReturnList == nil || len(functionType.ReturnList) == 0 {
		t := &Type{}
		t.Type = VariableTypeVoid // means no return ;
		t.Pos = pos
		return []*Type{t}
	}
	ret := make([]*Type, len(functionType.ReturnList))
	for k, v := range functionType.ReturnList {
		ret[k] = v.Type.Clone()
		ret[k].Pos = pos
	}
	return ret
}

func (functionType FunctionType) getParameterTypes() []*Type {
	ret := make([]*Type, len(functionType.ParameterList))
	for k, v := range functionType.ParameterList {
		ret[k] = v.Type
	}
	return ret
}

func (functionType *FunctionType) fitCallArgs(from *Pos, args *CallArgs,
	callArgsTypes []*Type, f *Function) (match bool, vArgs *CallVArgs, errs []error) {
	//trying to convert literal
	convertLiteralExpressionsToNeeds(*args, functionType.getParameterTypes(), callArgsTypes)
	errs = []error{}
	for _, v := range *args {
		if v.MayHaveMultiValue() && len(v.MultiValues) > 1 {
			errs = append(errs, fmt.Errorf("%s multi value in single value context",
				errMsgPrefix(from)))
			return
		}
	}
	if functionType.VArgs != nil {
		vArgs = &CallVArgs{}
		vArgs.NoArgs = true
		vArgs.Type = functionType.VArgs.Type
	}
	if len(callArgsTypes) > len(functionType.ParameterList) {
		if functionType.VArgs == nil {
			errMsg := fmt.Sprintf("%s too many paramaters to call\n", errMsgPrefix(from))
			errMsg += fmt.Sprintf("\thave %s\n", callHave(callArgsTypes))
			errMsg += fmt.Sprintf("\twant %s\n", callWant(functionType))
			errs = append(errs, fmt.Errorf(errMsg))
			return // no further check
		}
		v := functionType.VArgs
		for _, t := range callArgsTypes[len(functionType.ParameterList):] {
			if t == nil { // some error before
				return
			}
			if t.IsVArgs {
				if len(callArgsTypes[len(functionType.ParameterList):]) > 1 {
					errMsg := fmt.Sprintf("%s too many argument to call\n",
						errMsgPrefix(t.Pos))
					errMsg += fmt.Sprintf("\thave %s\n", callHave(callArgsTypes))
					errMsg += fmt.Sprintf("\twant %s\n", callWant(functionType))
					errs = append(errs, errors.New(errMsg))
					return
				}
				if false == v.Type.Equal(&errs, t) {
					errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(t.Pos),
						t.TypeString(), v.Type.TypeString()))
					return
				}
				vArgs.IsJavaArray = true
				continue
			}
			if false == v.Type.Array.Equal(&errs, t) {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix(t.Pos),
					t.TypeString(), v.Type.TypeString()))
				return
			}
		}
		vArgs.NoArgs = false
		k := len(functionType.ParameterList)
		vArgs.Length = len(callArgsTypes) - k
		vArgs.Expressions = (*args)[k:]
		*args = (*args)[:k]
		vArgs.Length = len(callArgsTypes) - k
	}

	if len(callArgsTypes) < len(functionType.ParameterList) {
		if f != nil && f.HaveDefaultValue && len(callArgsTypes) >= f.DefaultValueStartAt {
			for i := len(callArgsTypes); i < len(f.Type.ParameterList); i++ {
				*args = append(*args, f.Type.ParameterList[i].Expression)
			}
		} else { // no default value
			errMsg := fmt.Sprintf("%s too few paramaters to call\n", errMsgPrefix(from))
			errMsg += fmt.Sprintf("\thave %s\n", callHave(callArgsTypes))
			errMsg += fmt.Sprintf("\twant %s\n", callWant(functionType))
			errs = append(errs, fmt.Errorf(errMsg))
		}
	}
	for k, v := range functionType.ParameterList {
		if k < len(callArgsTypes) && callArgsTypes[k] != nil {
			if false == v.Type.Equal(&errs, callArgsTypes[k]) {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix((callArgsTypes)[k].Pos),
					callArgsTypes[k].TypeString(), v.Type.TypeString()))
			}
		}
	}
	match = len(errs) == 0
	return
}

type CallVArgs struct {
	Expressions []*Expression
	Length      int
	/*
		a := new int[](10)
		print(a...)
	*/
	IsJavaArray bool
	NoArgs      bool
	Type        *Type
}
