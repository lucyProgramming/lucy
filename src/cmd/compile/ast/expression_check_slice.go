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
		slice.Start.Type = EXPRESSION_TYPE_INT
		slice.Start.Data = int32(0)
	}

	startT, es := slice.Start.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}

	if startT != nil && startT.IsInteger() == false {
		*errs = append(*errs, fmt.Errorf("%s slice start must be integer,but '%s'",
			errMsgPrefix(slice.Start.Pos), startT.TypeString()))
	}
	if startT != nil && startT.Type == VARIABLE_TYPE_LONG {
		slice.Start.ConvertToNumber(VARIABLE_TYPE_INT)
	}
	if slice.End == nil {
		slice.End = &Expression{}
		slice.End.Pos = e.Pos
		slice.End.Type = EXPRESSION_TYPE_INT
		slice.End.Data = int32(-1) // special  , end == arr.end
	}
	endT, es := slice.End.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if endT != nil && endT.IsInteger() == false {
		*errs = append(*errs, fmt.Errorf("%s slice end must be integer,but '%s'",
			errMsgPrefix(slice.End.Pos), endT.TypeString()))
	}
	if endT != nil && endT.Type == VARIABLE_TYPE_LONG {
		slice.End.ConvertToNumber(VARIABLE_TYPE_INT)
	}

	t, es := slice.SliceOn.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	if t.Type != VARIABLE_TYPE_ARRAY {
		*errs = append(*errs, fmt.Errorf("%s cannot have slice on '%s'",
			errMsgPrefix(slice.SliceOn.Pos), t.TypeString()))
	}
	tt := t.Clone()
	tt.Pos = e.Pos
	return tt
}
