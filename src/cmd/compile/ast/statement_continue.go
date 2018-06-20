package ast

type StatementContinue struct {
	StatementFor *StatementFor
	Defers       []*StatementDefer
}

func (s *StatementContinue) mkDefers(block *Block) {
	if block.IsForBlock {
		s.Defers = append(s.Defers, block.Defers...)
		return
	}
	s.mkDefers(block.Outer)
}
