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
		if e.Type == ExpressionTypeNot {
			return &Type{
				Type: ExpressionTypeBool,
				Pos:  e.Pos,
			}
		}
		return nil
	}
	if e.Type == ExpressionTypeNot {
		if t.Type != VariableTypeBool {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression",
				errMsgPrefix(t.Pos)))
		}
	}
	if e.Type == ExpressionTypeNegative {
		if t.IsNumber() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on '%s'",
				errMsgPrefix(e.Pos), t.TypeString()))
		}
	}
	if e.Type == ExpressionTypeBitwiseNot {
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
