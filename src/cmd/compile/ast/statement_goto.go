package ast

import (
	"fmt"
)

type StatementGoTo struct {
	Defers         []*StatementDefer
	LabelName      string
	StatementLabel *StatementLabel
	Pos            *Pos
}

func (g *StatementGoTo) checkStatementGoTo(b *Block) error {
	label := b.searchLabel(g.LabelName)
	if label == nil {
		return fmt.Errorf("%s label named '%s' not found",
			g.Pos.ErrMsgPrefix(), g.LabelName)
	}
	g.StatementLabel = label
	g.mkDefers(b)
	return g.StatementLabel.Ready(g.Pos)
}

func (g *StatementGoTo) mkDefers(currentBlock *Block) {
	bs := []*Block{}
	for g.StatementLabel.Block != currentBlock {
		bs = append(bs, currentBlock)
		currentBlock = currentBlock.Outer
	}
	for _, b := range bs {
		if b.Defers != nil {
			g.Defers = append(g.Defers, b.Defers...)
		}
	}
}
