package ast

type GlobalVariable struct {
	Pos   Pos
	Name  string
	Init  *Expression
	Typ   *VariableType
	Value interface{}
}

type Const struct {
	AccessProperty
	Pos  Pos
	Name string
	Init *Expression
	Typ  *VariableType
	Data interface{}
}
