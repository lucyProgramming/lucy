package ast

import (
	"fmt"
)

func (this *Expression) checkBinaryExpression(block *Block, errs *[]error) (result *Type) {
	bin := this.Data.(*ExpressionBinary)
	left, es := bin.Left.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	right, es := bin.Right.checkSingleValueContextExpression(block)
	*errs = append(*errs, es...)
	if left != nil {
		if err := left.rightValueValid(); err != nil {
			*errs = append(*errs, err)
			return nil
		}
	}
	if right != nil {
		if err := right.rightValueValid(); err != nil {
			*errs = append(*errs, err)
			return nil
		}
	}

	// &&  ||
	if this.Type == ExpressionTypeLogicalOr ||
		this.Type == ExpressionTypeLogicalAnd {
		result = &Type{
			Type: VariableTypeBool,
			Pos:  this.Pos,
		}
		if left == nil || right == nil {
			return result
		}
		if left.Type != VariableTypeBool ||
			right.Type != VariableTypeBool {
			*errs = append(*errs, this.binaryWrongOpErr())
		}
		return result
	}
	// & |
	if this.Type == ExpressionTypeOr ||
		ExpressionTypeAnd == this.Type ||
		ExpressionTypeXor == this.Type {
		if left == nil || right == nil {
			if left != nil && left.IsNumber() {
				result := left.Clone()
				result.Pos = this.Pos
				return result
			}
			if right != nil && right.IsNumber() {
				result := right.Clone()
				result.Pos = this.Pos
				return result
			}
			return nil
		}
		if left.isInteger() == false || left.assignAble(errs, right) == false {
			*errs = append(*errs, this.binaryWrongOpErr())
		}
		result = left.Clone()
		result.Pos = this.Pos
		return result
	}
	if this.Type == ExpressionTypeLsh ||
		this.Type == ExpressionTypeRsh {
		if left == nil || right == nil {
			if left != nil && left.IsNumber() {
				result := left.Clone()
				result.Pos = this.Pos
				return result
			}
			return nil
		}
		if false == left.isInteger() ||
			right.isInteger() == false {
			*errs = append(*errs, this.binaryWrongOpErr())
		}
		if right.Type == VariableTypeLong {
			bin.Right.convertToNumberType(VariableTypeInt)
		}
		result = left.Clone()
		result.Pos = this.Pos
		return result
	}
	if this.Type == ExpressionTypeEq ||
		this.Type == ExpressionTypeNe ||
		this.Type == ExpressionTypeGe ||
		this.Type == ExpressionTypeGt ||
		this.Type == ExpressionTypeLe ||
		this.Type == ExpressionTypeLt {
		result = &Type{
			Type: VariableTypeBool,
			Pos:  this.Pos,
		}
		if left == nil || right == nil {
			return result
		}
		//number
		switch left.Type {
		case VariableTypeBool:
			if right.Type != VariableTypeBool || this.isEqOrNe() == false {
				*errs = append(*errs, this.binaryWrongOpErr())
			}
		case VariableTypeEnum:
			if left.assignAble(errs, right) == false {
				*errs = append(*errs, this.binaryWrongOpErr())
			}
		case VariableTypeByte:
			fallthrough
		case VariableTypeShort:
			fallthrough
		case VariableTypeChar:
			fallthrough
		case VariableTypeInt:
			fallthrough
		case VariableTypeFloat:
			fallthrough
		case VariableTypeLong:
			fallthrough
		case VariableTypeDouble:
			if (left.isInteger() && right.isInteger()) ||
				(left.isFloat() && right.isFloat()) {
				if left.assignAble(errs, right) == false {
					if left.Type < right.Type {
						bin.Left.convertToNumberType(right.Type)
					} else {
						bin.Right.convertToNumberType(left.Type)
					}
				}
			} else {
				*errs = append(*errs, this.binaryWrongOpErr())
			}
		case VariableTypeString:
			if right.Type == VariableTypeNull {
				if this.Type != ExpressionTypeEq && ExpressionTypeNe != this.Type {
					*errs = append(*errs, this.binaryWrongOpErr())

				}
			} else {
				if right.Type != VariableTypeString {
					*errs = append(*errs, this.binaryWrongOpErr())
				}
			}
		case VariableTypeNull:
			if right.IsPointer() == false || this.isEqOrNe() == false {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and '%s'",
					this.Pos.ErrMsgPrefix(),
					this.Op,
					right.TypeString()))
			}
		case VariableTypeMap:
			fallthrough
		case VariableTypeJavaArray:
			fallthrough
		case VariableTypeArray:
			fallthrough
		case VariableTypeObject:
			fallthrough
		case VariableTypeFunction:
			if left.assignAble(errs, right) == false || this.isEqOrNe() == false {
				*errs = append(*errs, this.binaryWrongOpErr())
			}
		default:
			*errs = append(*errs, this.binaryWrongOpErr())
		}
		return result
	}
	// + - * / %
	if this.Type == ExpressionTypeAdd ||
		this.Type == ExpressionTypeSub ||
		this.Type == ExpressionTypeMul ||
		this.Type == ExpressionTypeDiv ||
		this.Type == ExpressionTypeMod {
		if left == nil || right == nil {
			if left != nil {
				result := left.Clone()
				result.Pos = this.Pos
				return result
			}
			if right != nil {
				result := right.Clone()
				result.Pos = this.Pos
				return result
			}
			return nil
		}
		//check string first
		if left.Type == VariableTypeString ||
			right.Type == VariableTypeString { // string is always ok
			if this.Type != ExpressionTypeAdd {
				*errs = append(*errs, this.binaryWrongOpErr())
			}
			result = &Type{}
			result.Type = VariableTypeString
			result.Pos = this.Pos
			return result
		}
		if (left.isInteger() && right.isInteger()) ||
			(left.isFloat() && right.isFloat()) {
			if left.assignAble(errs, right) == false {
				if left.Type < right.Type {
					bin.Left.convertToNumberType(right.Type)
				} else {
					bin.Right.convertToNumberType(left.Type)
				}
			}
		} else {
			*errs = append(*errs, this.binaryWrongOpErr())
		}
		result = left.Clone()
		result.Pos = this.Pos
		return result
	}
	return nil
}
