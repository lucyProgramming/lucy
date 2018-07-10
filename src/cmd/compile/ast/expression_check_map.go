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
				*errs = append(*errs, fmt.Errorf("%s k is not right value valid", errMsgPrefix(v.Left.Pos)))
			}
			if noType && m.Type.Map.Key == nil {
				if kType.isTyped() == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for k", errMsgPrefix(v.Left.Pos)))
				} else {
					m.Type.Map.Key = kType
					mapK = m.Type.Map.Key
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
		if esNotEmpty(es) {
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
		if noType && m.Type.Map.Value == nil {
			if vType.isTyped() == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for v",
					errMsgPrefix(v.Left.Pos)))
			} else {
				m.Type.Map.Value = vType
				mapV = m.Type.Map.Value
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

	if m.Type.Map.Key == nil {
		m.Type.Map.Key = &Type{
			Type: VariableTypeVoid,
			Pos:  e.Pos,
		}
	}
	if m.Type.Map.Value == nil {
		m.Type.Map.Value = &Type{
			Type: VariableTypeVoid,
			Pos:  e.Pos,
		}
	}
	return m.Type
}
