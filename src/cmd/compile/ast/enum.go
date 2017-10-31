package ast

/*
	enum Day{
		Mondy = 1,
		Tuesday
	}
*/
type Enum struct {
	Pos   Pos
	Name  string
	Names []string
	Init  *Expression //should be a int expression
}
