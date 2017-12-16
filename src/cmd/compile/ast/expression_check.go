package ast

import (
	"fmt"
	"math"
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
		return nil, []error{fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error())}
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
		if tt == nil {
			panic(1)
		}
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
	case EXPRESSION_TYPE_CONST:
		return e.checkConstExpression(block, &errs), errs
	case EXPRESSION_TYPE_VAR:
		return e.checkVarExpression(block, &errs), errs
	case EXPRESSION_TYPE_FUNCTION_CALL:
		return e.checkFunctionCallExpression(block, &errs), errs
	case EXPRESSION_TYPE_NOT:
		tt := e.checkUnaryExpression(block, &errs)
		if tt != nil {
			return []*VariableType{tt}, errs
		} else {
			return nil, errs
		}
	case EXPRESSION_TYPE_NEGATIVE:
		tt := e.checkUnaryExpression(block, &errs)
		if tt != nil {
			return []*VariableType{tt}, errs
		} else {
			return nil, errs
		}
	default:
		panic(fmt.Sprintf("unhandled type inference:%s", e.OpName()))
	}
	return
}

func (e *Expression) checkUnaryExpression(block *Block, errs *[]error) *VariableType {
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
		if e.Typ == EXPRESSION_TYPE_NOT {
			return &VariableType{
				Typ: EXPRESSION_TYPE_BOOL,
			}
		} else {
			return &VariableType{
				Typ: EXPRESSION_TYPE_INT,
			}
		}
	}
	if e.Typ == EXPRESSION_TYPE_NOT {
		if t.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not(!) only works with bool expression", errMsgPrefix(e.Pos)))
		}
		return &VariableType{
			Typ: EXPRESSION_TYPE_BOOL,
		}
	}
	if e.Typ == EXPRESSION_TYPE_NEGATIVE {
		if !t.isNumber() {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on %s", errMsgPrefix(e.Pos), t.TypeString()))
		}
		tt := t.Clone()
		tt.Resource = nil
		return tt
	}
	panic("missing handle")
	return t
}

func (e *Expression) checkFunctionCallExpression(block *Block, errs *[]error) []*VariableType {
	call := e.Data.(*ExpressionFunctionCall)
	tt, es := call.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeValueContext(tt)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return mkVoidVariableTypes()
	}
	if t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s not a function", errMsgPrefix(call.Expression.Pos)))
		return mkVoidVariableTypes()
	}
	return e.checkFunctionCall(block, errs, t.Resource.Function, call.Args)
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, args []*Expression) []*VariableType {
	callargsTypes := []*VariableType{}
	for _, v := range args {
		tt, es := v.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if tt != nil {
			for _, vv := range tt {
				if vv.Typ == VARIABLE_TYPE_VOID {
					*errs = append(*errs, fmt.Errorf("%s function has no return value,cannot be used as right value", errMsgPrefix(v.Pos)))
					continue
				}
				callargsTypes = append(callargsTypes, vv)
			}
		}
	}
	if len(callargsTypes) > len(f.Typ.Parameters) {
		*errs = append(*errs, fmt.Errorf("%s too many paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}
	if len(callargsTypes) < len(f.Typ.Parameters) {
		*errs = append(*errs, fmt.Errorf("%s too few paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}
	for k, v := range f.Typ.Parameters {
		if k < len(callargsTypes) {
			if !v.Typ.typeCompatible(callargsTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s", errMsgPrefix(callargsTypes[k].Pos), v.Typ.TypeString(), callargsTypes[k].TypeString()))
			}
		}
	}
	ret := make([]*VariableType, len(f.Typ.Returns))
	for k, _ := range ret {
		ret[k] = f.Typ.Returns[k].Typ.Clone()
		ret[k].Resource = nil
	}
	return ret

}

func (e *Expression) checkVarExpression(block *Block, errs *[]error) []*VariableType {
	ts := []*VariableType{}
	vs := e.Data.(*ExpressionDeclareVariable)
	var err error
	for _, v := range vs.Vs {
		var t *VariableType
		if v.Expression != nil {
			tt, es := v.Expression.check(block)
			if errsNotEmpty(es) {
				*errs = append(*errs, es...)
			}
			t, err = e.mustBeValueContext(tt)
			if err != nil {
				*errs = append(*errs, err)
			}
		}
		err = v.Typ.resolve(block)
		if err != nil {
			*errs = append(*errs, err)
		} else {
			if t != nil && !v.Typ.typeCompatible(t) {

			}
		}
		v.Typ.Resource = &VariableTypeResource{}
		v.Typ.Resource.Var = v
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
		}
	}

	return ts

}
func (e *Expression) checkConstExpression(block *Block, errs *[]error) []*VariableType {
	cs := e.Data.(*ExpressionDeclareConsts)
	ts := []*VariableType{}
	for _, v := range cs.Cs {
		is, typ, value, err := v.Expression.getConstValue()
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(v.Pos), err.Error()))
		}
		if !is && err == nil {
			*errs = append(*errs, fmt.Errorf("%s const %v is not defined by const value", errMsgPrefix(v.Pos), v.Name))
		}
		if is {
			v.Expression.Typ = typ
			v.Expression.Data = value
		} else {
			v.Expression.Typ = EXPRESSION_TYPE_INT
			v.Expression.Data = math.MaxInt64
		}
		tt, _ := v.Expression.check(block)
		ts = append(ts, tt[0])
		v.Value = v.Expression.Data
		v.Typ = tt[0]
		v.Typ.Resource = &VariableTypeResource{}
		v.Typ.Resource.Const = v
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
		}
	}
	return ts
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
	valueTypes := []*VariableType{}
	values := binary.Right.Data.([]*Expression)
	for _, v := range values {
		ts, es := v.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if ts != nil {
			valueTypes = append(valueTypes, ts...)
		}
	}
	leftTypes := []*VariableType{}
	for _, v := range lefts {
		if v.Typ == EXPRESSION_TYPE_IDENTIFIER && v.Data.(string) == "_" {
			leftTypes = append(leftTypes, &VariableType{
				Typ: VARIABLE_TYPE_VOID,
			})
			continue
		}
		tt, es := v.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		t, err := e.mustBeValueContext(tt)
		if err != nil {
			*errs = append(*errs, err)
		}
		if t != nil {
			if !t.rightValueValid() {
				*errs = append(*errs, fmt.Errorf("%s cannot be used as right value", errMsgPrefix(t.Pos)))
			}
			leftTypes = append(leftTypes, t)
		}
	}
	if len(lefts) != len(valueTypes) {
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations", errMsgPrefix(e.Pos), len(valueTypes), len(lefts)))
	}
	for k, v := range lefts {
		if v.Typ == EXPRESSION_TYPE_IDENTIFIER && v.Data.(string) == "_" {
			continue
		}
		if k < len(leftTypes) {
			if err := leftTypes[k].assignAble(); err != nil {
				*errs = append(*errs, fmt.Errorf("%s cannot been used as left value", errMsgPrefix(e.Pos)))
			}
			if k < len(valueTypes) {
				if !leftTypes[k].typeCompatible(valueTypes[k]) {
					*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s", errMsgPrefix(e.Pos), leftTypes[k].TypeString(), valueTypes[k].TypeString()))
				}
			}
		}
	}
	return valueTypes
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
		f := d.(*Function)
		return f.VariableType, nil
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
		return &VariableType{
			Typ: VARIABLE_TYPE_LONG,
		}
	}
	// && AND ||
	if e.Typ == EXPRESSION_TYPE_LOGICAL_OR || EXPRESSION_TYPE_LOGICAL_AND == e.Typ {
		if t1.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(binary.Left.Pos)))
		}
		if t2.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(binary.Right.Pos)))
		}
		return &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
		}
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
		tt := t1.Clone()
		tt.Resource = nil
		return tt
	}
	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		if !t1.isNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(binary.Left.Pos)))
		}
		if !t2.isInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression", errMsgPrefix(binary.Right.Pos)))
		}
		tt := t1.Clone()
		tt.Resource = nil
		return tt
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
		tt := t1.Clone()
		tt.Resource = nil
		return tt
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
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(binary.Right.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
			}
		} else if t1.Typ == VARIABLE_TYPE_STRING {
			if e.Typ != EXPRESSION_TYPE_PLUS_ASSIGN {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm  '%s' on string", errMsgPrefix(binary.Right.Pos), e.OpName()))
			}
		} else {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
		}
		tt := t1.Clone()
		tt.Resource = nil
		return tt
	}
	panic("missing check")
	return nil
}
