package ast

type StatementBreak struct {
	Defers          []*StatementDefer
	StatementFor    *StatementFor
	StatementSwitch *StatementSwitch
}
