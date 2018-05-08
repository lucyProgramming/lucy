package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildLogical(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16, exits []*cg.JumpBackPatch) {
	exits = []*cg.JumpBackPatch{}
	bin := e.Data.(*ast.ExpressionBinary)
	maxstack, es := m.build(class, code, bin.Left, context, state)
	if es != nil {
		exits = append(exits, es...)
	}
	code.Codes[code.CodeLength] = cg.OP_dup
	code.CodeLength++
	if 2 > maxstack { // dup increment stack
		maxstack = 2
	}
	if e.Typ == ast.EXPRESSION_TYPE_LOGICAL_OR {
		// at this point,value is clear,leave 1 on stack
		exits = append(exits, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifne, code))
	} else { // and
		// at this point,value is clear,leave 0 on stack
		exits = append(exits, (&cg.JumpBackPatch{}).FromCode(cg.OP_ifeq, code))
	}
	code.Codes[code.CodeLength] = cg.OP_pop // pop 0
	code.CodeLength++
	stack, es := m.build(class, code, bin.Right, context, state)
	if es != nil {
		exits = append(exits, es...)
	}
	if stack > maxstack {
		maxstack = stack
	}

	return
}
