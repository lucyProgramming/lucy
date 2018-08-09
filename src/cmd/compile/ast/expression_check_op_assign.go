package ast

import "fmt"

func (e *Expression) checkOpAssignExpression(block *Block, errs *[]error) (t *Type) {
	bin := e.Data.(*ExpressionBinary)
	left := bin.Left.getLeftValue(block, errs)
	bin.Left.Value = left
	right, es := bin.Right.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if left == nil || right == nil {
		return
	}
	result := left.Clone()
	result.Pos = e.Pos
	if right.RightValueValid() == false {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not right value valid",
			errMsgPrefix(bin.Right.Pos), right.TypeString()))
		return result
	}

	if bin.Left.Type == ExpressionTypeIdentifier && e.IsStatementExpression == false {
		t := bin.Left.Data.(*ExpressionIdentifier)
		if t.Variable != nil {
			t.Variable.Used = true
		}
	}
	/*
		var  s string;
		s += "11111111";
	*/

	if left.Type == VariableTypeString {
		if right.Type != VariableTypeString || (e.Type != ExpressionTypePlusAssign) {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on string and '%s'",
				errMsgPrefix(e.Pos),
				e.OpName(),
				right.TypeString()))
		}
		return result
	}
	//number
	if e.Type == ExpressionTypePlusAssign ||
		e.Type == ExpressionTypeMinusAssign ||
		e.Type == ExpressionTypeMulAssign ||
		e.Type == ExpressionTypeDivAssign ||
		e.Type == ExpressionTypeModAssign {
		if left.Equal(errs, right) {
			return result
		}
		if left.IsInteger() && right.IsInteger() && bin.Right.IsLiteral() {
			bin.Right.ConvertToNumber(left.Type)
			return result
		}
		if left.IsFloat() && right.IsFloat() && bin.Right.IsLiteral() {
			bin.Right.ConvertToNumber(left.Type)
			return result
		}

	}
	if e.Type == ExpressionTypeAndAssign ||
		e.Type == ExpressionTypeOrAssign ||
		e.Type == ExpressionTypeXorAssign {
		if left.IsInteger() && left.Equal(errs, right) {
			return result
		}
	}
	if e.Type == ExpressionTypeLshAssign ||
		e.Type == ExpressionTypeRshAssign {
		if left.IsInteger() && right.IsInteger() {
			if right.Type == VariableTypeLong {
				bin.Right.ConvertToNumber(VariableTypeInt)
			}
			return result
		}
	}

	*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'",
		errMsgPrefix(e.Pos),
		e.OpName(),
		left.TypeString(),
		right.TypeString()))

	return result
}
