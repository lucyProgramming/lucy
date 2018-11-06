package ast

import (
	"fmt"
)

func (e *Expression) checkAssignExpression(block *Block, errs *[]error) *Type {
	bin := e.Data.(*ExpressionBinary)
	lefts := make([]*Expression, 1)
	if bin.Left.Type == ExpressionTypeList {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts[0] = bin.Left
	}
	values := bin.Right.Data.([]*Expression)
	for _, v := range values {
		v.Lefts = lefts
	}
	valueTypes := checkExpressions(block, values, errs, false)
	leftTypes := []*Type{}
	for _, v := range lefts {
		if v.IsIdentifier(UnderScore) {
			leftTypes = append(leftTypes, nil) // this is no assign situation
		} else {
			t := v.getLeftValue(block, errs)
			leftTypes = append(leftTypes, t) // append even if it`s nil
		}
	}
	convertExpressionsToNeeds(values, leftTypes, valueTypes)
	bin.Left.MultiValues = leftTypes
	if len(lefts) > len(valueTypes) { //expression length compare with value types is more appropriate
		pos := values[len(values)-1].Pos
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
			pos.ErrMsgPrefix(),
			len(valueTypes),
			len(lefts)))
	} else if len(lefts) < len(valueTypes) {
		pos := getExtraExpressionPos(values, len(lefts))
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
			pos.ErrMsgPrefix(),
			len(valueTypes),
			len(lefts)))
	}
	for k, v := range leftTypes {
		if v == nil { // get left value error or "_"
			continue
		}
		if k >= len(valueTypes) {
			continue
		}
		if valueTypes[k] == nil {
			continue
		}
		if false == leftTypes[k].assignAble(errs, valueTypes[k]) {
			*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
				errMsgPrefix(valueTypes[k].Pos),
				valueTypes[k].TypeString(), leftTypes[k].TypeString()))
		}
	}
	e.Data = &ExpressionAssign{
		Lefts:  lefts,
		Values: values,
	}
	voidReturn := mkVoidType(e.Pos)
	if len(lefts) > 1 {
		return voidReturn
	}
	if len(lefts) == 0 || leftTypes[0] == nil {
		return voidReturn
	}
	if e.IsStatementExpression == false {
		left := lefts[0]
		if left.Type == ExpressionTypeIdentifier {
			t := left.Data.(*ExpressionIdentifier)
			if t.Name == UnderScore {
				return voidReturn
			} else {
				if nil != t.Variable {
					t.Variable.Used = true
				}
			}
		}
	}
	// here is safe
	result := leftTypes[0].Clone()
	result.Pos = e.Pos
	return result
}
