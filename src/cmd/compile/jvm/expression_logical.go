package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16, exits []*cg.Exit) {
	bin := e.Data.(*ast.ExpressionBinary)
	maxStack, es := buildExpression.build(class, code, bin.Left, context, state)
	if es != nil {
		state.pushStack(class, bin.Left.ExpressionValue)
		context.MakeStackMap(code, state, code.CodeLength)
		state.popStack(1)
		fillOffsetForExits(es, code.CodeLength)
	}
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 2 > maxStack { // dup increment stack
		maxStack = 2
	}
	if e.Type == ast.ExpressionTypeLogicalOr {
		// at this point,value is clear,leave 1 on stack
		exits = append(exits, (&cg.Exit{}).FromCode(cg.OP_ifne, code))
	} else { //  &&
		// at this point,value is clear,leave 0 on stack
		exits = append(exits, (&cg.Exit{}).FromCode(cg.OP_ifeq, code))
	}
	code.Codes[code.CodeLength] = cg.OP_pop // pop 0
	code.CodeLength++
	stack, es := buildExpression.build(class, code, bin.Right, context, state)
	if es != nil {
		exits = append(exits, es...)
	}
	if stack > maxStack {
		maxStack = stack
	}
	return
}
