package ast

type Const struct {
	VariableDefinition
	Value interface{} // value base on type
}
