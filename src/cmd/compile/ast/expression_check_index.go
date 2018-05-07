package ast

import (
	"fmt"
)

func (e *Expression) checkIndexExpression(block *Block, errs *[]error) *VariableType {
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
	if t.Typ != VARIABLE_TYPE_ARRAY &&
		t.Typ != VARIABLE_TYPE_MAP &&
		t.Typ != VARIABLE_TYPE_JAVA_ARRAY {
		*errs = append(*errs, fmt.Errorf("%s cannot have 'index' on '%s'",
			errMsgPrefix(e.Pos), t.TypeString()))
		return nil
	}
	// array
	if t.Typ == VARIABLE_TYPE_ARRAY ||
		t.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		ts, es := index.Index.check(block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		indexType, err := index.Index.mustBeOneValueContext(ts)
		if err != nil {
			*errs = append(*errs, err)
		}
		if indexType != nil {
			if indexType.IsInteger() {
				if indexType.Typ == VARIABLE_TYPE_LONG {
					index.Index.ConvertToNumber(VARIABLE_TYPE_INT) //  convert to int
				}
			} else {
				*errs = append(*errs, fmt.Errorf("%s only integer can be used as index,but '%s'",
					errMsgPrefix(e.Pos), indexType.TypeString()))
			}
		}
		tt := t.ArrayType.Clone()
		tt.Pos = e.Pos
		return tt
	}
	// map
	ret := t.Map.V.Clone()
	ret.Pos = e.Pos
	ts, es = index.Index.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	indexType, err := index.Index.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if indexType == nil {
		return ret
	}
	if t.Map.K.Equal(indexType) == false {
		*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s' for index",
			errMsgPrefix(e.Pos), indexType.TypeString(), t.Map.K.TypeString()))
	}
	return ret

}
