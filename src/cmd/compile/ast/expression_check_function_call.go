package ast

import "fmt"

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
	if t.Typ == VARIABLE_TYPE_CLASS { // cast type
		ret := make([]*VariableType, 1)
		ret[0] = &VariableType{}
		ret[0].Typ = VARIABLE_TYPE_OBJECT
		ret[0].Class = t.Class
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument", errMsgPrefix(e.Pos)))
			return ret
		}
		e.Typ = EXPRESSION_TYPE_CONVERTION_TYPE
		convertType := &ExpressionTypeConvertion{}
		convertType.Expression = call.Args[0]
		convertType.Typ = ret[0]
		e.Data = convertType
		ts, es := call.Args[0].check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		t, err := call.Args[0].mustBeOneValueContext(ts)
		if err != nil {
			*errs = append(*errs, err)
		}
		if t == nil {
			return ret
		}
		if t.IsPrimitive() {
			*errs = append(*errs, fmt.Errorf("%s expression is primitive,cannot be cast to another type", errMsgPrefix(e.Pos)))
		}
		return ret
	}
	if t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s %s is not a function,but '%s'", errMsgPrefix(e.Pos), call.Expression.OpName(), t.TypeString()))
		t = &VariableType{
			Typ: VARIABLE_TYPE_VOID,
			Pos: e.Pos,
		}
		return []*VariableType{t}
	}
	call.Func = t.Function
	if t.Function.IsBuildin {
		return e.checkBuildinFunctionCall(block, errs, t.Function, call.Args)
	} else {
		return e.checkFunctionCall(block, errs, t.Function, &call.Args)
	}
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, args *CallArgs) []*VariableType {
	callargsTypes := checkExpressions(block, *args, errs)
	callargsTypes = checkRightValuesValid(callargsTypes, errs)
	if len(callargsTypes) > len(f.Typ.ParameterList) {
		errmsg := fmt.Sprintf("%s too many paramaters to call function '%s':\n", errMsgPrefix(e.Pos), f.Name)
		errmsg += fmt.Sprintf("\t have %s\n", f.badParameterMsg(f.Name, callargsTypes))
		errmsg += fmt.Sprintf("\t want %s\n", f.readableMsg())
		*errs = append(*errs, fmt.Errorf(errmsg))
	}
	ret := f.Typ.ReturnList.retTypes(e.Pos)
	if f.HaveDefaultValue {
		if len(callargsTypes) < f.DefaultValueStartAt {
			*errs = append(*errs, fmt.Errorf("%s too few paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
			return ret
		}
		for i := len(callargsTypes); i < len(f.Typ.ParameterList); i++ {
			*args = append(*args, f.Typ.ParameterList[i].Expression)
		}
	} else { // no default value
		if len(callargsTypes) < len(f.Typ.ParameterList) && len(*args) < len(f.Typ.ParameterList) && f.HaveDefaultValue == false {
			*errs = append(*errs, fmt.Errorf("%s too few paramaters to call function %s", errMsgPrefix(e.Pos), f.Name))
			return ret
		}
	}
	for k, v := range f.Typ.ParameterList {
		if k < len(callargsTypes) {
			if !v.Typ.TypeCompatible(callargsTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s type '%s' is not compatible with '%s'",
					errMsgPrefix((*args)[k].Pos),
					v.Typ.TypeString(),
					callargsTypes[k].TypeString()))
			}
		}
	}
	return ret
}
