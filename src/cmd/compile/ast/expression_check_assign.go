package ast

import (
	"fmt"
)

func (e *Expression) convertColonAssignAndVar2Assign(names []*Expression, values []*Expression) {
	e.Typ = EXPRESSION_TYPE_ASSIGN
	bin := &ExpressionBinary{}
	bin.Left = &Expression{}
	bin.Left.Typ = EXPRESSION_TYPE_LIST
	bin.Left.Data = names
	bin.Right = &Expression{}
	bin.Right.Typ = EXPRESSION_TYPE_LIST
	bin.Right.Data = values
	e.Data = bin
}

//when no error,convert to assign
func (e *Expression) checkColonAssignExpression(block *Block, errs *[]error) {
	bin := e.Data.(*ExpressionBinary)
	var names []*Expression
	if bin.Left.Typ == EXPRESSION_TYPE_IDENTIFIER {
		names = append(names, bin.Left)
	} else if bin.Left.Typ == EXPRESSION_TYPE_LIST {
		names = bin.Left.Data.([]*Expression)
	} else {
		*errs = append(*errs, fmt.Errorf("%s no names on the left", errMsgPrefix(e.Pos)))
		return
	}
	noErr := true
	values := bin.Right.Data.([]*Expression)
	ts := checkRightValuesValid(checkExpressions(block, values, errs), errs)
	if len(names) != len(ts) && len(ts) != 0 {
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(e.Pos),
			len(ts),
			len(names)))
		noErr = false
	}
	var err error
	noNewVaraible := true
	for k, v := range names {
		if v.Typ != EXPRESSION_TYPE_IDENTIFIER {
			*errs = append(*errs, fmt.Errorf("%s not a name on the left,but '%s'", errMsgPrefix(v.Pos), v.OpName()))
			noErr = false
			continue
		}
		identifier := v.Data.(*ExpressionIdentifer)
		if identifier.Name == NO_NAME_IDENTIFIER {
			continue
		}
		var variableType *VariableType
		if k < len(ts) && ts[k] != nil {
			variableType = ts[k]
		}

		if variable, ok := block.Vars[identifier.Name]; ok {
			if variableType != nil {
				if variable.Typ.TypeCompatible(ts[k]) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(ts[k].Pos),
						variable.Typ.TypeString(),
						ts[k].TypeString()))
					noErr = false
				}
			}
			identifier.Var = variable
		} else { // should be no error
			noNewVaraible = false
			vd := &VariableDefinition{}
			if k < len(ts) {
				vd.Typ = ts[k]
			}
			vd.Name = identifier.Name
			vd.Pos = v.Pos
			vd.Typ = variableType
			if vd.Typ == nil { // still cannot have type,new fake int type
				vd.Typ = &VariableType{}
				vd.Typ.Typ = VARIABLE_TYPE_VOID
				vd.Typ.Pos = v.Pos
			}
			err = block.insert(vd.Name, v.Pos, vd)
			identifier.Var = vd
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
			}
		}
	}
	if noNewVaraible {
		*errs = append(*errs, fmt.Errorf("%s no new variables to create", errMsgPrefix(e.Pos)))
		noErr = false
	}
	if noErr == false {
		return
	}
	e.convertColonAssignAndVar2Assign(names, values)

}

func (e *Expression) checkOpAssignExpression(block *Block, errs *[]error) (t *VariableType) {
	bin := e.Data.(*ExpressionBinary)
	t1, es := bin.Left.getLeftValue(block)
	bin.Left.VariableType = t1
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	ts, es := bin.Right.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t2, err := bin.Right.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t1 == nil || t2 == nil {
		return
	}
	//number
	if t1.IsNumber() {
		if !t2.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on number and '%s'",
				errMsgPrefix(e.Pos),
				e.OpName(),
				t2.TypeString()))
		}
	} else if t1.Typ == VARIABLE_TYPE_STRING {
		if t2.Typ != VARIABLE_TYPE_STRING {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on string and '%s'",
				errMsgPrefix(e.Pos),
				e.OpName(),
				t2.TypeString()))
		}
	} else {
		*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'",
			errMsgPrefix(e.Pos),
			e.OpName(),
			t1.TypeString(),
			t2.TypeString()))
	}
	tt := t1.Clone()
	tt.Pos = e.Pos
	return tt
}

func (e *Expression) checkAssignExpression(block *Block, errs *[]error) *VariableType {
	bin := e.Data.(*ExpressionBinary)
	lefts := make([]*Expression, 1)
	if bin.Left.Typ == EXPRESSION_TYPE_LIST {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts[0] = bin.Left
		bin.Left = &Expression{}
		bin.Left.Typ = EXPRESSION_TYPE_LIST
		bin.Left.Data = lefts // rewrite to list anyway
	}
	values := bin.Right.Data.([]*Expression)
	valueTypes := checkExpressions(block, values, errs)
	leftTypes := []*VariableType{}
	noAssign := true
	for _, v := range lefts {
		if v.Typ == EXPRESSION_TYPE_IDENTIFIER {
			name := v.Data.(*ExpressionIdentifer)
			if name.Name == NO_NAME_IDENTIFIER { // skip "_"
				leftTypes = append(leftTypes, nil) // this is no assign situation
				continue
			}
		}
		noAssign = false
		t, es := v.getLeftValue(block)
		v.VariableType = t
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			continue
		}
		leftTypes = append(leftTypes, t) // append even if it`s nil
	}
	bin.Left.VariableTypes = leftTypes
	if len(lefts) != len(valueTypes) { //expression length compare with value types is more appropriate
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
			errMsgPrefix(e.Pos),
			len(valueTypes),
			len(lefts)))
	}
	for k, v := range leftTypes {
		if v == nil { // get left value error or "_"
			continue
		}
		if k < len(valueTypes) {
			if !leftTypes[k].TypeCompatible(valueTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
					errMsgPrefix(e.Pos),
					valueTypes[k].TypeString(), leftTypes[k].TypeString()))
			}
		}
	}
	if noAssign {
		*errs = append(*errs, fmt.Errorf("%s no assign able expression on the left", errMsgPrefix(e.Pos)))
		return nil
	}
	if len(leftTypes) == 0 {
		return nil
	}
	if len(leftTypes) > 1 {
		return mkVoidType(e.Pos)
	}
	tt := leftTypes[0].Clone()
	tt.Pos = e.Pos
	return tt
}
