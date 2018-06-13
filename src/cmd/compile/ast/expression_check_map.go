package ast

import (
	"fmt"
)

func (e *Expression) checkMapExpression(block *Block, errs *[]error) *VariableType {
	m := e.Data.(*ExpressionMap)
	if m.Typ != nil {
		if err := m.Typ.resolve(block); err != nil {
			*errs = append(*errs, err)
		}
	}
	var mapK *VariableType
	var mapV *VariableType
	noType := m.Typ == nil
	if noType && len(m.KeyValuePairs) == 0 {
		*errs = append(*errs, fmt.Errorf("%s map literal has no type, no initiational values,cannot inference it`s type",
			errMsgPrefix(e.Pos)))
		return nil
	}
	if m.Typ == nil {
		m.Typ = &VariableType{}
		m.Typ.Pos = e.Pos
		m.Typ.Typ = VARIABLE_TYPE_MAP
	}
	if m.Typ.Map == nil {
		m.Typ.Map = &Map{}
	}
	for _, v := range m.KeyValuePairs {
		// map k
		kType, es := v.Left.checkSingleValueContextExpression(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if kType != nil {
			rightValueValid := kType.RightValueValid()
			if false == rightValueValid {
				*errs = append(*errs, fmt.Errorf("%s k is not right value valid", errMsgPrefix(v.Left.Pos)))
			}
			if noType && m.Typ.Map.K == nil {
				if kType.isTyped() == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for k", errMsgPrefix(v.Left.Pos)))
				} else {
					m.Typ.Map.K = kType
					mapK = m.Typ.Map.K
				}
			}
			if rightValueValid && mapK != nil {
				if mapK.Equal(errs, kType) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'", errMsgPrefix(v.Left.Pos),
						kType.TypeString(), mapK.TypeString()))
				}
			}
		}
		// map v
		vType, es := v.Right.checkSingleValueContextExpression(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if vType == nil {
			continue
		}
		if false == kType.RightValueValid() {
			*errs = append(*errs, fmt.Errorf("%s k is not right value valid",
				errMsgPrefix(v.Left.Pos)))
			continue
		}
		if noType && m.Typ.Map.V == nil {
			if vType.isTyped() == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for v",
					errMsgPrefix(v.Left.Pos)))
			} else {
				m.Typ.Map.V = vType
				mapV = m.Typ.Map.V
			}
		}
		if mapV != nil {
			if mapV.Equal(errs, vType) == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix(v.Right.Pos),
					vType.TypeString(), mapV.TypeString()))
			}
		}
	}

	if m.Typ.Map.K == nil {
		m.Typ.Map.K = &VariableType{
			Typ: VARIABLE_TYPE_VOID,
			Pos: e.Pos,
		}
	}
	if m.Typ.Map.V == nil {
		m.Typ.Map.V = &VariableType{
			Typ: VARIABLE_TYPE_VOID,
			Pos: e.Pos,
		}
	}
	return m.Typ
}
