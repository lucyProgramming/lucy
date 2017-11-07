package ast

import (
	"fmt"
	"github.com/756445638/sc/lex"
	"github.com/coreos/etcd/embed"
)

const (
	STATEMENT_TYPE_EXPRESSION = iota
	STATEMENT_TYPE_IF
	STATEMENT_TYPE_FOR
	STATEMENT_TYPE_CONTINUE
	STATEMENT_TYPE_RETURN
	STATEMENT_TYPE_BREAK
	STATEMENT_TYPE_SWITCH
	STATEMENT_TYPE_SKIP // skip this block
)

type Statement struct {
	Pos               Pos
	Typ               int
	StatementIf       *StatementIF
	Expression        *Expression // expression statment like a=123
	StatementFor      *StatementFor
	StatementReturn   *StatementReturn
	StatementTryCatch *StatementTryCatch
	StatmentSwitch    *StatmentSwitch
}

func (s *Statement) stateName() string {
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

func (s *Statement) check(b *Block) []error {
	errs := []error{}
	if b.InheritedAttribute.istop {
		if s.Typ == STATEMENT_TYPE_SKIP { //special case
			return errs // 0 length error
		}
	}
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		errs = append(errs, s.checkExpression(b)...)
	case STATEMENT_TYPE_IF:
		s.StatementIf.Block.inherite(b)
		errs = append(errs, s.StatementIf.check()...)
	case STATEMENT_TYPE_FOR:
		s.StatementIf.Block.inherite(b)
		errs = append(errs, s.StatementIf.check()...)
	case STATEMENT_TYPE_SWITCH:
		s.StatementIf.Block.inherite(b)
		errs = append(errs, s.StatementFor.check()...)
	case STATEMENT_TYPE_BREAK:
	case STATEMENT_TYPE_CONTINUE:
		if b.InheritedAttribute.infor {
			errs = append(errs, fmt.Errorf("%s %d:%d %s can`t in this scope", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn, s.stateName()))
		}
	default:
		return
	}
	return errs
}

func (s *Statement) checkExpression(b *Block) []error {
	errs := []error{}
	//func1() System.log("hello world")
	if EXPRESSION_TYPE_METHOD_CALL == s.Expression.Typ ||
		EXPRESSION_TYPE_FUNCTION_CALL == s.Expression.Typ {
		return
	}
	// i++ i-- ++i --i
	if EXPRESSION_TYPE_INCREMENT == s.Expression.Typ ||
		EXPRESSION_TYPE_DECREMENT == s.Expression.Typ ||
		EXPRESSION_TYPE_PRE_INCREMENT == s.Expression.Typ ||
		EXPRESSION_TYPE_PRE_DECREMENT == s.Expression.Typ {
		left := s.Expression.Data.(*Expression)     // left means left value
		if left.Typ == EXPRESSION_TYPE_IDENTIFIER { //naming
			item := b.searchSymbolicItemAlsoGlobalVar(left.Data.(string))
			if item == nil {
				errs = append(errs, fmt.Errorf("%s %d:%d varaible not found"))
			}
		}
		errs = append(errs, fmt.Errorf("%s %d:%d cannot apply ++ or -- on %s", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn, left.humanReadableString()))
		return errs
	}
	if EXPRESSION_TYPE_COLON_ASSIGN == s.Expression.Typ { //declare variable
		binary := s.Expression.Data.(*ExpressionBinary)
		if binary.Left.Typ != EXPRESSION_TYPE_IDENTIFIER {
			errs = append(errs, fmt.Errorf("%s %d:%d no name on the left", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn))
			return errs
		}
		name := binary.Left.Data.(string)
		if _, ok := b.SymbolicTable.itemsMap[name]; ok {
			errs = append(errs, fmt.Errorf("%s %d:%d no name on the left", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn))
			return errs
		}
		item := &SymbolicItem{}
		item.Name = name
		item.Typ = getTypeFromExpression(binary.Right)
		b.SymbolicTable.itemsMap[name] = item
		return errs
	}

	//
	if s.Expression.Typ == EXPRESSION_TYPE_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_PLUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MINUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MUL_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_DIV_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MOD_ASSIGN {

	}
	errs = append(errs, fmt.Errorf(""))
	return errs

}

type StatementTryCatch struct {
	TryBlock     *Block
	CatchBlock   *Block
	FinallyBlock *Block
}

type StatmentSwitch struct {
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
	Expression []*Expression
}

func (s *StatementReturn) check() []error {
	errs := []error{}
	return errs
}

type StatementFor struct {
	Init      *Expression
	Condition *Expression
	Post      *Expression
	Block     *Block
}

func (s *StatementFor) check() []error {
	errs := []error{}
	return errs
}

type StatementIF struct {
	Condition  *Expression
	Block      *Block
	ElseBlock  *Block
	ElseIfList []*StatementElseIf
}

func (s *StatementIF) check() []error {
	errs := []error{}
	return errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}

func (s *StatementElseIf) check() []error {
	errs := []error{}
	return errs
}
