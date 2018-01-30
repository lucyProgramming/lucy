package ast

import (
	"fmt"
)

func (e *Expression) checkVarExpression(block *Block, errs *[]error) {
	vs := e.Data.(*ExpressionDeclareVariable)
	args := e.checkExpressions(block, vs.Expressions, errs)
	args = e.checkRightValuesValid(args, errs)
	var err error
	for k, v := range vs.Vs {
		err = v.Typ.resolve(block)
		if err != nil {
			*errs = append(*errs, err)
		} else {
			if k < len(args) {
				if !v.Typ.typeCompatible(args[k]) {
					fmt.Errorf("%s cannot assign %s to %s", errMsgPrefix(args[k].Pos), args[k].TypeString(), v.Typ.TypeString())
				}
			}
		}
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
		}
	}

}
