package ast

import (
	"fmt"
)

type StatementGoTo struct {
	Defers         []*StatementDefer
	LabelName      string
	StatementLabel *StatementLabel
	Block          *Block
}

func (s *Statement) checkStatementGoTo(b *Block) error {
	label := b.searchLabel(s.StatementGoTo.LabelName)
	if label == nil {
		return fmt.Errorf("%s label named '%s' not found",
			errMsgPrefix(s.Pos), s.StatementGoTo.LabelName)
	}
	s.StatementGoTo.StatementLabel = label
	s.StatementGoTo.Block = b
	s.StatementGoTo.mkDefers()
	return s.StatementGoTo.StatementLabel.Ready(s.Pos)
}

func (s *StatementGoTo) mkDefers() {
	bs := []*Block{}
	bb := s.Block
	for s.StatementLabel.Block != bb {
		bs = append(bs, bb)
		bb = bb.Outer
	}
	if len(bs) == 0 {
		return
	}
	for _, b := range bs {
		if b.Defers != nil {
			s.Defers = append(s.Defers, b.Defers...)
		}
	}
}
