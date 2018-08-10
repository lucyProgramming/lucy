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
				Type: VariableTypeBool,
				Pos:  e.Pos,
			}
		}
		return nil
	}
	if unary.RightValueValid() == false {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not right value valid",
			errMsgPrefix(ee.Pos), unary.TypeString()))
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
				errMsgPrefix(unary.Pos), unary.TypeString()))
		}
	}
	if e.Type == ExpressionTypeBitwiseNot {
		if unary.IsInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '~' on '%s'",
				errMsgPrefix(unary.Pos), unary.TypeString()))
		}
	}
	result := unary.Clone()
	result.Pos = e.Pos
	return result
}

func (e *Expression) checkIncrementExpression(block *Block, errs *[]error) *Type {
	on := e.Data.(*Expression)
	increment := on.getLeftValue(block, errs)
	if increment == nil {
		return nil
	}
	if on.Type == ExpressionTypeIdentifier &&
		e.IsStatementExpression == false {
		on.Data.(*ExpressionIdentifier).Variable.Used = true
	}
	if !increment.IsNumber() {
		*errs = append(*errs, fmt.Errorf("%s cannot apply '%s' on '%s'",
			errMsgPrefix(on.Pos), on.OpName(), increment.TypeString()))
	}
	ret := increment.Clone()
	ret.Pos = e.Pos
	return ret
}
