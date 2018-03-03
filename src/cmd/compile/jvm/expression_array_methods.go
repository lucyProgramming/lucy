package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/common"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildArrayMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	switch call.Name {
	case common.ARRAY_METHOD_CAP,
		common.ARRAY_METHOD_SIZE,
		common.ARRAY_METHOD_START,
		common.ARRAY_METHOD_END:
		maxstack, _ = m.build(class, code, call.Expression, context)
		meta := ArrayMetas[call.Expression.VariableType.ArrayType.Typ]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.classname,
			Name:       call.Name,
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case common.ARRAY_METHOD_APPEND:
		maxstack, _ = m.build(class, code, call.Expression, context)
		meta := ArrayMetas[call.Expression.VariableType.ArrayType.Typ]
		for k, v := range call.Args {
			currentStack := uint16(1)
			if v.IsCall() && len(v.VariableTypes) > 0 {
				stack, _ := m.build(class, code, v, context)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				m.buildStoreArrayListAutoVar(code, context)
				for kk, t := range v.VariableTypes {
					currentStack := uint16(1)
					if k == len(call.Args)-1 && kk == len(v.VariableTypes)-1 {

					} else {
						code.Codes[code.CodeLength] = cg.OP_dup
						code.CodeLength++
						currentStack++
						if currentStack > maxstack {
							maxstack = currentStack
						}
					}
					if t := m.unPackArraylist(class, code, kk, t, context) + currentStack; t > maxstack {
						maxstack = t
					}
					if t := currentStack + t.JvmSlotSize(); t > maxstack {
						maxstack = t
					}
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      meta.classname,
						Name:       "append",
						Descriptor: meta.appendDescriptor,
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				}
				continue
			}
			if k != len(call.Args)-1 {
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				currentStack++
				if currentStack > maxstack {
					maxstack = currentStack
				}
			}
			stack, es := m.build(class, code, v, context)
			backPatchEs(es, code.CodeLength)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.classname,
				Name:       "append",
				Descriptor: meta.appendDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	default:
		panic("unkown method:" + call.Name)
	}
	return
}
