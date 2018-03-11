package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

/*
maxstack means return value stack size
*/
func (m *MakeExpression) buildFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Func.Isbuildin {
		return m.mkBuildinFunctionCall(class, code, e, context)
	}
	if call.Func.IsClosureFunction() == false {
		maxstack = m.buildCallArgs(class, code, call.Args, call.Func.Typ.ParameterList, context)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      call.Func.ClassMethod.Class.Name,
			Name:       call.Func.Name,
			Descriptor: Descriptor.methodDescriptor(call.Func),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	} else {
		//closure function call
	}
	if e.IsStatementExpression {
		if e.CallHasReturnValue() == false { // nothing to do
		} else if len(e.VariableTypes) == 1 {
			if 2 == e.VariableTypes[0].JvmSlotSize() {
				code.Codes[code.CodeLength] = cg.OP_pop2
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop
			}
			code.CodeLength++
		} else { // > 1
			code.Codes[code.CodeLength] = cg.OP_pop // arraylist on top
			code.CodeLength++
		}
	}

	if e.CallHasReturnValue() == false { // nothing

	} else if len(e.VariableTypes) == 1 {
		if t := e.VariableTypes[0].JvmSlotSize(); t > maxstack {
			maxstack = t
		}
	} else { // > 1
		if 1 > maxstack {
			maxstack = 1
		}
	}
	return
}
