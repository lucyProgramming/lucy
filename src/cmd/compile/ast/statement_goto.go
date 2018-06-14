package ast

import "fmt"

type StatementGoTo struct {
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
	return s.StatementGoTo.StatementLabel.Ready(s.Pos)
}
