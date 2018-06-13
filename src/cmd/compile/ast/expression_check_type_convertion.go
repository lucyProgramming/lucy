package ast

import (
	"fmt"
)

func (e *Expression) checkTypeConvertionExpression(block *Block, errs *[]error) *VariableType {
	conversion := e.Data.(*ExpressionTypeConversion)
	t, es := conversion.Expression.checkSingleValueContextExpression(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if t == nil {
		return nil
	}
	err := conversion.Typ.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	ret := conversion.Typ.Clone()
	ret.Pos = e.Pos

	if t.IsNumber() && conversion.Typ.IsNumber() {
		if conversion.Expression.IsLiteral() {
			conversion.Expression.convertNumberLiteralTo(conversion.Typ.Typ)
			//rewrite
			pos := e.Pos
			*e = *conversion.Expression
			e.Pos = pos // keep pos
		}
		return ret
	}

	// string([]byte)
	if conversion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_ARRAY && t.ArrayType.Typ == VARIABLE_TYPE_BYTE {
		return ret
	}
	// string(byte[])
	if conversion.Typ.Typ == VARIABLE_TYPE_STRING &&
		t.Typ == VARIABLE_TYPE_JAVA_ARRAY && t.ArrayType.Typ == VARIABLE_TYPE_BYTE {
		return ret
	}

	// []byte("hello world")
	if conversion.Typ.Typ == VARIABLE_TYPE_ARRAY && conversion.Typ.ArrayType.Typ == VARIABLE_TYPE_BYTE &&
		t.Typ == VARIABLE_TYPE_STRING {
		return ret
	}
	// byte[]("hello world")
	if conversion.Typ.Typ == VARIABLE_TYPE_JAVA_ARRAY && conversion.Typ.ArrayType.Typ == VARIABLE_TYPE_BYTE &&
		t.Typ == VARIABLE_TYPE_STRING {
		return ret
	}
	if conversion.Typ.validForTypeAssert() && t.IsPointer() {
		return ret
	}

	*errs = append(*errs, fmt.Errorf("%s cannot convert '%s' to '%s'",
		errMsgPrefix(e.Pos), t.TypeString(), conversion.Typ.TypeString()))
	return ret
}
