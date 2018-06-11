package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkSelectionExpression(block *Block, errs *[]error) (t *VariableType) {
	dot := e.Data.(*ExpressionSelection)
	t, es := dot.Expression.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	// dot
	if t.Typ != VARIABLE_TYPE_OBJECT &&
		t.Typ != VARIABLE_TYPE_CLASS &&
		t.Typ != VARIABLE_TYPE_PACKAGE {
		*errs = append(*errs, fmt.Errorf("%s cannot access field '%s' on '%s'",
			errMsgPrefix(e.Pos), dot.Name, t.TypeString()))
		return nil
	}
	var err error
	if t.Typ == VARIABLE_TYPE_PACKAGE {
		d, ok := t.Package.Block.NameExists(dot.Name)
		if ok == false {
			err = fmt.Errorf("%s '%s' not found", errMsgPrefix(e.Pos), dot.Name)
			*errs = append(*errs, err)
			return nil
		}
		switch d.(type) {
		case *VariableDefinition:
			v := d.(*VariableDefinition)
			tt := v.Typ.Clone()
			tt.Pos = e.Pos
			if (v.AccessFlags & cg.ACC_FIELD_PUBLIC) == 0 {
				err = fmt.Errorf("%s variable '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			dot.PackageVariable = v
			return tt
		case *Const:
			c := d.(*Const)
			e.fromConst(c) //
			tt := c.Typ.Clone()
			tt.Pos = e.Pos
			if c.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 {
				err = fmt.Errorf("%s const '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			return tt
		case *Class:
			c := d.(*Class)
			tt := &VariableType{}
			tt.Pos = e.Pos
			tt.Typ = VARIABLE_TYPE_CLASS
			tt.Class = c
			if (c.AccessFlags & cg.ACC_CLASS_PUBLIC) == 0 {
				err = fmt.Errorf("%s class '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			return tt
		case *EnumName:
			n := d.(*EnumName)
			if (n.Enum.AccessFlags & cg.ACC_CLASS_PUBLIC) == 0 {
				err = fmt.Errorf("%s enum '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			tt := &VariableType{}
			tt.Pos = e.Pos
			tt.Enum = n.Enum
			tt.EnumName = n
			tt.Typ = VARIABLE_TYPE_ENUM
			dot.EnumName = n
			return tt
		}
		err = fmt.Errorf("%s name '%s' cannot be used as right value", errMsgPrefix(e.Pos), dot.Name)
		*errs = append(*errs, err)
		return nil
	} else if t.Typ == VARIABLE_TYPE_OBJECT { // object
		if dot.Name == SUPER_FIELD_NAME {
			if t.Class.Name == JAVA_ROOT_CLASS {
				*errs = append(*errs, fmt.Errorf("%s '%s' is root class",
					errMsgPrefix(e.Pos), JAVA_ROOT_CLASS))
				return t
			}
			err = t.Class.loadSuperClass()
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
				return t
			}
			t := t.Clone()
			t.Pos = e.Pos
			t.Class = t.Class.SuperClass
			return t
		}
		field, err := t.Class.accessField(dot.Name, false)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		}
		if field != nil {
			if false == dot.Expression.isThis() && false == field.IsPublic() {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is private", errMsgPrefix(e.Pos),
					dot.Name))
			}
			if field.IsStatic() {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is static,cannot access by objectref",
					errMsgPrefix(e.Pos)))
			}
			t := field.Typ.Clone()
			t.Pos = e.Pos
			dot.Field = field
			return t
		}
	} else { // class
		if dot.Name == SUPER_FIELD_NAME {
			if t.Class.Name == JAVA_ROOT_CLASS {
				*errs = append(*errs, fmt.Errorf("%s '%s' is root class",
					errMsgPrefix(e.Pos), JAVA_ROOT_CLASS))
				return t
			}
			err = t.Class.loadSuperClass()
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
				return t
			}
			tt := t.Clone()
			tt.Pos = e.Pos
			tt.Class = tt.Class.SuperClass
			return tt
		}
		field, err := t.Class.accessField(dot.Name, false)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		}
		if field != nil {
			if field.IsPublic() == false && t.Class != block.InheritedAttribute.Class {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not public",
					errMsgPrefix(e.Pos),
					dot.Name))
			}
			if field.IsStatic() == false {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not static,should access by className",
					errMsgPrefix(e.Pos),
					dot.Name))
			}
			t := field.Typ.Clone()
			t.Pos = e.Pos
			dot.Field = field
			return t
		}
	}
	return nil
}
