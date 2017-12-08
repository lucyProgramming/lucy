package ast

type VariableDefinition struct {
	AccessFlags uint16 // public private or protected
	Pos         *Pos
	Expression  *Expression
	NameWithType
}

type Const struct {
	VariableDefinition
	Data interface{}
}
