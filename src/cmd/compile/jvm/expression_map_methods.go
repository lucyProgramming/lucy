package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildMapMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	maxstack, _ = m.build(class, code, call.Expression, context)
	switch call.Name {
	case "keyExist", "valueExist":
		variableType := call.Args[0].VariableType
		if call.Args[0].IsCall() {
			variableType = call.Args[0].VariableTypes[0]
		}
		if variableType.IsPointer() == false {
			PrimitiveObjectConverter.prepareStack(class, code, variableType)
		}
		stack, _ := m.build(class, code, call.Args[0], context)
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
		if variableType.IsPointer() == false {
			if t := 3 + stack; t > maxstack {
				maxstack = t
			}
			PrimitiveObjectConverter.putInObject(class, code, variableType)
		}
		//
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		code.CodeLength++
		if call.Name == "keyExist" {
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      java_hashmap_class,
				Name:       "containsKey",
				Descriptor: "(Ljava/lang/Object;)Z",
			}, code.Codes[code.CodeLength:code.CodeLength+2])
		} else {
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      java_hashmap_class,
				Name:       "containsValue",
				Descriptor: "(Ljava/lang/Object;)Z",
			}, code.Codes[code.CodeLength:code.CodeLength+2])
		}
		code.CodeLength += 2
	case "remove":

	case "removeAll":
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Name:       "clear",
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	default:
		panic("unkown" + call.Name)
	}
	return
}
