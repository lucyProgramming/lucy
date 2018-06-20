package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementFor struct {
	RangeAttr          *ForRangeAttr
	Exits              []*cg.Exit
	ContinueCodeOffset int
	Pos                *Position
	Init               *Expression
	Condition          *Expression
	After              *Expression
	Block              *Block
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
	if bin.Right.Type == EXPRESSION_TYPE_RANGE {
		rangeExpression = s.Condition.Data.(*Expression)
	} else if bin.Right.Type == EXPRESSION_TYPE_LIST {
		t := bin.Right.Data.([]*Expression)
		if len(t) > 1 {
			errs = append(errs, fmt.Errorf("%s for range statement only allow one argument on the right",
				errMsgPrefix(t[1].Pos)))
		}
		rangeExpression = t[0].Data.(*Expression)
	}
	rangeOn, es := rangeExpression.checkSingleValueContextExpression(s.Block)
	if errorsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if rangeOn == nil {
		return errs
	}
	if rangeOn.Type != VARIABLE_TYPE_ARRAY &&
		rangeOn.Type != VARIABLE_TYPE_JAVA_ARRAY &&
		rangeOn.Type != VARIABLE_TYPE_MAP {
		errs = append(errs, fmt.Errorf("%s cannot have range on '%s'",
			errMsgPrefix(rangeExpression.Pos), rangeOn.TypeString()))
		return errs
	}
	rangeExpression.ExpressionValue = rangeOn
	var lefts []*Expression
	if bin.Left.Type == EXPRESSION_TYPE_LIST {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts = []*Expression{bin.Left}
	}
	if len(lefts) > 2 {
		errs = append(errs, fmt.Errorf("%s cannot have more than 2 expressions on the left",
			errMsgPrefix(lefts[2].Pos)))
		lefts = lefts[:2]
	}
	modelKv := false
	if len(lefts) == 2 {
		modelKv = true
	}
	s.RangeAttr = &ForRangeAttr{}
	if s.Condition.Type == EXPRESSION_TYPE_ASSIGN {
		if modelKv {
			if false == lefts[0].IsNoNameIdentifier() {
				s.RangeAttr.ExpressionKey = lefts[0]
			}
			if false == lefts[1].IsNoNameIdentifier() {
				s.RangeAttr.ExpressionValue = lefts[1]
			}
		} else {
			if false == lefts[0].IsNoNameIdentifier() {
				s.RangeAttr.ExpressionValue = lefts[0]
			}
		}
	}
	s.RangeAttr.RangeOn = rangeExpression
	var err error
	if s.Condition.Type == EXPRESSION_TYPE_COLON_ASSIGN {
		if modelKv {
			if lefts[0].Type != EXPRESSION_TYPE_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left",
					errMsgPrefix(lefts[0].Pos)))
				return errs
			}
			if lefts[1].Type != EXPRESSION_TYPE_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left",
					errMsgPrefix(lefts[0].Pos)))
				return errs
			}
		} else {
			if lefts[0].Type != EXPRESSION_TYPE_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left",
					errMsgPrefix(lefts[0].Pos)))
				return errs
			}
		}
		var identifierK *ExpressionIdentifier
		var identifierV *ExpressionIdentifier
		var posK, posV *Position
		if modelKv {
			identifierK = lefts[0].Data.(*ExpressionIdentifier)
			identifierV = lefts[1].Data.(*ExpressionIdentifier)
			posK = lefts[0].Pos
			posV = lefts[1].Pos
		} else {
			identifierV = lefts[0].Data.(*ExpressionIdentifier)
			posV = lefts[0].Pos
		}

		if identifierV.Name != NO_NAME_IDENTIFIER {
			vd := &Variable{}
			if rangeOn.Type == VARIABLE_TYPE_ARRAY || rangeOn.Type == VARIABLE_TYPE_JAVA_ARRAY {
				vd.Type = rangeOn.ArrayType.Clone()
			} else {
				vd.Type = rangeOn.Map.Value.Clone()
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
		if modelKv && identifierK.Name != NO_NAME_IDENTIFIER {
			vd := &Variable{}
			var vt *Type
			if rangeOn.Type == VARIABLE_TYPE_ARRAY ||
				rangeOn.Type == VARIABLE_TYPE_JAVA_ARRAY {
				vt = &Type{}
				vt.Type = VARIABLE_TYPE_INT
			} else {
				vt = rangeOn.Map.Key.Clone()
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
	}

	if s.Condition.Type == EXPRESSION_TYPE_ASSIGN {
		var tk *Type
		if s.RangeAttr.ExpressionKey != nil {
			tk = s.RangeAttr.ExpressionKey.getLeftValue(s.Block, &errs)
			if tk == nil {
				return errs
			}
		}
		var tv *Type
		if s.RangeAttr.ExpressionValue != nil {
			tv = s.RangeAttr.ExpressionValue.getLeftValue(s.Block, &errs)
			if tv == nil {
				return errs
			}
		}
		var tkk, tvv *Type

		if rangeOn.Type == VARIABLE_TYPE_ARRAY ||
			rangeOn.Type == VARIABLE_TYPE_JAVA_ARRAY {
			tkk = &Type{
				Type: VARIABLE_TYPE_INT,
			}
			tvv = rangeOn.ArrayType
		} else {
			tkk = rangeOn.Map.Key
			tvv = rangeOn.Map.Value
		}
		if tk != nil {
			if tk.Equal(&errs, tkk) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for index",
					errMsgPrefix(s.RangeAttr.ExpressionKey.Pos), tk.TypeString(), tkk.TypeString())
				errs = append(errs, err)
				return errs
			}
		}
		if tv != nil {
			if tv.Equal(&errs, tvv) == false {
				err = fmt.Errorf("%s cannot use '%s' as '%s' for value destination",
					errMsgPrefix(s.RangeAttr.ExpressionKey.Pos), tk.TypeString(), tkk.TypeString())
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
	if s.Init == nil && s.After == nil && s.Condition != nil && s.Condition.canBeUsedForRange() { // for k,v := range arr
		return s.checkRange()
	}
	if s.Init != nil {
		s.Init.IsStatementExpression = true
		if s.Init.canBeUsedAsStatement() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.Init.Pos)))
		}
		_, es := s.Init.check(s.Block)
		if errorsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	if s.Condition != nil {
		if s.Condition.canBeUsedAsCondition() == false {
			errs = append(errs, fmt.Errorf("%s expression(%s) cannot used as condition",
				errMsgPrefix(s.Condition.Pos), s.Condition.OpName()))
		}
		t, es := s.Condition.checkSingleValueContextExpression(s.Block)
		if errorsNotEmpty(es) {
			errs = append(errs, es...)
		}
		if t != nil && t.Type != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s",
				errMsgPrefix(s.Condition.Pos), t.TypeString()))

		}
	}
	if s.After != nil {
		s.After.IsStatementExpression = true
		if s.After.canBeUsedAsStatement() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.After.Pos)))
		}
		_, es := s.After.check(s.Block)
		if errorsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	es := s.Block.checkStatements()
	if errorsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}
