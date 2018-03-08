package jvm

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, method *cg.MethodHighLevel, f *ast.Function) {
	context := &Context{}
	context.function = f
	context.mainclass = m.mainclass
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
		fmt.Println("%v\n", method.Code.Codes)
	}()
	method.Code.MaxLocals = f.Varoffset
	m.buildFunctionParameterAndReturnList(class, &method.Code, f.Typ, context)
	m.buildBlock(class, &method.Code, f.Block, context)
	return
}

func (m *MakeClass) buildFunctionParameterAndReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode, ft *ast.FunctionType, context *Context) {
	for _, v := range ft.ReturnList {
		if v.BeenCaptured {
			panic("111111111")
		}
		maxstack, es := m.MakeExpression.build(class, code, v.Expression, context)
		backPatchEs(es, code.CodeLength)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
		if v.Typ.IsNumber() && v.Typ.Typ != v.Expression.VariableType.Typ {
			m.MakeExpression.numberTypeConverter(code, v.Expression.VariableType.Typ, v.Typ.Typ)
		}
		copyOP(code, storeSimpleVarOp(v.Typ.Typ, v.LocalValOffset)...)
	}
}
