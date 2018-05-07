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
		buildChecker: func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*VariableType, pos *Pos) {
			meta := &BuildinFunctionPrintfMeta{}
			e.BuildinFunctionMeta = meta
			if len(args) == 0 || args[0] == nil {
				return // not error
			}
			if args[0].Typ == VARIABLE_TYPE_OBJECT {
				have, _ := args[0].Class.haveSuper("java/io/PrintStream")
				if have {
					_, err := e.Args[0].mustBeOneValueContext(e.Args[0].Values)
					if err != nil {
						*errs = append(*errs, err)
					} else {
						meta.Stream = e.Args[0]
						e.Args = e.Args[1:]
					}
				}
			}
		},
		IsBuildin: true,
		Name:      common.BUILD_IN_FUNCTION_PRINT,
	}
	catchBuildFunction := &Function{}
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
	catchBuildFunction.buildChecker = func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*VariableType, pos *Pos) {
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
					*errs = append(*errs, fmt.Errorf("%s load exception class failed,err:%v",
						errMsgPrefix(pos), err))
					return
				}
				ft.Typ.ReturnList[0].Typ.Class = c.(*Class)
				err = block.InheritedAttribute.Defer.registerExceptionClass(c.(*Class))
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(pos), err))
				}
				return
			} else {
				ft.Typ.ReturnList[0].Typ.Class = block.InheritedAttribute.Defer.ExceptionClass
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
		buildChecker: oneAnyTypeParameterChecker,
		IsBuildin:    true,
		Name:         common.BUILD_IN_FUNCTION_PANIC,
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_MONITORENTER] = &Function{
		buildChecker: monitorChecker,
		IsBuildin:    true,
		Name:         common.BUILD_IN_FUNCTION_MONITORENTER,
	}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_MONITOREXIT] = &Function{
		buildChecker: monitorChecker,
		IsBuildin:    true,
		Name:         common.BUILD_IN_FUNCTION_MONITOREXIT,
	}
	// sprintf
	sprintfBuildFunction := &Function{}
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_SPRINTF] = sprintfBuildFunction
	sprintfBuildFunction.Name = common.BUILD_IN_FUNCTION_SPRINTF
	sprintfBuildFunction.IsBuildin = true
	{
		sprintfBuildFunction.Typ.ReturnList = make([]*VariableDefinition, 1)
		sprintfBuildFunction.Typ.ReturnList[0] = &VariableDefinition{}
		sprintfBuildFunction.Typ.ReturnList[0].Name = "retrunValue"
		sprintfBuildFunction.Typ.ReturnList[0].Typ = &VariableType{}
		sprintfBuildFunction.Typ.ReturnList[0].Typ.Typ = VARIABLE_TYPE_STRING
	}
	sprintfBuildFunction.buildChecker = func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
		args []*VariableType, pos *Pos) {
		if len(args) == 0 {
			err := fmt.Errorf("%s '%s' expect one argument at lease",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_SPRINTF)
			*errs = append(*errs, err)
			return
		}
		if args[0] == nil {
			return
		}
		if args[0].Typ != VARIABLE_TYPE_STRING {
			err := fmt.Errorf("%s '%s' first argument must be string",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_SPRINTF)
			*errs = append(*errs, err)
			return
		}
		_, err := e.Args[0].mustBeOneValueContext(e.Args[0].Values)
		if err != nil {
			*errs = append(*errs, err)
			return
		}
		meta := &BuildinFunctionSprintfMeta{}
		e.BuildinFunctionMeta = meta
		meta.Format = e.Args[0]
		meta.ArgsLength = len(args) - 1
		e.Args = e.Args[1:]
	}
	// printf
	buildinFunctionsMap[common.BUILD_IN_FUNCTION_PRINTF] = &Function{
		buildChecker: func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
			args []*VariableType, pos *Pos) {
			meta := &BuildinFunctionPrintfMeta{}
			e.BuildinFunctionMeta = meta
			if len(args) == 0 {
				err := fmt.Errorf("%s '%s' expect one argument at least",
					errMsgPrefix(pos), common.BUILD_IN_FUNCTION_PRINTF)
				*errs = append(*errs, err)
				return
			}
			if args[0] == nil {
				return
			}

			if args[0].Typ == VARIABLE_TYPE_OBJECT {
				have, _ := args[0].Class.haveSuper("java/io/PrintStream")
				if have {
					_, err := e.Args[0].mustBeOneValueContext(e.Args[0].Values)
					if err != nil {
						*errs = append(*errs, err)
						return
					} else {
						meta.Stream = e.Args[0]
						e.Args = e.Args[1:]
						args = args[1:]
					}
				}
			}
			if len(args) == 0 {
				err := fmt.Errorf("%s missing format argument",
					errMsgPrefix(pos))
				*errs = append(*errs, err)
				return
			}
			if args[0] == nil {
				return
			}
			if args[0].Typ != VARIABLE_TYPE_STRING {
				err := fmt.Errorf("%s format must be string",
					errMsgPrefix(pos))
				*errs = append(*errs, err)
				return
			}
			_, err := e.Args[0].mustBeOneValueContext(e.Args[0].Values)
			if err != nil {
				*errs = append(*errs, err)
				return
			}
			meta.Format = e.Args[0]
			e.Args = e.Args[1:]
			meta.ArgsLength = len(args)
		},
		IsBuildin: true,
		Name:      common.BUILD_IN_FUNCTION_PRINTF,
	}
}

func monitorChecker(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
	args []*VariableType, pos *Pos) {
	if len(args) != 1 {
		*errs = append(*errs, fmt.Errorf("%s only expect one argument", errMsgPrefix(pos)))
		return
	}
	if args[0] == nil {
		return
	}
	if args[0].IsPointer() == false || args[0].Typ == VARIABLE_TYPE_STRING {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not valid type to call",
			errMsgPrefix(pos), args[0].TypeString()))
		return
	}
}
