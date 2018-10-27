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
				e.Pos.ErrMsgPrefix(), identifier.Name))
			return nil
		}
		if identifier.Name == THIS {
			*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
				e.Pos.ErrMsgPrefix(), THIS))
		}
		isCaptureVar := false
		d, err := block.searchIdentifier(e.Pos, identifier.Name, &isCaptureVar)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if d == nil {
			*errs = append(*errs, fmt.Errorf("%s '%s' not found",
				e.Pos.ErrMsgPrefix(), identifier.Name))
			return nil
		}
		switch d.(type) {
		case *Variable:
			v := d.(*Variable)
			if isCaptureVar {

				v.BeenCapturedAsLeftValue++
			}
			// variable is modifying , capture right value should not be ok
			// if no variable not change,after been captured, right value should ok too
			v.BeenCapturedAsLeftValue += v.BeenCapturedAsRightValue
			v.BeenCapturedAsRightValue = 0
			identifier.Variable = v
			result = identifier.Variable.Type.Clone()
			result.Pos = e.Pos
			e.Value = result
			return result
		default:
			*errs = append(*errs, fmt.Errorf("%s identifier '%s' is '%s' , cannot be used as left value",
				e.Pos.ErrMsgPrefix(), identifier.Name, block.identifierIsWhat(d)))
			return nil
		}
	case ExpressionTypeIndex:
		result = e.checkIndexExpression(block, errs)
		e.Value = result
		return result
	case ExpressionTypeSelection:
		selection := e.Data.(*ExpressionSelection)
		object, es := selection.Expression.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if object == nil {
			return nil
		}
		switch object.Type {
		case VariableTypeDynamicSelector:
			if selection.Name == SUPER {
				*errs = append(*errs, fmt.Errorf("%s access '%s' at '%s' not allow",
					e.Pos.ErrMsgPrefix(), SUPER, object.TypeString()))
				return nil
			}
			field, err := object.Class.getField(e.Pos, selection.Name, false)
			if err != nil {
				*errs = append(*errs, err)
			}
			if field == nil {
				return nil
			}
			selection.Field = field
			result = field.Type.Clone()
			result.Pos = e.Pos
			e.Value = result
			return result
		case VariableTypeObject, VariableTypeClass:
			field, err := object.Class.getField(e.Pos, selection.Name, false)
			if err != nil {
				*errs = append(*errs, err)
			}
			selection.Field = field
			if field != nil {
				err := selection.Expression.fieldAccessAble(block, field)
				if err != nil {
					*errs = append(*errs, err)
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
					e.Pos.ErrMsgPrefix(), object.Package.Name, selection.Name))
				return nil
			}
			switch variable.(type) {
			case *Variable:
				v := variable.(*Variable)
				if v.AccessFlags&cg.ACC_FIELD_PUBLIC == 0 &&
					object.Package.isSame(&PackageBeenCompile) == false {
					*errs = append(*errs, fmt.Errorf("%s '%s.%s' is private",
						e.Pos.ErrMsgPrefix(), object.Package.Name, selection.Name))
				}
				selection.PackageVariable = v
				result = v.Type.Clone()
				result.Pos = e.Pos
				e.Value = result
				return result
			default:
				*errs = append(*errs, fmt.Errorf("%s '%s' is not variable",
					e.Pos.ErrMsgPrefix(), selection.Name))
				return nil
			}
		case VariableTypeMagicFunction:
			v := object.Function.Type.searchName(selection.Name)
			if v == nil {
				err := fmt.Errorf("%s '%s' not found", e.Pos.ErrMsgPrefix(), selection.Name)
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
				e.Pos.ErrMsgPrefix(), selection.Name, object.TypeString()))
			return nil
		}
	default:
		*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
			e.Pos.ErrMsgPrefix(),
			e.Op))
		return nil
	}
	return nil
}
