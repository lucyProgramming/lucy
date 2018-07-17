package ast

import (
	"fmt"
)

func (e *Expression) checkUnaryExpression(block *Block, errs *[]error) *Type {
	ee := e.Data.(*Expression)
	unary, es := ee.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if unary == nil {
		if e.Type == ExpressionTypeNot {
			return &Type{
				Type: ExpressionTypeBool,
				Pos:  e.Pos,
			}
		}
		return nil
	}
	if e.Type == ExpressionTypeNot {
		if unary.Type != VariableTypeBool {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression",
				errMsgPrefix(unary.Pos)))
		}
	}
	if e.Type == ExpressionTypeNegative {
		if unary.IsNumber() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on '%s'",
				errMsgPrefix(e.Pos), unary.TypeString()))
		}
	}
	if e.Type == ExpressionTypeBitwiseNot {
		if unary.IsInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '~' on '%s'",
				errMsgPrefix(e.Pos), unary.TypeString()))
		}
	}
	ret := unary.Clone()
	ret.Pos = e.Pos
	return ret
}

func (e *Expression) checkIncrementExpression(block *Block, errs *[]error) *Type {
	on := e.Data.(*Expression)
	t := on.getLeftValue(block, errs)
	on.ExpressionValue = t
	if t == nil {
		return nil
	}
	if !t.IsNumber() {
		*errs = append(*errs, fmt.Errorf("%s cannot apply '%s' on '%s'",
			errMsgPrefix(on.Pos), on.OpName(), t.TypeString()))
	}
	ret := t.Clone()
	ret.Pos = e.Pos
	return ret
}
