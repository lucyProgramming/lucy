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
		*errs = append(*errs, fmt.Errorf("%s array literal has no type, no expression, cannot inference it`s type ",
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
		ts, es := v.check(block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if ts != nil {
			arr.Length += len(ts)
		}
		for _, t := range ts {
			if t == nil {
				continue
			}
			if noType && arr.Type == nil {
				if t.RightValueValid() && t.isTyped() {
					tt := t.Clone()
					tt.Pos = e.Pos
					arr.Type = &Type{}
					arr.Type.Type = VariableTypeArray
					arr.Type.Array = tt
					arr.Type.Pos = e.Pos
				} else {
					*errs = append(*errs, fmt.Errorf("%s right value '%s' untyped",
						errMsgPrefix(e.Pos), t.TypeString()))
				}
			}
			if arr.Type != nil {
				if arr.Type.Array.Equal(errs, t) == false {
					if noType {
						*errs = append(*errs, fmt.Errorf("%s array literal mix up '%s' and '%s'",
							errMsgPrefix(t.Pos), arr.Type.Array.TypeString(), t.TypeString()))
					} else {
						*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
							errMsgPrefix(t.Pos), t.TypeString(), arr.Type.Array.TypeString()))
					}
				}
			}
		}
	}
	if arr.Type == nil {
		return nil
	}
	tt := arr.Type.Clone()
	tt.Pos = e.Pos
	return tt
}
