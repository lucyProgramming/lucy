package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkSelectionExpression(block *Block, errs *[]error) *Type {
	selection := e.Data.(*ExpressionSelection)
	on, es := selection.Expression.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if on == nil {
		return nil
	}
	// dot
	if on.Type != VariableTypeObject &&
		on.Type != VariableTypeClass &&
		on.Type != VariableTypePackage {
		*errs = append(*errs, fmt.Errorf("%s cannot access field '%s' on '%s'",
			errMsgPrefix(e.Pos), selection.Name, on.TypeString()))
		return nil
	}
	var err error
	switch on.Type {
	case VariableTypePackage:
		d, ok := on.Package.Block.NameExists(selection.Name)
		if ok == false {
			err = fmt.Errorf("%s '%s' not found", errMsgPrefix(e.Pos), selection.Name)
			*errs = append(*errs, err)
			return nil
		}
		switch d.(type) {
		case *Variable:
			v := d.(*Variable)
			tt := v.Type.Clone()
			tt.Pos = e.Pos
			if (v.AccessFlags&cg.ACC_FIELD_PUBLIC) == 0 && on.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s variable '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			selection.PackageVariable = v
			return tt
		case *Constant:
			c := d.(*Constant)
			e.fromConst(c) //
			tt := c.Type.Clone()
			tt.Pos = e.Pos
			if c.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 && on.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s const '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			return tt
		case *Class:
			c := d.(*Class)
			tt := &Type{}
			tt.Pos = e.Pos
			tt.Type = VariableTypeClass
			tt.Class = c
			if (c.AccessFlags&cg.ACC_CLASS_PUBLIC) == 0 && on.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s class '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			return tt
		case *EnumName:
			n := d.(*EnumName)
			if (n.Enum.AccessFlags&cg.ACC_CLASS_PUBLIC) == 0 && on.Package.Name != PackageBeenCompile.Name {
				err = fmt.Errorf("%s enum '%s' is not public", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
			}
			tt := &Type{}
			tt.Pos = e.Pos
			tt.Enum = n.Enum
			tt.EnumName = n
			tt.Type = VariableTypeEnum
			selection.PackageEnumName = n
			return tt
		}
		err = fmt.Errorf("%s name '%s' cannot be used as right value", errMsgPrefix(e.Pos), selection.Name)
		*errs = append(*errs, err)
		return nil
	case VariableTypeObject:
		if selection.Name == SUPER {
			if on.Class.Name == JavaRootClass {
				*errs = append(*errs, fmt.Errorf("%s '%s' is root class",
					errMsgPrefix(e.Pos), JavaRootClass))
				return on
			}
			err = on.Class.loadSuperClass()
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
				return on
			}
			t := on.Clone()
			t.Pos = e.Pos
			t.Class = t.Class.SuperClass
			return t
		}
		if len(on.Class.Methods[selection.Name]) >= 1 {
			method := on.Class.Methods[selection.Name][0]
			selection.Method = method
			t := &Type{}
			t.Type = VariableTypeFunction
			t.Pos = e.Pos
			if method.IsStatic() {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is static,should access by className",
					errMsgPrefix(e.Pos),
					selection.Name))
			}
			t.FunctionType = &method.Function.Type
			return t
		}
		field, err := on.Class.accessField(selection.Name, false)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		}
		if field != nil {
			if false == selection.Expression.isThis() && false == field.IsPublic() {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is private", errMsgPrefix(e.Pos),
					selection.Name))
			}
			if field.IsStatic() {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is static,cannot access by className",
					errMsgPrefix(e.Pos), selection.Name))
			}
			t := field.Type.Clone()
			t.Pos = e.Pos
			selection.Field = field
			return t
		}
	case VariableTypeClass:
		if selection.Name == SUPER {
			if on.Class.Name == JavaRootClass {
				*errs = append(*errs, fmt.Errorf("%s '%s' is root class",
					errMsgPrefix(e.Pos), JavaRootClass))
				return on
			}
			err = on.Class.loadSuperClass()
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
				return on
			}
			tt := on.Clone()
			tt.Pos = e.Pos
			tt.Class = tt.Class.SuperClass
			return tt
		}
		if len(on.Class.Methods[selection.Name]) >= 1 {
			method := on.Class.Methods[selection.Name][0]
			selection.Method = method
			t := &Type{}
			t.Type = VariableTypeFunction
			t.Pos = e.Pos
			if method.IsStatic() == false {
				*errs = append(*errs, fmt.Errorf("%s method '%s' is not static,should access by object ref",
					errMsgPrefix(e.Pos),
					selection.Name))
			}
			t.FunctionType = &method.Function.Type
			return t
		}
		field, err := on.Class.accessField(selection.Name, false)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		}
		if field != nil {
			if field.IsPublic() == false && on.Class != block.InheritedAttribute.Class {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not public",
					errMsgPrefix(e.Pos),
					selection.Name))
			}
			if field.IsStatic() == false {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not static,should access by object ref",
					errMsgPrefix(e.Pos),
					selection.Name))
			}
			t := field.Type.Clone()
			t.Pos = e.Pos
			selection.Field = field
			return t
		}
	}
	return nil
}
