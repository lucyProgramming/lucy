package ast

import (
	"fmt"
)

func (e *Expression) checkBinaryExpression(block *Block, errs *[]error) (result *VariableType) {
	bin := e.Data.(*ExpressionBinary)
	ts1, err1 := bin.Left.check(block)
	ts2, err2 := bin.Right.check(block)
	if errsNotEmpty(err1) {
		*errs = append(*errs, err1...)
	}
	if errsNotEmpty(err2) {
		*errs = append(*errs, err2...)
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
		return nil
	}
	// &&  ||
	if e.Typ == EXPRESSION_TYPE_LOGICAL_OR || EXPRESSION_TYPE_LOGICAL_AND == e.Typ {
		if t1.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression,but '%s'",
				errMsgPrefix(bin.Left.Pos),
				t1.TypeString()))
		}
		if t2.Typ != VARIABLE_TYPE_BOOL {
			*errs = append(*errs, fmt.Errorf("%s not a bool expression,but '%s'",
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
	if e.Typ == EXPRESSION_TYPE_OR || EXPRESSION_TYPE_AND == e.Typ {
		if !t1.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(bin.Left.Pos)))
		}
		if !t2.IsNumber() {
			*errs = append(*errs, fmt.Errorf("%s not a number expression", errMsgPrefix(bin.Right.Pos)))
		}
		if t1.IsNumber() && t2.IsNumber() {
			if t1.Typ != t2.Typ {
				*errs = append(*errs, fmt.Errorf("%s cannot apply '&' or '|' on '%s' and '%s'",
					errMsgPrefix(bin.Right.Pos),
					t1.TypeString(),
					t2.TypeString()))
			}
		}
		result = t1.Clone()
		result.Pos = e.Pos
		return result
	}
	if e.Typ == EXPRESSION_TYPE_LEFT_SHIFT || e.Typ == EXPRESSION_TYPE_RIGHT_SHIFT {
		if !t1.IsInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression,but '%s'",
				errMsgPrefix(bin.Left.Pos),
				t1.TypeString()))
		}
		if !t2.IsInteger() {
			*errs = append(*errs, fmt.Errorf("%s not a integer expression,but '%s'",
				errMsgPrefix(bin.Right.Pos),
				t2.TypeString()))
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
			if !t2.IsNumber() {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VARIABLE_TYPE_STRING:
			if t2.Typ != VARIABLE_TYPE_STRING {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VARIABLE_TYPE_BOOL:
			if t2.Typ == VARIABLE_TYPE_BOOL {
				if e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE {
					*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
				}
			} else {
				*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			}
		case VARIABLE_TYPE_NULL:
			if t2.IsPointer() {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and '%s'(non-pointer)",
					errMsgPrefix(e.Pos),
					e.OpName(),
					t2.TypeString()))
			}
			if e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and 'pointer' ",
					errMsgPrefix(e.Pos),
					e.OpName()))
			}
		case VARIABLE_TYPE_ARRAY_INSTANCE:
			fallthrough
		case VARIABLE_TYPE_OBJECT:
			if t2.IsPointer() == false && t2.Typ != VARIABLE_TYPE_NULL {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'pointer' and '%s'(non-pointer)",
					errMsgPrefix(e.Pos),
					e.OpName(),
					t2.TypeString()))
			}
			if e.Typ != EXPRESSION_TYPE_EQ && e.Typ != EXPRESSION_TYPE_NE {
				*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on 'null' and 'pointer' ", errMsgPrefix(e.Pos), e.OpName()))
			}
		default:
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'", errMsgPrefix(e.Pos),
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
		if t1.Typ == VARIABLE_TYPE_STRING || t2.Typ == VARIABLE_TYPE_STRING { // string is always ok
			result = &VariableType{}
			result.Typ = VARIABLE_TYPE_STRING
			result.Pos = e.Pos
			return result
		}
		if t1.IsNumber() == false || t2.IsNumber() == false {
			*errs = append(*errs, e.wrongOpErr(t1.TypeString(), t2.TypeString()))
			result = t1.Clone()
			result.Pos = e.Pos
			return result
		}
		result = &VariableType{}
		result.Pos = e.Pos
		result.Typ = t1.NumberTypeConvertRule(t2)
		return result
	}
	panic("missing check" + e.OpName())
	return nil
}
