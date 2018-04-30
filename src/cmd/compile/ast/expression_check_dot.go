package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"path/filepath"
)

func (e *Expression) checkDotExpression(block *Block, errs *[]error) (t *VariableType) {
	dot := e.Data.(*ExpressionDot)
	ts, es := dot.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return nil
	}
	// dot
	if t.Typ != VARIABLE_TYPE_OBJECT &&
		t.Typ != VARIABLE_TYPE_CLASS &&
		t.Typ != VARIABLE_TYPE_PACKAGE {
		*errs = append(*errs, fmt.Errorf("%s cannot access '%s' on '%s'",
			errMsgPrefix(e.Pos), dot.Name, t.TypeString()))
		return nil
	}
	if t.Typ == VARIABLE_TYPE_PACKAGE {
		find := t.Package.Block.SearchByName(dot.Name)
		if find == nil {
			fmt.Println(t.Package.Block.Vars)
			err = fmt.Errorf("%s '%s' not found", errMsgPrefix(e.Pos), dot.Name)
			*errs = append(*errs, err)
			return nil
		}
		switch find.(type) {
		case *Function: // return function
			f := find.(*Function)
			tt := &VariableType{}
			tt.Typ = VARIABLE_TYPE_FUNCTION
			tt.Function = f
			tt.Pos = e.Pos
			if (f.AccessFlags & cg.ACC_METHOD_PUBLIC) == 0 {
				err = fmt.Errorf("%s function '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			return tt
		case *Const:
			t := find.(*Const)
			e.fromConst(t) //
			tt := t.Typ.Clone()
			tt.Pos = e.Pos
			if t.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 {
				err = fmt.Errorf("%s const '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			return tt
		case *Class:
			t := find.(*Class)
			tt := &VariableType{}
			tt.Pos = e.Pos
			tt.Typ = VARIABLE_TYPE_CLASS
			tt.Class = t
			if (t.AccessFlags & cg.ACC_CLASS_PUBLIC) == 0 {
				err = fmt.Errorf("%s class '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			return tt
		case *VariableDefinition:
			t := find.(*VariableDefinition)
			tt := t.Typ.Clone()
			tt.Pos = e.Pos
			if (t.AccessFlags & cg.ACC_FIELD_PUBLIC) == 0 {
				err = fmt.Errorf("%s variable '%s' is not public", errMsgPrefix(e.Pos), dot.Name)
				*errs = append(*errs, err)
			}
			dot.PackageVariable = t
			return tt
		case *VariableType:
			err = fmt.Errorf("%s name '%s' is a type,not a expression",
				errMsgPrefix(e.Pos), dot.Name)
			*errs = append(*errs, err)
			return nil
		default:
			err = fmt.Errorf("%s name '%s' is not a expression", errMsgPrefix(e.Pos), dot.Name)
			*errs = append(*errs, err)
			return nil
		}
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
			if !dot.Expression.isThis() && !field.IsPublic() {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is private", errMsgPrefix(e.Pos),
					dot.Name))
			}
			if field.IsStatic() {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is static,should access by className(%s)",
					errMsgPrefix(e.Pos),
					dot.Name, filepath.Base(t.Class.Name)))
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
			if field.IsPublic() == false && t.Class != block.InheritedAttribute.class {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not public",
					errMsgPrefix(e.Pos),
					dot.Name))
			}
			if field.IsStatic() == false {
				*errs = append(*errs, fmt.Errorf("%s field '%s' is not static,should access by objectref",
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
