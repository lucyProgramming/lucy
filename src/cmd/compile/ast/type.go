package ast

const (
	VARIABLE_TYPE_BOOL = iota
	VARIABLE_TYPE_BYTE
	VARIABLE_TYPE_INT
	VARIABLE_TYPE_FLOAT
	VARIABLE_TYPE_STRING
	VARIALBE_TYPE_FUNCTION
	VARIABLE_TYPE_CLASS       //new Person()
	VARIABLE_TYPE_COMBINATION // []int
)

type VariableType struct {
	Typ             int
	ClassName       string
	CombinationType *CombinationType
}

const (
	COMBINATION_TYPE_ARRAY = iota
)

type CombinationType struct {
	Typ COMBINATION_TYPE_ARRAY
}
