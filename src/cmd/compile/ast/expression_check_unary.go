package ast

import (
	"fmt"
)

func (e *Expression) checkUnaryExpression(block *Block, errs *[]error) *Type {
	ee := e.Data.(*Expression)
	t, es := ee.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		if e.Type == EXPRESSION_TYPE_NOT {
			return &Type{
				Type: EXPRESSION_TYPE_BOOL,
				Pos:  e.Pos,
			}
		} else {
			return &Type{
				Type: EXPRESSION_TYPE_INT,
				Pos:  e.Pos,
			}
		}
	}
	if e.Type == EXPRESSION_TYPE_NOT {
		if t.Type != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression",
				errMsgPrefix(t.Pos)))
		}
	}
	if e.Type == EXPRESSION_TYPE_NEGATIVE {
		if t.IsNumber() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on '%s'",
				errMsgPrefix(e.Pos), t.TypeString()))
		}
	}
	if e.Type == EXPRESSION_TYPE_BIT_NOT {
		if t.IsInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '~' on '%s'",
				errMsgPrefix(e.Pos), t.TypeString()))
		}
	}
	ret := t.Clone()
	ret.Pos = e.Pos
	return ret
}
func (e *Expression) checkIncrementExpression(block *Block, errs *[]error) *Type {
	ee := e.Data.(*Expression)
	t := ee.getLeftValue(block, errs)
	ee.ExpressionValue = t
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
