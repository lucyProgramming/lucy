package ast

import (
	"fmt"
)

/*
	in array type must eqaul
*/
func (e *Expression) checkArray(block *Block, errs *[]error) *VariableType {
	arr := e.Data.(*ExpressionArrayLiteral)
	if arr.Typ == nil && len(arr.Expressions) == 0 {
		*errs = append(*errs, fmt.Errorf("%s array literal has no type, no expression, cannot inference it`s type ",
			errMsgPrefix(e.Pos)))
		return nil
	}
	noType := true
	if arr.Typ != nil {
		noType = false
		err := arr.Typ.resolve(block)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
	}
	for _, v := range arr.Expressions {
		ts, es := v.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if ts != nil {
			arr.Length += len(ts)
		}
		for _, t := range ts {
			if t == nil {
				continue
			}
			if noType && arr.Typ == nil {
				if t.RightValueValid() && t.isTyped() {
					tt := t.Clone()
					tt.Pos = e.Pos
					arr.Typ = &VariableType{}
					arr.Typ.Typ = VARIABLE_TYPE_ARRAY
					arr.Typ.ArrayType = tt
					arr.Typ.Pos = e.Pos
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
			if arr.Typ != nil {
				if arr.Typ.ArrayType.Equal(errs, t) == false {
					if noType {
						*errs = append(*errs, fmt.Errorf("%s array literal mix up '%s' and '%s'",
							errMsgPrefix(t.Pos), arr.Typ.ArrayType.TypeString(), t.TypeString()))
					} else {
						*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
							errMsgPrefix(t.Pos), t.TypeString(), arr.Typ.ArrayType.TypeString()))
					}
				}
			}
		}
	}
	if arr.Typ == nil {
		return nil
	}
	tt := arr.Typ.Clone()
	tt.Pos = e.Pos
	return tt
}
