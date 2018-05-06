package ast

import "fmt"

type StatementGoto struct {
	Name           string
	StatementLable *StatementLable
}

func (s *Statement) checkStatementGoto(b *Block) error {
	lable := b.searchLable(s.StatementGoto.Name)
	if lable == nil {
		return fmt.Errorf("%s label named '%s' not found",
			errMsgPrefix(s.Pos), s.StatementGoto.Name)
	}
	s.StatementGoto.StatementLable = lable
	return s.StatementGoto.StatementLable.Ready(s.Pos)
}
