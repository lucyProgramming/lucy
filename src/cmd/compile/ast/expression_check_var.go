package ast

import (
	"fmt"
)

//	"fmt"

func (e *Expression) checkVarExpression(block *Block, errs *[]error) {
	vs := e.Data.(*ExpressionDeclareVariable)
	noErr := true
	names := []*Expression{}
	values := []*Expression{}
	var err error
	if vs.Values != nil && len(vs.Values) > 0 {
		valueTypes := checkRightValuesValid(checkExpressions(block, vs.Values, errs), errs)
		if len(valueTypes) != len(vs.Vs) {
			noErr = false
			*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
				errMsgPrefix(e.Pos),
				len(valueTypes),
				len(vs.Vs)))
		}
		for k, v := range vs.Vs {
			if v.Name == NO_NAME_IDENTIFIER {
				*errs = append(*errs, fmt.Errorf("%s '%s' is not a available name", errMsgPrefix(v.Pos), v.Name))
				noErr = false
				continue
			}
			err = v.Typ.resolve(block)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			err = block.insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			if k < len(valueTypes) {
				if valueTypes[k].TypeCompatible(vs.Vs[k].Typ) == false {
					err = fmt.Errorf("%s cannot assign  '%s' to '%s'",
						errMsgPrefix(valueTypes[k].Pos),
						valueTypes[k].TypeString(),
						vs.Vs[k].Typ.TypeString())
					*errs = append(*errs, err)
					noErr = false
					continue
				}
			}
			var nameExpression Expression
			nameExpression.Typ = EXPRESSION_TYPE_IDENTIFIER
			nameExpression.Pos = v.Pos
			identifier := &ExpressionIdentifer{}
			identifier.Name = v.Name
			identifier.Var = v
			nameExpression.Data = identifier
			names = append(names, &nameExpression)
		}
		values = vs.Values
	} else {
		for _, v := range vs.Vs {
			if v.Name == NO_NAME_IDENTIFIER {
				*errs = append(*errs, fmt.Errorf("%s '%s' is not a available name", errMsgPrefix(v.Pos), v.Name))
				noErr = false
				continue
			}
			err = v.Typ.resolve(block)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			err := block.insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			values = append(values, v.Typ.mkDefaultValueExpression())
			var nameExpression Expression
			nameExpression.Typ = EXPRESSION_TYPE_IDENTIFIER
			nameExpression.Pos = v.Pos
			identifier := &ExpressionIdentifer{}
			identifier.Name = v.Name
			identifier.Var = v
			nameExpression.Data = identifier
			names = append(names, &nameExpression)
		}
	}
	if noErr == false {
		return
	}
	e.convertColonAssignAndVar2Assign(names, values)
}
