package ast

import "fmt"

type StatementBreak struct {
	Defers              []*StatementDefer
	StatementFor        *StatementFor
	StatementSwitch     *StatementSwitch
	SwitchTemplateBlock *Block
}

func (b *StatementBreak) check(s *Statement, block *Block) []error {
	if block.InheritedAttribute.ForBreak == nil {
		return []error{fmt.Errorf("%s 'break' cannot in this scope", errMsgPrefix(s.Pos))}
	}
	if block.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'break' in 'defer'",
			errMsgPrefix(s.Pos))}
	}
	if t, ok := block.InheritedAttribute.ForBreak.(*StatementFor); ok {
		s.StatementBreak.StatementFor = t
	} else if t, ok := block.InheritedAttribute.ForBreak.(*StatementSwitch); ok {
		s.StatementBreak.StatementSwitch = t
	} else {
		s.StatementBreak.SwitchTemplateBlock = block.InheritedAttribute.ForBreak.(*Block)
	}
	s.StatementBreak.mkDefers(block)
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
		if block.IsSwitchTemplateBlock {
			b.Defers = append(b.Defers, block.Defers...)
			return
		}
		b.mkDefers(block.Outer)
	}
}
