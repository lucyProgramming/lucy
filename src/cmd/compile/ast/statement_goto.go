package ast

import "fmt"

type StatementGoto struct {
	Name           string
	Pos            *Pos
	StatementLable *StatementLable
}

func (s *Statement) checkStatementGoto(b *Block) error {
	t := b.SearchByName(s.StatementGoto.Name)
	if t == nil {
		return fmt.Errorf("%s label named '%s' not found",
			errMsgPrefix(s.StatementGoto.Pos), s.StatementGoto.Name)
	}
	if l, ok := t.(*StatementLable); ok == false {
		return fmt.Errorf("%s '%s' is not a lable",
			errMsgPrefix(s.StatementGoto.Pos), s.StatementGoto.Name)
	} else {
		s.StatementGoto.StatementLable = l
	}
	return s.StatementGoto.StatementLable.Ready(s.Pos)
}
