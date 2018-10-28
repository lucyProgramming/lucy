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

func (f *StatementFor) checkRange() []error {
	errs := []error{}
	//
	var rangeExpression *Expression
	bin := f.Condition.Data.(*ExpressionBinary)
	if bin.Right.Type == ExpressionTypeRange {
		rangeExpression = f.Condition.Data.(*Expression)
	} else if bin.Right.Type == ExpressionTypeList {
		t := bin.Right.Data.([]*Expression)
		if len(t) > 1 {
			errs = append(errs,
				fmt.Errorf("%s for range statement only allow one argument on the right",
					errMsgPrefix(t[1].Pos)))
		}
		rangeExpression = t[0].Data.(*Expression)
	}
	rangeOn, es := rangeExpression.checkSingleValueContextExpression(&f.initExpressionBlock)
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
	f.RangeAttr = &ForRangeAttr{}
	f.RangeAttr.RangeOn = rangeExpression
	var err error
	if f.Condition.Type == ExpressionTypeVarAssign {
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
			err = f.initExpressionBlock.Insert(identifierV.Name, f.Condition.Pos, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierV.Variable = vd
			f.RangeAttr.IdentifierValue = identifierV
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
			err = f.initExpressionBlock.Insert(identifierK.Name, posK, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierK.Variable = vd
			f.RangeAttr.IdentifierKey = identifierK
		}
	} else { // k,v = range arr
		if modelKv {
			if false == lefts[0].IsIdentifier(UnderScore) {
				f.RangeAttr.ExpressionKey = lefts[0]
			}
			if false == lefts[1].IsIdentifier(UnderScore) {
				f.RangeAttr.ExpressionValue = lefts[1]
			}
		} else {
			if false == lefts[0].IsIdentifier(UnderScore) {
				f.RangeAttr.ExpressionValue = lefts[0]
			}
		}
		var receiverKType *Type
		if f.RangeAttr.ExpressionKey != nil {
			receiverKType = f.RangeAttr.ExpressionKey.getLeftValue(&f.initExpressionBlock, &errs)
			if receiverKType == nil {
				return errs
			}
		}
		var receiverVType *Type
		if f.RangeAttr.ExpressionValue != nil {
			receiverVType = f.RangeAttr.ExpressionValue.getLeftValue(&f.initExpressionBlock, &errs)
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
					errMsgPrefix(f.RangeAttr.ExpressionKey.Pos),
					receiverKType.TypeString(), kType.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
		if receiverVType != nil {
			if receiverVType.assignAble(&errs, vType) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for value destination",
					errMsgPrefix(f.RangeAttr.ExpressionKey.Pos),
					receiverKType.TypeString(), kType.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
	}
	errs = append(errs, f.Block.check()...)
	return errs
}
func (f *StatementFor) check(block *Block) []error {
	f.initExpressionBlock.inherit(block)
	f.initExpressionBlock.InheritedAttribute.ForContinue = f
	f.initExpressionBlock.InheritedAttribute.ForBreak = f
	f.Block.inherit(&f.initExpressionBlock)
	errs := []error{}
	if f.Init == nil &&
		f.Increment == nil &&
		f.Condition != nil &&
		f.Condition.canBeUsedForRange() {
		// for k,v := range arr
		return f.checkRange()
	}
	if f.Init != nil {
		f.Init.IsStatementExpression = true
		if err := f.Init.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
		_, es := f.Init.check(&f.initExpressionBlock)
		errs = append(errs, es...)
	}
	if f.Condition != nil {
		if err := f.Condition.canBeUsedAsCondition(); err != nil {
			errs = append(errs, err)
		}
		t, es := f.Condition.checkSingleValueContextExpression(&f.initExpressionBlock)
		errs = append(errs, es...)
		if t != nil && t.Type != VariableTypeBool {
			errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s",
				errMsgPrefix(f.Condition.Pos), t.TypeString()))
		}
	}
	if f.Increment != nil {
		f.Increment.IsStatementExpression = true
		if err := f.Increment.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
		_, es := f.Increment.check(&f.initExpressionBlock)
		errs = append(errs, es...)
	}
	if len(errs) > 0 {
		return errs
	}
	errs = append(errs, f.Block.check()...)
	return errs
}
