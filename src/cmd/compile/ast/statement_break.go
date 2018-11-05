package ast

import "fmt"

type StatementBreak struct {
	Defers              []*StatementDefer
	StatementFor        *StatementFor
	StatementSwitch     *StatementSwitch
	SwitchTemplateBlock *Block
	Pos                 *Pos
}

func (this *StatementBreak) check(block *Block) []error {
	if block.InheritedAttribute.ForBreak == nil {
		return []error{fmt.Errorf("%s 'break' cannot in this scope", this.Pos.ErrMsgPrefix())}
	}
	if block.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'break' in 'defer'",
			this.Pos.ErrMsgPrefix())}
	}
	if t, ok := block.InheritedAttribute.ForBreak.(*StatementFor); ok {
		this.StatementFor = t
	} else if t, ok := block.InheritedAttribute.ForBreak.(*StatementSwitch); ok {
		this.StatementSwitch = t
	} else {
		this.SwitchTemplateBlock = block.InheritedAttribute.ForBreak.(*Block)
	}
	this.mkDefers(block)
	return nil
}

func (this *StatementBreak) mkDefers(block *Block) {
	if this.StatementFor != nil {
		if block.IsForBlock {
			this.Defers = append(this.Defers, block.Defers...)
			return
		}
		this.mkDefers(block.Outer)
		return
	} else if this.StatementSwitch != nil {
		//switch
		if block.IsSwitchBlock {
			this.Defers = append(this.Defers, block.Defers...)
			return
		}
		this.mkDefers(block.Outer)
	} else { //s.SwitchTemplateBlock != nil
		if block.IsWhenBlock {
			this.Defers = append(this.Defers, block.Defers...)
			return
		}
		this.mkDefers(block.Outer)
	}
}
