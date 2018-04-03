package ast

import "fmt"

func (e *Expression) checkDotExpression(block *Block, errs *[]error) (t *VariableType) {
	index := e.Data.(*ExpressionDot)
	ts, es := index.Expression.check(block)
	if errsNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	t, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}

	if t == nil {
		//try package
		return nil
	}

	// dot
	if t.Typ != VARIABLE_TYPE_OBJECT && t.Typ != VARIABLE_TYPE_CLASS &&
		t.Typ != VARIABLE_TYPE_PACKAGE {
		*errs = append(*errs, fmt.Errorf("%s cannot access field '%s' on '%s'", errMsgPrefix(e.Pos), index.Name, t.TypeString()))
		return nil
	}
	field, err := t.Class.accessField(index.Name, false)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(e.Pos), err.Error()))
	} else {
		if !index.Expression.isThisIdentifierExpression() && !field.IsPublic() {
			*errs = append(*errs, fmt.Errorf("%s field %s is private", errMsgPrefix(e.Pos),
				index.Name))
		}
	}
	if field != nil {
		t := field.Typ.Clone()
		t.Pos = e.Pos
		index.Field = field
		return t
	} else {
		return nil
	}
}
