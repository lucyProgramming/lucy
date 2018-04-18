package ast

import (
	"fmt"
)

func (e *Expression) checkIndexExpression(block *Block, errs *[]error) (t *VariableType) {
	index := e.Data.(*ExpressionIndex)
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
	if t.Typ != VARIABLE_TYPE_ARRAY && t.Typ != VARIABLE_TYPE_MAP {
		*errs = append(*errs, fmt.Errorf("%s cannot have 'index' on '%s'",
			errMsgPrefix(e.Pos), t.TypeString()))
		return nil
	}
	// array
	if t.Typ == VARIABLE_TYPE_ARRAY {
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
		tt := t.ArrayType.Clone()
		tt.Pos = e.Pos
		return tt
	}
	// map
	indexTs, es := index.Index.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	indexT, err := index.Index.mustBeOneValueContext(indexTs)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t != nil {
		if t.Map.K.Equal(indexT) == false {
			*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' for index",
				errMsgPrefix(e.Pos), indexT.TypeString(), t.Map.K.TypeString()))
		}
	}
	tt := t.Map.V.Clone()
	tt.Pos = e.Pos
	return tt

}
