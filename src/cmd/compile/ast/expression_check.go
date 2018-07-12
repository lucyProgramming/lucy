package ast

import (
	"fmt"
)

func (e *Expression) check(block *Block) (returnValueTypes []*Type, errs []error) {
	if e == nil {
		return nil, []error{}
	}
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
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeBool:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeBool,
				Pos:  e.Pos,
			},
		}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeByte:
		returnValueTypes = []*Type{{
			Type: VariableTypeByte,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeShort:
		returnValueTypes = []*Type{{
			Type: VariableTypeShort,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeInt:
		returnValueTypes = []*Type{{
			Type: VariableTypeInt,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeFloat:
		returnValueTypes = []*Type{{
			Type: VariableTypeFloat,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeDouble:
		returnValueTypes = []*Type{{
			Type: VariableTypeDouble,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeLong:
		returnValueTypes = []*Type{{
			Type: VariableTypeLong,
			Pos:  e.Pos,
		},
		}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeString:
		returnValueTypes = []*Type{{
			Type: VariableTypeString,
			Pos:  e.Pos,
		}}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeIdentifier:
		tt, err := e.checkIdentifierExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			e.ExpressionValue = tt
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
		e.ExpressionValue = tt
	case ExpressionTypeMap:
		tt := e.checkMapExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.ExpressionValue = tt
	case ExpressionTypeColonAssign:
		e.checkColonAssignExpression(block, &errs)
		e.ExpressionValue = mkVoidType(e.Pos)
		returnValueTypes = []*Type{e.ExpressionValue}
		if t, ok := e.Data.(*ExpressionDeclareVariable); ok && t != nil && len(t.Variables) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case ExpressionTypeAssign:
		tt := e.checkAssignExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.ExpressionValue = tt
		if e.Data.(*ExpressionBinary).Left.isListAndMoreThanNElements(1) {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
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
		e.ExpressionValue = tt
	case ExpressionTypeConst: // no return value
		errs = e.checkConstant(block)
		returnValueTypes = []*Type{mkVoidType(e.Pos)}
		e.ExpressionValue = returnValueTypes[0]
	case ExpressionTypeVar:
		e.checkVarExpression(block, &errs)
		returnValueTypes = []*Type{mkVoidType(e.Pos)}
		e.ExpressionValue = returnValueTypes[0]
		if t, ok := e.Data.(*ExpressionDeclareVariable); ok && t != nil && len(t.Variables) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case ExpressionTypeFunctionCall:
		returnValueTypes = e.checkFunctionCallExpression(block, &errs)
		e.ExpressionMultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			e.ExpressionValue = returnValueTypes[0]
		}
		if len(returnValueTypes) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case ExpressionTypeMethodCall:
		returnValueTypes = e.checkMethodCallExpression(block, &errs)
		e.ExpressionMultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			e.ExpressionValue = returnValueTypes[0]
		}
		if len(returnValueTypes) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		}
	case ExpressionTypeTypeAssert:
		returnValueTypes = e.checkTypeAssert(block, &errs)
		e.ExpressionMultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			e.ExpressionValue = returnValueTypes[0]
		}
		if len(returnValueTypes) > 1 {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
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
		e.ExpressionValue = tt
	case ExpressionTypeTernary:
		tt := e.checkTernaryExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		e.ExpressionValue = tt
	case ExpressionTypeIndex:
		tt := e.checkIndexExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.ExpressionValue = tt
		}
	case ExpressionTypeSelection:
		tt := e.checkSelectionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.ExpressionValue = tt
		}
	case ExpressionTypeCheckCast:
		tt := e.checkTypeConversionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.ExpressionValue = tt
		}
	case ExpressionTypeNew:
		tt := e.checkNewExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			e.ExpressionValue = tt
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
		e.ExpressionValue = tt
	case ExpressionTypeRange:
		errs = append(errs, fmt.Errorf("%s range is only work with 'for' statement",
			errMsgPrefix(e.Pos)))
	case ExpressionTypeSlice:
		tt := e.checkSlice(block, &errs)
		e.ExpressionValue = tt
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
	case ExpressionTypeArray:
		tt := e.checkArray(block, &errs)
		e.ExpressionValue = tt
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
	case ExpressionTypeFunctionLiteral:
		f := e.Data.(*Function)
		errs = f.check(block)
		f.IsClosureFunction = f.Closure.NotEmpty(f)
		if f.Name != "" {
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
	case ExpressionTypeList:
		errs = append(errs, fmt.Errorf("%s cannot have expression '%s' at this scope,"+
			"this may be cause be compiler error,please contact the author",
			errMsgPrefix(e.Pos), e.OpName()))
	case ExpressionTypeGlobal:
		returnValueTypes = make([]*Type, 1)
		returnValueTypes[0] = &Type{
			Type:    VariableTypePackage,
			Pos:     e.Pos,
			Package: &PackageBeenCompile,
		}
		e.ExpressionValue = returnValueTypes[0]
	default:
		panic(fmt.Sprintf("unhandled type:%v", e.OpName()))
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
	callArgsTypes := checkRightValuesValid(checkExpressions(block, call.Args, errs), errs)
	length := len(*errs)
	f.buildInFunctionChecker(f, e.Data.(*ExpressionFunctionCall), block, errs, callArgsTypes, e.Pos)
	if len(*errs) == length {
		//special case ,avoid null pointer
		return f.Type.getReturnTypes(e.Pos)
	}
	return nil //
}
func (e *Expression) checkSingleValueContextExpression(block *Block) (*Type, []error) {
	ts, es := e.check(block)
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		if es == nil {
			es = []error{err}
		} else {
			es = append(es, err)
		}
	}
	return t, es
}
