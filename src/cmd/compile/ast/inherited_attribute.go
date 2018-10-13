package ast

type InheritedAttribute struct {
	StatementOffset       int // should not inherite
	IsConstructionMethod  bool
	ForContinue           *StatementFor // if this statement is in for or not
	ForBreak              interface{}   // for or switch statement
	Function              *Function
	Class                 *Class
	Defer                 *StatementDefer
	ClassMethod           *ClassMethod
	ClassAndFunctionNames string
}
