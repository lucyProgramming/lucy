package ast

type GlobalVariable struct {
	AccessProperty
	Pos  Pos
	Name string
	Init *Expression
	Typ  *VariableType
}

type Const struct {
	AccessProperty
	Pos  Pos
	Name string
	Init *Expression
	Typ  *VariableType
	Data interface{}
}
