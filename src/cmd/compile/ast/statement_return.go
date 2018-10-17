package ast

import (
	"fmt"
)

type StatementReturn struct {
	Defers      []*StatementDefer
	Expressions []*Expression
}

func (r *StatementReturn) mkDefers(b *Block) {
	if b.IsFunctionBlock == false { // not top block
		r.mkDefers(b.Outer) // recursive
	}
	if b.Defers != nil {
		r.Defers = append(r.Defers, b.Defers...)
	}
}

func (r *StatementReturn) check(s *Statement, b *Block) []error {
	if b.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'return' in 'defer'",
			errMsgPrefix(s.Pos))}
	}
	errs := []error{}
	r.mkDefers(b)
	if len(r.Expressions) == 0 { // always ok
		return errs
	}
	returnValueTypes := checkExpressions(b, r.Expressions, &errs, false)
	rs := b.InheritedAttribute.Function.Type.ReturnList
	pos := r.Expressions[len(r.Expressions)-1].Pos
	if len(returnValueTypes) < len(rs) {
		errs = append(errs, fmt.Errorf("%s too few arguments to return", pos.ErrMsgPrefix()))
	} else if len(returnValueTypes) > len(rs) {
		errs = append(errs, fmt.Errorf("%s too many arguments to return", pos.ErrMsgPrefix()))
	}
	convertExpressionsToNeeds(r.Expressions,
		b.InheritedAttribute.Function.Type.mkCallReturnTypes(r.Expressions[0].Pos), returnValueTypes)
	for k, v := range rs {
		if k < len(returnValueTypes) && returnValueTypes[k] != nil {
			if false == v.Type.assignAble(&errs, returnValueTypes[k]) {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s' to return",
					returnValueTypes[k].Pos.ErrMsgPrefix(),
					returnValueTypes[k].TypeString(),
					v.Type.TypeString()))
			}
		}
	}
	return errs
}
