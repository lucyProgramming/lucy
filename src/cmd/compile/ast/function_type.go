package ast

import "fmt"

type FunctionType struct {
	parameterTypes []*Type
	ParameterList  ParameterList
	ReturnList     ReturnList
	VArgs          *Variable
}

func (functionType *FunctionType) isVargs() *Variable {
	if len(functionType.parameterTypes) == 0 {
		return nil
	}
	t := functionType.ParameterList[len(functionType.ParameterList)-1]
	if t.Type.IsVargs {
		return t
	} else {
		return nil
	}
}

func (functionType *FunctionType) NoReturnValue() bool {
	return len(functionType.ReturnList) == 0 ||
		functionType.ReturnList[0].Type.Type == VariableTypeVoid
}

type ParameterList []*Variable
type ReturnList []*Variable

func (functionType FunctionType) getReturnTypes(pos *Position) []*Type {
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
	if functionType.parameterTypes != nil {
		return functionType.parameterTypes
	}
	ret := make([]*Type, len(functionType.ParameterList))
	for k, v := range functionType.ParameterList {
		ret[k] = v.Type
	}
	functionType.parameterTypes = ret
	return ret
}

func (functionType *FunctionType) fitCallArgs(from *Position, args *CallArgs, block *Block) (fit bool, errs []error, vargs *CallVArgs) {
	errs = []error{}
	length := len(errs)
	callArgsTypes := checkExpressions(block, *args, &errs)
	if len(errs) != length {
		return
	}
	if len(callArgsTypes) > len(functionType.ParameterList) {
		if v := functionType.isVargs(); v != nil {
			for _, t := range callArgsTypes[len(functionType.ParameterList):] {
				if false == v.Type.Equal(&errs, t) {
					errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(t.Pos),
						t.TypeString(), v.Type.TypeString()))
					return
				}
			}
			parameterIndex := 0
			var e *Expression
			var k int
			for k, e = range *args {
				if e.MayHaveMultiValue() == false {
					parameterIndex++
				} else {
					if parameterIndex+len(e.MultiValues) >= len(functionType.parameterTypes) {
						errs = append(errs, fmt.Errorf("%s expression include arg and varg",
							errMsgPrefix(e.Pos)))
						return
					} else {
						parameterIndex += len(e.MultiValues)
					}
				}
				if parameterIndex == len(functionType.parameterTypes)-1 {
					break
				}
			}
			vargs = &CallVArgs{}
			vargs.Es = (*args)[k+1:]
			*args = (*args)[0 : k+1]
			vargs.Length = len(callArgsTypes) - len(functionType.parameterTypes)
		} else {
			errMsg := fmt.Sprintf("%s too many paramaters to call\n", errMsgPrefix(from))
			errMsg += fmt.Sprintf("\thave %s\n", callHave(callArgsTypes))
			errMsg += fmt.Sprintf("\twant %s\n", callWant(functionType.ParameterList))
			errs = append(errs, fmt.Errorf(errMsg))
		}
	}
	//trying to convert literal
	convertLiteralExpressionsToNeeds(*args, functionType.getParameterTypes(), callArgsTypes)
	if len(callArgsTypes) < len(functionType.ParameterList) {
		errMsg := fmt.Sprintf("%s too few paramaters to call\n", errMsgPrefix(from))
		errMsg += fmt.Sprintf("\thave %s\n", callHave(callArgsTypes))
		errMsg += fmt.Sprintf("\twant %s\n", callWant(functionType.ParameterList))
		errs = append(errs, fmt.Errorf(errMsg))
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
	fit = len(errs) == 0
	return
}

type CallVArgs struct {
	Es     []*Expression
	Length int
}
