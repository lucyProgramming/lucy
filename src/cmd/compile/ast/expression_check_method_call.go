package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func (e *Expression) checkMethodCallExpression(block *Block, errs *[]error) []*VariableType {
	call := e.Data.(*ExpressionMethodCall)
	object, es := call.Expression.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if object == nil {
		return nil
	}
	if object.Type == VARIABLE_TYPE_PACKAGE {
		d, exists := object.Package.Block.NameExists(call.Name)
		if exists == false {
			*errs = append(*errs, fmt.Errorf("%s function '%s' not found", errMsgPrefix(e.Pos), call.Name))
			return nil
		}
		switch d.(type) {
		case *Function:
			f := d.(*Function)
			//if f.TemplateFunction == nil {
			//	call.PackageFunction = f
			//	return e.checkFunctionCall(block, errs, f, &call.Args)
			//} else {
			// convert to function call
			e.Type = EXPRESSION_TYPE_FUNCTION_CALL
			call := (&ExpressionFunctionCall{}).FromMethodCall(e.Data.(*ExpressionMethodCall))
			call.Func = f
			e.Data = call
			return e.checkFunctionCall(block, errs, f, &call.Args)
			//}
		case *Class:
			//object cast
			class := d.(*Class)
			typeConversion := &ExpressionTypeConversion{}
			typeConversion.Type = &VariableType{}
			typeConversion.Type.Type = VARIABLE_TYPE_OBJECT
			typeConversion.Type.Pos = e.Pos
			typeConversion.Type.Class = class
			e.Type = EXPRESSION_TYPE_CHECK_CAST
			if len(call.Args) >= 1 {
				typeConversion.Expression = call.Args[0]
			}
			e.Data = typeConversion
			if len(call.Args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument", errMsgPrefix(e.Pos)))
				return []*VariableType{typeConversion.Type.Clone()}
			}
			return []*VariableType{e.checkTypeConversionExpression(block, errs)}
		case *VariableType:
			typeConversion := &ExpressionTypeConversion{}
			typeConversion.Type = object.Package.Block.TypeAlias[call.Name]
			e.Type = EXPRESSION_TYPE_CHECK_CAST
			if len(call.Args) >= 1 {
				typeConversion.Expression = call.Args[0]
			}
			e.Data = typeConversion
			if len(call.Args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s cast type expect 1 argument",
					errMsgPrefix(e.Pos)))
				return []*VariableType{typeConversion.Type}
			}
			return []*VariableType{e.checkTypeConversionExpression(block, errs)}
		default:
			*errs = append(*errs, fmt.Errorf("%s '%s' is not a function",
				errMsgPrefix(e.Pos), call.Name))
			return nil
		}
	}
	if object.Type == VARIABLE_TYPE_MAP {
		switch call.Name {
		case common.MAP_METHOD_KEY_EXISTS:
			ret := &VariableType{}
			ret.Pos = e.Pos
			ret.Type = VARIABLE_TYPE_BOOL
			if len(call.Args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s call '%s' expect one argument",
					errMsgPrefix(e.Pos), call.Name))
				return []*VariableType{ret}
			}
			matchKey := call.Name == common.MAP_METHOD_KEY_EXISTS
			t, es := call.Args[0].checkSingleValueContextExpression(block)
			if errorsNotEmpty(es) {
				*errs = append(*errs, es...)
			}
			if t == nil {
				return []*VariableType{ret}
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
			return []*VariableType{ret}
		case common.MAP_METHOD_REMOVE:
			ret := &VariableType{}
			ret.Pos = e.Pos
			ret.Type = VARIABLE_TYPE_VOID
			if len(call.Args) == 0 {
				*errs = append(*errs, fmt.Errorf("%s remove expect at last 1 argement",
					errMsgPrefix(e.Pos)))
			}
			for _, v := range call.Args {
				ts, es := v.check(block)
				if errorsNotEmpty(es) {
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
			return []*VariableType{ret}
		case common.MAP_METHOD_REMOVEALL:
			ret := &VariableType{}
			ret.Pos = e.Pos
			ret.Type = VARIABLE_TYPE_VOID
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s removeAll expect no arguments",
					errMsgPrefix(e.Pos)))
			}
			return []*VariableType{ret}
		case common.MAP_METHOD_SIZE:
			ret := &VariableType{}
			ret.Pos = e.Pos
			ret.Type = VARIABLE_TYPE_INT
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s too many argument to call '%s''",
					errMsgPrefix(e.Pos), call.Name))
			}
			return []*VariableType{ret}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on map", errMsgPrefix(e.Pos), call.Name))
			return nil
		}
		return nil
	}
	if object.Type == VARIABLE_TYPE_JAVA_ARRAY {
		switch call.Name {
		case common.ARRAY_METHOD_SIZE:
			t := &VariableType{}
			t.Type = VARIABLE_TYPE_INT
			t.Pos = e.Pos
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s method '%s' expect no arguments",
					errMsgPrefix(e.Pos), call.Name))
			}
			return []*VariableType{t}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on '%s'",
				errMsgPrefix(e.Pos), call.Name, object.TypeString()))
		}
		return nil
	}

	if object.Type == VARIABLE_TYPE_ARRAY {
		switch call.Name {
		case common.ARRAY_METHOD_SIZE,
			common.ARRAY_METHOD_CAP,   //for debug,remove when time is right
			common.ARRAY_METHOD_START, //for debug,remove when time is right
			common.ARRAY_METHOD_END:   //for debug,remove when time is right
			t := &VariableType{}
			t.Type = VARIABLE_TYPE_INT
			t.Pos = e.Pos
			if len(call.Args) > 0 {
				*errs = append(*errs, fmt.Errorf("%s too mamy argument to call,method '%s' expect no arguments",
					errMsgPrefix(e.Pos), call.Name))
			}
			return []*VariableType{t}
		case common.ARRAY_METHOD_APPEND, common.ARRAY_METHOD_APPEND_ALL:
			if len(call.Args) == 0 {
				*errs = append(*errs, fmt.Errorf("%s too few arguments to call %s,expect at least one argument",
					errMsgPrefix(e.Pos), call.Name))
			}
			for _, e := range call.Args {
				ts, es := e.check(block)
				if errorsNotEmpty(es) {
					*errs = append(*errs, es...)
				}
				for _, t := range ts {
					if call.Name == common.ARRAY_METHOD_APPEND {
						if object.ArrayType.Equal(errs, t) == false {
							*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' to call method '%s'",
								errMsgPrefix(t.Pos), t.TypeString(), object.ArrayType.TypeString(), call.Name))
						}
					} else {
						if object.Equal(errs, t) == false {
							*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' to call method '%s'",
								errMsgPrefix(t.Pos), t.TypeString(), object.ArrayType.TypeString(), call.Name))
						}
					}
				}
			}
			t := object.Clone()
			t.Pos = e.Pos
			return []*VariableType{t}
		default:
			*errs = append(*errs, fmt.Errorf("%s unkown call '%s' on array", errMsgPrefix(e.Pos), call.Name))
		}
		return nil
	}

	if object.Type == VARIABLE_TYPE_STRING {
		if err := loadJavaStringClass(e.Pos); err != nil {
			*errs = append(*errs, err)
			return nil
		}
		args := checkRightValuesValid(checkExpressions(block, call.Args, errs), errs)
		ms, matched, err := javaStringClass.accessMethod(e.Pos, errs, call.Name, args, nil, false)
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
			return ms[0].Func.Type.returnTypes(e.Pos)
		}
		if len(ms) == 0 {
			*errs = append(*errs, fmt.Errorf("%s method '%s' not found", errMsgPrefix(e.Pos), call.Name))
		} else {
			*errs = append(*errs, msNotMatchError(e.Pos, call.Name, ms, args))
		}
		return nil
	}
	if object.Type != VARIABLE_TYPE_OBJECT && object.Type != VARIABLE_TYPE_CLASS {
		*errs = append(*errs, fmt.Errorf("%s cannot make method call named '%s' on '%s'",
			errMsgPrefix(e.Pos), call.Name, object.TypeString()))
		return nil
	}
	// call father`s construction method
	if call.Name == SUPER_FIELD_NAME {
		if block.InheritedAttribute.IsConstruction == false ||
			block.IsFunctionTopBlock == false ||
			block.InheritedAttribute.StatementOffset != 0 {
			*errs = append(*errs, fmt.Errorf("%s call father`s constuction on must first statement of a constructon method",
				errMsgPrefix(e.Pos)))
			return nil
		}
		if object.Type != VARIABLE_TYPE_OBJECT {
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
		args := checkExpressions(block, call.Args, errs)
		args = checkRightValuesValid(args, errs)
		ms, matched, err := object.Class.SuperClass.matchConstructionFunction(e.Pos, errs, args, &call.Args)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			return nil
		}
		if matched {
			call.Name = "<init>"
			block.InheritedAttribute.Function.ConstructionMethodCalledByUser = true
			call.Method = ms[0]
			call.Class = object.Class.SuperClass
			ret := []*VariableType{&VariableType{}}
			ret[0].Type = VARIABLE_TYPE_VOID
			ret[0].Pos = e.Pos
			block.Statements[0].IsCallFatherConstructionStatement = true
			return ret
		}
		if len(ms) == 0 {
			*errs = append(*errs, fmt.Errorf("%s 'construction' not found",
				errMsgPrefix(e.Pos)))
		} else {
			*errs = append(*errs, msNotMatchError(e.Pos, "constructor", ms, args))
		}
		return nil
	}
	call.Class = object.Class
	args := checkExpressions(block, call.Args, errs)
	args = checkRightValuesValid(args, errs)
	ms, matched, err := object.Class.accessMethod(e.Pos, errs, call.Name, args, &call.Args, false)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
		return nil
	}
	if matched {
		if ms[0].IsStatic() {
			if object.Type != VARIABLE_TYPE_CLASS {
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
			if object.Type != VARIABLE_TYPE_OBJECT {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not static,shoule make call from object",
					errMsgPrefix(e.Pos), call.Name))
			}
		}
		call.Method = ms[0]
		return ms[0].Func.Type.returnTypes(e.Pos)
	}
	if len(ms) == 0 {
		*errs = append(*errs, fmt.Errorf("%s method '%s' not found", errMsgPrefix(e.Pos), call.Name))
	} else {
		*errs = append(*errs, msNotMatchError(e.Pos, call.Name, ms, args))
	}
	return nil
}
