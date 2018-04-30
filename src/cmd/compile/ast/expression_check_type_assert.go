package ast

import (
	"fmt"
)

/*
	in array type must eqaul
*/
func (e *Expression) checkTypeAssert(block *Block, errs *[]error) []*VariableType {
	assert := e.Data.(*ExpressionTypeAssert)
	objectTs, es := assert.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	object, err := assert.Expression.mustBeOneValueContext(objectTs)
	if err != nil {
		*errs = append(*errs, err)
	}
	if object == nil {
		return nil
	}
	if object.IsPrimitive() {
		*errs = append(*errs, fmt.Errorf("%s expression is primitive", errMsgPrefix(e.Pos)))
		return nil
	}
	err = assert.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if assert.Typ.Typ != VARIABLE_TYPE_OBJECT {
		*errs = append(*errs, fmt.Errorf("%s type is not a object", errMsgPrefix(e.Pos)))
		return nil
	}
	ret := make([]*VariableType, 2)
	ret[0] = &VariableType{}
	ret[0].Pos = e.Pos
	ret[0].Typ = VARIABLE_TYPE_OBJECT
	ret[0].Class = assert.Typ.Class
	ret[1] = &VariableType{}
	ret[1].Pos = e.Pos
	ret[1].Typ = VARIABLE_TYPE_BOOL // if  assert is ok
	return ret
}
