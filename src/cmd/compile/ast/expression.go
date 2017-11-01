package ast

import (
	"errors"
	"fmt"
)

const (
	EXPRESSION_TYPE_BOOL = iota
	EXPRESSION_TYPE_BYTE
	EXPRESSION_TYPE_INT
	EXPRESSION_TYPE_FLOAT
	EXPRESSION_TYPE_STRING
	EXPRESSION_TYPE_ARRAY
	EXPRESSION_TYPE_LOGICAL_OR
	EXPRESSION_TYPE_LOGICAL_AND
	EXPRESSION_TYPE_OR
	EXPRESSION_TYPE_AND
	EXPRESSION_TYPE_LEFT_SHIFT
	EXPRESSION_TYPE_RIGHT_SHIFT
	EXPRESSION_TYPE_ASSIGN
	EXPRESSION_TYPE_COLON_ASSIGN
	EXPRESSION_TYPE_PLUS_ASSIGN
	EXPRESSION_TYPE_MINUS_ASSIGN
	EXPRESSION_TYPE_MUL_ASSIGN
	EXPRESSION_TYPE_DIV_ASSIGN
	EXPRESSION_TYPE_MOD_ASSIGN
	EXPRESSION_TYPE_ASSIGN_FUNCTION
	EXPRESSION_TYPE_FUNCTION
	EXPRESSION_TYPE_EQ
	EXPRESSION_TYPE_NE
	EXPRESSION_TYPE_GE
	EXPRESSION_TYPE_GT
	EXPRESSION_TYPE_LE
	EXPRESSION_TYPE_LT
	EXPRESSION_TYPE_ADD
	EXPRESSION_TYPE_SUB
	EXPRESSION_TYPE_MUL
	EXPRESSION_TYPE_DIV
	EXPRESSION_TYPE_MOD
	EXPRESSION_TYPE_INDEX
	EXPRESSION_TYPE_METHOD_CALL
	EXPRESSION_TYPE_FUNCTION_CALL
	EXPRESSION_TYPE_INCREMENT
	EXPRESSION_TYPE_PRE_INCREMENT
	EXPRESSION_TYPE_PRE_DECREMENT
	EXPRESSION_TYPE_DECREMENT
	EXPRESSION_TYPE_NEGATIVE
	EXPRESSION_TYPE_NOT
	EXPRESSION_TYPE_IDENTIFIER
	EXPRESSION_TYPE_NULL
	EXPRESSION_TYPE_NEW
	EXPRESSION_TYPE_VAR
)

type Expression struct {
	Typ int
	/*
		BoolValue       bool
		IntValue        int64
		ByteValue       byte
		FloatValue      float64
		StringValue     string
		LeftExpression  *Expression
		RIghtExpression *Expression
	*/
	Data interface{} //
}

type ExpressionUnary Expression

type ExpressionFunctionCall struct {
	Name string //function name
	Args CallArgs
}
type ExpressionMethodCall struct {
	ClassName string
	ExpressionFunctionCall
}
type ExpressionBinary struct {
	Left  *Expression
	Right *Expression
}

//var (
//	ValueIsNotAConst = errors.New("value is not a const")
//)

func (binary *ExpressionBinary) getBinaryConstExpression() (is1 bool, typ1 int, value1 interface{}, err1 error, is2 bool, typ2 int, value2 interface{}, err2 error) {
	is1, typ1, value1, err1 = binary.Left.getConstValue()
	is2, typ2, value2, err2 = binary.Right.getConstValue()
	return
}

type getBinaryExpressionHandler func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int, value2 interface{}) (is bool, Typ int, Value interface{}, err error)

func (e *Expression) isNumber() bool {
	return e.Typ == EXPRESSION_TYPE_BYTE || e.Typ == EXPRESSION_TYPE_INT || e.Typ == EXPRESSION_TYPE_FLOAT
}
func expressionIsNumber(typ int) bool {
	return typ == EXPRESSION_TYPE_BYTE || typ == EXPRESSION_TYPE_INT || typ == EXPRESSION_TYPE_FLOAT
}

func (e *Expression) getBinaryExpressionConstValue(f getBinaryExpressionHandler) (is bool, Typ int, Value interface{}, err error) {
	binary := e.Data.(*ExpressionBinary)
	is1, typ1, value1, err1, is2, typ2, value2, err2 := binary.getBinaryConstExpression()
	if err1 != nil { //something is wrong
		err = err1
		return
	}
	if err2 != nil {
		err = err2
		return
	}
	return f(is1, typ1, value1, is2, typ2, value2)
}

//byte -> int
func (e *Expression) typeWider(typ1, typ2 int, value1, value2 interface{}) (t1 int, t2 int, v1 interface{}, v2 interface{}) { //
	if typ1 == typ2 {
		return typ1, typ2, value1, value2
	}
	if typ1 > typ2 {
		t1, t2 = typ1, typ1

	} else {
		t1, t2 = typ2, typ2
	}
	if t1 == typ1 { //typ1 has is wider
		v2 = e.typeConvertor(typ1, typ2, value2)
		v1 = value1
	} else {
		v1 = e.typeConvertor(typ2, typ1, value1)
		v2 = value2
	}
	return
}

func (e *Expression) typeConvertor(target int, origin int, v interface{}) interface{} {
	if target == EXPRESSION_TYPE_INT {
		switch origin {
		case EXPRESSION_TYPE_BYTE:
			return int64(v.(byte))
		case EXPRESSION_TYPE_INT:
			return v.(int64)
		case EXPRESSION_TYPE_FLOAT:
			panic("convert int to float")
		}
	}
	if target == EXPRESSION_TYPE_FLOAT {
		switch origin {
		case EXPRESSION_TYPE_BYTE:
			return int64(v.(byte))
		case EXPRESSION_TYPE_INT:
			return v.(int64)
		case EXPRESSION_TYPE_FLOAT:
			return v.(float64)
		}
	}
	panic(fmt.Sprintf("targt[%d] origin[%d] not handled", target, origin))
}

func float32IsZero(f float32) bool {
	return f < small_float && f > (-small_float)
}
func float64IsZero(f float64) bool {
	return f < small_float && f > (-small_float)
}

func (e *Expression) getConstValue() (is bool, Typ int, Value interface{}, err error) {
	if e.Typ == EXPRESSION_TYPE_BOOL ||
		e.Typ == EXPRESSION_TYPE_BYTE ||
		e.Typ == EXPRESSION_TYPE_INT ||
		e.Typ == EXPRESSION_TYPE_FLOAT ||
		e.Typ == EXPRESSION_TYPE_STRING {
		return true, e.Typ, e.Data, nil
	}
	// && and ||
	if e.Typ == EXPRESSION_TYPE_LOGICAL_AND || e.Typ == EXPRESSION_TYPE_LOGICAL_OR {
		return e.getBinaryExpressionConstValue(func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int, value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if is1 == false || is2 == false {
				is = false
				return
			}
			if typ1 != EXPRESSION_TYPE_BOOL || typ2 != EXPRESSION_TYPE_BOOL {
				err = errors.New("logical operation must apply to logical expressions")
				return
			}
			is = true
			Typ = EXPRESSION_TYPE_BOOL
			if e.Typ == EXPRESSION_TYPE_LOGICAL_AND {
				Value = value1.(bool) && value2.(bool)
			} else {
				Value = value1.(bool) || value2.(bool)
			}
			err = nil
			return
		})
	}
	// + - * / % algebra arithmetic
	if e.Typ == EXPRESSION_TYPE_ADD ||
		e.Typ == EXPRESSION_TYPE_SUB ||
		e.Typ == EXPRESSION_TYPE_MUL ||
		e.Typ == EXPRESSION_TYPE_DIV ||
		e.Typ == EXPRESSION_TYPE_MOD {
		return e.getBinaryExpressionConstValue(func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int, value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if is1 == false || is2 == false {
				is = false
				return
			}
			//string
			if typ1 == EXPRESSION_TYPE_STRING && typ2 == EXPRESSION_TYPE_STRING {
				is = true
				Typ = EXPRESSION_TYPE_STRING
				Value = value1.(string) + value2.(string)
				err = nil
				return
			}
			if expressionIsNumber(typ1) == false || expressionIsNumber(typ2) == false {
				err = errors.New("algebra operation must apply to number expressions or string+string or 'a'+'c' or 'a'-'c' ")
				return
			}
			typ1, typ2, value1, value2 = e.typeWider(typ1, typ2, value1, value2)
			if typ1 == EXPRESSION_TYPE_BYTE {
				is = true
				Typ = EXPRESSION_TYPE_BYTE
				switch e.Typ {
				case EXPRESSION_TYPE_ADD:
					Value = value1.(byte) + value2.(byte)
				case EXPRESSION_TYPE_SUB:
					Value = value1.(byte) - value2.(byte)
				case EXPRESSION_TYPE_MUL:
					Value = value1.(byte) * value2.(byte)
				case EXPRESSION_TYPE_DIV:
					if value2.(byte) == 0 {
						is = false
						err = fmt.Errorf("dividend is 0")
						return
					}
					Value = value1.(byte) / value2.(byte)
				case EXPRESSION_TYPE_MOD:
					if value2.(byte) == 0 {
						is = false
						err = fmt.Errorf("mod number is 0")
						return
					}
					Value = value1.(byte) % value2.(byte)
				}
				return
			}
			if typ1 == EXPRESSION_TYPE_INT {
				is = true
				Typ = EXPRESSION_TYPE_INT
				switch e.Typ {
				case EXPRESSION_TYPE_ADD:
					Value = value1.(int64) + value2.(int64)
				case EXPRESSION_TYPE_SUB:
					Value = value1.(int64) - value2.(int64)
				case EXPRESSION_TYPE_MUL:
					Value = value1.(int64) * value2.(int64)
				case EXPRESSION_TYPE_DIV:
					if value2.(int64) == 0 {
						is = false
						err = fmt.Errorf("dividend is 0")
						return
					}
					Value = value1.(int64) / value2.(int64)
				case EXPRESSION_TYPE_MOD:
					if value2.(int64) == 0 {
						is = false
						err = fmt.Errorf("mod number is 0")
						return
					}
					Value = value1.(int64) % value2.(int64)
				}
				return
			}
			if typ1 == EXPRESSION_TYPE_FLOAT {
				is = true
				Typ = EXPRESSION_TYPE_FLOAT
				switch e.Typ {
				case EXPRESSION_TYPE_ADD:
					Value = value1.(float64) + value2.(float64)
				case EXPRESSION_TYPE_SUB:
					Value = value1.(float64) - value2.(float64)
				case EXPRESSION_TYPE_MUL:
					Value = value1.(float64) * value2.(float64)
				case EXPRESSION_TYPE_DIV:
					if float64IsZero(value2.(float64)) == 0 {
						is = false
						err = fmt.Errorf("dividend is 0")
						return
					}
					Value = value1.(float64) / value2.(float64)
				case EXPRESSION_TYPE_MOD:
					if float64IsZero(value2.(float64)) == 0 {
						is = false
						err = fmt.Errorf("mod number is 0")
						return
					}
					Value = value1.(float64) % value2.(float64)
				}
				return
			}
			return
		})
	}
	// <<  >>
	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		return e.getBinaryExpressionConstValue(func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int, value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if is1 == false || is2 == false {
				is = false
				return
			}
			if typ1 != EXPRESSION_TYPE_INT || typ2 != EXPRESSION_TYPE_INT {
				err = errors.New("<< and >> operation must apply to number expressions,like 1<<10")
				return
			}
			is = true
			Typ = EXPRESSION_TYPE_INT
			err = nil
			if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT {
				Value = value1.(int64) << value2.(int64)
			} else {
				Value = value1.(int64) >> value2.(int64)
			}
			return
		})
	}
	// & |
	if e.Typ == EXPRESSION_TYPE_AND || e.Typ == EXPRESSION_TYPE_OR {
		return e.getBinaryExpressionConstValue(func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int, value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if is1 == false || is2 == false {
				is = false
				return
			}
			if (typ1 != EXPRESSION_TYPE_INT && typ1 != EXPRESSION_TYPE_BYTE) ||
				(typ2 != EXPRESSION_TYPE_INT && typ2 != EXPRESSION_TYPE_BYTE) {
				err = errors.New("& and | operation must apply to number expressions and byte")
				return
			}
			typ1, typ2, value1, value2 = e.typeWider(typ1, typ2, value1, value2)
			if typ1 == EXPRESSION_TYPE_INT {
				is = true
				Typ = EXPRESSION_TYPE_INT
				err = nil
				if EXPRESSION_TYPE_AND == e.Typ {
					e.Data = value1.(int64) & value2.(int64)
				} else {
					e.Data = value1.(int64) | value2.(int64)
				}
				return
			}
			if typ1 == EXPRESSION_TYPE_BYTE {
				is = true
				Typ = EXPRESSION_TYPE_BYTE
				err = nil
				if EXPRESSION_TYPE_AND == e.Typ {
					e.Data = value1.(byte) & value2.(byte)
				} else {
					e.Data = value1.(byte) | value2.(byte)
				}
				return
			}
			is = false
			return
		})
	}
	if e.Typ == EXPRESSION_TYPE_NOT {
		is, Typ, Value, err = Expression(e.Data.(*ExpressionUnary)).getConstValue()
		if err != nil {
			return
		}
		if is == false {
			return
		}
		if Typ != EXPRESSION_TYPE_BOOL {
			err = fmt.Errorf("!(not) can only apply to bool expression")
		}
		is = true
		Value = !Value.(bool)
		return

	}
	//  == !=
	if e.Typ == EXPRESSION_TYPE_EQ ||
		EXPRESSION_TYPE_NE == e.Typ {
		return e.getBinaryExpressionConstValue(func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int, value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if is1 == false || is2 == false {
				is = false
				return
			}
			if typ1 == EXPRESSION_TYPE_BOOL && typ2 == EXPRESSION_TYPE_BOOL {
				is = true
				Typ = EXPRESSION_TYPE_BOOL
				if EXPRESSION_TYPE_EQ == e.Typ {
					Value = value1.(bool) == value2.(bool)
				} else {
					Value = value1.(bool) != value2.(bool)
				}
			}

			return
		})

	}

	is = false
	return
}

type CallArgs []*Expression // f(1,2)　调用参数列表
