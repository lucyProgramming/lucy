package ast

import (
	"fmt"
)

type StatementReturn struct {
	Defers      []*Defer
	Expressions []*Expression
}

func (s *StatementReturn) mkDefers(b *Block) {
	if b.IsFunctionTopBlock == false { // not top block
		s.mkDefers(b.Outter) // recurvilly
	}
	if b.Defers != nil {
		s.Defers = append(s.Defers, b.Defers...)
	}
}

func (s *StatementReturn) check(b *Block) []error {
	s.mkDefers(b)
	if len(s.Expressions) == 0 {
		return nil
	}
	errs := make([]error, 0)
	returndValueTypes := checkRightValuesValid(checkExpressions(b, s.Expressions, &errs), &errs)
	pos := s.Expressions[len(s.Expressions)-1].Pos
	rs := b.InheritedAttribute.Function.Typ.ReturnList
	if len(returndValueTypes) < len(rs) {
		errs = append(errs, fmt.Errorf("%s too few arguments to return", errMsgPrefix(pos)))
	} else if len(returndValueTypes) > len(rs) {
		errs = append(errs, fmt.Errorf("%s too many arguments to return", errMsgPrefix(pos)))
	}
	convertLiteralExpressionsToNeeds(s.Expressions,
		b.InheritedAttribute.Function.Typ.retTypes(s.Expressions[0].Pos), returndValueTypes)
	for k, v := range rs {
		if k < len(returndValueTypes) && returndValueTypes[k] != nil {
			if !v.Typ.TypeCompatible(returndValueTypes[k]) {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s' to return",
					errMsgPrefix(returndValueTypes[k].Pos),
					returndValueTypes[k].TypeString(),
					v.Typ.TypeString()))
			}
		}
	}
	return errs
}
