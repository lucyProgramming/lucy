package ast

import (
	"fmt"
)

const (
	_ = iota
	STATEMENT_TYPE_EXPRESSION
	STATEMENT_TYPE_IF
	STATEMENT_TYPE_BLOCK
	STATEMENT_TYPE_FOR
	STATEMENT_TYPE_CONTINUE
	STATEMENT_TYPE_RETURN
	STATEMENT_TYPE_BREAK
	STATEMENT_TYPE_SWITCH
	STATEMENT_TYPE_SKIP // skip this block

)

type Statement struct {
	Pos               *Pos
	Typ               int
	StatementIf       *StatementIF
	Expression        *Expression // expression statment like a=123
	StatementFor      *StatementFor
	StatementReturn   *StatementReturn
	StatementSwitch   *StatementSwitch
	StatementBreak    *StatementBreak
	Block             *Block
	StatementContinue *StatementContinue
}

type StatementContinue struct {
	StatementFor *StatementFor
}
type StatementBreak struct {
	StatementFor    *StatementFor
	StatementSwitch *StatementSwitch
}

func (s *Statement) statementName() string {
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		return "expression statement"
	case STATEMENT_TYPE_IF:
		return "if statement"
	case STATEMENT_TYPE_FOR:
		return "for statement"
	case STATEMENT_TYPE_CONTINUE:
		return "continue statement"
	case STATEMENT_TYPE_BREAK:
		return "break statement"
	case STATEMENT_TYPE_SWITCH:
		return "switch statement"
	case STATEMENT_TYPE_SKIP:
		return "skip statement"
	}
	return ""
}

func (s *Statement) check(b *Block) []error { // b is father
	errs := []error{}
	if b.InheritedAttribute.istop {
		if s.Typ == STATEMENT_TYPE_SKIP { //special case
			return errs // 0 length error
		}
	}
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		errs = append(errs, s.checkStatementExpression(b)...)
	case STATEMENT_TYPE_IF:
		t, es := s.StatementIf.check(b)
		if len(es) > 0 {
			errs = append(errs, es...)
		}
		if t != nil {
			s.Typ = STATEMENT_TYPE_BLOCK
		}
		s.Block = t
	case STATEMENT_TYPE_FOR:
		errs = append(errs, s.StatementFor.check(b)...)
	case STATEMENT_TYPE_SWITCH:
		panic("........")
	case STATEMENT_TYPE_BREAK:
		if b.InheritedAttribute.StatementFor == nil && b.InheritedAttribute.StatementSwitch == nil {
			errs = append(errs, fmt.Errorf("%s %s can`t in this scope", errMsgPrefix(s.Pos), s.statementName()))
		} else {
			s.StatementBreak = &StatementBreak{}
			if b.InheritedAttribute.closed_is_for {
				s.StatementBreak.StatementFor = b.InheritedAttribute.StatementFor
			} else {
				s.StatementBreak.StatementSwitch = b.InheritedAttribute.StatementSwitch
			}
		}
	case STATEMENT_TYPE_CONTINUE:
		if b.InheritedAttribute.StatementFor == nil {
			errs = append(errs, fmt.Errorf("%s %s can`t in this scope", errMsgPrefix(s.Pos), s.statementName()))
		} else {
			if s.StatementContinue == nil {
				s.StatementContinue = &StatementContinue{b.InheritedAttribute.StatementFor}
			}
		}
	case STATEMENT_TYPE_RETURN:
		if b.InheritedAttribute.function == nil {
			errs = append(errs, fmt.Errorf("%s %s can`t in this scope", errMsgPrefix(s.Pos), s.statementName()))
			return errs
		}
		s.StatementReturn.Function = b.InheritedAttribute.function
		errs = append(errs, s.StatementReturn.check(b)...)
	default:
		panic("unkown type statement" + s.statementName())
	}
	return errs
}

func (s *Statement) checkStatementExpression(b *Block) (errs []error) {
	errs = []error{}
	if s.Expression.Typ == EXPRESSION_TYPE_COLON_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_FUNCTION_CALL ||
		s.Expression.Typ == EXPRESSION_TYPE_METHOD_CALL ||
		s.Expression.Typ == EXPRESSION_TYPE_FUNCTION ||
		s.Expression.Typ == EXPRESSION_TYPE_PLUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MINUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MUL_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_DIV_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MOD_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_INCREMENT ||
		s.Expression.Typ == EXPRESSION_TYPE_DECREMENT ||
		s.Expression.Typ == EXPRESSION_TYPE_PRE_INCREMENT ||
		s.Expression.Typ == EXPRESSION_TYPE_PRE_DECREMENT {
	} else {
		errs = append(errs, fmt.Errorf("%s expression evaluate but not used", errMsgPrefix(s.Expression.Pos)))
	}
	_, es := b.checkExpression_(s.Expression)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}

type StatementSwitch struct {
	BackPatchs          [][]byte
	Outter              *Block
	Condition           *Expression //switch
	StatmentSwitchCases []*StatmentSwitchCase
	Default             *Block
}

type StatmentSwitchCase struct {
	Match *Expression
	Block *Block
}

func (s *StatmentSwitchCase) check() []error {
	errs := []error{}
	return errs
}

type StatementReturn struct {
	Function    *Function
	Pos         *Pos // use some time
	Expressions []*Expression
}

func (s *StatementReturn) check(b *Block) []error {
	if len(s.Expressions) == 0 {
		return nil
	}
	errs := make([]error, 0)
	returndValueTypes := s.Expressions[0].checkExpressions(b, s.Expressions, &errs)
	pos := s.Expressions[len(s.Expressions)-1].Pos
	rs := b.InheritedAttribute.function.Typ.Returns
	if len(returndValueTypes) < len(rs) {
		errs = append(errs, fmt.Errorf("%s too few value to return", errMsgPrefix(pos)))
	}
	if len(returndValueTypes) > len(rs) {
		errs = append(errs, fmt.Errorf("%s too many value to return", errMsgPrefix(pos)))
	}

	for k, v := range rs {
		if k < len(returndValueTypes) {
			if !v.Typ.typeCompatible(returndValueTypes[k]) {
				errs = append(errs, fmt.Errorf("%s cannot use %s as %s to return", errMsgPrefix(returndValueTypes[k].Pos), returndValueTypes[k].TypeString(), v.Typ.TypeString()))
			}
		}
	}

	return errs
}

type StatementFor struct {
	Num        int
	BackPatchs [][]byte
	LoopBegin  uint16
	Pos        *Pos
	Init       *Expression
	Condition  *Expression
	Post       *Expression
	Block      *Block
}

func (s *StatementFor) check(block *Block) []error {
	s.Block.inherite(block)
	s.Block.InheritedAttribute.StatementFor = s
	s.Block.InheritedAttribute.closed_is_for = true
	errs := []error{}
	if s.Init != nil {
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
				errs = append(errs, fmt.Errorf("%s condition must be bool expression,but %s", errMsgPrefix(s.Condition.Pos), t.TypeString()))
			}
		}
	}
	if s.Post != nil {
		_, es := s.Block.checkExpression(s.Post)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	es := s.Block.check(nil)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}

type ElseIfList []*StatementElseIf

func (e ElseIfList) check(father *Block) []error {
	errs := make([]error, 0)
	var err error
	for _, v := range e {
		t, es := v.Block.checkExpression(v.Condition)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
			continue
		}
		if t.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(v.Condition.Pos), err))
		}
		errs = append(errs, v.Block.check(father)...)
	}
	return errs
}

type StatementIF struct {
	BackPatchs [][]byte
	Condition  *Expression
	Block      *Block
	ElseBlock  *Block
	ElseIfList ElseIfList
}

func (s *StatementIF) check(father *Block) (*Block, []error) {
	errs := []error{}
	//inherite
	s.Block.inherite(father)
	conditionType, es := s.Block.checkExpression(s.Condition)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if conditionType != nil {
		if conditionType.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression", errMsgPrefix(s.Condition.Pos)))
		}
	}
	// condition should a bool expression
	if s.Condition.Typ == EXPRESSION_TYPE_BOOL && s.Condition.Data.(bool) == true { // if(true){...}
		errs = append(errs, s.Block.check(nil)...)
		return s.Block, errs
	}
	if s.Condition.Typ == EXPRESSION_TYPE_BOOL && s.Condition.Data.(bool) == false { // if(false){}else{...}
		errs = append(errs, s.ElseBlock.check(father)...)
		return s.Block, errs
	}
	errs = append(errs, s.Block.check(nil)...)
	if s.ElseIfList != nil && len(s.ElseIfList) > 0 {
		errs = append(errs, s.ElseIfList.check(father)...)
	}
	if s.ElseBlock != nil {
		errs = append(errs, s.ElseBlock.check(father)...)
	}
	return nil, errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}

//func (s *StatementElseIf) check() []error {
//	errs := []error{}
//	return errs
//}
