package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementFor struct {
	RangeAttr        *StatmentForRangeAttr
	BackPatchs       []*cg.JumpBackPatch
	ContinueOPOffset int
	Pos              *Pos
	Init             *Expression
	Condition        *Expression
	Post             *Expression
	Block            *Block
}

type StatmentForRangeAttr struct {
	IdentifierK *ExpressionIdentifier
	IdentifierV *ExpressionIdentifier
	ExpressionK *Expression
	ExpressionV *Expression
	RangeOn     *Expression
}

func (s *StatementFor) checkRange() []error {
	errs := []error{}
	//
	var rangeExpression *Expression
	bin := s.Condition.Data.(*ExpressionBinary)
	if bin.Right.Typ == EXPRESSION_TYPE_RANGE {
		rangeExpression = s.Condition.Data.(*Expression)
	} else if bin.Right.Typ == EXPRESSION_TYPE_LIST {
		t := bin.Right.Data.([]*Expression)
		if len(t) > 1 {
			errs = append(errs, fmt.Errorf("%s for range statement only allow one argument on the right",
				errMsgPrefix(t[1].Pos)))
		}
		rangeExpression = t[0].Data.(*Expression)
	}
	rangeOn, es := rangeExpression.checkSingleValueContextExpression(s.Block)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if rangeOn == nil {
		return errs
	}
	if rangeOn.Typ != VARIABLE_TYPE_ARRAY &&
		rangeOn.Typ != VARIABLE_TYPE_JAVA_ARRAY &&
		rangeOn.Typ != VARIABLE_TYPE_MAP {
		errs = append(errs, fmt.Errorf("%s cannot have range on '%s'",
			errMsgPrefix(rangeExpression.Pos), rangeOn.TypeString()))
		return errs
	}
	rangeExpression.Value = rangeOn
	var lefts []*Expression
	if bin.Left.Typ == EXPRESSION_TYPE_LIST {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts = []*Expression{bin.Left}
	}
	if len(lefts) > 2 {
		errs = append(errs, fmt.Errorf("%s cannot have more than 2 expressions on the left",
			errMsgPrefix(lefts[2].Pos)))
		lefts = lefts[:2]
	}
	modelkv := false
	if len(lefts) == 2 {
		modelkv = true
	}
	s.RangeAttr = &StatmentForRangeAttr{}
	if s.Condition.Typ == EXPRESSION_TYPE_ASSIGN {
		if modelkv {
			if false == lefts[0].IsNoNameIdentifier() {
				s.RangeAttr.ExpressionK = lefts[0]
			}
			if false == lefts[1].IsNoNameIdentifier() {
				s.RangeAttr.ExpressionV = lefts[1]
			}
		} else {
			if false == lefts[0].IsNoNameIdentifier() {
				s.RangeAttr.ExpressionV = lefts[0]
			}
		}
	}
	s.RangeAttr.RangeOn = rangeExpression
	var err error
	if s.Condition.Typ == EXPRESSION_TYPE_COLON_ASSIGN {
		if modelkv {
			if lefts[0].Typ != EXPRESSION_TYPE_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left",
					errMsgPrefix(lefts[0].Pos)))
				return errs
			}
			if lefts[1].Typ != EXPRESSION_TYPE_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left",
					errMsgPrefix(lefts[0].Pos)))
				return errs
			}
		} else {
			if lefts[0].Typ != EXPRESSION_TYPE_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left",
					errMsgPrefix(lefts[0].Pos)))
				return errs
			}
		}
		var identifierK *ExpressionIdentifier
		var identifierV *ExpressionIdentifier
		var posk, posv *Pos
		if modelkv {
			identifierK = lefts[0].Data.(*ExpressionIdentifier)
			identifierV = lefts[1].Data.(*ExpressionIdentifier)
			posk = lefts[0].Pos
			posv = lefts[1].Pos
		} else {
			identifierV = lefts[0].Data.(*ExpressionIdentifier)
			posv = lefts[0].Pos
		}

		if identifierV.Name != NO_NAME_IDENTIFIER {
			vd := &VariableDefinition{}
			if rangeOn.Typ == VARIABLE_TYPE_ARRAY || rangeOn.Typ == VARIABLE_TYPE_JAVA_ARRAY {
				vd.Typ = rangeOn.ArrayType.Clone()
			} else {
				vd.Typ = rangeOn.Map.V.Clone()
			}
			vd.Pos = posv
			vd.Name = identifierV.Name
			err = s.Block.insert(identifierV.Name, s.Condition.Pos, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierV.Var = vd
			s.RangeAttr.IdentifierV = identifierV
		}
		if modelkv && identifierK.Name != NO_NAME_IDENTIFIER {
			vd := &VariableDefinition{}
			var vt *VariableType
			if rangeOn.Typ == VARIABLE_TYPE_ARRAY ||
				rangeOn.Typ == VARIABLE_TYPE_JAVA_ARRAY {
				vt = &VariableType{}
				vt.Typ = VARIABLE_TYPE_INT
			} else {
				vt = rangeOn.Map.K.Clone()
				vt.Pos = rangeOn.Pos
			}
			vd.Name = identifierK.Name
			vd.Typ = vt
			vd.Pos = posk
			err = s.Block.insert(identifierK.Name, posk, vd)
			if err != nil {
				errs = append(errs, err)
			}
			identifierK.Var = vd
			s.RangeAttr.IdentifierK = identifierK
		}
	}

	if s.Condition.Typ == EXPRESSION_TYPE_ASSIGN {
		var tk *VariableType
		if s.RangeAttr.ExpressionK != nil {
			tk = s.RangeAttr.ExpressionK.getLeftValue(s.Block, &errs)
			if tk == nil {
				return errs
			}
		}
		var tv *VariableType
		if s.RangeAttr.ExpressionV != nil {
			tv = s.RangeAttr.ExpressionV.getLeftValue(s.Block, &errs)
			if tv == nil {
				return errs
			}
		}
		var tkk, tvv *VariableType

		if rangeOn.Typ == VARIABLE_TYPE_ARRAY ||
			rangeOn.Typ == VARIABLE_TYPE_JAVA_ARRAY {
			tkk = &VariableType{
				Typ: VARIABLE_TYPE_INT,
			}
			tvv = rangeOn.ArrayType
		} else {
			tkk = rangeOn.Map.K
			tvv = rangeOn.Map.V
		}
		if tk != nil {
			if tk.Equal(&errs, tkk) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for index",
					errMsgPrefix(s.RangeAttr.ExpressionK.Pos), tk.TypeString(), tkk.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
		if tv != nil {
			if tv.Equal(&errs, tvv) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for value destination",
					errMsgPrefix(s.RangeAttr.ExpressionK.Pos), tk.TypeString(), tkk.TypeString())
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
	if s.Init == nil && s.Post == nil && s.Condition != nil && s.Condition.canBeUsedForRange() { // for k,v := range arr
		return s.checkRange()
	}
	if s.Init != nil {
		s.Init.IsStatementExpression = true
		if s.Init.canBeUsedAsStatement() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.Init.Pos)))
		}
		_, es := s.Init.check(s.Block)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	if s.Condition != nil {
		if s.Condition.canbeUsedAsCondition() == false {
			errs = append(errs, fmt.Errorf("%s expression(%s) cannot used as condition",
				errMsgPrefix(s.Condition.Pos), s.Condition.OpName()))
		}
		t, es := s.Condition.checkSingleValueContextExpression(s.Block)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		if t != nil && t.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s",
				errMsgPrefix(s.Condition.Pos), t.TypeString()))

		}
	}
	if s.Post != nil {
		s.Post.IsStatementExpression = true
		if s.Post.canBeUsedAsStatement() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.Post.Pos)))
		}
		_, es := s.Post.check(s.Block)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	es := s.Block.checkStatements()
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}
