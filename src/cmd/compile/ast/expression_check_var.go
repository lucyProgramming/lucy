package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkVarExpression(block *Block, errs *[]error) {
	ev := e.Data.(*ExpressionVar)
	if ev.Type != nil {
		if err := ev.Type.resolve(block); err != nil {
			*errs = append(*errs, err)
			return
		}
		if len(ev.Variables) == 1 {
			ev.Variables[0].Type = ev.Type
		} else {
			for _, v := range ev.Variables {
				v.Type = ev.Type.Clone()
			}
		}
	}
	if ev.Type == nil && len(ev.InitValues) == 0 {
		// looks impossible
		*errs = append(*errs,
			fmt.Errorf("%s expression var have not type and no initValues",
				errMsgPrefix(e.Pos)))
		return
	}
	var err error
	if len(ev.InitValues) > 0 {
		valueTypes := checkExpressions(block, ev.InitValues, errs, false)
		if ev.Type != nil {
			needs := make([]*Type, len(ev.Variables))
			for k, _ := range needs {
				needs[k] = ev.Type
			}
			convertExpressionsToNeeds(ev.InitValues, needs, valueTypes)
		}
		if len(ev.Variables) > len(valueTypes) {
			pos := ev.InitValues[len(ev.InitValues)-1].Pos // last Pos
			*errs = append(*errs, fmt.Errorf("%s too few values , assign %d values to %d destinations",
				pos.ErrMsgPrefix(),
				len(valueTypes),
				len(ev.Variables)))
		} else if len(ev.Variables) < len(valueTypes) {
			pos := getExtraExpressionPos(ev.InitValues, len(ev.Variables))
			*errs = append(*errs, fmt.Errorf("%s too many values , assign %d values to %d destinations",
				pos.ErrMsgPrefix(),
				len(valueTypes),
				len(ev.Variables)))
		}
		for k, v := range ev.Variables {
			if k < len(valueTypes) && valueTypes[k] != nil {
				if v.Type != nil {
					if v.Type.assignAble(errs, valueTypes[k]) == false {
						err = fmt.Errorf("%s cannot assign  '%s' to '%s'",
							errMsgPrefix(valueTypes[k].Pos),
							valueTypes[k].TypeString(),
							v.Type.TypeString())
						*errs = append(*errs, err)
						continue
					}
				} else {
					v.Type = valueTypes[k].Clone()
					v.Type.Pos = v.Pos
				}
			}
			if v.Type == nil {
				continue
			}
			if e.IsGlobal {
				err = PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
			} else {
				err = block.Insert(v.Name, v.Pos, v)
			}
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
		}
	} else {
		for _, v := range ev.Variables {
			var err error
			if e.IsGlobal {
				err = PackageBeenCompile.Block.Insert(v.Name, v.Pos, v)
			} else {
				err = block.Insert(v.Name, v.Pos, v)
			}
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
			ev.InitValues = append(ev.InitValues,
				v.Type.mkDefaultValueExpression())
		}
	}
	if e.IsGlobal {
		for _, v := range ev.Variables {
			v.IsGlobal = true
		}
		if e.IsPublic {
			for _, v := range ev.Variables {
				v.AccessFlags |= cg.AccFieldPublic
			}
		}
	}

}
