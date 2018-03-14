package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

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

func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, method *cg.MethodHighLevel, f *ast.Function) {
	context := &Context{}
	context.function = f
	context.mainclass = m.mainclass
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	code := &method.Code
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
		method.Code.MaxLocals = f.Varoffset // could  new slot when compile
	}()
	m.buildFunctionParameterAndReturnList(class, &method.Code, f.Typ, context)
	if f.AutoVarForReturnBecauseOfDefer != nil && f.HaveNoReturnValue() == false {
		if len(f.Typ.ReturnList) == 1 {
			switch f.Typ.ReturnList[0].Typ.Typ {
			case ast.VARIABLE_TYPE_BOOL:
				fallthrough
			case ast.VARIABLE_TYPE_BYTE:
				fallthrough
			case ast.VARIABLE_TYPE_SHORT:
				fallthrough
			case ast.VARIABLE_TYPE_INT:
				maxstack := loadLocalVar(class, code, f.Typ.ReturnList[0])
				if maxstack > code.MaxStack {
					code.MaxStack = maxstack
				}
				copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, f.AutoVarForReturnBecauseOfDefer.Offset)...)
			case ast.VARIABLE_TYPE_LONG:
				maxstack := loadLocalVar(class, code, f.Typ.ReturnList[0])
				if maxstack > code.MaxStack {
					code.MaxStack = maxstack
				}
				copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_LONG, f.AutoVarForReturnBecauseOfDefer.Offset)...)
			case ast.VARIABLE_TYPE_FLOAT:
				maxstack := loadLocalVar(class, code, f.Typ.ReturnList[0])
				if maxstack > code.MaxStack {
					code.MaxStack = maxstack
				}
				copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_LONG, f.AutoVarForReturnBecauseOfDefer.Offset)...)
			case ast.VARIABLE_TYPE_DOUBLE:
				maxstack := loadLocalVar(class, code, f.Typ.ReturnList[0])
				if maxstack > code.MaxStack {
					code.MaxStack = maxstack
				}
				copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_DOUBLE, f.AutoVarForReturnBecauseOfDefer.Offset)...)
			case ast.VARIABLE_TYPE_STRING:
				maxstack := loadLocalVar(class, code, f.Typ.ReturnList[0])
				if maxstack > code.MaxStack {
					code.MaxStack = maxstack
				}
				copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForReturnBecauseOfDefer.Offset)...)
			case ast.VARIABLE_TYPE_OBJECT:
				fallthrough
			case ast.VARIABLE_TYPE_MAP:
				fallthrough
			case ast.VARIABLE_TYPE_ARRAY: //[]int
				maxstack := loadLocalVar(class, code, f.Typ.ReturnList[0])
				if maxstack > code.MaxStack {
					code.MaxStack = maxstack
				}
				code.CodeLength++
				copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_LONG, f.AutoVarForReturnBecauseOfDefer.Offset)...)
			}
		} else { // >1
			code.Codes[code.CodeLength] = cg.OP_aconst_null
			code.CodeLength++
			copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForReturnBecauseOfDefer.Offset)...)
		}
	}

	context.firstCodeShouldUnderRecover = -1
	m.buildBlock(class, &method.Code, f.Block, context)
	return
}
