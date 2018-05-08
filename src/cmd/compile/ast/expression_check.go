package ast

import (
	"fmt"
)

func (e *Expression) check(block *Block) (Types []*VariableType, errs []error) {
	if e == nil {
		return nil, []error{}
	}
	_, err := e.constFold()
	if err != nil {
		return nil, []error{err}
	}

	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_NULL:
		Types = []*VariableType{
			{
				Typ: VARIABLE_TYPE_NULL,
				Pos: e.Pos,
			},
		}
		e.Value = Types[0]
	case EXPRESSION_TYPE_BOOL:
		Types = []*VariableType{
			{
				Typ: VARIABLE_TYPE_BOOL,
				Pos: e.Pos,
			},
		}
		e.Value = Types[0]
	case EXPRESSION_TYPE_BYTE:
		Types = []*VariableType{{
			Typ: VARIABLE_TYPE_BYTE,
			Pos: e.Pos,
		},
		}
		e.Value = Types[0]
	case EXPRESSION_TYPE_SHORT:
		Types = []*VariableType{{
			Typ: VARIABLE_TYPE_SHORT,
			Pos: e.Pos,
		},
		}
		e.Value = Types[0]

	case EXPRESSION_TYPE_INT:
		Types = []*VariableType{{
			Typ: VARIABLE_TYPE_INT,
			Pos: e.Pos,
		},
		}

		e.Value = Types[0]
	case EXPRESSION_TYPE_FLOAT:
		Types = []*VariableType{{
			Typ: VARIABLE_TYPE_FLOAT,
			Pos: e.Pos,
		},
		}
		e.Value = Types[0]
	case EXPRESSION_TYPE_DOUBLE:
		Types = []*VariableType{{
			Typ: VARIABLE_TYPE_DOUBLE,
			Pos: e.Pos,
		},
		}
		e.Value = Types[0]
	case EXPRESSION_TYPE_LONG:
		Types = []*VariableType{{
			Typ: VARIABLE_TYPE_LONG,
			Pos: e.Pos,
		},
		}
		e.Value = Types[0]
	case EXPRESSION_TYPE_STRING:
		Types = []*VariableType{{
			Typ: VARIABLE_TYPE_STRING,
			Pos: e.Pos,
		}}
		e.Value = Types[0]
	case EXPRESSION_TYPE_IDENTIFIER:
		tt, err := e.checkIdentiferExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			e.Value = tt
			Types = []*VariableType{tt}
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
			Types = []*VariableType{tt}
		}
		e.Value = tt
	case EXPRESSION_TYPE_MAP:
		tt := e.checkMapExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
		}
		e.Value = tt
	case EXPRESSION_TYPE_COLON_ASSIGN:
		e.checkColonAssignExpression(block, &errs)
		e.Value = mkVoidType(e.Pos)
		Types = []*VariableType{e.Value}
	case EXPRESSION_TYPE_ASSIGN:
		tt := e.checkAssignExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
		}
		e.Value = tt
	case EXPRESSION_TYPE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_DECREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_DECREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
		}
		e.Value = tt
	case EXPRESSION_TYPE_CONST: // no return value
		errs = e.checkConst(block)
		Types = []*VariableType{mkVoidType(e.Pos)}
		e.Value = Types[0]
	case EXPRESSION_TYPE_VAR:
		e.checkVarExpression(block, &errs)
		Types = []*VariableType{mkVoidType(e.Pos)}
		e.Value = Types[0]
	case EXPRESSION_TYPE_FUNCTION_CALL:
		Types = e.checkFunctionCallExpression(block, &errs)
		e.Values = Types
		if len(Types) > 0 {
			e.Value = Types[0]
		}
		if len(Types) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_METHOD_CALL:
		Types = e.checkMethodCallExpression(block, &errs)
		e.Values = Types
		if len(Types) > 0 {
			e.Value = Types[0]
		}
		if len(Types) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_TYPE_ASSERT:
		Types = e.checkTypeAssert(block, &errs)
		e.Values = Types
		if len(Types) > 0 {
			e.Value = Types[0]
		}
		if len(Types) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case EXPRESSION_TYPE_NOT:
		fallthrough
	case EXPRESSION_TYPE_NEGATIVE:
		fallthrough
	case EXPRESSION_TYPE_BITWISE_NOT:
		tt := e.checkUnaryExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
		}
		e.Value = tt
	case EXPRESSION_TYPE_INDEX:
		tt := e.checkIndexExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
			e.Value = tt
		}
	case EXPRESSION_TYPE_DOT:
		tt := e.checkDotExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
			e.Value = tt
		}
	case EXPRESSION_TYPE_CHECK_CAST:
		tt := e.checkTypeConvertionExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
			e.Value = tt
		}
	case EXPRESSION_TYPE_NEW:
		tt := e.checkNewExpression(block, &errs)
		if tt != nil {
			Types = []*VariableType{tt}
			e.Value = tt
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
			Types = []*VariableType{tt}
		}
		e.Value = tt
	case EXPRESSION_TYPE_RANGE:
		errs = append(errs, fmt.Errorf("%s range is only work with 'for' statement",
			errMsgPrefix(e.Pos)))
	case EXPRESSION_TYPE_SLICE:
		tt := e.checkSlice(block, &errs)
		e.Value = tt
		if tt != nil {
			Types = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_ARRAY:
		tt := e.checkArray(block, &errs)
		e.Value = tt
		if tt != nil {
			Types = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_FUNCTION:
		errs = append(errs, fmt.Errorf("%s cannot use function as  a expression",
			errMsgPrefix(e.Pos)))
	case EXPRESSION_TYPE_LIST:
		errs = append(errs, fmt.Errorf("%s cannot have expression '%s' at this scope,"+
			"this may be cause be compiler error,please contact the author",
			errMsgPrefix(e.Pos), e.OpName()))
	default:
		panic(fmt.Sprintf("unhandled type inference:%s", e.OpName()))
	}
	return Types, errs
}

func (e *Expression) mustBeOneValueContext(ts []*VariableType) (*VariableType, error) {
	if len(ts) == 0 {
		return nil, nil // no-type,no error
	}
	var err error
	if len(ts) > 1 {
		err = fmt.Errorf("%s multi value in single value context", errMsgPrefix(e.Pos))
	}
	return ts[0], err
}

func (e *Expression) checkBuildinFunctionCall(block *Block, errs *[]error, f *Function, args []*Expression) []*VariableType {
	callargsTypes := checkRightValuesValid(checkExpressions(block, args, errs), errs)
	length := len(*errs)
	f.buildChecker(f, e.Data.(*ExpressionFunctionCall), block, errs, callargsTypes, e.Pos)
	if len(*errs) == length {
		//special case ,avoid null pointer
		return f.Typ.retTypes(e.Pos)
	}
	return nil //
}
