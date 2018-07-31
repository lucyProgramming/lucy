package ast

import (
	"errors"
	"fmt"
)

type StatementSwitchTemplate struct {
	Pos                  *Pos
	Condition            *Type //switch
	StatementSwitchCases []*StatementSwitchTemplateCase
	Default              *Block
}

type StatementSwitchTemplateCase struct {
	Matches []*Type
	Block   *Block
}

/*
	switchStatement will be override
*/
func (s *StatementSwitchTemplate) check(block *Block, switchStatement *Statement) (errs []error) {
	errs = []error{}
	if s.Condition == nil { // must be a error must parse stage
		return errs
	}
	TName := s.Condition.Name
	if err := s.Condition.resolve(block); err != nil {
		errs = append(errs, err)
		return
	}
	var match *Type
	var matchBlock *Block
	typesChecked := []*Type{}
	checkIfAlreadyExist := func(ts []*Type, t *Type) *Type {
		for _, v := range ts {
			if v.StrictEqual(t) {
				return v
			}
		}
		return nil
	}
	for _, t := range s.StatementSwitchCases {
		for _, tt := range t.Matches {
			if err := tt.resolve(block); err != nil {
				errs = append(errs, err)
				continue
			}
			if exist := checkIfAlreadyExist(typesChecked, tt); exist != nil {
				errMsg := fmt.Sprintf("%s match '%s' already exist,first declared at:\n",
					errMsgPrefix(tt.Pos), tt.TypeString())
				errMsg += fmt.Sprintf("\t %s", errMsgPrefix(exist.Pos))
				errs = append(errs, errors.New(errMsg))
				return
			}
			typesChecked = append(typesChecked, tt)
			if s.Condition.StrictEqual(tt) == false {
				//no match here
				continue
			}
			// found
			if match == nil {
				match = tt
				matchBlock = t.Block
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	if match == nil {
		if s.Default == nil {
			errs = append(errs,
				fmt.Errorf("%s parameter type named '%s',resolve as '%s' has no match and no 'default block'",
					errMsgPrefix(s.Condition.Pos), TName, s.Condition.TypeString()))
			return
		} else {
			switchStatement.Type = StatementTypeBlock
			switchStatement.Block = s.Default
			switchStatement.Block.inherit(block)
			switchStatement.Block.IsSwitchTemplateBlock = true
			switchStatement.Block.InheritedAttribute.SwitchTemplateBlock = switchStatement.Block
			switchStatement.Block.InheritedAttribute.ForBreak = switchStatement.Block
			return switchStatement.Block.checkStatements()
		}
	}
	// let`s reWrite
	if matchBlock == nil {
		switchStatement.Type = StatementTypeNop
		return errs
	} else {
		switchStatement.Type = StatementTypeBlock
		switchStatement.Block = matchBlock
		switchStatement.Block.inherit(block)
		switchStatement.Block.IsSwitchTemplateBlock = true
		switchStatement.Block.InheritedAttribute.SwitchTemplateBlock = switchStatement.Block
		switchStatement.Block.InheritedAttribute.ForBreak = switchStatement.Block
		return append(errs, switchStatement.Block.checkStatements()...)
	}
}
