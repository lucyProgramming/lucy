package ast

func (e *Expression) arithmeticBinayConstFolder(bin *ExpressionBinary) (is bool, err error) {
	if bin.Left.isInteger() && bin.Right.isInteger() {
		switch bin.Left.Typ {
		case EXPRESSION_TYPE_BYTE:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_BYTE, bin.Left.Data.(byte), bin.Right.Data.(byte))
				e.Typ = EXPRESSION_TYPE_BYTE
			case EXPRESSION_TYPE_SHORT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_SHORT, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
				e.Typ = EXPRESSION_TYPE_SHORT
			case EXPRESSION_TYPE_INT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_INT, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
				e.Typ = EXPRESSION_TYPE_INT
			case EXPRESSION_TYPE_LONG:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_LONG, int64(bin.Left.Data.(byte)), bin.Right.Data.(int64))
				e.Typ = EXPRESSION_TYPE_LONG
			}
			return
		case EXPRESSION_TYPE_SHORT:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_SHORT, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
				e.Typ = EXPRESSION_TYPE_SHORT

			case EXPRESSION_TYPE_SHORT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_SHORT, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Typ = EXPRESSION_TYPE_SHORT

			case EXPRESSION_TYPE_INT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Typ = EXPRESSION_TYPE_INT

			case EXPRESSION_TYPE_LONG:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_LONG, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
				e.Typ = EXPRESSION_TYPE_LONG
			}
			return
		case EXPRESSION_TYPE_INT:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
				e.Typ = EXPRESSION_TYPE_INT
			case EXPRESSION_TYPE_SHORT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Typ = EXPRESSION_TYPE_INT
			case EXPRESSION_TYPE_INT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Typ = EXPRESSION_TYPE_INT
			case EXPRESSION_TYPE_LONG:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_LONG, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
				e.Typ = EXPRESSION_TYPE_LONG
			}
			return
		case EXPRESSION_TYPE_LONG:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), int64(bin.Right.Data.(byte)))
				e.Typ = EXPRESSION_TYPE_LONG
			case EXPRESSION_TYPE_SHORT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
				e.Typ = EXPRESSION_TYPE_LONG
			case EXPRESSION_TYPE_INT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
				e.Typ = EXPRESSION_TYPE_LONG
			case EXPRESSION_TYPE_LONG:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), bin.Right.Data.(int64))
				e.Typ = EXPRESSION_TYPE_LONG
			}
			return
		}
		return
	}
	if bin.Left.isFloat() && bin.Right.isFloat() {
		switch bin.Left.Typ {
		case EXPRESSION_TYPE_FLOAT:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_FLOAT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_FLOAT, bin.Left.Data.(float32), bin.Right.Data.(float32))
				e.Typ = EXPRESSION_TYPE_FLOAT
				return
			case EXPRESSION_TYPE_DOUBLE:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_DOUBLE, float64(bin.Left.Data.(float32)), bin.Right.Data.(float64))
				e.Typ = EXPRESSION_TYPE_DOUBLE
				return
			}
		case EXPRESSION_TYPE_DOUBLE:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_FLOAT:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_DOUBLE, bin.Left.Data.(float64), float64(bin.Right.Data.(float32)))
				e.Typ = EXPRESSION_TYPE_DOUBLE
				return
			case EXPRESSION_TYPE_DOUBLE:
				e.Data, is, err = e.numberTypeAlgebra(EXPRESSION_TYPE_DOUBLE, bin.Left.Data.(float64), bin.Right.Data.(float64))
				e.Typ = EXPRESSION_TYPE_DOUBLE
				return
			}
		}
		return
	}

	return
}

func (e *Expression) relationBinayConstFolder(bin *ExpressionBinary) (is bool, err error) {
	// true == true  false == false
	if bin.Left.Typ == EXPRESSION_TYPE_BOOL &&
		bin.Right.Typ == EXPRESSION_TYPE_BOOL &&
		(e.Typ == EXPRESSION_TYPE_EQ || e.Typ == EXPRESSION_TYPE_NE) {
		e.Data, _ = e.relationCompare(EXPRESSION_TYPE_BOOL, bin.Left.Data.(bool), bin.Right.Data.(bool))
		is = true
		e.Typ = EXPRESSION_TYPE_BOOL
		return
	}
	// null == null or null != nil
	if bin.Left.Typ == EXPRESSION_TYPE_NULL && bin.Right.Typ == EXPRESSION_TYPE_NULL &&
		(e.Typ == EXPRESSION_TYPE_EQ || e.Typ == EXPRESSION_TYPE_NE) {
		e.Data = e.Typ == EXPRESSION_TYPE_EQ
		is = true
		e.Typ = EXPRESSION_TYPE_BOOL
		return
	}
	// string and string
	if bin.Left.Typ == EXPRESSION_TYPE_STRING && bin.Right.Typ == EXPRESSION_TYPE_STRING {
		is = true
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data, _ = e.relationCompare(EXPRESSION_TYPE_STRING, bin.Left.Data, bin.Right.Data)
		return
	}
	if bin.Left.isInteger() && bin.Right.isInteger() {
		switch bin.Left.Typ {
		case EXPRESSION_TYPE_BYTE:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_BYTE, bin.Left.Data.(byte), bin.Right.Data.(byte))
			case EXPRESSION_TYPE_SHORT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_SHORT, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
			case EXPRESSION_TYPE_INT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_INT, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
			case EXPRESSION_TYPE_LONG:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_LONG, int64(bin.Left.Data.(byte)), bin.Right.Data.(int64))
			}
			is = true
			e.Typ = EXPRESSION_TYPE_BOOL
			return
		case EXPRESSION_TYPE_SHORT:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_SHORT, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
			case EXPRESSION_TYPE_SHORT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_SHORT, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case EXPRESSION_TYPE_INT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case EXPRESSION_TYPE_LONG:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_LONG, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
			}
			is = true
			e.Typ = EXPRESSION_TYPE_BOOL
			return
		case EXPRESSION_TYPE_INT:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
			case EXPRESSION_TYPE_SHORT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case EXPRESSION_TYPE_INT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_INT, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case EXPRESSION_TYPE_LONG:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_LONG, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
			}
			is = true
			e.Typ = EXPRESSION_TYPE_BOOL
			return
		case EXPRESSION_TYPE_LONG:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_BYTE:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), int64(bin.Right.Data.(byte)))
			case EXPRESSION_TYPE_SHORT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
			case EXPRESSION_TYPE_INT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
			case EXPRESSION_TYPE_LONG:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_LONG, bin.Left.Data.(int64), bin.Right.Data.(int64))
			}
			is = true
			e.Typ = EXPRESSION_TYPE_BOOL
			return
		}
		return
	}
	if bin.Left.isFloat() && bin.Right.isFloat() {
		switch bin.Left.Typ {
		case EXPRESSION_TYPE_FLOAT:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_FLOAT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_FLOAT, bin.Left.Data.(float32), bin.Right.Data.(float32))
			case EXPRESSION_TYPE_DOUBLE:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_DOUBLE, float64(bin.Left.Data.(float32)), bin.Right.Data.(float64))
			}
			is = true
			e.Typ = EXPRESSION_TYPE_BOOL
			return
		case EXPRESSION_TYPE_DOUBLE:
			switch bin.Right.Typ {
			case EXPRESSION_TYPE_FLOAT:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_DOUBLE, bin.Left.Data.(float64), float64(bin.Right.Data.(float32)))
			case EXPRESSION_TYPE_DOUBLE:
				e.Data, _ = e.relationCompare(EXPRESSION_TYPE_DOUBLE, bin.Left.Data.(float64), bin.Right.Data.(float64))
			}
			is = true
			e.Typ = EXPRESSION_TYPE_BOOL
			return
		}
		return
	}

	return
}

func (e *Expression) numberTypeAlgebra(typ int, value1, value2 interface{}) (value interface{}, support bool, err error) {
	support = true
	switch typ {
	case EXPRESSION_TYPE_BYTE:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = byte(value1.(byte) + value2.(byte))
		case EXPRESSION_TYPE_SUB:
			value = byte(value1.(byte) - value2.(byte))
		case EXPRESSION_TYPE_MUL:
			value = byte(value1.(byte) * value2.(byte))
		case EXPRESSION_TYPE_DIV:
			if value2.(byte) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = byte(value1.(byte) / value2.(byte))
			}
		case EXPRESSION_TYPE_MOD:
			if value2.(byte) == 0 {
				err = devisionByZeroErr(e.Pos)
			} else {
				value = byte(value1.(byte) % value2.(byte))
			}
		}
		return
	case EXPRESSION_TYPE_SHORT:
		switch e.Typ {
		case EXPRESSION_TYPE_ADD:
			value = value1.(int32) + value2.(int32)
		case EXPRESSION_TYPE_SUB:
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
		}
		return
	}
	support = false
	return
}

func (e *Expression) relationCompare(typ int, value1, value2 interface{}) (b, support bool) {
	support = true
	switch typ {
	case EXPRESSION_TYPE_BOOL:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(bool) == value2.(bool)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(bool) != value2.(bool)
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
	support = false
	return
}
