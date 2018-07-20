package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func (e *Expression) checkSlice(block *Block, errs *[]error) *Type {
	slice := e.Data.(*ExpressionSlice)
	//start
	if slice.Start == nil {
		slice.Start = &Expression{}
		slice.Start.Pos = e.Pos
		slice.Start.Type = ExpressionTypeInt
		slice.Start.Data = int32(0)
	}
	startType, es := slice.Start.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if startType != nil && startType.IsInteger() == false {
		*errs = append(*errs, fmt.Errorf("%s slice start must be integer,but '%s'",
			errMsgPrefix(slice.Start.Pos), startType.TypeString()))
	}
	if startType != nil && startType.Type == VariableTypeLong {
		slice.Start.ConvertToNumber(VariableTypeInt)
	}
	if slice.End != nil {
		endType, es := slice.End.checkSingleValueContextExpression(block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if endType != nil && endType.IsInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s slice end must be integer,but '%s'",
				errMsgPrefix(slice.End.Pos), endType.TypeString()))
		}
		if endType != nil && endType.Type == VariableTypeLong {
			slice.End.ConvertToNumber(VariableTypeInt)
		}
	} else {
		slice.End = &Expression{}
		slice.End.Type = ExpressionTypeFunctionCall
		slice.End.Pos = e.Pos
		slice.End.Value = &Type{
			Type: VariableTypeInt,
			Pos:  e.Pos,
		}
		call := &ExpressionFunctionCall{}
		call.Function = buildInFunctionsMap[common.BuildInFunctionLen]
		call.Args = []*Expression{slice.Expression}
		slice.End.Data = call
	}

	sliceOn, es := slice.Expression.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if sliceOn == nil {
		return nil
	}
	if sliceOn.Type != VariableTypeArray && sliceOn.Type != VariableTypeString {
		*errs = append(*errs, fmt.Errorf("%s cannot have slice on '%s'",
			errMsgPrefix(slice.Expression.Pos), sliceOn.TypeString()))
	}
	result := sliceOn.Clone()
	result.Pos = e.Pos
	return result
}
