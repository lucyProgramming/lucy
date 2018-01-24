package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildFunctionCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Func.Isbuildin {
		return m.mkBuildinFunctionCall(class, code, call, context)
	}
	return
}

func (m *MakeExpression) buildCallArgs(class *cg.ClassHighLevel, code *cg.AttributeCode, args []*ast.Expression, context *Context) (maxstack uint16) {
	currentStack := uint16(0)
	for _, e := range args {
		var variabletype *ast.VariableType
		if e.Typ == ast.EXPRESSION_TYPE_METHOD_CALL || e.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL {
			if len(e.VariableTypes) > 1 {
				m.buildStoreArrayListAutoVar(class, code, context)
				if t := currentStack + 1; t > maxstack {
					maxstack = t
				}
				for k, t := range e.VariableTypes {
					m.buildLoadArrayListAutoVar(class, code, context)
					switch k {
					case 0:
						code.Codes[code.CodeLength] = cg.OP_iconst_0
						code.CodeLength++
					case 1:
						code.Codes[code.CodeLength] = cg.OP_iconst_1
						code.CodeLength++
					case 2:
						code.Codes[code.CodeLength] = cg.OP_iconst_2
						code.CodeLength++
					case 3:
						code.Codes[code.CodeLength] = cg.OP_iconst_3
						code.CodeLength++
					case 4:
						code.Codes[code.CodeLength] = cg.OP_iconst_4
						code.CodeLength++
					case 5:
						code.Codes[code.CodeLength] = cg.OP_iconst_5
						code.CodeLength++
					default:
						if k > 255 {
							panic("over 255")
						}
						code.Codes[code.CodeLength] = cg.OP_bipush
						code.Codes[code.CodeLength+1] = byte(k)
						code.CodeLength += 2
					}
					if t := currentStack + 2; t > maxstack {
						maxstack = t
					}
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
						Class:      "java/util/ArrayList",
						Name:       "get",
						Descriptor: "(I)java/lang/Object",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
					switch t.Typ {
					case ast.VARIABLE_TYPE_BOOL:
						fallthrough
					case ast.VARIABLE_TYPE_BYTE:
						fallthrough
					case ast.VARIABLE_TYPE_SHORT:
						fallthrough
					case ast.VARIABLE_TYPE_INT:
						code.Codes[code.CodeLength] = cg.OP_invokevirtual
						class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
							Class:      "java/lang/Integer",
							Name:       "intValue",
							Descriptor: "()I",
						}, code.Codes[code.CodeLength+1:code.CodeLength+3])
						code.CodeLength += 3
					case ast.VARIABLE_TYPE_LONG:
						code.Codes[code.CodeLength] = cg.OP_invokevirtual
						class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
							Class:      "java/lang/Long",
							Name:       "longValue",
							Descriptor: "()J",
						}, code.Codes[code.CodeLength+1:code.CodeLength+3])
						code.CodeLength += 3
					case ast.VARIABLE_TYPE_FLOAT:
						code.Codes[code.CodeLength] = cg.OP_invokevirtual
						class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
							Class:      "java/lang/Float",
							Name:       "floatValue",
							Descriptor: "()F",
						}, code.Codes[code.CodeLength+1:code.CodeLength+3])
						code.CodeLength += 3
					case ast.VARIABLE_TYPE_DOUBLE:
						code.Codes[code.CodeLength] = cg.OP_invokevirtual
						class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
							Class:      "java/lang/Double",
							Name:       "doubleValue",
							Descriptor: "()D",
						}, code.Codes[code.CodeLength+1:code.CodeLength+3])
						code.CodeLength += 3
					case ast.VARIABLE_TYPE_STRING:
						fallthrough
					case ast.VARIABLE_TYPE_OBJECT:
						fallthrough
					case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
						//nothing to do
					}
					currentStack += t.JvmSlotSize()
				}
				continue
			}
			variabletype = e.VariableTypes[0]
		} else {
			variabletype = e.VariableType
		}
		ms, es := m.build(class, code, e, context)
		backPatchEs(es, code)
		if t := ms + currentStack; t > maxstack {
			maxstack = t
		}
		currentStack += variabletype.JvmSlotSize()
	}
	return
}

func (m *MakeExpression) buildMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	return
}
