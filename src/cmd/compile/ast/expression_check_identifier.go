package ast

import (
	"fmt"
)

func (e *Expression) checkIdentiferExpression(block *Block) (t *VariableType, err error) {
	identifer := e.Data.(*ExpressionIdentifer)
	d := block.SearchByName(identifer.Name)
	if d == nil {
		return nil, fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), identifer.Name)
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		f.Used = true
		tt := &VariableType{}
		tt.Typ = VARIABLE_TYPE_FUNCTION
		tt.Pos = e.Pos
		tt.Function = f
		return tt, nil
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		if t.Typ == nil { // in some case,variable defined wrong,could be nil
			return nil, nil
		}
		t.Used = true
		tt := t.Typ.Clone()
		tt.Pos = e.Pos
		identifer.Var = t
		return tt, nil
	case *Const:
		t := d.(*Const)
		t.Used = true
		e.fromConst(t)
		tt := t.Typ.Clone()
		tt.Pos = e.Pos
		return tt, nil
	//case *Enum:
	//	t := d.(*Enum)
	//	t.Used = true
	//	tt := t.VariableType.Clone()
	//	tt.Pos = e.Pos
	//	identifer.Enum = t
	//	return tt, nil
	//case *EnumName:
	//	t := d.(*EnumName)
	//	t.Enum.Used = true
	//	tt := t.Enum.VariableType.Clone()
	//	tt.Pos = e.Pos
	//	identifer.EnumName = t
	//	return tt, nil
	case *Class:
		t := &VariableType{}
		t.Typ = VARIABLE_TYPE_CLASS
		e.Pos = e.Pos
		t.Class = d.(*Class)
		return t, nil
	default:
		return nil, fmt.Errorf("%s identifier '%s' is not a expression", errMsgPrefix(e.Pos), identifer.Name)
	}
	return nil, nil
}

func (e *Expression) isThisIdentifierExpression() (is bool) {
	if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
		return
	}
	t := e.Data.(*ExpressionIdentifer)
	is = (t.Name == THIS)
	return
}
