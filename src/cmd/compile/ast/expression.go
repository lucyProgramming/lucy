package ast

import (
	"errors"
	"fmt"
)

const (
	_ = iota
	EXPRESSION_TYPE_BOOL
	EXPRESSION_TYPE_BYTE
	EXPRESSION_TYPE_INT
	EXPRESSION_TYPE_FLOAT
	EXPRESSION_TYPE_STRING
	EXPRESSION_TYPE_ARRAY // []bool{false,true}
	//binary expression
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
	//
	EXPRESSION_TYPE_INDEX // a["b"]
	EXPRESSION_TYPE_DOT   //a.b
	EXPRESSION_TYPE_METHOD_CALL
	EXPRESSION_TYPE_FUNCTION_CALL
	EXPRESSION_TYPE_INCREMENT
	EXPRESSION_TYPE_DECREMENT
	EXPRESSION_TYPE_PRE_INCREMENT
	EXPRESSION_TYPE_PRE_DECREMENT
	EXPRESSION_TYPE_NEGATIVE
	EXPRESSION_TYPE_NOT
	EXPRESSION_TYPE_IDENTIFIER
	EXPRESSION_TYPE_NULL
	EXPRESSION_TYPE_NEW
	EXPRESSION_TYPE_LIST
	EXPRESSION_TYPE_FUNCTION
	EXPRESSION_TYPE_VAR
	EXPRESSION_TYPE_CONST
	EXPRESSION_TYPE_CONVERTION_TYPE // []byte(str)
)

//receiver only one argument
func (e *Expression) typeName(typ ...int) string {
	t := e.Typ
	if len(typ) > 0 {
		t = typ[0]
	}
	switch t {
	case EXPRESSION_TYPE_BOOL:
		return "bool"
	case EXPRESSION_TYPE_BYTE:
		return "byte"
	case EXPRESSION_TYPE_INT:
		return "int"
	case EXPRESSION_TYPE_FLOAT:
		return "float"
	case EXPRESSION_TYPE_STRING:
		return "string"
	case EXPRESSION_TYPE_EQ:
		return "equal"
	case EXPRESSION_TYPE_NE:
		return "not equal"
	case EXPRESSION_TYPE_GE:
		return "greater than"
	case EXPRESSION_TYPE_GT:
		return "greater or equal"
	case EXPRESSION_TYPE_LE:
		return "less or equal"
	case EXPRESSION_TYPE_LT:
		return "less"
	case EXPRESSION_TYPE_ADD:
		return "add(+)"
	case EXPRESSION_TYPE_SUB:
		return "sub(-)"
	case EXPRESSION_TYPE_MUL:
		return "multiply(*)"
	case EXPRESSION_TYPE_DIV:
		return "divide(/)"
	case EXPRESSION_TYPE_MOD:
		return "mod(%)"
	}
	return ""
}

type Expression struct {
	Pos  *Pos
	Typ  int
	Data interface{}
}

type CallArgs []*Expression // f(1,2)　调用参数列表

type ExpressionFunctionCall struct {
	Expression *Expression
	Args       CallArgs
	Func       *Function
}

type ExpressionDeclareVariable struct {
	Vs []*VariableDefinition
}

type ExpressionDeclareConsts struct {
	Cs []*Const
}

type ExpressionTypeConvertion struct {
	Typ        *VariableType
	Expression *Expression
}

type ExpressionMethodCall struct {
	ExpressionFunctionCall
	Name   string
	Method *ClassMethod
}

type ExpressionNew struct {
	Typ  *VariableType
	Args CallArgs
}

type ExpressionBinary struct {
	Left  *Expression
	Right *Expression
}

type ExpressionArray struct {
	Typ        *VariableType
	Expression *Expression
}

/*
	literal value to float64
*/
func (e *Expression) literalValue2Float64() int64 {
	switch e.Typ {
	case EXPRESSION_TYPE_BYTE:
		return int64(e.Data.(byte))
	case EXPRESSION_TYPE_INT:
		return e.Data.(int64)
	case EXPRESSION_TYPE_FLOAT:
		return int64(e.Data.(float64))
	default:
		panic("unhandle convert to int64")
	}
}

/*
	literal value to float64
*/
func (e *Expression) literalValue2Int64() float64 {
	switch e.Typ {
	case EXPRESSION_TYPE_BYTE:
		return float64(e.Data.(byte))
	case EXPRESSION_TYPE_INT:
		return float64(e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return e.Data.(float64)
	default:
		panic("unhandle convert to float64")
	}
}

func (e *Expression) canBeCovert2Bool() (bool, error) {
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		return e.Data.(bool), nil
	case EXPRESSION_TYPE_BYTE:
		return e.Data.(byte) != 0, nil
	case EXPRESSION_TYPE_INT:
		return e.Data.(int64) != 0, nil
	case EXPRESSION_TYPE_FLOAT:
		return !float64IsZero(e.Data.(float64)), nil
	default:
		return false, fmt.Errorf("can not convert to bool")
	}
}

func (e *Expression) OpName() string {
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		return fmt.Sprintf("bool_value", e.Data.(bool))
	case EXPRESSION_TYPE_BYTE:
		return fmt.Sprintf("byte_value", e.Data.(byte))
	case EXPRESSION_TYPE_INT:
		return fmt.Sprintf("int_value", e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return fmt.Sprintf("float_value", e.Data.(float64))
	case EXPRESSION_TYPE_STRING:
		t := []byte(e.Data.(string))
		if len(t) > 10 {
			t = t[0:10]
		}
		t = append(t, []byte("...")...)
		return fmt.Sprintf("string(%s)", string(t))
	case EXPRESSION_TYPE_ARRAY:
		return "array_literal"
	case EXPRESSION_TYPE_LOGICAL_OR:
		return "||"
	case EXPRESSION_TYPE_LOGICAL_AND:
		return "&&"
	case EXPRESSION_TYPE_OR:
		return "|"
	case EXPRESSION_TYPE_AND:
		return "&&"
	case EXPRESSION_TYPE_LEFT_SHIFT:
		return "<<"
	case EXPRESSION_TYPE_RIGHT_SHIFT:
		return ">>"
	case EXPRESSION_TYPE_ASSIGN:
		return "="
	case EXPRESSION_TYPE_COLON_ASSIGN:
		return ":="
	case EXPRESSION_TYPE_PLUS_ASSIGN:
		return "+="
	case EXPRESSION_TYPE_MINUS_ASSIGN:
		return "-="
	case EXPRESSION_TYPE_MUL_ASSIGN:
		return "*="
	case EXPRESSION_TYPE_DIV_ASSIGN:
		return "/="
	case EXPRESSION_TYPE_MOD_ASSIGN:
		return "%="
	case EXPRESSION_TYPE_EQ:
		return "=="
	case EXPRESSION_TYPE_NE:
		return "!="
	case EXPRESSION_TYPE_GE:
		return ">="
	case EXPRESSION_TYPE_GT:
		return ">"
	case EXPRESSION_TYPE_LE:
		return "<="
	case EXPRESSION_TYPE_LT:
		return "<"
	case EXPRESSION_TYPE_ADD:
		return "+"
	case EXPRESSION_TYPE_SUB:
		return "-"
	case EXPRESSION_TYPE_MUL:
		return "*"
	case EXPRESSION_TYPE_DIV:
		return "/"
	case EXPRESSION_TYPE_MOD:
		return "%"
	case EXPRESSION_TYPE_INDEX: // a["b"]
		return "[]"
	case EXPRESSION_TYPE_DOT: //a.b
		return "."
	case EXPRESSION_TYPE_METHOD_CALL:
		return "method_call"
	case EXPRESSION_TYPE_FUNCTION_CALL:
		return "function_call"
	case EXPRESSION_TYPE_INCREMENT:
		return "++"
	case EXPRESSION_TYPE_DECREMENT:
		return "--"
	case EXPRESSION_TYPE_PRE_INCREMENT:
		return "++"
	case EXPRESSION_TYPE_PRE_DECREMENT:
		return "--"
	case EXPRESSION_TYPE_NEGATIVE:
		return "nagative"
	case EXPRESSION_TYPE_NOT:
		return "not"
	case EXPRESSION_TYPE_IDENTIFIER:
		return "identifier"
	case EXPRESSION_TYPE_NULL:
		return "null"
	case EXPRESSION_TYPE_NEW:
		return "new"
	case EXPRESSION_TYPE_LIST:
		return "expression_list"
	case EXPRESSION_TYPE_FUNCTION:
		return "function_literal"
	case EXPRESSION_TYPE_CONST:
		return "const"
	case EXPRESSION_TYPE_VAR:
		return "var"
	}
	panic("missing type")
}

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
func (e *Expression) typeWider(typ1, typ2 int, value1, value2 interface{}) (t1 int, t2 int, v1 interface{}, v2 interface{}, err error) { //
	if typ1 == typ2 {
		return typ1, typ2, value1, value2, nil
	}
	if typ1 > typ2 {
		t1, t2 = typ1, typ1

	} else {
		t1, t2 = typ2, typ2
	}
	if t1 == typ1 { //typ1 has is wider
		v2, err = e.typeConvertor(typ1, typ2, value2)
		v1 = value1
	} else {
		v1, err = e.typeConvertor(typ2, typ1, value1)
		v2 = value2
	}
	if err == nil {
		return
	}
	return typ1, typ2, value1, value2, err
}

func (e *Expression) typeConvertor(target int, origin int, v interface{}) (interface{}, error) {
	if target == EXPRESSION_TYPE_INT {
		switch origin {
		case EXPRESSION_TYPE_BYTE:
			return int64(v.(byte)), nil
		case EXPRESSION_TYPE_INT:
			return v.(int64), nil
		}
	}
	if target == EXPRESSION_TYPE_FLOAT {
		switch origin {
		case EXPRESSION_TYPE_BYTE:
			return int64(v.(byte)), nil
		case EXPRESSION_TYPE_INT:
			return v.(int64), nil
		case EXPRESSION_TYPE_FLOAT:
			return v.(float64), nil
		}
	}
	return nil, fmt.Errorf("targt[%d] origin[%d] not handled", target, origin)
}

func float32IsZero(f float32) bool {
	return float64IsZero(float64(f))
}
func float64IsZero(f float64) bool {
	return f < small_float && f > negative_small_float
}

func (e *Expression) relationnalCompare(typ int, value1, value2 interface{}) (b bool, err error) {
	fmt.Println("$$$$$$$$$$$", typ)
	fmt.Println(value1)
	fmt.Println(value2)

	switch typ {
	case EXPRESSION_TYPE_BOOL:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(bool) == value2.(bool)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(bool) != value2.(bool)
		}
		return
	case EXPRESSION_TYPE_BYTE:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(byte) == value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(byte) != value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(byte) > value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(byte) >= value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(byte) < value2.(byte)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(byte) <= value2.(byte)
		}
		return

	case EXPRESSION_TYPE_INT:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(int64) == value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(int64) != value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(int64) > value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(int64) >= value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(int64) < value2.(int64)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(int64) <= value2.(int64)
		}
		return

	case EXPRESSION_TYPE_FLOAT:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(float64) == value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(float64) != value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(float64) > value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(float64) >= value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(float64) < value2.(float64)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(float64) <= value2.(float64)
		}
		return

	case EXPRESSION_TYPE_STRING:
		if e.Typ == EXPRESSION_TYPE_EQ {
			b = value1.(string) == value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_NE {
			b = value1.(string) != value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_GT {
			b = value1.(string) > value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_GE {
			b = value1.(string) >= value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_LT {
			b = value1.(string) < value2.(string)
		} else if e.Typ == EXPRESSION_TYPE_LE {
			b = value1.(string) <= value2.(string)
		}
		return

	}

	return false, fmt.Errorf("can`t compare")
}

////this function no need to assign back
//func (e *Expression) foldConst() error {
//	is, typ, value, err := e.getConstValue() //something is error
//	if err != nil {
//		return err
//	}
//	if is { //overide
//		e.Typ = typ
//		e.Data = value
//		return nil
//	}
//	//this means below is is not a const definitely
//	//binary expression
//	if e.Typ == EXPRESSION_TYPE_LOGICAL_AND ||
//		e.Typ == EXPRESSION_TYPE_LOGICAL_OR ||
//		e.Typ == EXPRESSION_TYPE_ADD ||
//		e.Typ == EXPRESSION_TYPE_SUB ||
//		e.Typ == EXPRESSION_TYPE_MUL ||
//		e.Typ == EXPRESSION_TYPE_DIV ||
//		e.Typ == EXPRESSION_TYPE_MOD ||
//		e.Typ == EXPRESSION_TYPE_LEFT_SHIFT ||
//		e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT ||
//		e.Typ == EXPRESSION_TYPE_AND ||
//		e.Typ == EXPRESSION_TYPE_OR ||
//		e.Typ == EXPRESSION_TYPE_EQ ||
//		e.Typ == EXPRESSION_TYPE_NE ||
//		e.Typ == EXPRESSION_TYPE_GE ||
//		e.Typ == EXPRESSION_TYPE_GT ||
//		e.Typ == EXPRESSION_TYPE_LE ||
//		e.Typ == EXPRESSION_TYPE_LE {
//		binary := e.Data.(*ExpressionBinary)
//		err = binary.Left.constFold()
//		if err != nil {
//			return err
//		}
//		err = binary.Right.constFold()
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

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
			typ1, typ2, value1, value2, err = e.typeWider(typ1, typ2, value1, value2)
			if err != nil {
				return
			}
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
					if float64IsZero(value2.(float64)) {
						is = false
						err = fmt.Errorf("dividend is 0")
						return
					}
					Value = value1.(float64) / value2.(float64)
				case EXPRESSION_TYPE_MOD:
					is = false
					err = fmt.Errorf("can`t not apply mod(%) on float")
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
				Value = value1.(int64) << uint64(value2.(int64))
			} else {
				Value = value1.(int64) >> uint64(value2.(int64))
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
			typ1, typ2, value1, value2, err = e.typeWider(typ1, typ2, value1, value2)
			if err != nil {
				return
			}
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
		t := e.Data.(*Expression)
		is, Typ, Value, err = t.getConstValue()
		if err != nil {
			return
		}
		if is == false {
			return
		}
		if Typ != EXPRESSION_TYPE_BOOL {
			err = fmt.Errorf("!(not) can only apply to bool expression")
		} else {
			is = true
			Value = !Value.(bool)
		}
		return

	}
	//  == != > < >= <=
	if e.Typ == EXPRESSION_TYPE_EQ ||
		e.Typ == EXPRESSION_TYPE_NE ||
		e.Typ == EXPRESSION_TYPE_GE ||
		e.Typ == EXPRESSION_TYPE_GT ||
		e.Typ == EXPRESSION_TYPE_LE ||
		e.Typ == EXPRESSION_TYPE_LE {
		return e.getBinaryExpressionConstValue(func(is1 bool, typ1 int, value1 interface{}, is2 bool, typ2 int, value2 interface{}) (is bool, Typ int, Value interface{}, err error) {
			if is1 == false || is2 == false {
				is = false
				return
			}
			typ1, typ2, value1, value2, err = e.typeWider(typ1, typ2, value1, value2)
			if err != nil {
				err = fmt.Errorf("relation operation cannot apply to %s and %s", e.typeName(typ1), e.typeName(typ2))
				return
			}
			b, er := e.relationnalCompare(typ1, value1, value2)
			if er != nil {
				err = er
				return
			}
			is = true
			Value = b
			err = nil
			Typ = EXPRESSION_TYPE_BOOL
			return
		})
	}
	is = false
	err = nil
	return
}
