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
	EXPRESSION_TYPE_INDEX     // a["b"]
	EXPRESSION_TYPE_SELECTION //a.b
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
	EXPRESSION_TYPE_BIT_NOT
	//
	EXPRESSION_TYPE_IDENTIFIER
	EXPRESSION_TYPE_NEW
	EXPRESSION_TYPE_LIST
	EXPRESSION_TYPE_FUNCTION_LITERAL
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
	switch e.Type {
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
	case EXPRESSION_TYPE_SELECTION: //a.b
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
	case EXPRESSION_TYPE_BIT_NOT:
		return "~"
	case EXPRESSION_TYPE_IDENTIFIER:
		return fmt.Sprintf("identifier_%s", e.Data.(*ExpressionIdentifier).Name)
	case EXPRESSION_TYPE_NULL:
		return "null"
	case EXPRESSION_TYPE_NEW:
		return "new"
	case EXPRESSION_TYPE_LIST:
		return "expression_list"
	case EXPRESSION_TYPE_FUNCTION_LITERAL:
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
		return "conversion of type"
	case EXPRESSION_TYPE_TYPE_ASSERT:
		return "type assert"
	case EXPRESSION_TYPE_TYPE_ALIAS:
		return "type alias"
	default:
		return fmt.Sprintf("op[%d](missing handle)", e.Type)
	}
}

type Expression struct {
	Type                  int
	IsPublic              bool // only for global variable definition
	IsCompileAuto         bool // compile auto expression
	ExpressionValue       *Type
	ExpressionMultiValues []*Type
	Pos                   *Position
	Data                  interface{}
	IsStatementExpression bool
}

func (e *Expression) ConvertTo(t *Type) {
	c := &ExpressionTypeConversion{}
	c.Expression = &Expression{}
	*c.Expression = *e // copy
	c.Type = t
	e.ExpressionValue = t
	e.Type = EXPRESSION_TYPE_CHECK_CAST
	e.IsCompileAuto = true
	e.Data = c
}

func (e *Expression) ConvertToNumber(typ int) {
	if e.IsLiteral() {
		e.convertNumberLiteralTo(typ)
		e.ExpressionValue = &Type{
			Type: typ,
			Pos:  e.Pos,
		}
	} else {
		e.ConvertTo(&Type{
			Pos:  e.Pos,
			Type: typ,
		})
	}
}

type ExpressionTypeAssert ExpressionTypeConversion

/*
	const
*/
func (e *Expression) fromConst(c *Constant) {
	switch c.Type.Type {
	case VARIABLE_TYPE_BOOL:
		e.Type = EXPRESSION_TYPE_BOOL
		e.Data = c.Value.(bool)
	case VARIABLE_TYPE_BYTE:
		e.Type = EXPRESSION_TYPE_BYTE
		e.Data = c.Value.(byte)
	case VARIABLE_TYPE_SHORT:
		e.Type = EXPRESSION_TYPE_SHORT
		e.Data = c.Value.(int32)
	case VARIABLE_TYPE_INT:
		e.Type = EXPRESSION_TYPE_INT
		e.Data = c.Value.(int32)
	case VARIABLE_TYPE_LONG:
		e.Type = EXPRESSION_TYPE_LONG
		e.Data = c.Value.(int64)
	case VARIABLE_TYPE_FLOAT:
		e.Type = EXPRESSION_TYPE_FLOAT
		e.Data = c.Value.(float32)
	case VARIABLE_TYPE_DOUBLE:
		e.Type = EXPRESSION_TYPE_DOUBLE
		e.Data = c.Value.(float64)
	case VARIABLE_TYPE_STRING:
		e.Type = EXPRESSION_TYPE_STRING
		e.Data = c.Value.(string)
	}
}

type ExpressionTypeAlias struct {
	Name string
	Type *Type
	Pos  *Position
}

type ExpressionTernary struct {
	Selection *Expression
	True      *Expression
	False     *Expression
}

type ExpressionSlice struct {
	Array      *Expression
	Start, End *Expression
}

func (e *Expression) IsLiteral() bool {
	return e.Type == EXPRESSION_TYPE_BOOL ||
		e.Type == EXPRESSION_TYPE_STRING ||
		e.isNumber()
}

/*
	valid for condition
*/
func (e *Expression) canBeUsedAsCondition() bool {
	return e.Type == EXPRESSION_TYPE_NULL ||
		e.Type == EXPRESSION_TYPE_BOOL ||
		e.Type == EXPRESSION_TYPE_BYTE ||
		e.Type == EXPRESSION_TYPE_SHORT ||
		e.Type == EXPRESSION_TYPE_INT ||
		e.Type == EXPRESSION_TYPE_LONG ||
		e.Type == EXPRESSION_TYPE_FLOAT ||
		e.Type == EXPRESSION_TYPE_DOUBLE ||
		e.Type == EXPRESSION_TYPE_STRING ||
		e.Type == EXPRESSION_TYPE_ARRAY ||
		e.Type == EXPRESSION_TYPE_LOGICAL_OR ||
		e.Type == EXPRESSION_TYPE_LOGICAL_AND ||
		e.Type == EXPRESSION_TYPE_OR ||
		e.Type == EXPRESSION_TYPE_AND ||
		e.Type == EXPRESSION_TYPE_XOR ||
		e.Type == EXPRESSION_TYPE_LSH ||
		e.Type == EXPRESSION_TYPE_RSH ||
		e.Type == EXPRESSION_TYPE_ADD ||
		e.Type == EXPRESSION_TYPE_SUB ||
		e.Type == EXPRESSION_TYPE_MUL ||
		e.Type == EXPRESSION_TYPE_DIV ||
		e.Type == EXPRESSION_TYPE_MOD ||
		e.Type == EXPRESSION_TYPE_EQ ||
		e.Type == EXPRESSION_TYPE_NE ||
		e.Type == EXPRESSION_TYPE_GE ||
		e.Type == EXPRESSION_TYPE_GT ||
		e.Type == EXPRESSION_TYPE_LE ||
		e.Type == EXPRESSION_TYPE_LT ||
		e.Type == EXPRESSION_TYPE_INDEX ||
		e.Type == EXPRESSION_TYPE_SELECTION ||
		e.Type == EXPRESSION_TYPE_METHOD_CALL ||
		e.Type == EXPRESSION_TYPE_FUNCTION_CALL ||
		e.Type == EXPRESSION_TYPE_INCREMENT ||
		e.Type == EXPRESSION_TYPE_DECREMENT ||
		e.Type == EXPRESSION_TYPE_PRE_INCREMENT ||
		e.Type == EXPRESSION_TYPE_PRE_DECREMENT ||
		e.Type == EXPRESSION_TYPE_NEGATIVE ||
		e.Type == EXPRESSION_TYPE_NOT ||
		e.Type == EXPRESSION_TYPE_BIT_NOT ||
		e.Type == EXPRESSION_TYPE_IDENTIFIER ||
		e.Type == EXPRESSION_TYPE_NEW ||
		e.Type == EXPRESSION_TYPE_CHECK_CAST ||
		e.Type == EXPRESSION_TYPE_SLICE ||
		e.Type == EXPRESSION_TYPE_MAP ||
		e.Type == EXPRESSION_TYPE_TERNARY
}

func (e *Expression) canBeUsedAsStatement() bool {
	return e.Type == EXPRESSION_TYPE_COLON_ASSIGN ||
		e.Type == EXPRESSION_TYPE_ASSIGN ||
		e.Type == EXPRESSION_TYPE_FUNCTION_CALL ||
		e.Type == EXPRESSION_TYPE_METHOD_CALL ||
		e.Type == EXPRESSION_TYPE_FUNCTION_LITERAL ||
		e.Type == EXPRESSION_TYPE_PLUS_ASSIGN ||
		e.Type == EXPRESSION_TYPE_MINUS_ASSIGN ||
		e.Type == EXPRESSION_TYPE_MUL_ASSIGN ||
		e.Type == EXPRESSION_TYPE_DIV_ASSIGN ||
		e.Type == EXPRESSION_TYPE_MOD_ASSIGN ||
		e.Type == EXPRESSION_TYPE_AND_ASSIGN ||
		e.Type == EXPRESSION_TYPE_OR_ASSIGN ||
		e.Type == EXPRESSION_TYPE_XOR_ASSIGN ||
		e.Type == EXPRESSION_TYPE_LSH_ASSIGN ||
		e.Type == EXPRESSION_TYPE_RSH_ASSIGN ||
		e.Type == EXPRESSION_TYPE_INCREMENT ||
		e.Type == EXPRESSION_TYPE_DECREMENT ||
		e.Type == EXPRESSION_TYPE_PRE_INCREMENT ||
		e.Type == EXPRESSION_TYPE_PRE_DECREMENT ||
		e.Type == EXPRESSION_TYPE_VAR ||
		e.Type == EXPRESSION_TYPE_CONST
}

func (e *Expression) isNumber() bool {
	return e.isInteger() || e.isFloat()
}

func (e *Expression) isInteger() bool {
	return e.Type == EXPRESSION_TYPE_BYTE ||
		e.Type == EXPRESSION_TYPE_SHORT ||
		e.Type == EXPRESSION_TYPE_INT ||
		e.Type == EXPRESSION_TYPE_LONG
}
func (e *Expression) isFloat() bool {
	return e.Type == EXPRESSION_TYPE_FLOAT ||
		e.Type == EXPRESSION_TYPE_DOUBLE
}

/*
	check out this expression is increment or decrement
*/
func (e *Expression) IsIncrement() bool {
	return e.Type == EXPRESSION_TYPE_INCREMENT ||
		e.Type == EXPRESSION_TYPE_PRE_INCREMENT
}

func (e *Expression) isListAndMoreThanNElements(n int) bool {
	if e.Type != EXPRESSION_TYPE_LIST {
		return false
	}
	return len(e.Data.([]*Expression)) > n
}

func (e *Expression) HaveOnlyOneValue() bool {
	if e.MayHaveMultiValue() {
		return len(e.ExpressionMultiValues) == 1
	}
	return true
}

/*
	k,v := range arr
	k,v = range arr
*/
func (e *Expression) canBeUsedForRange() bool {
	if e.Type != EXPRESSION_TYPE_ASSIGN && e.Type != EXPRESSION_TYPE_COLON_ASSIGN {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
	if bin.Right.Type == EXPRESSION_TYPE_RANGE {
		return true
	}
	if bin.Right.Type == EXPRESSION_TYPE_LIST {
		t := bin.Right.Data.([]*Expression)
		if len(t) == 1 && t[0].Type == EXPRESSION_TYPE_RANGE {
			// bin.Right = t[0] // override
			return true
		}
	}
	return false
}

func (e *Expression) MayHaveMultiValue() bool {
	return e.Type == EXPRESSION_TYPE_FUNCTION_CALL ||
		e.Type == EXPRESSION_TYPE_METHOD_CALL ||
		e.Type == EXPRESSION_TYPE_TYPE_ASSERT
}

func (e *Expression) CallHasReturnValue() bool {
	return len(e.ExpressionMultiValues) >= 1 && e.ExpressionMultiValues[0].RightValueValid()
}

type CallArgs []*Expression // f(1,2)

type ExpressionFunctionCall struct {
	BuildInFunctionMeta      interface{} // for build function only
	Expression               *Expression
	Args                     CallArgs
	Function                 *Function
	ParameterTypes           []*Type // for template function
	TemplateFunctionCallPair *TemplateFunctionCallPair
	FunctionPointer          *FunctionType
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
	ParameterTypes  []*Type
}

type ExpressionDeclareVariable struct {
	Variables        []*Variable
	InitValues       []*Expression
	IfDeclaredBefore []bool // used for colon assign
}

func (e *ExpressionDeclareVariable) haveFunctionPointer() {
	//for _, v := range e.Variables {
	//	if v.Name == NO_NAME_IDENTIFIER {
	//		continue
	//	}
	//	if v.Type == nil || v.Type.Type != VARIABLE_TYPE_FUNCTION {
	//		continue
	//	}
	//	v.Type.FunctionType = &v.Type.Function.Type
	//}
}

type ExpressionTypeConversion struct {
	Type       *Type
	Expression *Expression
}

type ExpressionIdentifier struct {
	Name     string
	Variable *Variable
	Function *Function
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
	Field           *ClassField // expression is class or object
	PackageVariable *Variable   // expression is package
	PackageEnumName *EnumName   // expression is package
}

type ExpressionNew struct {
	Type                     *Type
	Args                     CallArgs
	Construction             *ClassMethod
	IsConvertJavaArray2Array bool
}

type ExpressionMap struct {
	Type          *Type
	KeyValuePairs []*ExpressionBinary
}

/*
for some general purpose
*/
type ExpressionBinary struct {
	Left  *Expression
	Right *Expression
}

type ExpressionArray struct {
	Type        *Type
	Expressions []*Expression
	Length      int
}

func (e *Expression) isThis() bool {
	if e.Type != EXPRESSION_TYPE_IDENTIFIER {
		return false
	}
	return e.Data.(*ExpressionIdentifier).Name == THIS
}

func (e *Expression) IsNoNameIdentifier() bool {
	if e.Type != EXPRESSION_TYPE_IDENTIFIER {
		return false
	}
	return e.Data.(*ExpressionIdentifier).Name == NO_NAME_IDENTIFIER
}
