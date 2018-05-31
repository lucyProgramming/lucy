package ast

import (
	"fmt"
)

func (e *Expression) checkTypeConvertionExpression(block *Block, errs *[]error) *VariableType {
	convertion := e.Data.(*ExpressionTypeConvertion)
	t, es := convertion.Expression.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	err := convertion.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	ret := convertion.Typ.Clone()
	ret.Pos = e.Pos

	if t.IsNumber() && convertion.Typ.IsNumber() {
		if convertion.Expression.IsLiteral() {
			convertion.Expression.convertNumberLiteralTo(convertion.Typ.Typ)
			//rewrite
			pos := e.Pos
			*e = *convertion.Expression
			e.Pos = pos // keep pos
		}
		return ret
	}

	// string([]byte)
	if convertion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_ARRAY && t.ArrayType.Typ == VARIABLE_TYPE_BYTE {
		return ret
	}
	// string(byte[])
	if convertion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_JAVA_ARRAY && t.ArrayType.Typ == VARIABLE_TYPE_BYTE {
		return ret
	}

	// []byte("hello world")
	if convertion.Typ.Typ == VARIABLE_TYPE_ARRAY && convertion.Typ.ArrayType.Typ == VARIABLE_TYPE_BYTE &&
		t.Typ == VARIABLE_TYPE_STRING {
		return ret
	}
	// byte[]("hello world")
	if convertion.Typ.Typ == VARIABLE_TYPE_JAVA_ARRAY && convertion.Typ.ArrayType.Typ == VARIABLE_TYPE_BYTE &&
		t.Typ == VARIABLE_TYPE_STRING {
		return ret
	}
	if convertion.Typ.validForTypeAssert() && t.IsPointer() {
		return ret
	}

	*errs = append(*errs, fmt.Errorf("%s cannot convert '%s' to '%s'",
		errMsgPrefix(e.Pos), t.TypeString(), convertion.Typ.TypeString()))
	return ret
}
