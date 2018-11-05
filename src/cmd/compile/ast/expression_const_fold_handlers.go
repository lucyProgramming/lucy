package ast

func (this *Expression) arithmeticBinaryConstFolder(bin *ExpressionBinary) (is bool, err error) {
	if bin.Left.Type != bin.Right.Type {
		return
	}
	switch bin.Left.Type {
	case ExpressionTypeByte:
		left := bin.Left.Data.(int64)
		right := bin.Right.Data.(int64)
		switch this.Type {
		case ExpressionTypeAdd:
			this.Data = left + right
		case ExpressionTypeSub:
			this.Data = left - right
		case ExpressionTypeMul:
			this.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		default:
			return false, nil
		}
		//if this.Type == ExpressionTypeAdd || this.Type == ExpressionTypeSub {
		//	if t := this.Data.(int64); (t >> 8) != 0 {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.byteExceeds(t))
		//	}
		//}
		this.Type = ExpressionTypeByte
		is = true
		return
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeInt:
		left := bin.Left.Data.(int64)
		right := bin.Right.Data.(int64)
		switch this.Type {
		case ExpressionTypeAdd:
			this.Data = left + right
		case ExpressionTypeSub:
			this.Data = left - right
		case ExpressionTypeMul:
			this.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		default:
			return false, nil
		}
		//switch this.Type {
		//case ExpressionTypeShort:
		//	if t := this.Data.(int64); (t >> 16) != 0 {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.shortExceeds(t))
		//	}
		//case ExpressionTypeChar:
		//	if t := this.Data.(int64); (t >> 16) != 0 {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.charExceeds(t))
		//	}
		//case ExpressionTypeInt:
		//	if t := this.Data.(int64); (t >> 32) != 0 {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.intExceeds(t))
		//	}
		//}
		this.Type = bin.Left.Type
		is = true
		return
	case ExpressionTypeLong:
		left := bin.Left.Data.(int64)
		right := bin.Right.Data.(int64)
		switch this.Type {
		case ExpressionTypeAdd:
			this.Data = left + right

		case ExpressionTypeSub:
			this.Data = left - right
		case ExpressionTypeMul:
			this.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		default:
			return false, nil
		}
		//if this.Type == ExpressionTypeAdd || this.Type == ExpressionTypeSub {
		//	if t := this.Data.(int64) < 0; t != (bin.Left.Data.(int64) < 0) && t != (bin.Left.Data.(int64) < 0) {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.longExceeds(this.Data.(int64)))
		//	}
		//}
		this.Type = ExpressionTypeLong
		is = true
		return
	case ExpressionTypeFloat:
		left := bin.Left.Data.(float32)
		right := bin.Right.Data.(float32)
		switch this.Type {
		case ExpressionTypeAdd:
			this.Data = left + right
		case ExpressionTypeSub:
			this.Data = left - right
		case ExpressionTypeMul:
			this.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		default:
			return false, nil
		}
		//if this.Type == ExpressionTypeAdd || this.Type == ExpressionTypeSub {
		//	if t := this.Data.(float32) < 0; t != (bin.Left.Data.(float32) < 0) && t != (bin.Left.Data.(float32) < 0) {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.floatExceeds())
		//	}
		//}
		this.Type = ExpressionTypeFloat
		is = true
		return
	case ExpressionTypeDouble:
		left := bin.Left.Data.(float64)
		right := bin.Right.Data.(float64)
		switch this.Type {
		case ExpressionTypeAdd:
			this.Data = left + right
		case ExpressionTypeSub:
			this.Data = left - right
		case ExpressionTypeMul:
			this.Data = left * right
		case ExpressionTypeDiv:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		case ExpressionTypeMod:
			if right == 0 {
				err = divisionByZeroErr(bin.Right.Pos)
			}
			this.Data = left
		default:
			return false, nil
		}
		//if this.Type == ExpressionTypeAdd || this.Type == ExpressionTypeSub {
		//	if t := this.Data.(float64) < 0; t != (bin.Left.Data.(float64) < 0) && t != (bin.Left.Data.(float64) < 0) {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.floatExceeds())
		//	}
		//}
		this.Type = ExpressionTypeDouble
		is = true
		return
	case ExpressionTypeString:
		left := bin.Left.Data.(string)
		right := bin.Right.Data.(string)
		if this.Type == ExpressionTypeAdd {
			if len(left)+len(right) < 65536 {
				this.Type = ExpressionTypeString
				this.Data = left + right
			} else {
				return false, nil
			}
		} else {
			return false, nil
		}
	}
	return
}

func (this *Expression) relationBinaryConstFolder(bin *ExpressionBinary) (is bool, err error) {
	if bin.Left.Type == ExpressionTypeBool &&
		bin.Right.Type == ExpressionTypeBool &&
		this.isEqOrNe() {
		if this.Type == ExpressionTypeEq {
			this.Data = bin.Left.Data.(bool) == bin.Right.Data.(bool)
		} else {
			this.Data = bin.Left.Data.(bool) != bin.Right.Data.(bool)
		}
		this.Type = ExpressionTypeBool
		return
	}
	if bin.Left.Type != bin.Right.Type {
		return false, nil
	}
	switch bin.Left.Type {
	case ExpressionTypeString:
		left := bin.Left.Data.(string)
		right := bin.Right.Data.(string)
		switch this.Type {
		case ExpressionTypeEq:
			this.Data = left == right
		case ExpressionTypeNe:
			this.Data = left != right
		case ExpressionTypeGe:
			this.Data = left >= right
		case ExpressionTypeGt:
			this.Data = left > right
		case ExpressionTypeLe:
			this.Data = left <= right
		case ExpressionTypeLt:
			this.Data = left < right
		}
		is = true
		this.Type = ExpressionTypeBool
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
		switch this.Type {
		case ExpressionTypeEq:
			this.Data = left == right
		case ExpressionTypeNe:
			this.Data = left != right
		case ExpressionTypeGe:
			this.Data = left >= right
		case ExpressionTypeGt:
			this.Data = left > right
		case ExpressionTypeLe:
			this.Data = left <= right
		case ExpressionTypeLt:
			this.Data = left < right
		}
		is = true
		this.Type = ExpressionTypeBool
		return
	case ExpressionTypeFloat:
		left := bin.Left.Data.(float32)
		right := bin.Right.Data.(float32)
		switch this.Type {
		case ExpressionTypeEq:
			this.Data = left == right
		case ExpressionTypeNe:
			this.Data = left != right
		case ExpressionTypeGe:
			this.Data = left >= right
		case ExpressionTypeGt:
			this.Data = left > right
		case ExpressionTypeLe:
			this.Data = left <= right
		case ExpressionTypeLt:
			this.Data = left < right
		}
		is = true
		this.Type = ExpressionTypeBool
		return
	case ExpressionTypeDouble:
		left := bin.Left.Data.(float64)
		right := bin.Right.Data.(float64)
		switch this.Type {
		case ExpressionTypeEq:
			this.Data = left == right
		case ExpressionTypeNe:
			this.Data = left != right
		case ExpressionTypeGe:
			this.Data = left >= right
		case ExpressionTypeGt:
			this.Data = left > right
		case ExpressionTypeLe:
			this.Data = left <= right
		case ExpressionTypeLt:
			this.Data = left < right
		}
		is = true
		this.Type = ExpressionTypeBool
		return
	}
	return
}
