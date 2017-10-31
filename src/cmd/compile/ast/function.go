package ast

type Function struct {
	Typ   FunctionType
	Name  string
	Block *Block
	Pos   Pos
}
type FunctionType struct {
	Parameters ParameterList
	Returns    ReturnList
}

type TypedNames struct {
	Name string
	Typ  VariableType
}

type Parameter struct {
	TypedNames
	Default *Expression //f(a int = 1) default parameter
}

type ParameterList []*Parameter
type ReturnList []*TypedNames
