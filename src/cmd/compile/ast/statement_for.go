package ast

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type StatementFor struct {
	StatmentForRangeAttr *StatmentForRangeAttr
	Num                  int
	BackPatchs           []*cg.JumpBackPatch
	ContinueBackPatchs   []*cg.JumpBackPatch
	LoopBegin            int
	ContinueOPOffset     int
	Pos                  *Pos
	Init                 *Expression
	Condition            *Expression
	Post                 *Expression
	Block                *Block
}

type StatmentForRangeAttr struct {
	ModelKV              bool
	IdentifierK          *ExpressionIdentifer
	IdentifierV          *ExpressionIdentifer
	ExpressionK          *Expression
	ExpressionV          *Expression
	Expression           *Expression
	Lefts                []*Expression
	Typ                  int
	AutoVarForRangeArray *AutoVarForRangeArray
	AutoVarForRangeMap   *AutoVarForRangeMap
}
type AutoVarForRangeMap struct {
	MapObject                uint16
	KeySets                  uint16
	KeySetsK, KeySetsKLength uint16
	K, V                     uint16
}

func (t *AutoVarForRangeMap) mkAutoVarForRange(f *Function, kt, vt *VariableType) {
	t.MapObject = f.Varoffset
	t.KeySets = f.Varoffset + 1
	t.KeySetsK = f.Varoffset + 2
	t.KeySetsKLength = f.Varoffset + 3
	f.Varoffset += 4

	t.K = f.Varoffset
	f.Varoffset += kt.JvmSlotSize()
	t.V = f.Varoffset
	f.Varoffset += vt.JvmSlotSize()
}

type AutoVarForRangeArray struct {
	K          uint16
	Elements   uint16
	Start, End uint16
	V          uint16
}

func (t *AutoVarForRangeArray) mkAutoVarForRange(f *Function, vt *VariableType) {
	t.K = f.Varoffset
	t.Elements = f.Varoffset + 1
	t.Start = f.Varoffset + 2
	t.End = f.Varoffset + 3
	f.Varoffset += 4
	t.V = f.Varoffset
	f.Varoffset += vt.JvmSlotSize()
}

func (s *StatementFor) checkRange() []error {
	errs := []error{}
	//
	var arrayExpression *Expression
	bin := s.Condition.Data.(*ExpressionBinary)
	if bin.Right.Typ == EXPRESSION_TYPE_RANGE {
		arrayExpression = s.Condition.Data.(*Expression)
	} else if bin.Right.Typ == EXPRESSION_TYPE_LIST {
		t := bin.Right.Data.([]*Expression)
		if len(t) > 1 {
			errs = append(errs, fmt.Errorf("%s for range statement only allow one argument on the right", errMsgPrefix(t[1].Pos)))
		}
		arrayExpression = t[0].Data.(*Expression)
	}
	ts, es := arrayExpression.check(s.Block)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	arrayt, err := arrayExpression.mustBeOneValueContext(ts)
	if err != nil {
		errs = append(errs, err)
	}
	if arrayt == nil {
		return errs
	}
	if arrayt.Typ != VARIABLE_TYPE_ARRAY_INSTANCE && arrayt.Typ != VARIABLE_TYPE_MAP {
		errs = append(errs, fmt.Errorf("%s cannot use '%s' for range,only allow 'array' and 'map'", errMsgPrefix(arrayExpression.Pos), arrayt.TypeString()))
		return errs
	}
	arrayExpression.VariableType = arrayt
	var lefts []*Expression
	if bin.Left.Typ == EXPRESSION_TYPE_LIST {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts = []*Expression{bin.Left}
	}
	if len(lefts) > 2 {
		errs = append(errs, fmt.Errorf("%s cannot have more than 2 expressions on the left", errMsgPrefix(lefts[2].Pos)))
		lefts = lefts[:2]
	}
	modelkv := false
	if len(lefts) == 2 {
		modelkv = true
	}
	s.StatmentForRangeAttr = &StatmentForRangeAttr{}
	s.StatmentForRangeAttr.ModelKV = modelkv
	s.StatmentForRangeAttr.Expression = arrayExpression
	if arrayt.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
		s.StatmentForRangeAttr.AutoVarForRangeArray = &AutoVarForRangeArray{}
		s.StatmentForRangeAttr.AutoVarForRangeArray.mkAutoVarForRange(s.Block.InheritedAttribute.function,
			arrayt.ArrayType)
	} else {
		s.StatmentForRangeAttr.AutoVarForRangeMap = &AutoVarForRangeMap{}
		s.StatmentForRangeAttr.AutoVarForRangeMap.mkAutoVarForRange(s.Block.InheritedAttribute.function,
			arrayt.Map.K, arrayt.Map.V)
	}

	s.StatmentForRangeAttr.Lefts = lefts
	if s.Condition.Typ == EXPRESSION_TYPE_COLON_ASSIGN {
		var identifier *ExpressionIdentifer
		var pos *Pos
		if lefts[0].Typ != EXPRESSION_TYPE_IDENTIFIER {
			errs = append(errs, fmt.Errorf("%s not a identifier on left", errMsgPrefix(lefts[0].Pos)))
		} else {
			identifier = lefts[0].Data.(*ExpressionIdentifer)
			pos = lefts[0].Pos
		}
		var identifier2 *ExpressionIdentifer
		var pos2 *Pos
		if modelkv {
			if lefts[1].Typ != EXPRESSION_TYPE_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left", errMsgPrefix(lefts[1].Pos)))
			} else {
				identifier2 = lefts[1].Data.(*ExpressionIdentifer)
				pos2 = lefts[1].Pos
			}
		}
		if modelkv {
			if identifier != nil {
				if identifier.Name == NO_NAME_IDENTIFIER {
					errs = append(errs, fmt.Errorf("%s not a valid name one left", errMsgPrefix(pos)))
				} else {
					vd := &VariableDefinition{}
					var vt *VariableType
					if arrayt.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
						vt = &VariableType{}
						vt.Typ = VARIABLE_TYPE_INT
					} else {
						vt = arrayt.Map.K.Clone()
						vt.Pos = arrayt.Pos
					}
					vd.Typ = vt
					vd.Pos = pos
					err = s.Block.insert(identifier.Name, pos, vd)
					if err != nil {
						errs = append(errs, err)
					}
					identifier.Var = vd
					s.StatmentForRangeAttr.IdentifierK = identifier
				}
			}
			if identifier2 != nil {
				if identifier2.Name == NO_NAME_IDENTIFIER {
					errs = append(errs, fmt.Errorf("%s not a valid name one left", errMsgPrefix(pos2)))
				} else {
					vd := &VariableDefinition{}
					if arrayt.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
						vd.Typ = arrayt.ArrayType.Clone()
					} else {
						vd.Typ = arrayt.Map.V.Clone()
					}
					vd.Pos = pos2
					err = s.Block.insert(identifier2.Name, s.Condition.Pos, vd)
					if err != nil {
						errs = append(errs, err)
					}
					identifier2.Var = vd
					s.StatmentForRangeAttr.IdentifierV = identifier2
				}
			}
		} else {
			if identifier != nil && identifier.Name == NO_NAME_IDENTIFIER {
				errs = append(errs, fmt.Errorf("%s not a identifier on left", errMsgPrefix(lefts[1].Pos)))
			}
			if identifier != nil {
				if identifier.Name == NO_NAME_IDENTIFIER {
					errs = append(errs, fmt.Errorf("%s not a valid name one left", errMsgPrefix(pos2)))
				} else {
					vd := &VariableDefinition{}
					if arrayt.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
						vd.Typ = arrayt.ArrayType.Clone()
					} else {
						vd.Typ = arrayt.Map.V.Clone()
					}
					vd.Typ.Pos = pos2
					vd.Pos = pos2
					err = s.Block.insert(identifier.Name, s.Condition.Pos, vd)
					if err != nil {
						errs = append(errs, err)
					}
					identifier.Var = vd
					s.StatmentForRangeAttr.IdentifierV = identifier
				}
			}
		}
	}
	if s.Condition.Typ == EXPRESSION_TYPE_ASSIGN {
		t1, es := lefts[0].getLeftValue(s.Block)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		var t2 *VariableType
		if modelkv {
			t2, es = lefts[0].getLeftValue(s.Block)
			if errsNotEmpty(es) {
				errs = append(errs, es...)
			}
		}
		if t1 == nil {
			goto checkBlock
		}
		if modelkv && t2 == nil {
			goto checkBlock
		}
		lefts[0].VariableType = t1
		if modelkv && t2 != nil {
			lefts[1].VariableType = t2
		}
		if arrayt.Typ == VARIABLE_TYPE_ARRAY_INSTANCE {
			if modelkv {
				if t1.IsInteger() == false {
					errs = append(errs, fmt.Errorf("%s index must be integer", errMsgPrefix(lefts[0].Pos)))
				}
				if t2.TypeCompatible(arrayt.ArrayType) == false {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(lefts[1].Pos), arrayt.ArrayType.TypeString(), t2.TypeString()))
				}

			} else { // v model
				if t1.TypeCompatible(arrayt.ArrayType) == false {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(lefts[1].Pos), arrayt.ArrayType.TypeString(), t2.TypeString()))
				}
			}
		} else {
			if modelkv {
				if false == t1.Equal(arrayt.Map.K) {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(lefts[1].Pos), arrayt.Map.K.TypeString(), t1.TypeString()))
				}
				if false == t2.Equal(arrayt.Map.V) {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(lefts[1].Pos), arrayt.Map.K.TypeString(), t2.TypeString()))
				}
			} else {
				if false == t1.Equal(arrayt.Map.V) {
					errs = append(errs, fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(lefts[1].Pos), arrayt.Map.K.TypeString(), t1.TypeString()))
				}
			}
		}
	}

checkBlock:
	errs = append(errs, s.Block.check()...)
	return errs
}
func (s *StatementFor) check(block *Block) []error {
	s.Block.inherite(block)
	s.Block.InheritedAttribute.StatementFor = s
	s.Block.InheritedAttribute.mostCloseForOrSwitchForBreak = s
	errs := []error{}
	if s.Init == nil && s.Post == nil && s.Condition != nil && s.Condition.canbeUsedAsForRange() { // for k,v := range arr
		return s.checkRange()
	}
	if s.Init != nil {
		s.Init.IsStatementExpression = true
		if s.Init.canBeUsedAsStatementExpression() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.Init.Pos)))
		}
		_, es := s.Block.checkExpression(s.Init)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	if s.Condition != nil {
		t, es := s.Block.checkExpression(s.Condition)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
		if t != nil {
			if t.Typ != VARIABLE_TYPE_BOOL {
				errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s",
					errMsgPrefix(s.Condition.Pos), t.TypeString()))
			}
		}
	}
	if s.Post != nil {
		s.Post.IsStatementExpression = true
		if s.Post.canBeUsedAsStatementExpression() == false {
			errs = append(errs, fmt.Errorf("%s cannot be used as statement", errMsgPrefix(s.Post.Pos)))
		}
		_, es := s.Block.checkExpression(s.Post)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	es := s.Block.check()
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}
