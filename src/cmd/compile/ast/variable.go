package ast

type VariableDefinition struct {
	AccessProperty
	Pos        *Pos
	Expression *Expression
	NameWithType
}

type Const struct {
	VariableDefinition
	Data interface{}
}
