package ast

import (
	"fmt"
)

type StatementGoTo struct {
	Defers         []*StatementDefer
	LabelName      string
	StatementLabel *StatementLabel
}

func (s *Statement) checkStatementGoTo(b *Block) error {
	label := b.searchLabel(s.StatementGoTo.LabelName)
	if label == nil {
		return fmt.Errorf("%s label named '%s' not found",
			errMsgPrefix(s.Pos), s.StatementGoTo.LabelName)
	}
	s.StatementGoTo.StatementLabel = label
	s.StatementGoTo.mkDefers(b)
	return s.StatementGoTo.StatementLabel.Ready(s.Pos)
}

func (s *StatementGoTo) mkDefers(currentBlock *Block) {
	bs := []*Block{}
	for s.StatementLabel.Block != currentBlock {
		bs = append(bs, currentBlock)
		currentBlock = currentBlock.Outer
	}
	for _, b := range bs {
		if b.Defers != nil {
			s.Defers = append(s.Defers, b.Defers...)
		}
	}
}
