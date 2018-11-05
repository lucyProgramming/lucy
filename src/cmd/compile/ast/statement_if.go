package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementIf struct {
	PrefixExpressions   []*Expression
	Condition           *Expression
	Pos                 *Pos
	initExpressionBlock Block
	Block               Block
	ElseIfList          []*StatementElseIf
	Else                *Block
	Exits               []*cg.Exit
}

func (this *StatementIf) check(father *Block) []error {
	this.initExpressionBlock.inherit(father)
	errs := []error{}
	for _, v := range this.PrefixExpressions {
		v.IsStatementExpression = true
		_, es := v.check(&this.initExpressionBlock)
		errs = append(errs, es...)
		if err := v.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
	}
	if this.Condition != nil {
		conditionType, es := this.Condition.checkSingleValueContextExpression(&this.initExpressionBlock)
		errs = append(errs, es...)
		if conditionType != nil &&
			conditionType.Type != VariableTypeBool {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression",
				this.Condition.Pos.ErrMsgPrefix()))
		}
		if err := this.Condition.canBeUsedAsCondition(); err != nil {
			errs = append(errs, err)
		}
	}
	this.Block.inherit(&this.initExpressionBlock)
	errs = append(errs, this.Block.check()...)
	for _, v := range this.ElseIfList {
		v.Block.inherit(&this.initExpressionBlock)
		if v.Condition != nil {
			conditionType, es := v.Condition.checkSingleValueContextExpression(v.Block)
			errs = append(errs, es...)
			if err := v.Condition.canBeUsedAsCondition(); err != nil {
				errs = append(errs, err)
			}
			if conditionType != nil &&
				conditionType.Type != VariableTypeBool {
				errs = append(errs,
					fmt.Errorf("%s condition is not a bool expression",
						conditionType.Pos.ErrMsgPrefix()))
			}
			errs = append(errs, v.Block.check()...)
		}
	}
	if this.Else != nil {
		this.Else.inherit(&this.initExpressionBlock)
		errs = append(errs, this.Else.check()...)
	}
	return errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}
