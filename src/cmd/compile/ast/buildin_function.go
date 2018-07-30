package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func init() {
	registerBuildInFunctions()
}

func registerBuildInFunctions() {
	buildInFunctionsMap[common.BuildInFunctionPrint] = &Function{
		buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Pos) {
			if len(e.ParameterTypes) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(pos)))
			}
			meta := &BuildInFunctionPrintfMeta{}
			e.BuildInFunctionMeta = meta
			if len(args) == 0 || args[0] == nil {
				return // not error
			}
			//if args[0].Type == VariableTypeObject {
			//	have, _ := args[0].Class.haveSuper("java/io/PrintStream")
			//	if have {
			//		_, err := e.Args[0].mustBeOneValueContext(e.Args[0].MultiValues)
			//		if err != nil {
			//			*errs = append(*errs, err)
			//		} else {
			//			meta.Stream = e.Args[0]
			//			e.Args = e.Args[1:]
			//		}
			//	}
			//}
		},
		IsBuildIn: true,
		Name:      common.BuildInFunctionPrint,
	}
	catch := &Function{}
	catch.IsBuildIn = true
	catch.Name = common.BuildInFunctionCatch
	buildInFunctionsMap[common.BuildInFunctionCatch] = catch
	{
		catch.Type.ReturnList = make([]*Variable, 1)
		catch.Type.ReturnList[0] = &Variable{}
		catch.Type.ReturnList[0].Name = "returnValue"
		catch.Type.ReturnList[0].Type = &Type{}
		catch.Type.ReturnList[0].Type.Type = VariableTypeObject
		catch.Type.ReturnList[0].Type.Class = &Class{}
		catch.Type.ReturnList[0].Type.Class.Name = DefaultExceptionClass
		catch.Type.ReturnList[0].Type.Class.NotImportedYet = true
	}
	catch.buildInFunctionChecker = func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Pos) {
		if len(e.ParameterTypes) > 0 {
			*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
				errMsgPrefix(pos)))
		}
		if block.InheritedAttribute.Defer == nil {
			*errs = append(*errs, fmt.Errorf("%s buildin function '%s' only allow in defer block",
				errMsgPrefix(pos), common.BuildInFunctionCatch))
			return
		}
		if len(args) > 1 {
			*errs = append(*errs, fmt.Errorf("%s build function '%s' expect at most 1 argument",
				errMsgPrefix(pos), common.BuildInFunctionCatch))
			return
		}
		if len(args) == 0 {
			// make default exception class
			// load java/lang/Exception this is default exception level to catch
			if block.InheritedAttribute.Defer.ExceptionClass == nil {
				c, err := ImportsLoader.LoadImport(DefaultExceptionClass)
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s load exception class failed,err:%v",
						errMsgPrefix(pos), err))
					return
				}
				f.Type.ReturnList[0].Type.Class = c.(*Class)
				err = block.InheritedAttribute.Defer.registerExceptionClass(c.(*Class))
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(pos), err))
				}
			} else {
				f.Type.ReturnList[0].Type.Class = block.InheritedAttribute.Defer.ExceptionClass

			}
			return
		}
		if args[0] == nil {
			return
		}
		if args[0].Type != VariableTypeObject {
			*errs = append(*errs, fmt.Errorf("%s build function '%s' expect a object ref argument",
				errMsgPrefix(pos), common.BuildInFunctionCatch))
			return
		}
		if has, _ := args[0].Class.haveSuper(JavaThrowableClass); has == false {
			*errs = append(*errs, fmt.Errorf("%s '%s' does not have super-class '%s'",
				errMsgPrefix(pos), args[0].Class.Name, JavaThrowableClass))
			return
		}
		err := block.InheritedAttribute.Defer.registerExceptionClass(args[0].Class)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(pos), err))
		}
	}
	buildInFunctionsMap[common.BuildInFunctionPanic] = &Function{
		buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall,
			block *Block, errs *[]error, args []*Type, pos *Pos) {
			if len(e.ParameterTypes) > 0 {
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
			if args[0].Type != VariableTypeObject {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for panic",
					errMsgPrefix(pos), args[0].TypeString()))
				return
			}
			if have, _ := args[0].Class.haveSuper(JavaThrowableClass); have == false {
				*errs = append(*errs, fmt.Errorf("%s cannot use '%s' for panic",
					errMsgPrefix(pos), args[0].TypeString()))
				return
			}
		},
		IsBuildIn: true,
		Name:      common.BuildInFunctionPanic,
	}
	buildInFunctionsMap[common.BuildInFunctionMonitorEnter] = &Function{
		buildInFunctionChecker: monitorChecker,
		IsBuildIn:              true,
		Name:                   common.BuildInFunctionMonitorEnter,
	}
	buildInFunctionsMap[common.BuildInFunctionMonitorExit] = &Function{
		buildInFunctionChecker: monitorChecker,
		IsBuildIn:              true,
		Name:                   common.BuildInFunctionMonitorExit,
	}
	// len
	buildInFunctionsMap[common.BuildInFunctionLen] = &Function{
		buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Pos) {
			if len(e.ParameterTypes) > 0 {
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
			if args[0].Type != VariableTypeArray && args[0].Type != VariableTypeJavaArray &&
				args[0].Type != VariableTypeMap && args[0].Type != VariableTypeString {
				*errs = append(*errs, fmt.Errorf("%s len expect 'array' or 'map' or 'string' argument",
					errMsgPrefix(pos)))
				return
			}
		},
		IsBuildIn: true,
		Name:      common.BuildInFunctionLen,
	}
	Len := buildInFunctionsMap[common.BuildInFunctionLen]
	Len.Type.ReturnList = make(ReturnList, 1)
	Len.Type.ReturnList[0] = &Variable{}
	Len.Type.ReturnList[0].Type = &Type{}
	Len.Type.ReturnList[0].Type.Type = VariableTypeInt
	// sprintf
	sprintf := &Function{}
	buildInFunctionsMap[common.BuildInFunctionSprintf] = sprintf
	sprintf.Name = common.BuildInFunctionSprintf
	sprintf.IsBuildIn = true
	{
		sprintf.Type.ReturnList = make([]*Variable, 1)
		sprintf.Type.ReturnList[0] = &Variable{}
		sprintf.Type.ReturnList[0].Name = "returnValue"
		sprintf.Type.ReturnList[0].Type = &Type{}
		sprintf.Type.ReturnList[0].Type.Type = VariableTypeString
	}
	sprintf.buildInFunctionChecker = func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
		args []*Type, pos *Pos) {
		if len(e.ParameterTypes) > 0 {
			*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
				errMsgPrefix(pos)))
		}
		if len(args) == 0 {
			err := fmt.Errorf("%s '%s' expect one argument at lease",
				errMsgPrefix(pos), common.BuildInFunctionSprintf)
			*errs = append(*errs, err)
			return
		}
		if args[0] == nil {
			return
		}
		if args[0].Type != VariableTypeString {
			err := fmt.Errorf("%s '%s' first argument must be string",
				errMsgPrefix(pos), common.BuildInFunctionSprintf)
			*errs = append(*errs, err)
			return
		}
		_, err := e.Args[0].mustBeOneValueContext(e.Args[0].MultiValues)
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
	buildInFunctionsMap[common.BuildInFunctionPrintf] = &Function{
		buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
			args []*Type, pos *Pos) {
			if len(e.ParameterTypes) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(pos)))
			}
			meta := &BuildInFunctionPrintfMeta{}
			e.BuildInFunctionMeta = meta
			if len(args) == 0 {
				err := fmt.Errorf("%s '%s' expect one argument at least",
					errMsgPrefix(pos), common.BuildInFunctionPrintf)
				*errs = append(*errs, err)
				return
			}
			if args[0] == nil {
				return
			}
			//if args[0].Type == VariableTypeObject {
			//	have, _ := args[0].Class.haveSuper("java/io/PrintStream")
			//	if have {
			//		_, err := e.Args[0].mustBeOneValueContext(e.Args[0].MultiValues)
			//		if err != nil {
			//			*errs = append(*errs, err)
			//			return
			//		} else {
			//			//meta.Stream = e.Args[0]
			//			e.Args = e.Args[1:]
			//			args = args[1:]
			//		}
			//	}
			//}
			if len(args) == 0 {
				err := fmt.Errorf("%s missing format argument",
					errMsgPrefix(pos))
				*errs = append(*errs, err)
				return
			}
			if args[0] == nil {
				return
			}
			if args[0].Type != VariableTypeString {
				err := fmt.Errorf("%s format must be string",
					errMsgPrefix(pos))
				*errs = append(*errs, err)
				return
			}
			_, err := e.Args[0].mustBeOneValueContext(e.Args[0].MultiValues)
			if err != nil {
				*errs = append(*errs, err)
				return
			}
			meta.Format = e.Args[0]
			e.Args = e.Args[1:]
			meta.ArgsLength = len(args)
		},
		IsBuildIn: true,
		Name:      common.BuildInFunctionPrintf,
	}
	buildInFunctionsMap[common.BuildInFunctionBlockHole] = &Function{
		buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
			args []*Type, pos *Pos) {
			// nothing to check
		},
		IsBuildIn: true,
		Name:      common.BuildInFunctionBlockHole,
	}
}

func monitorChecker(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
	args []*Type, pos *Pos) {
	if len(e.ParameterTypes) > 0 {
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
	if args[0].IsPointer() == false {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not valid type to call",
			errMsgPrefix(pos), args[0].TypeString()))
		return
	}
}
