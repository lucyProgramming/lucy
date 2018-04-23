package ast

import (
	"fmt"
)

func (e *Expression) checkTypeConvertionExpression(block *Block, errs *[]error) *VariableType {
	convertion := e.Data.(*ExpressionTypeConvertion)
	ts, es := convertion.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := convertion.Expression.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if t == nil {
		return nil
	}
	err = convertion.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
		return nil
	}
	if t.IsNumber() && convertion.Typ.IsNumber() {
		tt := convertion.Typ.Clone()
		tt.Pos = e.Pos
		return tt
	}
	if convertion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_ARRAY && t.ArrayType.Typ == VARIABLE_TYPE_BYTE {
		tt := convertion.Typ.Clone()
		tt.Pos = e.Pos
		return tt
	}
	if convertion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_OBJECT {
		tt := convertion.Typ.Clone()
		tt.Pos = e.Pos
		return tt
	}
	if convertion.Typ.Typ == VARIABLE_TYPE_ARRAY && convertion.Typ.ArrayType.Typ == VARIABLE_TYPE_BYTE &&
		t.Typ == VARIABLE_TYPE_STRING {
		tt := convertion.Typ.Clone()
		tt.Pos = e.Pos
		return tt
	}
	*errs = append(*errs, fmt.Errorf("%s cannot convert '%s' to '%s'",
		errMsgPrefix(e.Pos), t.TypeString(), convertion.Typ.TypeString()))
	return nil
}
