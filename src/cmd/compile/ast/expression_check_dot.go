package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
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
	if t.Typ != VARIABLE_TYPE_OBJECT && t.Typ != VARIABLE_TYPE_CLASS &&
		t.Typ != VARIABLE_TYPE_PACKAGE {
		*errs = append(*errs, fmt.Errorf("%s cannot access '%s' on '%s'", errMsgPrefix(e.Pos), dot.Name, t.TypeString()))
		return nil
	}
	if t.Typ == VARIABLE_TYPE_PACKAGE {
		find := t.Package.Block.SearchByName(dot.Name)
		if find == nil {
			err = fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), dot.Name)
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
				err = fmt.Errorf("%s function is not public", errMsgPrefix(e.Pos))
				*errs = append(*errs, err)
			}
			return tt
		case *Const:
			t := find.(*Const)
			e.fromConst(t) //
			tt := t.Typ.Clone()
			tt.Pos = e.Pos
			if t.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 {
				err = fmt.Errorf("%s function is not public", errMsgPrefix(e.Pos))
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
				err = fmt.Errorf("%s class is not public", errMsgPrefix(e.Pos))
				*errs = append(*errs, err)
			}
			return tt
		case *VariableDefinition:
			t := find.(*VariableDefinition)
			tt := t.Typ.Clone()
			tt.Pos = e.Pos
			if (t.AccessFlags & cg.ACC_FIELD_PUBLIC) == 0 {
				err = fmt.Errorf("%s variable is not public", errMsgPrefix(e.Pos))
				*errs = append(*errs, err)
			}
			dot.PackageVariableDefinition = t
			return tt
		case *VariableType:
			err = fmt.Errorf("%s name '%s' is a type,not a expression", errMsgPrefix(e.Pos), dot.Name)
			*errs = append(*errs, err)
			return nil
		default:
			err = fmt.Errorf("%s name is not a expression", errMsgPrefix(e.Pos), dot.Name)
			*errs = append(*errs, err)
			return nil
		}
	} else { // class or object
		field, err := t.Class.accessField(dot.Name, false)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		} else {
			if !dot.Expression.isThisIdentifierExpression() && !field.IsPublic() {
				*errs = append(*errs, fmt.Errorf("%s field %s is private", errMsgPrefix(e.Pos),
					dot.Name))
			}
		}
		if field != nil {
			t := field.Typ.Clone()
			t.Pos = e.Pos
			dot.Field = field
			return t
		} else {
			return nil
		}
	}
	return nil
}
