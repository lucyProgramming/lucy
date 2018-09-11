package ast

import (
	"fmt"
)

/*
	in array type must equal
*/
func (e *Expression) checkTypeAssert(block *Block, errs *[]error) []*Type {
	assert := e.Data.(*ExpressionTypeAssert)
	object, es := assert.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	if object.RightValueValid() == false {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not right value valid",
			errMsgPrefix(object.Pos), object.TypeString()))
		return nil
	}
	if object.IsPointer() == false {
		*errs = append(*errs, fmt.Errorf("%s expression is not pointer", errMsgPrefix(object.Pos)))
		return nil
	}
	err := assert.Type.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	result := make([]*Type, 2)
	result[0] = assert.Type.Clone()
	result[0].Pos = e.Pos
	result[1] = &Type{}
	result[1].Pos = e.Pos
	result[1].Type = VariableTypeBool // if  assert is ok
	if assert.Type.validForTypeAssertOrConversion() == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for type assertion",
			errMsgPrefix(assert.Type.Pos), assert.Type.TypeString()))
		return result
	}
	return result
}
