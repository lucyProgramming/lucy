package ast

import (
	"fmt"
)

func (this *Expression) checkSelectionExpression(block *Block, errs *[]error) *Type {
	selection := this.Data.(*ExpressionSelection)
	object, es := selection.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	switch object.Type {
	case VariableTypeMagicFunction:
		v := object.Function.Type.searchName(selection.Name)
		if v == nil {
			err := fmt.Errorf("%s '%s' not found",
				this.Pos.ErrMsgPrefix(), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
		this.Value = v.Type.Clone()
		this.Value.Pos = this.Pos
		this.Type = ExpressionTypeIdentifier
		identifier := &ExpressionIdentifier{}
		identifier.Name = selection.Name
		identifier.Variable = v
		this.Data = identifier
		return this.Value
	case VariableTypeDynamicSelector:
		if selection.Name == SUPER {
			*errs = append(*errs, fmt.Errorf("%s access '%s' at '%s' not allow",
				this.Pos.ErrMsgPrefix(), SUPER, object.TypeString()))
			return nil
		}
		access, err := object.Class.getFieldOrMethod(this.Pos, selection.Name, false)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if field, ok := access.(*ClassField); ok {
			selection.Field = field
			result := field.Type.Clone()
			result.Pos = this.Pos
			return result
		} else {
			method := access.(*ClassMethod)
			selection.Method = method
			result := &Type{
				Type:         VariableTypeFunction,
				FunctionType: &method.Function.Type,
				Pos:          this.Pos,
			}
			return result
		}
	case VariableTypePackage:
		d, ok := object.Package.Block.NameExists(selection.Name)
		if ok == false {
			err := fmt.Errorf("%s '%s' not found",
				this.Pos.ErrMsgPrefix(), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
		switch d.(type) {
		case *Variable:
			v := d.(*Variable)
			result := v.Type.Clone()
			result.Pos = this.Pos
			if v.isPublic() == false && object.Package.isSame(&PackageBeenCompile) == false {
				err := fmt.Errorf("%s variable '%s' is not public",
					this.Pos.ErrMsgPrefix(), selection.Name)
				*errs = append(*errs, err)
			}
			selection.PackageVariable = v
			return result
		case *Constant:
			c := d.(*Constant)
			this.fromConst(c) //
			result := c.Type.Clone()
			result.Pos = this.Pos
			if c.isPublic() == false && object.Package.isSame(&PackageBeenCompile) == false {
				err := fmt.Errorf("%s const '%s' is not public",
					this.Pos.ErrMsgPrefix(), selection.Name)
				*errs = append(*errs, err)
			}
			return result
		case *Class:
			c := d.(*Class)
			result := &Type{}
			result.Pos = this.Pos
			result.Type = VariableTypeClass
			result.Class = c
			if c.IsPublic() == false && object.Package.isSame(&PackageBeenCompile) == false {
				err := fmt.Errorf("%s class '%s' is not public",
					this.Pos.ErrMsgPrefix(), selection.Name)
				*errs = append(*errs, err)
			}
			return result
		case *EnumName:
			n := d.(*EnumName)
			if n.Enum.isPublic() == false && object.Package.isSame(&PackageBeenCompile) == false {
				err := fmt.Errorf("%s enum '%s' is not public",
					this.Pos.ErrMsgPrefix(), selection.Name)
				*errs = append(*errs, err)
			}
			result := &Type{}
			result.Pos = this.Pos
			result.Enum = n.Enum
			result.EnumName = n
			result.Type = VariableTypeEnum
			selection.PackageEnumName = n
			return result
		case *Function:
			f := d.(*Function)
			if f.IsPublic() == false && object.Package.isSame(&PackageBeenCompile) == false {
				err := fmt.Errorf("%s function '%s' is not public",
					this.Pos.ErrMsgPrefix(), selection.Name)
				*errs = append(*errs, err)
			}
			if f.TemplateFunction != nil {
				err := fmt.Errorf("%s function '%s' is a template function",
					this.Pos.ErrMsgPrefix(), selection.Name)
				*errs = append(*errs, err)
				return nil
			}
			result := &Type{}
			result.Pos = this.Pos
			result.Type = VariableTypeFunction
			result.FunctionType = &f.Type
			selection.PackageFunction = f
			return result
		default:
			err := fmt.Errorf("%s name '%s' cannot be used as right value",
				this.Pos.ErrMsgPrefix(), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
	case VariableTypeObject, VariableTypeClass:
		if selection.Name == SUPER {
			if object.Type == VariableTypeClass {
				*errs = append(*errs, fmt.Errorf("%s cannot access class`s super",
					object.Pos.ErrMsgPrefix()))
				return object
			}
			if object.Class.Name == JavaRootClass {
				*errs = append(*errs, fmt.Errorf("%s '%s' is root class",
					object.Pos.ErrMsgPrefix(), JavaRootClass))
				return object
			}
			err := object.Class.loadSuperClass(this.Pos)
			if err != nil {
				*errs = append(*errs, err)
				return object
			}
			if object.Class.SuperClass == nil {
				return object
			}
			result := object.Clone()
			result.Pos = this.Pos
			result.Class = result.Class.SuperClass
			return result
		}
		fieldOrMethod, err := object.Class.getFieldOrMethod(this.Pos, selection.Name, false)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if field, ok := fieldOrMethod.(*ClassField); ok {
			err := selection.Expression.fieldAccessAble(block, field)
			if err != nil {
				*errs = append(*errs, err)
			}
			result := field.Type.Clone()
			result.Pos = this.Pos
			selection.Field = field
			return result
		} else {
			method := fieldOrMethod.(*ClassMethod)
			err := selection.Expression.methodAccessAble(block, method)
			if err != nil {
				*errs = append(*errs, err)
			}
			selection.Method = method
			result := &Type{}
			result.Type = VariableTypeFunction
			result.FunctionType = &method.Function.Type
			result.Pos = this.Pos
			return result
		}

	default:
		*errs = append(*errs, fmt.Errorf("%s cannot access '%s' on '%s'",
			this.Pos.ErrMsgPrefix(), selection.Name, object.TypeString()))
		return nil
	}
	return nil
}
