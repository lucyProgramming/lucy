package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

func init() {
	registerBuildInFunctions()
}

func registerBuildInFunctions() {
	{
		//print
		buildInFunctionsMap[common.BuildInFunctionPrint] = &Function{
			buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Pos) {
				if len(e.ParameterTypes) > 0 {
					*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
						errMsgPrefix(e.ParameterTypes[0].Pos)))
				}
			},
			IsBuildIn: true,
		}
	}
	{
		//catch
		catch := &Function{}
		catch.IsBuildIn = true
		buildInFunctionsMap[common.BuildInFunctionCatch] = catch
		catch.Type.ReturnList = make([]*Variable, 1)
		catch.Type.ReturnList[0] = &Variable{}
		catch.Type.ReturnList[0].Name = "returnValue"
		catch.Type.ReturnList[0].Type = &Type{}
		catch.Type.ReturnList[0].Type.Type = VariableTypeObject
		catch.Type.ReturnList[0].Type.Class = &Class{}
		catch.Type.ReturnList[0].Type.Class.Name = DefaultExceptionClass
		catch.Type.ReturnList[0].Type.Class.NotImportedYet = true
		catch.buildInFunctionChecker = func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Pos) {
			if len(e.ParameterTypes) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(e.ParameterTypes[0].Pos)))
			}
			if block.InheritedAttribute.Defer == nil {
				*errs = append(*errs, fmt.Errorf("%s buildin function '%s' only allow in defer block",
					errMsgPrefix(pos), common.BuildInFunctionCatch))
				return
			}
			if len(e.Args) > 1 {
				*errs = append(*errs, fmt.Errorf("%s build function '%s' expect at most 1 argument",
					errMsgPrefix(e.Args[1].Pos), common.BuildInFunctionCatch))
				return
			}
			if len(e.Args) == 0 {
				if block.InheritedAttribute.Defer.ExceptionClass == nil {
					c, err := PackageBeenCompile.loadClass(DefaultExceptionClass)
					if err != nil {
						*errs = append(*errs, fmt.Errorf("%s load exception class failed,err:%v",
							errMsgPrefix(pos), err))
						return
					}
					f.Type.ReturnList[0].Type.Class = c
					err = block.InheritedAttribute.Defer.registerExceptionClass(c)
					if err != nil {
						*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(pos), err))
					}
				} else {
					f.Type.ReturnList[0].Type.Class = block.InheritedAttribute.Defer.ExceptionClass
				}
				return
			}
			if e.Args[0].Type != ExpressionTypeString {
				*errs = append(*errs, fmt.Errorf("%s build function '%s' expect string argument",
					errMsgPrefix(args[0].Pos), common.BuildInFunctionCatch))
			}
			className := e.Args[0].Data.(string)
			c, err := PackageBeenCompile.loadClass(className)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Args[0].Pos), err))
				return
			}
			have, _ := c.haveSuperClass(e.Args[0].Pos, JavaThrowableClass)
			if have == false {
				*errs = append(*errs, fmt.Errorf("%s '%s' does't have super class '%s'",
					errMsgPrefix(e.Args[0].Pos), className, JavaThrowableClass))
				return
			}
			err = block.InheritedAttribute.Defer.registerExceptionClass(c)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(args[0].Pos), err))
			}
		}
	}

	{
		//panic
		buildInFunctionsMap[common.BuildInFunctionPanic] = &Function{
			buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall,
				block *Block, errs *[]error, args []*Type, pos *Pos) {
				if len(e.ParameterTypes) > 0 {
					*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
						errMsgPrefix(e.ParameterTypes[0].Pos)))
				}
				if len(args) != 1 {
					*errs = append(*errs, fmt.Errorf("%s buildin function 'panic' expect one argument",
						errMsgPrefix(pos)))
					return
				}
				if len(args) == 0 || args[0] == nil {
					return
				}
				meta := &BuildInFunctionPanicMeta{}
				e.BuildInFunctionMeta = meta
				if args[0].Type == VariableTypeObject {
					meta.ArgThrowable, _ = args[0].Class.haveSuperClass(args[0].Pos, JavaThrowableClass)
				}
			},
			IsBuildIn: true,
		}
	}
	{
		buildInFunctionsMap[common.BuildInFunctionMonitorEnter] = &Function{
			buildInFunctionChecker: monitorChecker,
			IsBuildIn:              true,
		}
		buildInFunctionsMap[common.BuildInFunctionMonitorExit] = &Function{
			buildInFunctionChecker: monitorChecker,
			IsBuildIn:              true,
		}
	}
	{
		// len
		buildInFunctionsMap[common.BuildInFunctionLen] = &Function{
			buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Pos) {
				if len(e.ParameterTypes) > 0 {
					*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
						errMsgPrefix(e.ParameterTypes[0].Pos)))
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
		}
		Len := buildInFunctionsMap[common.BuildInFunctionLen]
		Len.Type.ReturnList = make(ReturnList, 1)
		Len.Type.ReturnList[0] = &Variable{}
		Len.Type.ReturnList[0].Type = &Type{}
		Len.Type.ReturnList[0].Type.Type = VariableTypeInt
	}
	{
		// sprintf
		sprintf := &Function{}
		buildInFunctionsMap[common.BuildInFunctionSprintf] = sprintf
		sprintf.IsBuildIn = true
		sprintf.Type.ReturnList = make([]*Variable, 1)
		sprintf.Type.ReturnList[0] = &Variable{}
		sprintf.Type.ReturnList[0].Name = "returnValue"
		sprintf.Type.ReturnList[0].Type = &Type{}
		sprintf.Type.ReturnList[0].Type.Type = VariableTypeString
		sprintf.buildInFunctionChecker = func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
			args []*Type, pos *Pos) {
			if len(e.ParameterTypes) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(e.ParameterTypes[0].Pos)))
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
	}

	{
		// assert
		typeOf := &Function{}
		buildInFunctionsMap[common.BuildInFunctionAssert] = typeOf
		typeOf.IsBuildIn = true
		typeOf.buildInFunctionChecker = func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
			args []*Type, pos *Pos) {
			if len(e.ParameterTypes) > 0 {
				*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
					errMsgPrefix(e.ParameterTypes[0].Pos)))
			}
			if len(args) == 0 {
				err := fmt.Errorf("%s '%s' expect one argument at lease",
					errMsgPrefix(pos), typeOf.Name)
				*errs = append(*errs, err)
				return
			}
			for _, a := range args {
				if a == nil {
					continue
				}
				if a.Type != VariableTypeBool {
					err := fmt.Errorf("%s not a bool expression",
						errMsgPrefix(pos))
					*errs = append(*errs, err)
				}
			}
		}
	}
	{
		// printf
		buildInFunctionsMap[common.BuildInFunctionPrintf] = &Function{
			buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
				args []*Type, pos *Pos) {
				if len(e.ParameterTypes) > 0 {
					*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
						errMsgPrefix(e.ParameterTypes[0].Pos)))
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
		}
	}

	buildInFunctionsMap[common.BuildInFunctionBlockHole] = &Function{
		buildInFunctionChecker: func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
			args []*Type, pos *Pos) {
		},
		IsBuildIn: true,
	}
	for k, v := range buildInFunctionsMap {
		v.Name = k
	}
}

func monitorChecker(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
	args []*Type, pos *Pos) {
	if len(e.ParameterTypes) > 0 {
		*errs = append(*errs, fmt.Errorf("%s buildin function expect no typed parameter",
			errMsgPrefix(e.ParameterTypes[0].Pos)))
	}
	if len(args) != 1 {
		pos := pos
		getFirstPosFromArgs(args[1:], &pos)
		*errs = append(*errs, fmt.Errorf("%s only expect one argument", errMsgPrefix(pos)))
		return
	}
	if args[0] == nil {
		return
	}
	if args[0].IsPointer() == false {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not valid type to call",
			errMsgPrefix(args[0].Pos), args[0].TypeString()))
		return
	}
}
