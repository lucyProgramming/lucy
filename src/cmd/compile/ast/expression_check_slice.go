package ast

import (
	"fmt"
)

func (e *Expression) checkSlice(block *Block, errs *[]error) *VariableType {
	slice := e.Data.(*ExpressionSlice)
	ts, es := slice.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := slice.Expression.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return nil
	}
	if t.Typ != VARIABLE_TYPE_ARRAY_INSTANCE {
		*errs = append(*errs, fmt.Errorf("%s cannot have slice on '%s'", errMsgPrefix(slice.Expression.Pos), t.TypeString()))
	}
	//start
	startTs, es := slice.Start.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	startT, err := slice.Start.mustBeOneValueContext(startTs)
	if err != nil {
		*errs = append(*errs, err)
	}
	if startT != nil && startT.IsInteger() == false {
		*errs = append(*errs, fmt.Errorf("%s slice start must be integer,nut '%s'", errMsgPrefix(slice.Start.Pos), startT.TypeString()))
	}
	//end
	endTs, es := slice.End.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	endT, err := slice.End.mustBeOneValueContext(endTs)
	if err != nil {
		*errs = append(*errs, err)
	}
	if endT != nil && endT.IsInteger() == false {
		*errs = append(*errs, fmt.Errorf("%s slice end must be integer,nut '%s'", errMsgPrefix(slice.End.Pos), endT.TypeString()))
	}
	tt := t.Clone()
	tt.Pos = e.Pos
	return tt
}
