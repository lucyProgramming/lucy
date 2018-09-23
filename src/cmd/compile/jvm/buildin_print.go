package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	function print
*/
func (buildExpression *BuildExpression) mkBuildInPrint(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	// get stream from stdout
	code.Codes[code.CodeLength] = cg.OP_getstatic
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      "java/lang/System",
		Field:      "out",
		Descriptor: "Ljava/io/PrintStream;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxStack = 1
	if len(call.Args) == 0 {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/io/PrintStream",
			Method:     "println",
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	length := len(state.Stacks)
	defer func() {
		// print have no return value,stack is empty
		state.popStack(len(state.Stacks) - length)
	}()
	state.pushStack(class, state.newObjectVariableType(javaPrintStreamClass))
	if len(call.Args) == 1 {
		stack := buildExpression.build(class, code, call.Args[0], context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}

		switch call.Args[0].Value.Type {
		case ast.VariableTypeBool:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Z)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeChar:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(C)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(I)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(J)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(F)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(D)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeString:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Ljava/lang/String;)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		default:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Ljava/lang/Object;)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}

		return
	}
	code.Codes[code.CodeLength] = cg.OP_ldc_w
	class.InsertStringConst("", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	state.pushStack(class, state.newObjectVariableType(javaStringClass))
	defer state.popStack(1)
	for k, v := range call.Args {
		variableType := v.Value
		stack := buildExpression.build(class, code, v, context, state)
		if t := 2 + stack; t > maxStack {
			maxStack = t
		}
		if t := 2 + buildExpression.stackTop2String(class, code, variableType, context, state); t > maxStack {
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "concat",
			Descriptor: "(Ljava/lang/String;)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if k != len(call.Args)-1 {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertStringConst(" ", code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			if 2 > maxStack {
				maxStack = 2
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      javaStringClass,
				Method:     "concat",
				Descriptor: "(Ljava/lang/String;)Ljava/lang/String;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
	// call println
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/io/PrintStream",
		Method:     "println",
		Descriptor: "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
