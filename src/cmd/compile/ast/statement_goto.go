package ast

import "fmt"

type StatementGoto struct {
	Name           string
	StatementLable *StatementLabel
}

func (s *Statement) checkStatementGoto(b *Block) error {
	label := b.searchLabel(s.StatementGoto.Name)
	if label == nil {
		return fmt.Errorf("%s label named '%s' not found",
			errMsgPrefix(s.Pos), s.StatementGoto.Name)
	}
	s.StatementGoto.StatementLable = label
	return s.StatementGoto.StatementLable.Ready(s.Pos)
}
