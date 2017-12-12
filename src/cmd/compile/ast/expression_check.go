package ast

import (
	"fmt"
)

func (e *Expression) check(block *Block) (t []*VariableType, errs []error) {
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
		t = []*VariableType{
			&VariableType{
				Typ: VARIABLE_TYPE_BOOL,
			},
		}
	case EXPRESSION_TYPE_BYTE:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_BYTE,
		},
		}
	case EXPRESSION_TYPE_INT:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_INT,
		},
		}
	case EXPRESSION_TYPE_FLOAT:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_FLOAT,
		},
		}
	case EXPRESSION_TYPE_STRING:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_STRING,
		}}

	case EXPRESSION_TYPE_IDENTIFIER:
		tt, err := e.checkIdentiferExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		return []*VariableType{tt}, errs
		//binaries
	case EXPRESSION_TYPE_LOGICAL_OR:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LOGICAL_AND:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_OR:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_AND:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LEFT_SHIFT:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_RIGHT_SHIFT:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_ASSIGN:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_PLUS_ASSIGN:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MINUS_ASSIGN:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MUL_ASSIGN:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_DIV_ASSIGN:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MOD_ASSIGN:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_EQ:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_NE:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_GE:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_GT:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LE:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LT:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_ADD:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_SUB:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MUL:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_DIV:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MOD:
		tt, errs := e.checkBinaryExpression(block)
		return []*VariableType{tt}, errs
	default:
		panic(fmt.Sprintf("unhandled type inference:%s", e.OpName()))
	}
	return
}

func (e *Expression) checkIdentiferExpression(block *Block) (t *VariableType, err error) {
	//identifer := e.Data.(*ExprssionIdentifier)
	//d, err := block.searchByName(identifer.Identifer)
	//if err != nil {
	//	return nil, err
	//}
	//switch d.(type) {
	//case []*Function:
	//	if len(d.([]*Function)) > 1 {
	//		return nil, fmt.Errorf("%s %s is defined as function multi times", errMsgPrefix(e.Pos), identifer.Identifer)
	//	}
	//	f := d.([]*Function)
	//	identifer.Function = f
	//	return &VariableType{
	//		Typ:          VARIALBE_TYPE_FUNCTION,
	//		FunctionType: f[0].Typ,
	//	}, nil
	//case *VariableDefinition:
	//	t := d.(*VariableDefinition)
	//	identifer.Variable = t
	//	return t.Typ, nil
	//case *Const:
	//	t := d.(*Const)
	//	identifer.Const = true
	//	return t.Typ, nil
	//case *Enum:
	//
	//}
	return nil, nil
}
func (e *Expression) checkBinaryExpression(block *Block) (t *VariableType, errs []error) {
	errs = []error{}
	binary := e.Data.(*ExpressionBinary)
	_, err1 := binary.Left.check(block)
	_, err2 := binary.Right.check(block)
	if errsNotEmpty(err1) {
		errs = append(errs, err1...)
	}
	if errsNotEmpty(err2) {
		errs = append(errs, err2...)
	}

	return nil, errs
}
