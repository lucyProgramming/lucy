package ast

import "fmt"

func (this *Expression) checkOpAssignExpression(block *Block, errs *[]error) (t *Type) {
	bin := this.Data.(*ExpressionBinary)
	if bin.Left.Type == ExpressionTypeList {
		list := bin.Left.Data.([]*Expression)
		if len(list) > 1 {
			*errs = append(*errs,
				fmt.Errorf("%s expect 1 expression on left",
					errMsgPrefix(this.Pos)))
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
	result.Pos = this.Pos
	if err := right.rightValueValid(); err != nil {
		*errs = append(*errs, err)
		return result
	}
	if bin.Left.Type == ExpressionTypeIdentifier &&
		this.IsStatementExpression == false {
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
			(this.Type == ExpressionTypePlusAssign) {
			return result
		}
	}
	//number
	if this.Type == ExpressionTypePlusAssign ||
		this.Type == ExpressionTypeMinusAssign ||
		this.Type == ExpressionTypeMulAssign ||
		this.Type == ExpressionTypeDivAssign ||
		this.Type == ExpressionTypeModAssign {
		if left.assignAble(errs, right) {
			return result
		}
		if left.isInteger() && right.isInteger() && bin.Right.isLiteral() {
			bin.Right.convertToNumberType(left.Type)
			return result
		}
		if left.isFloat() && right.isFloat() && bin.Right.isLiteral() {
			bin.Right.convertToNumberType(left.Type)
			return result
		}
	}
	if this.Type == ExpressionTypeAndAssign ||
		this.Type == ExpressionTypeOrAssign ||
		this.Type == ExpressionTypeXorAssign {
		if left.isInteger() && left.assignAble(errs, right) {
			return result
		}
	}
	if this.Type == ExpressionTypeLshAssign ||
		this.Type == ExpressionTypeRshAssign {
		if left.isInteger() && right.isInteger() {
			if right.Type == VariableTypeLong {
				bin.Right.convertToNumberType(VariableTypeInt)
			}
			return result
		}
	}
	*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'",
		this.Pos.ErrMsgPrefix(),
		this.Op,
		left.TypeString(),
		right.TypeString()))

	return result
}
