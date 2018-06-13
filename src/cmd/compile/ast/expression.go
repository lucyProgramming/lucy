package ast

import (
	"fmt"
)

const (
	_ = iota // start with 1
	//null
	EXPRESSION_TYPE_NULL
	// bool
	EXPRESSION_TYPE_BOOL
	// int types
	EXPRESSION_TYPE_BYTE
	EXPRESSION_TYPE_SHORT
	EXPRESSION_TYPE_INT
	EXPRESSION_TYPE_LONG
	EXPRESSION_TYPE_FLOAT
	EXPRESSION_TYPE_DOUBLE

	EXPRESSION_TYPE_STRING
	EXPRESSION_TYPE_ARRAY // []bool{false,true}
	//binary expression
	EXPRESSION_TYPE_LOGICAL_OR
	EXPRESSION_TYPE_LOGICAL_AND
	//
	EXPRESSION_TYPE_OR
	EXPRESSION_TYPE_AND
	EXPRESSION_TYPE_XOR
	EXPRESSION_TYPE_LSH
	EXPRESSION_TYPE_RSH
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
	EXPRESSION_TYPE_AND_ASSIGN
	EXPRESSION_TYPE_OR_ASSIGN
	EXPRESSION_TYPE_XOR_ASSIGN
	EXPRESSION_TYPE_LSH_ASSIGN
	EXPRESSION_TYPE_RSH_ASSIGN
	//
	EXPRESSION_TYPE_EQ
	EXPRESSION_TYPE_NE
	EXPRESSION_TYPE_GE
	EXPRESSION_TYPE_GT
	EXPRESSION_TYPE_LE
	EXPRESSION_TYPE_LT
	//

	//
	EXPRESSION_TYPE_INDEX  // a["b"]
	EXPRESSION_TYPE_SELECT //a.b
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
	EXPRESSION_TYPE_BITWISE_NOT
	//
	EXPRESSION_TYPE_IDENTIFIER
	EXPRESSION_TYPE_NEW
	EXPRESSION_TYPE_LIST
	EXPRESSION_TYPE_FUNCTION
	EXPRESSION_TYPE_VAR
	EXPRESSION_TYPE_CONST
	EXPRESSION_TYPE_CHECK_CAST // []byte(str)

	EXPRESSION_TYPE_RANGE // for range
	EXPRESSION_TYPE_SLICE // arr[0:2]
	EXPRESSION_TYPE_MAP   // map literal
	EXPRESSION_TYPE_TYPE_ALIAS
	EXPRESSION_TYPE_TYPE_ASSERT
	EXPRESSION_TYPE_TERNARY
)

func (e *Expression) OpName() string {
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		return fmt.Sprintf("%v", e.Data.(bool))
	case EXPRESSION_TYPE_BYTE:
		return fmt.Sprintf("%v", e.Data.(byte))
	case EXPRESSION_TYPE_SHORT:
		return fmt.Sprintf("%vs", e.Data.(int32))
	case EXPRESSION_TYPE_INT:
		return fmt.Sprintf("%v", e.Data.(int32))
	case EXPRESSION_TYPE_LONG:
		return fmt.Sprintf("%vL", e.Data.(int64))
	case EXPRESSION_TYPE_FLOAT:
		return fmt.Sprintf("%vf", e.Data.(float32))
	case EXPRESSION_TYPE_DOUBLE:
		return fmt.Sprintf("%vd", e.Data.(float64))
	case EXPRESSION_TYPE_STRING:
		return fmt.Sprintf("\"%v\"", e.Data)
	case EXPRESSION_TYPE_ARRAY:
		return "array_literal"
	case EXPRESSION_TYPE_LOGICAL_OR:
		return "||"
	case EXPRESSION_TYPE_LOGICAL_AND:
		return "&&"
	case EXPRESSION_TYPE_OR:
		return "|"
	case EXPRESSION_TYPE_AND:
		return "&"
	case EXPRESSION_TYPE_XOR:
		return "^"
	case EXPRESSION_TYPE_LSH:
		return "<<"
	case EXPRESSION_TYPE_RSH:
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
	case EXPRESSION_TYPE_AND_ASSIGN:
		return "&="
	case EXPRESSION_TYPE_OR_ASSIGN:
		return "|="
	case EXPRESSION_TYPE_LSH_ASSIGN:
		return "<<="
	case EXPRESSION_TYPE_RSH_ASSIGN:
		return ">>="
	case EXPRESSION_TYPE_XOR_ASSIGN:
		return "^="
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
		t := e.Data.(*ExpressionIndex)
		return fmt.Sprintf("%s[%s]", t.Expression.OpName(), t.Index.OpName())
	case EXPRESSION_TYPE_SELECT: //a.b
		t := e.Data.(*ExpressionSelection)
		return fmt.Sprintf("%s.%s", t.Expression.OpName(), t.Name)
	case EXPRESSION_TYPE_METHOD_CALL:
		t := e.Data.(*ExpressionMethodCall)
		return fmt.Sprintf("%s.%s()", t.Expression.OpName(), t.Name)
	case EXPRESSION_TYPE_FUNCTION_CALL:
		t := e.Data.(*ExpressionFunctionCall)
		return fmt.Sprintf("function_call(%s)", t.Expression.OpName())
	case EXPRESSION_TYPE_INCREMENT:
		return "++"
	case EXPRESSION_TYPE_DECREMENT:
		return "--"
	case EXPRESSION_TYPE_PRE_INCREMENT:
		return "++"
	case EXPRESSION_TYPE_PRE_DECREMENT:
		return "--"
	case EXPRESSION_TYPE_NEGATIVE:
		return "negative(-)"
	case EXPRESSION_TYPE_TERNARY:
		return "ternary(?:)"
	case EXPRESSION_TYPE_NOT:
		return "not(!)"
	case EXPRESSION_TYPE_BITWISE_NOT:
		return "~"
	case EXPRESSION_TYPE_IDENTIFIER:
		return fmt.Sprintf("identifier_%s", e.Data.(*ExpressionIdentifier).Name)
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
	case EXPRESSION_TYPE_RANGE:
		return "range"
	case EXPRESSION_TYPE_SLICE:
		return "slice"
	case EXPRESSION_TYPE_MAP:
		return "map_literal"
	case EXPRESSION_TYPE_CHECK_CAST:
		return "convertion of type"
	case EXPRESSION_TYPE_TYPE_ASSERT:
		return "type assert"
	case EXPRESSION_TYPE_TYPE_ALIAS:
		return "type alias"
	default:
		return fmt.Sprintf("op[%d](missing handle)", e.Typ)
	}
}

type Expression struct {
	Typ                   int
	IsPublic              bool // only for global var definition
	IsCompileAuto         bool // compile auto expression
	Value                 *VariableType
	Values                []*VariableType
	Pos                   *Pos
	Data                  interface{}
	IsStatementExpression bool
}

func (e *Expression) ConvertTo(t *VariableType) {
	c := &ExpressionTypeConversion{}
	c.Expression = &Expression{}
	*c.Expression = *e // copy
	c.Typ = t
	e.Value = t
	e.Typ = EXPRESSION_TYPE_CHECK_CAST
	e.IsCompileAuto = true
	e.Data = c
}

func (e *Expression) ConvertToNumber(typ int) {
	if e.IsLiteral() {
		e.convertNumberLiteralTo(typ)
		e.Value = &VariableType{
			Typ: typ,
			Pos: e.Pos,
		}
	} else {
		e.ConvertTo(&VariableType{
			Pos: e.Pos,
			Typ: typ,
		})
	}
}

type ExpressionTypeAssert ExpressionTypeConversion

/*
	const
*/
func (e *Expression) fromConst(c *Const) {
	switch c.Typ.Typ {
	case VARIABLE_TYPE_BOOL:
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data = c.Value.(bool)
	case VARIABLE_TYPE_BYTE:
		e.Typ = EXPRESSION_TYPE_BYTE
		e.Data = c.Value.(byte)
	case VARIABLE_TYPE_SHORT:
		e.Typ = EXPRESSION_TYPE_SHORT
		e.Data = c.Value.(int32)
	case VARIABLE_TYPE_INT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = c.Value.(int32)
	case VARIABLE_TYPE_LONG:
		e.Typ = EXPRESSION_TYPE_LONG
		e.Data = c.Value.(int64)
	case VARIABLE_TYPE_FLOAT:
		e.Typ = EXPRESSION_TYPE_FLOAT
		e.Data = c.Value.(float32)
	case VARIABLE_TYPE_DOUBLE:
		e.Typ = EXPRESSION_TYPE_DOUBLE
		e.Data = c.Value.(float64)
	case VARIABLE_TYPE_STRING:
		e.Typ = EXPRESSION_TYPE_STRING
		e.Data = c.Value.(string)
	}
}

type ExpressionTypeAlias struct {
	Name string
	Typ  *VariableType
	Pos  *Pos
}

type ExpressionTernary struct {
	Condition *Expression
	True      *Expression
	False     *Expression
}

type ExpressionSlice struct {
	Array      *Expression
	Start, End *Expression
}

func (e *Expression) IsLiteral() bool {
	return e.Typ == EXPRESSION_TYPE_BOOL ||
		e.Typ == EXPRESSION_TYPE_STRING ||
		e.isNumber()
}

/*
	valid for condition
*/
func (e *Expression) canbeUsedAsCondition() bool {
	return e.Typ == EXPRESSION_TYPE_NULL ||
		e.Typ == EXPRESSION_TYPE_BOOL ||
		e.Typ == EXPRESSION_TYPE_BYTE ||
		e.Typ == EXPRESSION_TYPE_SHORT ||
		e.Typ == EXPRESSION_TYPE_INT ||
		e.Typ == EXPRESSION_TYPE_LONG ||
		e.Typ == EXPRESSION_TYPE_FLOAT ||
		e.Typ == EXPRESSION_TYPE_DOUBLE ||
		e.Typ == EXPRESSION_TYPE_STRING ||
		e.Typ == EXPRESSION_TYPE_ARRAY ||
		e.Typ == EXPRESSION_TYPE_LOGICAL_OR ||
		e.Typ == EXPRESSION_TYPE_LOGICAL_AND ||
		e.Typ == EXPRESSION_TYPE_OR ||
		e.Typ == EXPRESSION_TYPE_AND ||
		e.Typ == EXPRESSION_TYPE_XOR ||
		e.Typ == EXPRESSION_TYPE_LSH ||
		e.Typ == EXPRESSION_TYPE_RSH ||
		e.Typ == EXPRESSION_TYPE_ADD ||
		e.Typ == EXPRESSION_TYPE_SUB ||
		e.Typ == EXPRESSION_TYPE_MUL ||
		e.Typ == EXPRESSION_TYPE_DIV ||
		e.Typ == EXPRESSION_TYPE_MOD ||
		e.Typ == EXPRESSION_TYPE_EQ ||
		e.Typ == EXPRESSION_TYPE_NE ||
		e.Typ == EXPRESSION_TYPE_GE ||
		e.Typ == EXPRESSION_TYPE_GT ||
		e.Typ == EXPRESSION_TYPE_LE ||
		e.Typ == EXPRESSION_TYPE_LT ||
		e.Typ == EXPRESSION_TYPE_INDEX ||
		e.Typ == EXPRESSION_TYPE_SELECT ||
		e.Typ == EXPRESSION_TYPE_METHOD_CALL ||
		e.Typ == EXPRESSION_TYPE_FUNCTION_CALL ||
		e.Typ == EXPRESSION_TYPE_INCREMENT ||
		e.Typ == EXPRESSION_TYPE_DECREMENT ||
		e.Typ == EXPRESSION_TYPE_PRE_INCREMENT ||
		e.Typ == EXPRESSION_TYPE_PRE_DECREMENT ||
		e.Typ == EXPRESSION_TYPE_NEGATIVE ||
		e.Typ == EXPRESSION_TYPE_NOT ||
		e.Typ == EXPRESSION_TYPE_BITWISE_NOT ||
		e.Typ == EXPRESSION_TYPE_IDENTIFIER ||
		e.Typ == EXPRESSION_TYPE_NEW ||
		e.Typ == EXPRESSION_TYPE_CHECK_CAST ||
		e.Typ == EXPRESSION_TYPE_SLICE ||
		e.Typ == EXPRESSION_TYPE_MAP ||
		e.Typ == EXPRESSION_TYPE_TERNARY
}

func (e *Expression) canBeUsedAsStatement() bool {
	return e.Typ == EXPRESSION_TYPE_COLON_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_FUNCTION_CALL ||
		e.Typ == EXPRESSION_TYPE_METHOD_CALL ||
		e.Typ == EXPRESSION_TYPE_FUNCTION ||
		e.Typ == EXPRESSION_TYPE_PLUS_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_MINUS_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_MUL_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_DIV_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_MOD_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_AND_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_OR_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_XOR_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_LSH_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_RSH_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_INCREMENT ||
		e.Typ == EXPRESSION_TYPE_DECREMENT ||
		e.Typ == EXPRESSION_TYPE_PRE_INCREMENT ||
		e.Typ == EXPRESSION_TYPE_PRE_DECREMENT ||
		e.Typ == EXPRESSION_TYPE_VAR ||
		e.Typ == EXPRESSION_TYPE_CONST
}

func (e *Expression) isNumber() bool {
	return e.isInteger() || e.isFloat()
}

func (e *Expression) isInteger() bool {
	return e.Typ == EXPRESSION_TYPE_BYTE ||
		e.Typ == EXPRESSION_TYPE_SHORT ||
		e.Typ == EXPRESSION_TYPE_INT ||
		e.Typ == EXPRESSION_TYPE_LONG
}
func (e *Expression) isFloat() bool {
	return e.Typ == EXPRESSION_TYPE_FLOAT ||
		e.Typ == EXPRESSION_TYPE_DOUBLE
}

/*
	check out this expression is increment or decrement
*/
func (e *Expression) IsSelfIncrement() bool {
	return e.Typ == EXPRESSION_TYPE_INCREMENT ||
		e.Typ == EXPRESSION_TYPE_PRE_INCREMENT
}

func (e *Expression) isListAndMoreThanIElements(i int) bool {
	if e.Typ != EXPRESSION_TYPE_LIST {
		return false
	}
	return len(e.Data.([]*Expression)) > i
}

func (e *Expression) HaveOnlyOneValue() bool {
	if e.MayHaveMultiValue() {
		return len(e.Values) == 1
	}
	return true
}

/*
	k,v := range arr
	k,v = range arr
*/
func (e *Expression) canBeUsedForRange() bool {
	if e.Typ != EXPRESSION_TYPE_ASSIGN && e.Typ != EXPRESSION_TYPE_COLON_ASSIGN {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
	if bin.Right.Typ == EXPRESSION_TYPE_RANGE {
		return true
	}
	if bin.Right.Typ == EXPRESSION_TYPE_LIST {
		t := bin.Right.Data.([]*Expression)
		if len(t) == 1 && t[0].Typ == EXPRESSION_TYPE_RANGE {
			// bin.Right = t[0] // override
			return true
		}
	}
	return false
}

func (e *Expression) MayHaveMultiValue() bool {
	return e.Typ == EXPRESSION_TYPE_FUNCTION_CALL ||
		e.Typ == EXPRESSION_TYPE_METHOD_CALL ||
		e.Typ == EXPRESSION_TYPE_TYPE_ASSERT
}

func (e *Expression) CallHasReturnValue() bool {
	return len(e.Values) >= 1 && e.Values[0].RightValueValid()
}

type CallArgs []*Expression // f(1,2)　调用参数列表

type ExpressionFunctionCall struct {
	BuildInFunctionMeta      interface{} // for build function only
	Expression               *Expression
	Args                     CallArgs
	Func                     *Function
	TypedParameters          []*VariableType // for template function
	TemplateFunctionCallPair *TemplateFunctionCallPair
}

func (e *ExpressionFunctionCall) FromMethodCall(call *ExpressionMethodCall) *ExpressionFunctionCall {
	e.Args = call.Args
	return e
}

type ExpressionMethodCall struct {
	Class           *Class //
	Expression      *Expression
	Args            CallArgs
	Name            string
	Method          *ClassMethod
	PackageFunction *Function // Expression is package
	TypedParameters []*VariableType
}

type ExpressionDeclareVariable struct {
	Variables       []*VariableDefinition
	Values          []*Expression
	IfDeclareBefore []bool // used for colon assign
}

type ExpressionTypeConversion struct {
	Typ        *VariableType
	Expression *Expression
}

type ExpressionIdentifier struct {
	Name     string
	Var      *VariableDefinition
	Func     *Function
	EnumName *EnumName
	Class    *Class
}

type ExpressionIndex struct {
	Expression *Expression
	Index      *Expression
}
type ExpressionSelection struct {
	Expression      *Expression
	Name            string
	Field           *ClassField         // expression is class or object
	PackageVariable *VariableDefinition // expression is package
	EnumName        *EnumName           // expression is package
}

type ExpressionNew struct {
	Typ                      *VariableType
	Args                     CallArgs
	Construction             *ClassMethod
	IsConvertJavaArray2Array bool
}

type ExpressionMap struct {
	Typ           *VariableType
	KeyValuePairs []*ExpressionBinary
}

// for general purpose
type ExpressionBinary struct {
	Left  *Expression
	Right *Expression
}

type ExpressionArrayLiteral struct {
	Typ         *VariableType
	Expressions []*Expression
	Length      int
}

func (e *Expression) isThis() bool {
	if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
		return false
	}
	return e.Data.(*ExpressionIdentifier).Name == THIS
}

func (e *Expression) IsNoNameIdentifier() bool {
	if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
		return false
	}
	return e.Data.(*ExpressionIdentifier).Name == NO_NAME_IDENTIFIER
}
