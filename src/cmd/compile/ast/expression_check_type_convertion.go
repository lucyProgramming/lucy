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
		if convertion.Expression.IsLiteral() {
			convertion.Expression.convertNumberLiteralTo(convertion.Typ.Typ)
			//rewrite
			pos := e.Pos
			*e = *convertion.Expression
			e.Pos = pos // keep pos
		}
		return tt
	}
	ret := convertion.Typ.Clone()
	ret.Pos = e.Pos
	// string(['h'] , 'e')
	if convertion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_ARRAY && t.ArrayType.Typ == VARIABLE_TYPE_BYTE {
		return ret
	}
	// string (new String(""))
	if convertion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_OBJECT {
		return ret
	}
	// []byte("hello world")
	if convertion.Typ.Typ == VARIABLE_TYPE_ARRAY && convertion.Typ.ArrayType.Typ == VARIABLE_TYPE_BYTE &&
		t.Typ == VARIABLE_TYPE_STRING {
		return ret
	}
	/*
		xxx(yyy)
	*/
	if convertion.Typ.Typ == VARIABLE_TYPE_OBJECT && t.IsPointer() {
		return ret
	}
	*errs = append(*errs, fmt.Errorf("%s cannot convert '%s' to '%s'",
		errMsgPrefix(e.Pos), t.TypeString(), convertion.Typ.TypeString()))
	return ret
}
