package ast

type Variable struct {
	Name string
	Init *Expression
}
type VariableList []*Variable

type Const struct {
	Name string
	Init *Expression
}
type ConstList []*Const
