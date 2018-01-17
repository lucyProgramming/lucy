package ast

import (
	"fmt"
	"math"
)

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
	case EXPRESSION_TYPE_NULL:
		t = []*VariableType{
			{
				Typ: VARIABLE_TYPE_NULL,
				Pos: e.Pos,
			},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_BOOL:
		t = []*VariableType{
			{
				Typ: VARIABLE_TYPE_BOOL,
				Pos: e.Pos,
			},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_BYTE:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_BYTE,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_INT:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_INT,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_FLOAT:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_DOUBLE,
			Pos: e.Pos,
		},
		}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_STRING:
		t = []*VariableType{{
			Typ: VARIABLE_TYPE_STRING,
			Pos: e.Pos,
		}}
		e.VariableType = t[0]
	case EXPRESSION_TYPE_IDENTIFIER:
		tt, err := e.checkIdentiferExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			e.VariableType = tt
			t = []*VariableType{tt}
		}
		//binaries
	case EXPRESSION_TYPE_LOGICAL_OR:
		fallthrough
	case EXPRESSION_TYPE_LOGICAL_AND:
		fallthrough
	case EXPRESSION_TYPE_OR:
		fallthrough
	case EXPRESSION_TYPE_AND:
		fallthrough
	case EXPRESSION_TYPE_LEFT_SHIFT:
		fallthrough
	case EXPRESSION_TYPE_RIGHT_SHIFT:
		fallthrough
	case EXPRESSION_TYPE_EQ:
		fallthrough
	case EXPRESSION_TYPE_NE:
		fallthrough
	case EXPRESSION_TYPE_GE:
		fallthrough
	case EXPRESSION_TYPE_GT:
		fallthrough
	case EXPRESSION_TYPE_LE:
		fallthrough
	case EXPRESSION_TYPE_LT:
		fallthrough
	case EXPRESSION_TYPE_ADD:
		fallthrough
	case EXPRESSION_TYPE_SUB:
		fallthrough
	case EXPRESSION_TYPE_MUL:
		fallthrough
	case EXPRESSION_TYPE_DIV:
		fallthrough
	case EXPRESSION_TYPE_MOD:
		tt := e.checkBinaryExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_COLON_ASSIGN:
		e.checkColonAssignExpression(block, &errs)
	case EXPRESSION_TYPE_ASSIGN:
		tt := e.checkAssignExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_DECREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_INCREMENT:
		fallthrough
	case EXPRESSION_TYPE_PRE_DECREMENT:
		tt := e.checkIncrementExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_CONST:
		e.checkConstExpression(block, &errs)
	case EXPRESSION_TYPE_VAR:
		e.checkVarExpression(block, &errs)
	case EXPRESSION_TYPE_FUNCTION_CALL:
		t = e.checkFunctionCallExpression(block, &errs)
		e.VariableTypes = t
	case EXPRESSION_TYPE_METHOD_CALL:
		t = e.checkMethodCallExpression(block, &errs)
		e.VariableTypes = t
	case EXPRESSION_TYPE_NOT:
		fallthrough
	case EXPRESSION_TYPE_NEGATIVE:
		tt := e.checkUnaryExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_INDEX:
		fallthrough
	case EXPRESSION_TYPE_DOT:
		tt := e.checkIndexExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_CONVERTION_TYPE:
		tt := e.checkTypeConvertionExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_NEW:
		tt := e.checkNewExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_PLUS_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_MINUS_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_MUL_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_DIV_ASSIGN:
		fallthrough
	case EXPRESSION_TYPE_MOD_ASSIGN:
		tt := e.checkOpAssignExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	case EXPRESSION_TYPE_FUNCTION:
		tt := e.checkFunctionExpression(block, &errs)
		if tt != nil {
			t = []*VariableType{tt}
		}
	default:
		panic(fmt.Sprintf("unhandled type inference:%s", e.OpName()))
	}
	return
}
func (e *Expression) checkFunctionExpression(block *Block, errs *[]error) *VariableType {
	f := e.Data.(*Function)
	*errs = append(*errs, f.check(block)...)
	return &VariableType{
		Pos:      e.Pos,
		Typ:      VARIABLE_TYPE_FUNCTION,
		Function: f,
	}
}

func (e *Expression) checkExpressions(block *Block, es []*Expression, errs *[]error) []*VariableType {
	ret := []*VariableType{}
	for _, v := range es {
		ts, e := v.check(block)
		if errsNotEmpty(e) {
			*errs = append(*errs, e...)
		}
		if ts != nil {
			ret = append(ret, ts...)
		}
	}
	return ret
}

func (e *Expression) mustBeOneValueContext(ts []*VariableType) (*VariableType, error) {
	if len(ts) == 0 {
		return nil, nil // no-type,no error
	}
	if len(ts) > 1 {
		return ts[0], fmt.Errorf("%s multi value in single value context", errMsgPrefix(e.Pos))
	}
	return ts[0], nil
}

func (e *Expression) checkNewExpression(block *Block, errs *[]error) *VariableType {
	no := e.Data.(*ExpressionNew)
	err := no.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		return nil
	}
	if no.Typ.Typ != VARIABLE_TYPE_CLASS {
		*errs = append(*errs, fmt.Errorf("%s only class type can be used by new", errMsgPrefix(e.Pos)))
		return nil
	}
	args := e.checkExpressions(block, no.Args, errs)

	f, accessable, err := no.Typ.Class.matchContructionFunction(args)
	if err != nil {
		*errs = append(*errs, err)
	} else {
		if !accessable {
			*errs = append(*errs, fmt.Errorf("%s construction method is private", errMsgPrefix(e.Pos)))
		}
	}
	no.Construction = f
	ret := &VariableType{}
	*ret = *no.Typ
	ret.Typ = VARIABLE_TYPE_OBJECT
	return ret
}
func (e *Expression) checkTypeConvertionExpression(block *Block, errs *[]error) *VariableType {
	c := e.Data.(*ExpressionTypeConvertion)
	ts, es := c.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return nil
	}

	return nil
}

func (e *Expression) checkUnaryExpression(block *Block, errs *[]error) *VariableType {
	ee := e.Data.(*Expression)
	ts, es := ee.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
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
		t := &VariableType{
			Typ: EXPRESSION_TYPE_BOOL,
			Pos: e.Pos,
		}
		e.VariableType = t
		return t
	}
	if e.Typ == EXPRESSION_TYPE_NEGATIVE {
		if !t.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '-' on '%s'", errMsgPrefix(e.Pos), t.TypeString()))
		}
		tt := t.Clone()
		tt.Pos = e.Pos
		e.VariableType = tt
		return tt
	}
	panic("missing handle")
	return t
}
func (e *Expression) checkMethodCallExpression(block *Block, errs *[]error) []*VariableType {
	call := e.Data.(*ExpressionMethodCall)
	ts, es := call.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
		switch call.Name {
		case "size":
			t = &VariableType{}
			t.Typ = VARIABLE_TYPE_INT
			t.Pos = e.Pos
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s too mamy argument to call 'size'", errMsgPrefix(e.Pos)))
			}
			return []*VariableType{t}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call %s on array", errMsgPrefix(e.Pos), call.Name))
		}
		return nil
	}
	if t.Typ != VARIABLE_TYPE_OBJECT {
		*errs = append(*errs, fmt.Errorf("%s cannot make method call on a none object", errMsgPrefix(e.Pos)))
	}
	args := e.checkExpressions(block, call.Args, errs)
	args = e.checkRightValues(args, errs)
	f, es := t.Class.accessMethod(call.Name, e.Pos, args)
	if errsNotEmpty(es) {
		*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err))
	} else {
		if !call.Expression.isThisIdentifierExpression() {
			*errs = append(*errs, fmt.Errorf("%s method  %s is not public", errMsgPrefix(e.Pos), call.Name))
		}
	}
	if f == nil {
		return nil
	}
	return args
}

func (e *Expression) checkRightValues(ts []*VariableType, errs *[]error) (ret []*VariableType) {
	ret = []*VariableType{}
	for _, v := range ts {
		if !v.rightValueValid() {
			*errs = append(*errs, fmt.Errorf("%s %s cannot used as right value", errMsgPrefix(v.Pos), v.TypeString()))
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

func (e *Expression) checkFunctionCallExpression(block *Block, errs *[]error) []*VariableType {
	call := e.Data.(*ExpressionFunctionCall)
	tt, es := call.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(tt)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return mkVoidVariableTypes(e.Pos)
	}
	if t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s not a function", errMsgPrefix(call.Expression.Pos)))
		return mkVoidVariableTypes(e.Pos)
	}
	call.Func = t.Function
	if t.Function.Isbuildin {
		return e.checkBuildinFunctionCall(block, errs, t.Function, call.Args)
	} else {
		return e.checkFunctionCall(block, errs, t.Function, call.Args)
	}
}

func (e *Expression) checkBuildinFunctionCall(block *Block, errs *[]error, f *Function, args []*Expression) []*VariableType {
	callargsTypes := e.checkRightValues(e.checkExpressions(block, args, errs), errs)
	if f.CallChecker != nil {
		f.CallChecker(errs, callargsTypes, e.Pos)
	} else {
		var t *VariableType
		if f.IsAnyNumberParameter {
			if f.Typ.Parameters != nil && len(f.Typ.Parameters) > 0 && f.Typ.Parameters[0] != nil {
				t = f.Typ.Parameters[0].Typ
			}
		}

		if !f.IsAnyNumberParameter {
			if len(callargsTypes) > len(f.Typ.Parameters) {
				*errs = append(*errs, fmt.Errorf("%s too many paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
			}
			if len(callargsTypes) < len(f.Typ.Parameters) && len(args) < len(f.Typ.Parameters) {
				*errs = append(*errs, fmt.Errorf("%s too few paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
			}
		}
		if f.IsAnyNumberParameter {
			for k, v := range callargsTypes {
				if t != nil {
					if !t.typeCompatible(v) {
						*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s", errMsgPrefix(args[k].Pos), v.TypeString(), t.TypeString()))
					}
				}
			}
		} else {
			for k, v := range f.Typ.Parameters {
				if k < len(callargsTypes) {
					if !v.Typ.typeCompatible(callargsTypes[k]) {
						*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s", errMsgPrefix(args[k].Pos), v.Typ.TypeString(), callargsTypes[k].TypeString()))
					}
				}

			}
		}
	}
	return f.Typ.Returns.retTypes(e.Pos)
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, args []*Expression) []*VariableType {
	callargsTypes := e.checkExpressions(block, args, errs)
	callargsTypes = e.checkRightValues(callargsTypes, errs)
	if len(callargsTypes) > len(f.Typ.Parameters) {
		*errs = append(*errs, fmt.Errorf("%s too many paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}
	if len(callargsTypes) < len(f.Typ.Parameters) && len(args) < len(f.Typ.Parameters) {
		*errs = append(*errs, fmt.Errorf("%s too few paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}

	for k, v := range f.Typ.Parameters {
		if k < len(callargsTypes) {
			if !v.Typ.typeCompatible(callargsTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s", errMsgPrefix(args[k].Pos), v.Typ.TypeString(), callargsTypes[k].TypeString()))
			}
		}
	}
	return f.Typ.Returns.retTypes(e.Pos)
}

func (e *Expression) checkVarExpression(block *Block, errs *[]error) {
	vs := e.Data.(*ExpressionDeclareVariable)
	args := e.checkExpressions(block, vs.Expressions, errs)
	args = e.checkRightValues(args, errs)
	var err error
	for k, v := range vs.Vs {
		err = v.Typ.resolve(block)
		if err != nil {
			*errs = append(*errs, err)
		} else {
			if k < len(args) {
				if !v.Typ.typeCompatible(args[k]) {
					fmt.Errorf("%s cannot assign %s to %s", errMsgPrefix(args[k].Pos), args[k].TypeString(), v.Typ.TypeString())
				}
			}
		}
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
		}
	}

}
func (e *Expression) checkConstExpression(block *Block, errs *[]error) {
	cs := e.Data.(*ExpressionDeclareConsts)
	for _, v := range cs.Cs {
		is, typ, value, err := v.Expression.getConstValue()
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(v.Pos), err.Error()))
		}
		if !is {
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
		v.Value = v.Expression.Data
		v.Typ = tt[0]
		err = block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
		}
	}
	return
}

func (e *Expression) checkIncrementExpression(block *Block, errs *[]error) *VariableType {
	ee := e.Data.(*Expression)
	t, es := ee.getLeftValue(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	if !t.IsNumber() {
		*errs = append(*errs, fmt.Errorf("%s cannot apply ++ or -- on %s", errMsgPrefix(ee.Pos), t.TypeString()))
	}
	tt := t.Clone()
	tt.Pos = e.Pos
	e.VariableType = tt
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
			if name.Name == "_" { // skip "_"
				lefts = append(lefts, nil)
				continue
			}
		}
		t, es := v.getLeftValue(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if t != nil {
			leftTypes = append(leftTypes, t)
		}
	}
	if len(lefts) != len(valueTypes) {
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations", errMsgPrefix(e.Pos), len(valueTypes), len(lefts)))
	}
	for k, v := range leftTypes {
		if v == nil {
			continue
		}

		if k < len(valueTypes) {
			if !leftTypes[k].typeCompatible(valueTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s", errMsgPrefix(e.Pos), leftTypes[k].TypeString(), valueTypes[k].TypeString()))
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
	ts := e.checkExpressions(block, values, errs)
	if len(names) != len(ts) {
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations", errMsgPrefix(e.Pos), len(ts), len(names)))
	}
	var err error
	var noNewVaraible bool
	for k, v := range names {
		if v.Typ != EXPRESSION_TYPE_IDENTIFIER {
			*errs = append(*errs, fmt.Errorf("%s not a name on the left", errMsgPrefix(v.Pos)))
			continue
		}
		name := v.Data.(*ExpressionIdentifer)
		if name.Name == "_" {

			continue
		}
		if variable, ok := block.Vars[name.Name]; ok {
			if k < len(ts) {
				if !variable.Typ.typeCompatible(ts[k]) {
					*errs = append(*errs, fmt.Errorf("%s type '%s' is not compatible with '%s'",
						errMsgPrefix(ts[k].Pos),
						variable.Typ.TypeString(),
						ts[k].TypeString()))
				}
			}
		} else { // should be no error
			noNewVaraible = true
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
	if !noNewVaraible {
		*errs = append(*errs, fmt.Errorf("%s no new variables to create", errMsgPrefix(e.Pos)))
	}

}

func (e *Expression) checkIdentiferExpression(block *Block) (t *VariableType, err error) {
	identifer := e.Data.(*ExpressionIdentifer)
	d := block.searchByName(identifer.Name)
	if d == nil { // search failed
		if block.InheritedAttribute.class != nil { // in class
			f, err := block.InheritedAttribute.class.accessField(identifer.Name)
			if err != nil {
				return nil, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err)
			} else {
				if f.isStatic() {

				} else {

				}
				e.Typ = EXPRESSION_TYPE_DOT
				t := &ExpressionIndex{}
				//				t.Expression =
				e.Data = t
			}
		}
	}
	if d == nil {
		return nil, fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), identifer.Name)
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		f.Used = true
		t = &f.VariableType
		tt := t.Clone()
		tt.Pos = e.Pos
		identifer.Func = f
		return tt, nil
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		t.Used = true
		tt := t.Typ.Clone()
		tt.Pos = e.Pos
		identifer.Var = t
		return tt, nil
	case *Const:
		t := d.(*Const)
		t.Used = true
		tt := t.Typ.Clone()
		tt.Pos = e.Pos
		identifer.Const = t
		return tt, nil
	case *Enum:
		t := d.(*Enum)
		t.Used = true
		tt := t.VariableType.Clone()
		tt.Pos = e.Pos
		identifer.Enum = t
		return tt, nil
	case *EnumName:
		t := d.(*EnumName)
		t.Enum.Used = true
		tt := t.Enum.VariableType.Clone()
		tt.Pos = e.Pos
		identifer.EnumName = t
		return tt, nil
	default:
		panic(1111111)
	}
	return nil, nil
}
func (e *Expression) getLeftValue(block *Block) (t *VariableType, errs []error) {
	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_IDENTIFIER:
		name := e.Data.(*ExpressionIdentifer)
		d := block.searchByName(name.Name)
		if d == nil {
			return nil, []error{fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), name.Name)}
		}
		switch d.(type) {
		case *VariableDefinition:
			return d.(*VariableDefinition).Typ, nil
		default:
			errs = append(errs, fmt.Errorf("%s identifier %s is not variable", errMsgPrefix(e.Pos), name.Name))
			return nil, []error{}
		}
	case EXPRESSION_TYPE_INDEX:
		return e.checkIndexExpression(block, &errs), errs
	case EXPRESSION_TYPE_DOT:
		return e.checkIndexExpression(block, &errs), errs
	default:
		errs = append(errs, fmt.Errorf("%s %s cannot be used as left value", errMsgPrefix(e.Pos), e.OpName()))
		return nil, errs
	}
}

func (e *Expression) isThisIdentifierExpression() (b bool) {
	if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
		return
	}
	t := e.Data.(*ExpressionIdentifer)
	b = (t.Name == THIS)
	return
}

func (e *Expression) checkIndexExpression(block *Block, errs *[]error) (t *VariableType) {
	index := e.Data.(*ExpressionIndex)
	f := func() *VariableType {
		ts, es := index.Expression.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		t, err := e.mustBeOneValueContext(ts)
		if err != nil {
			*errs = append(*errs, err)
		}
		if t == nil {
			return nil
		}
		if t.Typ != VARIABLE_TYPE_ARRAY_INSTANCE && VARIABLE_TYPE_OBJECT != t.Typ {
			op := "access"
			if e.Typ == EXPRESSION_TYPE_INDEX {
				op = "index"
			}
			*errs = append(*errs, fmt.Errorf("%s cannot %s on %s", errMsgPrefix(e.Pos), op, t.TypeString()))
			return nil
		}
		return t
	}
	obj := f()
	if obj == nil {
		return nil
	}
	if obj.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
		ts, es := index.Index.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		t, err := e.mustBeOneValueContext(ts)
		if err != nil {
			*errs = append(*errs, err)
		}
		if t != nil {
			if !t.IsInteger() {
				*errs = append(*errs, fmt.Errorf("%s only integer can be used as index", errMsgPrefix(e.Pos)))
			}
		}
		return obj.CombinationType
	}
	if obj.Typ == VARIABLE_TYPE_OBJECT {
		if e.Typ != EXPRESSION_TYPE_DOT {
			*errs = append(*errs, fmt.Errorf("%s object`s field can only access by '.'", errMsgPrefix(e.Pos)))
			return nil
		}
		f, err := obj.Class.accessField(index.Name)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		} else {
			if !index.Expression.isThisIdentifierExpression() && !f.isPublic() {
				*errs = append(*errs, fmt.Errorf("%s field %s is private", errMsgPrefix(e.Pos), index.Name))
			}
		}
		if f != nil {
			return f.Typ
		} else {
			return nil
		}
	}
	panic("111")
	return nil
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
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on number and '%s'", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
		}
	} else if t1.Typ == VARIABLE_TYPE_STRING {
		if t2.Typ != VARIABLE_TYPE_STRING {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on string and '%s'", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
		}
	} else {
		*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
	}
	tt := t1.Clone()
	tt.Pos = e.Pos
	return tt
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
	t1, err := e.mustBeOneValueContext(ts1)
	if err != nil {
		*errs = append(*errs, err)
	}
	t2, err := e.mustBeOneValueContext(ts2)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t1 == nil || t2 == nil {
		return &VariableType{
			Typ: VARIABLE_TYPE_LONG,
		}
	}
	// &&  ||
	if e.Typ == EXPRESSION_TYPE_LOGICAL_OR || EXPRESSION_TYPE_LOGICAL_AND == e.Typ {
		if t1.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression,but '%s'", errMsgPrefix(binary.Left.Pos), t1.TypeString()))
		}
		if t2.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression,but '%s'", errMsgPrefix(binary.Right.Pos), t2.TypeString()))
		}
		t = &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
			Pos: e.Pos,
		}
		e.VariableType = t
		return t
	}
	// & |
	if e.Typ == EXPRESSION_TYPE_OR || EXPRESSION_TYPE_AND == e.Typ {
		if !t1.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(binary.Left.Pos)))
		}
		if !t2.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(binary.Right.Pos)))
		}
		if t1.Typ != t2.Typ {
			*errs = append(*errs, fmt.Errorf("%s cannot apply '&' or '|' on '%s' and '%s'", errMsgPrefix(binary.Right.Pos), t1.TypeString(), t2.TypeString()))
		}
		tt := t1.Clone()
		tt.Pos = e.Pos
		e.VariableType = tt
		return tt
	}

	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		if !t1.IsInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression,but '%s'", errMsgPrefix(binary.Left.Pos), t1.TypeString()))
		}
		if !t2.IsInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression,but '%s'", errMsgPrefix(binary.Right.Pos), t2.TypeString()))
		}
		tt := t1.Clone()
		tt.Pos = e.Pos
		e.VariableType = tt
		return tt
	}
	if e.Typ == EXPRESSION_TYPE_EQ ||
		e.Typ == EXPRESSION_TYPE_NE ||
		e.Typ == EXPRESSION_TYPE_GE ||
		e.Typ == EXPRESSION_TYPE_GT ||
		e.Typ == EXPRESSION_TYPE_LE ||
		e.Typ == EXPRESSION_TYPE_LT {
		//number
		switch t1.Typ {
		case VARIABLE_TYPE_BYTE:
			fallthrough
		case VARIABLE_TYPE_SHORT:
			fallthrough
		case VARIABLE_TYPE_CHAR:
			fallthrough
		case VARIABLE_TYPE_INT:
			fallthrough
		case VARIABLE_TYPE_LONG:
			if !t2.IsNumber() {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'number' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
			}
		case VARIABLE_TYPE_STRING:
			if t2.Typ != VARIABLE_TYPE_STRING {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'string' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
			}
		case VARIABLE_TYPE_BOOL:
			if t2.Typ == VARIABLE_TYPE_BOOL {
				if e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE {
					*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'bool' and 'bool'", errMsgPrefix(e.Pos), e.OpName()))
				}
			} else {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'bool' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
			}
		case VARIABLE_TYPE_NULL:
			if t2.IsPointer() {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and '%s'(non-pointer)", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
			}
			if e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and 'pointer' ", errMsgPrefix(e.Pos), e.OpName()))
			}
		case VARIABLE_TYPE_ARRAY_INSTANCE:
			fallthrough
		case VARIABLE_TYPE_OBJECT:
			if t2.IsPointer() == false && t2.Typ != VARIABLE_TYPE_NULL {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'pointer' and '%s'(non-pointer)", errMsgPrefix(e.Pos), e.OpName(), t2.TypeString()))
			}
			if e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and 'pointer' ", errMsgPrefix(e.Pos), e.OpName()))
			}
		default:
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
		}
		t := &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
			Pos: e.Pos,
		}
		e.VariableType = t
		return t
	}
	if e.Typ == EXPRESSION_TYPE_ADD ||
		e.Typ == EXPRESSION_TYPE_SUB ||
		e.Typ == EXPRESSION_TYPE_MUL ||
		e.Typ == EXPRESSION_TYPE_DIV ||
		e.Typ == EXPRESSION_TYPE_MOD {
		if t1.IsNumber() {
			if !t2.IsNumber() {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(binary.Right.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
			}
		} else if t1.Typ == VARIABLE_TYPE_STRING {
			if e.Typ != EXPRESSION_TYPE_ADD || t2.Typ != VARIABLE_TYPE_STRING {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm  '%s' on 'string' and ", errMsgPrefix(binary.Right.Pos), e.OpName(), t2.TypeString()))
			}
			tt := t1.Clone()
			tt.Pos = e.Pos
			return tt
		} else {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(e.Pos), e.OpName(), t1.TypeString(), t2.TypeString()))
		}
		tt := &VariableType{}
		tt.Pos = e.Pos
		tt.Typ = t1.NumberTypeConvertRule(t2)
		return tt
	}
	panic("missing check" + e.OpName())
	return nil
}
