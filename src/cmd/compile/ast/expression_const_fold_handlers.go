package ast

func (e *Expression) arithmeticBinaryConstFolder(bin *ExpressionBinary) (is bool, err error) {
	if bin.Left.Type != bin.Right.Type {
		return
	}
	switch bin.Left.Type {
	case ExpressionTypeByte:
		left := bin.Left.Data.(byte)
		right := bin.Right.Data.(byte)
		switch e.Type {
		case ExpressionTypeAdd:
			e.Data = left + right
		case ExpressionTypeSub:
			e.Data = left - right
		case ExpressionTypeMul:
			e.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		default:
			return false, nil
		}
		e.Type = ExpressionTypeByte
		is = true
		return
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeInt:
		left := bin.Left.Data.(int32)
		right := bin.Right.Data.(int32)
		switch e.Type {
		case ExpressionTypeAdd:
			e.Data = left + right
		case ExpressionTypeSub:
			e.Data = left - right
		case ExpressionTypeMul:
			e.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		default:
			return false, nil
		}
		e.Type = bin.Left.Type
		is = true
		return
	case ExpressionTypeLong:
		left := bin.Left.Data.(int64)
		right := bin.Right.Data.(int64)
		switch e.Type {
		case ExpressionTypeAdd:
			e.Data = left + right
		case ExpressionTypeSub:
			e.Data = left - right
		case ExpressionTypeMul:
			e.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		default:
			return false, nil
		}
		e.Type = ExpressionTypeLong
		is = true
		return
	case ExpressionTypeFloat:
		left := bin.Left.Data.(float32)
		right := bin.Right.Data.(float32)
		switch e.Type {
		case ExpressionTypeAdd:
			e.Data = left + right
		case ExpressionTypeSub:
			e.Data = left - right
		case ExpressionTypeMul:
			e.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		default:
			return false, nil
		}
		e.Type = ExpressionTypeFloat
		is = true
		return
	case ExpressionTypeDouble:
		left := bin.Left.Data.(float64)
		right := bin.Right.Data.(float64)
		switch e.Type {
		case ExpressionTypeAdd:
			e.Data = left + right
		case ExpressionTypeSub:
			e.Data = left - right
		case ExpressionTypeMul:
			e.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			e.Data = left
		default:
			return false, nil
		}
		e.Type = ExpressionTypeDouble
		is = true
		return
	case ExpressionTypeString:
		left := bin.Left.Data.(string)
		right := bin.Right.Data.(string)
		if e.Type == ExpressionTypeAdd {
			e.Type = ExpressionTypeString
			e.Data = left + right
		} else {
			return false, nil
		}
	}
	return
}

func (e *Expression) relationBinaryConstFolder(bin *ExpressionBinary) (is bool, err error) {
	if bin.Left.Type == ExpressionTypeBool &&
		bin.Right.Type == ExpressionTypeBool &&
		e.isEqOrNe() {
		if e.Type == ExpressionTypeEq {
			e.Data = bin.Left.Data.(bool) == bin.Right.Data.(bool)
		} else {
			e.Data = bin.Left.Data.(bool) != bin.Right.Data.(bool)
		}
		e.Type = ExpressionTypeBool
		return
	}
	if bin.Left.Type != bin.Right.Type {
		return false, nil
	}
	switch bin.Left.Type {
	case ExpressionTypeString:
		left := bin.Left.Data.(string)
		right := bin.Right.Data.(string)
		switch e.Type {
		case ExpressionTypeEq:
			e.Data = left == right
		case ExpressionTypeNe:
			e.Data = left != right
		case ExpressionTypeGe:
			e.Data = left >= right
		case ExpressionTypeGt:
			e.Data = left > right
		case ExpressionTypeLe:
			e.Data = left <= right
		case ExpressionTypeLt:
			e.Data = left < right
		}
		is = true
		e.Type = ExpressionTypeBool
		return
	case ExpressionTypeByte:
		left := bin.Left.Data.(byte)
		right := bin.Right.Data.(byte)
		switch e.Type {
		case ExpressionTypeEq:
			e.Data = left == right
		case ExpressionTypeNe:
			e.Data = left != right
		case ExpressionTypeGe:
			e.Data = left >= right
		case ExpressionTypeGt:
			e.Data = left > right
		case ExpressionTypeLe:
			e.Data = left <= right
		case ExpressionTypeLt:
			e.Data = left < right
		}
		is = true
		e.Type = ExpressionTypeBool
		return
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeInt:
		left := bin.Left.Data.(int32)
		right := bin.Right.Data.(int32)
		switch e.Type {
		case ExpressionTypeEq:
			e.Data = left == right
		case ExpressionTypeNe:
			e.Data = left != right
		case ExpressionTypeGe:
			e.Data = left >= right
		case ExpressionTypeGt:
			e.Data = left > right
		case ExpressionTypeLe:
			e.Data = left <= right
		case ExpressionTypeLt:
			e.Data = left < right
		}
		is = true
		e.Type = ExpressionTypeBool
		return
	case ExpressionTypeLong:
		left := bin.Left.Data.(int64)
		right := bin.Right.Data.(int64)
		switch e.Type {
		case ExpressionTypeEq:
			e.Data = left == right
		case ExpressionTypeNe:
			e.Data = left != right
		case ExpressionTypeGe:
			e.Data = left >= right
		case ExpressionTypeGt:
			e.Data = left > right
		case ExpressionTypeLe:
			e.Data = left <= right
		case ExpressionTypeLt:
			e.Data = left < right
		}
		is = true
		e.Type = ExpressionTypeBool
		return
	case ExpressionTypeFloat:
		left := bin.Left.Data.(float32)
		right := bin.Right.Data.(float32)
		switch e.Type {
		case ExpressionTypeEq:
			e.Data = left == right
		case ExpressionTypeNe:
			e.Data = left != right
		case ExpressionTypeGe:
			e.Data = left >= right
		case ExpressionTypeGt:
			e.Data = left > right
		case ExpressionTypeLe:
			e.Data = left <= right
		case ExpressionTypeLt:
			e.Data = left < right
		}
		is = true
		e.Type = ExpressionTypeBool
		return
	case ExpressionTypeDouble:
		left := bin.Left.Data.(float64)
		right := bin.Right.Data.(float64)
		switch e.Type {
		case ExpressionTypeEq:
			e.Data = left == right
		case ExpressionTypeNe:
			e.Data = left != right
		case ExpressionTypeGe:
			e.Data = left >= right
		case ExpressionTypeGt:
			e.Data = left > right
		case ExpressionTypeLe:
			e.Data = left <= right
		case ExpressionTypeLt:
			e.Data = left < right
		}
		is = true
		e.Type = ExpressionTypeBool
		return
	}

	return
}
