package ast

type Class struct {
	Fields []*ClassField
}

type ClassMethod struct {
	Name string
}

type ClassField struct {
	Name string
	Typ  int
	Init *Expression
	Tag  string //for reflect

}

type FunctionType struct {
}
