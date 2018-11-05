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

func (this *StatementGoTo) checkStatementGoTo(b *Block) error {
	label := b.searchLabel(this.LabelName)
	if label == nil {
		return fmt.Errorf("%s label named '%s' not found",
			this.Pos.ErrMsgPrefix(), this.LabelName)
	}
	this.StatementLabel = label
	this.mkDefers(b)
	return this.StatementLabel.Ready(this.Pos)
}

func (this *StatementGoTo) mkDefers(currentBlock *Block) {
	bs := []*Block{}
	for this.StatementLabel.Block != currentBlock {
		bs = append(bs, currentBlock)
		currentBlock = currentBlock.Outer
	}
	for _, b := range bs {
		if b.Defers != nil {
			this.Defers = append(this.Defers, b.Defers...)
		}
	}
}
