package ast

import (
	"fmt"
)

func (e *Expression) checkColonAssignExpression(block *Block, errs *[]error) {
	binary := e.Data.(*ExpressionBinary)
	var names []*Expression
	if binary.Left.Typ == EXPRESSION_TYPE_IDENTIFIER {
		names = append(names, binary.Left)
	} else if binary.Left.Typ == EXPRESSION_TYPE_LIST {
		names = binary.Left.Data.([]*Expression)
	} else {
		*errs = append(*errs, fmt.Errorf("%s no name one the left", errMsgPrefix(e.Pos)))
	}
	values := binary.Right.Data.([]*Expression)
	ts := e.checkRightValuesValid(e.checkExpressions(block, values, errs), errs)
	if len(names) != len(ts) {
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(e.Pos),
			len(ts),
			len(names)))
	}
	var err error
	noNewVaraible := true
	for k, v := range names {
		if v.Typ != EXPRESSION_TYPE_IDENTIFIER {
			*errs = append(*errs, fmt.Errorf("%s not a name on the left", errMsgPrefix(v.Pos)))
			continue
		}
		name := v.Data.(*ExpressionIdentifer)
		if name.Name == NO_NAME_IDENTIFIER {
			continue
		}
		if variable, ok := block.Vars[name.Name]; ok {
			if k < len(ts) {
				if variable.Typ.typeCompatible(ts[k]) == false {
					*errs = append(*errs, fmt.Errorf("%s type '%s' is not compatible with '%s'",
						errMsgPrefix(ts[k].Pos),
						variable.Typ.TypeString(),
						ts[k].TypeString()))
				}
			}
		} else { // should be no error
			noNewVaraible = false
			vd := &VariableDefinition{}
			if k < len(ts) {
				vd.Typ = ts[k]
			}
			vd.Name = name.Name
			vd.Pos = v.Pos
			if k < len(ts) {
				vd.Typ = ts[k]
			}
			err = block.insert(vd.Name, v.Pos, vd)
			if err != nil {
				*errs = append(*errs, err)
			}
		}
	}
	if noNewVaraible == false {
		*errs = append(*errs, fmt.Errorf("%s no new variables to create", errMsgPrefix(e.Pos)))
	}
}

func (e *Expression) checkOpAssignExpression(block *Block, errs *[]error) (t *VariableType) {
	binary := e.Data.(*ExpressionBinary)
	t1, es := binary.Left.getLeftValue(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	ts, es := binary.Right.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t2, err := binary.Right.mustBeOneValueContext(ts)
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
	binary := e.Data.(*ExpressionBinary)
	lefts := make([]*Expression, 1)
	if binary.Left.Typ == EXPRESSION_TYPE_LIST {
		lefts = binary.Left.Data.([]*Expression)
	} else {
		lefts[0] = binary.Left
	}
	values := binary.Right.Data.([]*Expression)
	valueTypes := e.checkExpressions(block, values, errs)
	leftTypes := []*VariableType{}
	for _, v := range lefts {
		if v.Typ == EXPRESSION_TYPE_IDENTIFIER {
			name := v.Data.(*ExpressionIdentifer)
			if name.Name == NO_NAME_IDENTIFIER { // skip "_"
				lefts = append(lefts, nil)
				continue
			}
		}
		t, es := v.getLeftValue(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			continue
		}
		if t != nil {
			leftTypes = append(leftTypes, t)
		}
	}
	if len(lefts) != len(valueTypes) {
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
			errMsgPrefix(e.Pos),
			len(valueTypes),
			len(lefts)))
	}
	for k, v := range leftTypes {
		if v == nil {
			continue
		}
		if k < len(valueTypes) {
			if !leftTypes[k].typeCompatible(valueTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s",
					errMsgPrefix(e.Pos),
					leftTypes[k].TypeString(),
					valueTypes[k].TypeString()))
			}
		}
	}
	if len(leftTypes) > 1 {
		return nil
	}
	tt := leftTypes[0].Clone()
	tt.Pos = e.Pos
	e.VariableType = tt
	return tt
}
