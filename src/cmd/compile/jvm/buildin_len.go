package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) mkBuildInLen(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	a0 := call.Args[0]
	maxStack = buildExpression.build(class, code, a0, context, state)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 2 > maxStack {
		maxStack = 2
	}
	code.Codes[code.CodeLength] = cg.OP_ifnonnull
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 8)
	code.Codes[code.CodeLength+3] = cg.OP_pop
	code.Codes[code.CodeLength+4] = cg.OP_iconst_0
	code.CodeLength += 5
	noNullExit := (&cg.Exit{}).Init(cg.OP_goto, code)
	state.pushStack(class, a0.Value)
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	if a0.Value.Type == ast.VariableTypeJavaArray {
		code.Codes[code.CodeLength] = cg.OP_arraylength
		code.CodeLength++
	} else if a0.Value.Type == ast.VariableTypeArray {
		meta := ArrayMetas[a0.Value.Array.Type]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     "size",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else if a0.Value.Type == ast.VariableTypeMap {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      mapClass,
			Method:     "size",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else if a0.Value.Type == ast.VariableTypeString {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "length",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	writeExits([]*cg.Exit{noNullExit}, code.CodeLength)
	state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	if e.IsStatementExpression {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}
	return
}
