package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) getLeftValue(block *Block, errs *[]error) (result *Type) {
	switch e.Type {
	case ExpressionTypeIdentifier:
		identifier := e.Data.(*ExpressionIdentifier)
		if identifier.Name == NoNameIdentifier {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as left value",
				errMsgPrefix(e.Pos), identifier.Name))
			return nil
		}
		d, err := block.searchIdentifier(e.Pos, identifier.Name)
		if err != nil {
			*errs = append(*errs, err)
		}
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
			v := d.(*Variable)
			identifier.Variable = v
			result = identifier.Variable.Type.Clone()
			result.Pos = e.Pos
			e.Value = result
			return result
		default:
			*errs = append(*errs, fmt.Errorf("%s identifier '%s' is '%s' , cannot be used as left value",
				errMsgPrefix(e.Pos), identifier.Name, block.searchedIdentifierIs(d)))
			return nil
		}
	case ExpressionTypeIndex:
		result = e.checkIndexExpression(block, errs)
		e.Value = result
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
		case VariableTypeDynamicSelector:
			if selection.Name == SUPER {
				*errs = append(*errs, fmt.Errorf("%s access '%s' at '%s' not allow",
					errMsgPrefix(e.Pos), SUPER, object.TypeString()))
				return nil
			}
			field, err := object.Class.accessField(e.Pos, selection.Name, false)
			if err != nil {
				*errs = append(*errs, err)
			}
			if field != nil {
				selection.Field = field
				result = field.Type.Clone()
				result.Pos = e.Pos
				e.Value = result
				return result
			} else {
				return nil
			}
		case VariableTypeObject:
			field, err := object.Class.accessField(e.Pos, selection.Name, false)
			if err != nil {
				*errs = append(*errs, err)
			}
			selection.Field = field
			if field != nil {
				if field.IsStatic() {
					*errs = append(*errs, fmt.Errorf("%s field '%s' is static,should access by class",
						errMsgPrefix(e.Pos), selection.Name))
				}
				// not this and private
				if selection.Expression.IsIdentifier(THIS) == false {
					if (selection.Expression.Value.Class.LoadFromOutSide && field.IsPublic() == false) ||
						(selection.Expression.Value.Class.LoadFromOutSide == false && field.IsPrivate()) {
						*errs = append(*errs, fmt.Errorf("%s field '%s' is private",
							errMsgPrefix(e.Pos), selection.Name))
					}
				}
				result = field.Type.Clone()
				result.Pos = e.Pos
				e.Value = result
				return result
			}
			return nil
		case VariableTypeClass:
			field, err := object.Class.accessField(e.Pos, selection.Name, false)
			if err != nil {
				*errs = append(*errs, err)
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
				e.Value = result
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
			switch variable.(type) {
			case *Variable:
				v := variable.(*Variable)
				if v.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 && object.Package.Name != PackageBeenCompile.Name {
					*errs = append(*errs, fmt.Errorf("%s '%s.%s' is private",
						errMsgPrefix(e.Pos), object.Package.Name, selection.Name))
				}
				selection.PackageVariable = v
				result = v.Type.Clone()
				result.Pos = e.Pos
				e.Value = result
				return result
			default:
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' is not variable",
					errMsgPrefix(e.Pos), object.Package.Name, selection.Name))
				return nil
			}
		case VariableTypeMagicFunction:
			v := object.Function.Type.searchName(selection.Name)
			if v == nil {
				err := fmt.Errorf("%s '%s' not found", errMsgPrefix(e.Pos), selection.Name)
				*errs = append(*errs, err)
				return nil
			}
			e.Value = v.Type.Clone()
			e.Value.Pos = e.Pos
			e.Type = ExpressionTypeIdentifier
			identifier := &ExpressionIdentifier{}
			identifier.Name = selection.Name
			identifier.Variable = v
			e.Data = identifier
			result := v.Type.Clone()
			result.Pos = e.Pos
			return result
		default:
			*errs = append(*errs, fmt.Errorf("%s cannot access '%s' on '%s'",
				errMsgPrefix(e.Pos), selection.Name, object.TypeString()))
			return nil
		}
	default:
		*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
			errMsgPrefix(e.Pos),
			e.Description))
		return nil
	}
	return nil
}
