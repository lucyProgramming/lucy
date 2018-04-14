package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16, exits []*cg.JumpBackPatch) {
	exits = []*cg.JumpBackPatch{}
	bin := e.Data.(*ast.ExpressionBinary)
	maxstack, es := m.build(class, code, bin.Left, context, state)
	if es != nil {
		exits = append(exits, es...)
	}
	if e.Typ == ast.EXPRESSION_TYPE_LOGICAL_OR {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		if 2 > maxstack { // dup increment stack
			maxstack = 2
		}
		exit := (&cg.JumpBackPatch{}).FromCode(cg.OP_ifne, code) // at this point,value is clear,leave 1 on stack
		exits = append(exits, exit)
		code.Codes[code.CodeLength] = cg.OP_pop // pop 0
		code.CodeLength++
		stack, es := m.build(class, code, bin.Right, context, state)
		if es != nil {
			exits = append(exits, es...)
		}
		if stack > maxstack {
			maxstack = stack
		}
	} else { //and
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		if 2 > maxstack { // dup increment stack
			maxstack = 2
		}
		exit := (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code) // at this point,value is clear,leave 1 on stack
		exits = append(exits, exit)
		code.Codes[code.CodeLength] = cg.OP_pop // pop 1
		code.CodeLength++
		stack, es := m.build(class, code, bin.Right, context, state)
		if es != nil {
			exits = append(exits, es...)
		}
		if stack > maxstack {
			maxstack = stack
		}
	}
	return
}
