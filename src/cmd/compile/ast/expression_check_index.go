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
	if obj.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
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
				*errs = append(*errs, fmt.Errorf("%s only integer can be used as index",
					errMsgPrefix(e.Pos)))
			}
		}
		return obj.CombinationType
	}
	if obj.Typ == VARIABLE_TYPE_OBJECT {
		if e.Typ != EXPRESSION_TYPE_DOT {
			*errs = append(*errs, fmt.Errorf("%s object`s field can only access by '.'",
				errMsgPrefix(e.Pos)))
			return nil
		}
		f, err := obj.Class.accessField(index.Name)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
		} else {
			if !index.Expression.isThisIdentifierExpression() && !f.isPublic() {
				*errs = append(*errs, fmt.Errorf("%s field %s is private", errMsgPrefix(e.Pos),
					index.Name))
			}
		}
		if f != nil {
			return f.Typ
		} else {
			return nil
		}
	}
	panic("111")
	return nil
}
