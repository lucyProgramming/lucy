package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) getLeftValue(block *Block) (t *VariableType, errs []error) {
	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_IDENTIFIER:
		identifier := e.Data.(*ExpressionIdentifer)
		d := block.SearchByName(identifier.Name)
		if d == nil {
			return nil, []error{fmt.Errorf("%s '%s' not found",
				errMsgPrefix(e.Pos), identifier.Name)}
		}
		switch d.(type) {
		case *VariableDefinition:
			t := d.(*VariableDefinition)
			identifier.Var = t
			tt := identifier.Var.Typ.Clone()
			tt.Pos = e.Pos
			return tt, nil
		default:
			errs = append(errs, fmt.Errorf("%s identifier named '%s' is not variable",
				errMsgPrefix(e.Pos), identifier.Name))
			return nil, []error{}
		}
	case EXPRESSION_TYPE_INDEX:
		return e.checkIndexExpression(block, &errs), errs
	case EXPRESSION_TYPE_DOT:
		dot := e.Data.(*ExpressionDot)
		ts, es := dot.Expression.check(block)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		t, err := e.mustBeOneValueContext(ts)
		if err != nil {
			errs = append(errs, err)
		}
		if t == nil {
			return nil, errs
		}
		if t.Typ == VARIABLE_TYPE_OBJECT {
			field, err := t.Class.accessField(dot.Name, false)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			}
			dot.Field = field
			if field != nil {
				if field.IsStatic() {
					errs = append(errs, fmt.Errorf("%s field '%s' is static,should access by class",
						errMsgPrefix(e.Pos), dot.Name))
				}
				// not this and private
				if dot.Expression.isThis() == false && field.IsPrivate() {
					errs = append(errs, fmt.Errorf("%s field '%s' is private",
						errMsgPrefix(e.Pos), dot.Name))
				}
				tt := field.Typ.Clone()
				tt.Pos = e.Pos
				return tt, errs
			}
			return nil, errs
		} else if t.Typ == VARIABLE_TYPE_CLASS {
			field, err := t.Class.accessField(dot.Name, false)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			}
			dot.Field = field
			if field != nil {
				if field.IsStatic() == false {
					errs = append(errs, fmt.Errorf("%s field '%s' is not static,should access by instance",
						errMsgPrefix(e.Pos), dot.Name))
				}
				tt := field.Typ.Clone()
				tt.Pos = e.Pos
				return tt, errs
			}
			return nil, errs
		} else if t.Typ == VARIABLE_TYPE_PACKAGE {
			variable := t.Package.Block.SearchByName(dot.Name)
			if nil == variable {
				errs = append(errs, fmt.Errorf("%s '%s.%s' not found",
					errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				return nil, errs
			}
			if vd, ok := variable.(*VariableDefinition); ok == false {
				errs = append(errs, fmt.Errorf("%s '%s.%s' is not varible",
					errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				return nil, errs
			} else {
				if vd.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 {
					errs = append(errs, fmt.Errorf("%s '%s.%s' is private",
						errMsgPrefix(e.Pos), t.Package.Name, dot.Name))
				}
				dot.PackageVariable = vd
				tt := vd.Typ.Clone()
				tt.Pos = e.Pos
				return tt, errs
			}
		} else {
			errs = append(errs, fmt.Errorf("%s '%s' cannot be used as left value",
				errMsgPrefix(e.Pos),
				dot.Expression.OpName()))
			return nil, errs
		}
	default:
		errs = append(errs, fmt.Errorf("%s '%s' cannot be used as left value",
			errMsgPrefix(e.Pos),
			e.OpName()))
		return nil, errs
	}
	return nil, errs
}
