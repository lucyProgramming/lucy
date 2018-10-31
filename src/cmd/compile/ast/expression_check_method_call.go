package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func (e *Expression) checkMethodCallExpression(block *Block, errs *[]error) []*Type {
	call := e.Data.(*ExpressionMethodCall)
	object, es := call.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	// call father`s construction method
	if call.Name == SUPER && object.Type == VariableTypeObject {
		e.checkMethodCallExpressionOnSuper(block, errs, object)
		return []*Type{mkVoidType(e.Pos)}
	}
	switch object.Type {
	case VariableTypePackage:
		return e.checkMethodCallExpressionOnPackage(block, errs, object.Package)
	case VariableTypeMap:
		return e.checkMethodCallExpressionOnMap(block, errs, object.Map)
	case VariableTypeArray:
		return e.checkMethodCallExpressionOnArray(block, errs, object)
	case VariableTypeJavaArray:
		return e.checkMethodCallExpressionOnJavaArray(block, errs, object)
	case VariableTypeDynamicSelector:
		return e.checkMethodCallExpressionOnDynamicSelector(block, errs, object)
	case VariableTypeString:
		if err := loadJavaStringClass(e.Pos); err != nil {
			*errs = append(*errs, err)
			return nil
		}
		errsLength := len(*errs)
		args := checkExpressions(block, call.Args, errs, true)
		if len(*errs) > errsLength {
			return nil
		}
		ms, matched, err := javaStringClass.accessMethod(e.Pos, errs, call, args,
			false, nil)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if matched {
			call.Class = javaStringClass
			if false == call.Expression.IsIdentifier(ThisPointerName) &&
				ms[0].IsPublic() == false {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not public", e.Pos.ErrMsgPrefix(), call.Name))
			}
			call.Method = ms[0]
			return ms[0].Function.Type.mkCallReturnTypes(e.Pos)
		} else {
			*errs = append(*errs, methodsNotMatchError(e.Pos, call.Name, ms, args))
			return nil
		}

	case VariableTypeObject, VariableTypeClass:
		call.Class = object.Class
		errsLength := len(*errs)
		callArgTypes := checkExpressions(block, call.Args, errs, true)
		if len(*errs) > errsLength {
			return nil
		}
		if object.Class.IsInterface() {
			if object.Type == VariableTypeClass {
				*errs = append(*errs, fmt.Errorf("%s cannot make_node_objects call on interface '%s'",
					e.Pos.ErrMsgPrefix(), object.Class.Name))
				return nil
			}
			ms, matched, err :=
				object.Class.accessInterfaceObjectMethod(e.Pos, errs, call.Name, call, callArgTypes, false)
			if err != nil {
				*errs = append(*errs, err)
				return nil
			}
			if matched {
				if ms[0].IsStatic() {
					*errs = append(*errs, fmt.Errorf("%s method '%s' is static",
						e.Pos.ErrMsgPrefix(), call.Name))
				}
				call.Method = ms[0]
				return ms[0].Function.Type.mkCallReturnTypes(e.Pos)
			}
			*errs = append(*errs, methodsNotMatchError(e.Pos, call.Name, ms, callArgTypes))
			return nil
		}
		if len(call.ParameterTypes) > 0 {
			*errs = append(*errs, fmt.Errorf("%s method call expect no parameter types",
				errMsgPrefix(e.Pos)))
		}
		var fieldMethodHandler *ClassField
		ms, matched, err := object.Class.accessMethod(e.Pos, errs, call, callArgTypes,
			false, &fieldMethodHandler)
		if err != nil {
			*errs = append(*errs, err)
			if len(ms) > 0 {
				return ms[0].Function.Type.mkCallReturnTypes(e.Pos)
			}
			return nil
		}
		if fieldMethodHandler != nil {
			err := call.Expression.fieldAccessAble(block, fieldMethodHandler)
			if err != nil {
				*errs = append(*errs, err)
			}
			call.FieldMethodHandler = fieldMethodHandler
			return fieldMethodHandler.Type.FunctionType.mkCallReturnTypes(e.Pos)
		}
		if matched {
			m := ms[0]
			err := call.Expression.methodAccessAble(block, m)
			if err != nil {
				*errs = append(*errs, err)
			}
			call.Method = m
			return m.Function.Type.mkCallReturnTypes(e.Pos)
		}
		*errs = append(*errs, methodsNotMatchError(e.Pos, call.Name, ms, callArgTypes))
		return nil
	default:
		*errs = append(*errs, fmt.Errorf("%s cannot make_node_objects method call '%s' on '%s'",
			e.Pos.ErrMsgPrefix(), call.Name, object.TypeString()))
		return nil
	}
}

/*
	this.super()
*/
func (e *Expression) checkMethodCallExpressionOnSuper(block *Block, errs *[]error, object *Type) {
	call := e.Data.(*ExpressionMethodCall)
	if call.Expression.IsIdentifier(ThisPointerName) == false {
		*errs = append(*errs, fmt.Errorf("%s call father`s constuction must use 'thi.super()'",
			e.Pos.ErrMsgPrefix()))
		return
	}
	if block.InheritedAttribute.IsConstructionMethod == false ||
		block.IsFunctionBlock == false ||
		block.InheritedAttribute.StatementOffset != 0 {
		*errs = append(*errs,
			fmt.Errorf("%s call father`s constuction on must first statement of a constructon method",
				e.Pos.ErrMsgPrefix()))
		return
	}
	if object.Class.LoadFromOutSide {
		err := object.Class.loadSuperClass(e.Pos)
		if err != nil {
			*errs = append(*errs, err)
			return
		}
		if object.Class.SuperClass == nil {
			return
		}
	} else {
		if object.Class.SuperClass == nil {
			return
		}
	}
	errsLength := len(*errs)
	callArgsTypes := checkExpressions(block, call.Args, errs, true)
	if len(*errs) > errsLength {
		return
	}
	ms, matched, err := object.Class.SuperClass.accessConstructionFunction(e.Pos, errs,
		nil, call, callArgsTypes)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", e.Pos.ErrMsgPrefix(), err))
		return
	}
	if matched {
		m := ms[0]
		if err := object.Class.SuperClass.constructionMethodAccessAble(e.Pos, m); err != nil {
			*errs = append(*errs, err)
		}
		call.Name = "<init>"
		call.Method = m
		call.Class = object.Class.SuperClass
		block.Statements[0].IsCallFatherConstructionStatement = true
		block.InheritedAttribute.Function.CallFatherConstructionExpression = e
		return
	}
	*errs = append(*errs, methodsNotMatchError(e.Pos, object.TypeString(), ms, callArgsTypes))
}

func (e *Expression) checkMethodCallExpressionOnDynamicSelector(block *Block, errs *[]error, object *Type) []*Type {
	call := e.Data.(*ExpressionMethodCall)
	if call.Name == SUPER {
		*errs = append(*errs, fmt.Errorf("%s access '%s' at '%s' not allow",
			e.Pos.ErrMsgPrefix(), SUPER, object.TypeString()))
		return nil
	}
	var fieldMethodHandler *ClassField
	errsLength := len(*errs)
	callArgTypes := checkExpressions(block, call.Args, errs, true)
	if len(*errs) > errsLength {
		return nil
	}
	ms, matched, err := object.Class.accessMethod(e.Pos, errs, call, callArgTypes, false, &fieldMethodHandler)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if matched {
		if fieldMethodHandler != nil {
			call.FieldMethodHandler = fieldMethodHandler
			return fieldMethodHandler.Type.FunctionType.mkCallReturnTypes(e.Pos)
		} else {
			method := ms[0]
			call.Method = method
			return method.Function.Type.mkCallReturnTypes(e.Pos)
		}
	} else {
		*errs = append(*errs, methodsNotMatchError(e.Pos, call.Name, ms, callArgTypes))
	}
	return nil
}
func (e *Expression) checkMethodCallExpressionOnJavaArray(block *Block, errs *[]error, array *Type) []*Type {
	call := e.Data.(*ExpressionMethodCall)
	switch call.Name {
	case common.ArrayMethodSize:
		result := &Type{}
		result.Type = VariableTypeInt
		result.Pos = e.Pos
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s method '%s' expect no arguments",
				call.Args[0].Pos.ErrMsgPrefix(), call.Name))
		}
		return []*Type{result}
	default:
		*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on '%s'",
			e.Pos.ErrMsgPrefix(), call.Name, array.TypeString()))
	}
	return nil
}

func (e *Expression) checkMethodCallExpressionOnPackage(block *Block, errs *[]error, p *Package) []*Type {
	call := e.Data.(*ExpressionMethodCall)
	d, exists := p.Block.NameExists(call.Name)
	if exists == false {
		*errs = append(*errs, fmt.Errorf("%s function '%s' not found", e.Pos.ErrMsgPrefix(), call.Name))
		return nil
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		if f.IsPublic() == false &&
			p.isSame(&PackageBeenCompile) == false {
			*errs = append(*errs, fmt.Errorf("%s function '%s' is not public",
				e.Pos.ErrMsgPrefix(), call.Name))
		}
		if f.TemplateFunction != nil {
			// better convert to function call
			methodCall := e.Data.(*ExpressionMethodCall)
			functionCall := &ExpressionFunctionCall{}
			functionCall.Args = methodCall.Args
			functionCall.Function = f
			functionCall.ParameterTypes = methodCall.ParameterTypes
			e.Type = ExpressionTypeFunctionCall
			e.Data = functionCall
			return e.checkFunctionCall(block, errs, f, functionCall)
		} else {
			methodCall := e.Data.(*ExpressionMethodCall)
			methodCall.PackageFunction = f
			ret := f.Type.mkCallReturnTypes(e.Pos)
			errsLength := len(*errs)
			callArgsTypes := checkExpressions(block, methodCall.Args, errs, true)
			if len(*errs) > errsLength {
				return ret
			}
			var err error
			methodCall.VArgs, err = f.Type.fitArgs(e.Pos, &call.Args, callArgsTypes, f)
			if err != nil {
				*errs = append(*errs, err)
			}
			return ret
		}
	case *Variable:
		v := d.(*Variable)
		if v.isPublic() == false && p.isSame(&PackageBeenCompile) == false {
			*errs = append(*errs, fmt.Errorf("%s variable '%s' is not public",
				e.Pos.ErrMsgPrefix(), call.Name))
		}
		if v.Type.Type != VariableTypeFunction {
			*errs = append(*errs, fmt.Errorf("%s variable '%s' is not a function",
				e.Pos.ErrMsgPrefix(), call.Name))
			return nil
		}
		call := e.Data.(*ExpressionMethodCall)
		if len(call.ParameterTypes) > 0 {
			*errs = append(*errs, fmt.Errorf("%s variable '%s' cannot be a template fucntion",
				errMsgPrefix(call.ParameterTypes[0].Pos), call.Name))
		}
		ret := v.Type.FunctionType.mkCallReturnTypes(e.Pos)
		errsLength := len(*errs)
		callArgsTypes := checkExpressions(block, call.Args, errs, true)
		if len(*errs) > errsLength {
			return ret
		}
		vArgs, err := v.Type.FunctionType.fitArgs(e.Pos, &call.Args, callArgsTypes, nil)
		if err != nil {
			*errs = append(*errs, err)
		}

		call.PackageGlobalVariableFunction = v
		call.VArgs = vArgs
		return ret
	case *Class:
		//object cast
		class := d.(*Class)
		if class.IsPublic() == false && p.isSame(&PackageBeenCompile) == false {
			*errs = append(*errs, fmt.Errorf("%s class '%s' is not public",
				e.Pos.ErrMsgPrefix(), call.Name))
		}
		conversion := &ExpressionTypeConversion{}
		conversion.Type = &Type{}
		conversion.Type.Type = VariableTypeObject
		conversion.Type.Pos = e.Pos
		conversion.Type.Class = class
		e.Type = ExpressionTypeCheckCast
		if len(call.Args) >= 1 {
			conversion.Expression = call.Args[0]
		}
		e.Data = conversion
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument", e.Pos.ErrMsgPrefix()))
			return []*Type{conversion.Type.Clone()}
		}
		return []*Type{e.checkTypeConversionExpression(block, errs)}
	case *Type:
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
				e.Pos.ErrMsgPrefix()))
			result := p.Block.TypeAliases[call.Name].Clone()
			result.Pos = e.Pos
			return []*Type{result}
		}
		conversion := &ExpressionTypeConversion{}
		conversion.Type = p.Block.TypeAliases[call.Name]
		e.Type = ExpressionTypeCheckCast
		if len(call.Args) >= 1 {
			conversion.Expression = call.Args[0]
		}
		e.Data = conversion
		return []*Type{e.checkTypeConversionExpression(block, errs)}
	default:
		*errs = append(*errs, fmt.Errorf("%s '%s' is not a function",
			e.Pos.ErrMsgPrefix(), call.Name))
		return nil
	}
}
func (e *Expression) checkMethodCallExpressionOnArray(block *Block, errs *[]error, array *Type) []*Type {
	call := e.Data.(*ExpressionMethodCall)
	switch call.Name {
	case common.ArrayMethodSize,
		common.ArrayMethodCap,
		common.ArrayMethodStart,
		common.ArrayMethodEnd:
		result := &Type{}
		result.Type = VariableTypeInt
		result.Pos = e.Pos
		if len(call.Args) > 0 {
			*errs = append(*errs,
				fmt.Errorf("%s too mamy argument to call,method '%s' expect no arguments",
					call.Args[0].Pos.ErrMsgPrefix(), call.Name))
		}
		return []*Type{result}
	case common.ArrayMethodAppend,
		common.ArrayMethodAppendAll:
		if len(call.Args) == 0 {
			*errs = append(*errs,
				fmt.Errorf("%s too few arguments to call %s,expect at least one argument",
					e.Pos.ErrMsgPrefix(), call.Name))
		}
		ts := checkExpressions(block, call.Args, errs, true)
		for _, t := range ts {
			if t == nil {
				continue
			}
			if call.Name == common.ArrayMethodAppend {
				if array.Array.assignAble(errs, t) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' to call method '%s'",
						t.Pos.ErrMsgPrefix(), t.TypeString(), array.Array.TypeString(), call.Name))
				}
			} else {
				if array.assignAble(errs, t) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' to call method '%s'",
						t.Pos.ErrMsgPrefix(), t.TypeString(), array.TypeString(), call.Name))
				}
			}
		}
		result := &Type{}
		result.Type = VariableTypeVoid
		result.Pos = e.Pos
		return []*Type{result}
	case common.ArrayMethodGetUnderlyingArray:
		result := &Type{}
		result.Type = VariableTypeJavaArray
		result.Pos = e.Pos
		result.Array = array.Array.Clone()
		result.Array.Pos = e.Pos
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s too mamy argument to call,method '%s' expect no arguments",
				call.Args[0].Pos.ErrMsgPrefix(), call.Name))
		}
		return []*Type{result}
	default:
		*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on array", e.Pos.ErrMsgPrefix(), call.Name))
	}
	return nil
}
func (e *Expression) checkMethodCallExpressionOnMap(block *Block, errs *[]error, m *Map) []*Type {
	call := e.Data.(*ExpressionMethodCall)
	switch call.Name {
	case common.MapMethodKeyExist:
		ret := &Type{}
		ret.Pos = e.Pos
		ret.Type = VariableTypeBool
		if len(call.Args) != 1 {
			pos := e.Pos
			if len(call.Args) != 0 {
				pos = call.Args[1].Pos
			}
			*errs = append(*errs, fmt.Errorf("%s call '%s' expect one argument",
				pos.ErrMsgPrefix(), call.Name))
			return []*Type{ret}
		}
		t, es := call.Args[0].checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if t == nil {
			return []*Type{ret}
		}
		if false == m.K.assignAble(errs, t) {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
				t.Pos.ErrMsgPrefix(), t.TypeString(), m.K.TypeString()))
		}
		return []*Type{ret}
	case common.MapMethodRemove:
		ret := &Type{}
		ret.Pos = e.Pos
		ret.Type = VariableTypeVoid
		if len(call.Args) == 0 {
			*errs = append(*errs, fmt.Errorf("%s remove expect at last 1 argement",
				e.Pos.ErrMsgPrefix()))
			return []*Type{ret}
		}
		ts := checkExpressions(block, call.Args, errs, true)
		for _, t := range ts {
			if t == nil {
				continue
			}
			if m.K.assignAble(errs, t) == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' for map-key",
					t.Pos.ErrMsgPrefix(), t.TypeString(), m.K.TypeString()))
			}
		}
		return []*Type{ret}
	case common.MapMethodRemoveAll:
		ret := &Type{}
		ret.Pos = e.Pos
		ret.Type = VariableTypeVoid
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s '%s' expect no arguments",
				e.Pos.ErrMsgPrefix(), common.MapMethodRemoveAll))
		}
		return []*Type{ret}
	case common.MapMethodSize:
		ret := &Type{}
		ret.Pos = e.Pos
		ret.Type = VariableTypeInt
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s too many argument to call '%s''",
				call.Args[0].Pos.ErrMsgPrefix(), call.Name))
		}
		return []*Type{ret}
	default:
		*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on map",
			e.Pos.ErrMsgPrefix(), call.Name))
		return nil
	}
	return nil
}
