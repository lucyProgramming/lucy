package ast

import "fmt"

type StatementContinue struct {
	StatementFor *StatementFor
	Defers       []*StatementDefer
	Pos          *Pos
}

func (this *StatementContinue) check(block *Block) []error {
	if block.InheritedAttribute.ForContinue == nil {
		return []error{fmt.Errorf("%s 'continue' can`t in this scope",
			this.Pos.ErrMsgPrefix())}
	}
	if block.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'continue' in 'defer'",
			this.Pos.ErrMsgPrefix())}
	}
	this.StatementFor = block.InheritedAttribute.ForContinue
	this.mkDefers(block)
	return nil
}

func (this *StatementContinue) mkDefers(block *Block) {
	if block.IsForBlock {
		this.Defers = append(this.Defers, block.Defers...)
		return
	}
	this.mkDefers(block.Outer)
}
