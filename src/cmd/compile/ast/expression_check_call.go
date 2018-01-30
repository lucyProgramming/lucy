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
	if t.Typ != VARIABLE_TYPE_OBJECT && t.Typ != VARIABLE_TYPE_CLASS {
		*errs = append(*errs, fmt.Errorf("%s method call only can be made on 'object' or 'class'", errMsgPrefix(e.Pos)))
		return nil
	}

	args := e.checkExpressions(block, call.Args, errs)
	args = e.checkRightValuesValid(args, errs)
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
	if t == nil || t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), e.OpName()))
		t = &VariableType{
			Typ: VARIABLE_TYPE_VOID,
			Pos: e.Pos,
		}
		return []*VariableType{t}
	}
	if t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s %s is not a function", errMsgPrefix(e.Pos), e.OpName()))
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
	callargsTypes := e.checkExpressions(block, args, errs)
	callargsTypes = e.checkRightValuesValid(callargsTypes, errs)
	if len(callargsTypes) > len(f.Typ.ParameterList) {
		*errs = append(*errs, fmt.Errorf("%s too many paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}
	if len(callargsTypes) < len(f.Typ.ParameterList) && len(args) < len(f.Typ.ParameterList) {
		*errs = append(*errs, fmt.Errorf("%s too few paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
	}
	for k, v := range f.Typ.ParameterList {
		if k < len(callargsTypes) {
			if !v.Typ.typeCompatible(callargsTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s type %s is not compatible with %s",
					errMsgPrefix(args[k].Pos),
					v.Typ.TypeString(),
					callargsTypes[k].TypeString()))
			}
		}
	}
	return f.Typ.ReturnList.retTypes(e.Pos)
}
