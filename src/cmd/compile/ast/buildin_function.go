package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
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
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_MONITORENTER] = &buildFunction{
		checker: monitorChecker,
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_MONITOREXIT] = &buildFunction{
		checker: monitorChecker,
	}
}

func monitorChecker(block *Block, errs *[]error, args []*VariableType, pos *Pos) {
	if len(args) != 1 {
		*errs = append(*errs, fmt.Errorf("%s only expect one argument", errMsgPrefix(pos)))
		return
	}
	if args[0].IsPointer() == false || args[0].Typ == VARIABLE_TYPE_STRING {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not valid type to call",
			errMsgPrefix(pos), args[0].TypeString()))
		return
	}
}
