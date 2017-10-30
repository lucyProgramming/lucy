package ast

/*
	enum Day{
		Mondy = 1,
		Tuesday
	}
*/
type Enum struct {
	Names []string
	Init  *Expression //should be a int expression
}
