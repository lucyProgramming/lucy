package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

type CallChecker func(block *Block, errs *[]error, args []*VariableType, returnList ReturnList, pos *Pos)

type buildFunction struct {
	args       []*VariableDefinition
	returnList []*VariableDefinition
	checker    CallChecker
}

func init() {
	registerBuildinFunctions()
}

func registerBuildinFunctions() {
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_PRINT] = &buildFunction{
		checker: func(block *Block, errs *[]error, args []*VariableType, returnList ReturnList, pos *Pos) {},
	}
	catchBuildFunction := &buildFunction{}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_CATCH] = catchBuildFunction
	{
		catchBuildFunction.returnList = make([]*VariableDefinition, 1)
		catchBuildFunction.returnList[0] = &VariableDefinition{}
		catchBuildFunction.returnList[0].Name = "retrunValue"
		catchBuildFunction.returnList[0].Typ = &VariableType{}
		catchBuildFunction.returnList[0].Typ.Typ = VARIABLE_TYPE_OBJECT
		//class is going to make value by checker
	}
	catchBuildFunction.checker = func(block *Block, errs *[]error, args []*VariableType, returnList ReturnList, pos *Pos) {
		if block.InheritedAttribute.Defer == nil {
			*errs = append(*errs, fmt.Errorf("%s buildin function '%s' only allow in defer block",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_CATCH))
			return
		}
		if block.IsFunctionTopBlock == false {
			*errs = append(*errs, fmt.Errorf("%s buildin function '%s' only can be use in function top level block",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_CATCH))
			return
		}
		if len(args) > 1 {
			*errs = append(*errs, fmt.Errorf("%s build function '%s' expect at most 1 argument",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_CATCH))
			return
		}
		if len(args) == 0 {
			// make default exception class
			// load java/lang/Exception this is default exception level to catch
			if block.InheritedAttribute.Defer.ExceptionClass == nil {
				_, c, err := NameLoader.LoadName(DEFAULT_EXCEPTION_CLASS)
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s  load exception class failed,err:%v",
						errMsgPrefix(pos), err))
					return
				}
				returnList[0].Typ.Class = c.(*Class)
				err = block.InheritedAttribute.Defer.registerExceptionClass(c.(*Class))
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(pos), err))
				}
				return
			} else {
				returnList[0].Typ.Class = block.InheritedAttribute.Defer.ExceptionClass
				return
			}
		}
		if args[0].Typ != VARIABLE_TYPE_CLASS {
			*errs = append(*errs, fmt.Errorf("%s build function '%s' expect class",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_CATCH))
			return
		}
		if has, _ := args[0].Class.haveSuper(JAVA_THROWABLE_CLASS); has == false {
			*errs = append(*errs, fmt.Errorf("%s '%s' does not have super-class '%s'",
				errMsgPrefix(pos), args[0].Class.Name, JAVA_THROWABLE_CLASS))
			return
		}
		err := block.InheritedAttribute.Defer.registerExceptionClass(args[0].Class)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(pos), err))
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

func monitorChecker(block *Block, errs *[]error, args []*VariableType, returnList ReturnList, pos *Pos) {
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
