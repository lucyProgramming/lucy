package ast

func (e *Expression) arithmeticBinaryConstFolder(bin *ExpressionBinary) (is bool, err error) {
	if bin.Left.isInteger() && bin.Right.isInteger() {
		switch bin.Left.Type {
		case ExpressionTypeByte:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeByte, bin.Left.Data.(byte), bin.Right.Data.(byte))
				e.Type = ExpressionTypeByte
			case ExpressionTypeShort:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeShort, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
				e.Type = ExpressionTypeShort
			case ExpressionTypeInt:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeInt, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
				e.Type = ExpressionTypeInt
			case ExpressionTypeLong:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeLong, int64(bin.Left.Data.(byte)), bin.Right.Data.(int64))
				e.Type = ExpressionTypeLong
			}
			return
		case ExpressionTypeShort:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeShort, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
				e.Type = ExpressionTypeShort
			case ExpressionTypeShort:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeShort, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Type = ExpressionTypeShort
			case ExpressionTypeInt:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeInt, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Type = ExpressionTypeInt
			case ExpressionTypeLong:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeLong, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
				e.Type = ExpressionTypeLong
			}
			return
		case ExpressionTypeInt:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeInt, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
				e.Type = ExpressionTypeInt
			case ExpressionTypeShort:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeInt, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Type = ExpressionTypeInt
			case ExpressionTypeInt:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeInt, bin.Left.Data.(int32), bin.Right.Data.(int32))
				e.Type = ExpressionTypeInt
			case ExpressionTypeLong:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeLong, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
				e.Type = ExpressionTypeLong
			}
			return
		case ExpressionTypeLong:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeLong, bin.Left.Data.(int64), int64(bin.Right.Data.(byte)))
				e.Type = ExpressionTypeLong
			case ExpressionTypeShort:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeLong, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
				e.Type = ExpressionTypeLong
			case ExpressionTypeInt:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeLong, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
				e.Type = ExpressionTypeLong
			case ExpressionTypeLong:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeLong, bin.Left.Data.(int64), bin.Right.Data.(int64))
				e.Type = ExpressionTypeLong
			}
			return
		}
		return
	}
	if bin.Left.isFloat() && bin.Right.isFloat() {
		switch bin.Left.Type {
		case ExpressionTypeFloat:
			switch bin.Right.Type {
			case ExpressionTypeFloat:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeFloat, bin.Left.Data.(float32), bin.Right.Data.(float32))
				e.Type = ExpressionTypeFloat
				return
			case ExpressionTypeDouble:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeDouble, float64(bin.Left.Data.(float32)), bin.Right.Data.(float64))
				e.Type = ExpressionTypeDouble
				return
			}
		case ExpressionTypeDouble:
			switch bin.Right.Type {
			case ExpressionTypeFloat:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeDouble, bin.Left.Data.(float64), float64(bin.Right.Data.(float32)))
				e.Type = ExpressionTypeDouble
				return
			case ExpressionTypeDouble:
				e.Data, is, err = e.numberTypeAlgebra(ExpressionTypeDouble, bin.Left.Data.(float64), bin.Right.Data.(float64))
				e.Type = ExpressionTypeDouble
				return
			}
		}
		return
	}

	return
}

func (e *Expression) relationBinaryConstFolder(bin *ExpressionBinary) (is bool, err error) {
	// true == true  false == false
	if bin.Left.Type == ExpressionTypeBool &&
		bin.Right.Type == ExpressionTypeBool &&
		(e.Type == ExpressionTypeEq || e.Type == ExpressionTypeNe) {
		e.Data, _ = e.relationCompare(ExpressionTypeBool, bin.Left.Data.(bool), bin.Right.Data.(bool))
		is = true
		e.Type = ExpressionTypeBool
		return
	}
	// null == null or null != nil
	if bin.Left.Type == ExpressionTypeNull && bin.Right.Type == ExpressionTypeNull &&
		(e.Type == ExpressionTypeEq || e.Type == ExpressionTypeNe) {
		e.Data = e.Type == ExpressionTypeEq
		is = true
		e.Type = ExpressionTypeBool
		return
	}
	// string and string
	if bin.Left.Type == ExpressionTypeString && bin.Right.Type == ExpressionTypeString {
		is = true
		e.Type = ExpressionTypeBool
		e.Data, _ = e.relationCompare(ExpressionTypeString, bin.Left.Data, bin.Right.Data)
		return
	}
	if bin.Left.isInteger() && bin.Right.isInteger() {
		switch bin.Left.Type {
		case ExpressionTypeByte:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, _ = e.relationCompare(ExpressionTypeByte, bin.Left.Data.(byte), bin.Right.Data.(byte))
			case ExpressionTypeShort:
				e.Data, _ = e.relationCompare(ExpressionTypeShort, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
			case ExpressionTypeInt:
				e.Data, _ = e.relationCompare(ExpressionTypeInt, int32(bin.Left.Data.(byte)), bin.Right.Data.(int32))
			case ExpressionTypeLong:
				e.Data, _ = e.relationCompare(ExpressionTypeLong, int64(bin.Left.Data.(byte)), bin.Right.Data.(int64))
			}
			is = true
			e.Type = ExpressionTypeBool
			return
		case ExpressionTypeShort:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, _ = e.relationCompare(ExpressionTypeShort, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
			case ExpressionTypeShort:
				e.Data, _ = e.relationCompare(ExpressionTypeShort, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case ExpressionTypeInt:
				e.Data, _ = e.relationCompare(ExpressionTypeInt, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case ExpressionTypeLong:
				e.Data, _ = e.relationCompare(ExpressionTypeLong, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
			}
			is = true
			e.Type = ExpressionTypeBool
			return
		case ExpressionTypeInt:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, _ = e.relationCompare(ExpressionTypeInt, bin.Left.Data.(int32), int32(bin.Right.Data.(byte)))
			case ExpressionTypeShort:
				e.Data, _ = e.relationCompare(ExpressionTypeInt, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case ExpressionTypeInt:
				e.Data, _ = e.relationCompare(ExpressionTypeInt, bin.Left.Data.(int32), bin.Right.Data.(int32))
			case ExpressionTypeLong:
				e.Data, _ = e.relationCompare(ExpressionTypeLong, int64(bin.Left.Data.(int32)), bin.Right.Data.(int64))
			}
			is = true
			e.Type = ExpressionTypeBool
			return
		case ExpressionTypeLong:
			switch bin.Right.Type {
			case ExpressionTypeByte:
				e.Data, _ = e.relationCompare(ExpressionTypeLong, bin.Left.Data.(int64), int64(bin.Right.Data.(byte)))
			case ExpressionTypeShort:
				e.Data, _ = e.relationCompare(ExpressionTypeLong, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
			case ExpressionTypeInt:
				e.Data, _ = e.relationCompare(ExpressionTypeLong, bin.Left.Data.(int64), int64(bin.Right.Data.(int32)))
			case ExpressionTypeLong:
				e.Data, _ = e.relationCompare(ExpressionTypeLong, bin.Left.Data.(int64), bin.Right.Data.(int64))
			}
			is = true
			e.Type = ExpressionTypeBool
			return
		}
		return
	}
	if bin.Left.isFloat() && bin.Right.isFloat() {
		switch bin.Left.Type {
		case ExpressionTypeFloat:
			switch bin.Right.Type {
			case ExpressionTypeFloat:
				e.Data, _ = e.relationCompare(ExpressionTypeFloat, bin.Left.Data.(float32), bin.Right.Data.(float32))
			case ExpressionTypeDouble:
				e.Data, _ = e.relationCompare(ExpressionTypeDouble, float64(bin.Left.Data.(float32)), bin.Right.Data.(float64))
			}
			is = true
			e.Type = ExpressionTypeBool
			return
		case ExpressionTypeDouble:
			switch bin.Right.Type {
			case ExpressionTypeFloat:
				e.Data, _ = e.relationCompare(ExpressionTypeDouble, bin.Left.Data.(float64), float64(bin.Right.Data.(float32)))
			case ExpressionTypeDouble:
				e.Data, _ = e.relationCompare(ExpressionTypeDouble, bin.Left.Data.(float64), bin.Right.Data.(float64))
			}
			is = true
			e.Type = ExpressionTypeBool
			return
		}
		return
	}

	return
}

func (e *Expression) numberTypeAlgebra(typ ExpressionKind, value1, value2 interface{}) (value interface{}, support bool, err error) {
	support = true
	switch typ {
	case ExpressionTypeByte:
		switch e.Type {
		case ExpressionTypeAdd:
			value = byte(value1.(byte) + value2.(byte))
		case ExpressionTypeSub:
			value = byte(value1.(byte) - value2.(byte))
		case ExpressionTypeMul:
			value = byte(value1.(byte) * value2.(byte))
		case ExpressionTypeDiv:
			if value2.(byte) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = byte(value1.(byte) / value2.(byte))
			}
		case ExpressionTypeMod:
			if value2.(byte) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = byte(value1.(byte) % value2.(byte))
			}
		}
		return
	case ExpressionTypeShort:
		switch e.Type {
		case ExpressionTypeAdd:
			value = value1.(int32) + value2.(int32)
		case ExpressionTypeSub:
			value = value1.(int32) - value2.(int32)
		case ExpressionTypeMul:
			value = value1.(int32) * value2.(int32)
		case ExpressionTypeDiv:
			if value2.(int32) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) / value2.(int32)
			}
		case ExpressionTypeMod:
			if value2.(int32) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) % value2.(int32)
			}
		}
		return
	case ExpressionTypeInt:
		switch e.Type {
		case ExpressionTypeAdd:
			value = value1.(int32) + value2.(int32)
		case ExpressionTypeSub:
			value = value1.(int32) - value2.(int32)
		case ExpressionTypeMul:
			value = value1.(int32) * value2.(int32)
		case ExpressionTypeDiv:
			if value2.(int32) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) / value2.(int32)
			}
		case ExpressionTypeMod:
			if value2.(int32) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(int32) % value2.(int32)
			}
		}
		return
	case ExpressionTypeLong:
		switch e.Type {
		case ExpressionTypeAdd:
			value = value1.(int64) + value2.(int64)
		case ExpressionTypeSub:
			value = value1.(int64) - value2.(int64)
		case ExpressionTypeMul:
			value = value1.(int64) * value2.(int64)
		case ExpressionTypeDiv:
			if value2.(int64) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(int64) / value2.(int64)
			}
		case ExpressionTypeMod:
			if value2.(int64) == 0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(int64) % value2.(int64)
			}
		}
		return
	case ExpressionTypeFloat:
		switch e.Type {
		case ExpressionTypeAdd:
			value = value1.(float32) + value2.(float32)
		case ExpressionTypeSub:
			value = value1.(float32) - value2.(float32)
		case ExpressionTypeMul:
			value = value1.(float32) * value2.(float32)
		case ExpressionTypeDiv:
			if value2.(float32) == 0.0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(float32) / value2.(float32)
			}
		}
		return
	case ExpressionTypeDouble:
		switch e.Type {
		case ExpressionTypeAdd:
			value = value1.(float64) + value2.(float64)
		case ExpressionTypeSub:
			value = value1.(float64) - value2.(float64)
		case ExpressionTypeMul:
			value = value1.(float64) * value2.(float64)
		case ExpressionTypeDiv:
			if value2.(float64) == 0.0 {
				err = divisionByZeroErr(e.Pos)
			} else {
				value = value1.(float64) / value2.(float64)
			}
		}
		return
	}
	support = false
	return
}

func (e *Expression) relationCompare(typ ExpressionKind, value1, value2 interface{}) (b, support bool) {
	support = true
	switch typ {
	case ExpressionTypeBool:
		if e.Type == ExpressionTypeEq {
			b = value1.(bool) == value2.(bool)
		} else if e.Type == ExpressionTypeNe {
			b = value1.(bool) != value2.(bool)
		}
		return
	case ExpressionTypeByte:
		if e.Type == ExpressionTypeEq {
			b = value1.(byte) == value2.(byte)
		} else if e.Type == ExpressionTypeNe {
			b = value1.(byte) != value2.(byte)
		} else if e.Type == ExpressionTypeGt {
			b = value1.(byte) > value2.(byte)
		} else if e.Type == ExpressionTypeGe {
			b = value1.(byte) >= value2.(byte)
		} else if e.Type == ExpressionTypeLt {
			b = value1.(byte) < value2.(byte)
		} else if e.Type == ExpressionTypeLe {
			b = value1.(byte) <= value2.(byte)
		}
		return
	case ExpressionTypeInt:
		if e.Type == ExpressionTypeEq {
			b = value1.(int32) == value2.(int32)
		} else if e.Type == ExpressionTypeNe {
			b = value1.(int32) != value2.(int32)
		} else if e.Type == ExpressionTypeGt {
			b = value1.(int32) > value2.(int32)
		} else if e.Type == ExpressionTypeGe {
			b = value1.(int32) >= value2.(int32)
		} else if e.Type == ExpressionTypeLt {
			b = value1.(int32) < value2.(int32)
		} else if e.Type == ExpressionTypeLe {
			b = value1.(int32) <= value2.(int32)
		}
		return
	case ExpressionTypeLong:
		if e.Type == ExpressionTypeEq {
			b = value1.(int64) == value2.(int64)
		} else if e.Type == ExpressionTypeNe {
			b = value1.(int64) != value2.(int64)
		} else if e.Type == ExpressionTypeGt {
			b = value1.(int64) > value2.(int64)
		} else if e.Type == ExpressionTypeGe {
			b = value1.(int64) >= value2.(int64)
		} else if e.Type == ExpressionTypeLt {
			b = value1.(int64) < value2.(int64)
		} else if e.Type == ExpressionTypeLe {
			b = value1.(int64) <= value2.(int64)
		}
		return
	case ExpressionTypeFloat:
		if e.Type == ExpressionTypeEq {
			b = value1.(float32) == value2.(float32)
		} else if e.Type == ExpressionTypeNe {
			b = value1.(float32) != value2.(float32)
		} else if e.Type == ExpressionTypeGt {
			b = value1.(float32) > value2.(float32)
		} else if e.Type == ExpressionTypeGe {
			b = value1.(float32) >= value2.(float32)
		} else if e.Type == ExpressionTypeLt {
			b = value1.(float32) < value2.(float32)
		} else if e.Type == ExpressionTypeLe {
			b = value1.(float32) <= value2.(float32)
		}
		return
	case ExpressionTypeDouble:
		if e.Type == ExpressionTypeEq {
			b = value1.(float64) == value2.(float64)
		} else if e.Type == ExpressionTypeNe {
			b = value1.(float64) != value2.(float64)
		} else if e.Type == ExpressionTypeGt {
			b = value1.(float64) > value2.(float64)
		} else if e.Type == ExpressionTypeGe {
			b = value1.(float64) >= value2.(float64)
		} else if e.Type == ExpressionTypeLt {
			b = value1.(float64) < value2.(float64)
		} else if e.Type == ExpressionTypeLe {
			b = value1.(float64) <= value2.(float64)
		}
		return
	case ExpressionTypeString:
		if e.Type == ExpressionTypeEq {
			b = value1.(string) == value2.(string)
		} else if e.Type == ExpressionTypeNe {
			b = value1.(string) != value2.(string)
		} else if e.Type == ExpressionTypeGt {
			b = value1.(string) > value2.(string)
		} else if e.Type == ExpressionTypeGe {
			b = value1.(string) >= value2.(string)
		} else if e.Type == ExpressionTypeLt {
			b = value1.(string) < value2.(string)
		} else if e.Type == ExpressionTypeLe {
			b = value1.(string) <= value2.(string)
		}
		return
	}
	support = false
	return
}
