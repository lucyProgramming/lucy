package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) mkBuildInSprintf(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	// format,must be string
	call := e.Data.(*ast.ExpressionFunctionCall)
	meta := call.BuildInFunctionMeta.(*ast.BuildInFunctionSprintfMeta)
	maxStack = buildExpression.build(class, code, meta.Format, context, state)
	state.pushStack(class, state.newObjectVariableType(javaStringClass))
	loadInt32(class, code, int32(meta.ArgsLength))
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst("java/lang/Object", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	currentStack := uint16(2)
	if currentStack > maxStack {
		maxStack = currentStack
	}
	objectArray := &ast.Type{}
	objectArray.Type = ast.VariableTypeJavaArray
	objectArray.Array = state.newObjectVariableType(javaRootClass)
	state.pushStack(class, objectArray)
	index := int32(0)
	for _, v := range call.Args {
		currentStack = 2
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, index)
		currentStack += 2
		state.pushStack(class, objectArray)
		state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
		stack := buildExpression.build(class, code, v, context, state)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Value.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.Value)
		}
		code.Codes[code.CodeLength] = cg.OP_aastore
		code.CodeLength++
		index++
		state.popStack(2)
	}
	code.Codes[code.CodeLength] = cg.OP_invokestatic
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringClass,
		Method:     "format",
		Descriptor: "(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if e.IsStatementExpression {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}
	return
}
