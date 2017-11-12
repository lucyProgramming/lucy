package ast

type VariableDefinition struct {
	SymbolicItem
	AccessProperty
	Pos  *Pos
	Init *Expression
}

type Const struct {
	AccessProperty
	Pos  Pos
	Name string
	Init *Expression
	Typ  *VariableType
	Data interface{}
}
