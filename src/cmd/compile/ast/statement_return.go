package ast

import (
	"fmt"
)

type StatementReturn struct {
	Function    *Function
	Pos         *Pos // use some time
	Expressions []*Expression
}

func (s *StatementReturn) check(b *Block) []error {
	s.Function = b.InheritedAttribute.function
	if len(b.InheritedAttribute.function.Typ.ReturnList) > 0 && len(s.Expressions) == 0 {
		s.Expressions = make([]*Expression, len(b.InheritedAttribute.function.Typ.ReturnList))
		for k, v := range b.InheritedAttribute.function.Typ.ReturnList {
			identifer := &ExpressionIdentifer{
				Name: v.Name,
			}
			s.Expressions[k] = &Expression{
				Data: identifer,
				Typ:  EXPRESSION_TYPE_IDENTIFIER,
			}
		}
	}
	if len(s.Expressions) == 0 {
		return nil
	}
	errs := make([]error, 0)
	returndValueTypes := checkExpressions(b, s.Expressions, &errs)
	pos := s.Expressions[len(s.Expressions)-1].Pos
	rs := b.InheritedAttribute.function.Typ.ReturnList
	if len(returndValueTypes) < len(rs) {
		errs = append(errs, fmt.Errorf("%s too few arguments to return", errMsgPrefix(pos)))
	}
	if len(returndValueTypes) > len(rs) {
		errs = append(errs, fmt.Errorf("%s too many arguments to return", errMsgPrefix(pos)))
	}
	for k, v := range rs {
		if k < len(returndValueTypes) {
			if !v.Typ.TypeCompatible(returndValueTypes[k]) {
				errs = append(errs, fmt.Errorf("%s cannot use %s as %s to return",
					errMsgPrefix(returndValueTypes[k].Pos),
					returndValueTypes[k].TypeString(),
					v.Typ.TypeString()))
			}
		}
	}
	return errs
}
