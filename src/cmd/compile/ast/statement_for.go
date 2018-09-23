package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementFor struct {
	RangeAttr          *ForRangeAttr
	Exits              []*cg.Exit
	ContinueCodeOffset int
	Pos                *Pos
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

func (s *StatementFor) checkRange() []error {
	errs := []error{}
	//
	var rangeExpression *Expression
	bin := s.Condition.Data.(*ExpressionBinary)
	if bin.Right.Type == ExpressionTypeRange {
		rangeExpression = s.Condition.Data.(*Expression)
	} else if bin.Right.Type == ExpressionTypeList {
		t := bin.Right.Data.([]*Expression)
		if len(t) > 1 {
			errs = append(errs, fmt.Errorf("%s for range statement only allow one argument on the right",
				errMsgPrefix(t[1].Pos)))
		}
		rangeExpression = t[0].Data.(*Expression)
	}
	rangeOn, es := rangeExpression.checkSingleValueContextExpression(s.Block.Outer)
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
		errs = append(errs, fmt.Errorf("%s cannot have more than 2 expressions on the left",
			errMsgPrefix(lefts[2].Pos)))
		lefts = lefts[0:2]
	}
	modelKv := len(lefts) == 2
	s.RangeAttr = &ForRangeAttr{}
	s.RangeAttr.RangeOn = rangeExpression
	var err error
	if s.Condition.Type == ExpressionTypeVarAssign {
		for _, v := range lefts {
			if v.Type != ExpressionTypeIdentifier {
				errs = append(errs, fmt.Errorf("%s not a identifier on left",
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
		if identifierV.Name != NoNameIdentifier {
			vd := &Variable{}
			if rangeOn.Type == VariableTypeArray ||
				rangeOn.Type == VariableTypeJavaArray {
				vd.Type = rangeOn.Array.Clone()
			} else {
				vd.Type = rangeOn.Map.V.Clone()
			}
			vd.Pos = posV
			vd.Name = identifierV.Name
			err = s.Block.Insert(identifierV.Name, s.Condition.Pos, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierV.Variable = vd
			s.RangeAttr.IdentifierValue = identifierV
		}
		if modelKv &&
			identifierK.Name != NoNameIdentifier {
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
			err = s.Block.Insert(identifierK.Name, posK, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierK.Variable = vd
			s.RangeAttr.IdentifierKey = identifierK
		}
	} else { // k,v = range arr
		if modelKv {
			if false == lefts[0].IsIdentifier(NoNameIdentifier) {
				s.RangeAttr.ExpressionKey = lefts[0]
			}
			if false == lefts[1].IsIdentifier(NoNameIdentifier) {
				s.RangeAttr.ExpressionValue = lefts[1]
			}
		} else {
			if false == lefts[0].IsIdentifier(NoNameIdentifier) {
				s.RangeAttr.ExpressionValue = lefts[0]
			}
		}
		var receiverKType *Type
		if s.RangeAttr.ExpressionKey != nil {
			receiverKType = s.RangeAttr.ExpressionKey.getLeftValue(s.Block, &errs)
			if receiverKType == nil {
				return errs
			}
		}
		var receiverVType *Type
		if s.RangeAttr.ExpressionValue != nil {
			receiverVType = s.RangeAttr.ExpressionValue.getLeftValue(s.Block, &errs)
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
					errMsgPrefix(s.RangeAttr.ExpressionKey.Pos), receiverKType.TypeString(), kType.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
		if receiverVType != nil {
			if receiverVType.assignAble(&errs, vType) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for value destination",
					errMsgPrefix(s.RangeAttr.ExpressionKey.Pos), receiverKType.TypeString(), kType.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
	}
	errs = append(errs, s.Block.checkStatements()...)
	return errs
}
func (s *StatementFor) check(block *Block) []error {
	s.Block.inherit(block)
	s.Block.InheritedAttribute.StatementFor = s
	s.Block.InheritedAttribute.ForBreak = s
	errs := []error{}
	if s.Init == nil &&
		s.Increment == nil &&
		s.Condition != nil &&
		s.Condition.canBeUsedForRange() {
		// for k,v := range arr
		return s.checkRange()
	}
	if s.Init != nil {
		s.Init.IsStatementExpression = true
		if err := s.Init.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
		_, es := s.Init.check(s.Block)
		errs = append(errs, es...)
	}
	if s.Condition != nil {
		if err := s.Condition.canBeUsedAsCondition(); err != nil {
			errs = append(errs, err)
		}
		t, es := s.Condition.checkSingleValueContextExpression(s.Block)
		errs = append(errs, es...)
		if t != nil && t.Type != VariableTypeBool {
			errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s",
				errMsgPrefix(s.Condition.Pos), t.TypeString()))
		}
	}
	if s.Increment != nil {
		s.Increment.IsStatementExpression = true
		if err := s.Increment.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
		_, es := s.Increment.check(s.Block)
		errs = append(errs, es...)
	}
	if len(errs) > 0 {
		return errs
	}
	errs = append(errs, s.Block.checkStatements()...)
	return errs
}
