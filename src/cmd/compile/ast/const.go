package ast

type Const struct {
	VariableDefinition
	Data interface{}
}
