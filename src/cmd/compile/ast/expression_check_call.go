package ast

import (
	"fmt"
)

func (e *Expression) checkMethodCallExpression(block *Block, errs *[]error) []*VariableType {
	call := e.Data.(*ExpressionMethodCall)
	ts, es := call.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	object, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if object.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
		switch call.Name {
		case "size", "start", "end", "cap":
			t := &VariableType{}
			t.Typ = VARIABLE_TYPE_INT
			t.Pos = e.Pos
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s too mamy argument to call,method '%s' expect no arguments",
					errMsgPrefix(e.Pos), call.Name))
			}
			return []*VariableType{t}
		case "append":
			if len(call.Args) == 0 {
				*errs = append(*errs, fmt.Errorf("%s too mamy argument to call,method '%s' expect no arguments",
					errMsgPrefix(e.Pos), call.Name))
			}
			for _, e := range call.Args {
				_, es := e.check(block)
				if errsNotEmpty(es) {
					*errs = append(*errs, es...)
				}
				//				for _, t := range ts {

				//				}
			}
			t := &VariableType{}
			t.Typ = VARIABLE_TYPE_VOID
			t.Pos = e.Pos
			return []*VariableType{t}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on array", errMsgPrefix(e.Pos), call.Name))
		}
		return nil
	}
	if object.Typ != VARIABLE_TYPE_OBJECT && object.Typ != VARIABLE_TYPE_CLASS {
		*errs = append(*errs, fmt.Errorf("%s cannot make method call named '%s' on '%s'", errMsgPrefix(e.Pos), call.Name, object.TypeString()))
		return nil
	}
	args := checkExpressions(block, call.Args, errs)
	args = checkRightValuesValid(args, errs)
	f, es := object.Class.accessMethod(call.Name, e.Pos, args)
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
		*errs = append(*errs, fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), call.Expression.OpName()))
		t = &VariableType{
			Typ: VARIABLE_TYPE_VOID,
			Pos: e.Pos,
		}
		return nil
	}
	if t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s %s is not a function", errMsgPrefix(e.Pos), call.Expression.OpName()))
		t = &VariableType{
			Typ: VARIABLE_TYPE_VOID,
			Pos: e.Pos,
		}
		return []*VariableType{t}
	}
	call.Func = t.Function
	if t.Function.Isbuildin {
		return e.checkBuildinFunctionCall(block, errs, t.Function, call.Args)
	} else {
		return e.checkFunctionCall(block, errs, t.Function, call.Args)
	}
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, args []*Expression) []*VariableType {
	callargsTypes := checkExpressions(block, args, errs)
	callargsTypes = checkRightValuesValid(callargsTypes, errs)
	if len(callargsTypes) > len(f.Typ.ParameterList) {
		*errs = append(*errs, fmt.Errorf("%s too many paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}
	if len(callargsTypes) < len(f.Typ.ParameterList) && len(args) < len(f.Typ.ParameterList) {
		*errs = append(*errs, fmt.Errorf("%s too few paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}
	for k, v := range f.Typ.ParameterList {
		if k < len(callargsTypes) {
			if !v.Typ.TypeCompatible(callargsTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s",
					errMsgPrefix(args[k].Pos),
					v.Typ.TypeString(),
					callargsTypes[k].TypeString()))
			}
		}
	}
	return f.Typ.ReturnList.retTypes(e.Pos)
}
