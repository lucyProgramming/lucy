package ast

import (
	"errors"
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
			fmt.Errorf("%s map literal has no type and no initiational values,"+
				"cannot inference it`s type",
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
	longMap := make(map[int64]*Pos)
	floatMap := make(map[float32]*Pos)
	doubleMap := make(map[float64]*Pos)
	stringMap := make(map[string]*Pos)
	for _, v := range m.KeyValuePairs {
		// map k
		kType, es := v.Key.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if kType != nil {
			if err := kType.rightValueValid(); err != nil {
				*errs = append(*errs, err)
				continue
			}
			if noType && m.Type.Map.K == nil {
				if err := kType.isTyped(); err != nil {
					*errs = append(*errs, err)
				} else {
					m.Type.Map.K = kType
					mapK = m.Type.Map.K
				}
			}
			if mapK != nil {
				if mapK.assignAble(errs, kType) == false {
					if noType {
						*errs = append(*errs, fmt.Errorf("%s mix '%s' and '%s' for map value",
							errMsgPrefix(v.Key.Pos),
							kType.TypeString(), mapK.TypeString()))
					} else {
						*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
							errMsgPrefix(v.Key.Pos),
							kType.TypeString(), mapK.TypeString()))
					}
				}
			}
		}
		if m.Type.Map.K != nil &&
			v.Key.isLiteral() &&
			m.Type.Map.K.Type == v.Key.Value.Type {
			errMsg := func(pos *Pos, first *Pos, which interface{}) error {
				errMsg := fmt.Sprintf("%s  '%v' duplicate key,first declared at:\n",
					pos.ErrMsgPrefix(), which)
				errMsg += fmt.Sprintf("\t%s", errMsgPrefix(first))
				return errors.New(errMsg)
			}
			switch m.Type.Map.K.Type {
			case VariableTypeByte:
				fallthrough
			case VariableTypeChar:
				fallthrough
			case VariableTypeShort:
				fallthrough
			case VariableTypeInt:
				fallthrough
			case VariableTypeLong:
				value := v.Key.Data.(int64)
				if first, ok := longMap[value]; ok {
					*errs = append(*errs, errMsg(v.Key.Pos, first, v.Key.Data))
				} else {
					longMap[value] = v.Key.Pos
				}
			case VariableTypeFloat:
				value := v.Key.Data.(float32)
				if first, ok := floatMap[value]; ok {
					*errs = append(*errs, errMsg(v.Key.Pos, first, v.Key.Data))
				} else {
					floatMap[value] = v.Key.Pos
				}
			case VariableTypeDouble:
				value := v.Key.Data.(float64)
				if first, ok := doubleMap[value]; ok {
					*errs = append(*errs, errMsg(v.Key.Pos, first, v.Key.Data))
				} else {
					doubleMap[value] = v.Key.Pos
				}
			case VariableTypeString:
				value := v.Key.Data.(string)
				if first, ok := stringMap[value]; ok {
					*errs = append(*errs, errMsg(v.Key.Pos, first, v.Key.Data))
				} else {
					stringMap[value] = v.Key.Pos
				}
			}
		}
		// map v
		vType, es := v.Value.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if vType == nil {
			continue
		}
		if err := kType.rightValueValid(); err != nil {
			*errs = append(*errs, err)
			continue
		}
		if noType && m.Type.Map.V == nil {
			if err := vType.isTyped(); err != nil {
				*errs = append(*errs, err)
			} else {
				m.Type.Map.V = vType
				mapV = m.Type.Map.V
			}
		}
		if mapV != nil {
			if mapV.assignAble(errs, vType) == false {
				if noType {
					*errs = append(*errs, fmt.Errorf("%s mix '%s' and '%s' for map key",
						errMsgPrefix(v.Value.Pos),
						vType.TypeString(), mapV.TypeString()))
				} else {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(v.Value.Pos),
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
