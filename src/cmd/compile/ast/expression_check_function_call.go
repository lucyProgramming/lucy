package ast

import "fmt"

func (e *Expression) checkFunctionCallExpression(block *Block, errs *[]error) []*VariableType {
	ret := []*VariableType{mkVoidType(e.Pos)}
	call := e.Data.(*ExpressionFunctionCall)
	t, es := call.Expression.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
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
		f := f // override f
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
				if !v.Typ.Equal(errs, callargsTypes[k]) {
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

func (e *Expression) checkTemplateFunctionCall(block *Block, errs *[]error,
	argTypes []*VariableType, f *Function) (ret *Function) {
	call := e.Data.(*ExpressionFunctionCall)
	typeParameters := make(map[string]*VariableType)
	tArgs := []*VariableType{}
	tIndexes := []int{}
	for k, v := range f.Typ.ParameterList {
		if v == nil || v.Typ == nil || v.Typ.Typ != VARIABLE_TYPE_T {
			continue
		}
		if k >= len(argTypes) || argTypes[k] == nil {
			*errs = append(*errs, fmt.Errorf("%s missing typed parameter,index at %d",
				errMsgPrefix(e.Pos), k))
			return
		}
		if v.Typ.canBeBindWith(argTypes[k]) == false {
			*errs = append(*errs, fmt.Errorf("%s cannot bind '%s' to '%s'",
				errMsgPrefix(e.Pos), v.Typ.TypeString(), argTypes[k].TypeString()))
		}
		name := v.Typ.Name
		typeParameters[name] = v.Typ
		tArgs = append(tArgs, argTypes[k])
		tIndexes = append(tIndexes, k)
	}
	tps := call.TypedParameters
	retTypes := []*VariableType{}
	retIndexes := []int{}
	for k, v := range f.Typ.ReturnList {
		if v == nil || v.Typ == nil || v.Typ.Typ != VARIABLE_TYPE_T {
			continue
		}
		if len(tps) == 0 || tps[0] == nil {
			*errs = append(*errs, fmt.Errorf("%s missing typed return value,index at %d",
				errMsgPrefix(e.Pos), k))
			continue
		}
		if v.Typ.canBeBindWith(tps[0]) == false {
			*errs = append(*errs, fmt.Errorf("%s cannot bind '%s' to '%s'",
				errMsgPrefix(e.Pos), v.Typ.TypeString(), tps[0].TypeString()))
		}
		name := v.Typ.Name
		retTypes = append(retTypes, tps[0])
		typeParameters[name] = tps[0]
		retIndexes = append(retIndexes, k)
		tps = tps[1:]
	}
	call.TemplateFunctionCallPair = f.TemplateFunction.insert(tArgs, retTypes, ret, errs)
	if call.TemplateFunctionCallPair.Function == nil { // not called before,make the binds
		t, es := f.clone(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			return nil
		}
		t.TemplateFunction = nil
		call.TemplateFunctionCallPair.Function = t
		for k, v := range tArgs {
			index := tIndexes[k]
			pos := f.Typ.ParameterList[index].Pos //
			*call.TemplateFunctionCallPair.Function.Typ.ParameterList[index].Typ =
				*v
			call.TemplateFunctionCallPair.Function.Typ.ParameterList[index].Typ.Pos = pos
		}
		for k, v := range retTypes {
			index := retIndexes[k]
			pos := f.Typ.ReturnList[index].Pos //
			*call.TemplateFunctionCallPair.Function.Typ.ReturnList[index].Typ =
				*v
			call.TemplateFunctionCallPair.Function.Typ.ReturnList[index].Typ.Pos = pos
			call.TemplateFunctionCallPair.Function.Typ.ReturnList[index].Expression =
				v.mkDefaultValueExpression()
		}
	}
	ret = call.TemplateFunctionCallPair.Function
	// when all ok ,ret is not a template function any more
	if len(tps) > 0 {
		*errs = append(*errs, fmt.Errorf("%s to many typed parameter to call template function.",
			errMsgPrefix(e.Pos)))
	}
	return ret
}
