package ast

func (s *Statement) checkStatementExpression(block *Block) []error {
	var errs []error
	s.Expression.IsStatementExpression = true
	if err := s.Expression.canBeUsedAsStatement(); err != nil {
		errs = append(errs, err)
	}
	_, es := s.Expression.check(block)
	errs = append(errs, es...)
	return errs
}
