package ast

type FunctionType struct {
	parameterTypes []*Type
	ParameterList  ParameterList
	ReturnList     ReturnList
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
