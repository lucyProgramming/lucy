package ast

type buildFunction struct {
	args    []*VariableDefinition
	returns []*VariableDefinition
	checker func(errs *[]error, args []*VariableType, pos *Pos)
}

func init() {
	buildinFunctionsMap["print"] = &buildFunction{
		checker: func(errs *[]error, args []*VariableType, pos *Pos) {},
	}
	buildinFunctionsMap["panic"] = &buildFunction{
		checker: oneAnyTypeParameterChecker,
	}
	buildinFunctionsMap["recover"] = &buildFunction{
		checker: oneAnyTypeParameterChecker,
	}
}
