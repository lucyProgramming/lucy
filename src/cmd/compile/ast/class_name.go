package ast

type ClassNameDefinition struct {
	Name, BinaryName string // short name for class
	Pos              *Pos
}
