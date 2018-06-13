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
	if t.Type == VARIABLE_TYPE_CLASS { // cast type
		typeConversion := &ExpressionTypeConversion{}
		typeConversion.Type = &VariableType{}
		typeConversion.Type.Type = VARIABLE_TYPE_OBJECT
		typeConversion.Type.Class = t.Class
		typeConversion.Type.Pos = e.Pos
		ret := []*VariableType{typeConversion.Type}
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
				errMsgPrefix(e.Pos)))
			return ret
		}
		e.Type = EXPRESSION_TYPE_CHECK_CAST
		typeConversion.Expression = call.Args[0]
		e.Data = typeConversion
		e.checkTypeConvertionExpression(block, errs)
		return ret
	}
	if t.Type != VARIABLE_TYPE_FUNCTION {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not a function,but '%s'",
			errMsgPrefix(e.Pos),
			call.Expression.OpName(), t.TypeString()))
		return nil
	}
	call.Func = t.Function
	if t.Function.IsBuildIn {
		return e.checkBuildInFunctionCall(block, errs, t.Function, call.Args)
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
	if len(callArgsTypes) > len(f.Type.ParameterList) {
		errMsg := fmt.Sprintf("%s too many paramaters to call function '%s':\n", errMsgPrefix(e.Pos), f.Name)
		errMsg += fmt.Sprintf("\thave %s\n", f.badParameterMsg(f.Name, callArgsTypes))
		errMsg += fmt.Sprintf("\twant %s\n", f.readableMsg())
		*errs = append(*errs, fmt.Errorf(errMsg))
	}
	//trying to convert literal
	var ret []*VariableType
	convertLiteralExpressionsToNeeds(*args, f.Type.getParameterTypes(), callArgsTypes)
	if f.TemplateFunction == nil {
		ret = f.Type.retTypes(e.Pos)
	} else {
		ret = tf.Type.retTypes(e.Pos)
	}
	{
		f := f // override f
		if f.TemplateFunction != nil {
			f = tf
		}
		if len(callArgsTypes) < len(f.Type.ParameterList) {
			if f.HaveDefaultValue && len(callArgsTypes) >= f.DefaultValueStartAt {
				for i := len(callArgsTypes); i < len(f.Type.ParameterList); i++ {
					*args = append(*args, f.Type.ParameterList[i].Expression)
				}
			} else { // no default value
				errMsg := fmt.Sprintf("%s too few paramaters to call function '%s'\n", errMsgPrefix(e.Pos), f.Name)
				errMsg += fmt.Sprintf("\thave %s\n", f.badParameterMsg(f.Name, callArgsTypes))
				errMsg += fmt.Sprintf("\twant %s\n", f.readableMsg())
				*errs = append(*errs, fmt.Errorf(errMsg))
				return ret
			}
		}
		for k, v := range f.Type.ParameterList {
			if k < len(callArgsTypes) && callArgsTypes[k] != nil {
				if !v.Type.Equal(errs, callArgsTypes[k]) {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix((*args)[k].Pos),
						callArgsTypes[k].TypeString(), v.Type.TypeString()))
				}
			}
		}
	}
	return ret
}

func (e *Expression) checkTemplateFunctionCall(block *Block, errs *[]error,
	argTypes []*VariableType, f *Function) (ret *Function) {
	call := e.Data.(*ExpressionFunctionCall)
	parameterType := make(map[string]*VariableType)
	for k, v := range f.Type.ParameterList {
		if v == nil || v.Type == nil || len(v.Type.haveParameterType()) == 0 {
			continue
		}
		if k >= len(argTypes) || argTypes[k] == nil {
			*errs = append(*errs, fmt.Errorf("%s missing typed parameter,index at %d",
				errMsgPrefix(e.Pos), k))
			return
		}
		if err := v.Type.canBeBindWithType(parameterType, argTypes[k]); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v",
				errMsgPrefix(argTypes[k].Pos), err))
			return
		}
	}
	tps := call.TypedParameters
	for k, v := range f.Type.ReturnList {
		if v == nil || v.Type == nil || len(v.Type.haveParameterType()) == 0 {
			continue
		}
		if len(tps) == 0 || tps[0] == nil {
			//trying already have
			if err := v.Type.canBeBindWithTypedParameters(parameterType); err == nil {
				//very good no error
				continue
			}
			*errs = append(*errs, fmt.Errorf("%s missing typed return value,index at %d",
				errMsgPrefix(e.Pos), k))
			return
		}
		if err := v.Type.canBeBindWithType(parameterType, tps[0]); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v",
				errMsgPrefix(tps[0].Pos), err))
			return nil
		}
		tps = tps[1:]
	}
	call.TemplateFunctionCallPair = f.TemplateFunction.insert(parameterType, ret, errs)
	if call.TemplateFunctionCallPair.Function == nil { // not called before,make the binds
		cloneFunction, es := f.clone()
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			return nil
		}
		cloneFunction.TemplateFunction = nil
		call.TemplateFunctionCallPair.Function = cloneFunction
		cloneFunction.parameterType = parameterType
		for _, v := range cloneFunction.Type.ParameterList {
			if len(v.Type.haveParameterType()) > 0 {
				v.Type.bindWithTypedParameters(parameterType)
			}
		}
		for _, v := range cloneFunction.Type.ReturnList {
			if len(v.Type.haveParameterType()) > 0 {
				v.Type.bindWithTypedParameters(parameterType)
			}
		}
		//check this function
		if cloneFunction.Block.Functions == nil {
			cloneFunction.Block.Functions = make(map[string]*Function)
		}
		cloneFunction.Block.Functions[cloneFunction.Name] = cloneFunction
		cloneFunction.Block.inherit(&PackageBeenCompile.Block)
		cloneFunction.Block.InheritedAttribute.Function = cloneFunction
		cloneFunction.checkParametersAndReturns(errs)
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
