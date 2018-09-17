package ast

import (
	"fmt"
)

func (e *Expression) check(block *Block) (returnValueTypes []*Type, errs []error) {
	if e == nil {
		return nil, []error{}
	}
	block.InheritedAttribute.Function.ExpressionCount++
	_, err := e.constantFold()
	if err != nil {
		return nil, []error{err}
	}
	errs = []error{}
	switch e.Type {
	case ExpressionTypeNull:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeNull,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeDot:
		if block.InheritedAttribute.Class == nil {
			errs = []error{fmt.Errorf("%s '%s' must in class scope",
				errMsgPrefix(e.Pos), e.Description)}
		} else {
			returnValueTypes = []*Type{
				{
					Type:  VariableTypeDynamicSelector,
					Pos:   e.Pos,
					Class: block.InheritedAttribute.Class,
				},
			}
			e.Value = returnValueTypes[0]
		}
	case ExpressionTypeBool:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeBool,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeByte:
		returnValueTypes = []*Type{{
			Type: VariableTypeByte,
			Pos:  e.Pos,
		},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeShort:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeShort,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeInt:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeInt,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeChar:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeChar,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeFloat:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeFloat,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeDouble:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeDouble,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeLong:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeLong,
				Pos:  e.Pos,
			},
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeString:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeString,
				Pos:  e.Pos,
			}}
		e.Value = returnValueTypes[0]
	case ExpressionTypeIdentifier:
		tt, err := e.checkIdentifierExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			e.Value = tt
			returnValueTypes = []*Type{tt}
		}
		//binaries
	case ExpressionTypeLogicalOr:
		fallthrough
	case ExpressionTypeLogicalAnd:
		fallthrough
	case ExpressionTypeOr:
		fallthrough
	case ExpressionTypeAnd:
		fallthrough
	case ExpressionTypeXor:
		fallthrough
	case ExpressionTypeLsh:
		fallthrough
	case ExpressionTypeRsh:
		fallthrough
	case ExpressionTypeEq:
		fallthrough
	case ExpressionTypeNe:
		fallthrough
	case ExpressionTypeGe:
		fallthrough
	case ExpressionTypeGt:
		fallthrough
	case ExpressionTypeLe:
		fallthrough
	case ExpressionTypeLt:
		fallthrough
	case ExpressionTypeAdd:
		fallthrough
	case ExpressionTypeSub:
		fallthrough
	case ExpressionTypeMul:
		fallthrough
	case ExpressionTypeDiv:
		fallthrough
	case ExpressionTypeMod:
		tt := e.checkBinaryExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.Value = tt
	case ExpressionTypeMap:
		tt := e.checkMapExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.Value = tt
	case ExpressionTypeVarAssign:
		e.checkVarAssignExpression(block, &errs)
		e.Value = mkVoidType(e.Pos)
		returnValueTypes = []*Type{e.Value}
	case ExpressionTypeAssign:
		tt := e.checkAssignExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.Value = tt
	case ExpressionTypeIncrement:
		fallthrough
	case ExpressionTypeDecrement:
		fallthrough
	case ExpressionTypePrefixIncrement:
		fallthrough
	case ExpressionTypePrefixDecrement:
		tt := e.checkIncrementExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.Value = tt
	case ExpressionTypeConst: // no return value
		errs = e.checkConstant(block)
		returnValueTypes = []*Type{mkVoidType(e.Pos)}
		e.Value = returnValueTypes[0]
	case ExpressionTypeVar:
		e.checkVarExpression(block, &errs)
		returnValueTypes = []*Type{mkVoidType(e.Pos)}
		e.Value = returnValueTypes[0]
	case ExpressionTypeFunctionCall:
		returnValueTypes = e.checkFunctionCallExpression(block, &errs)
		e.MultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			e.Value = returnValueTypes[0]
		}
	case ExpressionTypeMethodCall:
		returnValueTypes = e.checkMethodCallExpression(block, &errs)
		e.MultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			e.Value = returnValueTypes[0]
		}
	case ExpressionTypeTypeAssert:
		returnValueTypes = e.checkTypeAssert(block, &errs)
		e.MultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			e.Value = returnValueTypes[0]
		}
	case ExpressionTypeNot:
		fallthrough
	case ExpressionTypeNegative:
		fallthrough
	case ExpressionTypeBitwiseNot:
		tt := e.checkUnaryExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.Value = tt
	case ExpressionTypeQuestion:
		tt := e.checkQuestionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.Value = tt
	case ExpressionTypeIndex:
		tt := e.checkIndexExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.Value = tt
		}
	case ExpressionTypeSelection:
		tt := e.checkSelectionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.Value = tt
		}
	case ExpressionTypeCheckCast:
		tt := e.checkTypeConversionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.Value = tt
		}
	case ExpressionTypeNew:
		tt := e.checkNewExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.Value = tt
		}
	case ExpressionTypePlusAssign:
		fallthrough
	case ExpressionTypeMinusAssign:
		fallthrough
	case ExpressionTypeMulAssign:
		fallthrough
	case ExpressionTypeDivAssign:
		fallthrough
	case ExpressionTypeModAssign:
		fallthrough
	case ExpressionTypeAndAssign:
		fallthrough
	case ExpressionTypeOrAssign:
		fallthrough
	case ExpressionTypeLshAssign:
		fallthrough
	case ExpressionTypeRshAssign:
		fallthrough
	case ExpressionTypeXorAssign:
		tt := e.checkOpAssignExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.Value = tt
	case ExpressionTypeRange:
		errs = append(errs, fmt.Errorf("%s range is only work with 'for' statement",
			errMsgPrefix(e.Pos)))
	case ExpressionTypeSlice:
		tt := e.checkSlice(block, &errs)
		e.Value = tt
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
	case ExpressionTypeArray:
		tt := e.checkArray(block, &errs)
		e.Value = tt
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
	case ExpressionTypeFunctionLiteral:
		f := e.Data.(*Function)
		if e.IsStatementExpression == false && f.Name != "" {
			errs = append(errs,
				fmt.Errorf("%s function literal named '%s' expect no name", errMsgPrefix(e.Pos), f.Name))
		}
		es := f.check(block)
		errs = append(errs, es...)
		f.IsClosureFunction = f.Closure.NotEmpty(f)
		if e.IsStatementExpression {
			err := block.Insert(f.Name, f.Pos, f)
			if err != nil {
				errs = append(errs, err)
			}
		}
		returnValueTypes = make([]*Type, 1)
		returnValueTypes[0] = &Type{
			Type:         VariableTypeFunction,
			Pos:          e.Pos,
			FunctionType: &f.Type,
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeList:
		errs = append(errs, fmt.Errorf("%s cannot have expression '%s' at this scope,"+
			"this may be cause be compiler error,please contact the author",
			errMsgPrefix(e.Pos), e.Description))
	case ExpressionTypeGlobal:
		returnValueTypes = make([]*Type, 1)
		returnValueTypes[0] = &Type{
			Type:    VariableTypePackage,
			Pos:     e.Pos,
			Package: &PackageBeenCompile,
		}
		e.Value = returnValueTypes[0]
	case ExpressionTypeParenthesis:
		*e = *e.Data.(*Expression) // override
		return e.check(block)
	case ExpressionTypeVArgs:
		var t *Type
		t, errs = e.Data.(*Expression).checkSingleValueContextExpression(block)
		if len(errs) > 0 {
			return returnValueTypes, errs
		}
		e.Value = t
		returnValueTypes = []*Type{t}
		if t == nil {
			return
		}
		if t.Type != VariableTypeJavaArray {
			errs = append(errs, fmt.Errorf("%s cannot pack non java array to variable-length arguments",
				errMsgPrefix(e.Pos)))
			return
		}
		t.IsVArgs = true
	default:
		panic(fmt.Sprintf("unhandled type:%v", e.Description))
	}
	return returnValueTypes, errs
}

func (e *Expression) mustBeOneValueContext(ts []*Type) (*Type, error) {
	if len(ts) == 0 {
		return nil, nil // no-type,no error
	}
	var err error
	if len(ts) > 1 {
		err = fmt.Errorf("%s multi value in single value context", errMsgPrefix(e.Pos))
	}
	return ts[0], err
}

func (e *Expression) checkBuildInFunctionCall(block *Block, errs *[]error, f *Function, call *ExpressionFunctionCall) []*Type {
	callArgsTypes := checkExpressions(block, call.Args, errs, true)
	if f.LoadedFromCorePackage {
		if f.TemplateFunction != nil {
			tf := e.checkTemplateFunctionCall(block, errs, callArgsTypes, f)
			if tf == nil {
				return nil
			}
			var err error
			call.VArgs, err = tf.Type.fitArgs(e.Pos, &call.Args, callArgsTypes, tf)
			if err != nil {
				*errs = append(*errs, err)
			}
			return tf.Type.mkReturnTypes(e.Pos)
		} else {
			var err error
			call.VArgs, err = f.Type.fitArgs(e.Pos, &call.Args, callArgsTypes, f)
			if err != nil {
				*errs = append(*errs, err)
			}
			return f.Type.mkReturnTypes(e.Pos)
		}
	}

	length := len(*errs)
	f.buildInFunctionChecker(f, e.Data.(*ExpressionFunctionCall), block, errs, callArgsTypes, e.Pos)
	if len(*errs) == length {
		//special case ,avoid null pointer
		return f.Type.mkReturnTypes(e.Pos)
	}
	return nil //
}
func (e *Expression) checkSingleValueContextExpression(block *Block) (*Type, []error) {
	ts, es := e.check(block)
	ret, err := e.mustBeOneValueContext(ts)
	if err != nil {
		if es == nil {
			es = []error{err}
		} else {
			es = append(es, err)
		}
	}
	return ret, es
}

func (e *Expression) methodAccessAble(block *Block, m *ClassMethod, errs *[]error) {
	if e.Value.Type == VariableTypeObject {
		if m.IsStatic() {
			*errs = append(*errs, fmt.Errorf("%s method '%s' is static,shoule make call from class",
				errMsgPrefix(e.Pos), m.Function.Name))
		}
		if false == e.IsIdentifier(THIS) {
			if (e.Value.Class.LoadFromOutSide && m.IsPublic() == false) ||
				(e.Value.Class.LoadFromOutSide == false && m.IsPrivate() == true) {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not accessable",
					errMsgPrefix(e.Pos), m.Function.Name))
			}
		}
	} else {
		if m.IsStatic() == false {
			*errs = append(*errs, fmt.Errorf("%s method '%s' is not static,shoule make call from objectref",
				errMsgPrefix(e.Pos), m.Function.Name))
		}
		if e.Value.Class != block.InheritedAttribute.Class {
			if (e.Value.Class.LoadFromOutSide && m.IsPublic() == false) ||
				(e.Value.Class.LoadFromOutSide == false && m.IsPrivate() == true) {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not accessable",
					errMsgPrefix(e.Pos), m.Function.Name))
			}
		}
	}
}

func (e *Expression) fieldAccessAble(block *Block, fieldMethodHandler *ClassField, errs *[]error) {
	if e.Value.Type == VariableTypeObject {
		if fieldMethodHandler.IsStatic() {
			*errs = append(*errs, fmt.Errorf("%s method '%s' is static,shoule make call from class",
				errMsgPrefix(e.Pos), fieldMethodHandler.Name))
		}
		if false == e.IsIdentifier(THIS) {
			if (e.Value.Class.LoadFromOutSide && fieldMethodHandler.IsPublic() == false) ||
				(e.Value.Class.LoadFromOutSide == false && fieldMethodHandler.IsPrivate() == true) {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not accessable",
					errMsgPrefix(e.Pos), fieldMethodHandler.Name))
			}
		}
	} else { // class
		if fieldMethodHandler.IsStatic() == false {
			*errs = append(*errs, fmt.Errorf("%s method '%s' is not static,shoule make call from objectref",
				errMsgPrefix(e.Pos), fieldMethodHandler.Name))
		}
		if e.Value.Class != block.InheritedAttribute.Class {
			if (e.Value.Class.LoadFromOutSide && fieldMethodHandler.IsPublic() == false) ||
				(e.Value.Class.LoadFromOutSide == false && fieldMethodHandler.IsPrivate() == true) {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not accessable",
					errMsgPrefix(e.Pos), fieldMethodHandler.Name))
			}
		}
	}
}
