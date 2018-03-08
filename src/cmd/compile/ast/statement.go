package ast

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
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
	STATEMENT_TYPE_LABLE
	STATEMENT_TYPE_GOTO
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
	StatmentLable     *StatementLable
	StatementGoto     *StatementGoto
}

type StatementGoto struct {
	Name           string
	Pos            *Pos
	StatementLable *StatementLable
}

type StatementLable struct {
	Name        string
	Pos         *Pos
	BackPatches []*cg.JumpBackPatch
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
	case STATEMENT_TYPE_LABLE:
		return "'lable statement'"
	case STATEMENT_TYPE_GOTO:
		return "'goto statement'"
	}
	return ""
}

func (s *Statement) check(b *Block) []error { // b is father
	if b.InheritedAttribute.function.isPackageBlockFunction {
		if s.Typ == STATEMENT_TYPE_SKIP { //special case
			return nil // 0 length error
		}
	}
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		return s.checkStatementExpression(b)
	case STATEMENT_TYPE_IF:
		return s.StatementIf.check(b)
	case STATEMENT_TYPE_FOR:
		return s.StatementFor.check(b)
	case STATEMENT_TYPE_SWITCH:
		panic("........")
	case STATEMENT_TYPE_BREAK:
		if b.InheritedAttribute.StatementFor == nil && b.InheritedAttribute.StatementSwitch == nil {
			return []error{fmt.Errorf("%s %s can`t in this scope", errMsgPrefix(s.Pos), s.statementName())}
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
			return []error{fmt.Errorf("%s %s can`t in this scope",
				errMsgPrefix(s.Pos), s.statementName())}
		} else {
			if s.StatementContinue == nil {
				s.StatementContinue = &StatementContinue{}
			}
			if s.StatementContinue.StatementFor == nil {
				s.StatementContinue.StatementFor = b.InheritedAttribute.StatementFor
			}
		}
	case STATEMENT_TYPE_RETURN:
		if b.InheritedAttribute.function == nil {
			return []error{fmt.Errorf("%s %s can`t in this scope",
				errMsgPrefix(s.Pos), s.statementName())}
		}
		return s.StatementReturn.check(b)
	case STATEMENT_TYPE_LABLE:
	case STATEMENT_TYPE_GOTO:
		err := s.checkStatementGoto(b)
		if err != nil {
			return []error{err}
		}
	default:
		panic("unkown type statement" + s.statementName())
	}
	return nil
}

func (s *Statement) checkStatementGoto(b *Block) error {
	t := b.searchByName(s.StatementGoto.Name)
	if t == nil {
		return fmt.Errorf("%s label named '%s' not found", errMsgPrefix(s.StatementGoto.Pos), s.StatementGoto.Name)
	}
	if l, ok := t.(*StatementLable); ok == false {
		return fmt.Errorf("%s '%s' is not a lable", errMsgPrefix(s.StatementGoto.Pos), s.StatementGoto.Name)
	} else {
		s.StatementGoto.StatementLable = l
	}
	return nil
}
func (s *Statement) checkStatementExpression(b *Block) []error {
	errs := []error{}
	if s.Expression.canBeUsedAsStatementExpression() {
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
	BackPatchs          []*cg.JumpBackPatch
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
