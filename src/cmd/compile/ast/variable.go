package ast

type GlobalVariable struct {
	Pos  Pos
	Name string
	Init *Expression
}

type Const struct {
	Pos  Pos
	Name string
	Init *Expression
}
