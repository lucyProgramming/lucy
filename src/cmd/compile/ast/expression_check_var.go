package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkVarExpression(block *Block, errs *[]error) {
	vs := e.Data.(*ExpressionDeclareVariable)
	noErr := true
	var err error
	vs.IfDeclareBefor = make([]bool, len(vs.Vs)) // all create this time
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
				*errs = append(*errs, fmt.Errorf("%s '%s' is not a available name",
					errMsgPrefix(v.Pos), v.Name))
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
			if k < len(valueTypes) && valueTypes[k] != nil {
				if vs.Vs[k].Typ.TypeCompatible(valueTypes[k]) == false {
					err = fmt.Errorf("%s cannot assign  '%s' to '%s'",
						errMsgPrefix(valueTypes[k].Pos),
						valueTypes[k].TypeString(),
						vs.Vs[k].Typ.TypeString())
					*errs = append(*errs, err)
					noErr = false
					continue
				}
			}
			if e.IsPublic {
				v.AccessFlags |= cg.ACC_FIELD_PUBLIC
			}
		}
	} else {
		for _, v := range vs.Vs {
			if v.Name == NO_NAME_IDENTIFIER {
				*errs = append(*errs, fmt.Errorf("%s '%s' is not a available name",
					errMsgPrefix(v.Pos), v.Name))
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
			vs.Values = append(vs.Values, v.Typ.mkDefaultValueExpression())
			if e.IsPublic {
				v.AccessFlags |= cg.ACC_FIELD_PUBLIC
			}
		}
	}
	if noErr == false {
		return
	}
}
