package ast

import (
	"fmt"
)

func (this *Expression) checkQuestionExpression(block *Block, errs *[]error) *Type {
	question := this.Data.(*ExpressionQuestion)
	condition, es := question.Selection.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if condition != nil {
		if condition.Type != VariableTypeBool {
			*errs = append(*errs,
				fmt.Errorf("%s not a bool expression",
					condition.Pos.ErrMsgPrefix()))
		}
		if err := question.Selection.canBeUsedAsCondition(); err != nil {
			*errs = append(*errs, err)
		}
	}
	tType, es := question.True.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if tType == nil {
		return nil
	}
	if err := tType.rightValueValid(); err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if err := tType.isTyped(); err != nil {
		*errs = append(*errs, err)
		return nil
	}
	ret := tType.Clone()
	ret.Pos = this.Pos
	fType, es := question.False.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if fType != nil &&
		tType.assignAble(errs, fType) == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
			fType.Pos.ErrMsgPrefix(), fType.TypeString(), tType.TypeString()))
	}
	return ret
}
