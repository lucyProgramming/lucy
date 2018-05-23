package ast

type InheritedAttribute struct {
	StatementOffset   int
	IsConstruction    bool
	StatementFor      *StatementFor // if this statement is in for or not
	StatementSwitch   *StatementSwitch
	statementForBreak interface{} // for or switch statement
	Function          *Function
	Class             *Class
	Defer             *Defer
}
