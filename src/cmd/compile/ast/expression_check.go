package ast

import (
	"fmt"
)

func (e *Expression) check(block *Block) (t *VariableType, errs []error) {
	is, typ, data, err := e.getConstValue()
	if err != nil {
		return nil, []error{err}
	}
	if is {
		e.Typ = typ
		e.Data = data
	}
	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		t = &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
		}
	case EXPRESSION_TYPE_BYTE:
		t = &VariableType{
			Typ: VARIABLE_TYPE_BYTE,
		}
	case EXPRESSION_TYPE_INT:
		t = &VariableType{
			Typ: VARIABLE_TYPE_INT,
		}
	case EXPRESSION_TYPE_FLOAT:
		t = &VariableType{
			Typ: VARIABLE_TYPE_FLOAT,
		}
	case EXPRESSION_TYPE_STRING:
		t = &VariableType{
			Typ: VARIABLE_TYPE_STRING,
		}
	case EXPRESSION_TYPE_IDENTIFIER:
		t, err = e.checkIdentiferExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		return t, errs
		//binaries
	case EXPRESSION_TYPE_LOGICAL_OR:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_LOGICAL_AND:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_OR:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_AND:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_LEFT_SHIFT:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_RIGHT_SHIFT:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_ASSIGN:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_COLON_ASSIGN:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_PLUS_ASSIGN:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_MINUS_ASSIGN:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_MUL_ASSIGN:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_DIV_ASSIGN:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_MOD_ASSIGN:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_EQ:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_NE:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_GE:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_GT:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_LE:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_LT:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_ADD:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_SUB:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_MUL:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_DIV:
		return e.checkBinaryExpression(block)
	case EXPRESSION_TYPE_MOD:
		return e.checkBinaryExpression(block)
	default:
		panic(fmt.Sprintf("unhandled type inference:%s", e.OpName()))
	}
	return
}

func (e *Expression) checkIdentiferExpression(block *Block) (t *VariableType, err error) {
	name := e.Data.(string)
	d, err := block.searchByName(name)
	if err != nil {
		return nil, err
	}
	switch d.(type) {
	case []*Function:
		if len(d.([]*Function)) > 1 {
			return nil, fmt.Errorf("%s %s is defined as function more than one time", errMsgPrefix(e.Pos), name)
		}
		return &VariableType{
			Typ:          VARIALBE_TYPE_FUNCTION,
			FunctionType: (d.([]*Function))[0].Typ,
		}, nil
	case *VariableDefinition:
		return nil, nil
	}
	return nil, nil
}
func (e *Expression) checkBinaryExpression(block *Block) (t *VariableType, errs []error) {
	errs = []error{}
	binary := e.Data.(*ExpressionBinary)
	t1, err := binary.Left.check(block)
	t2, err := binary.Right.check(block)

	return nil, nil
}
