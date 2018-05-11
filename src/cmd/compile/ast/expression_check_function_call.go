package ast

import "fmt"

func (e *Expression) checkFunctionCallExpression(block *Block, errs *[]error) []*VariableType {
	ret := []*VariableType{mkVoidType(e.Pos)}
	call := e.Data.(*ExpressionFunctionCall)
	ts, es := call.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return ret
	}
	if t.Typ == VARIABLE_TYPE_CLASS { // cast type
		convertType := &ExpressionTypeConvertion{}
		convertType.Typ = &VariableType{}
		convertType.Typ.Typ = VARIABLE_TYPE_OBJECT
		convertType.Typ.Class = t.Class
		convertType.Typ.Pos = e.Pos
		ret := []*VariableType{convertType.Typ}
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
				errMsgPrefix(e.Pos)))
			return ret
		}
		e.Typ = EXPRESSION_TYPE_CHECK_CAST
		convertType.Expression = call.Args[0]
		e.Data = convertType
		e.checkTypeConvertionExpression(block, errs)
		return ret
	}
	if t.Typ != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not a function,but '%s'",
			errMsgPrefix(e.Pos),
			call.Expression.OpName(), t.TypeString()))
		return ret
	}
	call.Func = t.Function
	if t.Function.IsBuildin {
		return e.checkBuildinFunctionCall(block, errs, t.Function, call.Args)
	} else {
		ret = e.checkFunctionCall(block, errs, t.Function, &call.Args)
		return ret
	}
}

func (e *Expression) checkTemplateFunctionCall(block *Block, errs *[]error,
	argTypes []*VariableType, f *Function) (ret *Function) {
	call := e.Data.(*ExpressionFunctionCall)
	ret, es := f.clone(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
		return ret
	}
	typeParameters := make(map[string]*VariableType)
	for k, v := range ret.Typ.ParameterList {
		if v == nil && v.Typ == nil && v.Typ.Typ != VARIABLE_TYPE_T {
			continue
		}
		if k > len(argTypes) || argTypes[k] == nil {
			*errs = append(*errs, fmt.Errorf("%s missing %d typed parameter", k))
			return
		}
		pos := v.Typ.Pos // keep pos
		name := v.Typ.Name
		*v.Typ = *argTypes[k]
		v.Typ.Pos = pos
		typeParameters[name] = v.Typ
	}
	tps := call.TypedParameters
	retTypes := []*VariableType{}
	for k, v := range ret.Typ.ReturnList {
		if v.Typ == nil {
			continue
		}
		if v == nil && v.Typ == nil && v.Typ.Typ != VARIABLE_TYPE_T {
			retTypes = append(retTypes, v.Typ)
			continue
		}
		if len(tps) == 0 || tps[0] == nil {
			*errs = append(*errs, fmt.Errorf("%s missing %d return type", k))
			continue
		}
		name := v.Typ.Name
		pos := v.Typ.Pos // keep pos
		*v.Typ = *tps[0]
		v.Typ.Pos = pos
		tps = tps[1:]
		retTypes = append(retTypes, v.Typ)
		typeParameters[name] = v.Typ
		v.Expression = v.Typ.mkDefaultValueExpression()
	}
	call.TemplateFunctionCallPair = f.TemplateFunction.insert(argTypes, retTypes, ret)
	ret.TypeParameters = typeParameters
	// when all ok ,tf is not  template function
	ret.TemplateFunction = nil
	return ret
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, args *CallArgs) []*VariableType {
	callargsTypes := checkExpressions(block, *args, errs)
	callargsTypes = checkRightValuesValid(callargsTypes, errs)
	ret := []*VariableType{mkVoidType(e.Pos)}
	var tf *Function
	if f.TemplateFunction != nil {
		length := len(*errs)
		//rewrite
		tf = e.checkTemplateFunctionCall(block, errs, callargsTypes, f)
		if len(*errs) != length { // if no
			return ret
		}
	}
	if len(callargsTypes) > len(f.Typ.ParameterList) {
		errmsg := fmt.Sprintf("%s too many paramaters to call function '%s':\n", errMsgPrefix(e.Pos), f.Name)
		errmsg += fmt.Sprintf("\thave %s\n", f.badParameterMsg(f.Name, callargsTypes))
		errmsg += fmt.Sprintf("\twant %s\n", f.readableMsg())
		*errs = append(*errs, fmt.Errorf(errmsg))
	}
	//trying to convert literal
	convertLiteralExpressionsToNeeds(*args, f.Typ.needParameterTypes(), callargsTypes)
	if f.TemplateFunction == nil {
		ret = f.Typ.retTypes(e.Pos)
	} else {
		ret = tf.Typ.retTypes(e.Pos)
	}
	{
		f := f
		if f.TemplateFunction != nil {
			f = tf
		}
		if len(callargsTypes) < len(f.Typ.ParameterList) {
			if f.HaveDefaultValue && len(callargsTypes) >= f.DefaultValueStartAt {
				for i := len(callargsTypes); i < len(f.Typ.ParameterList); i++ {
					*args = append(*args, f.Typ.ParameterList[i].Expression)
				}
			} else { // no default value
				errmsg := fmt.Sprintf("%s too few paramaters to call function '%s'\n", errMsgPrefix(e.Pos), f.Name)
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
	}
	if f.TemplateFunction != nil {
		if tf.Block.Funcs == nil {
			tf.Block.Funcs = make(map[string]*Function)
		}
		tf.Block.Funcs[tf.Name] = tf //incase recursively
		tf.checkBlock(errs)
	}
	return ret
}
