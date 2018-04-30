package ast

import (
	"fmt"
)

func (e *Expression) getBinaryExpressionConstValue(f binaryConstFolder) (is bool, Typ int, Value interface{}, err error) {
	binary := e.Data.(*ExpressionBinary)
	is1, typ1, value1, err1, is2, typ2, value2, err2 := binary.getBinaryConstExpression()
	if err1 != nil { //something is wrong
		err = err1
		return
	}
	if err2 != nil {
		err = err2
		return
	}
	if is1 == false || is2 == false {
		is = false
		err = nil
		return
	}
	return f(is1, typ1, value1, is2, typ2, value2)
}

func (e *Expression) wrongOpErr(typ1, typ2 string) error {
	return fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
		errMsgPrefix(e.Pos),
		e.OpName(),
		typ1,
		typ2)
}

//byte -> int
func (e *Expression) typeWider(typ1, typ2 int, value1, value2 interface{}) (t int, v1 interface{}, v2 interface{}, err error) { //
	if typ1 == typ2 {
		return typ1, value1, value2, nil
	}
	if typ1 > typ2 {
		t = typ1
	} else {
		t = typ2
	}
	if t == typ1 { //typ1 has is wider
		v2, err = e.typeConvertor(t, typ2, value2)
		v1 = value1
	} else {
		v1, err = e.typeConvertor(t, typ1, value1)
		v2 = value2
	}
	if err == nil {
		return
	}
	return t, value1, value2, err
}

func (e *Expression) typeConvertor(target int, origin int, v interface{}) (interface{}, error) {
	switch target {
	case EXPRESSION_TYPE_SHORT:
		if origin == EXPRESSION_TYPE_BYTE {
			return int32(v.(byte)), nil
		}
	case EXPRESSION_TYPE_INT:
		if origin == EXPRESSION_TYPE_BYTE {
			return int32(v.(byte)), nil
		} else if origin == EXPRESSION_TYPE_SHORT {
			// short
		}
	case EXPRESSION_TYPE_LONG:
		if origin == EXPRESSION_TYPE_BYTE {
			return int64(v.(byte)), nil
		} else if origin == EXPRESSION_TYPE_SHORT || origin == EXPRESSION_TYPE_INT {
			// short
			return int64(v.(int32)), nil
		}
	case EXPRESSION_TYPE_FLOAT:
		if origin == EXPRESSION_TYPE_BYTE {
			return float32(v.(byte)), nil
		} else if origin == EXPRESSION_TYPE_SHORT || origin == EXPRESSION_TYPE_INT {
			// short
			return float32(v.(int32)), nil
		} else if origin == EXPRESSION_TYPE_LONG {
			return float32(v.(int64)), nil
		}
	case EXPRESSION_TYPE_DOUBLE:
		if origin == EXPRESSION_TYPE_BYTE {
			return float64(v.(byte)), nil
		} else if origin == EXPRESSION_TYPE_SHORT || origin == EXPRESSION_TYPE_INT {
			// short
			return float64(v.(int32)), nil
		} else if origin == EXPRESSION_TYPE_LONG {
			return float64(v.(int64)), nil
		} else if origin == EXPRESSION_TYPE_FLOAT {
			return float64(v.(float32)), nil
		}
	}
	return nil, e.wrongOpErr(e.OpName(origin), e.OpName(target))
}

func (e *Expression) getConstValue() (is bool, Typ int, Value interface{}, err error) {
	if e.IsLiteral() {
		return true, e.Typ, e.Data, nil
	}
	// !
	if e.Typ == EXPRESSION_TYPE_NOT {
		ee := e.Data.(*Expression)
		is, Typ, Value, err = ee.getConstValue()
		if err != nil || is == false {
			return
		}
		if Typ != EXPRESSION_TYPE_BOOL {
			err = fmt.Errorf("%s cannot apply '!' on a non-bool expression", errMsgPrefix(e.Pos))
			return
		}
		Value = !Value.(bool)
		return
	}
	if e.Typ == EXPRESSION_TYPE_NEGATIVE {
		ee := e.Data.(*Expression)
		is, Typ, Value, err = ee.getConstValue()
		if err != nil || is == false {
			return
		}
		if e.IsNumber(Typ) == false {
			is = false
			err = fmt.Errorf("%s cannot apply '-' on '%s'", errMsgPrefix(e.Pos), e.OpName(Typ))
			return
		}
		if Typ == EXPRESSION_TYPE_BYTE {
			is = false
			err = fmt.Errorf("%s cannot apply '-' on 'byte'", errMsgPrefix(e.Pos))
			return
		}
		switch Typ {
		case EXPRESSION_TYPE_BYTE:
			Value = -Value.(byte)
		case EXPRESSION_TYPE_SHORT:
			Value = -Value.(int32)
		case EXPRESSION_TYPE_INT:
			Value = -Value.(int32)
		case EXPRESSION_TYPE_LONG:
			Value = -Value.(int64)
		case EXPRESSION_TYPE_FLOAT:
			Value = -Value.(float32)
		case EXPRESSION_TYPE_DOUBLE:
			Value = -Value.(float64)
		}
		return
	}
	// && and ||
	if e.Typ == EXPRESSION_TYPE_LOGICAL_AND || e.Typ == EXPRESSION_TYPE_LOGICAL_OR {
		f := func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int,
			value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if typ1 != EXPRESSION_TYPE_BOOL || typ2 != EXPRESSION_TYPE_BOOL {
				err = fmt.Errorf("%s logical operation must apply to bool expressions", errMsgPrefix(e.Pos))
				return
			}
			is = true
			Typ = EXPRESSION_TYPE_BOOL
			if e.Typ == EXPRESSION_TYPE_LOGICAL_AND {
				Value = value1.(bool) && value2.(bool)
			} else {
				Value = value1.(bool) || value2.(bool)
			}
			err = nil
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	// + - * / % algebra arithmetic
	if e.Typ == EXPRESSION_TYPE_ADD ||
		e.Typ == EXPRESSION_TYPE_SUB ||
		e.Typ == EXPRESSION_TYPE_MUL ||
		e.Typ == EXPRESSION_TYPE_DIV ||
		e.Typ == EXPRESSION_TYPE_MOD {
		f := func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int,
			value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if typ1 == EXPRESSION_TYPE_STRING || typ2 == EXPRESSION_TYPE_STRING {
				if e.Typ != EXPRESSION_TYPE_ADD {
					err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
					return
				}
				if typ1 == EXPRESSION_TYPE_STRING || typ2 == EXPRESSION_TYPE_STRING {
					if typ1 != EXPRESSION_TYPE_STRING || typ2 != EXPRESSION_TYPE_STRING {
						err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
					}
					return
				}
				is = true
				Value = value1.(string) + value2.(string)
				return
			} else if e.IsNumber(typ1) || e.IsNumber(typ2) {
				if e.IsNumber(typ1) == false || e.IsNumber(typ2) == false {
					err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
				}
				Typ, value1, value2, err = e.typeWider(typ1, typ2, value1, value2)
				if err != nil {
					return
				}
				Value, err = e.numberTypeAlgebra(Typ, value1, value2)
				if err == nil {
					is = true
				} else {
					is = false
				}
				return
			}
			err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	// <<  >>
	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		f := func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int,
			value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if typ2 != EXPRESSION_TYPE_INT || e.IsNumber(typ1) == false {
				err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
				return
			}
			switch typ1 {
			case EXPRESSION_TYPE_BYTE:
				if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
					Value = value1.(byte) << uint64(value2.(byte))
				} else {
					Value = value1.(byte) >> uint64(value2.(byte))
				}
			case EXPRESSION_TYPE_INT:
				if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
					Value = value1.(int32) << uint64(value2.(int32))
				} else {
					Value = value1.(int32) >> uint64(value2.(int32))
				}
			case EXPRESSION_TYPE_LONG:
				if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
					Value = value1.(int64) << uint64(value2.(int64))
				} else {
					Value = value1.(int64) >> uint64(value2.(int64))
				}
			case EXPRESSION_TYPE_FLOAT:
				err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
				return
			}
			Typ = typ1
			is = true
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	// & |
	if e.Typ == EXPRESSION_TYPE_AND || e.Typ == EXPRESSION_TYPE_OR {
		f := func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int,
			value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if e.IsNumber(typ1) == false || e.IsNumber(typ2) == false {
				err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
				return
			}
			Typ, value1, value2, err = e.typeWider(typ1, typ2, value1, value2)
			if err != nil {
				return
			}
			switch Typ {
			case EXPRESSION_TYPE_BYTE:
				if EXPRESSION_TYPE_AND == e.Typ {
					e.Data = value1.(byte) & value2.(byte)
				} else {
					e.Data = value1.(byte) | value2.(byte)
				}
			case EXPRESSION_TYPE_INT:
				if EXPRESSION_TYPE_AND == e.Typ {
					e.Data = value1.(int64) & value2.(int64)
				} else {
					e.Data = value1.(int64) | value2.(int64)
				}
			default:
				err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
				return
			}
			is = true
			Typ = typ1
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	if e.Typ == EXPRESSION_TYPE_NOT {
		t := e.Data.(*Expression)
		is, Typ, Value, err = t.getConstValue()
		if err != nil {
			return
		}
		if is == false {
			return
		}
		if Typ != EXPRESSION_TYPE_BOOL {
			err = fmt.Errorf("!(not) can only apply to bool expression")
		} else {
			is = true
			Value = !Value.(bool)
			Typ = EXPRESSION_TYPE_BOOL
		}
		return
	}
	//  == != > < >= <=
	if e.Typ == EXPRESSION_TYPE_EQ ||
		e.Typ == EXPRESSION_TYPE_NE ||
		e.Typ == EXPRESSION_TYPE_GE ||
		e.Typ == EXPRESSION_TYPE_GT ||
		e.Typ == EXPRESSION_TYPE_LE ||
		e.Typ == EXPRESSION_TYPE_LE {
		f := func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int,
			value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if is1 == false || is2 == false {
				is = false
				return
			}
			if typ1 == EXPRESSION_TYPE_NULL || typ2 == EXPRESSION_TYPE_NULL {
				if (typ1 != VARIABLE_TYPE_NULL || typ2 != VARIABLE_TYPE_NULL) ||
					(e.Typ != EXPRESSION_TYPE_EQ || EXPRESSION_TYPE_NE != e.Typ) {
					err = e.wrongOpErr(e.OpName(typ1), e.OpName(typ2))
					return
				}
				is = true
				Typ = EXPRESSION_TYPE_BOOL
				Value = e.Typ == EXPRESSION_TYPE_EQ
				return
			}
			Typ, value1, value2, err = e.typeWider(typ1, typ2, value1, value2)
			if err != nil {
				err = fmt.Errorf("%s relation operation cannot apply to '%s' and '%s'",
					errMsgPrefix(e.Pos), e.OpName(typ1), e.OpName(typ2))
				return
			}
			b, er := e.relationCompare(Typ, value1, value2)
			if er != nil {
				err = er
				return
			}
			is = true
			Value = b
			err = nil
			Typ = EXPRESSION_TYPE_BOOL
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	is = false
	err = nil
	return
}

func (e *Expression) numberTypeAlgebra(typ int, value1, value2 interface{}) (value interface{}, err error) {
	switch typ {
	case EXPRESSION_TYPE_BYTE:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = value1.(byte) + value2.(byte)
		case EXPRESSION_TYPE_SUB:
			value = value1.(byte) - value2.(byte)
		case EXPRESSION_TYPE_MUL:
			value = value1.(byte) * value2.(byte)
		case EXPRESSION_TYPE_DIV:
			if value2.(byte) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(byte) / value2.(byte)
			}
		case EXPRESSION_TYPE_MOD:
			if value2.(byte) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(byte) % value2.(byte)
			}
		}
		return
	case EXPRESSION_TYPE_SHORT:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = value1.(int32) + value2.(int32)
		case EXPRESSION_TYPE_SUB:
			fmt.Println(value1.(int32), value2.(int32))
			value = value1.(int32) - value2.(int32)
		case EXPRESSION_TYPE_MUL:
			value = value1.(int32) * value2.(int32)
		case EXPRESSION_TYPE_DIV:
			if value2.(int32) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) / value2.(int32)
			}
		case EXPRESSION_TYPE_MOD:
			if value2.(int32) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) % value2.(int32)
			}
		}
		return
	case EXPRESSION_TYPE_INT:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = value1.(int32) + value2.(int32)
		case EXPRESSION_TYPE_SUB:
			fmt.Println(value1.(int32), value2.(int32))
			value = value1.(int32) - value2.(int32)
		case EXPRESSION_TYPE_MUL:
			value = value1.(int32) * value2.(int32)
		case EXPRESSION_TYPE_DIV:
			if value2.(int32) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) / value2.(int32)
			}
		case EXPRESSION_TYPE_MOD:
			if value2.(int32) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) % value2.(int32)
			}
		}
		return
	case EXPRESSION_TYPE_LONG:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = value1.(int64) + value2.(int64)
		case EXPRESSION_TYPE_SUB:
			fmt.Println(value1.(int64), value2.(int64))
			value = value1.(int64) - value2.(int64)
		case EXPRESSION_TYPE_MUL:
			value = value1.(int64) * value2.(int64)
		case EXPRESSION_TYPE_DIV:
			if value2.(int64) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(int64) / value2.(int64)
			}
		case EXPRESSION_TYPE_MOD:
			if value2.(int64) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(int64) % value2.(int64)
			}
		}
		return
	case EXPRESSION_TYPE_FLOAT:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = value1.(float32) + value2.(float32)
		case EXPRESSION_TYPE_SUB:
			value = value1.(float32) - value2.(float32)
		case EXPRESSION_TYPE_MUL:
			value = value1.(float32) * value2.(float32)
		case EXPRESSION_TYPE_DIV:
			if value2.(float32) == 0.0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(float32) / value2.(float32)
			}
		case EXPRESSION_TYPE_MOD:
			return nil, fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
				errMsgPrefix(e.Pos), e.OpName(), e.OpName(typ), e.OpName(typ))
		}
		return
	case EXPRESSION_TYPE_DOUBLE:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = value1.(float64) + value2.(float64)
		case EXPRESSION_TYPE_SUB:
			value = value1.(float64) - value2.(float64)
		case EXPRESSION_TYPE_MUL:
			value = value1.(float64) * value2.(float64)
		case EXPRESSION_TYPE_DIV:
			if value2.(float64) == 0.0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = value1.(float64) / value2.(float64)
			}
		case EXPRESSION_TYPE_MOD:
			return nil, fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
				errMsgPrefix(e.Pos), e.OpName(), e.OpName(typ), e.OpName(typ))
		}
		return
	}
	return nil, fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
		errMsgPrefix(e.Pos), e.OpName(), e.OpName(typ), e.OpName(typ))
}

func (e *Expression) relationCompare(typ int, value1, value2 interface{}) (b bool, err error) {
	switch typ {
	case EXPRESSION_TYPE_BOOL:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(bool) == value2.(bool)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(bool) != value2.(bool)
		} else {
			return false, fmt.Errorf("%s cannot apply '%s' on 'bool' and 'bool'",
				errMsgPrefix(e.Pos), e.OpName(typ))
		}
		return
	case EXPRESSION_TYPE_BYTE:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(byte) == value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(byte) != value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(byte) > value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(byte) >= value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(byte) < value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(byte) <= value2.(byte)
		}
		return
	case EXPRESSION_TYPE_INT:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(int32) == value2.(int32)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(int32) != value2.(int32)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(int32) > value2.(int32)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(int32) >= value2.(int32)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(int32) < value2.(int32)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(int32) <= value2.(int32)
		}
		return
	case EXPRESSION_TYPE_LONG:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(int64) == value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(int64) != value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(int64) > value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(int64) >= value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(int64) < value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(int64) <= value2.(int64)
		}
		return
	case EXPRESSION_TYPE_FLOAT:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(float32) == value2.(float32)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(float32) != value2.(float32)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(float32) > value2.(float32)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(float32) >= value2.(float32)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(float32) < value2.(float32)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(float32) <= value2.(float32)
		}
		return
	case EXPRESSION_TYPE_DOUBLE:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(float64) == value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(float64) != value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(float64) > value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(float64) >= value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(float64) < value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(float64) <= value2.(float64)
		}
		return
	case EXPRESSION_TYPE_STRING:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(string) == value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(string) != value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(string) > value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(string) >= value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(string) < value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(string) <= value2.(string)
		}
		return

	}
	return false, fmt.Errorf("%s can`t apply '%s' on '%s' and '%s'",
		errMsgPrefix(e.Pos), e.OpName(), e.OpName(typ), e.OpName(typ))
}
