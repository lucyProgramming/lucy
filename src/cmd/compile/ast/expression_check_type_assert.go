package ast

import (
	"fmt"
)

/*
	in array type must equal
*/
func (this *Expression) checkTypeAssert(block *Block, errs *[]error) []*Type {
	assert := this.Data.(*ExpressionTypeAssert)
	object, es := assert.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	if err := object.rightValueValid(); err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if object.IsPointer() == false {
		*errs = append(*errs,
			fmt.Errorf("%s expression is not pointer",
				errMsgPrefix(object.Pos)))
		return nil
	}
	err := assert.Type.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if assert.Type.validForTypeAssertOrConversion() == false {
		*errs = append(*errs,
			fmt.Errorf("%s cannot use '%s' for type assertion",
				errMsgPrefix(assert.Type.Pos), assert.Type.TypeString()))
		return nil
	}
	var result []*Type
	if len(this.Lefts) > 1 {
		assert.MultiValueContext = true
		result = make([]*Type, 2)
		result[0] = assert.Type.Clone()
		result[0].Pos = this.Pos
		result[1] = &Type{}
		result[1].Pos = this.Pos
		result[1].Type = VariableTypeBool // if  assert is ok
	} else {
		result = make([]*Type, 1)
		result[0] = assert.Type.Clone()
		result[0].Pos = this.Pos
	}
	return result
}
