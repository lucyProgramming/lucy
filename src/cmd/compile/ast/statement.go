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
		return "'expression statement'"
	case STATEMENT_TYPE_IF:
		return "'if statement'"
	case STATEMENT_TYPE_FOR:
		return "'for statement'"
	case STATEMENT_TYPE_CONTINUE:
		return "'continue statement'"
	case STATEMENT_TYPE_BREAK:
		return "'break statement'"
	case STATEMENT_TYPE_SWITCH:
		return "'switch statement'"
	case STATEMENT_TYPE_SKIP:
		return "'skip statement'"
	}
	return ""
}

func (s *Statement) check(b *Block) []error { // b is father
	errs := []error{}
	if b.InheritedAttribute.function.isPackageBlockFunction {
		if s.Typ == STATEMENT_TYPE_SKIP { //special case
			return errs // 0 length error
		}
	}
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		errs = append(errs, s.checkStatementExpression(b)...)
	case STATEMENT_TYPE_IF:
		errs = append(errs, s.StatementIf.check(b)...)
	case STATEMENT_TYPE_FOR:
		errs = append(errs, s.StatementFor.check(b)...)
	case STATEMENT_TYPE_SWITCH:
		panic("........")
	case STATEMENT_TYPE_BREAK:
		if b.InheritedAttribute.StatementFor == nil && b.InheritedAttribute.StatementSwitch == nil {
			errs = append(errs, fmt.Errorf("%s %s can`t in this scope", errMsgPrefix(s.Pos), s.statementName()))
		} else {

			s.StatementBreak = &StatementBreak{}
			if f, ok := b.InheritedAttribute.mostCloseForOrSwitchForBreak.(*StatementFor); ok {
				s.StatementBreak.StatementFor = f
			} else {
				s.StatementBreak.StatementSwitch = b.InheritedAttribute.mostCloseForOrSwitchForBreak.(*StatementSwitch)
			}
		}
	case STATEMENT_TYPE_CONTINUE:
		if b.InheritedAttribute.StatementFor == nil {
			errs = append(errs, fmt.Errorf("%s %s can`t in this scope",
				errMsgPrefix(s.Pos), s.statementName()))
		} else {
			if s.StatementContinue == nil {
				s.StatementContinue = &StatementContinue{b.InheritedAttribute.StatementFor}
			}
		}
	case STATEMENT_TYPE_RETURN:
		if b.InheritedAttribute.function == nil {
			errs = append(errs, fmt.Errorf("%s %s can`t in this scope",
				errMsgPrefix(s.Pos), s.statementName()))
			return errs
		}
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
		err := fmt.Errorf("%s expression '%s' evaluate but not used",
			errMsgPrefix(s.Expression.Pos), s.Expression.OpName())
		errs = append(errs, err)
	}
	s.Expression.IsStatementExpression = true
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
