package ast

import (
	"fmt"
)

func (e *Expression) checkNewExpression(block *Block, errs *[]error) *VariableType {
	no := e.Data.(*ExpressionNew)
	err := no.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		return nil
	}
	if no.Typ.Typ == VARIABLE_TYPE_CLASS {
		args := checkExpressions(block, no.Args, errs)
		f, accessable, err := no.Typ.Class.matchContructionFunction(args)
		if err != nil {
			*errs = append(*errs, err)
		} else {
			if !accessable {
				*errs = append(*errs, fmt.Errorf("%s construction method is private", errMsgPrefix(e.Pos)))
			}
		}
		no.Construction = f
		ret := &VariableType{}
		*ret = *no.Typ
		ret.Typ = VARIABLE_TYPE_OBJECT
		ret.Pos = e.Pos
		return ret
	} else {
		*errs = append(*errs, fmt.Errorf("%s only class type can be used by new", errMsgPrefix(e.Pos)))
		return nil
	}
}
