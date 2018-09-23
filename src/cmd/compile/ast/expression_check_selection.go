package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkSelectionExpression(block *Block, errs *[]error) *Type {
	selection := e.Data.(*ExpressionSelection)
	object, es := selection.Expression.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if object == nil {
		return nil
	}
	switch object.Type {
	case VariableTypeMagicFunction:
		v := object.Function.Type.searchName(selection.Name)
		if v == nil {
			err := fmt.Errorf("%s '%s' not found", errMsgPrefix(selection.Pos), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
		e.Value = v.Type.Clone()
		e.Value.Pos = selection.Pos
		e.Type = ExpressionTypeIdentifier
		identifier := &ExpressionIdentifier{}
		identifier.Name = selection.Name
		identifier.Variable = v
		e.Data = identifier
		return e.Value
	case VariableTypeDynamicSelector:
		if selection.Name == SUPER {
			*errs = append(*errs, fmt.Errorf("%s access '%s' at '%s' not allow",
				errMsgPrefix(selection.Pos), SUPER, object.TypeString()))
			return nil
		}
		access, err := object.Class.getFieldOrMethod(selection.Pos, selection.Name, false)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if field, ok := access.(*ClassField); ok {
			selection.Field = field
			result := field.Type.Clone()
			result.Pos = selection.Pos
			return result
		} else {
			method := access.(*ClassMethod)
			selection.Method = method
			result := &Type{
				Type:         VariableTypeFunction,
				FunctionType: &method.Function.Type,
				Pos:          selection.Pos,
			}
			return result
		}
	case VariableTypePackage:
		d, ok := object.Package.Block.NameExists(selection.Name)
		if ok == false {
			err := fmt.Errorf("%s '%s' not found",
				errMsgPrefix(selection.Pos), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
		switch d.(type) {
		case *Variable:
			v := d.(*Variable)
			result := v.Type.Clone()
			result.Pos = selection.Pos
			if (v.AccessFlags&cg.ACC_FIELD_PUBLIC) == 0 &&
				object.Package.Name != PackageBeenCompile.Name {
				err := fmt.Errorf("%s variable '%s' is not public",
					errMsgPrefix(selection.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			selection.PackageVariable = v
			return result
		case *Constant:
			c := d.(*Constant)
			e.fromConst(c) //
			result := c.Type.Clone()
			result.Pos = selection.Pos
			if c.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 &&
				object.Package.Name != PackageBeenCompile.Name {
				err := fmt.Errorf("%s const '%s' is not public",
					errMsgPrefix(selection.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			return result
		case *Class:
			c := d.(*Class)
			result := &Type{}
			result.Pos = selection.Pos
			result.Type = VariableTypeClass
			result.Class = c
			if (c.AccessFlags&cg.ACC_CLASS_PUBLIC) == 0 &&
				object.Package.Name != PackageBeenCompile.Name {
				err := fmt.Errorf("%s class '%s' is not public",
					errMsgPrefix(selection.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			return result
		case *EnumName:
			n := d.(*EnumName)
			if (n.Enum.AccessFlags&cg.ACC_CLASS_PUBLIC) == 0 &&
				object.Package.Name != PackageBeenCompile.Name {
				err := fmt.Errorf("%s enum '%s' is not public",
					errMsgPrefix(selection.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			result := &Type{}
			result.Pos = selection.Pos
			result.Enum = n.Enum
			result.EnumName = n
			result.Type = VariableTypeEnum
			selection.PackageEnumName = n
			return result
		case *Function:
			f := d.(*Function)
			if (f.AccessFlags&cg.ACC_METHOD_PUBLIC) == 0 &&
				object.Package.Name != PackageBeenCompile.Name {
				err := fmt.Errorf("%s enum '%s' is not public",
					errMsgPrefix(selection.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			result := &Type{}
			result.Pos = selection.Pos
			result.Type = VariableTypeFunction
			result.FunctionType = &f.Type
			selection.PackageFunction = f
			return result
		default:
			err := fmt.Errorf("%s name '%s' cannot be used as right value",
				errMsgPrefix(selection.Pos), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
	case VariableTypeObject, VariableTypeClass:
		if selection.Expression.Value.Type == VariableTypeObject {
			if selection.Name == SUPER {
				if object.Class.Name == JavaRootClass {
					*errs = append(*errs, fmt.Errorf("%s '%s' is root class",
						errMsgPrefix(selection.Pos), JavaRootClass))
					return object
				}
				err := object.Class.loadSuperClass(selection.Pos)
				if err != nil {
					*errs = append(*errs, err)
					return object
				}
				result := object.Clone()
				result.Pos = selection.Pos
				result.Class = result.Class.SuperClass
				return result
			}
		}

		fieldOrMethod, err := object.Class.getFieldOrMethod(selection.Pos, selection.Name, false)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if field, ok := fieldOrMethod.(*ClassField); ok {
			selection.Expression.fieldAccessAble(block, field, errs)
			result := field.Type.Clone()
			result.Pos = selection.Pos
			selection.Field = field
			return result
		} else {
			method := fieldOrMethod.(*ClassMethod)
			selection.Expression.methodAccessAble(block, method, errs)
			selection.Method = method
			result := &Type{}
			result.Type = VariableTypeFunction
			result.FunctionType = &method.Function.Type
			result.Pos = selection.Pos
			return result
		}

	default:
		*errs = append(*errs, fmt.Errorf("%s cannot access '%s' on '%s'",
			errMsgPrefix(selection.Pos), selection.Name, object.TypeString()))
		return nil
	}
	return nil
}
