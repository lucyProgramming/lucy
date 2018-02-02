package ast

type CallChecker func(block *Block, errs *[]error, args []*VariableType, pos *Pos)

type buildFunction struct {
	args    []*VariableDefinition
	returns []*VariableDefinition
	checker CallChecker
}

func init() {
	buildinFunctionsMap["print"] = &buildFunction{
		checker: func(block *Block, errs *[]error, args []*VariableType, pos *Pos) {
			block.InheritedAttribute.function.mkArrayListVarForMultiReturn()
		},
	}
	buildinFunctionsMap["panic"] = &buildFunction{
		checker: oneAnyTypeParameterChecker,
	}
	buildinFunctionsMap["recover"] = &buildFunction{
		checker: oneAnyTypeParameterChecker,
	}
}
