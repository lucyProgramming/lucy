package ast

type StatementBreak struct {
	Defers              []*StatementDefer
	StatementFor        *StatementFor
	StatementSwitch     *StatementSwitch
	SwitchTemplateBlock *Block
}

func (s *StatementBreak) mkDefers(block *Block) {
	if s.StatementFor != nil {
		if block.IsForBlock {
			s.Defers = append(s.Defers, block.Defers...)
			return
		}
		s.mkDefers(block.Outer)
		return
	} else if s.StatementSwitch != nil {
		//switch
		if block.IsSwitchBlock {
			s.Defers = append(s.Defers, block.Defers...)
			return
		}
		s.mkDefers(block.Outer)
	} else { //  s.SwitchTemplateBlock != nil
		if block.IsSwitchTemplateBlock {
			s.Defers = append(s.Defers, block.Defers...)
			return
		}
		s.mkDefers(block.Outer)
	}
}
