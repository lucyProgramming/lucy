package ast

type InheritedAttribute struct {
	StatementOffset       int // should not inherite
	IsConstructionMethod  bool
	StatementFor          *StatementFor // if this statement is in for or not
	StatementSwitch       *StatementSwitch
	SwitchTemplateBlock   *Block
	ForBreak              interface{} // for or switch statement
	Function              *Function
	Class                 *Class
	Defer                 *StatementDefer
	ClassMethod           *ClassMethod
	ClassAndFunctionNames string
}
