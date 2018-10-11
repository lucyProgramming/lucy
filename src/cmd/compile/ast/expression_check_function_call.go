package ast

import (
	"fmt"
)

func (e *Expression) checkFunctionCallExpression(block *Block, errs *[]error) []*Type {
	call := e.Data.(*ExpressionFunctionCall)
	if call.Expression.Type == ExpressionTypeIdentifier {
		identifier := call.Expression.Data.(*ExpressionIdentifier)
		d, err := block.searchIdentifier(call.Expression.Pos, identifier.Name)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		switch d.(type) {
		case *Function:
			f := d.(*Function)
			call.Function = f
			//if f.IsBuildIn {
			//	return e.checkBuildInFunctionCall(block, errs, f, call)
			//} else {
			return e.checkFunctionCall(block, errs, f, call)
			//}
		case *Type:
			typeConversion := &ExpressionTypeConversion{}
			typeConversion.Type = d.(*Type)
			if len(call.Args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
					e.Pos.errMsgPrefix()))
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
		case *Class:
			typeConversion := &ExpressionTypeConversion{}
			typeConversion.Type = &Type{}
			typeConversion.Type.Type = VariableTypeObject
			typeConversion.Type.Class = d.(*Class)
			typeConversion.Type.Pos = e.Pos
			if len(call.Args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
					e.Pos.errMsgPrefix()))
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
		case *Variable:
			v := d.(*Variable)
			if v.Type.Type != VariableTypeFunction {
				*errs = append(*errs, fmt.Errorf("%s '%s' is not a function , but '%s' ",
					call.Expression.Pos.errMsgPrefix(), v.Name, v.Type.TypeString()))
				return nil
			}
			call.Expression.Value = &Type{
				Pos:          e.Pos,
				Type:         VariableTypeFunction,
				FunctionType: v.Type.FunctionType,
			}
			identifier.Variable = v
			return e.checkFunctionPointerCall(block, errs, v.Type.FunctionType, call)
		default:
			*errs = append(*errs, fmt.Errorf("%s cannot make call on '%s'",
				call.Expression.Pos.errMsgPrefix(), block.identifierIsWhat(d)))
			return nil
		}
	}
	functionPointer, es := call.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if functionPointer == nil {
		return nil
	}
	if functionPointer.Type != VariableTypeFunction {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not a function , but '%s'",
			errMsgPrefix(e.Pos),
			call.Expression.Description, functionPointer.TypeString()))
		return nil
	}
	if call.Expression.Type == ExpressionTypeFunctionLiteral {
		/*
			fn() {

			}()
			no name function is statement too
		*/
		functionPointer.Function = call.Expression.Data.(*Function)
		call.Expression.IsStatementExpression = true
	}
	return e.checkFunctionPointerCall(block, errs, functionPointer.FunctionType, call)
}

func (e *Expression) checkFunctionPointerCall(block *Block, errs *[]error, ft *FunctionType, call *ExpressionFunctionCall) []*Type {
	callArgsTypes := checkExpressions(block, call.Args, errs, true)
	ret := ft.mkCallReturnTypes(e.Pos)
	var err error
	call.VArgs, err = ft.fitArgs(e.Pos, &call.Args, callArgsTypes, nil)
	if err != nil {
		*errs = append(*errs, err)
	}
	return ret
}

func (e *Expression) checkFunctionCall(block *Block, errs *[]error, f *Function, call *ExpressionFunctionCall) []*Type {
	callArgsTypes := checkExpressions(block, call.Args, errs, true)
	var tf *Function
	if f.TemplateFunction != nil {
		length := len(*errs)
		//rewrite
		tf = e.checkTemplateFunctionCall(block, errs, callArgsTypes, f)
		if len(*errs) != length { // if no
			return nil
		}
		ret := tf.Type.mkCallReturnTypes(e.Pos)
		var err error
		call.VArgs, err = tf.Type.fitArgs(e.Pos, &call.Args, callArgsTypes, tf)
		if err != nil {
			*errs = append(*errs, err)
		}
		return ret
	} else { // not template function
		if f.IsBuildIn {
			if f.LoadedFromCorePackage {
				var err error
				call.VArgs, err = f.Type.fitArgs(e.Pos, &call.Args, callArgsTypes, f)
				if err != nil {
					*errs = append(*errs, err)
				}
				return f.Type.mkCallReturnTypes(e.Pos)
			} else {
				length := len(*errs)
				f.buildInFunctionChecker(f, e.Data.(*ExpressionFunctionCall), block, errs, callArgsTypes, e.Pos)
				if len(*errs) == length {
					//special case ,avoid null pointer
					return f.Type.mkCallReturnTypes(e.Pos)
				}
				return nil //
			}
		} else {
			if len(call.ParameterTypes) > 0 {
				*errs = append(*errs, fmt.Errorf("%s function is not a template function",
					errMsgPrefix(e.Pos)))
			}
			ret := f.Type.mkCallReturnTypes(e.Pos)
			var err error
			call.VArgs, err = f.Type.fitArgs(e.Pos, &call.Args, callArgsTypes, f)
			if err != nil {
				*errs = append(*errs, err)
			}
			return ret
		}
	}
}

func (e *Expression) checkTemplateFunctionCall(block *Block, errs *[]error,
	argTypes []*Type, f *Function) (ret *Function) {
	call := e.Data.(*ExpressionFunctionCall)
	parameterTypes := make(map[string]*Type)
	parameterTypeArray := []*Type{}
	for k, v := range f.Type.ParameterList {
		if v == nil ||
			v.Type == nil ||
			len(v.Type.getParameterType()) == 0 {
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
		t := v.Type.Clone()
		t.bindWithParameterTypes(parameterTypes)
		parameterTypeArray = append(parameterTypeArray, t)
	}
	tps := call.ParameterTypes
	for k, v := range f.Type.ReturnList {
		if v == nil || v.Type == nil || len(v.Type.getParameterType()) == 0 {
			continue
		}
		if len(tps) == 0 || tps[0] == nil {
			//trying already have
			if err := v.Type.canBeBindWithParameterTypes(parameterTypes); err == nil {
				//very good no error
				t := v.Type.Clone()
				t.bindWithParameterTypes(parameterTypes)
				parameterTypeArray = append(parameterTypeArray, t)
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
		t := v.Type.Clone()
		t.bindWithParameterTypes(parameterTypes)
		parameterTypeArray = append(parameterTypeArray, t)
		tps = tps[1:]
	}
	call.TemplateFunctionCallPair = f.TemplateFunction.insert(parameterTypeArray)
	if call.TemplateFunctionCallPair.Function == nil { // not called before,make the binds
		cloneFunction, es := f.clone()
		if len(es) > 0 {
			*errs = append(*errs, es...)
			return nil
		}
		cloneFunction.TemplateFunction = nil
		call.TemplateFunctionCallPair.Function = cloneFunction
		cloneFunction.parameterTypes = parameterTypes
		for _, v := range cloneFunction.Type.ParameterList {
			if len(v.Type.getParameterType()) > 0 {
				v.Type = parameterTypeArray[0]
				parameterTypeArray = parameterTypeArray[1:]
			}
		}
		for _, v := range cloneFunction.Type.ReturnList {
			if len(v.Type.getParameterType()) > 0 {
				v.Type = parameterTypeArray[0]
				parameterTypeArray = parameterTypeArray[1:]
			}
		}
		//check this function
		cloneFunction.Block.inherit(&PackageBeenCompile.Block)
		if cloneFunction.Block.Functions == nil {
			cloneFunction.Block.Functions = make(map[string]*Function)
		}
		cloneFunction.Block.Functions[cloneFunction.Name] = cloneFunction
		cloneFunction.Block.InheritedAttribute.Function = cloneFunction
		cloneFunction.checkParametersAndReturns(errs, true, false)
		cloneFunction.checkBlock(errs)
	}
	ret = call.TemplateFunctionCallPair.Function
	// when all ok ,ret is not a template function any more
	if len(tps) > 0 {
		*errs = append(*errs, fmt.Errorf("%s to many parameter type to call template function",
			errMsgPrefix(e.Pos)))
	}
	return ret
}
