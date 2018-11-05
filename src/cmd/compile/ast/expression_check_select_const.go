package ast

import "fmt"

func (this *Expression) checkSelectConstExpression(block *Block, errs *[]error) *Type {
	selection := this.Data.(*ExpressionSelection)
	object, es := selection.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	if object.Type != VariableTypeClass {
		*errs = append(*errs, fmt.Errorf("%s not a class , but '%s'",
			object.Pos.ErrMsgPrefix(), object.TypeString()))
		return nil
	}
	if object.Class.Block.Constants == nil ||
		object.Class.Block.Constants[selection.Name] == nil {
		*errs = append(*errs, fmt.Errorf("%s const '%s' not found",
			this.Pos.ErrMsgPrefix(), selection.Name))
		return nil
	}
	c := object.Class.Block.Constants[selection.Name]
	this.fromConst(c)
	result := c.Type.Clone()
	result.Pos = this.Pos
	return result
}
