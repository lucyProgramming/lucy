package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkVarExpression(block *Block, errs *[]error) {
	ev := e.Data.(*ExpressionVar)
	if ev.Type == nil {
		return
	}
	if err := ev.Type.resolve(block); err != nil {
		*errs = append(*errs, err)
		return
	}
	for _, v := range ev.Variables {
		v.Type = ev.Type.Clone()
	}
	noErr := true
	var err error
	if len(ev.InitValues) > 0 {
		valueTypes := checkExpressions(block, ev.InitValues, errs, false)
		{
			needs := make([]*Type, len(ev.Variables))
			for k, _ := range needs {
				needs[k] = ev.Type
			}
			convertExpressionsToNeeds(ev.InitValues, needs, valueTypes)
		}
		if len(valueTypes) != len(ev.Variables) {
			noErr = false
			*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
				errMsgPrefix(e.Pos),
				len(valueTypes),
				len(ev.Variables)))
		}
		for k, v := range ev.Variables {
			if v.Name == NoNameIdentifier {
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
				if ev.Variables[k].Type.Equal(errs, valueTypes[k]) == false {
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
		for _, v := range ev.Variables {
			if v.Name == NoNameIdentifier {
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
			ev.InitValues = append(ev.InitValues, v.Type.mkDefaultValueExpression())
			if e.IsPublic {
				v.AccessFlags |= cg.ACC_FIELD_PUBLIC
			}
		}
	}
	if noErr == false {
		return
	}
}
