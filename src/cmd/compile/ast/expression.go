package ast

import (
	"fmt"
)

const (
	_ = iota
	//value type
	EXPRESSION_TYPE_NULL
	EXPRESSION_TYPE_BOOL
	EXPRESSION_TYPE_BYTE
	EXPRESSION_TYPE_INT
	EXPRESSION_TYPE_FLOAT
	EXPRESSION_TYPE_DOUBLE
	EXPRESSION_TYPE_LONG
	EXPRESSION_TYPE_STRING
	EXPRESSION_TYPE_ARRAY // []bool{false,true}
	//binary expression
	EXPRESSION_TYPE_LOGICAL_OR
	EXPRESSION_TYPE_LOGICAL_AND
	//
	EXPRESSION_TYPE_OR
	EXPRESSION_TYPE_AND
	EXPRESSION_TYPE_LEFT_SHIFT
	EXPRESSION_TYPE_RIGHT_SHIFT
	EXPRESSION_TYPE_ADD
	EXPRESSION_TYPE_SUB
	EXPRESSION_TYPE_MUL
	EXPRESSION_TYPE_DIV
	EXPRESSION_TYPE_MOD
	//
	EXPRESSION_TYPE_ASSIGN
	EXPRESSION_TYPE_COLON_ASSIGN
	//
	EXPRESSION_TYPE_PLUS_ASSIGN
	EXPRESSION_TYPE_MINUS_ASSIGN
	EXPRESSION_TYPE_MUL_ASSIGN
	EXPRESSION_TYPE_DIV_ASSIGN
	EXPRESSION_TYPE_MOD_ASSIGN
	//
	EXPRESSION_TYPE_EQ
	EXPRESSION_TYPE_NE
	EXPRESSION_TYPE_GE
	EXPRESSION_TYPE_GT
	EXPRESSION_TYPE_LE
	EXPRESSION_TYPE_LT
	//

	//
	EXPRESSION_TYPE_INDEX // a["b"]
	EXPRESSION_TYPE_DOT   //a.b
	//
	EXPRESSION_TYPE_METHOD_CALL
	EXPRESSION_TYPE_FUNCTION_CALL
	//
	EXPRESSION_TYPE_INCREMENT
	EXPRESSION_TYPE_DECREMENT
	EXPRESSION_TYPE_PRE_INCREMENT
	EXPRESSION_TYPE_PRE_DECREMENT
	//
	EXPRESSION_TYPE_NEGATIVE
	EXPRESSION_TYPE_NOT
	//
	EXPRESSION_TYPE_IDENTIFIER
	EXPRESSION_TYPE_NEW
	EXPRESSION_TYPE_LIST
	EXPRESSION_TYPE_FUNCTION
	EXPRESSION_TYPE_VAR
	EXPRESSION_TYPE_CONST
	EXPRESSION_TYPE_CONVERTION_TYPE // []byte(str)
)

func (e *Expression) IsLiteral() bool {
	return e.Typ == EXPRESSION_TYPE_NULL ||
		e.Typ == EXPRESSION_TYPE_BOOL ||
		e.Typ == EXPRESSION_TYPE_BYTE ||
		e.Typ == EXPRESSION_TYPE_INT ||
		e.Typ == EXPRESSION_TYPE_FLOAT ||
		e.Typ == EXPRESSION_TYPE_STRING
}

type Expression struct {
	VariableType          *VariableType   //
	VariableTypes         []*VariableType // functioncall or methodcall can with multi results
	Pos                   *Pos
	Typ                   int
	Data                  interface{}
	IsStatementExpression bool
}

type CallArgs []*Expression // f(1,2)　调用参数列表

type ExpressionFunctionCall struct {
	Expression *Expression
	Args       CallArgs
	Func       *Function
}

type ExpressionDeclareVariable struct {
	Vs          []*VariableDefinition
	Expressions []*Expression
}

type ExpressionDeclareConsts struct {
	Cs []*Const
}

type ExpressionTypeConvertion struct {
	Typ        *VariableType
	Expression *Expression
}

type ExpressionIdentifer struct {
	Name string
	Var  *VariableDefinition
	Func *Function
	//enumas
	Enum     *Enum
	EnumName *EnumName
	//class
	Class *Class
}

type ExpressionIndex struct {
	Expression *Expression
	Index      *Expression
	Name       string
	Field      *ClassField
}

type ExpressionMethodCall struct {
	Expression *Expression
	Args       CallArgs
	Name       string
	Method     *ClassMethod
}

type ExpressionNew struct {
	Typ          *VariableType
	Args         CallArgs
	Construction *ClassMethod
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

func (e *Expression) OpName(typ ...int) string {
	t := e.Typ
	if len(typ) > 0 {
		t = typ[0]
	}
	switch t {
	case EXPRESSION_TYPE_BOOL:
		return fmt.Sprintf("bool(%v)", e.Data.(bool))
	case EXPRESSION_TYPE_BYTE:
		return fmt.Sprintf("byte(%v)", e.Data.(byte))
	case EXPRESSION_TYPE_INT:
		return fmt.Sprintf("int(%v)", e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return fmt.Sprintf("float(%v)", e.Data.(float64))
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
		return fmt.Sprintf("identifier(%s)", e.Data.(*ExpressionIdentifer).Name)
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

/*
	take one argument
*/
func (e *Expression) isNumber(typ ...int) bool {
	t := e.Typ
	if len(typ) > 0 {
		t = typ[0]
	}
	return t == EXPRESSION_TYPE_BYTE || t == EXPRESSION_TYPE_INT || t == EXPRESSION_TYPE_FLOAT
}
