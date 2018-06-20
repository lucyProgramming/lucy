package ast

import (
	"fmt"
)

func (e *Expression) checkIndexExpression(block *Block, errs *[]error) *Type {
	index := e.Data.(*ExpressionIndex)
	t, es := index.Expression.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}

	if t == nil {
		return nil
	}
	if t.Type != VARIABLE_TYPE_ARRAY &&
		t.Type != VARIABLE_TYPE_MAP &&
		t.Type != VARIABLE_TYPE_JAVA_ARRAY {
		*errs = append(*errs, fmt.Errorf("%s cannot have 'index' on '%s'",
			errMsgPrefix(e.Pos), t.TypeString()))
		return nil
	}
	// array
	if t.Type == VARIABLE_TYPE_ARRAY ||
		t.Type == VARIABLE_TYPE_JAVA_ARRAY {
		indexType, es := index.Index.checkSingleValueContextExpression(block)
		if errorsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if indexType != nil {
			if indexType.IsInteger() {
				if indexType.Type == VARIABLE_TYPE_LONG {
					index.Index.ConvertToNumber(VARIABLE_TYPE_INT) //  convert to int
				}
			} else {
				*errs = append(*errs, fmt.Errorf("%s only integer can be used as index,but '%s'",
					errMsgPrefix(e.Pos), indexType.TypeString()))
			}
		}
		tt := t.ArrayType.Clone()
		tt.Pos = e.Pos
		return tt
	}
	// map
	ret := t.Map.Value.Clone()
	ret.Pos = e.Pos
	indexType, es := index.Index.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}

	if indexType == nil {
		return ret
	}
	if t.Map.Key.Equal(errs, indexType) == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' for index",
			errMsgPrefix(e.Pos), indexType.TypeString(), t.Map.Key.TypeString()))
	}
	return ret

}
