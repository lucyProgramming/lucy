package ast

import (
	"fmt"
	"math"
)

func (e *Expression) checkConstExpression(block *Block, errs *[]error) {
	cs := e.Data.(*ExpressionDeclareConsts)
	for _, v := range cs.Cs {
		is, typ, value, err := v.Expression.getConstValue()
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(v.Pos), err.Error()))
		}
		if !is {
			*errs = append(*errs, fmt.Errorf("%s const %v is not defined by const value", errMsgPrefix(v.Pos), v.Name))
		}
		if is {
			v.Expression.Typ = typ
			v.Expression.Data = value
		} else {
			v.Expression.Typ = EXPRESSION_TYPE_INT
			v.Expression.Data = math.MaxInt64
		}
		tt, _ := v.Expression.check(block)
		v.Typ = tt[0]
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
		}

	}
	return
}
