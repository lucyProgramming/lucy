package ast

import (
	"fmt"
)

const (
	_ = iota // start with 1
	//null
	ExpressionTypeNull
	// bool
	ExpressionTypeBool
	// int types
	ExpressionTypeByte
	ExpressionTypeShort
	ExpressionTypeInt
	ExpressionTypeLong
	ExpressionTypeFloat
	ExpressionTypeDouble

	ExpressionTypeString
	ExpressionTypeArray // []bool{false,true}
	//binary expression
	ExpressionTypeLogicalOr
	ExpressionTypeLogicalAnd
	//
	ExpressionTypeOr
	ExpressionTypeAnd
	ExpressionTypeXor
	ExpressionTypeLsh
	ExpressionTypeRsh
	ExpressionTypeAdd
	ExpressionTypeSub
	ExpressionTypeMul
	ExpressionTypeDiv
	ExpressionTypeMod
	//
	ExpressionTypeAssign
	ExpressionTypeColonAssign
	//
	ExpressionTypePlusAssign
	ExpressionTypeMinusAssign
	ExpressionTypeMulAssign
	ExpressionTypeDivAssign
	ExpressionTypeModAssign
	ExpressionTypeAndAssign
	ExpressionTypeOrAssign
	ExpressionTypeXorAssign
	ExpressionTypeLshAssign
	ExpressionTypeRshAssign
	//
	ExpressionTypeEq
	ExpressionTypeNe
	ExpressionTypeGe
	ExpressionTypeGt
	ExpressionTypeLe
	ExpressionTypeLt
	//

	//
	ExpressionTypeIndex     // a["b"]
	ExpressionTypeSelection //a.b
	//
	ExpressionTypeMethodCall
	ExpressionTypeFunctionCall
	//
	ExpressionTypeIncrement
	ExpressionTypeDecrement
	ExpressionTypePrefixIncrement
	ExpressionTypePrefixDecrement
	//
	ExpressionTypeNegative
	ExpressionTypeNot
	ExpressionTypeBitwiseNot
	//
	ExpressionTypeIdentifier
	ExpressionTypeNew
	ExpressionTypeList
	ExpressionTypeFunctionLiteral
	ExpressionTypeVar
	ExpressionTypeConst
	ExpressionTypeCheckCast // []byte(str)

	ExpressionTypeRange // for range
	ExpressionTypeSlice // arr[0:2]
	ExpressionTypeMap   // map literal
	ExpressionTypeTypeAlias
	ExpressionTypeTypeAssert
	ExpressionTypeTernary
	ExpressionTypeGlobal
)

func (e *Expression) OpName() string {
	switch e.Type {
	case ExpressionTypeBool:
		return fmt.Sprintf("%v", e.Data.(bool))
	case ExpressionTypeByte:
		return fmt.Sprintf("%v", e.Data.(byte))
	case ExpressionTypeShort:
		return fmt.Sprintf("%vs", e.Data.(int32))
	case ExpressionTypeInt:
		return fmt.Sprintf("%v", e.Data.(int32))
	case ExpressionTypeLong:
		return fmt.Sprintf("%vL", e.Data.(int64))
	case ExpressionTypeFloat:
		return fmt.Sprintf("%vf", e.Data.(float32))
	case ExpressionTypeDouble:
		return fmt.Sprintf("%vd", e.Data.(float64))
	case ExpressionTypeString:
		return fmt.Sprintf("\"%v\"", e.Data)
	case ExpressionTypeArray:
		return "array_literal"
	case ExpressionTypeLogicalOr:
		return "||"
	case ExpressionTypeLogicalAnd:
		return "&&"
	case ExpressionTypeOr:
		return "|"
	case ExpressionTypeAnd:
		return "&"
	case ExpressionTypeXor:
		return "^"
	case ExpressionTypeLsh:
		return "<<"
	case ExpressionTypeRsh:
		return ">>"
	case ExpressionTypeAssign:
		return "="
	case ExpressionTypeColonAssign:
		return ":="
	case ExpressionTypePlusAssign:
		return "+="
	case ExpressionTypeMinusAssign:
		return "-="
	case ExpressionTypeMulAssign:
		return "*="
	case ExpressionTypeDivAssign:
		return "/="
	case ExpressionTypeModAssign:
		return "%="
	case ExpressionTypeAndAssign:
		return "&="
	case ExpressionTypeOrAssign:
		return "|="
	case ExpressionTypeLshAssign:
		return "<<="
	case ExpressionTypeRshAssign:
		return ">>="
	case ExpressionTypeXorAssign:
		return "^="
	case ExpressionTypeEq:
		return "=="
	case ExpressionTypeNe:
		return "!="
	case ExpressionTypeGe:
		return ">="
	case ExpressionTypeGt:
		return ">"
	case ExpressionTypeLe:
		return "<="
	case ExpressionTypeLt:
		return "<"
	case ExpressionTypeAdd:
		return "+"
	case ExpressionTypeSub:
		return "-"
	case ExpressionTypeMul:
		return "*"
	case ExpressionTypeDiv:
		return "/"
	case ExpressionTypeMod:
		return "%"
	case ExpressionTypeIndex: // a["b"]
		t := e.Data.(*ExpressionIndex)
		return fmt.Sprintf("%s[%s]", t.Expression.OpName(), t.Index.OpName())
	case ExpressionTypeSelection: //a.b
		t := e.Data.(*ExpressionSelection)
		return fmt.Sprintf("%s.%s", t.Expression.OpName(), t.Name)
	case ExpressionTypeMethodCall:
		t := e.Data.(*ExpressionMethodCall)
		return fmt.Sprintf("%s.%s()", t.Expression.OpName(), t.Name)
	case ExpressionTypeFunctionCall:
		t := e.Data.(*ExpressionFunctionCall)
		return fmt.Sprintf("function_call(%s)", t.Expression.OpName())
	case ExpressionTypeIncrement:
		return "++"
	case ExpressionTypeDecrement:
		return "--"
	case ExpressionTypePrefixIncrement:
		return "++"
	case ExpressionTypePrefixDecrement:
		return "--"
	case ExpressionTypeNegative:
		return "negative(-)"
	case ExpressionTypeTernary:
		return "ternary(?:)"
	case ExpressionTypeNot:
		return "not(!)"
	case ExpressionTypeBitwiseNot:
		return "~"
	case ExpressionTypeIdentifier:
		return fmt.Sprintf("identifier_%s", e.Data.(*ExpressionIdentifier).Name)
	case ExpressionTypeNull:
		return "null"
	case ExpressionTypeNew:
		return "new"
	case ExpressionTypeList:
		return "expression_list"
	case ExpressionTypeFunctionLiteral:
		return "function_literal"
	case ExpressionTypeConst:
		return "const"
	case ExpressionTypeVar:
		return "var"
	case ExpressionTypeRange:
		return "range"
	case ExpressionTypeSlice:
		return "slice"
	case ExpressionTypeMap:
		return "map_literal"
	case ExpressionTypeCheckCast:
		return "conversion of type"
	case ExpressionTypeTypeAssert:
		return "type assert"
	case ExpressionTypeTypeAlias:
		return "type alias"
	case ExpressionTypeGlobal:
		return "global"
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
	e.Type = ExpressionTypeCheckCast
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
	case VariableTypeBool:
		e.Type = ExpressionTypeBool
		e.Data = c.Value.(bool)
	case VariableTypeByte:
		e.Type = ExpressionTypeByte
		e.Data = c.Value.(byte)
	case VariableTypeShort:
		e.Type = ExpressionTypeShort
		e.Data = c.Value.(int32)
	case VariableTypeInt:
		e.Type = ExpressionTypeInt
		e.Data = c.Value.(int32)
	case VariableTypeLong:
		e.Type = ExpressionTypeLong
		e.Data = c.Value.(int64)
	case VariableTypeFloat:
		e.Type = ExpressionTypeFloat
		e.Data = c.Value.(float32)
	case VariableTypeDouble:
		e.Type = ExpressionTypeDouble
		e.Data = c.Value.(float64)
	case VariableTypeString:
		e.Type = ExpressionTypeString
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
	return e.Type == ExpressionTypeBool ||
		e.Type == ExpressionTypeString ||
		e.isNumber()
}

/*
	valid for condition
*/
func (e *Expression) canBeUsedAsCondition() bool {
	return e.Type == ExpressionTypeNull ||
		e.Type == ExpressionTypeBool ||
		e.Type == ExpressionTypeByte ||
		e.Type == ExpressionTypeShort ||
		e.Type == ExpressionTypeInt ||
		e.Type == ExpressionTypeLong ||
		e.Type == ExpressionTypeFloat ||
		e.Type == ExpressionTypeDouble ||
		e.Type == ExpressionTypeString ||
		e.Type == ExpressionTypeArray ||
		e.Type == ExpressionTypeLogicalOr ||
		e.Type == ExpressionTypeLogicalAnd ||
		e.Type == ExpressionTypeOr ||
		e.Type == ExpressionTypeAnd ||
		e.Type == ExpressionTypeXor ||
		e.Type == ExpressionTypeLsh ||
		e.Type == ExpressionTypeRsh ||
		e.Type == ExpressionTypeAdd ||
		e.Type == ExpressionTypeSub ||
		e.Type == ExpressionTypeMul ||
		e.Type == ExpressionTypeDiv ||
		e.Type == ExpressionTypeMod ||
		e.Type == ExpressionTypeEq ||
		e.Type == ExpressionTypeNe ||
		e.Type == ExpressionTypeGe ||
		e.Type == ExpressionTypeGt ||
		e.Type == ExpressionTypeLe ||
		e.Type == ExpressionTypeLt ||
		e.Type == ExpressionTypeIndex ||
		e.Type == ExpressionTypeSelection ||
		e.Type == ExpressionTypeMethodCall ||
		e.Type == ExpressionTypeFunctionCall ||
		e.Type == ExpressionTypeIncrement ||
		e.Type == ExpressionTypeDecrement ||
		e.Type == ExpressionTypePrefixIncrement ||
		e.Type == ExpressionTypePrefixDecrement ||
		e.Type == ExpressionTypeNegative ||
		e.Type == ExpressionTypeNot ||
		e.Type == ExpressionTypeBitwiseNot ||
		e.Type == ExpressionTypeIdentifier ||
		e.Type == ExpressionTypeNew ||
		e.Type == ExpressionTypeCheckCast ||
		e.Type == ExpressionTypeSlice ||
		e.Type == ExpressionTypeMap ||
		e.Type == ExpressionTypeTernary
}

func (e *Expression) canBeUsedAsStatement() bool {
	return e.Type == ExpressionTypeColonAssign ||
		e.Type == ExpressionTypeAssign ||
		e.Type == ExpressionTypeFunctionCall ||
		e.Type == ExpressionTypeMethodCall ||
		e.Type == ExpressionTypeFunctionLiteral ||
		e.Type == ExpressionTypePlusAssign ||
		e.Type == ExpressionTypeMinusAssign ||
		e.Type == ExpressionTypeMulAssign ||
		e.Type == ExpressionTypeDivAssign ||
		e.Type == ExpressionTypeModAssign ||
		e.Type == ExpressionTypeAndAssign ||
		e.Type == ExpressionTypeOrAssign ||
		e.Type == ExpressionTypeXorAssign ||
		e.Type == ExpressionTypeLshAssign ||
		e.Type == ExpressionTypeRshAssign ||
		e.Type == ExpressionTypeIncrement ||
		e.Type == ExpressionTypeDecrement ||
		e.Type == ExpressionTypePrefixIncrement ||
		e.Type == ExpressionTypePrefixDecrement ||
		e.Type == ExpressionTypeVar ||
		e.Type == ExpressionTypeConst
}

func (e *Expression) isNumber() bool {
	return e.isInteger() || e.isFloat()
}

func (e *Expression) isInteger() bool {
	return e.Type == ExpressionTypeByte ||
		e.Type == ExpressionTypeShort ||
		e.Type == ExpressionTypeInt ||
		e.Type == ExpressionTypeLong
}
func (e *Expression) isFloat() bool {
	return e.Type == ExpressionTypeFloat ||
		e.Type == ExpressionTypeDouble
}

/*
	check out this expression is increment or decrement
*/
func (e *Expression) IsIncrement() bool {
	return e.Type == ExpressionTypeIncrement ||
		e.Type == ExpressionTypePrefixIncrement
}

func (e *Expression) isListAndMoreThanNElements(n int) bool {
	if e.Type != ExpressionTypeList {
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
	if e.Type != ExpressionTypeAssign && e.Type != ExpressionTypeColonAssign {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
	if bin.Right.Type == ExpressionTypeRange {
		return true
	}
	if bin.Right.Type == ExpressionTypeList {
		t := bin.Right.Data.([]*Expression)
		if len(t) == 1 && t[0].Type == ExpressionTypeRange {
			// bin.Right = t[0] // override
			return true
		}
	}
	return false
}

func (e *Expression) MayHaveMultiValue() bool {
	return e.Type == ExpressionTypeFunctionCall ||
		e.Type == ExpressionTypeMethodCall ||
		e.Type == ExpressionTypeTypeAssert
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
	Class          *Class //
	Expression     *Expression
	Args           CallArgs
	Name           string
	Method         *ClassMethod
	ParameterTypes []*Type
	//PackageFunction *Function // Expression is package
}

type ExpressionDeclareVariable struct {
	Variables        []*Variable
	InitValues       []*Expression
	IfDeclaredBefore []bool // used for colon assign
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
}

type ExpressionIndex struct {
	Expression *Expression
	Index      *Expression
}
type ExpressionSelection struct {
	Expression      *Expression
	Name            string
	Field           *ClassField  // expression is class or object
	Method          *ClassMethod // pack to method handle
	PackageVariable *Variable    // expression is package
	PackageEnumName *EnumName    // expression is package
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
	if e.Type != ExpressionTypeIdentifier {
		return false
	}
	return e.Data.(*ExpressionIdentifier).Name == THIS
}

func (e *Expression) IsNoNameIdentifier() bool {
	if e.Type != ExpressionTypeIdentifier {
		return false
	}
	return e.Data.(*ExpressionIdentifier).Name == NoNameIdentifier
}
