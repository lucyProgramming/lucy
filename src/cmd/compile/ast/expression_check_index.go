package ast

import (
	"fmt"
)

func (e *Expression) checkIndexExpression(block *Block, errs *[]error) *Type {
	index := e.Data.(*ExpressionIndex)
	t, es := index.Expression.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	if t.Type != VariableTypeArray &&
		t.Type != VariableTypeMap &&
		t.Type != VariableTypeJavaArray {
		*errs = append(*errs, fmt.Errorf("%s cannot have 'index' on '%s'",
			errMsgPrefix(e.Pos), t.TypeString()))
		return nil
	}
	// array
	if t.Type == VariableTypeArray ||
		t.Type == VariableTypeJavaArray {
		indexType, es := index.Index.checkSingleValueContextExpression(block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if indexType != nil {
			if indexType.IsInteger() {
				if indexType.Type == VariableTypeLong {
					index.Index.ConvertToNumber(VariableTypeInt) //  convert to int
				}
			} else {
				*errs = append(*errs, fmt.Errorf("%s only integer can be used as index,but '%s'",
					errMsgPrefix(e.Pos), indexType.TypeString()))
			}
		}
		tt := t.Array.Clone()
		tt.Pos = e.Pos
		return tt
	} else {
		// map
		ret := t.Map.Value.Clone()
		ret.Pos = e.Pos
		indexType, es := index.Index.checkSingleValueContextExpression(block)
		if esNotEmpty(es) {
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
}
