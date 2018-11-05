package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func (this *Expression) checkMethodCallExpression(block *Block, errs *[]error) []*Type {
	call := this.Data.(*ExpressionMethodCall)
	object, es := call.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	// call father`s construction method
	if call.Name == SUPER && object.Type == VariableTypeObject {
		this.checkMethodCallExpressionOnSuper(block, errs, object)
		return []*Type{mkVoidType(this.Pos)}
	}
	switch object.Type {
	case VariableTypePackage:
		return this.checkMethodCallExpressionOnPackage(block, errs, object.Package)
	case VariableTypeMap:
		return this.checkMethodCallExpressionOnMap(block, errs, object.Map)
	case VariableTypeArray:
		return this.checkMethodCallExpressionOnArray(block, errs, object)
	case VariableTypeJavaArray:
		return this.checkMethodCallExpressionOnJavaArray(block, errs, object)
	case VariableTypeDynamicSelector:
		if call.Name == "finalize" {
			*errs = append(*errs, fmt.Errorf("%s cannot call '%s'", this.Pos.ErrMsgPrefix(), call.Name))
			return nil
		}
		return this.checkMethodCallExpressionOnDynamicSelector(block, errs, object)
	case VariableTypeString:
		if call.Name == "finalize" {
			*errs = append(*errs, fmt.Errorf("%s cannot call '%s'", this.Pos.ErrMsgPrefix(), call.Name))
			return nil
		}
		if err := loadJavaStringClass(this.Pos); err != nil {
			*errs = append(*errs, err)
			return nil
		}
		errsLength := len(*errs)
		args := checkExpressions(block, call.Args, errs, true)
		if len(*errs) > errsLength {
			return nil
		}
		ms, matched, err := javaStringClass.accessMethod(this.Pos, errs, call, args,
			false, nil)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if matched {
			call.Class = javaStringClass
			if false == call.Expression.IsIdentifier(ThisPointerName) &&
				ms[0].IsPublic() == false {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not public", this.Pos.ErrMsgPrefix(), call.Name))
			}
			call.Method = ms[0]
			return ms[0].Function.Type.mkCallReturnTypes(this.Pos)
		} else {
			*errs = append(*errs, methodsNotMatchError(this.Pos, call.Name, ms, args))
			return nil
		}

	case VariableTypeObject, VariableTypeClass:
		if call.Name == "finalize" {
			*errs = append(*errs, fmt.Errorf("%s cannot call '%s'", this.Pos.ErrMsgPrefix(), call.Name))
			return nil
		}
		call.Class = object.Class
		errsLength := len(*errs)
		callArgTypes := checkExpressions(block, call.Args, errs, true)
		if len(*errs) > errsLength {
			return nil
		}
		if object.Class.IsInterface() {
			if object.Type == VariableTypeClass {
				*errs = append(*errs, fmt.Errorf("%s cannot make_node_objects call on interface '%s'",
					this.Pos.ErrMsgPrefix(), object.Class.Name))
				return nil
			}
			ms, matched, err :=
				object.Class.accessInterfaceObjectMethod(this.Pos, errs, call.Name, call, callArgTypes, false)
			if err != nil {
				*errs = append(*errs, err)
				return nil
			}
			if matched {
				if ms[0].IsStatic() {
					*errs = append(*errs, fmt.Errorf("%s method '%s' is static",
						this.Pos.ErrMsgPrefix(), call.Name))
				}
				call.Method = ms[0]
				return ms[0].Function.Type.mkCallReturnTypes(this.Pos)
			}
			*errs = append(*errs, methodsNotMatchError(this.Pos, call.Name, ms, callArgTypes))
			return nil
		}
		if len(call.ParameterTypes) > 0 {
			*errs = append(*errs, fmt.Errorf("%s method call expect no parameter types",
				errMsgPrefix(this.Pos)))
		}
		var fieldMethodHandler *ClassField
		ms, matched, err := object.Class.accessMethod(this.Pos, errs, call, callArgTypes,
			false, &fieldMethodHandler)
		if err != nil {
			*errs = append(*errs, err)
			if len(ms) > 0 {
				return ms[0].Function.Type.mkCallReturnTypes(this.Pos)
			}
			return nil
		}
		if fieldMethodHandler != nil {
			err := call.Expression.fieldAccessAble(block, fieldMethodHandler)
			if err != nil {
				*errs = append(*errs, err)
			}
			call.FieldMethodHandler = fieldMethodHandler
			return fieldMethodHandler.Type.FunctionType.mkCallReturnTypes(this.Pos)
		}
		if matched {
			m := ms[0]
			err := call.Expression.methodAccessAble(block, m)
			if err != nil {
				*errs = append(*errs, err)
			}
			call.Method = m
			return m.Function.Type.mkCallReturnTypes(this.Pos)
		}
		*errs = append(*errs, methodsNotMatchError(this.Pos, call.Name, ms, callArgTypes))
		return nil
	default:
		*errs = append(*errs, fmt.Errorf("%s cannot make_node_objects method call '%s' on '%s'",
			this.Pos.ErrMsgPrefix(), call.Name, object.TypeString()))
		return nil
	}
}

/*
	this.super()
*/
func (this *Expression) checkMethodCallExpressionOnSuper(
	block *Block,
	errs *[]error,
	object *Type) {
	call := this.Data.(*ExpressionMethodCall)
	if call.Expression.IsIdentifier(ThisPointerName) == false {
		*errs = append(*errs, fmt.Errorf("%s call father`s constuction must use 'thi.super()'",
			this.Pos.ErrMsgPrefix()))
		return
	}
	if block.InheritedAttribute.IsConstructionMethod == false ||
		block.IsFunctionBlock == false ||
		block.InheritedAttribute.StatementOffset != 0 {
		*errs = append(*errs,
			fmt.Errorf("%s call father`s constuction on must first statement of a constructon method",
				this.Pos.ErrMsgPrefix()))
		return
	}
	if object.Class.LoadFromOutSide {
		err := object.Class.loadSuperClass(this.Pos)
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
	ms, matched, err := object.Class.SuperClass.accessConstructionMethod(this.Pos, errs,
		nil, call, callArgsTypes)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", this.Pos.ErrMsgPrefix(), err))
		return
	}
	if matched {
		m := ms[0]
		if err := object.Class.SuperClass.constructionMethodAccessAble(this.Pos, m); err != nil {
			*errs = append(*errs, err)
		}
		call.Name = "<init>"
		call.Method = m
		call.Class = object.Class.SuperClass
		block.Statements[0].IsCallFatherConstructionStatement = true
		block.InheritedAttribute.Function.CallFatherConstructionExpression = this
		return
	}
	*errs = append(*errs, methodsNotMatchError(this.Pos, object.TypeString(), ms, callArgsTypes))
}

func (this *Expression) checkMethodCallExpressionOnDynamicSelector(block *Block, errs *[]error, object *Type) []*Type {
	call := this.Data.(*ExpressionMethodCall)
	if call.Name == SUPER {
		*errs = append(*errs, fmt.Errorf("%s access '%s' at '%s' not allow",
			this.Pos.ErrMsgPrefix(), SUPER, object.TypeString()))
		return nil
	}
	var fieldMethodHandler *ClassField
	errsLength := len(*errs)
	callArgTypes := checkExpressions(block, call.Args, errs, true)
	if len(*errs) > errsLength {
		return nil
	}
	ms, matched, err := object.Class.accessMethod(this.Pos, errs, call, callArgTypes, false, &fieldMethodHandler)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if matched {
		if fieldMethodHandler != nil {
			call.FieldMethodHandler = fieldMethodHandler
			return fieldMethodHandler.Type.FunctionType.mkCallReturnTypes(this.Pos)
		} else {
			method := ms[0]
			call.Method = method
			return method.Function.Type.mkCallReturnTypes(this.Pos)
		}
	} else {
		*errs = append(*errs, methodsNotMatchError(this.Pos, call.Name, ms, callArgTypes))
	}
	return nil
}
func (this *Expression) checkMethodCallExpressionOnJavaArray(block *Block, errs *[]error, array *Type) []*Type {
	call := this.Data.(*ExpressionMethodCall)
	switch call.Name {
	case common.ArrayMethodSize:
		result := &Type{}
		result.Type = VariableTypeInt
		result.Pos = this.Pos
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s method '%s' expect no arguments",
				call.Args[0].Pos.ErrMsgPrefix(), call.Name))
		}
		return []*Type{result}
	default:
		*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on '%s'",
			this.Pos.ErrMsgPrefix(), call.Name, array.TypeString()))
	}
	return nil
}

func (this *Expression) checkMethodCallExpressionOnPackage(
	block *Block,
	errs *[]error,
	p *Package) []*Type {
	call := this.Data.(*ExpressionMethodCall)
	d, exists := p.Block.NameExists(call.Name)
	if exists == false {
		*errs = append(*errs, fmt.Errorf("%s function '%s' not found", this.Pos.ErrMsgPrefix(), call.Name))
		return nil
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		if f.IsPublic() == false &&
			p.isSame(&PackageBeenCompile) == false {
			*errs = append(*errs, fmt.Errorf("%s function '%s' is not public",
				this.Pos.ErrMsgPrefix(), call.Name))
		}
		if f.TemplateFunction != nil {
			// better convert to function call
			methodCall := this.Data.(*ExpressionMethodCall)
			functionCall := &ExpressionFunctionCall{}
			functionCall.Args = methodCall.Args
			functionCall.Function = f
			functionCall.ParameterTypes = methodCall.ParameterTypes
			this.Type = ExpressionTypeFunctionCall
			this.Data = functionCall
			return this.checkFunctionCall(block, errs, f, functionCall)
		} else {
			methodCall := this.Data.(*ExpressionMethodCall)
			methodCall.PackageFunction = f
			ret := f.Type.mkCallReturnTypes(this.Pos)
			errsLength := len(*errs)
			callArgsTypes := checkExpressions(block, methodCall.Args, errs, true)
			if len(*errs) > errsLength {
				return ret
			}
			var err error
			methodCall.VArgs, err = f.Type.fitArgs(this.Pos, &call.Args, callArgsTypes, f)
			if err != nil {
				*errs = append(*errs, err)
			}
			return ret
		}
	case *Variable:
		v := d.(*Variable)
		if v.isPublic() == false && p.isSame(&PackageBeenCompile) == false {
			*errs = append(*errs, fmt.Errorf("%s variable '%s' is not public",
				this.Pos.ErrMsgPrefix(), call.Name))
		}
		if v.Type.Type != VariableTypeFunction {
			*errs = append(*errs, fmt.Errorf("%s variable '%s' is not a function",
				this.Pos.ErrMsgPrefix(), call.Name))
			return nil
		}
		call := this.Data.(*ExpressionMethodCall)
		if len(call.ParameterTypes) > 0 {
			*errs = append(*errs, fmt.Errorf("%s variable '%s' cannot be a template fucntion",
				errMsgPrefix(call.ParameterTypes[0].Pos), call.Name))
		}
		ret := v.Type.FunctionType.mkCallReturnTypes(this.Pos)
		errsLength := len(*errs)
		callArgsTypes := checkExpressions(block, call.Args, errs, true)
		if len(*errs) > errsLength {
			return ret
		}
		vArgs, err := v.Type.FunctionType.fitArgs(this.Pos, &call.Args, callArgsTypes, nil)
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
				this.Pos.ErrMsgPrefix(), call.Name))
		}
		conversion := &ExpressionTypeConversion{}
		conversion.Type = &Type{}
		conversion.Type.Type = VariableTypeObject
		conversion.Type.Pos = this.Pos
		conversion.Type.Class = class
		this.Type = ExpressionTypeCheckCast
		if len(call.Args) >= 1 {
			conversion.Expression = call.Args[0]
		}
		this.Data = conversion
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument", this.Pos.ErrMsgPrefix()))
			return []*Type{conversion.Type.Clone()}
		}
		return []*Type{this.checkTypeConversionExpression(block, errs)}
	case *Type:
		if len(call.Args) != 1 {
			*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
				this.Pos.ErrMsgPrefix()))
			result := p.Block.TypeAliases[call.Name].Clone()
			result.Pos = this.Pos
			return []*Type{result}
		}
		conversion := &ExpressionTypeConversion{}
		conversion.Type = p.Block.TypeAliases[call.Name]
		this.Type = ExpressionTypeCheckCast
		if len(call.Args) >= 1 {
			conversion.Expression = call.Args[0]
		}
		this.Data = conversion
		return []*Type{this.checkTypeConversionExpression(block, errs)}
	default:
		*errs = append(*errs, fmt.Errorf("%s '%s' is not a function",
			this.Pos.ErrMsgPrefix(), call.Name))
		return nil
	}
}
func (this *Expression) checkMethodCallExpressionOnArray(
	block *Block,
	errs *[]error,
	array *Type) []*Type {
	call := this.Data.(*ExpressionMethodCall)
	switch call.Name {
	case common.ArrayMethodSize,
		common.ArrayMethodCap,
		common.ArrayMethodStart,
		common.ArrayMethodEnd:
		result := &Type{}
		result.Type = VariableTypeInt
		result.Pos = this.Pos
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
					this.Pos.ErrMsgPrefix(), call.Name))
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
		result.Pos = this.Pos
		return []*Type{result}
	case common.ArrayMethodGetUnderlyingArray:
		result := &Type{}
		result.Type = VariableTypeJavaArray
		result.Pos = this.Pos
		result.Array = array.Array.Clone()
		result.Array.Pos = this.Pos
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s too mamy argument to call,method '%s' expect no arguments",
				call.Args[0].Pos.ErrMsgPrefix(), call.Name))
		}
		return []*Type{result}
	default:
		*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on array", this.Pos.ErrMsgPrefix(), call.Name))
	}
	return nil
}
func (this *Expression) checkMethodCallExpressionOnMap(
	block *Block,
	errs *[]error,
	m *Map) []*Type {
	call := this.Data.(*ExpressionMethodCall)
	switch call.Name {
	case common.MapMethodKeyExist:
		ret := &Type{}
		ret.Pos = this.Pos
		ret.Type = VariableTypeBool
		if len(call.Args) != 1 {
			pos := this.Pos
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
		ret.Pos = this.Pos
		ret.Type = VariableTypeVoid
		if len(call.Args) == 0 {
			*errs = append(*errs, fmt.Errorf("%s remove expect at last 1 argement",
				this.Pos.ErrMsgPrefix()))
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
		ret.Pos = this.Pos
		ret.Type = VariableTypeVoid
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s '%s' expect no arguments",
				this.Pos.ErrMsgPrefix(), common.MapMethodRemoveAll))
		}
		return []*Type{ret}
	case common.MapMethodSize:
		ret := &Type{}
		ret.Pos = this.Pos
		ret.Type = VariableTypeInt
		if len(call.Args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s too many argument to call '%s''",
				call.Args[0].Pos.ErrMsgPrefix(), call.Name))
		}
		return []*Type{ret}
	default:
		*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on map",
			this.Pos.ErrMsgPrefix(), call.Name))
		return nil
	}
	return nil
}
