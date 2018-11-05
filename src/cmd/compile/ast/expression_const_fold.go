package ast

import (
	"fmt"
)

func (this *Expression) getBinaryExpressionConstValue(folder binaryConstFolder) (is bool, err error) {
	bin := this.Data.(*ExpressionBinary)
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

func (this *Expression) binaryWrongOpErr() error {
	var typ1, typ2 string
	bin := this.Data.(*ExpressionBinary)
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
		this.Pos.ErrMsgPrefix(),
		this.Op,
		typ1,
		typ2)
}

//
//func (this *Expression) byteExceeds(t int64) error {
//	this.Data = int64(byte(t))
//	return fmt.Errorf("%s constant %d exceeds [-128 , 127 ]", this.Pos.ErrMsgPrefix(), t)
//}
//func (this *Expression) shortExceeds(t int64) error {
//	this.Data = int64(int16(t))
//	return fmt.Errorf("%s constant %d exceeds [-32768 , 32767 ]", this.Pos.ErrMsgPrefix(), t)
//}
//func (this *Expression) charExceeds(t int64) error {
//	this.Data = int64(uint16(t))
//	return fmt.Errorf("%s constant %d exceeds [0 , 65535 ]", this.Pos.ErrMsgPrefix(), t)
//}
//func (this *Expression) intExceeds(t int64) error {
//	this.Data = int64(int32(t))
//	return fmt.Errorf("%s constant %d exceeds [-32768 , 32767 ]",
//		this.Pos.ErrMsgPrefix(), t)
//}
//func (this *Expression) longExceeds(t int64) error {
//	return fmt.Errorf("%s constant exceeds [-9223372036854775808 , 9223372036854775807 ]",
//		this.Pos.ErrMsgPrefix())
//}
//func (this *Expression) floatExceeds() error {
//	return fmt.Errorf("%s float constant exceeds", this.Pos.ErrMsgPrefix())
//}
//func (this *Expression) doubleExceeds() error {
//	return fmt.Errorf("%s double constant exceeds", this.Pos.ErrMsgPrefix())
//}

func (this *Expression) constantFold() (is bool, err error) {
	if this.isLiteral() {
		//if this.checkRangeCalled {
		//	return true, nil
		//}
		//this.checkRangeCalled = true
		//switch this.Type {
		//case ExpressionTypeByte:
		//	t := this.Data.(int64)
		//	if this.AsSubForNegative == nil {
		//		if t > int64(math.MaxInt8) {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.byteExceeds(t))
		//		}
		//	} else {
		//		if t > (int64(math.MaxInt8) + 1) {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.byteExceeds(t))
		//		}
		//		this.AsSubForNegative.Data = -t
		//		this.AsSubForNegative.Type = ExpressionTypeByte
		//	}
		//case ExpressionTypeShort:
		//	t := this.Data.(int64)
		//	if this.AsSubForNegative == nil {
		//		if t > int64(math.MaxInt16) {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.shortExceeds(t))
		//		}
		//	} else {
		//		if t > (int64(math.MaxInt16) + 1) {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.shortExceeds(t))
		//		}
		//		this.AsSubForNegative.Data = -t
		//		this.AsSubForNegative.Type = ExpressionTypeShort
		//	}
		//case ExpressionTypeChar:
		//	t := this.Data.(int64)
		//	if t > int64(math.MaxUint16) {
		//		PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.charExceeds(t))
		//	}
		//	this.Data = t
		//case ExpressionTypeInt:
		//	t := this.Data.(int64)
		//	if this.AsSubForNegative == nil {
		//		if t > int64(math.MaxInt32) {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.intExceeds(t))
		//		}
		//	} else {
		//		if t > (int64(math.MaxInt32) + 1) {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.intExceeds(t))
		//		}
		//		this.AsSubForNegative.Data = -t
		//		this.AsSubForNegative.Type = ExpressionTypeInt
		//	}
		//case ExpressionTypeLong:
		//	t := this.Data.(int64)
		//	if this.AsSubForNegative == nil {
		//		if t>>63 != 0 {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors)
		//		}
		//	} else {
		//		if (t>>63 != 0) &&
		//			(t<<1) != 0 {
		//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.longExceeds(t))
		//		}
		//		this.AsSubForNegative.Data = -this.Data.(int64)
		//		this.AsSubForNegative.Type = ExpressionTypeLong
		//	}
		//}
		return true, nil
	}
	// ~
	if this.Type == ExpressionTypeBitwiseNot {
		ee := this.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		if ee.isInteger() == false {
			err = fmt.Errorf("%s cannot apply '^' on a non-integer expression",
				errMsgPrefix(this.Pos))
			return
		}
		this.Type = ee.Type
		switch ee.Type {
		case ExpressionTypeByte:
			this.Data = ^ee.Data.(int64)
		case ExpressionTypeChar:
			this.Data = ^ee.Data.(int64)
		case ExpressionTypeShort:
			this.Data = ^ee.Data.(int64)
		case ExpressionTypeInt:
			this.Data = ^ee.Data.(int64)
		case ExpressionTypeLong:
			this.Data = ^ee.Data.(int64)
		}
	}
	// !
	if this.Type == ExpressionTypeNot {
		ee := this.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		if ee.Type != ExpressionTypeBool {
			err = fmt.Errorf("%s cannot apply '!' on a non-bool expression",
				errMsgPrefix(this.Pos))
			return
		}
		this.Type = ExpressionTypeBool
		this.Data = !ee.Data.(bool)
		return
	}
	// -
	if this.Type == ExpressionTypeNegative {
		ee := this.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		switch ee.Type {
		case ExpressionTypeFloat:
			is = true
			this.Data = -ee.Data.(float32)
			this.Type = ExpressionTypeFloat
			return
		case ExpressionTypeDouble:
			is = true
			this.Data = -ee.Data.(float64)
			this.Type = ExpressionTypeDouble
			return
		}
	}
	// && and ||
	if this.Type == ExpressionTypeLogicalAnd || this.Type == ExpressionTypeLogicalOr {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.Type != ExpressionTypeBool ||
				bin.Right.Type != ExpressionTypeBool {
				err = this.binaryWrongOpErr()
				return
			}
			is = true
			if this.Type == ExpressionTypeLogicalAnd {
				this.Data = bin.Left.Data.(bool) && bin.Right.Data.(bool)
			} else {
				this.Data = bin.Left.Data.(bool) || bin.Right.Data.(bool)
			}
			this.Type = ExpressionTypeBool
			return
		}
		return this.getBinaryExpressionConstValue(f)
	}
	// + - * / % algebra arithmetic
	if this.Type == ExpressionTypeAdd ||
		this.Type == ExpressionTypeSub ||
		this.Type == ExpressionTypeMul ||
		this.Type == ExpressionTypeDiv ||
		this.Type == ExpressionTypeMod {
		is, err = this.getBinaryExpressionConstValue(this.arithmeticBinaryConstFolder)
		return
	}
	// <<  >>
	if this.Type == ExpressionTypeLsh || this.Type == ExpressionTypeRsh {
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
				if this.Type == ExpressionTypeLsh {
					this.Data = bin.Left.Data.(int64) << byte(bin.Right.getLongValue())
				} else {
					this.Data = bin.Left.Data.(int64) >> byte(bin.Right.getLongValue())
				}
			}
			//if this.Type == ExpressionTypeLsh {
			//	switch bin.Left.Type {
			//	case ExpressionTypeByte:
			//		if t := this.Data.(int64); (t >> 8) != 0 {
			//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.byteExceeds(t))
			//		}
			//	case ExpressionTypeShort:
			//		if t := this.Data.(int64); (t >> 16) != 0 {
			//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.shortExceeds(t))
			//		}
			//	case ExpressionTypeChar:
			//		if t := this.Data.(int64); (t >> 16) != 0 {
			//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.charExceeds(t))
			//		}
			//	case ExpressionTypeInt:
			//		if t := this.Data.(int64); (t >> 32) != 0 {
			//			PackageBeenCompile.errors = append(PackageBeenCompile.errors, this.intExceeds(t))
			//		}
			//	}
			//}
			this.Type = bin.Left.Type
			return
		}
		return this.getBinaryExpressionConstValue(f)
	}
	// & | ^
	if this.Type == ExpressionTypeAnd ||
		this.Type == ExpressionTypeOr ||
		this.Type == ExpressionTypeXor {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.isInteger() == false || bin.Right.isInteger() == false ||
				bin.Left.Type != bin.Right.Type {
				return // not integer or type not equal
			}
			switch bin.Left.Type {
			case ExpressionTypeByte:
				if this.Type == ExpressionTypeAnd {
					this.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if this.Type == ExpressionTypeOr {
					this.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					this.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			case ExpressionTypeShort:
				if this.Type == ExpressionTypeAnd {
					this.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if this.Type == ExpressionTypeOr {
					this.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					this.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			case ExpressionTypeChar:
				if this.Type == ExpressionTypeAnd {
					this.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if this.Type == ExpressionTypeOr {
					this.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					this.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			case ExpressionTypeInt:
				if this.Type == ExpressionTypeAnd {
					this.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if this.Type == ExpressionTypeOr {
					this.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					this.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			case ExpressionTypeLong:
				if this.Type == ExpressionTypeAnd {
					this.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if this.Type == ExpressionTypeOr {
					this.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					this.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			}
			is = true
			this.Type = bin.Left.Type
			return
		}
		return this.getBinaryExpressionConstValue(f)
	}
	if this.Type == ExpressionTypeNot {
		ee := this.Data.(*Expression)
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
		this.Type = ExpressionTypeBool
		this.Data = !ee.Data.(bool)
		return
	}
	//  == != > < >= <=
	if this.Type == ExpressionTypeEq ||
		this.Type == ExpressionTypeNe ||
		this.Type == ExpressionTypeGe ||
		this.Type == ExpressionTypeGt ||
		this.Type == ExpressionTypeLe ||
		this.Type == ExpressionTypeLt {
		return this.getBinaryExpressionConstValue(this.relationBinaryConstFolder)
	}
	return
}

func (this *Expression) getLongValue() int64 {
	if this.isNumber() == false {
		panic("not number")
	}
	switch this.Type {
	case ExpressionTypeByte:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		fallthrough
	case ExpressionTypeLong:
		return this.Data.(int64)
	case ExpressionTypeFloat:
		return int64(this.Data.(float32))
	case ExpressionTypeDouble:
		return int64(this.Data.(float64))
	}
	panic("no match")
}

func (this *Expression) getDoubleValue() float64 {
	if this.isNumber() == false {
		panic("not number")
	}
	switch this.Type {
	case ExpressionTypeByte:
		fallthrough
	case ExpressionTypeChar:
		fallthrough
	case ExpressionTypeShort:
		fallthrough
	case ExpressionTypeInt:
		fallthrough
	case ExpressionTypeLong:
		return float64(this.Data.(int64))
	case ExpressionTypeFloat:
		return float64(this.Data.(float32))
	case ExpressionTypeDouble:
		return this.Data.(float64)
	}
	panic("no match")
}

func (this *Expression) convertLiteralToNumberType(to VariableTypeKind) {
	if this.isNumber() == false {
		panic("not a number")
	}
	switch to {
	case VariableTypeByte:
		this.Data = this.getLongValue()
		this.Type = ExpressionTypeByte
	case VariableTypeShort:
		this.Data = this.getLongValue()
		this.Type = ExpressionTypeShort
	case VariableTypeChar:
		this.Data = this.getLongValue()
		this.Type = ExpressionTypeChar
	case VariableTypeInt:
		this.Data = this.getLongValue()
		this.Type = ExpressionTypeInt
	case VariableTypeLong:
		this.Data = this.getLongValue()
		this.Type = ExpressionTypeLong
	case VariableTypeFloat:
		this.Data = float32(this.getDoubleValue())
		this.Type = ExpressionTypeFloat
	case VariableTypeDouble:
		this.Data = this.getDoubleValue()
		this.Type = ExpressionTypeDouble
	}
}
