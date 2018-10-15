package ast

import "fmt"

type StatementContinue struct {
	StatementFor *StatementFor
	Defers       []*StatementDefer
}

func (c *StatementContinue) check(s *Statement, block *Block) []error {
	if block.InheritedAttribute.ForContinue == nil {
		return []error{fmt.Errorf("%s 'continue' can`t in this scope",
			errMsgPrefix(s.Pos))}
	}
	if block.InheritedAttribute.Defer != nil {
		return []error{fmt.Errorf("%s cannot has 'continue' in 'defer'",
			errMsgPrefix(s.Pos))}
	}
	s.StatementContinue.StatementFor = block.InheritedAttribute.ForContinue
	s.StatementContinue.mkDefers(block)
	return nil
}

func (c *StatementContinue) mkDefers(block *Block) {
	if block.IsForBlock {
		c.Defers = append(c.Defers, block.Defers...)
		return
	}
	c.mkDefers(block.Outer)
}
