package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func init() {
	registerBuildInFunctions()
}

func registerBuildInFunctions() {
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_PRINT] = &Function{
		buildInFunctionChecker: func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*VariableType, pos *Pos) {
			if len(e.TypedParameters) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(pos)))
			}
			meta := &BuildInFunctionPrintfMeta{}
			e.BuildInFunctionMeta = meta
			if len(args) == 0 || args[0] == nil {
				return // not error
			}
			if args[0].Type == VARIABLE_TYPE_OBJECT {
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
		IsBuildIn: true,
		Name:      common.BUILD_IN_FUNCTION_PRINT,
	}
	catchBuildFunction := &Function{}
	catchBuildFunction.IsBuildIn = true
	catchBuildFunction.Name = common.BUILD_IN_FUNCTION_CATCH
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_CATCH] = catchBuildFunction
	{
		catchBuildFunction.Type.ReturnList = make([]*VariableDefinition, 1)
		catchBuildFunction.Type.ReturnList[0] = &VariableDefinition{}
		catchBuildFunction.Type.ReturnList[0].Name = "returnValue"
		catchBuildFunction.Type.ReturnList[0].Type = &VariableType{}
		catchBuildFunction.Type.ReturnList[0].Type.Type = VARIABLE_TYPE_OBJECT
		catchBuildFunction.Type.ReturnList[0].Type.Class = &Class{}
		catchBuildFunction.Type.ReturnList[0].Type.Class.Name = DEFAULT_EXCEPTION_CLASS
		catchBuildFunction.Type.ReturnList[0].Type.Class.NotImportedYet = true
		//class is going to make value by checker
	}
	catchBuildFunction.buildInFunctionChecker = func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*VariableType, pos *Pos) {
		if len(e.TypedParameters) > 0 {
			*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
				errMsgPrefix(pos)))
		}
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
				c, err := ImportsLoader.LoadImport(DEFAULT_EXCEPTION_CLASS)
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s load exception class failed,err:%v",
						errMsgPrefix(pos), err))
					return
				}
				ft.Type.ReturnList[0].Type.Class = c.(*Class)
				err = block.InheritedAttribute.Defer.registerExceptionClass(c.(*Class))
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(pos), err))
				}
			} else {
				ft.Type.ReturnList[0].Type.Class = block.InheritedAttribute.Defer.ExceptionClass

			}
			return
		}
		if args[0] == nil {
			return
		}
		if args[0].Type != VARIABLE_TYPE_OBJECT {
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
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_PANIC] = &Function{
		buildInFunctionChecker: func(ft *Function, e *ExpressionFunctionCall,
			block *Block, errs *[]error, args []*VariableType, pos *Pos) {
			if len(e.TypedParameters) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(pos)))
			}
			if len(args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s buildin function 'panic' expect one argument",
					errMsgPrefix(pos)))
				return
			}
			if len(args) == 0 || args[0] == nil {
				return
			}
			if args[0].Type != VARIABLE_TYPE_OBJECT {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for panic",
					errMsgPrefix(pos), args[0].TypeString()))
				return
			}
			if have, _ := args[0].Class.haveSuper(JAVA_THROWABLE_CLASS); have == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for panic",
					errMsgPrefix(pos), args[0].TypeString()))
				return
			}
		},
		IsBuildIn: true,
		Name:      common.BUILD_IN_FUNCTION_PANIC,
	}
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_MONITORENTER] = &Function{
		buildInFunctionChecker: monitorChecker,
		IsBuildIn:              true,
		Name:                   common.BUILD_IN_FUNCTION_MONITORENTER,
	}
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_MONITOREXIT] = &Function{
		buildInFunctionChecker: monitorChecker,
		IsBuildIn:              true,
		Name:                   common.BUILD_IN_FUNCTION_MONITOREXIT,
	}
	// len
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_LEN] = &Function{
		buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*VariableType, pos *Pos) {
			if len(e.TypedParameters) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(pos)))
			}
			if len(args) != 1 {
				*errs = append(*errs, fmt.Errorf("%s expect one argument", errMsgPrefix(pos)))
				return
			}
			if args[0] == nil {
				return
			}
			if args[0].Type != VARIABLE_TYPE_ARRAY && args[0].Type != VARIABLE_TYPE_JAVA_ARRAY &&
				args[0].Type != VARIABLE_TYPE_MAP && args[0].Type != VARIABLE_TYPE_STRING {
				*errs = append(*errs, fmt.Errorf("%s len expect 'array' or 'map' or 'string' argument",
					errMsgPrefix(pos)))
				return
			}
		},
		IsBuildIn: true,
		Name:      common.BUILD_IN_FUNCTION_LEN,
	}
	lenFunction := buildInFunctionsMap[common.BUILD_IN_FUNCTION_LEN]
	lenFunction.Type.ReturnList = make(ReturnList, 1)
	lenFunction.Type.ReturnList[0] = &VariableDefinition{}
	lenFunction.Type.ReturnList[0].Type = &VariableType{}
	lenFunction.Type.ReturnList[0].Type.Type = VARIABLE_TYPE_INT
	// sprintf
	sprintfBuildFunction := &Function{}
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_SPRINTF] = sprintfBuildFunction
	sprintfBuildFunction.Name = common.BUILD_IN_FUNCTION_SPRINTF
	sprintfBuildFunction.IsBuildIn = true
	{
		sprintfBuildFunction.Type.ReturnList = make([]*VariableDefinition, 1)
		sprintfBuildFunction.Type.ReturnList[0] = &VariableDefinition{}
		sprintfBuildFunction.Type.ReturnList[0].Name = "retrunValue"
		sprintfBuildFunction.Type.ReturnList[0].Type = &VariableType{}
		sprintfBuildFunction.Type.ReturnList[0].Type.Type = VARIABLE_TYPE_STRING
	}
	sprintfBuildFunction.buildInFunctionChecker = func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
		args []*VariableType, pos *Pos) {
		if len(e.TypedParameters) > 0 {
			*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
				errMsgPrefix(pos)))
		}
		if len(args) == 0 {
			err := fmt.Errorf("%s '%s' expect one argument at lease",
				errMsgPrefix(pos), common.BUILD_IN_FUNCTION_SPRINTF)
			*errs = append(*errs, err)
			return
		}
		if args[0] == nil {
			return
		}
		if args[0].Type != VARIABLE_TYPE_STRING {
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
		meta := &BuildInFunctionSprintfMeta{}
		e.BuildInFunctionMeta = meta
		meta.Format = e.Args[0]
		meta.ArgsLength = len(args) - 1
		e.Args = e.Args[1:]
	}
	// printf
	buildInFunctionsMap[common.BUILD_IN_FUNCTION_PRINTF] = &Function{
		buildInFunctionChecker: func(ft *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
			args []*VariableType, pos *Pos) {
			if len(e.TypedParameters) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(pos)))
			}
			meta := &BuildInFunctionPrintfMeta{}
			e.BuildInFunctionMeta = meta
			if len(args) == 0 {
				err := fmt.Errorf("%s '%s' expect one argument at least",
					errMsgPrefix(pos), common.BUILD_IN_FUNCTION_PRINTF)
				*errs = append(*errs, err)
				return
			}
			if args[0] == nil {
				return
			}
			if args[0].Type == VARIABLE_TYPE_OBJECT {
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
			if args[0].Type != VARIABLE_TYPE_STRING {
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
		IsBuildIn: true,
		Name:      common.BUILD_IN_FUNCTION_PRINTF,
	}
}

func monitorChecker(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
	args []*VariableType, pos *Pos) {
	if len(e.TypedParameters) > 0 {
		*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
			errMsgPrefix(pos)))
	}
	if len(args) != 1 {
		*errs = append(*errs, fmt.Errorf("%s only expect one argument", errMsgPrefix(pos)))
		return
	}
	if args[0] == nil {
		return
	}
	if args[0].IsPointer() == false || args[0].Type == VARIABLE_TYPE_STRING {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not valid type to call",
			errMsgPrefix(pos), args[0].TypeString()))
		return
	}
}
