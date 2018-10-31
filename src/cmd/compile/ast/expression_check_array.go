package ast

import (
	"fmt"
)

/*
	in array type must equal
*/
func (e *Expression) checkArray(block *Block, errs *[]error) *Type {
	arr := e.Data.(*ExpressionArray)
	if arr.Type == nil && len(arr.Expressions) == 0 {
		*errs = append(*errs,
			fmt.Errorf("%s array literal has no type and no expression, cannot inference it`s type ",
				errMsgPrefix(e.Pos)))
		return nil
	}
	var noType = true
	if arr.Type != nil {
		noType = false
		err := arr.Type.resolve(block)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
	}

	for _, v := range arr.Expressions {
		eType, es := v.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if eType == nil {
			continue
		}
		if arr.Type != nil &&
			noType == false {
			convertExpressionToNeed(v, arr.Type.Array, eType)
			eType = v.Value
		}
		if noType && arr.Type == nil {
			if err := eType.isTyped(); err == nil {
				arr.Type = &Type{}
				arr.Type.Type = VariableTypeArray
				arr.Type.Array = eType.Clone()
				arr.Type.Pos = e.Pos
			} else {
				*errs = append(*errs, err)
			}
		}
		if arr.Type != nil {
			if arr.Type.Array.assignAble(errs, eType) == false {
				if noType {
					*errs = append(*errs, fmt.Errorf("%s array literal mix up '%s' and '%s'",
						errMsgPrefix(eType.Pos), arr.Type.Array.TypeString(), eType.TypeString()))
				} else {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(eType.Pos), eType.TypeString(), arr.Type.Array.TypeString()))
				}
			}
		}
	}
	if arr.Type == nil {
		return nil
	}
	result := arr.Type.Clone()
	result.Pos = e.Pos
	return result
}
