package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementIf struct {
	PrefixExpressions []*Expression
	Condition         *Expression
	ConditionBlock    Block
	Block             Block
	ElseIfList        []*StatementElseIf
	ElseBlock         *Block
	Exits             []*cg.Exit
}

func (s *StatementIf) check(father *Block) []error {
	s.ConditionBlock.inherit(father)
	errs := []error{}
	for _, v := range s.PrefixExpressions {
		v.IsStatementExpression = true
		_, es := v.check(&s.ConditionBlock)
		if esNotEmpty(es) {
			errs = append(errs, es...)
		}
		if v.canBeUsedAsStatement() == false {
			err := fmt.Errorf("%s expression '%s' evaluate but not used",
				errMsgPrefix(v.Pos), v.Description)
			errs = append(errs, err)
		}
	}
	if s.Condition != nil {
		conditionType, es := s.Condition.checkSingleValueContextExpression(&s.ConditionBlock)
		if esNotEmpty(es) {
			errs = append(errs, es...)
		}
		if conditionType != nil && conditionType.Type != VariableTypeBool {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
				errMsgPrefix(s.Condition.Pos)))
		}
		if s.Condition.canBeUsedAsCondition() == false {
			errs = append(errs, fmt.Errorf("%s expression '%s' cannot used as condition",
				errMsgPrefix(s.Condition.Pos), s.Condition.Description))
		}
	}
	s.Block.inherit(&s.ConditionBlock)
	errs = append(errs, s.Block.checkStatements()...)
	for _, v := range s.ElseIfList {
		v.Block.inherit(&s.ConditionBlock)
		if v.Condition.canBeUsedAsCondition() == false {
			errs = append(errs, fmt.Errorf("%s expression '%s' cannot used as condition",
				errMsgPrefix(s.Condition.Pos), v.Condition.Description))
		}
		conditionType, es := v.Condition.checkSingleValueContextExpression(v.Block)
		if esNotEmpty(es) {
			errs = append(errs, es...)
		}
		if conditionType != nil && conditionType.Type != VariableTypeBool {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
				errMsgPrefix(s.Condition.Pos)))
		}
		errs = append(errs, v.Block.checkStatements()...)
	}
	if s.ElseBlock != nil {
		s.ElseBlock.inherit(&s.ConditionBlock)
		errs = append(errs, s.ElseBlock.checkStatements()...)
	}
	return errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}
