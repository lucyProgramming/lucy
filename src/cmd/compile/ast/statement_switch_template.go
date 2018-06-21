package ast

import "fmt"

type StatementSwitchTemplate struct {
	Pos                  *Position
	Condition            *Type //switch
	StatementSwitchCases []*StatementSwitchTemplateCase
	Default              *Block
}

type StatementSwitchTemplateCase struct {
	Matches []*Type
	Block   *Block
}

func (s *StatementSwitchTemplate) check(block *Block, statement *Statement) (errs []error) {
	errs = []error{}
	TName := s.Condition.Name
	if err := s.Condition.resolve(block); err != nil {
		errs = append(errs, err)
		return
	}
	var match *Type
	var matchBlock *Block
	typesChecked := []*Type{}
	es := []error{}
	checkIfAlreadyExist := func(ts []*Type, t *Type) bool {
		for _, v := range ts {
			if v.Equal(&es, t) {
				return true
			}
		}
		return false
	}
	for _, t := range s.StatementSwitchCases {
		for _, tt := range t.Matches {
			if err := tt.resolve(block); err != nil {
				errs = append(errs, err)
				continue
			}
			if checkIfAlreadyExist(typesChecked, tt) {
				errs = append(errs, fmt.Errorf("%s match '%s' already exist",
					errMsgPrefix(tt.Pos), tt.TypeString()))
				return
			}
			if s.Condition.Equal(&es, tt) == false {
				continue
			}
			// found
			match = tt
			matchBlock = t.Block
		}
	}
	if len(errs) > 0 {
		return errs
	}
	if match == nil {
		if s.Default == nil {
			errs = append(errs, fmt.Errorf("%s parameter type named '%s',resolve as '%s' has no match",
				errMsgPrefix(s.Pos), TName, s.Condition.TypeString()))
			return
		}
		statement.Type = STATEMENT_TYPE_BLOCK
		statement.Block = s.Default
		return statement.Block.checkStatements()
	}
	// let`s reWrite
	if matchBlock == nil {
		statement.Type = STATEMENT_TYPE_NOP
	} else {
		statement.Type = STATEMENT_TYPE_BLOCK
		statement.Block = matchBlock
		return append(errs, statement.Block.checkStatements()...)
	}
	return
}