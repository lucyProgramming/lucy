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
	err := conversion.Type.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	ret := conversion.Type.Clone()
	ret.Pos = e.Pos

	if t.IsNumber() && conversion.Type.IsNumber() {
		if conversion.Expression.IsLiteral() {
			conversion.Expression.convertNumberLiteralTo(conversion.Type.Type)
			//rewrite
			pos := e.Pos
			*e = *conversion.Expression
			e.Pos = pos // keep pos
		}
		return ret
	}

	// string([]byte)
	if conversion.Type.Type == VARIABLE_TYPE_STRING &&
		t.Type == VARIABLE_TYPE_ARRAY && t.ArrayType.Type == VARIABLE_TYPE_BYTE {
		return ret
	}
	// string(byte[])
	if conversion.Type.Type == VARIABLE_TYPE_STRING &&
		t.Type == VARIABLE_TYPE_JAVA_ARRAY && t.ArrayType.Type == VARIABLE_TYPE_BYTE {
		return ret
	}

	// []byte("hello world")
	if conversion.Type.Type == VARIABLE_TYPE_ARRAY && conversion.Type.ArrayType.Type == VARIABLE_TYPE_BYTE &&
		t.Type == VARIABLE_TYPE_STRING {
		return ret
	}
	// byte[]("hello world")
	if conversion.Type.Type == VARIABLE_TYPE_JAVA_ARRAY && conversion.Type.ArrayType.Type == VARIABLE_TYPE_BYTE &&
		t.Type == VARIABLE_TYPE_STRING {
		return ret
	}
	if conversion.Type.validForTypeAssertOrConversion() && t.IsPointer() {
		return ret
	}

	*errs = append(*errs, fmt.Errorf("%s cannot convert '%s' to '%s'",
		errMsgPrefix(e.Pos), t.TypeString(), conversion.Type.TypeString()))
	return ret
}
