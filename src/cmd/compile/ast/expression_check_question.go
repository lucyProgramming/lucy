package ast

import (
	"fmt"
)

func (e *Expression) checkQuestionExpression(block *Block, errs *[]error) *Type {
	question := e.Data.(*ExpressionQuestion)
	condition, es := question.Selection.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if condition != nil {
		if condition.Type != VariableTypeBool {
			*errs = append(*errs,
				fmt.Errorf("%s not a bool expression",
					errMsgPrefix(question.Selection.Pos)))
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
	ret.Pos = e.Pos
	fType, es := question.False.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if fType == nil {
		return ret
	}
	if tType.assignAble(errs, fType) == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
			errMsgPrefix(question.False.Pos), fType.TypeString(), tType.TypeString()))
	}
	return ret
}
