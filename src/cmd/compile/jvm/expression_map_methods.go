package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildMapMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	maxstack, _ = m.build(class, code, call.Expression, context, state)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	hashMapVerifyType := state.newObjectVariableType(java_hashmap_class)
	state.Stacks = append(state.Stacks,
		state.newStackMapVerificationTypeInfo(class, hashMapVerifyType)...)
	switch call.Name {
	case common.MAP_METHOD_KEY_EXISTS:
		variableType := call.Args[0].VariableType
		if call.Args[0].MayHaveMultiValue() {
			variableType = call.Args[0].VariableTypes[0]
		}
		stack, _ := m.build(class, code, call.Args[0], context, state)
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
		if variableType.IsPointer() == false {
			primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		code.CodeLength++
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Method:     "containsKey",
			Descriptor: "(Ljava/lang/Object;)Z",
		}, code.Codes[code.CodeLength:code.CodeLength+2])
		code.CodeLength += 2
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	case common.MAP_METHOD_REMOVE:
		currentStack := uint16(1)
		callRemove := func() {
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      java_hashmap_class,
				Method:     "remove",
				Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = cg.OP_pop
			code.CodeLength += 4
		}
		for k, v := range call.Args {
			currentStack = 1
			if v.MayHaveMultiValue() && len(v.VariableTypes) > 1 {
				stack, _ := m.build(class, code, v, context, state)
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
						state.Stacks = append(state.Stacks,
							state.newStackMapVerificationTypeInfo(class, hashMapVerifyType)...)
					}
					//load
					m.buildLoadArrayListAutoVar(code, context)
					loadInt32(class, code, int32(kk))
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      java_arrylist_class,
						Method:     "get",
						Descriptor: "(I)Ljava/lang/Object;",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
					//remove
					callRemove()
					state.popStack(1)
				}
				continue
			}
			variableType := v.VariableType
			if v.MayHaveMultiValue() {
				variableType = v.VariableTypes[0]
			}
			if k == len(call.Args)-1 {
			} else { // not last one
				code.Codes[code.CodeLength] = cg.OP_dup
				currentStack++
				if currentStack > maxstack {
					maxstack = currentStack
				}
				state.Stacks = append(state.Stacks,
					state.newStackMapVerificationTypeInfo(class, hashMapVerifyType)...)
			}
			stack, _ := m.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			if variableType.IsPointer() == false {
				primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, variableType)
			}
			//call remove
			callRemove()
			state.popStack(1)
		}
	case common.MAP_METHOD_REMOVEALL:
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Method:     "clear",
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case common.MAP_METHOD_SIZE:
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_hashmap_class,
			Method:     "size",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}
