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
	noErr := true
	values := bin.Right.Data.([]*Expression)
	assignTypes := checkExpressions(block, values, errs, false)
	if len(lefts) > len(assignTypes) {
		pos := values[len(values)-1].Pos
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(pos),
			len(assignTypes),
			len(lefts)))
		noErr = false
	} else if len(lefts) < len(assignTypes) {
		pos := e.Pos
		getFirstPosFromArgs(assignTypes[len(lefts):], &pos)
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(pos),
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
			if t.Equal(errs, variableType) == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix(t.Pos), variableType.TypeString(), t.TypeString()))
			}
			continue
		}
		identifier := v.Data.(*ExpressionIdentifier)
		if identifier.Name == NoNameIdentifier {
			continue
		}
		if variable, ok := block.Variables[identifier.Name]; ok {
			if variableType != nil {
				if variable.Type.Equal(errs, variableType) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(assignTypes[k].Pos),
						variable.Type.TypeString(),
						variableType.TypeString()))
					noErr = false
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
			vd.Type = variableType.Clone()
			vd.Type.Pos = e.Pos
			if variableType.isTyped() == false {
				*errs = append(*errs, fmt.Errorf("%s '%s' not typed",
					errMsgPrefix(v.Pos), variableType.TypeString()))
			}
			err = block.Insert(vd.Name, v.Pos, vd)
			identifier.Variable = vd
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			if e.IsPublic { // only use when is is global
				vd.AccessFlags |= cg.ACC_FIELD_PUBLIC
			}
		}
	}
	if noNewVariable {
		*errs = append(*errs, fmt.Errorf("%s no new variables to create", errMsgPrefix(e.Pos)))
		noErr = false
	}
	if noErr == false {
		return
	}
	// no error,rewrite data
	e.Data = assign
}
