package ast

import "fmt"

func (e *Expression) checkOpAssignExpression(block *Block, errs *[]error) (t *Type) {
	bin := e.Data.(*ExpressionBinary)
	if bin.Left.Type == ExpressionTypeList {
		list := bin.Left.Data.([]*Expression)
		if len(list) > 1 {
			*errs = append(*errs,
				fmt.Errorf("%s expect 1 expression on left",
					errMsgPrefix(e.Pos)))
		}
		bin.Left = list[0]
	}

	left := bin.Left.getLeftValue(block, errs)
	right, es := bin.Right.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if left == nil || right == nil {
		return
	}
	result := left.Clone()
	result.Pos = e.Pos
	if err := right.rightValueValid(); err != nil {
		*errs = append(*errs, err)
		return result
	}
	if bin.Left.Type == ExpressionTypeIdentifier &&
		e.IsStatementExpression == false {
		/*
			var a = 1
			print(a += 1)
		*/
		t := bin.Left.Data.(*ExpressionIdentifier)
		if t.Variable != nil {
			t.Variable.Used = true
		}
	}
	convertExpressionToNeed(bin.Right, left, right)
	/*
		var  s string
		s += "11111111"
	*/
	if left.Type == VariableTypeString {
		if right.Type == VariableTypeString &&
			(e.Type == ExpressionTypePlusAssign) {
			return result
		}
	}
	//number
	if e.Type == ExpressionTypePlusAssign ||
		e.Type == ExpressionTypeMinusAssign ||
		e.Type == ExpressionTypeMulAssign ||
		e.Type == ExpressionTypeDivAssign ||
		e.Type == ExpressionTypeModAssign {
		if left.assignAble(errs, right) {
			return result
		}
		if left.isInteger() && right.isInteger() && bin.Right.isLiteral() {
			bin.Right.convertToNumber(left.Type)
			return result
		}
		if left.isFloat() && right.isFloat() && bin.Right.isLiteral() {
			bin.Right.convertToNumber(left.Type)
			return result
		}
	}
	if e.Type == ExpressionTypeAndAssign ||
		e.Type == ExpressionTypeOrAssign ||
		e.Type == ExpressionTypeXorAssign {
		if left.isInteger() && left.assignAble(errs, right) {
			return result
		}
	}
	if e.Type == ExpressionTypeLshAssign ||
		e.Type == ExpressionTypeRshAssign {
		if left.isInteger() && right.isInteger() {
			if right.Type == VariableTypeLong {
				bin.Right.convertToNumber(VariableTypeInt)
			}
			return result
		}
	}
	*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'",
		e.Pos.ErrMsgPrefix(),
		e.Op,
		left.TypeString(),
		right.TypeString()))

	return result
}
