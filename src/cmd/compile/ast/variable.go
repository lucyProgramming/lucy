package ast

type VariableDefinition struct {
	AccessProperty
	Pos  *Pos
	Init *Expression
	NameWithType
}

type Const struct {
	VariableDefinition
	Data interface{}
}
