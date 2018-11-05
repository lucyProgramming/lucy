package ast

import (
	"fmt"
)

func (this *Expression) check(block *Block) (returnValueTypes []*Type, errs []error) {
	if this == nil {
		return nil, []error{}
	}
	_, err := this.constantFold()
	if err != nil {
		return nil, []error{err}
	}
	errs = []error{}
	switch this.Type {
	case ExpressionTypeNull:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeNull,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeDot:
		if block.InheritedAttribute.Class == nil {
			errs = []error{fmt.Errorf("%s '%s' must in class scope",
				this.Pos.ErrMsgPrefix(), this.Op)}
		} else {
			returnValueTypes = []*Type{
				{
					Type:  VariableTypeDynamicSelector,
					Pos:   this.Pos,
					Class: block.InheritedAttribute.Class,
				},
			}
			this.Value = returnValueTypes[0]
		}
	case ExpressionTypeBool:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeBool,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeByte:
		returnValueTypes = []*Type{{
			Type: VariableTypeByte,
			Pos:  this.Pos,
		},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeShort:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeShort,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeInt:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeInt,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeChar:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeChar,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeFloat:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeFloat,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeDouble:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeDouble,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeLong:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeLong,
				Pos:  this.Pos,
			},
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeString:
		returnValueTypes = []*Type{
			{
				Type: VariableTypeString,
				Pos:  this.Pos,
			}}
		this.Value = returnValueTypes[0]
	case ExpressionTypeIdentifier:
		tt, err := this.checkIdentifierExpression(block)
		if err != nil {
			errs = append(errs, err)
		}
		if tt != nil {
			this.Value = tt
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
		length := len(errs)
		tt := this.checkBinaryExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		if len(errs) == length { // no error
			if ee := this.binaryExpressionDependOnSub(); ee != nil {
				*this = *ee
			}
		}
		this.Value = tt
	case ExpressionTypeMap:
		tt := this.checkMapExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		this.Value = tt
	case ExpressionTypeVarAssign:
		this.checkVarAssignExpression(block, &errs)
		this.Value = mkVoidType(this.Pos)
		returnValueTypes = []*Type{this.Value}
	case ExpressionTypeAssign:
		tt := this.checkAssignExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		this.Value = tt
	case ExpressionTypeIncrement:
		fallthrough
	case ExpressionTypeDecrement:
		fallthrough
	case ExpressionTypePrefixIncrement:
		fallthrough
	case ExpressionTypePrefixDecrement:
		tt := this.checkIncrementExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		this.Value = tt
	case ExpressionTypeConst: // no return value
		errs = this.checkConstant(block)
		returnValueTypes = []*Type{mkVoidType(this.Pos)}
		this.Value = returnValueTypes[0]
	case ExpressionTypeVar:
		this.checkVarExpression(block, &errs)
		returnValueTypes = []*Type{mkVoidType(this.Pos)}
		this.Value = returnValueTypes[0]
	case ExpressionTypeFunctionCall:
		returnValueTypes = this.checkFunctionCallExpression(block, &errs)
		this.MultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			this.Value = returnValueTypes[0]
		}
	case ExpressionTypeMethodCall:
		returnValueTypes = this.checkMethodCallExpression(block, &errs)
		this.MultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			this.Value = returnValueTypes[0]
		}
	case ExpressionTypeTypeAssert:
		returnValueTypes = this.checkTypeAssert(block, &errs)
		this.MultiValues = returnValueTypes
		if len(returnValueTypes) > 0 {
			this.Value = returnValueTypes[0]
		}
	case ExpressionTypeNot:
		fallthrough
	case ExpressionTypeNegative:
		fallthrough
	case ExpressionTypeBitwiseNot:
		tt := this.checkUnaryExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		this.Value = tt
	case ExpressionTypeQuestion:
		tt := this.checkQuestionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		this.Value = tt
	case ExpressionTypeIndex:
		tt := this.checkIndexExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			this.Value = tt
		}
	case ExpressionTypeSelection:
		tt := this.checkSelectionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			this.Value = tt
		}
	case ExpressionTypeSelectionConst:
		tt := this.checkSelectConstExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			this.Value = tt
		}
	case ExpressionTypeCheckCast:
		tt := this.checkTypeConversionExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			this.Value = tt
		}
	case ExpressionTypeNew:
		tt := this.checkNewExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
			this.Value = tt
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
		tt := this.checkOpAssignExpression(block, &errs)
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
		this.Value = tt
	case ExpressionTypeRange:
		errs = append(errs, fmt.Errorf("%s range is only work with 'for' statement",
			errMsgPrefix(this.Pos)))
	case ExpressionTypeSlice:
		tt := this.checkSlice(block, &errs)
		this.Value = tt
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
	case ExpressionTypeArray:
		tt := this.checkArray(block, &errs)
		this.Value = tt
		if tt != nil {
			returnValueTypes = []*Type{tt}
		}
	case ExpressionTypeFunctionLiteral:
		f := this.Data.(*Function)
		PackageBeenCompile.statementLevelFunctions =
			append(PackageBeenCompile.statementLevelFunctions, f)
		if this.IsStatementExpression {
			err := block.Insert(f.Name, f.Pos, f)
			if err != nil {
				errs = append(errs, err)
			}
		}
		es := f.check(block)
		errs = append(errs, es...)
		returnValueTypes = make([]*Type, 1)
		returnValueTypes[0] = &Type{
			Type:         VariableTypeFunction,
			Pos:          this.Pos,
			FunctionType: &f.Type,
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeList:
		errs = append(errs,
			fmt.Errorf("%s cannot have expression '%s' at this scope,"+
				"this may be cause by the compiler error,please contact the author",
				this.Pos.ErrMsgPrefix(), this.Op))
	case ExpressionTypeGlobal:
		returnValueTypes = make([]*Type, 1)
		returnValueTypes[0] = &Type{
			Type:    VariableTypePackage,
			Pos:     this.Pos,
			Package: &PackageBeenCompile,
		}
		this.Value = returnValueTypes[0]
	case ExpressionTypeParenthesis:
		*this = *this.Data.(*Expression) // override
		return this.check(block)
	case ExpressionTypeVArgs:
		var t *Type
		t, errs = this.Data.(*Expression).checkSingleValueContextExpression(block)
		if len(errs) > 0 {
			return returnValueTypes, errs
		}
		this.Value = t
		returnValueTypes = []*Type{t}
		if t == nil {
			return
		}
		if t.Type != VariableTypeJavaArray {
			errs = append(errs, fmt.Errorf("%s cannot pack non java array to variable-length arguments",
				errMsgPrefix(this.Pos)))
			return
		}
		t.IsVariableArgs = true
	default:
		panic(fmt.Sprintf("unhandled type:%v", this.Op))
	}
	return returnValueTypes, errs
}

func (this *Expression) mustBeOneValueContext(ts []*Type) (*Type, error) {
	if len(ts) == 0 {
		return nil, nil // no-type,no error
	}
	var err error
	if len(ts) > 1 {
		err = fmt.Errorf("%s multi value in single value context", errMsgPrefix(this.Pos))
	}
	return ts[0], err
}

func (this *Expression) checkSingleValueContextExpression(block *Block) (*Type, []error) {
	ts, es := this.check(block)
	ret, err := this.mustBeOneValueContext(ts)
	if err != nil {
		if es == nil {
			es = []error{err}
		} else {
			es = append(es, err)
		}
	}
	return ret, es
}

func (this *Expression) methodAccessAble(block *Block, method *ClassMethod) error {
	if this.Value.Type == VariableTypeObject {
		if method.IsStatic() {
			return fmt.Errorf("%s method '%s' is static",
				this.Pos.ErrMsgPrefix(), method.Function.Name)
		}
		if false == this.IsIdentifier(ThisPointerName) {
			if this.Value.Class.LoadFromOutSide {
				if this.Value.Class.IsPublic() == false {
					return fmt.Errorf("%s class '%s' is not public",
						this.Pos.ErrMsgPrefix(), this.Value.Class.Name)
				}
				if method.IsPublic() == false {
					return fmt.Errorf("%s method '%s' is not public",
						this.Pos.ErrMsgPrefix(), method.Function.Name)
				}
			} else {
				if method.IsPrivate() {
					return fmt.Errorf("%s method '%s' is private",
						this.Pos.ErrMsgPrefix(), method.Function.Name)
				}
			}
		}
	} else {
		if method.IsStatic() == false {
			return fmt.Errorf("%s method '%s' is a instance method",
				this.Pos.ErrMsgPrefix(), method.Function.Name)
		}
		if this.Value.Class != block.InheritedAttribute.Class {
			if this.Value.Class.LoadFromOutSide {
				if this.Value.Class.IsPublic() == false {
					return fmt.Errorf("%s class '%s' is not public",
						this.Pos.ErrMsgPrefix(), this.Value.Class.Name)
				}
				if method.IsPublic() == false {
					return fmt.Errorf("%s method '%s' is not public",
						this.Pos.ErrMsgPrefix(), method.Function.Name)
				}
			} else {
				if method.IsPrivate() {
					return fmt.Errorf("%s method '%s' is private",
						this.Pos.ErrMsgPrefix(), method.Function.Name)
				}
			}
		}
	}
	return nil
}

func (this *Expression) fieldAccessAble(block *Block, field *ClassField) error {
	if this.Value.Type == VariableTypeObject {
		if field.IsStatic() {
			return fmt.Errorf("%s field '%s' is static",
				this.Pos.ErrMsgPrefix(), field.Name)
		}
		if false == this.IsIdentifier(ThisPointerName) {
			if this.Value.Class.LoadFromOutSide {
				if this.Value.Class.IsPublic() == false {
					return fmt.Errorf("%s class '%s' is not public",
						this.Pos.ErrMsgPrefix(), this.Value.Class.Name)
				}
				if field.IsPublic() == false {
					return fmt.Errorf("%s field '%s' is not public",
						this.Pos.ErrMsgPrefix(), field.Name)
				}
			} else {
				if field.IsPrivate() {
					return fmt.Errorf("%s field '%s' is private",
						this.Pos.ErrMsgPrefix(), field.Name)
				}
			}
		}
	} else { // class
		if field.IsStatic() == false {
			return fmt.Errorf("%s field '%s' is not static",
				this.Pos.ErrMsgPrefix(), field.Name)
		}
		if this.Value.Class != block.InheritedAttribute.Class {
			if this.Value.Class.LoadFromOutSide {
				if this.Value.Class.IsPublic() == false {
					return fmt.Errorf("%s class '%s' is not public",
						this.Pos.ErrMsgPrefix(), this.Value.Class.Name)
				}
				if field.IsPublic() == false {
					return fmt.Errorf("%s field '%s' is not public",
						this.Pos.ErrMsgPrefix(), field.Name)
				}
			} else {
				if field.IsPrivate() {
					return fmt.Errorf("%s field '%s' is private",
						this.Pos.ErrMsgPrefix(), field.Name)
				}
			}
		}
	}
	return nil
}
