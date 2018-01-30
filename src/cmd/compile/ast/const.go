package ast

type Const struct {
	VariableDefinition
	BoolValue    bool
	Int64Value   int64
	Float64Value float64
	StringValue  string
}
