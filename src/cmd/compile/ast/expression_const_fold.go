package ast

import (
	"fmt"
	"math"
)

func (e *Expression) getBinaryExpressionConstValue(folder binaryConstFolder) (is bool, err error) {
	bin := e.Data.(*ExpressionBinary)
	is1, err1 := bin.Left.constantFold()
	is2, err2 := bin.Right.constantFold()
	if err1 != nil { //something is wrong
		err = err1
		return
	}
	if err2 != nil {
		err = err2
		return
	}
	if is1 == false ||
		is2 == false {
		is = false
		err = nil
		return
	}
	return folder(bin)
}

type binaryConstFolder func(bin *ExpressionBinary) (is bool, err error)

func (e *Expression) makeWrongOpErr(typ1, typ2 string) error {
	return fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
		e.Pos.ErrMsgPrefix(),
		e.Op,
		typ1,
		typ2)
}

func (e *Expression) constantFold() (is bool, err error) {
	if e.isLiteral() {
		switch e.Type {
		case ExpressionTypeByte:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t > int64(math.MaxInt8) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant %d exceeds [-128 , 127 ]", e.Pos.ErrMsgPrefix(), t))
				}
				e.Data = byte(t)
			} else {
				if t > (int64(math.MaxInt8) + 1) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant %d exceeds [-128 , 127 ]", e.Pos.ErrMsgPrefix(), -t))
				}
				if t == (int64(math.MaxInt8) + 1) {
					e.AsSubForNegative.Data = byte(1 << 7)
				} else {
					e.AsSubForNegative.Data = -byte(t)
				}
				e.AsSubForNegative.Type = ExpressionTypeByte
			}
		case ExpressionTypeShort:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t > int64(math.MaxInt16) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant %d exceeds [-32768 , 32767 ]",
							e.Pos.ErrMsgPrefix(), t))
				}
				e.Data = int32(t)
			} else {
				if t > (int64(math.MaxInt16) + 1) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant %d exceeds [-128 , 127 ]",
							e.Pos.ErrMsgPrefix(), -t))
				}
				if t == (int64(math.MaxInt16) + 1) {
					e.AsSubForNegative.Data = int32(math.MinInt16)
				} else {
					e.AsSubForNegative.Data = -int32(t)
				}
				e.AsSubForNegative.Type = ExpressionTypeShort
			}
		case ExpressionTypeChar:
			t := e.Data.(int64)
			if t > int64(math.MaxUint16) {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors,
					fmt.Errorf("%s constant %d exceeds [0 , 65535 ]",
						e.Pos.ErrMsgPrefix(), t))
			}
			e.Data = int32(t)
		case ExpressionTypeInt:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t > int64(math.MaxInt32) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant %d exceeds [-32768 , 32767 ]",
							e.Pos.ErrMsgPrefix(), t))
				}
				e.Data = int32(t)
			} else {
				if t > (int64(math.MaxInt32) + 1) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant %d exceeds [-2147483648 , 2147483647 ]",
							e.Pos.ErrMsgPrefix(), -t))
				}
				if t == (int64(math.MaxInt32) + 1) {
					e.AsSubForNegative.Data = int32(math.MinInt32)
				} else {
					e.AsSubForNegative.Data = -int32(t)
				}
				e.AsSubForNegative.Type = ExpressionTypeInt

			}
		case ExpressionTypeLong:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t>>63 != 0 {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant  exceeds [-9223372036854775808 , 9223372036854775807 ]",
							e.Pos.ErrMsgPrefix()))
				}
			} else {
				if (t>>63 != 0) &&
					(t<<1) != 0 {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors,
						fmt.Errorf("%s constant exceeds [-9223372036854775808 , 9223372036854775807 ]",
							e.Pos.ErrMsgPrefix()))
				}
				e.AsSubForNegative.Data = -e.Data.(int64)
				e.AsSubForNegative.Type = ExpressionTypeLong
			}
		}
		return true, nil
	}
	// ~
	if e.Type == ExpressionTypeBitwiseNot {
		ee := e.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		if ee.isInteger() == false {
			err = fmt.Errorf("%s cannot apply '^' on a non-integer expression",
				errMsgPrefix(e.Pos))
			return
		}
		e.Type = ee.Type
		switch ee.Type {
		case ExpressionTypeByte:
			e.Data = ^ee.Data.(byte)
		case ExpressionTypeChar:
			e.Data = ^ee.Data.(int32)
		case ExpressionTypeShort:
			e.Data = ^ee.Data.(int32)
		case ExpressionTypeInt:
			e.Data = ^ee.Data.(int32)
		case ExpressionTypeLong:
			e.Data = ^ee.Data.(int64)
		}
	}
	// !
	if e.Type == ExpressionTypeNot {
		ee := e.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		if ee.Type != ExpressionTypeBool {
			err = fmt.Errorf("%s cannot apply '!' on a non-bool expression",
				errMsgPrefix(e.Pos))
			return
		}
		e.Type = ExpressionTypeBool
		e.Data = !ee.Data.(bool)
		return
	}
	// -
	if e.Type == ExpressionTypeNegative {
		ee := e.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		switch ee.Type {
		case ExpressionTypeFloat:
			is = true
			e.Data = -ee.Data.(float32)
			e.Type = ExpressionTypeFloat
			return
		case ExpressionTypeDouble:
			is = true
			e.Data = -ee.Data.(float64)
			e.Type = ExpressionTypeDouble
			return
		}
	}
	// && and ||
	if e.Type == ExpressionTypeLogicalAnd || e.Type == ExpressionTypeLogicalOr {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.Type != ExpressionTypeBool ||
				bin.Right.Type != ExpressionTypeBool {
				err = e.makeWrongOpErr(bin.Left.Op, bin.Right.Op)
				return
			}
			is = true
			if e.Type == ExpressionTypeLogicalAnd {
				e.Data = bin.Left.Data.(bool) && bin.Right.Data.(bool)
			} else {
				e.Data = bin.Left.Data.(bool) || bin.Right.Data.(bool)
			}
			e.Type = ExpressionTypeBool
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	// + - * / % algebra arithmetic
	if e.Type == ExpressionTypeAdd ||
		e.Type == ExpressionTypeSub ||
		e.Type == ExpressionTypeMul ||
		e.Type == ExpressionTypeDiv ||
		e.Type == ExpressionTypeMod {
		is, err = e.getBinaryExpressionConstValue(e.arithmeticBinaryConstFolder)
		return
	}

	// <<  >>
	if e.Type == ExpressionTypeLsh || e.Type == ExpressionTypeRsh {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.isInteger() == false || bin.Right.isInteger() == false {
				return
			}
			switch bin.Left.Type {
			case ExpressionTypeByte:
				if e.Type == ExpressionTypeLsh {
					e.Data = byte(bin.Left.Data.(byte) << bin.Right.getByteValue())
				} else {
					e.Data = byte(bin.Left.Data.(byte) >> bin.Right.getByteValue())
				}
			case ExpressionTypeShort:
				if e.Type == ExpressionTypeLsh {
					e.Data = int32(bin.Left.Data.(int32) << bin.Right.getByteValue())
				} else {
					e.Data = int32(bin.Left.Data.(int32) >> bin.Right.getByteValue())
				}
			case ExpressionTypeChar:
				if e.Type == ExpressionTypeLsh {
					e.Data = int32(bin.Left.Data.(int32) << bin.Right.getByteValue())
				} else {
					e.Data = int32(bin.Left.Data.(int32) >> bin.Right.getByteValue())
				}
			case ExpressionTypeInt:
				if e.Type == ExpressionTypeLsh {
					e.Data = int32(bin.Left.Data.(int32) << bin.Right.getByteValue())
				} else {
					e.Data = int32(bin.Left.Data.(int32) >> bin.Right.getByteValue())
				}
			case ExpressionTypeLong:
				if e.Type == ExpressionTypeLsh {
					e.Data = int64(bin.Left.Data.(int64) << bin.Right.getByteValue())
				} else {
					e.Data = int64(bin.Left.Data.(int64) >> bin.Right.getByteValue())
				}
			}
			e.Type = bin.Left.Type
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	// & | ^
	if e.Type == ExpressionTypeAnd ||
		e.Type == ExpressionTypeOr ||
		e.Type == ExpressionTypeXor {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.isInteger() == false || bin.Right.isInteger() == false ||
				bin.Left.Type != bin.Right.Type {
				return // not integer or type not equal
			}
			switch bin.Left.Type {
			case ExpressionTypeByte:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(byte) & bin.Right.Data.(byte)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(byte) | bin.Right.Data.(byte)
				} else {
					e.Data = bin.Left.Data.(byte) ^ bin.Right.Data.(byte)
				}
			case ExpressionTypeShort:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(int32) & bin.Right.Data.(int32)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int32) | bin.Right.Data.(int32)
				} else {
					e.Data = bin.Left.Data.(int32) ^ bin.Right.Data.(int32)
				}
			case ExpressionTypeChar:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(int32) & bin.Right.Data.(int32)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int32) | bin.Right.Data.(int32)
				} else {
					e.Data = bin.Left.Data.(int32) ^ bin.Right.Data.(int32)
				}
			case ExpressionTypeInt:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(int32) & bin.Right.Data.(int32)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int32) | bin.Right.Data.(int32)
				} else {
					e.Data = bin.Left.Data.(int32) ^ bin.Right.Data.(int32)
				}
			case ExpressionTypeLong:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					e.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			}
			is = true
			e.Type = bin.Left.Type
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	if e.Type == ExpressionTypeNot {
		ee := e.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil {
			return
		}
		if is == false {
			return
		}
		if ee.Type != ExpressionTypeBool {
			return false, fmt.Errorf("!(not) can only apply to bool expression")
		}
		is = true
		e.Type = ExpressionTypeBool
		e.Data = !ee.Data.(bool)
		return
	}
	//  == != > < >= <=
	if e.Type == ExpressionTypeEq ||
		e.Type == ExpressionTypeNe ||
		e.Type == ExpressionTypeGe ||
		e.Type == ExpressionTypeGt ||
		e.Type == ExpressionTypeLe ||
		e.Type == ExpressionTypeLt {
		return e.getBinaryExpressionConstValue(e.relationBinaryConstFolder)
	}
	return
}

func (e *Expression) getByteValue() byte {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		return e.Data.(byte)
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		return byte(e.Data.(int32))
	case ExpressionTypeLong:
		return byte(e.Data.(int64))
	case ExpressionTypeFloat:
		return byte(e.Data.(float32))
	case ExpressionTypeDouble:
		return byte(e.Data.(float64))
	}
	return 0
}

func (e *Expression) getShortValue() int32 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		return int32(e.Data.(byte))
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		return int32(e.Data.(int32))
	case ExpressionTypeLong:
		return int32(e.Data.(int64))
	case ExpressionTypeFloat:
		return int32(e.Data.(float32))
	case ExpressionTypeDouble:
		return int32(e.Data.(float64))
	}
	return 0
}

func (e *Expression) getCharValue() int32 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		return int32(e.Data.(byte))
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		return int32(e.Data.(int32))
	case ExpressionTypeLong:
		return int32(e.Data.(int64))
	case ExpressionTypeFloat:
		return int32(e.Data.(float32))
	case ExpressionTypeDouble:
		return int32(e.Data.(float64))
	}
	return 0
}
func (e *Expression) getIntValue() int32 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		return int32(e.Data.(byte))
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		return int32(e.Data.(int32))
	case ExpressionTypeLong:
		return int32(e.Data.(int64))
	case ExpressionTypeFloat:
		return int32(e.Data.(float32))
	case ExpressionTypeDouble:
		return int32(e.Data.(float64))
	}
	return 0
}

func (e *Expression) getLongValue() int64 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		return int64(e.Data.(byte))
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		return int64(e.Data.(int32))
	case ExpressionTypeLong:
		return int64(e.Data.(int64))
	case ExpressionTypeFloat:
		return int64(e.Data.(float32))
	case ExpressionTypeDouble:
		return int64(e.Data.(float64))
	}
	return 0
}
func (e *Expression) getFloatValue() float32 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		return float32(e.Data.(byte))
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		return float32(e.Data.(int32))
	case ExpressionTypeLong:
		return float32(e.Data.(int64))
	case ExpressionTypeFloat:
		return float32(e.Data.(float32))
	case ExpressionTypeDouble:
		return float32(e.Data.(float64))
	}
	return 0
}

func (e *Expression) getDoubleValue() float64 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		return float64(e.Data.(byte))
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		return float64(e.Data.(int32))
	case ExpressionTypeLong:
		return float64(e.Data.(int64))
	case ExpressionTypeFloat:
		return float64(e.Data.(float32))
	case ExpressionTypeDouble:
		return float64(e.Data.(float64))
	}
	return 0
}

func (e *Expression) convertLiteralToNumberType(to VariableTypeKind) {
	if e.isNumber() == false {
		panic("not a number")
	}
	switch to {
	case VariableTypeByte:
		e.Data = e.getByteValue()
		e.Type = ExpressionTypeByte
	case VariableTypeShort:
		e.Data = e.getShortValue()
		e.Type = ExpressionTypeShort
	case VariableTypeChar:
		e.Data = e.getCharValue()
		e.Type = ExpressionTypeChar
	case VariableTypeInt:
		e.Data = e.getIntValue()
		e.Type = ExpressionTypeInt
	case VariableTypeLong:
		e.Data = e.getLongValue()
		e.Type = ExpressionTypeLong
	case VariableTypeFloat:
		e.Data = e.getFloatValue()
		e.Type = ExpressionTypeFloat
	case VariableTypeDouble:
		e.Data = e.getDoubleValue()
		e.Type = ExpressionTypeDouble
	}
}
