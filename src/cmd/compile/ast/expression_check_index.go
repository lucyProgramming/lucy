package ast

import (
	"fmt"
)

func (e *Expression) checkIndexExpression(block *Block, errs *[]error) *Type {
	index := e.Data.(*ExpressionIndex)
	on, es := index.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if on == nil {
		return nil
	}
	switch on.Type {
	case VariableTypeArray, VariableTypeJavaArray:
		indexType, es := index.Index.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
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
		result := on.Array.Clone()
		result.Pos = e.Pos
		return result
	case VariableTypeMap:
		result := on.Map.V.Clone()
		result.Pos = e.Pos
		indexType, es := index.Index.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if indexType == nil {
			return result
		}
		if on.Map.K.Equal(errs, indexType) == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' for index",
				errMsgPrefix(e.Pos), indexType.TypeString(), on.Map.K.TypeString()))
		}
		return result
	default:
		*errs = append(*errs, fmt.Errorf("%s cannot index '%s'",
			errMsgPrefix(e.Pos), on.TypeString()))
		return nil
	}
}
