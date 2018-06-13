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
		s.mkDefers(b.Outer) // recurvilly
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
	returnValueTypes := checkRightValuesValid(checkExpressions(b, s.Expressions, &errs), &errs)
	pos := s.Expressions[len(s.Expressions)-1].Pos
	rs := b.InheritedAttribute.Function.Type.ReturnList
	if len(returnValueTypes) < len(rs) {
		errs = append(errs, fmt.Errorf("%s too few arguments to return", errMsgPrefix(pos)))
	} else if len(returnValueTypes) > len(rs) {
		errs = append(errs, fmt.Errorf("%s too many arguments to return", errMsgPrefix(pos)))
	}
	convertLiteralExpressionsToNeeds(s.Expressions,
		b.InheritedAttribute.Function.Type.retTypes(s.Expressions[0].Pos), returnValueTypes)
	for k, v := range rs {
		if k < len(returnValueTypes) && returnValueTypes[k] != nil {
			if false == v.Type.Equal(&errs, returnValueTypes[k]) {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s' to return",
					errMsgPrefix(returnValueTypes[k].Pos),
					returnValueTypes[k].TypeString(),
					v.Type.TypeString()))
			}
		}
	}
	return errs
}
