package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildSlice(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	slice := e.Data.(*ast.ExpressionSlice)
	maxstack, _ = m.build(class, code, slice.Expression, context, nil)
	// build start
	if slice.Start != nil {
		stack, _ := m.build(class, code, slice.Start, context, nil)
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
		if slice.Start.VariableType.Typ == ast.VARIABLE_TYPE_LONG {
			m.numberTypeConverter(code, ast.VARIABLE_TYPE_LONG, ast.VARIABLE_TYPE_INT)
		}
	} else {
		code.Codes[code.CodeLength] = cg.OP_iconst_0
		code.CodeLength++
	}
	if slice.End != nil {
		stack, _ := m.build(class, code, slice.End, context, nil)
		if slice.End.VariableType.Typ == ast.VARIABLE_TYPE_LONG {
			m.numberTypeConverter(code, ast.VARIABLE_TYPE_LONG, ast.VARIABLE_TYPE_INT)
		}
		if t := 2 + stack; t > maxstack {
			maxstack = t
		}
	} else {
		code.Codes[code.CodeLength] = cg.OP_iconst_m1
		code.CodeLength++
		if 3 > maxstack { // stack top ->  ...arrayref start end
			maxstack = 3
		}
	}
	meta := ArrayMetas[e.VariableType.ArrayType.Typ]
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.classname,
		Method:     "slice",
		Descriptor: meta.sliceDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
