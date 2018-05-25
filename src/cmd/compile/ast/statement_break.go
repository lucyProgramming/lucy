package ast

type StatementBreak struct {
	Defers          []*Defer
	StatementFor    *StatementFor
	StatementSwitch *StatementSwitch
}
