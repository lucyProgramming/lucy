package ast

type ExpressionTypeKind int

const (
	_                             ExpressionTypeKind = iota // start with 1
	ExpressionTypeNull                                      // null
	ExpressionTypeBool                                      // true or false
	ExpressionTypeByte                                      // 'a' or 97b
	ExpressionTypeShort                                     // 100s
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
)

type Expression struct {
	Type ExpressionTypeKind
	/*
		only for global variable definition
		public hello := "hai...."
	*/
	IsPublic              bool
	IsCompileAuto         bool // compile auto expression
	Value                 *Type
	MultiValues           []*Type
	Pos                   *Pos
	Data                  interface{}
	IsStatementExpression bool
	Description           string
}

func (e *Expression) IsRelation() bool {
	return e.Type == ExpressionTypeEq ||
		e.Type == ExpressionTypeNe ||
		e.Type == ExpressionTypeGe ||
		e.Type == ExpressionTypeGt ||
		e.Type == ExpressionTypeLe ||
		e.Type == ExpressionTypeLt
}

/*
	1 > 2
	'a' > 'b'
	1s > 2s
*/
func (e *Expression) IsIntCompare() bool {
	if e.IsRelation() == false {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
	i1 := bin.Left.Value.IsInteger() && bin.Left.Value.Type != VariableTypeLong
	i2 := bin.Right.Value.IsInteger() && bin.Right.Value.Type != VariableTypeLong
	return i1 && i2
}

/*
	a == null
*/
func (e *Expression) IsNullCompare() bool {
	if e.IsRelation() == false {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
	return bin.Left.Type == ExpressionTypeNull ||
		bin.Right.Type == ExpressionTypeNull
}

/*
	a > "b"
*/
func (e *Expression) IsStringCompare() bool {
	if e.IsRelation() == false {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
	return bin.Left.Value.Type == VariableTypeString
}

/*
	var a ,b []int
	a == b
*/
func (e *Expression) IsPointerCompare() bool {
	if e.IsRelation() == false {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
	return bin.Left.Value.IsPointer()
}

func (e *Expression) ConvertTo(to *Type) {
	c := &ExpressionTypeConversion{}
	c.Expression = &Expression{}
	c.Expression.Description = "compilerAuto"
	*c.Expression = *e // copy
	c.Type = to
	e.Value = to
	e.Type = ExpressionTypeCheckCast
	e.IsCompileAuto = true
	e.Data = c
}

func (e *Expression) ConvertToNumber(typ VariableTypeKind) {
	if e.IsLiteral() {
		e.convertNumberLiteralTo(typ)
		e.Value = &Type{
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
	const spread
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

type ExpressionQuestion struct {
	Selection *Expression
	True      *Expression
	False     *Expression
}

type ExpressionSlice struct {
	ExpressionOn *Expression
	Start, End   *Expression
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
		e.Type == ExpressionTypeQuestion
}

func (e *Expression) canBeUsedAsStatement() bool {
	return e.Type == ExpressionTypeVarAssign ||
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
	if e.Type == ExpressionTypeIncrement ||
		e.Type == ExpressionTypePrefixIncrement ||
		e.Type == ExpressionTypeDecrement ||
		e.Type == ExpressionTypePrefixDecrement {
	} else {
		panic("not increment or decrement at all")
	}
	return e.Type == ExpressionTypeIncrement ||
		e.Type == ExpressionTypePrefixIncrement
}

func (e *Expression) isListAndMoreThanNElements(n int) bool {
	if e.Type != ExpressionTypeList {
		return false
	}
	return len(e.Data.([]*Expression)) > n
}

/*
	k,v := range arr
	k,v = range arr
*/
func (e *Expression) canBeUsedForRange() bool {
	if e.Type != ExpressionTypeAssign && e.Type != ExpressionTypeVarAssign {
		return false
	}
	bin := e.Data.(*ExpressionBinary)
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

func (e *Expression) HaveMultiValue() bool {
	if e.Type == ExpressionTypeFunctionCall ||
		e.Type == ExpressionTypeMethodCall ||
		e.Type == ExpressionTypeTypeAssert {
		return len(e.MultiValues) > 1
	}
	return false
}

type CallArgs []*Expression // f(1,2)

type ExpressionFunctionCall struct {
	BuildInFunctionMeta      interface{} // for build in function only
	Expression               *Expression
	Args                     CallArgs
	VArgs                    *CallVArgs
	Function                 *Function
	ParameterTypes           []*Type // for template function
	TemplateFunctionCallPair *TemplateFunctionCallPair
	FunctionPointer          *FunctionType
}

type ExpressionMethodCall struct {
	Class              *Class // for object or class
	Expression         *Expression
	Args               CallArgs
	VArgs              *CallVArgs
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
	VArgs        *CallVArgs
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
	Length      int // elements length
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
