package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementFor struct {
	RangeAttr           *ForRangeAttr
	Exits               []*cg.Exit
	ContinueCodeOffset  int
	Pos                 *Pos
	initExpressionBlock Block
	/*
		for i := 0 ; i < 10 ;i ++ {

		}
	*/
	Init      *Expression
	Condition *Expression
	Increment *Expression
	Block     *Block
}

type ForRangeAttr struct {
	IdentifierKey   *ExpressionIdentifier
	IdentifierValue *ExpressionIdentifier
	ExpressionKey   *Expression
	ExpressionValue *Expression
	RangeOn         *Expression
}

func (this *StatementFor) checkRange() []error {
	errs := []error{}
	//
	var rangeExpression *Expression
	bin := this.Condition.Data.(*ExpressionBinary)
	if bin.Right.Type == ExpressionTypeRange {
		rangeExpression = this.Condition.Data.(*Expression)
	} else if bin.Right.Type == ExpressionTypeList {
		t := bin.Right.Data.([]*Expression)
		if len(t) > 1 {
			errs = append(errs,
				fmt.Errorf("%s for range statement only allow one argument on the right",
					errMsgPrefix(t[1].Pos)))
		}
		rangeExpression = t[0].Data.(*Expression)
	}
	rangeOn, es := rangeExpression.checkSingleValueContextExpression(&this.initExpressionBlock)
	errs = append(errs, es...)
	if rangeOn == nil {
		return errs
	}
	if rangeOn.Type == VariableTypeString {
		// byte[]("")
		conversion := &ExpressionTypeConversion{}
		conversion.Type = &Type{
			Type: VariableTypeJavaArray,
			Array: &Type{
				Type: VariableTypeByte,
				Pos:  rangeOn.Pos,
			},
			Pos: rangeOn.Pos,
		}
		conversion.Expression = rangeExpression
		bs := &Expression{
			Type: ExpressionTypeCheckCast,
			Data: conversion,
			Pos:  rangeOn.Pos,
		}
		bs.Value = conversion.Type
		rangeExpression = bs
		rangeOn = conversion.Type
	}
	if rangeOn.Type != VariableTypeArray &&
		rangeOn.Type != VariableTypeJavaArray &&
		rangeOn.Type != VariableTypeMap {
		errs = append(errs, fmt.Errorf("%s cannot range on '%s'",
			errMsgPrefix(rangeExpression.Pos), rangeOn.TypeString()))
		return errs
	}
	var lefts []*Expression
	if bin.Left.Type == ExpressionTypeList {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts = []*Expression{bin.Left}
	}
	if len(lefts) > 2 {
		errs = append(errs,
			fmt.Errorf("%s cannot have more than 2 expressions on the left",
				errMsgPrefix(lefts[2].Pos)))
		lefts = lefts[0:2]
	}
	modelKv := len(lefts) == 2
	this.RangeAttr = &ForRangeAttr{}
	this.RangeAttr.RangeOn = rangeExpression
	var err error
	if this.Condition.Type == ExpressionTypeVarAssign {
		for _, v := range lefts {
			if v.Type != ExpressionTypeIdentifier {
				errs = append(errs,
					fmt.Errorf("%s not a identifier on left",
						errMsgPrefix(v.Pos)))
				return errs
			}
		}
		var identifierK *ExpressionIdentifier
		var identifierV *ExpressionIdentifier
		var posK, posV *Pos
		if modelKv {
			identifierK = lefts[0].Data.(*ExpressionIdentifier)
			identifierV = lefts[1].Data.(*ExpressionIdentifier)
			posK = lefts[0].Pos
			posV = lefts[1].Pos
		} else {
			identifierV = lefts[0].Data.(*ExpressionIdentifier)
			posV = lefts[0].Pos
		}
		if identifierV.Name != UnderScore {
			vd := &Variable{}
			if rangeOn.Type == VariableTypeArray ||
				rangeOn.Type == VariableTypeJavaArray {
				vd.Type = rangeOn.Array.Clone()
			} else {
				vd.Type = rangeOn.Map.V.Clone()
			}
			vd.Pos = posV
			vd.Name = identifierV.Name
			err = this.initExpressionBlock.Insert(identifierV.Name, this.Condition.Pos, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierV.Variable = vd
			this.RangeAttr.IdentifierValue = identifierV
		}
		if modelKv &&
			identifierK.Name != UnderScore {
			vd := &Variable{}
			var vt *Type
			if rangeOn.Type == VariableTypeArray ||
				rangeOn.Type == VariableTypeJavaArray {
				vt = &Type{}
				vt.Type = VariableTypeInt
			} else {
				vt = rangeOn.Map.K.Clone()
				vt.Pos = rangeOn.Pos
			}
			vd.Name = identifierK.Name
			vd.Type = vt
			vd.Pos = posK
			err = this.initExpressionBlock.Insert(identifierK.Name, posK, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierK.Variable = vd
			this.RangeAttr.IdentifierKey = identifierK
		}
	} else { // k,v = range arr
		if modelKv {
			if false == lefts[0].IsIdentifier(UnderScore) {
				this.RangeAttr.ExpressionKey = lefts[0]
			}
			if false == lefts[1].IsIdentifier(UnderScore) {
				this.RangeAttr.ExpressionValue = lefts[1]
			}
		} else {
			if false == lefts[0].IsIdentifier(UnderScore) {
				this.RangeAttr.ExpressionValue = lefts[0]
			}
		}
		var receiverKType *Type
		if this.RangeAttr.ExpressionKey != nil {
			receiverKType = this.RangeAttr.ExpressionKey.getLeftValue(&this.initExpressionBlock, &errs)
			if receiverKType == nil {
				return errs
			}
		}
		var receiverVType *Type
		if this.RangeAttr.ExpressionValue != nil {
			receiverVType = this.RangeAttr.ExpressionValue.getLeftValue(&this.initExpressionBlock, &errs)
			if receiverVType == nil {
				return errs
			}
		}
		var kType, vType *Type
		if rangeOn.Type == VariableTypeArray ||
			rangeOn.Type == VariableTypeJavaArray {
			kType = &Type{
				Type: VariableTypeInt,
			}
			vType = rangeOn.Array
		} else {
			kType = rangeOn.Map.K
			vType = rangeOn.Map.V
		}
		if receiverKType != nil {
			if receiverKType.assignAble(&errs, kType) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for index",
					errMsgPrefix(this.RangeAttr.ExpressionKey.Pos),
					receiverKType.TypeString(), kType.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
		if receiverVType != nil {
			if receiverVType.assignAble(&errs, vType) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for value destination",
					errMsgPrefix(this.RangeAttr.ExpressionKey.Pos),
					receiverKType.TypeString(), kType.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
	}
	errs = append(errs, this.Block.check()...)
	return errs
}
func (this *StatementFor) check(block *Block) []error {
	this.initExpressionBlock.inherit(block)
	this.initExpressionBlock.InheritedAttribute.ForContinue = this
	this.initExpressionBlock.InheritedAttribute.ForBreak = this
	this.Block.inherit(&this.initExpressionBlock)
	errs := []error{}
	if this.Init == nil &&
		this.Increment == nil &&
		this.Condition != nil &&
		this.Condition.canBeUsedForRange() {
		// for k,v := range arr
		return this.checkRange()
	}
	if this.Init != nil {
		this.Init.IsStatementExpression = true
		if err := this.Init.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
		_, es := this.Init.check(&this.initExpressionBlock)
		errs = append(errs, es...)
	}
	if this.Condition != nil {
		if err := this.Condition.canBeUsedAsCondition(); err != nil {
			errs = append(errs, err)
		}
		t, es := this.Condition.checkSingleValueContextExpression(&this.initExpressionBlock)
		errs = append(errs, es...)
		if t != nil && t.Type != VariableTypeBool {
			errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s",
				errMsgPrefix(this.Condition.Pos), t.TypeString()))
		}
	}
	if this.Increment != nil {
		this.Increment.IsStatementExpression = true
		if err := this.Increment.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
		_, es := this.Increment.check(&this.initExpressionBlock)
		errs = append(errs, es...)
	}
	if len(errs) > 0 {
		return errs
	}
	errs = append(errs, this.Block.check()...)
	return errs
}
