package ast

import (
	"fmt"
)

func (e *Expression) checkSlice(block *Block, errs *[]error) *Type {
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
	if startT != nil && startT.Type == VariableTypeLong {
		slice.Start.ConvertToNumber(VariableTypeInt)
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
	if endT != nil && endT.Type == VariableTypeLong {
		slice.End.ConvertToNumber(VariableTypeInt)
	}

	t, es := slice.Array.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	if t.Type != VariableTypeArray {
		*errs = append(*errs, fmt.Errorf("%s cannot have slice on '%s'",
			errMsgPrefix(slice.Array.Pos), t.TypeString()))
	}
	tt := t.Clone()
	tt.Pos = e.Pos
	return tt
}
