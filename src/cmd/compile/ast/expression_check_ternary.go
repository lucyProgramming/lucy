package ast

import (
	"fmt"
)

func (e *Expression) checkTernaryExpression(block *Block, errs *[]error) *VariableType {
	ternary := e.Data.(*ExpressionTernary)
	condition, es := ternary.Condition.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if condition != nil {
		if condition.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(e.Pos)))
		}
		if ternary.Condition.canbeUsedAsCondition() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as condition",
				errMsgPrefix(e.Pos), e.OpName()))
		}
	}
	True, es := ternary.True.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	False, es := ternary.False.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if True != nil && False != nil && True.Equal(False) == false {
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
