package ast

import (
	"fmt"
)

type StatementGoTo struct {
	Defers         []*StatementDefer
	LabelName      string
	StatementLabel *StatementLabel
}

func (g *StatementGoTo) checkStatementGoTo(pos *Pos, b *Block) error {
	label := b.searchLabel(g.LabelName)
	if label == nil {
		return fmt.Errorf("%s label named '%s' not found",
			pos.ErrMsgPrefix(), g.LabelName)
	}
	g.StatementLabel = label
	g.mkDefers(b)
	return g.StatementLabel.Ready(pos)
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
