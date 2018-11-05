package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *Expression) getLeftValue(block *Block, errs *[]error) (result *Type) {
	switch this.Type {
	case ExpressionTypeIdentifier:
		identifier := this.Data.(*ExpressionIdentifier)
		if identifier.Name == UnderScore {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as left value",
				this.Pos.ErrMsgPrefix(), identifier.Name))
			return nil
		}
		if identifier.Name == ThisPointerName {
			*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
				this.Pos.ErrMsgPrefix(), ThisPointerName))
		}
		isCaptureVar := false
		d, err := block.searchIdentifier(this.Pos, identifier.Name, &isCaptureVar)
		if err != nil {
			*errs = append(*errs, err)
			return nil
		}
		if d == nil {
			*errs = append(*errs, fmt.Errorf("%s '%s' not found",
				this.Pos.ErrMsgPrefix(), identifier.Name))
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
			result.Pos = this.Pos
			this.Value = result
			return result
		default:
			*errs = append(*errs, fmt.Errorf("%s identifier '%s' is '%s' , cannot be used as left value",
				this.Pos.ErrMsgPrefix(), identifier.Name, block.identifierIsWhat(d)))
			return nil
		}
	case ExpressionTypeIndex:
		result = this.checkIndexExpression(block, errs)
		this.Value = result
		return result
	case ExpressionTypeSelection:
		selection := this.Data.(*ExpressionSelection)
		object, es := selection.Expression.checkSingleValueContextExpression(block)
		*errs = append(*errs, es...)
		if object == nil {
			return nil
		}
		switch object.Type {
		case VariableTypeDynamicSelector:
			if selection.Name == SUPER {
				*errs = append(*errs, fmt.Errorf("%s access '%s' at '%s' not allow",
					this.Pos.ErrMsgPrefix(), SUPER, object.TypeString()))
				return nil
			}
			field, err := object.Class.getField(this.Pos, selection.Name, false)
			if err != nil {
				*errs = append(*errs, err)
			}
			if field == nil {
				return nil
			}
			selection.Field = field
			result = field.Type.Clone()
			result.Pos = this.Pos
			this.Value = result
			return result
		case VariableTypeObject, VariableTypeClass:
			field, err := object.Class.getField(this.Pos, selection.Name, false)
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
				result.Pos = this.Pos
				this.Value = result
				return result
			}
			return nil
		case VariableTypePackage:
			variable, exists := object.Package.Block.NameExists(selection.Name)
			if exists == false {
				*errs = append(*errs, fmt.Errorf("%s '%s.%s' not found",
					this.Pos.ErrMsgPrefix(), object.Package.Name, selection.Name))
				return nil
			}
			switch variable.(type) {
			case *Variable:
				v := variable.(*Variable)
				if v.AccessFlags&cg.AccFieldPublic == 0 &&
					object.Package.isSame(&PackageBeenCompile) == false {
					*errs = append(*errs, fmt.Errorf("%s '%s.%s' is private",
						this.Pos.ErrMsgPrefix(), object.Package.Name, selection.Name))
				}
				selection.PackageVariable = v
				result = v.Type.Clone()
				result.Pos = this.Pos
				this.Value = result
				return result
			default:
				*errs = append(*errs, fmt.Errorf("%s '%s' is not variable",
					this.Pos.ErrMsgPrefix(), selection.Name))
				return nil
			}
		case VariableTypeMagicFunction:
			v := object.Function.Type.searchName(selection.Name)
			if v == nil {
				err := fmt.Errorf("%s '%s' not found", this.Pos.ErrMsgPrefix(), selection.Name)
				*errs = append(*errs, err)
				return nil
			}
			this.Value = v.Type.Clone()
			this.Value.Pos = this.Pos
			this.Type = ExpressionTypeIdentifier
			identifier := &ExpressionIdentifier{}
			identifier.Name = selection.Name
			identifier.Variable = v
			this.Data = identifier
			result := v.Type.Clone()
			result.Pos = this.Pos
			return result
		default:
			*errs = append(*errs, fmt.Errorf("%s cannot access '%s' on '%s'",
				this.Pos.ErrMsgPrefix(), selection.Name, object.TypeString()))
			return nil
		}
	default:
		*errs = append(*errs, fmt.Errorf("%s '%s' cannot be used as left value",
			this.Pos.ErrMsgPrefix(),
			this.Op))
		return nil
	}
	return nil
}
