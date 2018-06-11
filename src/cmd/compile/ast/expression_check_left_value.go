package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) getLeftValue(block *Block, errs *[]error) (t *VariableType) {
	switch e.Typ {
	case EXPRESSION_TYPE_IDENTIFIER:
		identifier := e.Data.(*ExpressionIdentifier)
		d, _ := block.SearchByName(identifier.Name)
		if d == nil {
			*errs = append(*errs, fmt.Errorf("%s '%s' not found",
				errMsgPrefix(e.Pos), identifier.Name))
			return nil
		}
		switch d.(type) {
		case *VariableDefinition:
			if identifier.Name == THIS {
				*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
					errMsgPrefix(e.Pos), THIS))
			}
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
	case EXPRESSION_TYPE_SELECT:
		dot := e.Data.(*ExpressionSelection)
		t, es := dot.Expression.checkSingleValueContextExpression(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
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
				if dot.Expression.isThis() == false && field.IsPrivate() {
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
			variable, exists := t.Package.Block.NameExists(dot.Name)
			if exists == false {
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' not found",
					errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				return nil
			}
			if vd, ok := variable.(*VariableDefinition); ok && vd != nil {
				if vd.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 {
					*errs = append(*errs, fmt.Errorf("%s '%s.%s' is private",
						errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				}
				dot.PackageVariable = vd
				tt := vd.Typ.Clone()
				tt.Pos = e.Pos
				return tt
			} else {
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' is not variable",
					errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				return nil
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
