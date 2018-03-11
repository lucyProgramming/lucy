package ast

import "github.com/756445638/lucy/src/cmd/compile/common"

type CallChecker func(block *Block, errs *[]error, args []*VariableType, pos *Pos)

type buildFunction struct {
	args    []*VariableDefinition
	returns []*VariableDefinition
	checker CallChecker
}

func init() {
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_PRINT] = &buildFunction{
		checker: func(block *Block, errs *[]error, args []*VariableType, pos *Pos) {
			block.InheritedAttribute.Function.mkAutoVarForMultiReturn()
		},
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_CATCH] = &buildFunction{
		checker: oneAnyTypeParameterChecker,
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_PANIC] = &buildFunction{
		checker: oneAnyTypeParameterChecker,
	}
}
