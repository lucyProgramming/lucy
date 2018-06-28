package ast

type StatementBreak struct {
	Defers          []*StatementDefer
	StatementFor    *StatementFor
	StatementSwitch *StatementSwitch
}

func (s *StatementBreak) mkDefers(block *Block) {
	if s.StatementFor != nil {
		if block.IsForBlock {
			s.Defers = append(s.Defers, block.Defers...)
			return
		}
		s.mkDefers(block.Outer)
		return
	} else {
		// switch
		if block.IsSwitchStatementTopBlock {
			s.Defers = append(s.Defers, block.Defers...)
			return
		}
		s.mkDefers(block.Outer)
	}
}
