package ast

type FunctionType struct {
	ParameterList ParameterList
	ReturnList    ReturnList
}

type ParameterList []*Variable
type ReturnList []*Variable

func (ft FunctionType) returnTypes(pos *Position) []*Type {
	if ft.ReturnList == nil || len(ft.ReturnList) == 0 {
		t := &Type{}
		t.Type = VARIABLE_TYPE_VOID // means no return;
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
