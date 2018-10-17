package ast

import "fmt"

func (e *Expression) checkSelectConstExpression(block *Block, errs *[]error) *Type {
	selection := e.Data.(*ExpressionSelection)
	object, es := selection.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	if object.Type != VariableTypeClass {
		*errs = append(*errs, fmt.Errorf("%s not a class , but '%s'",
			e.Pos.ErrMsgPrefix(), object.TypeString()))
		return nil
	}
	if object.Class.Block.Constants == nil ||
		object.Class.Block.Constants[selection.Name] == nil {
		*errs = append(*errs, fmt.Errorf("%s const '%s' not found",
			e.Pos.ErrMsgPrefix(), selection.Name))
		return nil
	}
	c := object.Class.Block.Constants[selection.Name]
	e.fromConst(c)
	result := c.Type.Clone()
	result.Pos = e.Pos
	return result
}
