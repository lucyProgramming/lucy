package ast

import (
	"fmt"
)

func (e *Expression) checkUnaryExpression(block *Block, errs *[]error) *VariableType {
	ee := e.Data.(*Expression)
	t, es := ee.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}

	if t == nil {
		if e.Typ == EXPRESSION_TYPE_NOT {
			return &VariableType{
				Typ: EXPRESSION_TYPE_BOOL,
				Pos: e.Pos,
			}
		} else {
			return &VariableType{
				Typ: EXPRESSION_TYPE_INT,
				Pos: e.Pos,
			}
		}
	}
	if e.Typ == EXPRESSION_TYPE_NOT {
		if t.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression",
				errMsgPrefix(t.Pos)))
		}
	}
	if e.Typ == EXPRESSION_TYPE_NEGATIVE {
		if t.IsNumber() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on '%s'",
				errMsgPrefix(e.Pos), t.TypeString()))
		}
	}
	if e.Typ == EXPRESSION_TYPE_BITWISE_NOT {
		if t.IsInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '~' on '%s'",
				errMsgPrefix(e.Pos), t.TypeString()))
		}
	}
	ret := t.Clone()
	ret.Pos = e.Pos
	return ret
}
func (e *Expression) checkIncrementExpression(block *Block, errs *[]error) *VariableType {
	ee := e.Data.(*Expression)
	t := ee.getLeftValue(block, errs)
	ee.Value = t
	if t == nil {
		return nil
	}
	if !t.IsNumber() {
		*errs = append(*errs, fmt.Errorf("%s cannot apply '++' or '--' on '%s'",
			errMsgPrefix(ee.Pos), t.TypeString()))
	}
	tt := t.Clone()
	tt.Pos = e.Pos
	return tt
}
