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

func (gt *StatementGoTo) checkStatementGoTo(b *Block) error {
	label := b.searchLabel(gt.LabelName)
	if label == nil {
		return fmt.Errorf("%s label named '%s' not found",
			gt.Pos.ErrMsgPrefix(), gt.LabelName)
	}
	gt.StatementLabel = label
	gt.mkDefers(b)
	return gt.StatementLabel.Ready(gt.Pos)
}

func (gt *StatementGoTo) mkDefers(currentBlock *Block) {
	bs := []*Block{}
	for gt.StatementLabel.Block != currentBlock {
		bs = append(bs, currentBlock)
		currentBlock = currentBlock.Outer
	}
	for _, b := range bs {
		if b.Defers != nil {
			gt.Defers = append(gt.Defers, b.Defers...)
		}
	}
}
