package ast

import (
	"fmt"
)

/*
	in array type must eqaul
*/
func (e *Expression) checkArray(block *Block, errs *[]error) *VariableType {
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
		if errorsNotEmpty(es) {
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
					arr.Type = &VariableType{}
					arr.Type.Type = VARIABLE_TYPE_ARRAY
					arr.Type.ArrayType = tt
					arr.Type.Pos = e.Pos
				} else {
					if t.RightValueValid() {
						*errs = append(*errs, fmt.Errorf("%s right value '%s' untyped",
							errMsgPrefix(e.Pos), t.TypeString()))
					} else {
						*errs = append(*errs, fmt.Errorf("%s right value '%s' invalid",
							errMsgPrefix(e.Pos), t.TypeString()))
					}
				}
			}
			if arr.Type != nil {
				if arr.Type.ArrayType.Equal(errs, t) == false {
					if noType {
						*errs = append(*errs, fmt.Errorf("%s array literal mix up '%s' and '%s'",
							errMsgPrefix(t.Pos), arr.Type.ArrayType.TypeString(), t.TypeString()))
					} else {
						*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
							errMsgPrefix(t.Pos), t.TypeString(), arr.Type.ArrayType.TypeString()))
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
