package ast

/*
	enum {
		Mondy = 1,
		Tuesday
	}
*/

type EnumNames struct {
	Enum  *Enum
	Name  string
	Pos   Pos
	Value int64
}

type Enum struct {
	Name  string
	Pos   Pos
	Names []*EnumNames
	Init  *Expression //should be a int expression
	Value int64
}
