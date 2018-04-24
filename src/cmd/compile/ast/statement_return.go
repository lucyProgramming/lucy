package ast

import (
	"fmt"
)

type StatementReturn struct {
	//Function    *Function
	Pos         *Pos // use some time
	Expressions []*Expression
}

func (s *StatementReturn) check(b *Block) []error {
	//s.Function = b.InheritedAttribute.Function
	if len(s.Expressions) == 0 {
		return nil
	}
	errs := make([]error, 0)
	returndValueTypes := checkExpressions(b, s.Expressions, &errs)
	pos := s.Expressions[len(s.Expressions)-1].Pos
	rs := b.InheritedAttribute.Function.Typ.ReturnList
	if len(returndValueTypes) < len(rs) {
		errs = append(errs, fmt.Errorf("%s too few arguments to return", errMsgPrefix(pos)))
	}
	if len(returndValueTypes) > len(rs) {
		errs = append(errs, fmt.Errorf("%s too many arguments to return", errMsgPrefix(pos)))
	}
	for k, v := range rs {
		if k < len(returndValueTypes) && returndValueTypes[k] != nil {
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
