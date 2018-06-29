package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) mkBuildInLen(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	maxStack, _ = buildExpression.build(class, code, call.Args[0], context, state)
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 2 > maxStack {
		maxStack = 2
	}
	exit := (&cg.Exit{}).FromCode(cg.OP_ifnull, code)
	//binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 3)
	if call.Args[0].ExpressionValue.Type == ast.VariableTypeJavaArray {
		code.Codes[code.CodeLength] = cg.OP_arraylength
		code.CodeLength++
	} else if call.Args[0].ExpressionValue.Type == ast.VariableTypeArray {
		meta := ArrayMetas[call.Args[0].ExpressionValue.Array.Type]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     "size",
			Descriptor: "()I",
		},
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else if call.Args[0].ExpressionValue.Type == ast.VariableTypeMap {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaMapClass,
			Method:     "size",
			Descriptor: "()I",
		},
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else if call.Args[0].ExpressionValue.Type == ast.VariableTypeString {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "length",
			Descriptor: "()I",
		},
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	fillOffsetForExits([]*cg.Exit{exit}, code.CodeLength+3)
	state.pushStack(class, call.Args[0].ExpressionValue)
	context.MakeStackMap(code, state, code.CodeLength+3)
	state.popStack(1)
	code.Codes[code.CodeLength] = cg.OP_goto
	binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 5)
	code.Codes[code.CodeLength+3] = cg.OP_pop
	code.Codes[code.CodeLength+4] = cg.OP_iconst_0
	code.CodeLength += 5
	state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
	context.MakeStackMap(code, state, code.CodeLength)
	state.popStack(1)
	if e.IsStatementExpression {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}
	return
}
