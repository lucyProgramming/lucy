package ast

import (
	"fmt"
)

func (e *Expression) checkIndexExpression(block *Block, errs *[]error) (t *VariableType) {
	index := e.Data.(*ExpressionIndex)
	f := func() *VariableType {
		ts, es := index.Expression.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		t, err := e.mustBeOneValueContext(ts)
		if err != nil {
			*errs = append(*errs, err)
		}
		if t == nil {
			return nil
		}
		if t.Typ != VARIABLE_TYPE_ARRAY_INSTANCE && VARIABLE_TYPE_OBJECT != t.Typ {
			op := "access"
			if e.Typ == EXPRESSION_TYPE_INDEX {
				op = "index"
			}
			*errs = append(*errs, fmt.Errorf("%s cannot %s on %s", errMsgPrefix(e.Pos), op, t.TypeString()))
			return nil
		}
		return t
	}
	obj := f()
	if obj == nil {
		return nil
	}
	if e.Typ == EXPRESSION_TYPE_INDEX { // index
		if obj.Typ != VARIABLE_TYPE_ARRAY_INSTANCE {
			*errs = append(*errs, fmt.Errorf("%s expresson is not a array", errMsgPrefix(e.Pos)))
			return nil
		}
		ts, es := index.Index.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		t, err := e.mustBeOneValueContext(ts)
		if err != nil {
			*errs = append(*errs, err)
		}
		if t != nil {
			if !t.IsInteger() {
				*errs = append(*errs, fmt.Errorf("%s only integer can be used as index,but '%s'",
					errMsgPrefix(e.Pos), t.TypeString()))
			}
		}
		tt := obj.CombinationType.Clone()
		tt.Pos = e.Pos
		return tt
	}
	// dot
	if obj.Typ != VARIABLE_TYPE_OBJECT && obj.Typ != VARIABLE_TYPE_CLASS {
		*errs = append(*errs, fmt.Errorf("%s cannot access field '%s' on '%s'", errMsgPrefix(e.Pos), index.Name, obj.TypeString()))
		return nil
	}
	if e.Typ != EXPRESSION_TYPE_DOT {
		*errs = append(*errs, fmt.Errorf("%s object`s field can only access by '.'",
			errMsgPrefix(e.Pos)))
		return nil
	}
	field, err := obj.Class.accessField(index.Name)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
	} else {
		if !index.Expression.isThisIdentifierExpression() && !field.isPublic() {
			*errs = append(*errs, fmt.Errorf("%s field %s is private", errMsgPrefix(e.Pos),
				index.Name))
		}
	}
	if field != nil {
		return field.Typ
	} else {
		return nil
	}
	return nil
}
