package ast

import "github.com/astaxie/beego/logs/es"

const (
	STATEMENT_TYPE_EXPRESSION = iota
	STATEMENT_TYPE_IF
	STATEMENT_TYPE_FOR
	STATEMENT_TYPE_CONTINUE
	STATEMENT_TYPE_RETURN
	STATEMENT_TYPE_BREAK
	STATEMENT_TYPE_SWITCH
	STATEMENT_TYPE_PASS // skip this block
	STATEMENT_TYPE_EXIT       //exit program
)

type Statement struct {
	Typ               int
	StatementIf       *StatementIF
	Expression        *Expression // expression statment like a=123
	StatementFor      *StatementFor
	StatementReturn   *StatementReturn
	StatementTryCatch *StatementTryCatch
	StatmentSwitch    *StatmentSwitch
}

func (s *Statement) check(b *Block) []error {
	if b.istop {
		if s.Typ
	}

	errs := []error{}
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		errs = append(errs, s.checkExpression(b)...)
	case STATEMENT_TYPE_IF:
		s.StatementIf.Block.Outter = b
		errs = append(errs, s.checkIf()...)
	case STATEMENT_TYPE_FOR:
		s.StatementFor.Block.Outter = b
		errs = append(errs, s.checkFor()...)
	case STATEMENT_TYPE_SWITCH:
		s.StatmentSwitch.Outter = b
		errs = append(errs, s.checkSwitch()...)
	case STATEMENT_TYPE_BREAK:

	default:

	}
	return errs
}

func (s *Statement) checkExpression(b *Block) []error {
	errs := []error{}
	if s.Expression.Typ == EXPRESSION_TYPE_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_COLON_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_PLUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MINUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MUL_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_DIV_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MOD_ASSIGN {
		return nil
	}
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

type StatementReturn struct {
	Expression []*Expression
}
type StatementFor struct {
	Init      *Expression
	Condition *Expression
	Post      *Expression
	Block     *Block
}

type StatementIF struct {
	Condition  *Expression
	Block      *Block
	ElseBlock  *Block
	ElseIfList []*StatementElseIf
}
type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}
