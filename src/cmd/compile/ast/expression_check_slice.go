package ast

import (
	"fmt"
)

func (e *Expression) checkSlice(block *Block, errs *[]error) *VariableType {
	slice := e.Data.(*ExpressionSlice)
	//start
	if slice.Start == nil {
		slice.Start = &Expression{}
		slice.Start.Pos = e.Pos
		slice.Start.Typ = EXPRESSION_TYPE_INT
		slice.Start.Data = int32(0)
	}

	startTs, es := slice.Start.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	startT, err := slice.Start.mustBeOneValueContext(startTs)
	if err != nil {
		*errs = append(*errs, err)
	}
	if startT != nil && startT.IsInteger() == false {
		*errs = append(*errs, fmt.Errorf("%s slice start must be integer,but '%s'",
			errMsgPrefix(slice.Start.Pos), startT.TypeString()))
	}
	if startT != nil && startT.Typ == VARIABLE_TYPE_LONG {
		slice.Start.ConvertToNumber(VARIABLE_TYPE_INT)
	}
	if slice.End == nil {
		slice.End = &Expression{}
		slice.End.Pos = e.Pos
		slice.End.Typ = EXPRESSION_TYPE_INT
		slice.End.Data = int32(-1) // special
	}
	endTs, es := slice.End.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	endT, err := slice.End.mustBeOneValueContext(endTs)
	if err != nil {
		*errs = append(*errs, err)
	}
	if endT != nil && endT.IsInteger() == false {
		*errs = append(*errs, fmt.Errorf("%s slice end must be integer,but '%s'",
			errMsgPrefix(slice.End.Pos), endT.TypeString()))
	}
	if endT != nil && endT.Typ == VARIABLE_TYPE_LONG {
		slice.End.ConvertToNumber(VARIABLE_TYPE_INT)
	}

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
	if t.Typ != VARIABLE_TYPE_ARRAY {
		*errs = append(*errs, fmt.Errorf("%s cannot have slice on '%s'",
			errMsgPrefix(slice.Expression.Pos), t.TypeString()))
	}
	tt := t.Clone()
	tt.Pos = e.Pos
	return tt
}
