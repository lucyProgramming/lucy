package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkVarAssignExpression(block *Block, errs *[]error) {
	bin := e.Data.(*ExpressionBinary)
	var lefts []*Expression
	if bin.Left.Type == ExpressionTypeList {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts = []*Expression{bin.Left}
	}
	values := bin.Right.Data.([]*Expression)
	for _, v := range values {
		v.Lefts = lefts
	}
	assignTypes := checkExpressions(block, values, errs, false)
	if len(lefts) > len(assignTypes) {
		pos := values[len(values)-1].Pos
		*errs = append(*errs, fmt.Errorf("%s too few values , assign %d values to %d destinations",
			pos.ErrMsgPrefix(),
			len(assignTypes),
			len(lefts)))
	} else if len(lefts) < len(assignTypes) {
		pos := getExtraExpressionPos(values, len(lefts))
		*errs = append(*errs, fmt.Errorf("%s too many values , assign %d values to %d destinations",
			pos.ErrMsgPrefix(),
			len(assignTypes),
			len(lefts)))
	}
	var err error
	noNewVariable := true
	assign := &ExpressionVarAssign{}
	assign.Lefts = lefts
	assign.InitValues = values
	assign.IfDeclaredBefore = make([]bool, len(lefts))
	for k, v := range lefts {
		var variableType *Type = nil
		if k < len(assignTypes) {
			variableType = assignTypes[k]
		}
		if v.Type != ExpressionTypeIdentifier {
			t := v.getLeftValue(block, errs)
			if t == nil || variableType == nil {
				continue
			}
			if t.assignAble(errs, variableType) == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix(t.Pos), variableType.TypeString(), t.TypeString()))
			}
			continue
		}
		identifier := v.Data.(*ExpressionIdentifier)
		if identifier.Name == UnderScore {
			continue
		}
		if variable, ok := block.Variables[identifier.Name]; ok {
			if variableType != nil {
				if variable.Type.assignAble(errs, variableType) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(assignTypes[k].Pos),
						variable.Type.TypeString(),
						variableType.TypeString()))
				}
			}
			identifier.Variable = variable
			assign.IfDeclaredBefore[k] = true
		} else { // should be no error
			noNewVariable = false
			vd := &Variable{}
			if k < len(assignTypes) {
				vd.Type = assignTypes[k]
			}
			vd.Name = identifier.Name
			vd.Pos = v.Pos
			if variableType == nil {
				continue
			}
			vd.Comment = identifier.Comment
			vd.Type = variableType.Clone()
			vd.Type.Pos = e.Pos
			if err := variableType.isTyped(); err != nil {
				*errs = append(*errs, err)
			}
			if e.IsGlobal {
				err = PackageBeenCompile.Block.Insert(vd.Name, v.Pos, vd)
				vd.IsGlobal = true
			} else {
				err = block.Insert(vd.Name, v.Pos, vd)
			}
			identifier.Variable = vd
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
			if e.IsPublic { // only use when is is global
				vd.AccessFlags |= cg.ACC_FIELD_PUBLIC
			}
		}
	}
	if noNewVariable {
		*errs = append(*errs, fmt.Errorf("%s no new variables to create", errMsgPrefix(e.Pos)))
	}
	e.Data = assign
}
