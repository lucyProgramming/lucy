package ast

type StatementContinue struct {
	StatementFor *StatementFor
	Defers       []*Defer
}
