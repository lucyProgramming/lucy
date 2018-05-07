package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) getLeftValue(block *Block, errs *[]error) (t *VariableType) {
	switch e.Typ {
	case EXPRESSION_TYPE_IDENTIFIER:
		identifier := e.Data.(*ExpressionIdentifer)
		d := block.SearchByName(identifier.Name)
		if d == nil {
			*errs = append(*errs, fmt.Errorf("%s '%s' not found",
				errMsgPrefix(e.Pos), identifier.Name))
			return nil
		}
		switch d.(type) {
		case *VariableDefinition:
			t := d.(*VariableDefinition)
			identifier.Var = t
			tt := identifier.Var.Typ.Clone()
			tt.Pos = e.Pos
			return tt
		default:
			*errs = append(*errs, fmt.Errorf("%s identifier named '%s' is not variable",
				errMsgPrefix(e.Pos), identifier.Name))
			return nil
		}
	case EXPRESSION_TYPE_INDEX:
		return e.checkIndexExpression(block, errs)
	case EXPRESSION_TYPE_DOT:
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
		if t.Typ == VARIABLE_TYPE_OBJECT {
			field, err := t.Class.accessField(dot.Name, false)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			}
			dot.Field = field
			if field != nil {
				if field.IsStatic() {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is static,should access by class",
						errMsgPrefix(e.Pos), dot.Name))
				}
				// not this and private
				if dot.Expression.IsThis() == false && field.IsPrivate() {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is private",
						errMsgPrefix(e.Pos), dot.Name))
				}
				tt := field.Typ.Clone()
				tt.Pos = e.Pos
				return tt
			}
			return nil
		} else if t.Typ == VARIABLE_TYPE_CLASS {
			field, err := t.Class.accessField(dot.Name, false)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			}
			dot.Field = field
			if field != nil {
				if field.IsStatic() == false {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is not static,should access by instance",
						errMsgPrefix(e.Pos), dot.Name))
				}
				tt := field.Typ.Clone()
				tt.Pos = e.Pos
				return tt
			}
			return nil
		} else if t.Typ == VARIABLE_TYPE_PACKAGE {
			if false == t.Package.Block.nameExists(dot.Name) {
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' not found",
					errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				return nil
			}
			variable := t.Package.Block.SearchByName(dot.Name)
			if vd, ok := variable.(*VariableDefinition); ok == false {
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' is not variable",
					errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				return nil
			} else {
				if vd.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 {
					*errs = append(*errs, fmt.Errorf("%s '%s.%s' is private",
						errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				}
				dot.PackageVariable = vd
				tt := vd.Typ.Clone()
				tt.Pos = e.Pos
				return tt
			}
		} else {
			*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
				errMsgPrefix(e.Pos),
				dot.Expression.OpName()))
			return nil
		}
	default:
		*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
			errMsgPrefix(e.Pos),
			e.OpName()))
		return nil
	}
	return nil
}
