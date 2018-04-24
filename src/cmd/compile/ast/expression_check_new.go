package ast

import (
	"fmt"
)

func (e *Expression) checkNewExpression(block *Block, errs *[]error) *VariableType {
	no := e.Data.(*ExpressionNew)
	err := no.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		fmt.Println(err)
		return nil
	}
	if no.Typ.Typ == VARIABLE_TYPE_MAP {
		return e.checkNewMapExpression(block, no, errs)
	}
	if no.Typ.Typ == VARIABLE_TYPE_ARRAY {
		return e.checkNewArrayExpression(block, no, errs)
	}
	ret := &VariableType{}
	*ret = *no.Typ
	ret.Typ = VARIABLE_TYPE_OBJECT
	ret.Pos = e.Pos
	// new object
	if no.Typ.Typ != VARIABLE_TYPE_OBJECT {
		*errs = append(*errs, fmt.Errorf("%s cannot have new on type '%s'",
			errMsgPrefix(e.Pos), no.Typ.TypeString()))
		return ret
	}
	args := checkExpressions(block, no.Args, errs)
	ms, matched, err := no.Typ.Class.matchContructionFunction(args, &no.Args)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
		return ret
	}
	if matched {
		no.Construction = ms[0]
		return ret
	}
	if len(ms) == 0 {
		*errs = append(*errs, fmt.Errorf("%s  'construction' not found",
			errMsgPrefix(e.Pos)))
	} else {
		*errs = append(*errs, msNotMatchError(e.Pos, "constructor", ms, args))
	}
	return ret
}

func (e *Expression) checkNewMapExpression(block *Block, newMap *ExpressionNew,
	errs *[]error) *VariableType {
	if len(newMap.Args) > 0 {
		*errs = append(*errs, fmt.Errorf("%s new map expect no arguments",
			errMsgPrefix(newMap.Args[0].Pos)))
	}
	tt := newMap.Typ.Clone()
	tt.Pos = e.Pos
	return tt
}

func (e *Expression) checkNewArrayExpression(block *Block, newArray *ExpressionNew,
	errs *[]error) *VariableType {
	ret := newArray.Typ.Clone() // clone the type
	ret.Pos = e.Pos
	if len(newArray.Args) > 1 { // 0 and 1 is accpect
		*errs = append(*errs,
			fmt.Errorf("%s new array must have one int argument",
				errMsgPrefix(newArray.Args[0].Pos)))
		newArray.Args = []*Expression{} // reset to 0,continue to analyse
	}
	if len(newArray.Args) == 0 { // 0 is default
		ee := &Expression{}
		ee.Typ = EXPRESSION_TYPE_INT
		ee.Data = int32(0)
		newArray.Args = []*Expression{ee}
	}
	ts := checkRightValuesValid(checkExpressions(block, newArray.Args, errs), errs)
	amount, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if amount == nil {
		return ret
	}
	if amount.Typ != VARIABLE_TYPE_INT {
		*errs = append(*errs, fmt.Errorf("%s argument must be 'int',but '%s'",
			errMsgPrefix(amount.Pos), amount.TypeString()))
	}
	//no further checks
	return ret
}
