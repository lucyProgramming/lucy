package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkVarAssignExpression(block *Block, errs *[]error) {
	bin := e.Data.(*ExpressionBinary)
	var names []*Expression
	if bin.Left.Type == ExpressionTypeList {
		names = bin.Left.Data.([]*Expression)
	} else {
		names = []*Expression{bin.Left}
	}
	noErr := true
	values := bin.Right.Data.([]*Expression)
	assignTypes := checkExpressions(block, values, errs, false)
	if len(names) > len(assignTypes) {
		pos := e.Pos
		getLastPosFromArgs(assignTypes, &pos)
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(pos),
			len(assignTypes),
			len(names)))
		noErr = false
	} else if len(names) < len(assignTypes) {
		pos := e.Pos
		getFirstPosFromArgs(assignTypes[len(names):], &pos)
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(pos),
			len(assignTypes),
			len(names)))
	}
	var err error
	noNewVariable := true
	assign := &ExpressionVarAssign{}
	assign.InitValues = values
	for k, v := range names {
		if v.Type != ExpressionTypeIdentifier {
			*errs = append(*errs, fmt.Errorf("%s not a name on the left,but '%s'",
				errMsgPrefix(v.Pos), v.OpName()))
			noErr = false
			continue
		}
		identifier := v.Data.(*ExpressionIdentifier)
		if identifier.Name == NoNameIdentifier {
			vd := &Variable{}
			vd.Name = identifier.Name
			assign.Variables = append(assign.Variables, vd)
			assign.IfDeclaredBefore = append(assign.IfDeclaredBefore, false)
			continue
		}
		var variableType *Type
		if k < len(assignTypes) {
			variableType = assignTypes[k]
		}
		if variable, ok := block.Variables[identifier.Name]; ok {
			if variableType != nil {
				if variable.Type.Equal(errs, assignTypes[k]) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(assignTypes[k].Pos),
						variable.Type.TypeString(),
						assignTypes[k].TypeString()))
					noErr = false
				}
			}
			identifier.Variable = variable
			assign.Variables = append(assign.Variables, variable)
			assign.IfDeclaredBefore = append(assign.IfDeclaredBefore, true)
		} else { // should be no error
			noNewVariable = false
			vd := &Variable{}
			if k < len(assignTypes) {
				vd.Type = assignTypes[k]
			}
			vd.Name = identifier.Name
			vd.Pos = v.Pos
			if variableType != nil {
				vd.Type = variableType.Clone()
				vd.Type.Pos = e.Pos
				if variableType.isTyped() == false {
					*errs = append(*errs, fmt.Errorf("%s '%s' not typed",
						errMsgPrefix(v.Pos), variableType.TypeString()))
				}
			} else {
				vd.Type = &Type{}
				vd.Type.Type = VariableTypeVoid
				vd.Type.Pos = v.Pos
			}
			err = block.Insert(vd.Name, v.Pos, vd)
			identifier.Variable = vd
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			assign.Variables = append(assign.Variables, vd)
			assign.IfDeclaredBefore = append(assign.IfDeclaredBefore, false)
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
