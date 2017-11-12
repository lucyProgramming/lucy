package ast

type VariableDefinition struct {
	SymbolicItem
	AccessProperty
	Pos  *Pos
	Init *Expression
}

type Const struct {
	VariableDefinition
	Data interface{}
}
