package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func init() {
	registerBuildinFunctions()
}

func registerBuildinFunctions() {
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_PRINT] = &Function{
		callChecker: func(block *Block, errs *[]error, args []*VariableType, returnList ReturnList, pos *Pos) {},
		IsBuildin:   true,
		Typ:         &FunctionType{},
		Name:        common.BUILD_IN_FUNCTION_PRINT,
	}
	catchBuildFunction := &Function{}
	catchBuildFunction.Typ = &FunctionType{}
	catchBuildFunction.IsBuildin = true
	catchBuildFunction.Name = common.BUILD_IN_FUNCTION_CATCH
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_CATCH] = catchBuildFunction
	{
		catchBuildFunction.Typ.ReturnList = make([]*VariableDefinition, 1)
		catchBuildFunction.Typ.ReturnList[0] = &VariableDefinition{}
		catchBuildFunction.Typ.ReturnList[0].Name = "retrunValue"
		catchBuildFunction.Typ.ReturnList[0].Typ = &VariableType{}
		catchBuildFunction.Typ.ReturnList[0].Typ.Typ = VARIABLE_TYPE_OBJECT
		//class is going to make value by checker
	}
	catchBuildFunction.callChecker = func(block *Block, errs *[]error, args []*VariableType, returnList ReturnList, pos *Pos) {
		if block.InheritedAttribute.Defer == nil ||
			block.InheritedAttribute.Defer.allowCatch == false {
			*errs = append(*errs, fmt.Errorf("%s buildin function '%s' only allow in defer block",
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
			*errs = append(*errs, fmt.Errorf("%s build function '%s' expect a object ref argument",
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
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_PANIC] = &Function{
		callChecker: oneAnyTypeParameterChecker,
		IsBuildin:   true,
		Typ:         &FunctionType{},
		Name:        common.BUILD_IN_FUNCTION_PANIC,
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_MONITORENTER] = &Function{
		callChecker: monitorChecker,
		IsBuildin:   true,
		Typ:         &FunctionType{},
		Name:        common.BUILD_IN_FUNCTION_MONITORENTER,
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_MONITOREXIT] = &Function{
		callChecker: monitorChecker,
		IsBuildin:   true,
		Typ:         &FunctionType{},
		Name:        common.BUILD_IN_FUNCTION_MONITOREXIT,
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
