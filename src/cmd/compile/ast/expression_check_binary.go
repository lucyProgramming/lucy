package ast

import (
	"fmt"
)

func (e *Expression) checkBinaryExpression(block *Block, errs *[]error) (result *VariableType) {
	bin := e.Data.(*ExpressionBinary)
	ts1, es := bin.Left.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	ts2, es := bin.Right.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	var err error
	t1, err := e.mustBeOneValueContext(ts1)
	if err != nil {
		*errs = append(*errs, err)
	}
	t2, err := e.mustBeOneValueContext(ts2)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t1 == nil || t2 == nil {
		var tt *VariableType
		if t1 != nil {
			tt = t1.Clone()
			tt.Pos = e.Pos
			return tt
		}
		if t2 != nil {
			tt = t2.Clone()
			tt.Pos = e.Pos
			return tt
		}
		return tt
	}
	// &&  ||
	if e.Typ == EXPRESSION_TYPE_LOGICAL_OR ||
		EXPRESSION_TYPE_LOGICAL_AND == e.Typ {
		if t1.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression on left,but '%s'",
				errMsgPrefix(bin.Left.Pos),
				t1.TypeString()))
		}
		if t2.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression on right,but '%s'",
				errMsgPrefix(bin.Right.Pos),
				t2.TypeString()))
		}
		result = &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
			Pos: e.Pos,
		}
		return result
	}
	// & |
	if e.Typ == EXPRESSION_TYPE_OR ||
		EXPRESSION_TYPE_AND == e.Typ ||
		EXPRESSION_TYPE_XOR == e.Typ {
		if !t1.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression",
				errMsgPrefix(bin.Left.Pos)))
		}
		if !t2.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression",
				errMsgPrefix(bin.Right.Pos)))
		}
		if t1.IsNumber() && t2.IsNumber() {
			if t1.Equal(t2) { //force to equal
				*errs = append(*errs, fmt.Errorf("%s cannot apply '%s' on '%s' and '%s'",
					errMsgPrefix(e.Pos), e.OpName(),
					t1.TypeString(),
					t2.TypeString()))
			}
		}
		result = t1.Clone()
		result.Pos = e.Pos
		return result
	}

	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT ||
		e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		if !t1.IsInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression,but '%s'",
				errMsgPrefix(bin.Left.Pos),
				t1.TypeString()))
		}
		if !t2.IsInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression,but '%s'",
				errMsgPrefix(bin.Right.Pos),
				t2.TypeString()))
		} else { // integer
			if t2.Typ == VARIABLE_TYPE_LONG {
				bin.Right.ConvertToNumber(VARIABLE_TYPE_INT)
			}
		}
		result = t1.Clone()
		result.Pos = e.Pos
		return result
	}
	if e.Typ == EXPRESSION_TYPE_EQ ||
		e.Typ == EXPRESSION_TYPE_NE ||
		e.Typ == EXPRESSION_TYPE_GE ||
		e.Typ == EXPRESSION_TYPE_GT ||
		e.Typ == EXPRESSION_TYPE_LE ||
		e.Typ == EXPRESSION_TYPE_LT {
		//number
		switch t1.Typ {
		case VARIABLE_TYPE_BOOL:
			if t2.Typ != VARIABLE_TYPE_BOOL || (e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE) {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VARIABLE_TYPE_ENUM:
			if t1.Equal(t2) == false || (e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE) {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VARIABLE_TYPE_BYTE:
			fallthrough
		case VARIABLE_TYPE_SHORT:
			fallthrough
		case VARIABLE_TYPE_INT:
			fallthrough
		case VARIABLE_TYPE_FLOAT:
			fallthrough
		case VARIABLE_TYPE_LONG:
			fallthrough
		case VARIABLE_TYPE_DOUBLE:
			if (t1.IsInteger() && t2.IsInteger()) ||
				(t1.IsFloat() && t2.IsFloat()) {
				if t1.Equal(t2) == false {
					if bin.Left.IsLiteral() == false && bin.Right.IsLiteral() == false {
						*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
					} else {
						if bin.Left.IsLiteral() {
							bin.Left.ConvertToNumber(t2.Typ)
						} else {
							bin.Right.ConvertToNumber(t1.Typ)
						}
					}
				}
			} else {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VARIABLE_TYPE_STRING:
			if t1.Equal(t2) == false {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VARIABLE_TYPE_NULL:
			if t2.IsPointer() == false || (e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE) {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and '%s'",
					errMsgPrefix(e.Pos),
					e.OpName(),
					t2.TypeString()))
			}
		case VARIABLE_TYPE_MAP:
			fallthrough
		case VARIABLE_TYPE_JAVA_ARRAY:
			fallthrough
		case VARIABLE_TYPE_ARRAY:
			fallthrough
		case VARIABLE_TYPE_OBJECT:
			if t1.Equal(t2) == false || (e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE) {
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
		t := &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
			Pos: e.Pos,
		}
		return t
	}
	if e.Typ == EXPRESSION_TYPE_ADD ||
		e.Typ == EXPRESSION_TYPE_SUB ||
		e.Typ == EXPRESSION_TYPE_MUL ||
		e.Typ == EXPRESSION_TYPE_DIV ||
		e.Typ == EXPRESSION_TYPE_MOD {
		//check string first
		if t1.Typ == VARIABLE_TYPE_STRING || t2.Typ == VARIABLE_TYPE_STRING { // string is always ok
			if e.Typ != EXPRESSION_TYPE_ADD {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
			result = &VariableType{}
			result.Typ = VARIABLE_TYPE_STRING
			result.Pos = e.Pos
			return result
		}
		if (t1.IsInteger() && t2.IsInteger()) ||
			(t1.IsFloat() && t2.IsFloat()) {
			if t1.Equal(t2) == false {
				if bin.Left.IsLiteral() == false && bin.Right.IsLiteral() == false {
					*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
				} else {
					if bin.Left.IsLiteral() {
						bin.Left.ConvertToNumber(t2.Typ)
					} else {
						bin.Right.ConvertToNumber(t1.Typ)
					}
				}
			}
		} else {
			*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
		}

		result = t1.Clone()
		result.Pos = e.Pos
		return result
	}
	return nil
}
