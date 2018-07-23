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
	noType := true
	if arr.Type != nil {
		noType = false
		err := arr.Type.resolve(block)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
	}
	for _, v := range arr.Expressions {
		eTypes, es := v.check(block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if eTypes != nil {
			arr.Length += len(eTypes)
		}
		for _, eType := range eTypes {
			if eType == nil {
				continue
			}
			if noType && arr.Type == nil {
				if eType.RightValueValid() && eType.isTyped() {
					tmp := eType.Clone()
					tmp.Pos = e.Pos
					arr.Type = &Type{}
					arr.Type.Type = VariableTypeArray
					arr.Type.Array = tmp
					arr.Type.Pos = e.Pos
				} else {
					*errs = append(*errs, fmt.Errorf("%s right value '%s' untyped",
						errMsgPrefix(e.Pos), eType.TypeString()))
				}
			}
			if arr.Type != nil {
				if arr.Type.Array.Equal(errs, eType) == false {
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
	}
	if arr.Type == nil {
		return nil
	}
	result := arr.Type.Clone()
	result.Pos = e.Pos
	return result
}
