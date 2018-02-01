package ast

import (
	"fmt"
)

type StatementIF struct {
	BackPatchs [][]byte
	Condition  *Expression
	Block      *Block
	ElseBlock  *Block
	ElseIfList []*StatementElseIf
}

func (s *StatementIF) check(father *Block) []error {
	errs := []error{}
	conditionType, es := s.Block.checkExpression(s.Condition)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if conditionType != nil {
		if conditionType.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
				errMsgPrefix(s.Condition.Pos)))
		}
	}
	errs = append(errs, s.Block.check(father)...)
	if s.ElseIfList != nil && len(s.ElseIfList) > 0 {
		for _, v := range s.ElseIfList {
			conditionType, es := s.Block.checkExpression(s.Condition)
			if errsNotEmpty(es) {
				errs = append(errs, es...)
			}
			if conditionType != nil {
				if conditionType.Typ != VARIABLE_TYPE_BOOL {
					errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
						errMsgPrefix(s.Condition.Pos)))
				}
			}
			errs = append(errs, v.Block.check(father)...)
		}
	}
	if s.ElseBlock != nil {
		errs = append(errs, s.ElseBlock.check(father)...)
	}
	return errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}
