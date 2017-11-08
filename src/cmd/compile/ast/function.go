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

type TypedName struct {
	Name string
	Typ  VariableType
}

type Parameter struct {
	TypedName
	Default *Expression //f(a int = 1) default parameter
}

type ParameterList []*Parameter
type ReturnList []*TypedName
