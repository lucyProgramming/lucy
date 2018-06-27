package ast

import (
	"fmt"
)

func (e *Expression) check(block *Block) (Types []*Type, errs []error) {
	if e == nil {
		return nil, []error{}
	}
	_, err := e.constantFold()
	if err != nil {
		return nil, []error{err}
	}

	errs = []error{}
	switch e.Type {
	case ExpressionNull:
		Types = []*Type{
			{
				Type: VariableTypeNull,
				Pos:  e.Pos,
			},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_BOOL:
		Types = []*Type{
			{
				Type: VariableTypeBool,
				Pos:  e.Pos,
			},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_BYTE:
		Types = []*Type{{
			Type: VariableTypeByte,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_SHORT:
		Types = []*Type{{
			Type: VariableTypeShort,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_INT:
		Types = []*Type{{
			Type: VariableTypeInt,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_FLOAT:
		Types = []*Type{{
			Type: VariableTypeFloat,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_DOUBLE:
		Types = []*Type{{
			Type: VariableTypeDouble,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_LONG:
		Types = []*Type{{
			Type: VariableTypeLong,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_STRING:
		Types = []*Type{{
			Type: VariableTypeString,
			Pos:  e.Pos,
		}}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_IDENTIFIER:
		tt, err := e.checkIdentifierExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			e.ExpressionValue = tt
			Types = []*Type{tt}
		}
		//binaries
	case EXPRESSION_TYPE_LOGICAL_OR:
		fallthrough
	case EXPRESSION_TYPE_LOGICAL_AND:
		fallthrough
	case EXPRESSION_TYPE_OR:
		fallthrough
	case EXPRESSION_TYPE_AND:
		fallthrough
	case EXPRESSION_TYPE_XOR:
		fallthrough
	case EXPRESSION_TYPE_LSH:
		fallthrough
	case EXPRESSION_TYPE_RSH:
		fallthrough
	case EXPRESSION_TYPE_EQ:
		fallthrough
	case EXPRESSION_TYPE_NE:
		fallthrough
	case EXPRESSION_TYPE_GE:
		fallthrough
	case EXPRESSION_TYPE_GT:
		fallthrough
	case EXPRESSION_TYPE_LE:
		fallthrough
	case EXPRESSION_TYPE_LT:
		fallthrough
	case EXPRESSION_TYPE_ADD:
		fallthrough
	case EXPRESSION_TYPE_SUB:
		fallthrough
	case EXPRESSION_TYPE_MUL:
		fallthrough
	case EXPRESSION_TYPE_DIV:
		fallthrough
	case EXPRESSION_TYPE_MOD:
		tt := e.checkBinaryExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
		}
		e.ExpressionValue = tt
	case EXPRESSION_TYPE_MAP:
		tt := e.checkMapExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
		}
		e.ExpressionValue = tt
	case EXPRESSION_TYPE_COLON_ASSIGN:
		e.checkColonAssignExpression(block, &errs)
		e.ExpressionValue = mkVoidType(e.Pos)
		Types = []*Type{e.ExpressionValue}
		if t, ok := e.Data.(*ExpressionDeclareVariable); ok && t != nil && len(t.Variables) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_ASSIGN:
		tt := e.checkAssignExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
		}
		e.ExpressionValue = tt
		if e.Data.(*ExpressionBinary).Left.isListAndMoreThanNElements(1) {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_DECREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_DECREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
		}
		e.ExpressionValue = tt
	case EXPRESSION_TYPE_CONST: // no return value
		errs = e.checkConstant(block)
		Types = []*Type{mkVoidType(e.Pos)}
		e.ExpressionValue = Types[0]
	case EXPRESSION_TYPE_VAR:
		e.checkVarExpression(block, &errs)
		Types = []*Type{mkVoidType(e.Pos)}
		e.ExpressionValue = Types[0]
		if t, ok := e.Data.(*ExpressionDeclareVariable); ok && t != nil && len(t.Variables) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_FUNCTION_CALL:
		Types = e.checkFunctionCallExpression(block, &errs)
		e.ExpressionMultiValues = Types
		if len(Types) > 0 {
			e.ExpressionValue = Types[0]
		}
		if len(Types) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_METHOD_CALL:
		Types = e.checkMethodCallExpression(block, &errs)
		e.ExpressionMultiValues = Types
		if len(Types) > 0 {
			e.ExpressionValue = Types[0]
		}
		if len(Types) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_TYPE_ASSERT:
		Types = e.checkTypeAssert(block, &errs)
		e.ExpressionMultiValues = Types
		if len(Types) > 0 {
			e.ExpressionValue = Types[0]
		}
		if len(Types) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_NOT:
		fallthrough
	case EXPRESSION_TYPE_NEGATIVE:
		fallthrough
	case EXPRESSION_TYPE_BIT_NOT:
		tt := e.checkUnaryExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
		}
		e.ExpressionValue = tt
	case EXPRESSION_TYPE_TERNARY:
		tt := e.checkTernaryExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
		}
		e.ExpressionValue = tt
	case EXPRESSION_TYPE_INDEX:
		tt := e.checkIndexExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
			e.ExpressionValue = tt
		}
	case EXPRESSION_TYPE_SELECTION:
		tt := e.checkSelectionExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
			e.ExpressionValue = tt
		}
	case EXPRESSION_TYPE_CHECK_CAST:
		tt := e.checkTypeConversionExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
			e.ExpressionValue = tt
		}
	case EXPRESSION_TYPE_NEW:
		tt := e.checkNewExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
			e.ExpressionValue = tt
		}
	case EXPRESSION_TYPE_PLUS_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_MINUS_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_MUL_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_DIV_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_MOD_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_AND_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_OR_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_LSH_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_RSH_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_XOR_ASSIGN:
		tt := e.checkOpAssignExpression(block, &errs)
		if tt != nil {
			Types = []*Type{tt}
		}
		e.ExpressionValue = tt
	case EXPRESSION_TYPE_RANGE:
		errs = append(errs, fmt.Errorf("%s range is only work with 'for' statement",
			errMsgPrefix(e.Pos)))
	case EXPRESSION_TYPE_SLICE:
		tt := e.checkSlice(block, &errs)
		e.ExpressionValue = tt
		if tt != nil {
			Types = []*Type{tt}
		}
	case EXPRESSION_TYPE_ARRAY:
		tt := e.checkArray(block, &errs)
		e.ExpressionValue = tt
		if tt != nil {
			Types = []*Type{tt}
		}
	case EXPRESSION_TYPE_FUNCTION_LITERAL:
		f := e.Data.(*Function)
		errs = f.check(block)
		f.IsClosureFunction = f.Closure.NotEmpty(f)
		if f.Name != "" {
			err := block.Insert(f.Name, f.Pos, f)
			if err != nil {
				errs = append(errs, err)
			}
		}
		Types = make([]*Type, 1)
		Types[0] = &Type{
			Type:         VariableTypeFunction,
			Pos:          e.Pos,
			FunctionType: &f.Type,
		}
	case EXPRESSION_TYPE_LIST:
		errs = append(errs, fmt.Errorf("%s cannot have expression '%s' at this scope,"+
			"this may be cause be compiler error,please contact the author",
			errMsgPrefix(e.Pos), e.OpName()))
	default:
		panic(fmt.Sprintf("unhandled type:%v", e.OpName()))
	}
	return Types, errs
}

func (e *Expression) mustBeOneValueContext(ts []*Type) (*Type, error) {
	if len(ts) == 0 {
		return nil, nil // no-type,no error
	}
	var err error
	if len(ts) > 1 {
		err = fmt.Errorf("%s multi value in single value context", errMsgPrefix(e.Pos))
	}
	return ts[0], err
}

func (e *Expression) checkBuildInFunctionCall(block *Block, errs *[]error, f *Function, call *ExpressionFunctionCall) []*Type {
	callArgsTypes := checkRightValuesValid(checkExpressions(block, call.Args, errs), errs)
	length := len(*errs)
	f.buildInFunctionChecker(f, e.Data.(*ExpressionFunctionCall), block, errs, callArgsTypes, e.Pos)
	if len(*errs) == length {
		//special case ,avoid null pointer
		return f.Type.returnTypes(e.Pos)
	}
	return nil //
}
func (e *Expression) checkSingleValueContextExpression(block *Block) (*Type, []error) {
	ts, es := e.check(block)
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		if es == nil {
			es = []error{err}
		} else {
			es = append(es, err)
		}
	}
	return t, es
}
