package ast

type TypeAlias struct {
	Name string
	Type *Type
	Pos  *Pos
}
