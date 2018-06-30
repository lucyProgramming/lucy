package ast

import (
	"fmt"
)

func (e *Expression) checkTypeConversionExpression(block *Block, errs *[]error) *Type {
	conversion := e.Data.(*ExpressionTypeConversion)
	t, es := conversion.Expression.checkSingleValueContextExpression(block)
	if errorsNotEmpty(es) {
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
	if conversion.Type.Type == VariableTypeString &&
		t.Type == VariableTypeArray && t.Array.Type == VariableTypeByte {
		return ret
	}
	// string(byte[])
	if conversion.Type.Type == VariableTypeString &&
		t.Type == VariableTypeJavaArray && t.Array.Type == VariableTypeByte {
		return ret
	}

	// []byte("hello world")
	if conversion.Type.Type == VariableTypeArray && conversion.Type.Array.Type == VariableTypeByte &&
		t.Type == VariableTypeString {
		return ret
	}
	// byte[]("hello world")
	if conversion.Type.Type == VariableTypeJavaArray && conversion.Type.Array.Type == VariableTypeByte &&
		t.Type == VariableTypeString {
		return ret
	}
	if conversion.Type.validForTypeAssertOrConversion() && t.IsPointer() {
		return ret
	}
	*errs = append(*errs, fmt.Errorf("%s cannot convert '%s' to '%s'",
		errMsgPrefix(e.Pos), t.TypeString(), conversion.Type.TypeString()))
	return ret
}
