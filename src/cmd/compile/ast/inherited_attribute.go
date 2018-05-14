package ast

type InheritedAttribute struct {
	StatementOffset   int
	IsConstruction    bool
	StatementFor      *StatementFor // if this statement is in for or not
	StatementSwitch   *StatementSwitch
	statementForBreak interface{} // for or switch statement
	Function          *Function
	class             *Class
	Defer             *Defer
}
