package ast

import (
	"fmt"
)

func (e *Expression) checkIdentiferExpression(block *Block) (t *VariableType, err error) {
	identifer := e.Data.(*ExpressionIdentifer)
	d := block.searchByName(identifer.Name)
	if d == nil {
		return nil, fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), identifer.Name)
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		f.Used = true
		tt := f.VariableType.Clone()
		tt.Pos = e.Pos
		identifer.Func = f
		return tt, nil
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		if t.Typ == nil { // in some case,variable defined wrong,could be nil
			return nil, nil
		}
		t.CaptureLevel = 0 //when caputre is done,level reset to zero
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
	case *Enum:
		t := d.(*Enum)
		t.Used = true
		tt := t.VariableType.Clone()
		tt.Pos = e.Pos
		identifer.Enum = t
		return tt, nil
	case *EnumName:
		t := d.(*EnumName)
		t.Enum.Used = true
		tt := t.Enum.VariableType.Clone()
		tt.Pos = e.Pos
		identifer.EnumName = t
		return tt, nil
	default:
		return nil, fmt.Errorf("%s identifier '%s' is not a expression", errMsgPrefix(e.Pos), identifer.Name)
	}
	return nil, nil
}

func (e *Expression) fromConst(c *Const) {
	switch c.Typ.Typ {
	case VARIABLE_TYPE_BOOL:
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data = c.Data.(bool)
	case VARIABLE_TYPE_BYTE:
		e.Typ = EXPRESSION_TYPE_BYTE
		e.Data = c.Data.(byte)
	case VARIABLE_TYPE_INT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = c.Data.(int32)
	case VARIABLE_TYPE_LONG:
		e.Typ = EXPRESSION_TYPE_LONG
		e.Data = c.Data.(int64)
	case VARIABLE_TYPE_FLOAT:
		e.Typ = EXPRESSION_TYPE_FLOAT
		e.Data = c.Data.(float32)
	case VARIABLE_TYPE_DOUBLE:
		e.Typ = EXPRESSION_TYPE_DOUBLE
		e.Data = c.Data.(float64)
	case VARIABLE_TYPE_STRING:
		e.Typ = EXPRESSION_TYPE_STRING
		e.Data = c.Data.(string)
	}
}

func (e *Expression) isThisIdentifierExpression() (b bool) {
	if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
		return
	}
	t := e.Data.(*ExpressionIdentifer)
	b = (t.Name == THIS)
	return
}
