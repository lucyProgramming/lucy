package ast

import (
	"fmt"
)

type StatementReturn struct {
	Defers      []*StatementDefer
	Expressions []*Expression
	Pos         *Pos
}

func (this *StatementReturn) mkDefers(b *Block) {
	if b.IsFunctionBlock == false { // not top block
		this.mkDefers(b.Outer) // recursive
	}
	if b.Defers != nil {
		this.Defers = append(this.Defers, b.Defers...)
	}
}

func (this *StatementReturn) check(b *Block) []error {
	if b.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'return' in 'defer'",
			this.Pos.ErrMsgPrefix())}
	}
	errs := []error{}
	this.mkDefers(b)
	if len(this.Expressions) == 0 { // always ok
		return errs
	}
	returnValueTypes := checkExpressions(b, this.Expressions, &errs, false)
	rs := b.InheritedAttribute.Function.Type.ReturnList
	pos := this.Expressions[len(this.Expressions)-1].Pos
	if len(returnValueTypes) < len(rs) {
		errs = append(errs, fmt.Errorf("%s too few arguments to return", pos.ErrMsgPrefix()))
	} else if len(returnValueTypes) > len(rs) {
		errs = append(errs, fmt.Errorf("%s too many arguments to return", pos.ErrMsgPrefix()))
	}
	convertExpressionsToNeeds(this.Expressions,
		b.InheritedAttribute.Function.Type.mkCallReturnTypes(this.Expressions[0].Pos), returnValueTypes)
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
