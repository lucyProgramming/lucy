package ast

type FunctionType struct {
	ParameterList ParameterList
	ReturnList    ReturnList
}

type ParameterList []*VariableDefinition
type ReturnList []*VariableDefinition

func (ft FunctionType) retTypes(pos *Pos) []*VariableType {
	if ft.ReturnList == nil || len(ft.ReturnList) == 0 {
		t := &VariableType{}
		t.Type = VARIABLE_TYPE_VOID // means no return;
		t.Pos = pos
		return []*VariableType{t}
	}
	ret := make([]*VariableType, len(ft.ReturnList))
	for k, v := range ft.ReturnList {
		ret[k] = v.Type.Clone()
		ret[k].Pos = pos
	}
	return ret
}

func (ft FunctionType) getParameterTypes() []*VariableType {
	ret := make([]*VariableType, len(ft.ParameterList))
	for k, v := range ft.ParameterList {
		ret[k] = v.Type
	}
	return ret
}
