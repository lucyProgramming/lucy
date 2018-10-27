package ast

import "fmt"

type StatementContinue struct {
	StatementFor *StatementFor
	Defers       []*StatementDefer
	Pos          *Pos
}

func (c *StatementContinue) check(block *Block) []error {
	if block.InheritedAttribute.ForContinue == nil {
		return []error{fmt.Errorf("%s 'continue' can`t in this scope",
			c.Pos.ErrMsgPrefix())}
	}
	if block.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'continue' in 'defer'",
			c.Pos.ErrMsgPrefix())}
	}
	c.StatementFor = block.InheritedAttribute.ForContinue
	c.mkDefers(block)
	return nil
}

func (c *StatementContinue) mkDefers(block *Block) {
	if block.IsForBlock {
		c.Defers = append(c.Defers, block.Defers...)
		return
	}
	c.mkDefers(block.Outer)
}
