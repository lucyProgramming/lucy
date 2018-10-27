package ast

import "fmt"

type StatementBreak struct {
	Defers              []*StatementDefer
	StatementFor        *StatementFor
	StatementSwitch     *StatementSwitch
	SwitchTemplateBlock *Block
	Pos                 *Pos
}

func (b *StatementBreak) check(block *Block) []error {
	if block.InheritedAttribute.ForBreak == nil {
		return []error{fmt.Errorf("%s 'break' cannot in this scope", b.Pos.ErrMsgPrefix())}
	}
	if block.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'break' in 'defer'",
			b.Pos.ErrMsgPrefix())}
	}
	if t, ok := block.InheritedAttribute.ForBreak.(*StatementFor); ok {
		b.StatementFor = t
	} else if t, ok := block.InheritedAttribute.ForBreak.(*StatementSwitch); ok {
		b.StatementSwitch = t
	} else {
		b.SwitchTemplateBlock = block.InheritedAttribute.ForBreak.(*Block)
	}
	b.mkDefers(block)
	return nil
}

func (b *StatementBreak) mkDefers(block *Block) {
	if b.StatementFor != nil {
		if block.IsForBlock {
			b.Defers = append(b.Defers, block.Defers...)
			return
		}
		b.mkDefers(block.Outer)
		return
	} else if b.StatementSwitch != nil {
		//switch
		if block.IsSwitchBlock {
			b.Defers = append(b.Defers, block.Defers...)
			return
		}
		b.mkDefers(block.Outer)
	} else { //s.SwitchTemplateBlock != nil
		if block.IsWhenBlock {
			b.Defers = append(b.Defers, block.Defers...)
			return
		}
		b.mkDefers(block.Outer)
	}
}
