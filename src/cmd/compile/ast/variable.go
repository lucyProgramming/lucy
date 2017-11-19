package ast

type VariableDefinition struct {
	Access     int // public private or protected
	Pos        *Pos
	Expression *Expression
	NameWithType
}

type Const struct {
	VariableDefinition
	Data interface{}
}
