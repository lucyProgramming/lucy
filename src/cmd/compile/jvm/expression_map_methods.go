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
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		code.CodeLength++
		name := "containsKey"
		if call.Name != "keyExist" {
			name = "containsValue"
		}
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Name:       name,
			Descriptor: "(Ljava/lang/Object;)Z",
		}, code.Codes[code.CodeLength:code.CodeLength+2])
		code.CodeLength += 2
	case "remove":
		currentStack := uint16(1)
		removeMethodName := "remove"
		removeDescriptor := "(Ljava/lang/Object;)Ljava/lang/Object;"
		for k, v := range call.Args {
			currentStack = 1
			if v.IsCall() && len(v.VariableTypes) > 0 {
				stack, _ := m.build(class, code, v, context)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				m.buildStoreArrayListAutoVar(code, context) // store to temp
				for kk, _ := range v.VariableTypes {
					currentStack = 1
					if k == len(call.Args)-1 && kk == len(v.VariableTypes)-1 {
					} else {
						code.Codes[code.CodeLength] = cg.OP_dup
						code.CodeLength++
						currentStack++
					}
					//load
					m.buildLoadArrayListAutoVar(code, context)
					loadInt32(code, class, int32(kk))
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      java_arrylist_class,
						Name:       "get",
						Descriptor: "(I)Ljava/lang/Object;",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
					//remove
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      java_hashmap_class,
						Name:       removeMethodName,
						Descriptor: removeDescriptor,
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.Codes[code.CodeLength+3] = cg.OP_pop
					code.CodeLength += 4
				}
				continue
			}
			variableType := v.VariableType
			if v.IsCall() {
				variableType = v.VariableTypes[0]
			}
			if k != len(call.Args)-1 { // not last one
				code.Codes[code.CodeLength] = cg.OP_dup
				currentStack++
				if currentStack > maxstack {
					maxstack = currentStack
				}
			}
			if variableType.IsPointer() == false {
				currentStack += PrimitiveObjectConverter.prepareStack(class, code, variableType)
				if currentStack > maxstack {
					maxstack = currentStack
				}
			}
			stack, _ := m.build(class, code, v, context)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			if variableType.IsPointer() == false {
				PrimitiveObjectConverter.putInObject(class, code, variableType)
			}
			//call remove
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      java_hashmap_class,
				Name:       removeMethodName,
				Descriptor: removeDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = cg.OP_pop
			code.CodeLength += 4
		}
	case "removeAll":
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Name:       "clear",
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case "size":
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Name:       "size",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	default:
		panic("unkown" + call.Name)
	}
	return
}
