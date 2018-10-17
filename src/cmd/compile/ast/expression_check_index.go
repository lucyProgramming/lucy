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
	case VariableTypeArray,
		VariableTypeJavaArray:
		indexType, es := index.Index.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if indexType != nil {
			if indexType.isInteger() {
				if indexType.Type == VariableTypeLong {
					index.Index.convertToNumber(VariableTypeInt) //  convert to int
				}
			} else {
				*errs = append(*errs,
					fmt.Errorf("%s only integer can be used as index,but '%s'",
						index.Index.Pos.ErrMsgPrefix(), indexType.TypeString()))
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
		if on.Map.K.assignAble(errs, indexType) == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' for index",
				index.Index.Pos.ErrMsgPrefix(), indexType.TypeString(), on.Map.K.TypeString()))
		}
		return result
	case VariableTypeString:
		indexType, es := index.Index.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if indexType != nil {
			if indexType.isInteger() {
				if indexType.Type == VariableTypeLong {
					index.Index.convertToNumber(VariableTypeInt) //  convert to int
				}
			} else {
				*errs = append(*errs, fmt.Errorf("%s only integer can be used as index,but '%s'",
					index.Index.Pos.ErrMsgPrefix(), indexType.TypeString()))
			}
		}
		result := &Type{
			Type: VariableTypeByte,
			Pos:  e.Pos,
		}
		return result
	default:
		*errs = append(*errs, fmt.Errorf("%s cannot index '%s'",
			on.Pos.ErrMsgPrefix(), on.TypeString()))
		return nil
	}
}
