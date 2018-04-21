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
	var mapk *VariableType
	var mapv *VariableType
	noType := m.Typ == nil
	if noType && len(m.Values) == 0 {
		*errs = append(*errs, fmt.Errorf("%s map literal has no type, no initiational values,cannot inference it`s type ", errMsgPrefix(e.Pos)))
		goto end
	}
	if m.Typ == nil {
		m.Typ = &VariableType{}
		m.Typ.Typ = VARIABLE_TYPE_MAP
	}
	if m.Typ.Map == nil {
		m.Typ.Map = &Map{}
	}
	for k, v := range m.Values {
		// map k
		ktypes, es := v.Left.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		ktype, err := v.Left.mustBeOneValueContext(ktypes)
		if err != nil {
			*errs = append(*errs, err)
		}
		if ktype != nil {
			rightValueValid := ktype.RightValueValid()
			if false == rightValueValid {
				*errs = append(*errs, fmt.Errorf("%s k is not right value valid", errMsgPrefix(v.Left.Pos)))
			}
			if noType && k == 0 {
				if ktype.isTyped() == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for k", errMsgPrefix(v.Left.Pos)))
				} else {
					if noType && k == 0 {
						m.Typ.Map.K = ktype
						mapk = m.Typ.Map.K
					}
				}
			}
			if rightValueValid && mapk != nil {
				if mapk.Equal(ktype) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'", errMsgPrefix(v.Left.Pos),
						ktype.TypeString(), mapk.TypeString()))
				}
			}
		}
		// map v
		vtypes, es := v.Right.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		vtype, err := v.Right.mustBeOneValueContext(vtypes)
		if err != nil {
			*errs = append(*errs, err)
		}
		if vtype == nil {
			continue
		}
		if false == ktype.RightValueValid() {
			*errs = append(*errs, fmt.Errorf("%s k is not right value valid", errMsgPrefix(v.Left.Pos)))
			continue
		}
		if noType && k == 0 {
			if vtype.isTyped() == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use untyped value for v", errMsgPrefix(v.Left.Pos)))
			} else {
				if noType && k == 0 {
					m.Typ.Map.V = vtype
					mapv = m.Typ.Map.V
				}
			}
		}
		if mapv != nil {
			if mapv.Equal(vtype) == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'", errMsgPrefix(v.Right.Pos),
					vtype.TypeString(), mapv.TypeString()))
			}
		}
	}

end:
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
