package ast

import (
	"fmt"
)

func (e *Expression) mustBeValueContext(ts []*VariableType) (*VariableType, error) {
	if len(ts) == 0 {
		return nil, nil // no-type,no error
	}
	if len(ts) > 1 {
		return ts[0], fmt.Errorf("multi value in single value context ")
	}
	return ts[0], nil
}

func (e *Expression) check(block *Block) (t []*VariableType, errs []error) {
	is, typ, data, err := e.getConstValue()
	if err != nil {
		return nil, []error{err}
	}
	if is {
		e.Typ = typ
		e.Data = data
	}
	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		t = []*VariableType{
			&VariableType{
				Typ: VARIABLE_TYPE_BOOL,
				Pos: e.Pos,
			},
		}
	case EXPRESSION_TYPE_BYTE:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_BYTE,
			Pos: e.Pos,
		},
		}
	case EXPRESSION_TYPE_INT:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_INT,
			Pos: e.Pos,
		},
		}
	case EXPRESSION_TYPE_FLOAT:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_FLOAT,
			Pos: e.Pos,
		},
		}
	case EXPRESSION_TYPE_STRING:
		t = []*VariableType{&VariableType{
			Typ: VARIABLE_TYPE_STRING,
			Pos: e.Pos,
		}}

	case EXPRESSION_TYPE_IDENTIFIER:
		tt, err := e.checkIdentiferExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			fmt.Println("1")
			return []*VariableType{tt}, errs
		} else {
			return nil, errs
		}
		//binaries
	case EXPRESSION_TYPE_LOGICAL_OR:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LOGICAL_AND:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_OR:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_AND:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LEFT_SHIFT:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_RIGHT_SHIFT:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_PLUS_ASSIGN:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MINUS_ASSIGN:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MUL_ASSIGN:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_DIV_ASSIGN:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MOD_ASSIGN:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_EQ:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_NE:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_GE:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_GT:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LE:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_LT:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_ADD:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_SUB:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MUL:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_DIV:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_MOD:
		tt := e.checkBinaryExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_COLON_ASSIGN:
		return e.checkColonAssignExpression(block, &errs), errs
	case EXPRESSION_TYPE_ASSIGN:
		return e.checkAssignExpression(block, &errs), errs
	case EXPRESSION_TYPE_INCREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_DECREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_PRE_INCREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		return []*VariableType{tt}, errs
	case EXPRESSION_TYPE_PRE_DECREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		return []*VariableType{tt}, errs
	default:
		panic(fmt.Sprintf("unhandled type inference:%s", e.OpName()))
	}
	return
}

func (e *Expression) checkIncrementExpression(block *Block, errs *[]error) *VariableType {
	ee := e.Data.(*Expression)
	ts, es := ee.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return nil
	}
	if !t.isNumber() {
		*errs = append(*errs, fmt.Errorf("%s cannot apply ++ or -- on %s", errMsgPrefix(ee.Pos), t.TypeString()))
	}
	if err = t.assignAble(); err != nil {
		*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(ee.Pos), err.Error()))
	}
	tt := t.Clone()
	tt.Resource = nil
	return tt
}

func (e *Expression) checkAssignExpression(block *Block, errs *[]error) (ts []*VariableType) {
	binary := e.Data.(*ExpressionBinary)
	lefts := make([]*Expression, 1)
	if binary.Left.Typ == EXPRESSION_TYPE_LIST {
		lefts = binary.Left.Data.([]*Expression)
	} else {
		lefts[0] = binary.Left
	}

	return nil
}

func (e *Expression) checkColonAssignExpression(block *Block, errs *[]error) (ts []*VariableType) {
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
	ts = []*VariableType{}
	for _, v := range values {
		tt, es := v.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if tt != nil {
			ts = append(ts, tt...)
		}
	}
	if len(names) > 0 {
		if len(names) != len(ts) {
			*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations", errMsgPrefix(e.Pos), len(ts), len(names)))
		}
	}
	var err error
	for i := 0; i < len(names); i++ {
		if names[i].Typ != EXPRESSION_TYPE_IDENTIFIER {
			*errs = append(*errs, fmt.Errorf("%s not a name on the left", errMsgPrefix(names[i].Pos)))
			continue
		}
		if names[i].Data.(string) == "_" {
			continue
		}
		vd := &VariableDefinition{}
		vd.Name = names[i].Data.(string)
		vd.Pos = names[i].Pos
		if i >= 0 && i < len(ts) {
			vd.Typ = ts[i]
			ts[i].Resource = &VariableTypeResource{}
			ts[i].Resource.Var = vd
		}
		err = block.insert(vd.Name, names[i].Pos, vd)
		if err != nil {
			*errs = append(*errs, err)
		}
	}

	return nil
}

func (e *Expression) checkIdentiferExpression(block *Block) (t *VariableType, err error) {
	name := e.Data.(string)
	d, err := block.searchByName(name)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err)
	}
	switch d.(type) {
	case *Function:
		//		f := d.(*Function)
		return nil, nil
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		return t.Typ, nil
	case *Const:
		t := d.(*Const)
		return t.Typ, nil
	case *Enum:
		t := d.(*Enum)
		return t.VariableType, nil
	default:
		panic(1111111)
	}
	return nil, nil
}

func (e *Expression) checkBinaryExpression(block *Block, errs *[]error) (t *VariableType) {
	binary := e.Data.(*ExpressionBinary)
	ts1, err1 := binary.Left.check(block)
	ts2, err2 := binary.Right.check(block)
	if errsNotEmpty(err1) {
		*errs = append(*errs, err1...)
	}
	if errsNotEmpty(err2) {
		*errs = append(*errs, err2...)
	}
	var err error
	t1, err := e.mustBeValueContext(ts1)
	if err != nil {
		*errs = append(*errs, err)
	}
	t2, err := e.mustBeValueContext(ts2)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t1 == nil || t2 == nil {
		return
	}
	// && AND ||
	if e.Typ == EXPRESSION_TYPE_LOGICAL_OR || EXPRESSION_TYPE_LOGICAL_AND == e.Typ {
		if t1.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(binary.Left.Pos)))
		}
		if t2.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(binary.Right.Pos)))
		}
		return t1
	}
	if e.Typ == EXPRESSION_TYPE_OR || EXPRESSION_TYPE_AND == e.Typ {
		if !t1.isNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(binary.Left.Pos)))
		}
		if !t2.isNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(binary.Right.Pos)))
		}
		if t1.isNumber() && t2.isNumber() {
			if t1.Typ != t2.Typ {
				*errs = append(*errs, fmt.Errorf("%s %s does not match %s", errMsgPrefix(e.Pos), t1.TypeString(), t2.TypeString()))
			}
		}
		return t1
	}
	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		if !t1.isNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(binary.Left.Pos)))
		}
		if !t2.isInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression", errMsgPrefix(binary.Right.Pos)))
		}
		return t1
	}
	if e.Typ == EXPRESSION_TYPE_PLUS_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_MINUS_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_MUL_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_DIV_ASSIGN ||
		e.Typ == EXPRESSION_TYPE_MOD_ASSIGN {
		//cannot be assign
		if err := t1.assignAble(); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		}

		if t1.isNumber() {
			if !t2.isNumber() {
				*errs = append(*errs, fmt.Errorf("%s not a number on the right of the equation", errMsgPrefix(e.Pos)))
			}
		}
		if t1.Typ == VARIABLE_TYPE_STRING {
			if e.Typ != EXPRESSION_TYPE_PLUS_ASSIGN {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm  '%s' on string", errMsgPrefix(e.Pos), e.OpName()))
			}
		}
		return t1
	}
	if e.Typ == EXPRESSION_TYPE_EQ ||
		e.Typ == EXPRESSION_TYPE_NE ||
		e.Typ == EXPRESSION_TYPE_GE ||
		e.Typ == EXPRESSION_TYPE_GT ||
		e.Typ == EXPRESSION_TYPE_LE ||
		e.Typ == EXPRESSION_TYPE_LT {
		//number
		if t1.isNumber() {
			if !t2.isNumber() {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on number and '%s'", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
			}
		} else if t1.Typ == VARIABLE_TYPE_STRING {
			if t2.Typ != VARIABLE_TYPE_STRING {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on string and '%s'", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
			}
		} else {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
		}
		return &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
		}
	}
	if e.Typ == EXPRESSION_TYPE_ADD ||
		e.Typ == EXPRESSION_TYPE_SUB ||
		e.Typ == EXPRESSION_TYPE_MUL ||
		e.Typ == EXPRESSION_TYPE_DIV ||
		e.Typ == EXPRESSION_TYPE_MOD {
		if t1.isNumber() {
			if !t2.isNumber() {
				*errs = append(*errs, fmt.Errorf("%s not a number on the right of the equation", errMsgPrefix(e.Pos)))
			}
		} else if t1.Typ == VARIABLE_TYPE_STRING {
			if e.Typ != EXPRESSION_TYPE_PLUS_ASSIGN {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm  '%s' on string", errMsgPrefix(e.Pos), e.OpName()))
			}
		} else {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
		}
		return t1
	}
	panic("missing check")
	return nil
}
