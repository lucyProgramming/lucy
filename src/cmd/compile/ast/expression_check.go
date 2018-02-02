package ast

import (
	"fmt"
)

func (e *Expression) check(block *Block) (t []*VariableType, errs []error) {
	is, typ, data, err := e.getConstValue()
	fmt.Println(e.OpName())
	if err != nil {
		return nil, []error{err}
	}
	if is {
		e.Typ = typ
		e.Data = data
	}
	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_NULL:
		t = []*VariableType{
			{
				Typ: VARIABLE_TYPE_NULL,
				Pos: e.Pos,
			},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_BOOL:
		t = []*VariableType{
			{
				Typ: VARIABLE_TYPE_BOOL,
				Pos: e.Pos,
			},
		}
		e.VariableType = t[0]

	case EXPRESSION_TYPE_BYTE:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_BYTE,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_INT:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_INT,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_FLOAT:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_FLOAT,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_DOUBLE:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_DOUBLE,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_LONG:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_LONG,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_STRING:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_STRING,
			Pos: e.Pos,
		}}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_IDENTIFIER:
		tt, err := e.checkIdentiferExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			e.VariableType = tt
			t = []*VariableType{tt}
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
	case EXPRESSION_TYPE_LEFT_SHIFT:
		fallthrough
	case EXPRESSION_TYPE_RIGHT_SHIFT:
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
			t = []*VariableType{tt}
		}
		e.VariableType = tt
	case EXPRESSION_TYPE_COLON_ASSIGN:
		e.checkColonAssignExpression(block, &errs)
	case EXPRESSION_TYPE_ASSIGN:
		tt := e.checkAssignExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
		e.VariableType = tt
	case EXPRESSION_TYPE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_DECREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_DECREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
		e.VariableType = tt
	case EXPRESSION_TYPE_CONST:
		e.checkConstExpression(block, &errs)
	case EXPRESSION_TYPE_VAR:
		e.checkVarExpression(block, &errs)
	case EXPRESSION_TYPE_FUNCTION_CALL:
		t = e.checkFunctionCallExpression(block, &errs)
		e.VariableTypes = t
		if len(t) > 1 {
			block.InheritedAttribute.function.mkArrayListVarForMultiReturn()
		}
	case EXPRESSION_TYPE_METHOD_CALL:
		t = e.checkMethodCallExpression(block, &errs)
		e.VariableTypes = t
		if len(t) > 1 {
			block.InheritedAttribute.function.mkArrayListVarForMultiReturn()
		}
	case EXPRESSION_TYPE_NOT:
		fallthrough
	case EXPRESSION_TYPE_NEGATIVE:
		tt := e.checkUnaryExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
		e.VariableType = tt
	case EXPRESSION_TYPE_INDEX:
		fallthrough
	case EXPRESSION_TYPE_DOT:
		tt := e.checkIndexExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_CONVERTION_TYPE:
		tt := e.checkTypeConvertionExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_NEW:
		tt := e.checkNewExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
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
		tt := e.checkOpAssignExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
		e.VariableType = tt
	case EXPRESSION_TYPE_FUNCTION:
		tt := e.checkFunctionExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}

	default:
		panic(fmt.Sprintf("unhandled type inference:%s", e.OpName()))
	}
	return t, errs
}
func (e *Expression) checkFunctionExpression(block *Block, errs *[]error) *VariableType {
	f := e.Data.(*Function)
	f.MkVariableType()
	*errs = append(*errs, f.check(block)...)
	return &f.VariableType
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

func (e *Expression) checkTypeConvertionExpression(block *Block, errs *[]error) *VariableType {
	c := e.Data.(*ExpressionTypeConvertion)
	ts, es := c.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return nil
	}
	return nil
}

func (e *Expression) checkBuildinFunctionCall(block *Block, errs *[]error, f *Function, args []*Expression) []*VariableType {
	callargsTypes := checkRightValuesValid(checkExpressions(block, args, errs), errs)
	f.callchecker(block, errs, callargsTypes, e.Pos)
	return f.Typ.ReturnList.retTypes(e.Pos)
}
