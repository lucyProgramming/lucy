package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkMethodCallExpression(block *Block, errs *[]error) []*Type {
	call := e.Data.(*ExpressionMethodCall)
	object, es := call.Expression.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if object == nil {
		return nil
	}
	if object.Type == VariableTypePackage {
		d, exists := object.Package.Block.NameExists(call.Name)
		if exists == false {
			*errs = append(*errs, fmt.Errorf("%s function '%s' not found", errMsgPrefix(e.Pos), call.Name))
			return nil
		}
		switch d.(type) {
		case *Function:
			f := d.(*Function)
			if f.isPublic() == false && object.Package.Name != PackageBeenCompile.Name {
				*errs = append(*errs, fmt.Errorf("%s function '%s' is not public",
					errMsgPrefix(e.Pos), call.Name))
			}
			call := e.Data.(*ExpressionMethodCall)
			callArgsTypes := checkRightValuesValid(checkExpressions(block, call.Args, errs), errs)
			_, vArgs, es := f.Type.fitCallArgs(e.Pos, &call.Args, callArgsTypes, f)
			ret := f.Type.getReturnTypes(e.Pos)
			if esNotEmpty(es) {
				*errs = append(*errs, es...)
			}
			call.PackageFunction = f
			call.VArgs = vArgs
			return ret
		case *Class:
			//object cast
			class := d.(*Class)
			if class.IsPublic() == false && object.Package.Name != PackageBeenCompile.Name {
				*errs = append(*errs, fmt.Errorf("%s class '%s' is not public",
					errMsgPrefix(e.Pos), call.Name))
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
				*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument", errMsgPrefix(e.Pos)))
				return []*Type{conversion.Type.Clone()}
			}
			return []*Type{e.checkTypeConversionExpression(block, errs)}
		case *Type:
			conversion := &ExpressionTypeConversion{}
			conversion.Type = object.Package.Block.TypeAliases[call.Name]
			e.Type = ExpressionTypeCheckCast
			if len(call.Args) >= 1 {
				conversion.Expression = call.Args[0]
			}
			e.Data = conversion
			if len(call.Args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
					errMsgPrefix(e.Pos)))
				return []*Type{conversion.Type}
			}
			return []*Type{e.checkTypeConversionExpression(block, errs)}
		case *Variable:
			v := d.(*Variable)
			if (v.AccessFlags&cg.ACC_FIELD_PUBLIC) == 0 && object.Package.Name != PackageBeenCompile.Name {
				*errs = append(*errs, fmt.Errorf("%s variable '%s' is not public",
					errMsgPrefix(e.Pos), call.Name))
			}
			if v.Type.Type != VariableTypeFunction {
				*errs = append(*errs, fmt.Errorf("%s variable '%s' is not a function",
					errMsgPrefix(e.Pos), call.Name))
				return nil
			}
			call := e.Data.(*ExpressionMethodCall)
			callArgsTypes := checkRightValuesValid(checkExpressions(block, call.Args, errs), errs)
			_, vArgs, es := v.Type.FunctionType.fitCallArgs(e.Pos, &call.Args, callArgsTypes, nil)
			ret := v.Type.FunctionType.getReturnTypes(e.Pos)
			if esNotEmpty(es) {
				*errs = append(*errs, es...)
			}
			call.PackageGlobalVariableFunctionHandler = v
			call.VArgs = vArgs
			return ret
		default:
			*errs = append(*errs, fmt.Errorf("%s '%s' is not a function",
				errMsgPrefix(e.Pos), call.Name))
			return nil
		}
	}
	if object.Type == VariableTypeMap {
		switch call.Name {
		case common.MapMethodKeyExists:
			ret := &Type{}
			ret.Pos = e.Pos
			ret.Type = VariableTypeBool
			if len(call.Args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s call '%s' expect one argument",
					errMsgPrefix(e.Pos), call.Name))
				return []*Type{ret}
			}
			matchKey := call.Name == common.MapMethodKeyExists
			t, es := call.Args[0].checkSingleValueContextExpression(block)
			if esNotEmpty(es) {
				*errs = append(*errs, es...)
			}
			if t == nil {
				return []*Type{ret}
			}
			if matchKey {
				if false == object.Map.K.Equal(errs, t) {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(e.Pos), t.TypeString(), object.Map.K.TypeString()))
				}
			} else {
				if false == object.Map.V.Equal(errs, t) {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(e.Pos), t.TypeString(), object.Map.V.TypeString()))
				}
			}
			return []*Type{ret}
		case common.MapMethodRemove:
			ret := &Type{}
			ret.Pos = e.Pos
			ret.Type = VariableTypeVoid
			if len(call.Args) == 0 {
				*errs = append(*errs, fmt.Errorf("%s remove expect at last 1 argement",
					errMsgPrefix(e.Pos)))
			}
			for _, v := range call.Args {
				ts, es := v.check(block)
				if esNotEmpty(es) {
					*errs = append(*errs, es...)
				}
				for _, t := range ts {
					if t == nil {
						continue
					}
					if object.Map.K.Equal(errs, t) == false {
						*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' for key",
							errMsgPrefix(e.Pos), t.TypeString(), object.Map.K.TypeString()))
					}
				}
			}
			return []*Type{ret}
		case common.MapMethodRemoveAll:
			ret := &Type{}
			ret.Pos = e.Pos
			ret.Type = VariableTypeVoid
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s removeAll expect no arguments",
					errMsgPrefix(e.Pos)))
			}
			return []*Type{ret}
		case common.MapMethodSize:
			ret := &Type{}
			ret.Pos = e.Pos
			ret.Type = VariableTypeInt
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s too many argument to call '%s''",
					errMsgPrefix(e.Pos), call.Name))
			}
			return []*Type{ret}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on map", errMsgPrefix(e.Pos), call.Name))
			return nil
		}
		return nil
	}
	if object.Type == VariableTypeJavaArray {
		switch call.Name {
		case common.ArrayMethodSize:
			t := &Type{}
			t.Type = VariableTypeInt
			t.Pos = e.Pos
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s method '%s' expect no arguments",
					errMsgPrefix(e.Pos), call.Name))
			}
			return []*Type{t}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on '%s'",
				errMsgPrefix(e.Pos), call.Name, object.TypeString()))
		}
		return nil
	}

	if object.Type == VariableTypeArray {
		switch call.Name {
		case common.ArrayMethodSize,
			common.ArrayMethodCap,   //for debug,remove when time is right
			common.ArrayMethodStart, //for debug,remove when time is right
			common.ArrayMethodEnd:   //for debug,remove when time is right
			t := &Type{}
			t.Type = VariableTypeInt
			t.Pos = e.Pos
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s too mamy argument to call,method '%s' expect no arguments",
					errMsgPrefix(e.Pos), call.Name))
			}
			return []*Type{t}
		case common.ArrayMethodAppend, common.ArrayMethodAppendAll:
			if len(call.Args) == 0 {
				*errs = append(*errs, fmt.Errorf("%s too few arguments to call %s,expect at least one argument",
					errMsgPrefix(e.Pos), call.Name))
			}
			for _, e := range call.Args {
				ts, es := e.check(block)
				if esNotEmpty(es) {
					*errs = append(*errs, es...)
				}
				for _, t := range ts {
					if call.Name == common.ArrayMethodAppend {
						if object.Array.Equal(errs, t) == false {
							*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' to call method '%s'",
								errMsgPrefix(t.Pos), t.TypeString(), object.Array.TypeString(), call.Name))
						}
					} else {
						if object.Equal(errs, t) == false {
							*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' to call method '%s'",
								errMsgPrefix(t.Pos), t.TypeString(), object.Array.TypeString(), call.Name))
						}
					}
				}
			}
			t := object.Clone() // support
			t.Pos = e.Pos
			return []*Type{t}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on array", errMsgPrefix(e.Pos), call.Name))
		}
		return nil
	}

	if object.Type == VariableTypeString {
		if err := loadJavaStringClass(e.Pos); err != nil {
			*errs = append(*errs, err)
			return nil
		}
		//var fieldMethodHandler *ClassField
		args := checkRightValuesValid(checkExpressions(block, call.Args, errs), errs)
		ms, matched, err := javaStringClass.accessMethod(e.Pos, errs, call.Name, call, args, false, nil)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if matched {
			call.Class = javaStringClass
			if false == call.Expression.isThis() &&
				ms[0].IsPublic() == false {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not public", errMsgPrefix(e.Pos), call.Name))
			}
			call.Method = ms[0]
			return ms[0].Function.Type.getReturnTypes(e.Pos)
		}
		if len(ms) == 0 {
			*errs = append(*errs, fmt.Errorf("%s method '%s' not found", errMsgPrefix(e.Pos), call.Name))
		} else {
			*errs = append(*errs, msNotMatchError(e.Pos, call.Name, ms, args))
		}
		return nil
	}

	// call father`s construction method
	if call.Name == SUPER {
		if block.InheritedAttribute.IsConstructionMethod == false ||
			block.IsFunctionBlock == false ||
			block.InheritedAttribute.StatementOffset != 0 {
			*errs = append(*errs, fmt.Errorf("%s call father`s constuction on must first statement of a constructon method",
				errMsgPrefix(e.Pos)))
			return nil
		}
		if object.Type != VariableTypeObject {
			*errs = append(*errs, fmt.Errorf("%s cannot call father`s constuction on '%s'",
				errMsgPrefix(e.Pos), object.TypeString()))
			return nil
		}
		if call.Expression.isThis() == false {
			*errs = append(*errs, fmt.Errorf("%s call father`s constuction must use 'this'",
				errMsgPrefix(e.Pos)))
			return nil
		}
		err := object.Class.loadSuperClass()
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			return nil
		}
		callArgsTypes := checkExpressions(block, call.Args, errs)
		callArgsTypes = checkRightValuesValid(callArgsTypes, errs)
		ms, matched, err := object.Class.SuperClass.matchConstructionFunction(e.Pos, errs, nil, call, callArgsTypes)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			return nil
		}
		if block.InheritedAttribute.ClassMethod.isCompilerAuto && matched == false {
			//
			*errs = append(*errs, fmt.Errorf("%s compile auto constuction method cannnot match appropriate father`s constuction",
				errMsgPrefix(e.Pos)))
			return nil
		}
		if matched {
			call.Name = "<init>"
			call.Method = ms[0]
			call.Class = object.Class.SuperClass
			ret := []*Type{&Type{}}
			ret[0].Type = VariableTypeVoid
			ret[0].Pos = e.Pos
			block.Statements[0].IsCallFatherConstructionStatement = true
			block.InheritedAttribute.Function.CallFatherConstructionExpression = e
			return ret
		}
		if len(ms) == 0 {
			*errs = append(*errs, fmt.Errorf("%s 'construction' not found",
				errMsgPrefix(e.Pos)))
		} else {
			*errs = append(*errs, msNotMatchError(e.Pos, "constructor", ms, callArgsTypes))
		}
		return nil
	}
	if object.Type != VariableTypeObject && object.Type != VariableTypeClass {
		*errs = append(*errs, fmt.Errorf("%s cannot make method call named '%s' on '%s'",
			errMsgPrefix(e.Pos), call.Name, object.TypeString()))
		return nil
	}
	call.Class = object.Class
	callArgTypes := checkExpressions(block, call.Args, errs)
	callArgTypes = checkRightValuesValid(callArgTypes, errs)
	if object.Class.IsInterface() {
		if object.Type == VariableTypeClass {
			*errs = append(*errs, fmt.Errorf("%s cannot make call on interface '%s'",
				errMsgPrefix(e.Pos), object.Class.Name))
			return nil
		}
		ms, matched, err := object.Class.accessInterfaceMethod(e.Pos, errs, call.Name, call, callArgTypes, false)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			return nil
		}
		if matched {
			if ms[0].IsStatic() {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is static",
					errMsgPrefix(e.Pos), call.Name))
			}
			call.Method = ms[0]
			return ms[0].Function.Type.getReturnTypes(e.Pos)
		}
		if len(ms) == 0 {
			*errs = append(*errs, fmt.Errorf("%s method '%s' not found", errMsgPrefix(e.Pos), call.Name))
		} else {
			*errs = append(*errs, msNotMatchError(e.Pos, call.Name, ms, callArgTypes))
		}
		return nil
	}
	if len(call.ParameterTypes) > 0 {
		*errs = append(*errs, fmt.Errorf("%s method call expect no parameter types",
			errMsgPrefix(e.Pos)))
	}
	var fieldMethodHandler *ClassField
	ms, matched, err := object.Class.accessMethod(e.Pos, errs, call.Name, call, callArgTypes, false, &fieldMethodHandler)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
		return nil
	}
	if fieldMethodHandler != nil {
		if fieldMethodHandler.IsStatic() {
			if object.Type != VariableTypeClass {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is static,shoule make call from class",
					errMsgPrefix(e.Pos), call.Name))
			}
			if fieldMethodHandler.IsPublic() == false && object.Class != block.InheritedAttribute.Class {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not public", errMsgPrefix(e.Pos), call.Name))
			}
		} else {
			if false == call.Expression.isThis() &&
				fieldMethodHandler.IsPublic() == false {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not public", errMsgPrefix(e.Pos), call.Name))
			}
			if object.Type != VariableTypeObject {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not static,shoule make call from object",
					errMsgPrefix(e.Pos), call.Name))
			}
		}
		call.FieldMethodHandler = fieldMethodHandler
		return fieldMethodHandler.Type.FunctionType.getReturnTypes(e.Pos)
	}
	if matched {
		if ms[0].IsStatic() {
			if object.Type != VariableTypeClass {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is static,shoule make call from class",
					errMsgPrefix(e.Pos), call.Name))
			}
			if ms[0].IsPublic() == false && object.Class != block.InheritedAttribute.Class {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not public", errMsgPrefix(e.Pos), call.Name))
			}
		} else {
			if false == call.Expression.isThis() &&
				ms[0].IsPublic() == false {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not public", errMsgPrefix(e.Pos), call.Name))
			}
			if object.Type != VariableTypeObject {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not static,shoule make call from object",
					errMsgPrefix(e.Pos), call.Name))
			}
		}
		call.Method = ms[0]
		return ms[0].Function.Type.getReturnTypes(e.Pos)
	}
	if len(ms) == 0 {
		*errs = append(*errs, fmt.Errorf("%s method '%s' not found", errMsgPrefix(e.Pos), call.Name))
	} else {
		*errs = append(*errs, msNotMatchError(e.Pos, call.Name, ms, callArgTypes))
	}
	return nil
}
