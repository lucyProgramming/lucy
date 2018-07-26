package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) getLeftValue(block *Block, errs *[]error) (result *Type) {
	switch e.Type {
	case ExpressionTypeIdentifier:
		identifier := e.Data.(*ExpressionIdentifier)
		d, _ := block.searchIdentifier(identifier.Name)
		if d == nil {
			*errs = append(*errs, fmt.Errorf("%s '%s' not found",
				errMsgPrefix(e.Pos), identifier.Name))
			return nil
		}
		switch d.(type) {
		case *Variable:
			if identifier.Name == THIS {
				*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
					errMsgPrefix(e.Pos), THIS))
			}
			t := d.(*Variable)
			identifier.Variable = t
			result = identifier.Variable.Type.Clone()
			result.Pos = e.Pos
			return result
		default:
			*errs = append(*errs, fmt.Errorf("%s identifier named '%s' is not variable",
				errMsgPrefix(e.Pos), identifier.Name))
			return nil
		}
	case ExpressionTypeIndex:
		result = e.checkIndexExpression(block, errs)
		return result
	case ExpressionTypeSelection:
		selection := e.Data.(*ExpressionSelection)
		object, es := selection.Expression.checkSingleValueContextExpression(block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if object == nil {
			return nil
		}
		switch object.Type {
		case VariableTypeObject:
			field, err := object.Class.accessField(selection.Name, false)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			}
			selection.Field = field
			if field != nil {
				if field.IsStatic() {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is static,should access by class",
						errMsgPrefix(e.Pos), selection.Name))
				}
				// not this and private
				if selection.Expression.isThis() == false {
					if (selection.Expression.Value.Class.LoadFromOutSide && field.IsPublic() == false) ||
						(selection.Expression.Value.Class.LoadFromOutSide == false && field.IsPrivate()) {
						*errs = append(*errs, fmt.Errorf("%s field '%s' is private",
							errMsgPrefix(e.Pos), selection.Name))
					}
				}
				result = field.Type.Clone()
				result.Pos = e.Pos
				return result
			}
			return nil
		case VariableTypeClass:
			field, err := object.Class.accessField(selection.Name, false)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
			}
			selection.Field = field
			if field != nil {
				if field.IsStatic() == false {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is not static,should access by instance",
						errMsgPrefix(e.Pos), selection.Name))
				}
				if block.InheritedAttribute.Class != selection.Expression.Value.Class {
					if (selection.Expression.Value.Class.LoadFromOutSide && field.IsPublic() == false) ||
						(selection.Expression.Value.Class.LoadFromOutSide == false && field.IsPrivate()) {
						*errs = append(*errs, fmt.Errorf("%s field '%s' is private",
							errMsgPrefix(e.Pos), selection.Name))
					}
				}
				result = field.Type.Clone()
				result.Pos = e.Pos
				return result
			}
			return nil
		case VariableTypePackage:
			variable, exists := object.Package.Block.NameExists(selection.Name)
			if exists == false {
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' not found",
					errMsgPrefix(e.Pos), object.Package.Name, selection.Name))
				return nil
			}
			if v, ok := variable.(*Variable); ok && v != nil {
				if v.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 && object.Package.Name != PackageBeenCompile.Name {
					*errs = append(*errs, fmt.Errorf("%s '%s.%s' is private",
						errMsgPrefix(e.Pos), object.Package.Name, selection.Name))
				}
				selection.PackageVariable = v
				result = v.Type.Clone()
				result.Pos = e.Pos
				return result
			} else {
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' is not variable",
					errMsgPrefix(e.Pos), object.Package.Name, selection.Name))
				return nil
			}
		default:
			*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
				errMsgPrefix(e.Pos),
				selection.Expression.OpName()))
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
