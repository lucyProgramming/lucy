package ast

import (
	"fmt"
)

func (e *Expression) getLeftValue(block *Block) (t *VariableType, errs []error) {
	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_IDENTIFIER:
		name := e.Data.(*ExpressionIdentifer)
		d := block.searchByName(name.Name)
		if d == nil {
			return nil, []error{fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), name.Name)}
		}
		switch d.(type) {
		case *VariableDefinition:
			t := d.(*VariableDefinition)
			t.CaptureLevel = 0
			return t.Typ, nil
		default:
			errs = append(errs, fmt.Errorf("%s identifier %s is not variable",
				errMsgPrefix(e.Pos), name.Name))
			return nil, []error{}
		}
	case EXPRESSION_TYPE_INDEX:
		return e.checkIndexExpression(block, &errs), errs
	case EXPRESSION_TYPE_DOT:
		return e.checkIndexExpression(block, &errs), errs
	default:
		errs = append(errs, fmt.Errorf("%s %s cannot be used as left value",
			errMsgPrefix(e.Pos),
			e.OpName()))
		return nil, errs
	}
}
