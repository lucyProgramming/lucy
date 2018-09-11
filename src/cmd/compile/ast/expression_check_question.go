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
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(question.Selection.Pos)))
		}
		if question.Selection.canBeUsedAsCondition() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as condition",
				errMsgPrefix(question.Selection.Pos), question.Selection.Description))
		}
	}
	True, es := question.True.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if True != nil {
		if True.RightValueValid() == false {
			*errs = append(*errs, fmt.Errorf("%s not right value valid",
				errMsgPrefix(question.True.Pos)))
			return nil
		}
		if True.isTyped() == false {
			*errs = append(*errs, fmt.Errorf("%s not typed",
				errMsgPrefix(question.True.Pos)))
			return nil
		}
	}
	False, es := question.False.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if True != nil &&
		False != nil &&
		True.assignAble(errs, False) == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
			errMsgPrefix(question.False.Pos), False.TypeString(), True.TypeString()))
	}
	if True != nil {
		result := True.Clone()
		result.Pos = e.Pos
		return result
	}
	if False != nil {
		result := False.Clone()
		result.Pos = e.Pos
		return result
	}
	return nil
}
