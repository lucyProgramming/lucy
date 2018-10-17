package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func (e *Expression) checkSlice(block *Block, errs *[]error) *Type {
	on := e.Data.(*ExpressionSlice)
	sliceOn, es := on.ExpressionOn.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if sliceOn == nil {
		return nil
	}
	if sliceOn.Type != VariableTypeArray &&
		sliceOn.Type != VariableTypeString {
		*errs = append(*errs, fmt.Errorf("%s cannot have slice on '%s'",
			errMsgPrefix(sliceOn.Pos), sliceOn.TypeString()))
	}
	//start
	if on.Start == nil {
		on.Start = &Expression{}
		on.Start.Pos = e.Pos
		on.Start.Description = "compilerAuto"
		on.Start.Type = ExpressionTypeInt
		on.Start.Data = int32(0)
	}
	startType, es := on.Start.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if startType != nil {
		if startType.isInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for slice",
				errMsgPrefix(startType.Pos), startType.TypeString()))
		} else {
			if startType.Type == VariableTypeLong {
				on.Start.convertToNumber(VariableTypeInt)
			}
		}
	}
	if on.End != nil {
		endType, es := on.End.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if endType != nil &&
			endType.isInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for slice",
				errMsgPrefix(endType.Pos), endType.TypeString()))
		}
		if endType != nil &&
			endType.Type == VariableTypeLong {
			on.End.convertToNumber(VariableTypeInt)
		}
	} else {
		on.End = &Expression{}
		on.End.Description = "compilerAuto"
		on.End.Type = ExpressionTypeFunctionCall
		on.End.Pos = e.Pos
		on.End.Value = &Type{
			Type: VariableTypeInt,
			Pos:  e.Pos,
		}
		call := &ExpressionFunctionCall{}
		call.Function = buildInFunctionsMap[common.BuildInFunctionLen]
		call.Args = []*Expression{on.ExpressionOn}
		on.End.Data = call
	}
	result := sliceOn.Clone()
	result.Pos = e.Pos
	return result
}
