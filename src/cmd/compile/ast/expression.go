package ast

import (
	"errors"
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

var (
	ValueIsNotAConst = errors.New("value is not a const")
)

func (binary *ExpressionBinary) getBinaryConstExpression() (is1 bool, typ1 int, value1 interface{}, err1 error, is2 bool, typ2 int, value2 interface{}, err2 error) {
	is1, typ1, value1, err1 = binary.Left.getConstValue()
	is2, typ2, value2, err2 = binary.Right.getConstValue()
	return
}

//
//func (e *Expression) getConstValueLogicExpression() (is bool, Typ int, Value interface{}, err error) {
//	binary := e.Data.(*ExpressionBinary)
//	is1, typ1, value1, err1, is2, typ2, value2, err2 := binary.getBinaryConstExpression()
//	if err1 != nil { //something is wrong
//		err = err1
//		return
//	}
//	if err2 != nil {
//		err = err2
//		return
//	}
//	if is1 == false || is2 == false {
//		is = false
//		return
//	}
//	if typ1 != EXPRESSION_TYPE_BOOL || typ2 != EXPRESSION_TYPE_BOOL {
//		err = errors.New("logical operation must apply to logical expression")
//		return
//	}
//	is = true
//	Typ = EXPRESSION_TYPE_BOOL
//	Value = value1.(bool) && value2.(bool)
//	err = nil
//	return
//}

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
			if expressionIsNumber(typ1) == false || expressionIsNumber(typ2) == false {
				err = errors.New("algebra operation must apply to number expressions")
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

	is = false
	return

}

type CallArgs []*Expression // f(1,2)　调用参数列表
