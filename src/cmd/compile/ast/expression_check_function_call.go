package ast

import "fmt"

func (e *Expression) checkFunctionCallExpression(block *Block, errs *[]error) []*VariableType {
	call := e.Data.(*ExpressionFunctionCall)
	t, es := call.Expression.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
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
		return nil
	}
	call.Func = t.Function
	if t.Function.IsBuildIn {
		return e.checkBuildinFunctionCall(block, errs, t.Function, call.Args)
	} else {
		return e.checkFunctionCall(block, errs, t.Function, &call.Args)

	}
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, args *CallArgs) []*VariableType {
	callArgsTypes := checkExpressions(block, *args, errs)
	callArgsTypes = checkRightValuesValid(callArgsTypes, errs)
	var tf *Function
	if f.TemplateFunction != nil {
		length := len(*errs)
		//rewrite
		tf = e.checkTemplateFunctionCall(block, errs, callArgsTypes, f)
		if len(*errs) != length { // if no
			return nil
		}
	} else { // not template function
		call := e.Data.(*ExpressionFunctionCall)
		if len(call.TypedParameters) > 0 {
			*errs = append(*errs, fmt.Errorf("%s function is not a template function,cannot not have typed parameters",
				errMsgPrefix(e.Pos)))
		}
	}
	if len(callArgsTypes) > len(f.Typ.ParameterList) {
		errMsg := fmt.Sprintf("%s too many paramaters to call function '%s':\n", errMsgPrefix(e.Pos), f.Name)
		errMsg += fmt.Sprintf("\thave %s\n", f.badParameterMsg(f.Name, callArgsTypes))
		errMsg += fmt.Sprintf("\twant %s\n", f.readableMsg())
		*errs = append(*errs, fmt.Errorf(errMsg))
	}
	//trying to convert literal
	var ret []*VariableType
	convertLiteralExpressionsToNeeds(*args, f.Typ.needParameterTypes(), callArgsTypes)
	if f.TemplateFunction == nil {
		ret = f.Typ.retTypes(e.Pos)
	} else {
		ret = tf.Typ.retTypes(e.Pos)
	}
	{
		f := f // override f
		if f.TemplateFunction != nil {
			f = tf
		}
		if len(callArgsTypes) < len(f.Typ.ParameterList) {
			if f.HaveDefaultValue && len(callArgsTypes) >= f.DefaultValueStartAt {
				for i := len(callArgsTypes); i < len(f.Typ.ParameterList); i++ {
					*args = append(*args, f.Typ.ParameterList[i].Expression)
				}
			} else { // no default value
				errMsg := fmt.Sprintf("%s too few paramaters to call function '%s'\n", errMsgPrefix(e.Pos), f.Name)
				errMsg += fmt.Sprintf("\thave %s\n", f.badParameterMsg(f.Name, callArgsTypes))
				errMsg += fmt.Sprintf("\twant %s\n", f.readableMsg())
				*errs = append(*errs, fmt.Errorf(errMsg))
				return ret
			}
		}
		for k, v := range f.Typ.ParameterList {
			if k < len(callArgsTypes) && callArgsTypes[k] != nil {
				if !v.Typ.Equal(errs, callArgsTypes[k]) {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix((*args)[k].Pos),
						callArgsTypes[k].TypeString(), v.Typ.TypeString()))
				}
			}
		}
	}
	return ret
}

func (e *Expression) checkTemplateFunctionCall(block *Block, errs *[]error,
	argTypes []*VariableType, f *Function) (ret *Function) {
	call := e.Data.(*ExpressionFunctionCall)
	typedParameters := make(map[string]*VariableType)
	for k, v := range f.Typ.ParameterList {
		if v == nil || v.Typ == nil || len(v.Typ.haveT()) == 0 {
			continue
		}
		if k >= len(argTypes) || argTypes[k] == nil {
			//			//trying already have
			//			if err := v.Typ.canBeBindWithTypedParameters(typedParameters); err == nil {
			//				continue
			//			}
			*errs = append(*errs, fmt.Errorf("%s missing typed parameter,index at %d",
				errMsgPrefix(e.Pos), k))
			return
		}
		if err := v.Typ.canBebindWithType(typedParameters, argTypes[k]); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v",
				errMsgPrefix(argTypes[k].Pos), err))
			return
		}
	}
	tps := call.TypedParameters
	for k, v := range f.Typ.ReturnList {
		if v == nil || v.Typ == nil || len(v.Typ.haveT()) == 0 {
			continue
		}
		if len(tps) == 0 || tps[0] == nil {
			//trying already have
			if err := v.Typ.canBeBindWithTypedParameters(typedParameters); err == nil {
				//very good no error
				continue
			}
			*errs = append(*errs, fmt.Errorf("%s missing typed return value,index at %d",
				errMsgPrefix(e.Pos), k))
			return
		}
		if err := v.Typ.canBebindWithType(typedParameters, tps[0]); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v",
				errMsgPrefix(tps[0].Pos), err))
			return nil
		}
		tps = tps[1:]
	}
	call.TemplateFunctionCallPair = f.TemplateFunction.insert(typedParameters, ret, errs)
	if call.TemplateFunctionCallPair.Function == nil { // not called before,make the binds
		cloneFunction, es := f.clone()
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			return nil
		}
		cloneFunction.TemplateFunction = nil
		call.TemplateFunctionCallPair.Function = cloneFunction
		cloneFunction.TypeParameters = typedParameters
		for _, v := range cloneFunction.Typ.ParameterList {
			if len(v.Typ.haveT()) > 0 {
				v.Typ.bindWithTypedParameters(typedParameters)
			}
		}
		for _, v := range cloneFunction.Typ.ReturnList {
			if len(v.Typ.haveT()) > 0 {
				v.Typ.bindWithTypedParameters(typedParameters)
			}
		}
		//check this function
		if cloneFunction.Block.Funcs == nil {
			cloneFunction.Block.Funcs = make(map[string]*Function)
		}
		cloneFunction.Block.Funcs[cloneFunction.Name] = cloneFunction
		cloneFunction.Block.inherit(&PackageBeenCompile.Block)
		cloneFunction.Block.InheritedAttribute.Function = cloneFunction
		cloneFunction.checkParametersAndRetuns(errs)
		cloneFunction.checkBlock(errs)
	}
	ret = call.TemplateFunctionCallPair.Function
	// when all ok ,ret is not a template function any more
	if len(tps) > 0 {
		*errs = append(*errs, fmt.Errorf("%s to many typed parameter to call template function",
			errMsgPrefix(e.Pos)))
	}
	return ret
}
