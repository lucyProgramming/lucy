package ast

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type StatementSwitch struct {
	Pos                 *Pos
	BackPatchs          []*cg.JumpBackPatch
	Condition           *Expression //switch
	StatmentSwitchCases []*StatmentSwitchCase
	Default             *Block
}

type StatmentSwitchCase struct {
	Matches []*Expression
	Block   Block
}

func (s *StatementSwitch) check(b *Block) []error {
	errs := []error{}
	conditionType, es := b.checkExpression(s.Condition)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if conditionType == nil {
		return errs
	}
	if conditionType.Typ == VARIABLE_TYPE_BOOL {
		errs = append(errs, fmt.Errorf("%s bool not allow for switch", errMsgPrefix(conditionType.Pos)))
		return errs
	}
	if len(s.StatmentSwitchCases) == 0 {
		errs = append(errs, fmt.Errorf("%s switch statement has no cases", errMsgPrefix(s.Pos)))
	}
	for _, v := range s.StatmentSwitchCases {
		for _, e := range v.Matches {
			t, es := b.checkExpression(e)
			if errsNotEmpty(es) {
				errs = append(errs, es...)
			}
			if conditionType.Equal(t) == false {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'", errMsgPrefix(e.Pos), t.TypeString(), conditionType.TypeString()))

			}
		}
		v.Block.inherite(b)
		errs = append(errs, v.Block.check()...)
	}
	if s.Default != nil {
		s.Default.inherite(b)
		errs = append(errs, s.Default.check()...)
	}
	return errs
}
