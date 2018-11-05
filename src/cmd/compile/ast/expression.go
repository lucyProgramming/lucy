package ast

import (
	"fmt"
)

type ExpressionTypeKind int

const (
	_                             ExpressionTypeKind = iota // start with 1
	ExpressionTypeNull                                      // null
	ExpressionTypeBool                                      // true or false
	ExpressionTypeByte                                      // 'a' or 97b
	ExpressionTypeShort                                     // 100s
	ExpressionTypeChar                                      // '\u0000'
	ExpressionTypeInt                                       // 100
	ExpressionTypeLong                                      // 100L
	ExpressionTypeFloat                                     // 1.0
	ExpressionTypeDouble                                    // 1.0d
	ExpressionTypeString                                    // "hello world"
	ExpressionTypeArray                                     // []bool{false,true}
	ExpressionTypeLogicalOr                                 // a || b
	ExpressionTypeLogicalAnd                                // a && b
	ExpressionTypeOr                                        // a | b
	ExpressionTypeAnd                                       // a & b
	ExpressionTypeXor                                       // a ^b
	ExpressionTypeLsh                                       // a << b
	ExpressionTypeRsh                                       // a >> b
	ExpressionTypeAdd                                       // a + b
	ExpressionTypeSub                                       // a - b
	ExpressionTypeMul                                       // a * b
	ExpressionTypeDiv                                       // a / b
	ExpressionTypeMod                                       // a % b
	ExpressionTypeAssign                                    // a = b
	ExpressionTypeVarAssign                                 // a := b
	ExpressionTypePlusAssign                                // a += b
	ExpressionTypeMinusAssign                               // a -= b
	ExpressionTypeMulAssign                                 // a *= b
	ExpressionTypeDivAssign                                 // a /= b
	ExpressionTypeModAssign                                 // a %= b
	ExpressionTypeAndAssign                                 // a &= b
	ExpressionTypeOrAssign                                  // a |= b
	ExpressionTypeXorAssign                                 // a ^= b
	ExpressionTypeLshAssign                                 // a <<= b
	ExpressionTypeRshAssign                                 // a >>= b
	ExpressionTypeEq                                        // a == b
	ExpressionTypeNe                                        // a != b
	ExpressionTypeGe                                        // a >= b
	ExpressionTypeGt                                        // a > b
	ExpressionTypeLe                                        // a <= b
	ExpressionTypeLt                                        // a < b
	ExpressionTypeIndex                                     // a["b"]
	ExpressionTypeSelection                                 // a.b
	ExpressionTypeSelectionConst                            // ::
	ExpressionTypeMethodCall                                // a.b()
	ExpressionTypeFunctionCall                              // a()
	ExpressionTypeIncrement                                 // a++
	ExpressionTypeDecrement                                 // a--
	ExpressionTypePrefixIncrement                           // ++ a
	ExpressionTypePrefixDecrement                           // -- a
	ExpressionTypeNegative                                  // -a
	ExpressionTypeNot                                       // !a
	ExpressionTypeBitwiseNot                                // ~a
	ExpressionTypeIdentifier                                // a
	ExpressionTypeNew                                       // new []int(10)
	ExpressionTypeList                                      // a,b := "hello","world"
	ExpressionTypeFunctionLiteral                           // fn() { print("hello world"); }
	ExpressionTypeVar                                       // var a,b int
	ExpressionTypeConst                                     // const a = "hello world"
	ExpressionTypeCheckCast                                 // []byte(str)
	ExpressionTypeRange                                     // for range
	ExpressionTypeSlice                                     // arr[0:2]
	ExpressionTypeMap                                       // map literal
	ExpressionTypeTypeAssert                                // a.(Object)
	ExpressionTypeQuestion                                  // true ? a : b
	ExpressionTypeGlobal                                    // global.XXX
	ExpressionTypeParenthesis                               // ( a )
	ExpressionTypeVArgs                                     // a ...
	ExpressionTypeDot                                       // .
)

type Expression struct {
	//checkRangeCalled bool
	Type ExpressionTypeKind
	/*
		only for global variable definition
		public hello := "hai...."
	*/
	IsPublic              bool // for global
	IsGlobal              bool
	IsCompileAuto         bool // compile auto expression
	Value                 *Type
	MultiValues           []*Type
	Pos                   *Pos
	Data                  interface{}
	IsStatementExpression bool
	Op                    string
	Lefts                 []*Expression // left values
	AsSubForNegative      *Expression
}

func (this *Expression) binaryExpressionDependOnSub() *Expression {
	switch this.Type {
	case ExpressionTypeAdd:
		bin := this.Data.(*ExpressionBinary)
		// 0 + a
		if bin.Left.isNumber() && bin.Left.getDoubleValue() == 0 {
			return bin.Right
		}
		// a + 0
		if bin.Right.isNumber() && bin.Right.getDoubleValue() == 0 {
			return bin.Left
		}
	case ExpressionTypeSub:
		// a - 0
		bin := this.Data.(*ExpressionBinary)
		if bin.Right.isNumber() && bin.Right.getDoubleValue() == 0 {
			return bin.Left
		}

	case ExpressionTypeMul:
		// a * 0 == 0
		bin := this.Data.(*ExpressionBinary)
		if bin.Right.isNumber() && bin.Right.getDoubleValue() == 0 {
			return bin.Right
		}
		// 0 * a == 0
		if bin.Left.isNumber() && bin.Left.getDoubleValue() == 0 {
			return bin.Left
		}
		// a * 1 == a
		if bin.Right.isNumber() && bin.Right.getDoubleValue() == 1 {
			return bin.Left
		}
		// 1 * a == a
		if bin.Left.isNumber() && bin.Left.getDoubleValue() == 1 {
			return bin.Right
		}
	case ExpressionTypeDiv:
		// a / 1 == a
		bin := this.Data.(*ExpressionBinary)
		if bin.Right.isNumber() && bin.Right.getDoubleValue() == 1 {
			return bin.Left
		}
	case ExpressionTypeEq:
		bin := this.Data.(*ExpressionBinary)
		if bin.Left.Value.Type == VariableTypeBool {
			// true == a
			if bin.Left.isBoolLiteral(true) {
				return bin.Right
			}
			// a == true
			if bin.Right.isBoolLiteral(true) {
				return bin.Left
			}
		}
	case ExpressionTypeNe:
		bin := this.Data.(*ExpressionBinary)
		if bin.Left.Value.Type == VariableTypeBool {
			// false != a
			if bin.Left.isBoolLiteral(false) {
				return bin.Right
			}
			// a != false
			if bin.Right.isBoolLiteral(false) {
				return bin.Left
			}
		}
	case ExpressionTypeLogicalAnd:
		bin := this.Data.(*ExpressionBinary)
		// true && a
		if bin.Left.isBoolLiteral(true) {
			return bin.Right
		}
		// a && true
		if bin.Right.isBoolLiteral(true) {
			return bin.Left
		}
	case ExpressionTypeLogicalOr:
		bin := this.Data.(*ExpressionBinary)
		// false || a
		if bin.Left.isBoolLiteral(false) {
			return bin.Right
		}
		// a || false
		if bin.Right.isBoolLiteral(false) {
			return bin.Left
		}
	}
	return nil
}

func (this *Expression) isBoolLiteral(b bool) bool {
	if this.Type != ExpressionTypeBool {
		return false
	}
	return this.Data.(bool) == b
}

func (this *Expression) isRelation() bool {
	return this.Type == ExpressionTypeEq ||
		this.Type == ExpressionTypeNe ||
		this.Type == ExpressionTypeGe ||
		this.Type == ExpressionTypeGt ||
		this.Type == ExpressionTypeLe ||
		this.Type == ExpressionTypeLt
}

/*
	1 > 2
	'a' > 'b'
	1s > 2s
*/
func (this *Expression) Is2IntCompare() bool {
	if this.isRelation() == false {
		return false
	}
	bin := this.Data.(*ExpressionBinary)
	i1 := bin.Left.Value.isInteger() && bin.Left.Value.Type != VariableTypeLong
	i2 := bin.Right.Value.isInteger() && bin.Right.Value.Type != VariableTypeLong
	return i1 && i2
}

/*
	a == null
*/
func (this *Expression) IsCompare2Null() bool {
	if this.isRelation() == false {
		return false
	}
	bin := this.Data.(*ExpressionBinary)
	return bin.Left.Type == ExpressionTypeNull ||
		bin.Right.Type == ExpressionTypeNull
}

/*
	a > "b"
*/
func (this *Expression) Is2StringCompare() bool {
	if this.isRelation() == false {
		return false
	}
	bin := this.Data.(*ExpressionBinary)
	return bin.Left.Value.Type == VariableTypeString
}

/*
	var a ,b []int
	a == b
*/
func (this *Expression) Is2PointerCompare() bool {
	if this.isRelation() == false {
		return false
	}
	bin := this.Data.(*ExpressionBinary)
	return bin.Left.Value.IsPointer()
}

func (this *Expression) convertTo(to *Type) {
	c := &ExpressionTypeConversion{}
	c.Expression = &Expression{}
	c.Expression.Op = "checkcast"
	*c.Expression = *this // copy
	c.Type = to
	this.Value = to
	this.Type = ExpressionTypeCheckCast
	this.IsCompileAuto = true
	this.Data = c
}

func (this *Expression) convertToNumberType(typ VariableTypeKind) {
	if this.isLiteral() {
		this.convertLiteralToNumberType(typ)
		this.Value = &Type{
			Type: typ,
			Pos:  this.Pos,
		}
	} else {
		this.convertTo(&Type{
			Pos:  this.Pos,
			Type: typ,
		})
	}
}

type ExpressionTypeAssert struct {
	ExpressionTypeConversion
	MultiValueContext bool
}

/*
	const spread
*/
func (this *Expression) fromConst(c *Constant) {
	this.Op = c.Name
	switch c.Type.Type {
	case VariableTypeBool:
		this.Type = ExpressionTypeBool
		this.Data = c.Value.(bool)
	case VariableTypeByte:
		this.Type = ExpressionTypeByte
		this.Data = c.Value.(int64)
	case VariableTypeShort:
		this.Type = ExpressionTypeShort
		this.Data = c.Value.(int64)
	case VariableTypeChar:
		this.Type = ExpressionTypeChar
		this.Data = c.Value.(int64)
	case VariableTypeInt:
		this.Type = ExpressionTypeInt
		this.Data = c.Value.(int64)
	case VariableTypeLong:
		this.Type = ExpressionTypeLong
		this.Data = c.Value.(int64)
	case VariableTypeFloat:
		this.Type = ExpressionTypeFloat
		this.Data = c.Value.(float32)
	case VariableTypeDouble:
		this.Type = ExpressionTypeDouble
		this.Data = c.Value.(float64)
	case VariableTypeString:
		this.Type = ExpressionTypeString
		this.Data = c.Value.(string)
	}
}

type ExpressionQuestion struct {
	Selection *Expression
	True      *Expression
	False     *Expression
}

type ExpressionSlice struct {
	ExpressionOn *Expression
	Start, End   *Expression
}

func (this *Expression) isLiteral() bool {
	return this.Type == ExpressionTypeBool ||
		this.Type == ExpressionTypeString ||
		this.isNumber()
}

/*
	valid for condition
*/
func (this *Expression) canBeUsedAsCondition() error {
	if this.Type == ExpressionTypeNull ||
		this.Type == ExpressionTypeBool ||
		this.Type == ExpressionTypeByte ||
		this.Type == ExpressionTypeShort ||
		this.Type == ExpressionTypeInt ||
		this.Type == ExpressionTypeLong ||
		this.Type == ExpressionTypeFloat ||
		this.Type == ExpressionTypeDouble ||
		this.Type == ExpressionTypeString ||
		this.Type == ExpressionTypeArray ||
		this.Type == ExpressionTypeLogicalOr ||
		this.Type == ExpressionTypeLogicalAnd ||
		this.Type == ExpressionTypeOr ||
		this.Type == ExpressionTypeAnd ||
		this.Type == ExpressionTypeXor ||
		this.Type == ExpressionTypeLsh ||
		this.Type == ExpressionTypeRsh ||
		this.Type == ExpressionTypeAdd ||
		this.Type == ExpressionTypeSub ||
		this.Type == ExpressionTypeMul ||
		this.Type == ExpressionTypeDiv ||
		this.Type == ExpressionTypeMod ||
		this.Type == ExpressionTypeEq ||
		this.Type == ExpressionTypeNe ||
		this.Type == ExpressionTypeGe ||
		this.Type == ExpressionTypeGt ||
		this.Type == ExpressionTypeLe ||
		this.Type == ExpressionTypeLt ||
		this.Type == ExpressionTypeIndex ||
		this.Type == ExpressionTypeSelection ||
		this.Type == ExpressionTypeMethodCall ||
		this.Type == ExpressionTypeFunctionCall ||
		this.Type == ExpressionTypeIncrement ||
		this.Type == ExpressionTypeDecrement ||
		this.Type == ExpressionTypePrefixIncrement ||
		this.Type == ExpressionTypePrefixDecrement ||
		this.Type == ExpressionTypeNegative ||
		this.Type == ExpressionTypeNot ||
		this.Type == ExpressionTypeBitwiseNot ||
		this.Type == ExpressionTypeIdentifier ||
		this.Type == ExpressionTypeNew ||
		this.Type == ExpressionTypeCheckCast ||
		this.Type == ExpressionTypeSlice ||
		this.Type == ExpressionTypeMap ||
		this.Type == ExpressionTypeQuestion {
		return nil
	}
	return fmt.Errorf("%s cannot use '%s' as condition",
		this.Pos.ErrMsgPrefix(), this.Op)
}

func (this *Expression) canBeUsedAsStatement() error {
	if this.Type == ExpressionTypeVarAssign ||
		this.Type == ExpressionTypeAssign ||
		this.Type == ExpressionTypeFunctionCall ||
		this.Type == ExpressionTypeMethodCall ||
		this.Type == ExpressionTypeFunctionLiteral ||
		this.Type == ExpressionTypePlusAssign ||
		this.Type == ExpressionTypeMinusAssign ||
		this.Type == ExpressionTypeMulAssign ||
		this.Type == ExpressionTypeDivAssign ||
		this.Type == ExpressionTypeModAssign ||
		this.Type == ExpressionTypeAndAssign ||
		this.Type == ExpressionTypeOrAssign ||
		this.Type == ExpressionTypeXorAssign ||
		this.Type == ExpressionTypeLshAssign ||
		this.Type == ExpressionTypeRshAssign ||
		this.Type == ExpressionTypeIncrement ||
		this.Type == ExpressionTypeDecrement ||
		this.Type == ExpressionTypePrefixIncrement ||
		this.Type == ExpressionTypePrefixDecrement ||
		this.Type == ExpressionTypeVar ||
		this.Type == ExpressionTypeConst {
		return nil
	}
	return fmt.Errorf("%s expression '%s' evaluate but not used",
		this.Pos.ErrMsgPrefix(), this.Op)
}

func (this *Expression) isNumber() bool {
	return this.isInteger() ||
		this.isFloat()
}

func (this *Expression) isInteger() bool {
	return this.Type == ExpressionTypeByte ||
		this.Type == ExpressionTypeShort ||
		this.Type == ExpressionTypeInt ||
		this.Type == ExpressionTypeLong ||
		this.Type == ExpressionTypeChar
}
func (this *Expression) isFloat() bool {
	return this.Type == ExpressionTypeFloat ||
		this.Type == ExpressionTypeDouble
}

func (this *Expression) isEqOrNe() bool {
	return this.Type == ExpressionTypeEq ||
		this.Type == ExpressionTypeNe
}

/*
	check out this expression is increment or decrement
*/
func (this *Expression) IsIncrement() bool {
	if this.Type == ExpressionTypeIncrement ||
		this.Type == ExpressionTypePrefixIncrement ||
		this.Type == ExpressionTypeDecrement ||
		this.Type == ExpressionTypePrefixDecrement {
	} else {
		panic("not increment or decrement at all")
	}
	return this.Type == ExpressionTypeIncrement ||
		this.Type == ExpressionTypePrefixIncrement
}

/*
	k,v := range arr
	k,v = range arr
*/
func (this *Expression) canBeUsedForRange() bool {
	if this.Type != ExpressionTypeAssign &&
		this.Type != ExpressionTypeVarAssign {
		return false
	}
	bin := this.Data.(*ExpressionBinary)
	if bin.Right.Type == ExpressionTypeRange {
		return true
	}
	if bin.Right.Type == ExpressionTypeList {
		t := bin.Right.Data.([]*Expression)
		if len(t) == 1 && t[0].Type == ExpressionTypeRange {
			return true
		}
	}
	return false
}

func (this *Expression) HaveMultiValue() bool {
	if this.Type == ExpressionTypeFunctionCall ||
		this.Type == ExpressionTypeMethodCall ||
		this.Type == ExpressionTypeTypeAssert {
		return len(this.MultiValues) > 1
	}
	return false

}

type CallArgs []*Expression // f(1,2)

type ExpressionFunctionCall struct {
	BuildInFunctionMeta      interface{} // for build in function only
	Expression               *Expression
	Args                     CallArgs
	VArgs                    *CallVariableArgs
	Function                 *Function
	ParameterTypes           []*Type // for template function
	TemplateFunctionCallPair *TemplateFunctionInstance
}

type ExpressionMethodCall struct {
	Class              *Class // for object or class
	Expression         *Expression
	Args               CallArgs
	VArgs              *CallVariableArgs
	Name               string
	Method             *ClassMethod
	FieldMethodHandler *ClassField
	/*
		unSupport !!!!!!
	*/
	ParameterTypes                []*Type
	PackageFunction               *Function
	PackageGlobalVariableFunction *Variable
}

type ExpressionVar struct {
	Type       *Type
	Variables  []*Variable
	InitValues []*Expression
}

type ExpressionVarAssign struct {
	Lefts            []*Expression
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
	Comment  string
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
	PackageFunction *Function    // expression is package , pack function to method handle
	PackageVariable *Variable    // expression is package , get package variable
	PackageEnumName *EnumName    // expression is package , get enumName
}

type ExpressionNew struct {
	Type         *Type
	Args         CallArgs
	Construction *ClassMethod
	VArgs        *CallVariableArgs
}

type ExpressionMap struct {
	Type          *Type
	KeyValuePairs []*ExpressionKV
}

type ExpressionKV struct {
	Key   *Expression
	Value *Expression
}

/*
	for some general purpose
*/
type ExpressionBinary struct {
	Left  *Expression
	Right *Expression
}

// for package jvm
type ExpressionAssign struct {
	Lefts  []*Expression
	Values []*Expression
}

type ExpressionArray struct {
	Type        *Type
	Expressions []*Expression
}

func (this *Expression) IsIdentifier(identifier string) bool {
	if this.Type != ExpressionTypeIdentifier {
		return false
	}
	return this.Data.(*ExpressionIdentifier).Name == identifier
}
