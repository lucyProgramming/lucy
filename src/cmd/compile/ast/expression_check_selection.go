package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkSelectionExpression(block *Block, errs *[]error) *Type {
	selection := e.Data.(*ExpressionSelection)
	object, es := selection.Expression.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if object == nil {
		return nil
	}
	// dot
	if object.Type != VariableTypeObject &&
		object.Type != VariableTypeClass &&
		object.Type != VariableTypePackage {
		*errs = append(*errs, fmt.Errorf("%s cannot access field '%s' on '%s'",
			errMsgPrefix(e.Pos), selection.Name, object.TypeString()))
		return nil
	}
	var err error
	switch object.Type {
	case VariableTypePackage:
		d, ok := object.Package.Block.NameExists(selection.Name)
		if ok == false {
			err = fmt.Errorf("%s '%s' not found", errMsgPrefix(e.Pos), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
		switch d.(type) {
		case *Variable:
			v := d.(*Variable)
			result := v.Type.Clone()
			result.Pos = e.Pos
			if (v.AccessFlags&cg.ACC_FIELD_PUBLIC) == 0 && object.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s variable '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			selection.PackageVariable = v
			return result
		case *Constant:
			c := d.(*Constant)
			e.fromConst(c) //
			result := c.Type.Clone()
			result.Pos = e.Pos
			if c.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 && object.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s const '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			return result
		case *Class:
			c := d.(*Class)
			result := &Type{}
			result.Pos = e.Pos
			result.Type = VariableTypeClass
			result.Class = c
			if (c.AccessFlags&cg.ACC_CLASS_PUBLIC) == 0 && object.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s class '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			return result
		case *EnumName:
			n := d.(*EnumName)
			if (n.Enum.AccessFlags&cg.ACC_CLASS_PUBLIC) == 0 && object.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s enum '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			result := &Type{}
			result.Pos = e.Pos
			result.Enum = n.Enum
			result.EnumName = n
			result.Type = VariableTypeEnum
			selection.PackageEnumName = n
			return result
		case *Function:
			f := d.(*Function)
			if (f.AccessFlags&cg.ACC_METHOD_PUBLIC) == 0 && object.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s enum '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			result := &Type{}
			result.Pos = e.Pos
			result.Type = VariableTypeFunction
			result.FunctionType = &f.Type
			selection.PackageFunction = f
			return result
		default:
			err = fmt.Errorf("%s name '%s' cannot be used as right value", errMsgPrefix(e.Pos), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
	case VariableTypeObject:
		if selection.Name == SUPER {
			if object.Class.Name == JavaRootClass {
				*errs = append(*errs, fmt.Errorf("%s '%s' is root class",
					errMsgPrefix(e.Pos), JavaRootClass))
				return object
			}
			err = object.Class.loadSuperClass(e.Pos)
			if err != nil {
				*errs = append(*errs, err)
				return object
			}
			result := object.Clone()
			result.Pos = e.Pos
			result.Class = result.Class.SuperClass
			return result
		}
		fieldOrMethod, err := object.Class.getFieldOrMethod(e.Pos, selection.Name, false)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if field, ok := fieldOrMethod.(*ClassField); ok {
			if selection.Expression.isThis() == false {
				if (selection.Expression.Value.Class.LoadFromOutSide && field.IsPublic() == false) ||
					(selection.Expression.Value.Class.LoadFromOutSide == false && field.IsPrivate()) {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is private",
						errMsgPrefix(e.Pos), selection.Name))
				}
			}
			if field.IsStatic() {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is static,cannot access by className",
					errMsgPrefix(e.Pos), selection.Name))
			}
			result := field.Type.Clone()
			result.Pos = e.Pos
			selection.Field = field
			return result
		} else {
			method := fieldOrMethod.(*ClassMethod)
			if method.IsStatic() {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is static,should access by className",
					errMsgPrefix(e.Pos),
					selection.Name))
			}
			if selection.Expression.isThis() == false {
				if (selection.Expression.Value.Class.LoadFromOutSide && method.IsPublic() == false) ||
					(selection.Expression.Value.Class.LoadFromOutSide == false && method.IsPrivate()) {
					*errs = append(*errs, fmt.Errorf("%s method '%s' is private",
						errMsgPrefix(e.Pos), selection.Name))
				}
			}
			selection.Method = method
			result := &Type{}
			result.Type = VariableTypeFunction
			result.FunctionType = &method.Function.Type
			result.Pos = e.Pos
			return result
		}
	case VariableTypeClass:
		if selection.Name == SUPER {
			*errs = append(*errs, fmt.Errorf("%s cannot access super on class named '%s'",
				errMsgPrefix(e.Pos), object.Class.Name))
			return nil
		}
		fieldOrMethod, err := object.Class.getFieldOrMethod(e.Pos, selection.Name, false)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if field, ok := fieldOrMethod.(*ClassField); ok {
			if block.InheritedAttribute.Class != selection.Expression.Value.Class {
				if (selection.Expression.Value.Class.LoadFromOutSide && field.IsPublic() == false) ||
					(selection.Expression.Value.Class.LoadFromOutSide == false && field.IsPrivate()) {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is private",
						errMsgPrefix(e.Pos), selection.Name))
				}
			}
			if field.IsStatic() == false {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not static,should access by object ref",
					errMsgPrefix(e.Pos),
					selection.Name))
			}
			result := field.Type.Clone()
			result.Pos = e.Pos
			selection.Field = field
			return result
		} else {
			method := fieldOrMethod.(*ClassMethod)
			if block.InheritedAttribute.Class != selection.Expression.Value.Class {
				if (selection.Expression.Value.Class.LoadFromOutSide && method.IsPublic() == false) ||
					(selection.Expression.Value.Class.LoadFromOutSide == false && method.IsPrivate()) {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is private",
						errMsgPrefix(e.Pos), selection.Name))
				}
			}
			selection.Method = method
			result := &Type{}
			result.Type = VariableTypeFunction
			result.Pos = e.Pos
			if method.IsStatic() == false {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not static,should access by object ref",
					errMsgPrefix(e.Pos),
					selection.Name))
			}
			result.FunctionType = &method.Function.Type
			return result
		}
	}
	return nil
}
