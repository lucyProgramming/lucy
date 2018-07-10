package ast

import (
	"fmt"
)

func (e *Expression) checkBinaryExpression(block *Block, errs *[]error) (result *Type) {
	bin := e.Data.(*ExpressionBinary)
	t1, es := bin.Left.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t2, es := bin.Right.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}

	// &&  ||
	if e.Type == ExpressionTypeLogicalOr ||
		ExpressionTypeLogicalAnd == e.Type {
		result = &Type{
			Type: VariableTypeBool,
			Pos:  e.Pos,
		}
		if t1 == nil || t2 == nil {
			return result
		}
		if t1.Type != VariableTypeBool || t2.Type != VariableTypeBool {
			*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
		}

		return result
	}
	// & |
	if e.Type == ExpressionTypeOr ||
		ExpressionTypeAnd == e.Type ||
		ExpressionTypeXor == e.Type {
		if t1 == nil || t2 == nil {
			if t1 != nil {
				tt := t1.Clone()
				tt.Pos = e.Pos
				return tt
			}
			if t2 != nil {
				tt := t2.Clone()
				tt.Pos = e.Pos
				return tt
			}
			return nil
		}
		if t1.IsInteger() == false || t1.Equal(errs, t2) == false {
			*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
		}
		result = t1.Clone()
		result.Pos = e.Pos
		return result
	}

	if e.Type == ExpressionTypeLsh ||
		e.Type == ExpressionTypeRsh {
		if t1 == nil || t2 == nil {
			if t1 != nil {
				tt := t1.Clone()
				tt.Pos = e.Pos
				return tt
			}
			if t2 != nil {
				tt := t2.Clone()
				tt.Pos = e.Pos
				return tt
			}
			return nil
		}
		if false == t1.IsInteger() || t2.IsInteger() == false {
			*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
		} else {
			if t2.Type == VariableTypeLong {
				bin.Right.ConvertToNumber(VariableTypeInt)
			}
		}
		result = t1.Clone()
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
		if t1 == nil || t2 == nil {
			return result
		}
		//number
		switch t1.Type {
		case VariableTypeBool:
			if t2.Type != VariableTypeBool || (e.Type != ExpressionTypeEq && e.Type != ExpressionTypeNe) {
				*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VariableTypeEnum:
			if t1.Equal(errs, t2) == false || (e.Type != ExpressionTypeEq && e.Type != ExpressionTypeNe) {
				*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VariableTypeByte:
			fallthrough
		case VariableTypeShort:
			fallthrough
		case VariableTypeInt:
			fallthrough
		case VariableTypeFloat:
			fallthrough
		case VariableTypeLong:
			fallthrough
		case VariableTypeDouble:
			if (t1.IsInteger() && t2.IsInteger()) ||
				(t1.IsFloat() && t2.IsFloat()) {
				if t1.Equal(errs, t2) == false {
					if bin.Left.IsLiteral() == false && bin.Right.IsLiteral() == false {
						*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
					} else {
						if bin.Left.IsLiteral() {
							bin.Left.ConvertToNumber(t2.Type)
						} else {
							bin.Right.ConvertToNumber(t1.Type)
						}
					}
				}
			} else {
				*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VariableTypeString:
			if t2.Type == VariableTypeNull {
				if e.Type != ExpressionTypeEq && ExpressionTypeNe != e.Type {
					*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))

				}
			} else {
				if t2.Type != VariableTypeString {
					*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
				}
			}

		case VariableTypeNull:
			if t2.IsPointer() == false || (e.Type != ExpressionTypeEq && e.Type != ExpressionTypeNe) {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and '%s'",
					errMsgPrefix(e.Pos),
					e.OpName(),
					t2.TypeString()))
			}
		case VariableTypeMap:
			fallthrough
		case VariableTypeJavaArray:
			fallthrough
		case VariableTypeArray:
			fallthrough
		case VariableTypeObject:
			if t1.Equal(errs, t2) == false || (e.Type != ExpressionTypeEq && e.Type != ExpressionTypeNe) {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'",
					errMsgPrefix(e.Pos),
					e.OpName(),
					t1.TypeString(),
					t2.TypeString()))
			}
		default:
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'",
				errMsgPrefix(e.Pos),
				e.OpName(),
				t1.TypeString(),
				t2.TypeString()))
		}

		return result
	}
	//
	if e.Type == ExpressionTypeAdd ||
		e.Type == ExpressionTypeSub ||
		e.Type == ExpressionTypeMul ||
		e.Type == ExpressionTypeDiv ||
		e.Type == ExpressionTypeMod {
		if t1 == nil || t2 == nil {
			if t1 != nil {
				tt := t1.Clone()
				tt.Pos = e.Pos
				return tt
			}
			if t2 != nil {
				tt := t2.Clone()
				tt.Pos = e.Pos
				return tt
			}
			return nil
		}
		//check string first
		if t1.Type == VariableTypeString || t2.Type == VariableTypeString { // string is always ok
			if e.Type != ExpressionTypeAdd {
				*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
			}
			result = &Type{}
			result.Type = VariableTypeString
			result.Pos = e.Pos
			return result
		}
		if (t1.IsInteger() && t2.IsInteger()) ||
			(t1.IsFloat() && t2.IsFloat()) {
			if t1.Equal(errs, t2) == false {
				if bin.Left.IsLiteral() == false && bin.Right.IsLiteral() == false {
					*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
				} else {
					if bin.Left.IsLiteral() {
						bin.Left.ConvertToNumber(t2.Type)
					} else {
						bin.Right.ConvertToNumber(t1.Type)
					}
				}
			}
		} else {
			*errs = append(*errs, e.mkWrongOpErr(t1.TypeString(), t2.TypeString()))
		}
		result = t1.Clone()
		result.Pos = e.Pos
		return result
	}
	return nil
}
