package ast

import (
	"fmt"
)

func (e *Expression) checkFunctionCallExpression(block *Block, errs *[]error) []*Type {
	call := e.Data.(*ExpressionFunctionCall)
	callExpression, es := call.Expression.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if callExpression == nil {
		return nil
	}
	if callExpression.Type == VariableTypeClass { // cast type
		typeConversion := &ExpressionTypeConversion{}
		typeConversion.Type = &Type{}
		typeConversion.Type.Type = VariableTypeObject
		typeConversion.Type.Class = callExpression.Class
		typeConversion.Type.Pos = e.Pos
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
				errMsgPrefix(e.Pos)))
			return nil
		}
		e.Type = ExpressionTypeCheckCast
		typeConversion.Expression = call.Args[0]
		e.Data = typeConversion
		ret := e.checkTypeConversionExpression(block, errs)
		if ret == nil {
			return nil
		}
		return []*Type{ret}
	}
	if callExpression.Type == VariableTypeTypeAlias {
		typeConversion := &ExpressionTypeConversion{}
		typeConversion.Type = callExpression.AliasType
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
				errMsgPrefix(e.Pos)))
			return nil
		}
		e.Type = ExpressionTypeCheckCast
		typeConversion.Expression = call.Args[0]
		e.Data = typeConversion
		ret := e.checkTypeConversionExpression(block, errs)
		if ret == nil {
			return nil
		}
		return []*Type{ret}
	}
	if callExpression.Type != VariableTypeFunction {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not a function,but '%s'",
			errMsgPrefix(e.Pos),
			call.Expression.OpName(), callExpression.TypeString()))
		return nil
	}
	if callExpression.Function != nil {
		call.Function = callExpression.Function
		if callExpression.Function.IsBuildIn {
			return e.checkBuildInFunctionCall(block, errs, callExpression.Function, call)
		} else {
			return e.checkFunctionCall(block, errs, callExpression.Function, call)
		}
	}
	return e.checkFunctionPointerCall(block, errs, callExpression.FunctionType, call)
}

func (e *Expression) checkFunctionPointerCall(block *Block, errs *[]error, ft *FunctionType, call *ExpressionFunctionCall) []*Type {
	callArgsTypes := checkExpressions(block, call.Args, errs)
	callArgsTypes = checkRightValuesValid(callArgsTypes, errs)
	if len(call.ParameterTypes) > 0 {
		*errs = append(*errs, fmt.Errorf("%s function is not a template function,cannot not have typed parameters",
			errMsgPrefix(e.Pos)))
	}

	if len(callArgsTypes) > len(ft.ParameterList) {
		errMsg := fmt.Sprintf("%s too many paramaters to call\n", errMsgPrefix(e.Pos))
		errMsg += fmt.Sprintf("\thave %s\n", functionPointerCallHave(callArgsTypes))
		errMsg += fmt.Sprintf("\twant %s\n", functionPointerCallWant(ft.ParameterList))
		*errs = append(*errs, fmt.Errorf(errMsg))
	}
	//trying to convert literal
	var ret []*Type
	convertLiteralExpressionsToNeeds(call.Args, ft.getParameterTypes(), callArgsTypes)
	if len(callArgsTypes) < len(ft.ParameterList) {
		errMsg := fmt.Sprintf("%s too few paramaters to call\n", errMsgPrefix(e.Pos))
		errMsg += fmt.Sprintf("\thave %s\n", functionPointerCallHave(callArgsTypes))
		errMsg += fmt.Sprintf("\twant %s\n", functionPointerCallWant(ft.ParameterList))
		*errs = append(*errs, fmt.Errorf(errMsg))
		return ret
	}
	for k, v := range ft.ParameterList {
		if k < len(callArgsTypes) && callArgsTypes[k] != nil {
			if false == v.Type.Equal(errs, callArgsTypes[k]) {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix((call.Args)[k].Pos),
					callArgsTypes[k].TypeString(), v.Type.TypeString()))
			}
		}
	}
	return ret
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, call *ExpressionFunctionCall) []*Type {
	callArgsTypes := checkExpressions(block, call.Args, errs)
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
		if len(call.ParameterTypes) > 0 {
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
	var ret []*Type
	convertLiteralExpressionsToNeeds(call.Args, f.Type.getParameterTypes(), callArgsTypes)
	if f.TemplateFunction == nil {
		ret = f.Type.returnTypes(e.Pos)
	} else {
		ret = tf.Type.returnTypes(e.Pos)
	}
	{
		f := f // override f
		if f.TemplateFunction != nil {
			f = tf
		}
		if len(callArgsTypes) < len(f.Type.ParameterList) {
			if f.HaveDefaultValue && len(callArgsTypes) >= f.DefaultValueStartAt {
				for i := len(callArgsTypes); i < len(f.Type.ParameterList); i++ {
					call.Args = append(call.Args, f.Type.ParameterList[i].Expression)
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
						errMsgPrefix((callArgsTypes)[k].Pos),
						callArgsTypes[k].TypeString(), v.Type.TypeString()))
				}
			}
		}
	}
	return ret
}

func (e *Expression) checkTemplateFunctionCall(block *Block, errs *[]error,
	argTypes []*Type, f *Function) (ret *Function) {
	call := e.Data.(*ExpressionFunctionCall)
	parameterTypes := make(map[string]*Type)
	for k, v := range f.Type.ParameterList {
		if v == nil || v.Type == nil || len(v.Type.haveParameterType()) == 0 {
			continue
		}
		if k >= len(argTypes) || argTypes[k] == nil {
			*errs = append(*errs, fmt.Errorf("%s missing typed parameter,index at %d",
				errMsgPrefix(e.Pos), k))
			return
		}
		if err := v.Type.canBeBindWithType(parameterTypes, argTypes[k]); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v",
				errMsgPrefix(argTypes[k].Pos), err))
			return
		}
	}
	tps := call.ParameterTypes
	for k, v := range f.Type.ReturnList {
		if v == nil || v.Type == nil || len(v.Type.haveParameterType()) == 0 {
			continue
		}
		if len(tps) == 0 || tps[0] == nil {
			//trying already have
			if err := v.Type.canBeBindWithParameterTypes(parameterTypes); err == nil {
				//very good no error
				continue
			}
			*errs = append(*errs, fmt.Errorf("%s missing typed return value,index at %d",
				errMsgPrefix(e.Pos), k))
			return
		}
		if err := v.Type.canBeBindWithType(parameterTypes, tps[0]); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v",
				errMsgPrefix(tps[0].Pos), err))
			return nil
		}
		tps = tps[1:]
	}
	call.TemplateFunctionCallPair = f.TemplateFunction.insert(parameterTypes, ret, errs)
	if call.TemplateFunctionCallPair.Function == nil { // not called before,make the binds
		cloneFunction, es := f.clone()
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
			return nil
		}
		cloneFunction.TemplateFunction = nil
		call.TemplateFunctionCallPair.Function = cloneFunction
		cloneFunction.parameterTypes = parameterTypes
		for _, v := range cloneFunction.Type.ParameterList {
			if len(v.Type.haveParameterType()) > 0 {
				v.Type.bindWithParameterTypes(parameterTypes)
			}
		}
		for _, v := range cloneFunction.Type.ReturnList {
			if len(v.Type.haveParameterType()) > 0 {
				v.Type.bindWithParameterTypes(parameterTypes)
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
		*errs = append(*errs, fmt.Errorf("%s to many parameter type  to call template function",
			errMsgPrefix(e.Pos)))
	}
	return ret
}
