package jvm

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context, state *StackMapState) {
	if len(b.Defers) > 0 { // should be more defers when compile
		context.Defers = append(context.Defers, b.Defers...)
	}
	for _, s := range b.Statements {
		maxstack := m.buildStatement(class, code, b, s, context, state)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		if len(state.Stacks) > 0 {
			panic(fmt.Sprintf("stack is not empty:%d", len(state.Stacks)))
		}
	}

	if b.IsFunctionTopBlock == false && len(b.Defers) > 0 {
		index := len(b.Defers) - 1
		for index >= 0 {
			ss := (&StackMapState{}).FromLast(state)
			m.buildBlock(class, code, &b.Defers[index].Block, context, state)
			index--
			state.addTop(ss)
		}
	}
	return
}
