package ast

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/common"
)

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
	b := &buildFunction{}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_CATCH] = b
	b.checker = func(block *Block, errs *[]error, args []*VariableType, pos *Pos) {
		if len(args) > 0 {
			*errs = append(*errs, fmt.Errorf("%s build function '%s' expect no args",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_CATCH))
		}
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_PANIC] = &buildFunction{
		checker: oneAnyTypeParameterChecker,
	}
}
