package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementIf struct {
	PrefixExpressions []*Expression
	Condition         *Expression
	ConditionBlock    Block
	TrueBlock         Block
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
		errs = append(errs, es...)
		if err := v.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
	}
	if s.Condition != nil {
		conditionType, es := s.Condition.checkSingleValueContextExpression(&s.ConditionBlock)
		errs = append(errs, es...)
		if conditionType != nil && conditionType.Type != VariableTypeBool {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
				errMsgPrefix(s.Condition.Pos)))
		}
		if err := s.Condition.canBeUsedAsCondition(); err != nil {
			errs = append(errs, err)
		}
	}
	s.TrueBlock.inherit(&s.ConditionBlock)
	errs = append(errs, s.TrueBlock.checkStatementsAndUnused()...)
	for _, v := range s.ElseIfList {
		v.Block.inherit(&s.ConditionBlock)
		if v.Condition != nil {
			conditionType, es := v.Condition.checkSingleValueContextExpression(v.Block)
			errs = append(errs, es...)
			if err := v.Condition.canBeUsedAsCondition(); err != nil {
				errs = append(errs, err)
			}
			if conditionType != nil &&
				conditionType.Type != VariableTypeBool {
				errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
					errMsgPrefix(s.Condition.Pos)))
			}
			errs = append(errs, v.Block.checkStatementsAndUnused()...)
		}
	}
	if s.ElseBlock != nil {
		s.ElseBlock.inherit(&s.ConditionBlock)
		errs = append(errs, s.ElseBlock.checkStatementsAndUnused()...)
	}
	return errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}
