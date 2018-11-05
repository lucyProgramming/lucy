package ast

import (
	"errors"
	"fmt"
)

type StatementWhen struct {
	Pos       *Pos
	Condition *Type
	Cases     []*StatementWhenCase
	Default   *Block
}

type StatementWhenCase struct {
	Matches []*Type
	Block   *Block
}

/*
	switchStatement will be override
*/
func (this *StatementWhen) check(block *Block, switchStatement *Statement) (errs []error) {
	if this.Condition == nil { // must be a error must parse stage
		return nil
	}
	errs = []error{}
	if len(this.Condition.getParameterType(&block.InheritedAttribute.Function.Type)) == 0 {
		errs = append(errs, fmt.Errorf("%s '%s' constains no parameter type",
			this.Condition.Pos.ErrMsgPrefix(), this.Condition.TypeString()))
		return errs
	}
	if err := this.Condition.resolve(block); err != nil {
		errs = append(errs, err)
		return
	}
	var match *Type
	var matchBlock *Block
	typesChecked := []*Type{}
	checkExists := func(ts []*Type, t *Type) *Type {
		for _, v := range ts {
			if v.Equal(t) {
				return v
			}
		}
		return nil
	}
	for _, t := range this.Cases {
		for _, tt := range t.Matches {
			if err := tt.resolve(block); err != nil {
				errs = append(errs, err)
				continue
			}
			if exist := checkExists(typesChecked, tt); exist != nil {
				errMsg := fmt.Sprintf("%s match '%s' already exist,first declared at:\n",
					errMsgPrefix(tt.Pos), tt.TypeString())
				errMsg += fmt.Sprintf("\t %s", errMsgPrefix(exist.Pos))
				errs = append(errs, errors.New(errMsg))
				return
			}
			typesChecked = append(typesChecked, tt)
			if this.Condition.Equal(tt) == false {
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
		if this.Default == nil {
			errs = append(errs,
				fmt.Errorf("%s condition resolve as '%s' has no match and no 'default block'",
					errMsgPrefix(this.Condition.Pos), this.Condition.TypeString()))
		} else {
			switchStatement.Type = StatementTypeBlock
			switchStatement.Block = this.Default
			switchStatement.Block.inherit(block)
			switchStatement.Block.IsWhenBlock = true
			switchStatement.Block.InheritedAttribute.ForBreak = switchStatement.Block
			errs = append(errs, switchStatement.Block.check()...)
		}
		return
	}
	// let`s reWrite
	if matchBlock == nil {
		switchStatement.Type = StatementTypeNop
		return errs
	} else {
		switchStatement.Type = StatementTypeBlock
		switchStatement.Block = matchBlock
		switchStatement.Block.inherit(block)
		switchStatement.Block.IsWhenBlock = true
		switchStatement.Block.InheritedAttribute.ForBreak = switchStatement.Block
		return append(errs, switchStatement.Block.check()...)
	}
}
