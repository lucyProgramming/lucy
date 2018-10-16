package ast

import (
	"fmt"
)

func (e *Expression) checkBinaryExpression(block *Block, errs *[]error) (result *Type) {
	bin := e.Data.(*ExpressionBinary)
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
	if e.Type == ExpressionTypeLogicalOr ||
		ExpressionTypeLogicalAnd == e.Type {
		result = &Type{
			Type: VariableTypeBool,
			Pos:  e.Pos,
		}
		if left == nil || right == nil {
			return result
		}
		if left.Type != VariableTypeBool || right.Type != VariableTypeBool {
			*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
		}
		return result
	}
	// & |
	if e.Type == ExpressionTypeOr ||
		ExpressionTypeAnd == e.Type ||
		ExpressionTypeXor == e.Type {
		if left == nil || right == nil {
			if left != nil && left.IsNumber() {
				result := left.Clone()
				result.Pos = e.Pos
				return result
			}
			if right != nil && right.IsNumber() {
				result := right.Clone()
				result.Pos = e.Pos
				return result
			}
			return nil
		}
		if left.IsInteger() == false || left.assignAble(errs, right) == false {
			*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
		}
		result = left.Clone()
		result.Pos = e.Pos
		return result
	}
	if e.Type == ExpressionTypeLsh ||
		e.Type == ExpressionTypeRsh {
		if left == nil || right == nil {
			if left != nil && left.IsNumber() {
				result := left.Clone()
				result.Pos = e.Pos
				return result
			}
			return nil
		}
		if false == left.IsInteger() || right.IsInteger() == false {
			*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
		}
		if right.Type == VariableTypeLong {
			bin.Right.convertToNumber(VariableTypeInt)
		}
		result = left.Clone()
		result.Pos = e.Pos
		return result
	}
	if e.Type == ExpressionTypeEq ||
		e.Type == ExpressionTypeNe ||
		e.Type == ExpressionTypeGe ||
		e.Type == ExpressionTypeGt ||
		e.Type == ExpressionTypeLe ||
		e.Type == ExpressionTypeLt {
		result = &Type{
			Type: VariableTypeBool,
			Pos:  e.Pos,
		}
		if left == nil || right == nil {
			return result
		}
		//number
		switch left.Type {
		case VariableTypeBool:
			if right.Type != VariableTypeBool || e.isEqOrNe() == false {
				*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
			}
		case VariableTypeEnum:
			if left.assignAble(errs, right) == false {
				*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
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
			if (left.IsInteger() && right.IsInteger()) ||
				(left.IsFloat() && right.IsFloat()) {
				if left.assignAble(errs, right) == false {
					switch {
					case bin.Left.isLiteral() == false && bin.Right.isLiteral() == false:
						*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
					case bin.Left.isLiteral() && bin.Right.isLiteral():
						if bin.Right.Value.Type > bin.Left.Value.Type {
							bin.Left.convertToNumber(bin.Right.Value.Type)
						} else {
							bin.Right.convertToNumber(bin.Left.Value.Type)
						}
					default:
						if bin.Left.isLiteral() {
							bin.Left.convertToNumber(right.Type)
						} else {
							bin.Right.convertToNumber(left.Type)
						}
					}
				}
			} else {
				*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
			}
		case VariableTypeString:
			if right.Type == VariableTypeNull {
				if e.Type != ExpressionTypeEq && ExpressionTypeNe != e.Type {
					*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))

				}
			} else {
				if right.Type != VariableTypeString {
					*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
				}
			}
		case VariableTypeNull:
			if right.IsPointer() == false || e.isEqOrNe() == false {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and '%s'",
					errMsgPrefix(e.Pos),
					e.Description,
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
			if left.assignAble(errs, right) == false || e.isEqOrNe() == false {
				*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
			}
		default:
			*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
		}
		return result
	}
	// + - * / %
	if e.Type == ExpressionTypeAdd ||
		e.Type == ExpressionTypeSub ||
		e.Type == ExpressionTypeMul ||
		e.Type == ExpressionTypeDiv ||
		e.Type == ExpressionTypeMod {
		if left == nil || right == nil {
			if left != nil {
				result := left.Clone()
				result.Pos = e.Pos
				return result
			}
			if right != nil {
				result := right.Clone()
				result.Pos = e.Pos
				return result
			}
			return nil
		}
		//check string first
		if left.Type == VariableTypeString ||
			right.Type == VariableTypeString { // string is always ok
			if e.Type != ExpressionTypeAdd {
				*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
			}
			result = &Type{}
			result.Type = VariableTypeString
			result.Pos = e.Pos
			return result
		}
		if (left.IsInteger() && right.IsInteger()) ||
			(left.IsFloat() && right.IsFloat()) {
			if left.assignAble(errs, right) == false {
				switch {
				case bin.Left.isLiteral() == false && bin.Right.isLiteral() == false:
					*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
				case bin.Left.isLiteral() && bin.Right.isLiteral():
					if bin.Right.Value.Type > bin.Left.Value.Type {
						bin.Left.convertToNumber(bin.Right.Value.Type)
					} else {
						bin.Right.convertToNumber(bin.Left.Value.Type)
					}
				default:
					if bin.Left.isLiteral() {
						bin.Left.convertToNumber(right.Type)
					} else {
						bin.Right.convertToNumber(left.Type)
					}
				}
			}
		} else {
			*errs = append(*errs, e.makeWrongOpErr(left.TypeString(), right.TypeString()))
		}
		result = left.Clone()
		result.Pos = e.Pos
		return result
	}
	return nil
}
