package ast

const (
	STATEMENT_TYPE_EXPRESSION = iota
	STATEMENT_TYPE_IF
	STATEMENT_TYPE_FOR
	STATEMENT_TYPE_CONTINUE
	STATEMENT_TYPE_RETURN
	STATEMENT_TYPE_BREAK
	STATEMENT_TYPE_SWITCH
)

type Block struct {
	Statements []*Statement
}

type Statement struct {
	Typ               int
	StatementIf       *StatementIF
	Expression        *Expression // expression statment like a=123
	StatementFor      *StatementFor
	StatementReturn   *StatementReturn
	StatementTryCatch *StatementTryCatch
	StatmentSwitch    *StatmentSwitch
}

type StatementTryCatch struct {
	TryBlock     *Block
	CatchBlock   *Block
	FinallyBlock *Block
}

type StatmentSwitch struct {
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

type Block struct {
	Statments []*Statement
}
