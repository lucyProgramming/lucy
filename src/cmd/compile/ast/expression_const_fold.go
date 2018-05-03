package ast

import (
	"fmt"
)

func (e *Expression) getBinaryExpressionConstValue(f binaryConstFolder) (is bool, err error) {
	bin := e.Data.(*ExpressionBinary)
	is1, err1 := bin.Left.getConstValue()
	is2, err2 := bin.Right.getConstValue()
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

func (e *Expression) wrongOpErr(typ1, typ2 string) error {
	return fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
		errMsgPrefix(e.Pos),
		e.OpName(),
		typ1,
		typ2)
}

func (e *Expression) getConstValue() (is bool, err error) {
	if e.IsLiteral() {
		return true, nil
	}
	// !
	if e.Typ == EXPRESSION_TYPE_NOT {
		ee := e.Data.(*Expression)
		is, err = ee.getConstValue()
		if err != nil || is == false {
			return
		}
		if ee.Typ != EXPRESSION_TYPE_BOOL {
			err = fmt.Errorf("%s cannot apply '!' on a non-bool expression", errMsgPrefix(e.Pos))
			return
		}
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data = !ee.Data.(bool)
		return
	}
	if e.Typ == EXPRESSION_TYPE_NEGATIVE {
		ee := e.Data.(*Expression)
		is, err = ee.getConstValue()
		if err != nil || is == false {
			return
		}
		if ee.isNumber() == false {
			is = false
			err = fmt.Errorf("%s cannot apply '-' on '%s'", errMsgPrefix(e.Pos), ee.OpName())
			return
		}
		e.Typ = ee.Typ
		switch ee.Typ {
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
	if e.Typ == EXPRESSION_TYPE_LOGICAL_AND || e.Typ == EXPRESSION_TYPE_LOGICAL_OR {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.Typ != EXPRESSION_TYPE_BOOL ||
				bin.Right.Typ != EXPRESSION_TYPE_BOOL {
				err = e.wrongOpErr(bin.Left.OpName(), bin.Right.OpName())
				return
			}
			is = true
			e.Typ = EXPRESSION_TYPE_BOOL
			if e.Typ == EXPRESSION_TYPE_LOGICAL_AND {
				e.Data = bin.Left.Data.(bool) && bin.Right.Data.(bool)
			} else {
				e.Data = bin.Left.Data.(bool) || bin.Right.Data.(bool)
			}
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
		is, err = e.getBinaryExpressionConstValue(e.arithmeticBinayConstFolder)
		return
	}
	// <<  >>
	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.isInteger() == false || bin.Right.isInteger() == false {
				return
			}
			switch bin.Left.Typ {
			case EXPRESSION_TYPE_BYTE:
				if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
					e.Data = byte(bin.Left.Data.(byte) << bin.Right.getByteValue())
				} else {
					e.Data = byte(bin.Left.Data.(byte) >> bin.Right.getByteValue())
				}
			case EXPRESSION_TYPE_SHORT:
				if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
					e.Data = int32(bin.Left.Data.(int32) << bin.Right.getByteValue())
				} else {
					e.Data = int32(bin.Left.Data.(int32) >> bin.Right.getByteValue())
				}
			case EXPRESSION_TYPE_INT:
				if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
					e.Data = int32(bin.Left.Data.(int32) << bin.Right.getByteValue())
				} else {
					e.Data = int32(bin.Left.Data.(int32) >> bin.Right.getByteValue())
				}
			case EXPRESSION_TYPE_LONG:
				if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
					e.Data = int64(bin.Left.Data.(int64) << bin.Right.getByteValue())
				} else {
					e.Data = int64(bin.Left.Data.(int64) >> bin.Right.getByteValue())
				}
			}
			e.Typ = bin.Left.Typ
			return
		}

		return e.getBinaryExpressionConstValue(f)
	}
	// & | ^
	if e.Typ == EXPRESSION_TYPE_AND ||
		e.Typ == EXPRESSION_TYPE_OR ||
		e.Typ == EXPRESSION_TYPE_XOR {
		f := func(bin *ExpressionBinary) (is bool, err error) {
			if bin.Left.isInteger() == false || bin.Right.isInteger() == false ||
				bin.Left.Typ != bin.Right.Typ {
				return // not integer or type not equal
			}
			e.Typ = bin.Left.Typ
			switch bin.Left.Typ {
			case EXPRESSION_TYPE_BYTE:
				if e.Typ == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(byte) & bin.Right.Data.(byte)
				} else if e.Typ == EXPRESSION_TYPE_OR {
					e.Data = bin.Left.Data.(byte) | bin.Right.Data.(byte)
				} else {
					e.Data = bin.Left.Data.(byte) ^ bin.Right.Data.(byte)
				}
			case EXPRESSION_TYPE_SHORT:
				if e.Typ == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(int32) & bin.Right.Data.(int32)
				} else if e.Typ == EXPRESSION_TYPE_OR {
					e.Data = bin.Left.Data.(int32) | bin.Right.Data.(int32)
				} else {
					e.Data = bin.Left.Data.(int32) ^ bin.Right.Data.(int32)
				}
			case EXPRESSION_TYPE_INT:
				if e.Typ == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(int32) & bin.Right.Data.(int32)
				} else if e.Typ == EXPRESSION_TYPE_OR {
					e.Data = bin.Left.Data.(int32) | bin.Right.Data.(int32)
				} else {
					e.Data = bin.Left.Data.(int32) ^ bin.Right.Data.(int32)
				}
			case EXPRESSION_TYPE_LONG:
				if e.Typ == EXPRESSION_TYPE_AND {
					e.Data = bin.Left.Data.(int64) & bin.Right.Data.(int64)
				} else if e.Typ == EXPRESSION_TYPE_OR {
					e.Data = bin.Left.Data.(int64) | bin.Right.Data.(int64)
				} else {
					e.Data = bin.Left.Data.(int64) ^ bin.Right.Data.(int64)
				}
			}
			is = true
			return
		}
		return e.getBinaryExpressionConstValue(f)
	}
	if e.Typ == EXPRESSION_TYPE_NOT {
		ee := e.Data.(*Expression)
		is, err = ee.getConstValue()
		if err != nil {
			return
		}
		if is == false {
			return
		}
		if ee.Typ != EXPRESSION_TYPE_BOOL {
			return false, fmt.Errorf("!(not) can only apply to bool expression")
		}
		is = true
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data = !ee.Data.(bool)
		return
	}
	//  == != > < >= <=
	if e.Typ == EXPRESSION_TYPE_EQ ||
		e.Typ == EXPRESSION_TYPE_NE ||
		e.Typ == EXPRESSION_TYPE_GE ||
		e.Typ == EXPRESSION_TYPE_GT ||
		e.Typ == EXPRESSION_TYPE_LE ||
		e.Typ == EXPRESSION_TYPE_LT {
		return e.getBinaryExpressionConstValue(e.relationBinayConstFolder)
	}
	return
}
