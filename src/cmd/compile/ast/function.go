package ast

type Function struct {
	Typ   FunctionType
	Name  string
	Block *Block
	Pos   Pos
}

func (f *Function) check(b *Block) []erros {
	if b != nil {
		f.Block.inherite(b)
	}
	errs := make([]error, 0)
	checkParaMeterAndRetuns(errs)
}

func (f *FunctionType) checkParaMeterAndRetuns(errs []error) {

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
