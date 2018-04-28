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
		return nil
	}
	if t.Typ == VARIABLE_TYPE_CLASS { // cast type
		ret := make([]*VariableType, 1)
		ret[0] = &VariableType{}
		ret[0].Typ = VARIABLE_TYPE_OBJECT
		ret[0].Class = t.Class
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
				errMsgPrefix(e.Pos)))
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
			*errs = append(*errs, fmt.Errorf("%s expression is primitive(non-pointer),cannot be cast to pointer",
				errMsgPrefix(e.Pos)))
		}
		return ret
	}
	if t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s %s is not a function,but '%s'",
			errMsgPrefix(e.Pos),
			call.Expression.OpName(), t.TypeString()))
		return nil
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
		errmsg += fmt.Sprintf("\thave %s\n", f.badParameterMsg(f.Name, callargsTypes))
		errmsg += fmt.Sprintf("\twant %s\n", f.readableMsg())
		*errs = append(*errs, fmt.Errorf(errmsg))
	}
	ret := f.Typ.ReturnList.retTypes(e.Pos)
	if len(callargsTypes) < len(f.Typ.ParameterList) {
		if f.HaveDefaultValue && len(callargsTypes) >= f.DefaultValueStartAt {
			for i := len(callargsTypes); i < len(f.Typ.ParameterList); i++ {
				*args = append(*args, f.Typ.ParameterList[i].Expression)
			}
		} else { // no default value
			errmsg := fmt.Sprintf("%s too few paramaters to call function %s\n", errMsgPrefix(e.Pos), f.Name)
			errmsg += fmt.Sprintf("\thave %s\n", f.badParameterMsg(f.Name, callargsTypes))
			errmsg += fmt.Sprintf("\twant %s\n", f.readableMsg())
			*errs = append(*errs, fmt.Errorf(errmsg))
			return ret
		}
	}
	for k, v := range f.Typ.ParameterList {
		if k < len(callargsTypes) && callargsTypes[k] != nil {
			if !v.Typ.TypeCompatible(callargsTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix((*args)[k].Pos),
					callargsTypes[k].TypeString(), v.Typ.TypeString()))
			}
		}
	}
	return ret
}
