package ast

import (
	"fmt"
)

func (e *Expression) checkMapExpression(block *Block, errs *[]error) *Type {
	m := e.Data.(*ExpressionMap)
	if m.Type != nil {
		if err := m.Type.resolve(block); err != nil {
			*errs = append(*errs, err)
		}
	}
	var mapK *Type
	var mapV *Type
	noType := m.Type == nil
	if noType && len(m.KeyValuePairs) == 0 {
		*errs = append(*errs,
			fmt.Errorf("%s map literal has no type and no initiational values,cannot inference it`s type",
				errMsgPrefix(e.Pos)))
		return nil
	}
	if m.Type == nil {
		m.Type = &Type{}
		m.Type.Pos = e.Pos
		m.Type.Type = VariableTypeMap
	}
	if m.Type.Map == nil {
		m.Type.Map = &Map{}
	}
	for _, v := range m.KeyValuePairs {
		// map k
		kType, es := v.Left.checkSingleValueContextExpression(block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if kType != nil {
			rightValueValid := kType.RightValueValid()
			if false == rightValueValid {
				*errs = append(*errs, fmt.Errorf("%s k is not right value valid",
					errMsgPrefix(v.Left.Pos)))
			}
			if noType && m.Type.Map.K == nil {
				if kType.isTyped() == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for k",
						errMsgPrefix(v.Left.Pos)))
				} else {
					m.Type.Map.K = kType
					mapK = m.Type.Map.K
				}
			}
			if rightValueValid && mapK != nil {
				if mapK.Equal(errs, kType) == false {
					if noType {
						*errs = append(*errs, fmt.Errorf("%s mix '%s' and '%s' for map value",
							errMsgPrefix(v.Left.Pos),
							kType.TypeString(), mapK.TypeString()))
					} else {
						*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
							errMsgPrefix(v.Left.Pos),
							kType.TypeString(), mapK.TypeString()))
					}
				}
			}
		}
		// map v
		vType, es := v.Right.checkSingleValueContextExpression(block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if vType == nil {
			continue
		}
		if false == kType.RightValueValid() {
			*errs = append(*errs, fmt.Errorf("%s k is not right value valid",
				errMsgPrefix(v.Right.Pos)))
			continue
		}
		if noType && m.Type.Map.V == nil {
			if vType.isTyped() == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for v",
					errMsgPrefix(v.Right.Pos)))
			} else {
				m.Type.Map.V = vType
				mapV = m.Type.Map.V
			}
		}
		if mapV != nil {
			if mapV.Equal(errs, vType) == false {
				if noType {
					*errs = append(*errs, fmt.Errorf("%s mix '%s' and '%s' for map key",
						errMsgPrefix(v.Right.Pos),
						vType.TypeString(), mapV.TypeString()))
				} else {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(v.Right.Pos),
						vType.TypeString(), mapV.TypeString()))
				}
			}
		}
	}
	if m.Type.Map.K == nil {
		m.Type.Map.K = &Type{
			Type: VariableTypeVoid,
			Pos:  e.Pos,
		}
	}
	if m.Type.Map.V == nil {
		m.Type.Map.V = &Type{
			Type: VariableTypeVoid,
			Pos:  e.Pos,
		}
	}
	return m.Type
}
