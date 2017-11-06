package ast

type GlobalVariable struct {
	Pos  Pos
	Name string
	Init *Expression
	Typ  string //string a; Person b Person is a class name
}

type Const struct {
	Pos  Pos
	Name string
	Init *Expression
}
