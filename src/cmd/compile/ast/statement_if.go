package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementIF struct {
	PreExpressions []*Expression
	Condition      *Expression
	CondtionBlock  Block
	Block          Block
	ElseIfList     []*StatementElseIf
	ElseBlock      *Block
	BackPatchs     []*cg.JumpBackPatch
}

func (s *StatementIF) check(father *Block) []error {
	s.CondtionBlock.inherite(father)
	errs := []error{}
	for _, v := range s.PreExpressions {
		v.IsStatementExpression = true
		_, es := v.check(&s.CondtionBlock)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		if v.canBeUsedAsStatement() == false {
			err := fmt.Errorf("%s expression '%s' evaluate but not used",
				errMsgPrefix(v.Pos), v.OpName())
			errs = append(errs, err)
			continue
		}
	}
	if s.Condition != nil {
		conditionType, es := s.Condition.checkSingleValueContextExpression(&s.CondtionBlock)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		if s.Condition.canbeUsedAsCondition() == false {
			errs = append(errs, fmt.Errorf("%s expression '%s' cannot used as condition",
				errMsgPrefix(s.Condition.Pos), s.Condition.OpName()))
		}
		if conditionType != nil && conditionType.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
				errMsgPrefix(s.Condition.Pos)))
		}
	}

	s.Block.inherite(&s.CondtionBlock)
	errs = append(errs, s.Block.checkStatements()...)
	for _, v := range s.ElseIfList {
		v.Block.inherite(&s.CondtionBlock)
		if v.Condition.canbeUsedAsCondition() == false {
			errs = append(errs, fmt.Errorf("%s expression '%s' cannot used as condition",
				errMsgPrefix(s.Condition.Pos), v.Condition.OpName()))
		}
		conditionType, es := v.Condition.checkSingleValueContextExpression(v.Block)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		if conditionType != nil && conditionType.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
				errMsgPrefix(s.Condition.Pos)))
		}
		errs = append(errs, v.Block.checkStatements()...)
	}
	if s.ElseBlock != nil {
		s.ElseBlock.inherite(&s.CondtionBlock)
		errs = append(errs, s.ElseBlock.checkStatements()...)
	}
	return errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}
