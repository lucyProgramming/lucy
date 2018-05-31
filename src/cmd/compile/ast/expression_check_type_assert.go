package ast

import (
	"fmt"
)

/*
	in array type must eqaul
*/
func (e *Expression) checkTypeAssert(block *Block, errs *[]error) []*VariableType {
	assert := e.Data.(*ExpressionTypeAssert)
	object, es := assert.Expression.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}

	if object == nil {
		return nil
	}
	if object.IsPointer() == false {
		*errs = append(*errs, fmt.Errorf("%s expression is not pointer", errMsgPrefix(e.Pos)))
		return nil
	}
	err := assert.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	ret := make([]*VariableType, 2)
	ret[0] = &VariableType{}
	ret[0] = assert.Typ.Clone()
	ret[0].Pos = e.Pos
	ret[1] = &VariableType{}
	ret[1].Pos = e.Pos
	ret[1].Typ = VARIABLE_TYPE_BOOL // if  assert is ok
	if assert.Typ.validForTypeAssert() == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for type assertion",
			errMsgPrefix(e.Pos), assert.Typ.TypeString()))
		return ret
	}
	return ret
}
