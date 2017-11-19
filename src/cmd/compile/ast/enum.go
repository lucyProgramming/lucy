package ast

/*
	enum Day{
		Monday = 1,
		Tuesday
	}
*/

type EnumName struct {
	Enum  *Enum
	Name  string
	Pos   *Pos
	Value int64
}

type Enum struct {
	Access int // public private or protected
	Name   string
	Pos    *Pos
	Names  []*EnumName
	Init   *Expression //should be a int expression
}
