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
	if no.Typ.Typ == VARIABLE_TYPE_ARRAY {
		return e.checkNewArrayExpression(block, no, errs)
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

func (e *Expression) checkNewArrayExpression(block *Block, newArray *ExpressionNew, errs *[]error) *VariableType {
	ret := &VariableType{}
	*ret = *newArray.Typ
	ret.Typ = VARIABLE_TYPE_ARRAY_INSTANCE
	ret.Pos = e.Pos
	if len(newArray.Args) > 1 { // 0 and 1 is accpect
		*errs = append(*errs, fmt.Errorf("%s new array must have one int argument"))
		newArray.Args = []*Expression{} // reset to 0,continue to analyse
	}
	if len(newArray.Args) == 0 {
		ee := &Expression{}
		ee.Typ = EXPRESSION_TYPE_INT
		ee.Data = int32(0)
	}
	ts := checkRightValuesValid(checkExpressions(block, newArray.Args, errs), errs)
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t.Typ != VARIABLE_TYPE_INT {
		*errs = append(*errs, fmt.Errorf("%s argument must be 'int'", errMsgPrefix(t.Pos)))
	}
	//no further checks
	return ret
}
