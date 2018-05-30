package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementFor struct {
	RangeAttr          *StatmentForRangeAttr
	Num                int
	BackPatchs         []*cg.JumpBackPatch
	ContinueBackPatchs []*cg.JumpBackPatch
	ContinueOPOffset   int
	Pos                *Pos
	Init               *Expression
	Condition          *Expression
	Post               *Expression
	Block              *Block
}

type StatmentForRangeAttr struct {
	ModelKV     bool
	IdentifierK *ExpressionIdentifer
	IdentifierV *ExpressionIdentifer
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
	s.RangeAttr.ModelKV = modelkv
	if s.Condition.Typ == EXPRESSION_TYPE_ASSIGN {
		if modelkv {
			s.RangeAttr.ExpressionK = lefts[0]
			s.RangeAttr.ExpressionV = lefts[1]
		} else {
			s.RangeAttr.ExpressionV = lefts[0]
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
		var identifierK *ExpressionIdentifer
		var identifierV *ExpressionIdentifer
		var posk, posv *Pos
		if modelkv {
			identifierK = lefts[0].Data.(*ExpressionIdentifer)
			identifierV = lefts[1].Data.(*ExpressionIdentifer)
			posk = lefts[0].Pos
			posv = lefts[1].Pos
		} else {
			identifierV = lefts[0].Data.(*ExpressionIdentifer)
			posv = lefts[0].Pos
		}
		s.RangeAttr.IdentifierV = identifierV
		s.RangeAttr.IdentifierK = identifierK
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
		}
	}

	if s.Condition.Typ == EXPRESSION_TYPE_ASSIGN {
		t1 := lefts[0].getLeftValue(s.Block, &errs)
		lefts[0].Value = t1
		var t2 *VariableType
		if modelkv {
			t2 = lefts[1].getLeftValue(s.Block, &errs)
			lefts[1].Value = t2
		}
		if t1 == nil {
			return errs
		}
		if modelkv && t2 == nil {
			return errs
		}
		lefts[0].Value = t1
		if modelkv && t2 != nil {
			lefts[1].Value = t2
		}
		if rangeOn.Typ == VARIABLE_TYPE_ARRAY ||
			rangeOn.Typ == VARIABLE_TYPE_JAVA_ARRAY {
			if modelkv {
				if t1.IsInteger() == false {
					errs = append(errs, fmt.Errorf("%s index must be integer", errMsgPrefix(lefts[0].Pos)))
					return errs
				}
				if t2.Equal(rangeOn.ArrayType) == false {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(lefts[1].Pos), rangeOn.ArrayType.TypeString(), t2.TypeString()))
					return errs
				}

			} else { // v model
				if t1.Equal(rangeOn.ArrayType) == false {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(lefts[1].Pos), rangeOn.ArrayType.TypeString(), t2.TypeString()))
					return errs
				}
			}
		} else { // map type
			if modelkv {
				if false == t1.Equal(rangeOn.Map.K) {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(lefts[1].Pos), rangeOn.Map.K.TypeString(), t1.TypeString()))
					return errs

				}
				if false == t2.Equal(rangeOn.Map.V) {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(lefts[1].Pos), rangeOn.Map.K.TypeString(), t2.TypeString()))
					return errs

				}
			} else {
				if false == t1.Equal(rangeOn.Map.V) {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(lefts[1].Pos), rangeOn.Map.K.TypeString(), t1.TypeString()))
					return errs

				}
			}
		}
	}
	errs = append(errs, s.Block.checkStatements()...)
	return errs
}
func (s *StatementFor) check(block *Block) []error {
	s.Block.inherite(block)
	s.Block.InheritedAttribute.StatementFor = s
	s.Block.InheritedAttribute.statementForBreak = s
	errs := []error{}
	if s.Init == nil && s.Post == nil && s.Condition != nil && s.Condition.canbeUsedForRange() { // for k,v := range arr
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
