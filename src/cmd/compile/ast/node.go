package ast

type Node struct {
	LNode *Node
	RNode *Node
}

type Expression struct {
	Typ             int
	BoolValue       bool
	IntValue        int64
	ByteValue       byte
	FloatValue      float64
	StringValue     string
	LeftExpression  *Expression
	RIghtExpression *Expression
}

type ExpressionFunctionCall struct {
	Name string //function name
	Args []*Expression
}
