package ast

type InheritedAttribute struct {
	StatementOffset      int
	IsConstructionMethod bool
	StatementFor         *StatementFor // if this statement is in for or not
	StatementSwitch      *StatementSwitch
	ForBreak             interface{} // for or switch statement
	Function             *Function
	Class                *Class
	Defer                *StatementDefer
}
