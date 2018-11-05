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

func (e *Expression) binaryWrongOpErr() error {
	var typ1, typ2 string
	bin := e.Data.(*ExpressionBinary)
	if bin.Left.Value != nil {
		typ1 = bin.Left.Value.TypeString()
	} else {
		typ1 = bin.Left.Op
	}
	if bin.Right.Value != nil {
		typ2 = bin.Right.Value.TypeString()
	} else {
		typ2 = bin.Right.Op
	}
	return fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
		e.Pos.ErrMsgPrefix(),
		e.Op,
		typ1,
		typ2)
}

func (e *Expression) byteExceeds(t int64) error {
	e.Data = int64(byte(t))
	return fmt.Errorf("%s constant %d exceeds [-128 , 127 ]", e.Pos.ErrMsgPrefix(), t)
}
func (e *Expression) shortExceeds(t int64) error {
	e.Data = int64(int16(t))
	return fmt.Errorf("%s constant %d exceeds [-32768 , 32767 ]", e.Pos.ErrMsgPrefix(), t)
}
func (e *Expression) charExceeds(t int64) error {
	e.Data = int64(uint16(t))
	return fmt.Errorf("%s constant %d exceeds [0 , 65535 ]", e.Pos.ErrMsgPrefix(), t)
}
func (e *Expression) intExceeds(t int64) error {
	e.Data = int64(int32(t))
	return fmt.Errorf("%s constant %d exceeds [-32768 , 32767 ]",
		e.Pos.ErrMsgPrefix(), t)
}
func (e *Expression) longExceeds(t int64) error {
	return fmt.Errorf("%s constant exceeds [-9223372036854775808 , 9223372036854775807 ]",
		e.Pos.ErrMsgPrefix())
}
func (e *Expression) floatExceeds() error {
	return fmt.Errorf("%s float constant exceeds", e.Pos.ErrMsgPrefix())
}
func (e *Expression) doubleExceeds() error {
	return fmt.Errorf("%s double constant exceeds", e.Pos.ErrMsgPrefix())
}

func (e *Expression) constantFold() (is bool, err error) {
	if e.isLiteral() {
		if e.checkRangeCalled {
			return true, nil
		}
		e.checkRangeCalled = true
		switch e.Type {
		case ExpressionTypeByte:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t > int64(math.MaxInt8) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.byteExceeds(t))
				}
			} else {
				if t > (int64(math.MaxInt8) + 1) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.byteExceeds(t))
				}
				e.AsSubForNegative.Data = -t
				e.AsSubForNegative.Type = ExpressionTypeByte
			}
		case ExpressionTypeShort:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t > int64(math.MaxInt16) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.shortExceeds(t))
				}
			} else {
				if t > (int64(math.MaxInt16) + 1) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.shortExceeds(t))
				}
				e.AsSubForNegative.Data = -t
				e.AsSubForNegative.Type = ExpressionTypeShort
			}
		case ExpressionTypeChar:
			t := e.Data.(int64)
			if t > int64(math.MaxUint16) {
				PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.charExceeds(t))
			}
			e.Data = t
		case ExpressionTypeInt:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t > int64(math.MaxInt32) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.intExceeds(t))
				}
			} else {
				if t > (int64(math.MaxInt32) + 1) {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.intExceeds(t))
				}
				e.AsSubForNegative.Data = -t
				e.AsSubForNegative.Type = ExpressionTypeInt
			}
		case ExpressionTypeLong:
			t := e.Data.(int64)
			if e.AsSubForNegative == nil {
				if t>>63 != 0 {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors)
				}
			} else {
				if (t>>63 != 0) &&
					(t<<1) != 0 {
					PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.longExceeds(t))
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
			e.Data = ^ee.Data.(int64)
		case ExpressionTypeChar:
			e.Data = ^ee.Data.(int64)
		case ExpressionTypeShort:
			e.Data = ^ee.Data.(int64)
		case ExpressionTypeInt:
			e.Data = ^ee.Data.(int64)
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
				err = e.binaryWrongOpErr()
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
				fallthrough
			case ExpressionTypeShort:
				fallthrough
			case ExpressionTypeChar:
				fallthrough
			case ExpressionTypeInt:
				fallthrough
			case ExpressionTypeLong:
				if e.Type == ExpressionTypeLsh {
					e.Data = bin.Left.Data.(int64) << byte(bin.Right.getLongValue())
				} else {
					e.Data = bin.Left.Data.(int64) >> byte(bin.Right.getLongValue())
				}
			}
			if e.Type == ExpressionTypeLsh {
				switch bin.Left.Type {
				case ExpressionTypeByte:
					if t := e.Data.(int64); (t >> 8) != 0 {
						PackageBeenCompile.errors = append(PackageBeenCompile.errors, e.byteExceeds(t))
					}
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
					e.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					e.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			case ExpressionTypeShort:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					e.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			case ExpressionTypeChar:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					e.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			case ExpressionTypeInt:
				if e.Type == ExpressionTypeAnd {
					e.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if e.Type == ExpressionTypeOr {
					e.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					e.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
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

func (e *Expression) getLongValue() int64 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		fallthrough
	case ExpressionTypeLong:
		return e.Data.(int64)
	case ExpressionTypeFloat:
		return int64(e.Data.(float32))
	case ExpressionTypeDouble:
		return int64(e.Data.(float64))
	}
	panic("no match")
}

func (e *Expression) getDoubleValue() float64 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case ExpressionTypeByte:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		fallthrough
	case ExpressionTypeLong:
		return float64(e.Data.(int64))
	case ExpressionTypeFloat:
		return float64(e.Data.(float32))
	case ExpressionTypeDouble:
		return e.Data.(float64)
	}
	panic("no match")
}

func (e *Expression) convertLiteralToNumberType(to VariableTypeKind) {
	if e.isNumber() == false {
		panic("not a number")
	}
	switch to {
	case VariableTypeByte:
		e.Data = e.getLongValue()
		e.Type = ExpressionTypeByte
	case VariableTypeShort:
		e.Data = e.getLongValue()
		e.Type = ExpressionTypeShort
	case VariableTypeChar:
		e.Data = e.getLongValue()
		e.Type = ExpressionTypeChar
	case VariableTypeInt:
		e.Data = e.getLongValue()
		e.Type = ExpressionTypeInt
	case VariableTypeLong:
		e.Data = e.getLongValue()
		e.Type = ExpressionTypeLong
	case VariableTypeFloat:
		e.Data = float32(e.getDoubleValue())
		e.Type = ExpressionTypeFloat
	case VariableTypeDouble:
		e.Data = e.getDoubleValue()
		e.Type = ExpressionTypeDouble
	}
}
