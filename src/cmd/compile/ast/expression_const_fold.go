package ast

import (
	"fmt"
)

func (e *Expression) getBinaryExpressionConstValue(f binaryConstFolder) (is bool, err error) {
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
	if is1 == false || is2 == false {
		is = false
		err = nil
		return
	}
	return f(bin)
}

type binaryConstFolder func(bin *ExpressionBinary) (is bool, err error)

func (e *Expression) mkWrongOpErr(typ1, typ2 string) error {
	return fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
		errMsgPrefix(e.Pos),
		e.OpName(),
		typ1,
		typ2)
}

func (e *Expression) constantFold() (is bool, err error) {
	if e.IsLiteral() {
		return true, nil
	}
	// ~
	if e.Type == EXPRESSION_TYPE_BITWISE_NOT {
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
		case EXPRESSION_TYPE_BYTE:
			e.Data = ^ee.Data.(byte)
		case EXPRESSION_TYPE_SHORT:
			e.Data = ^ee.Data.(int32)
		case EXPRESSION_TYPE_INT:
			e.Data = ^ee.Data.(int32)
		case EXPRESSION_TYPE_LONG:
			e.Data = ^ee.Data.(int64)
		}
	}
	// !
	if e.Type == EXPRESSION_TYPE_NOT {
		ee := e.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		if ee.Type != EXPRESSION_TYPE_BOOL {
			err = fmt.Errorf("%s cannot apply '!' on a non-bool expression",
				errMsgPrefix(e.Pos))
			return
		}
		e.Type = EXPRESSION_TYPE_BOOL
		e.Data = !ee.Data.(bool)
		return
	}
	if e.Type == EXPRESSION_TYPE_NEGATIVE {
		ee := e.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil || is == false {
			return
		}
		if ee.isNumber() == false {
			is = false
			err = fmt.Errorf("%s cannot apply '-' on '%s'",
				errMsgPrefix(e.Pos), ee.OpName())
			return
		}
		e.Type = ee.Type
		switch ee.Type {
		case EXPRESSION_TYPE_BYTE:
			e.Data = -ee.Data.(byte)
		case EXPRESSION_TYPE_SHORT:
			e.Data = -ee.Data.(int32)
		case EXPRESSION_TYPE_INT:
			e.Data = -ee.Data.(int32)
		case EXPRESSION_TYPE_LONG:
			e.Data = -ee.Data.(int64)
		case EXPRESSION_TYPE_FLOAT:
			e.Data = -ee.Data.(float32)
		case EXPRESSION_TYPE_DOUBLE:
			e.Data = -ee.Data.(float64)
		}
		return
	}
	// && and ||
	if e.Type == EXPRESSION_TYPE_LOGICAL_AND || e.Type == EXPRESSION_TYPE_LOGICAL_OR {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.Type != EXPRESSION_TYPE_BOOL ||
				bin.Right.Type != EXPRESSION_TYPE_BOOL {
				err = e.mkWrongOpErr(bin.Left.OpName(), bin.Right.OpName())
				return
			}
			is = true
			e.Type = EXPRESSION_TYPE_BOOL
			if e.Type == EXPRESSION_TYPE_LOGICAL_AND {
				e.Data = bin.Left.Data.(bool) && bin.Right.Data.(bool)
			} else {
				e.Data = bin.Left.Data.(bool) || bin.Right.Data.(bool)
			}
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	// + - * / % algebra arithmetic
	if e.Type == EXPRESSION_TYPE_ADD ||
		e.Type == EXPRESSION_TYPE_SUB ||
		e.Type == EXPRESSION_TYPE_MUL ||
		e.Type == EXPRESSION_TYPE_DIV ||
		e.Type == EXPRESSION_TYPE_MOD {
		is, err = e.getBinaryExpressionConstValue(e.arithmeticBinaryConstFolder)
		return
	}
	// <<  >>
	if e.Type == EXPRESSION_TYPE_LSH || e.Type == EXPRESSION_TYPE_RSH {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.isInteger() == false || bin.Right.isInteger() == false {
				return
			}
			switch bin.Left.Type {
			case EXPRESSION_TYPE_BYTE:
				if e.Type == EXPRESSION_TYPE_LSH {
					e.Data = byte(bin.Left.Data.(byte) << bin.Right.getByteValue())
				} else {
					e.Data = byte(bin.Left.Data.(byte) >> bin.Right.getByteValue())
				}
			case EXPRESSION_TYPE_SHORT:
				if e.Type == EXPRESSION_TYPE_LSH {
					e.Data = int32(bin.Left.Data.(int32) << bin.Right.getByteValue())
				} else {
					e.Data = int32(bin.Left.Data.(int32) >> bin.Right.getByteValue())
				}
			case EXPRESSION_TYPE_INT:
				if e.Type == EXPRESSION_TYPE_LSH {
					e.Data = int32(bin.Left.Data.(int32) << bin.Right.getByteValue())
				} else {
					e.Data = int32(bin.Left.Data.(int32) >> bin.Right.getByteValue())
				}
			case EXPRESSION_TYPE_LONG:
				if e.Type == EXPRESSION_TYPE_LSH {
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
	if e.Type == EXPRESSION_TYPE_AND ||
		e.Type == EXPRESSION_TYPE_OR ||
		e.Type == EXPRESSION_TYPE_XOR {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.isInteger() == false || bin.Right.isInteger() == false ||
				bin.Left.Type != bin.Right.Type {
				return // not integer or type not equal
			}
			switch bin.Left.Type {
			case EXPRESSION_TYPE_BYTE:
				if e.Type == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(byte) & bin.Right.Data.(byte)
				} else if e.Type == EXPRESSION_TYPE_OR {
					e.Data = bin.Left.Data.(byte) | bin.Right.Data.(byte)
				} else {
					e.Data = bin.Left.Data.(byte) ^ bin.Right.Data.(byte)
				}
			case EXPRESSION_TYPE_SHORT:
				if e.Type == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(int32) & bin.Right.Data.(int32)
				} else if e.Type == EXPRESSION_TYPE_OR {
					e.Data = bin.Left.Data.(int32) | bin.Right.Data.(int32)
				} else {
					e.Data = bin.Left.Data.(int32) ^ bin.Right.Data.(int32)
				}
			case EXPRESSION_TYPE_INT:
				if e.Type == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(int32) & bin.Right.Data.(int32)
				} else if e.Type == EXPRESSION_TYPE_OR {
					e.Data = bin.Left.Data.(int32) | bin.Right.Data.(int32)
				} else {
					e.Data = bin.Left.Data.(int32) ^ bin.Right.Data.(int32)
				}
			case EXPRESSION_TYPE_LONG:
				if e.Type == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if e.Type == EXPRESSION_TYPE_OR {
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
	if e.Type == EXPRESSION_TYPE_NOT {
		ee := e.Data.(*Expression)
		is, err = ee.constantFold()
		if err != nil {
			return
		}
		if is == false {
			return
		}
		if ee.Type != EXPRESSION_TYPE_BOOL {
			return false, fmt.Errorf("!(not) can only apply to bool expression")
		}
		is = true
		e.Type = EXPRESSION_TYPE_BOOL
		e.Data = !ee.Data.(bool)
		return
	}
	//  == != > < >= <=
	if e.Type == EXPRESSION_TYPE_EQ ||
		e.Type == EXPRESSION_TYPE_NE ||
		e.Type == EXPRESSION_TYPE_GE ||
		e.Type == EXPRESSION_TYPE_GT ||
		e.Type == EXPRESSION_TYPE_LE ||
		e.Type == EXPRESSION_TYPE_LT {
		return e.getBinaryExpressionConstValue(e.relationBinaryConstFolder)
	}
	return
}

func (e *Expression) getByteValue() byte {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case EXPRESSION_TYPE_BYTE:
		return e.Data.(byte)
	case EXPRESSION_TYPE_SHORT:
		fallthrough
	case EXPRESSION_TYPE_INT:
		return byte(e.Data.(int32))
	case EXPRESSION_TYPE_LONG:
		return byte(e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return byte(e.Data.(float32))
	case EXPRESSION_TYPE_DOUBLE:
		return byte(e.Data.(float64))
	}
	return 0
}
func (e *Expression) getShortValue() int32 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case EXPRESSION_TYPE_BYTE:
		return int32(e.Data.(byte))
	case EXPRESSION_TYPE_SHORT:
		fallthrough
	case EXPRESSION_TYPE_INT:
		return int32(e.Data.(int32))
	case EXPRESSION_TYPE_LONG:
		return int32(e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return int32(e.Data.(float32))
	case EXPRESSION_TYPE_DOUBLE:
		return int32(e.Data.(float64))
	}
	return 0
}
func (e *Expression) getIntValue() int32 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case EXPRESSION_TYPE_BYTE:
		return int32(e.Data.(byte))
	case EXPRESSION_TYPE_SHORT:
		fallthrough
	case EXPRESSION_TYPE_INT:
		return int32(e.Data.(int32))
	case EXPRESSION_TYPE_LONG:
		return int32(e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return int32(e.Data.(float32))
	case EXPRESSION_TYPE_DOUBLE:
		return int32(e.Data.(float64))
	}
	return 0
}

func (e *Expression) getLongValue() int64 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case EXPRESSION_TYPE_BYTE:
		return int64(e.Data.(byte))
	case EXPRESSION_TYPE_SHORT:
		fallthrough
	case EXPRESSION_TYPE_INT:
		return int64(e.Data.(int32))
	case EXPRESSION_TYPE_LONG:
		return int64(e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return int64(e.Data.(float32))
	case EXPRESSION_TYPE_DOUBLE:
		return int64(e.Data.(float64))
	}
	return 0
}
func (e *Expression) getFloatValue() float32 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case EXPRESSION_TYPE_BYTE:
		return float32(e.Data.(byte))
	case EXPRESSION_TYPE_SHORT:
		fallthrough
	case EXPRESSION_TYPE_INT:
		return float32(e.Data.(int32))
	case EXPRESSION_TYPE_LONG:
		return float32(e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return float32(e.Data.(float32))
	case EXPRESSION_TYPE_DOUBLE:
		return float32(e.Data.(float64))
	}
	return 0
}

func (e *Expression) getDoubleValue() float64 {
	if e.isNumber() == false {
		panic("not number")
	}
	switch e.Type {
	case EXPRESSION_TYPE_BYTE:
		return float64(e.Data.(byte))
	case EXPRESSION_TYPE_SHORT:
		fallthrough
	case EXPRESSION_TYPE_INT:
		return float64(e.Data.(int32))
	case EXPRESSION_TYPE_LONG:
		return float64(e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return float64(e.Data.(float32))
	case EXPRESSION_TYPE_DOUBLE:
		return float64(e.Data.(float64))
	}
	return 0
}

func (e *Expression) convertNumberLiteralTo(t int) {
	if e.isNumber() == false {
		panic("...")
	}
	switch t {
	case VARIABLE_TYPE_BYTE:
		e.Data = e.getByteValue()
		e.Type = EXPRESSION_TYPE_BYTE
	case VARIABLE_TYPE_SHORT:
		e.Data = e.getShortValue()
		e.Type = EXPRESSION_TYPE_SHORT
	case VARIABLE_TYPE_INT:
		e.Data = e.getIntValue()
		e.Type = EXPRESSION_TYPE_INT
	case VARIABLE_TYPE_LONG:
		e.Data = e.getLongValue()
		e.Type = EXPRESSION_TYPE_LONG
	case VARIABLE_TYPE_FLOAT:
		e.Data = e.getFloatValue()
		e.Type = EXPRESSION_TYPE_FLOAT
	case VARIABLE_TYPE_DOUBLE:
		e.Data = e.getDoubleValue()
		e.Type = EXPRESSION_TYPE_DOUBLE
	}
}
