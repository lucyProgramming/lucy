package ast

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type StatementFor struct {
	Num                int
	BackPatchs         []*cg.JumpBackPatch
	ContinueBackPatchs []*cg.JumpBackPatch
	LoopBegin          uint16
	ContinueOPOffset   uint16
	Pos                *Pos
	Init               *Expression
	Condition          *Expression
	Post               *Expression
	Block              *Block
}

func (s *StatementFor) check(block *Block) []error {
	s.Block.inherite(block)
	s.Block.InheritedAttribute.StatementFor = s
	s.Block.InheritedAttribute.mostCloseForOrSwitchForBreak = s
	errs := []error{}
	if s.Init != nil {
		s.Init.IsStatementExpression = true
		if s.Init.canBeUsedAsStatementExpression() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.Init.Pos)))
		}
		_, es := s.Block.checkExpression(s.Init)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	if s.Condition != nil {
		t, es := s.Block.checkExpression(s.Condition)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		if t != nil {
			if t.Typ != VARIABLE_TYPE_BOOL {
				errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s",
					errMsgPrefix(s.Condition.Pos), t.TypeString()))
			}
		}
	}
	if s.Post != nil {
		s.Post.IsStatementExpression = true
		if s.Post.canBeUsedAsStatementExpression() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.Post.Pos)))
		}
		_, es := s.Block.checkExpression(s.Post)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	es := s.Block.check()
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}
