package ast

import (
	"fmt"
)

func (this *Expression) checkUnaryExpression(block *Block, errs *[]error) *Type {
	ee := this.Data.(*Expression)
	unary, es := ee.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if unary == nil {
		// !a , looks like a bool
		if this.Type == ExpressionTypeNot {
			return &Type{
				Type: VariableTypeBool,
				Pos:  this.Pos,
			}
		}
		return nil
	}
	if err := unary.rightValueValid(); err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if this.Type == ExpressionTypeNot {
		if unary.Type != VariableTypeBool {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression , but '%s'",
				unary.Pos.ErrMsgPrefix(), unary.TypeString()))
		}
	}
	if this.Type == ExpressionTypeNegative {
		if unary.IsNumber() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on '%s'",
				unary.Pos.ErrMsgPrefix(), unary.TypeString()))
		}
	}
	if this.Type == ExpressionTypeBitwiseNot {
		if unary.isInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '~' on '%s'",
				unary.Pos.ErrMsgPrefix(), unary.TypeString()))
		}
	}
	result := unary.Clone()
	result.Pos = this.Pos
	return result
}

func (this *Expression) checkIncrementExpression(block *Block, errs *[]error) *Type {
	on := this.Data.(*Expression)
	increment := on.getLeftValue(block, errs)
	if increment == nil {
		return nil
	}
	if on.Type == ExpressionTypeIdentifier &&
		this.IsStatementExpression == false {
		/*
			special case
			fn1(a++)
		*/
		t := on.Data.(*ExpressionIdentifier)
		if t.Variable != nil {
			t.Variable.Used = true
		}
	}
	if false == increment.IsNumber() {
		*errs = append(*errs,
			fmt.Errorf("%s cannot apply '%s' on '%s'",
				on.Pos.ErrMsgPrefix(), on.Op, increment.TypeString()))
	}
	result := increment.Clone()
	result.Pos = this.Pos
	return result
}
