package ast

import (
	"fmt"
)

func (e *Expression) checkConstExpression(block *Block, errs *[]error) {
	cs := e.Data.(*ExpressionDeclareConsts)
	for k, v := range cs.Consts {
		if k > len(cs.Expressions) {
			break
		}
		v.Expression = cs.Expressions[k]
		is, typ, value, err := v.Expression.getConstValue()
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(v.Pos), err.Error()))
		}
		if is == false {
			*errs = append(*errs, fmt.Errorf("%s const named '%s' is not defined by const value",
				errMsgPrefix(v.Pos), v.Name))
			continue
		}
		v.Value = value
		v.Expression.Typ = typ
		v.Expression.Data = value
		tt, _ := v.Expression.check(block)
		v.Typ = tt[0]
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
		}
	}
	return
}
