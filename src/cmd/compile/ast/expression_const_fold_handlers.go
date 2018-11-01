package ast

func (e *Expression) arithmeticBinaryConstFolder(bin *ExpressionBinary) (is bool, err error) {
	if bin.Left.Type != bin.Right.Type {
		return
	}
	switch bin.Left.Type {
	case ExpressionTypeByte:
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
		if e.Type == ExpressionTypeAdd || e.Type == ExpressionTypeSub {
			if t := e.Data.(int64); (t >> 8) != 0 {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.byteExceeds(t))
			}
		}
		e.Type = ExpressionTypeByte
		is = true
		return
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeInt:
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
		switch e.Type {
		case ExpressionTypeShort:
			if t := e.Data.(int64); (t >> 16) != 0 {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.shortExceeds(t))
			}
		case ExpressionTypeChar:
			if t := e.Data.(int64); (t >> 16) != 0 {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.charExceeds(t))
			}
		case ExpressionTypeInt:
			if t := e.Data.(int64); (t >> 32) != 0 {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.intExceeds(t))
			}
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
		if e.Type == ExpressionTypeAdd || e.Type == ExpressionTypeSub {
			if t := e.Data.(int64) < 0; t != (bin.Left.Data.(int64) < 0) && t != (bin.Left.Data.(int64) < 0) {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.longExceeds(e.Data.(int64)))
			}
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
		if e.Type == ExpressionTypeAdd || e.Type == ExpressionTypeSub {
			if t := e.Data.(float32) < 0; t != (bin.Left.Data.(float32) < 0) && t != (bin.Left.Data.(float32) < 0) {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.floatExceeds())
			}
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
		if e.Type == ExpressionTypeAdd || e.Type == ExpressionTypeSub {
			if t := e.Data.(float64) < 0; t != (bin.Left.Data.(float64) < 0) && t != (bin.Left.Data.(float64) < 0) {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.floatExceeds())
			}
		}
		e.Type = ExpressionTypeDouble
		is = true
		return
	case ExpressionTypeString:
		left := bin.Left.Data.(string)
		right := bin.Right.Data.(string)
		if e.Type == ExpressionTypeAdd {
			if len(left)+len(right) < 65536 {
				e.Type = ExpressionTypeString
				e.Data = left + right
			} else {
				return false, nil
			}
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
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeInt:
		fallthrough
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
