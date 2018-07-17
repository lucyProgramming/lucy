package ast

import (
	"fmt"
)

func (e *Expression) checkTernaryExpression(block *Block, errs *[]error) *Type {
	question := e.Data.(*ExpressionQuestion)
	condition, es := question.Selection.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if condition != nil {
		if condition.Type != VariableTypeBool {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(e.Pos)))
		}
		if question.Selection.canBeUsedAsCondition() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as condition",
				errMsgPrefix(e.Pos), e.OpName()))
		}
	}
	True, es := question.True.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	False, es := question.False.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if True != nil && False != nil && True.Equal(errs, False) == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
			errMsgPrefix(e.Pos), False.TypeString(), True.TypeString()))
	}
	if True != nil {
		tt := True.Clone()
		tt.Pos = e.Pos
		return tt
	}
	if False != nil {
		tt := False.Clone()
		tt.Pos = e.Pos
		return tt
	}
	return nil
}
