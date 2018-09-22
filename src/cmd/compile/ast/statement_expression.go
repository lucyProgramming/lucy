package ast

import "fmt"

func (s *Statement) checkStatementExpression(block *Block) []error {
	errs := []error{}
	s.Expression.IsStatementExpression = true
	if s.Expression.canBeUsedAsStatement() == false {
		err := fmt.Errorf("%s expression '%s' evaluate but not used",
			errMsgPrefix(s.Expression.Pos), s.Expression.Description)
		errs = append(errs, err)
	}
	_, es := s.Expression.check(block)
	errs = append(errs, es...)
	return errs
}
