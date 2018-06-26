package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkVarExpression(block *Block, errs *[]error) {
	vs := e.Data.(*ExpressionDeclareVariable)
	noErr := true
	var err error
	vs.IfDeclaredBefore = make([]bool, len(vs.Variables)) // all create this time
	if vs.InitValues != nil && len(vs.InitValues) > 0 {
		valueTypes := checkRightValuesValid(checkExpressions(block, vs.InitValues, errs), errs)
		if len(valueTypes) != len(vs.Variables) {
			noErr = false
			*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
				errMsgPrefix(e.Pos),
				len(valueTypes),
				len(vs.Variables)))
		}
		for k, v := range vs.Variables {
			if v.Name == NO_NAME_IDENTIFIER {
				*errs = append(*errs, fmt.Errorf("%s '%s' is not a available name",
					errMsgPrefix(v.Pos), v.Name))
				noErr = false
				continue
			}
			err = v.Type.resolve(block)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			err = block.Insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			if k < len(valueTypes) && valueTypes[k] != nil {
				if vs.Variables[k].Type.Equal(errs, valueTypes[k]) == false {
					err = fmt.Errorf("%s cannot assign  '%s' to '%s'",
						errMsgPrefix(valueTypes[k].Pos),
						valueTypes[k].TypeString(),
						v.Type.TypeString())
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
		for _, v := range vs.Variables {
			if v.Name == NO_NAME_IDENTIFIER {
				*errs = append(*errs, fmt.Errorf("%s '%s' is not a available name",
					errMsgPrefix(v.Pos), v.Name))
				noErr = false
				continue
			}
			err = v.Type.resolve(block)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			err := block.Insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			vs.InitValues = append(vs.InitValues, v.Type.mkDefaultValueExpression())
			if e.IsPublic {
				v.AccessFlags |= cg.ACC_FIELD_PUBLIC
			}
		}
	}
	if noErr == false {
		return
	}

	//vs.insertFunctionPointer()
}
