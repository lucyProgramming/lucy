package ast

import (
	"fmt"
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
			sliceOn.Pos.ErrMsgPrefix(), sliceOn.TypeString()))
	}
	//start
	if on.Start == nil {
		on.Start = &Expression{}
		on.Start.Pos = e.Pos
		on.Start.Op = "intLiteral"
		on.Start.Type = ExpressionTypeInt
		on.Start.Data = int64(0)
	}
	startType, es := on.Start.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if startType != nil {
		if startType.isInteger() == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for startIndex",
				startType.Pos.ErrMsgPrefix(), startType.TypeString()))
		} else {
			if startType.Type == VariableTypeLong {
				on.Start.convertToNumberType(VariableTypeInt)
			}
			if on.Start.isLiteral() {
				startIndexValue := on.Start.getLongValue()
				if startIndexValue < 0 {
					*errs = append(*errs,
						fmt.Errorf("%s startIndex '%d' is negative",
							startType.Pos.ErrMsgPrefix(), startIndexValue))
				}
			}
		}
	}
	if on.End != nil {
		endType, es := on.End.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if endType != nil {
			if endType.isInteger() == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for endIndex",
					endType.Pos.ErrMsgPrefix(), endType.TypeString()))
			} else {
				if endType.Type == VariableTypeLong {
					on.End.convertToNumberType(VariableTypeInt)
				}
				if on.End.isLiteral() {
					endIndexValue := on.End.getLongValue()
					if endIndexValue < 0 {
						*errs = append(*errs,
							fmt.Errorf("%s endIndex '%d' is negative",
								endType.Pos.ErrMsgPrefix(), endIndexValue))
					}
					if startType != nil &&
						startType.isInteger() &&
						on.Start.isLiteral() {
						if on.Start.getLongValue() > endIndexValue {
							*errs = append(*errs,
								fmt.Errorf("%s startIndex '%d' is greater than endIndex '%d'",
									endType.Pos.ErrMsgPrefix(), on.Start.getLongValue(), endIndexValue))
						}
					}
				}
			}
		}
	}
	result := sliceOn.Clone()
	result.Pos = e.Pos
	return result
}
