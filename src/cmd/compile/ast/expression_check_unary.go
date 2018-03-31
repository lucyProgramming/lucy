package ast

import (
	"fmt"
)

func (e *Expression) checkUnaryExpression(block *Block, errs *[]error) *VariableType {
	ee := e.Data.(*Expression)
	ts, es := ee.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		if e.Typ == EXPRESSION_TYPE_NOT {
			return &VariableType{
				Typ: EXPRESSION_TYPE_BOOL,
				Pos: e.Pos,
			}
		}
		return nil
	}
	if e.Typ == EXPRESSION_TYPE_NOT {
		if t.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(t.Pos)))
		}
		t := &VariableType{
			Typ: EXPRESSION_TYPE_BOOL,
			Pos: e.Pos,
		}
		return t
	} else if e.Typ == EXPRESSION_TYPE_NEGATIVE {
		if t.IsNumber() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on '%s'", errMsgPrefix(e.Pos), t.TypeString()))
		}
		tt := t.Clone()
		tt.Pos = e.Pos
		return tt
	}
	return nil
}
func (e *Expression) checkIncrementExpression(block *Block, errs *[]error) *VariableType {
	ee := e.Data.(*Expression)
	t, es := ee.getLeftValue(block)
	ee.VariableType = t
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	if !t.IsNumber() {
		*errs = append(*errs, fmt.Errorf("%s cannot apply ++ or -- on '%s'",
			errMsgPrefix(ee.Pos), t.TypeString()))
		return nil
	}
	tt := t.Clone()
	tt.Pos = e.Pos
	return tt
}
