package ast

func (this *Statement) checkStatementExpression(block *Block) []error {
	var errs []error
	this.Expression.IsStatementExpression = true
	if err := this.Expression.canBeUsedAsStatement(); err != nil {
		errs = append(errs, err)
	}
	_, es := this.Expression.check(block)
	errs = append(errs, es...)
	return errs
}
