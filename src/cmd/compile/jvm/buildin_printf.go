package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	function printf
*/
func (m *MakeExpression) mkBuildinPrintf(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	call := e.Data.(*ast.ExpressionFunctionCall)
	meta := call.BuildinFunctionMeta.(*ast.BuildinFunctionPrintfMeta)
	if meta.Stream == nil {
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      "java/lang/System",
			Field:      "out",
			Descriptor: "Ljava/io/PrintStream;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxStack = 1
	} else { // get stream from args
		maxStack, _ = m.build(class, code, meta.Stream, context, state)
	}
	state.pushStack(class, state.newObjectVariableType(java_print_stream_class))
	stack, _ := m.build(class, code, meta.Format, context, state)
	if t := 1 + stack; t > maxStack {
		maxStack = t
	}
	state.pushStack(class, state.newObjectVariableType(java_string_class))
	loadInt32(class, code, int32(meta.ArgsLength))
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst("java/lang/Object", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	currentStack := uint16(3)
	if currentStack > maxStack {
		maxStack = currentStack
	}
	objectArray := &ast.VariableType{}
	objectArray.Typ = ast.VARIABLE_TYPE_JAVA_ARRAY
	objectArray.ArrayType = state.newObjectVariableType(java_root_class)
	state.pushStack(class, objectArray)

	index := int32(0)
	for _, v := range call.Args {
		if v.MayHaveMultiValue() && len(v.Values) > 1 {
			currentStack = 3
			stack, _ := m.build(class, code, v, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			// store in temp var
			multiValuePacker.storeArrayListAutoVar(code, context)
			for kk, _ := range v.Values {
				currentStack = 3
				code.Codes[code.CodeLength] = cg.OP_dup
				code.CodeLength++
				loadInt32(class, code, index)
				currentStack += 2
				stack = multiValuePacker.unPackObject(class, code, kk, context)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				code.Codes[code.CodeLength] = cg.OP_aastore
				code.CodeLength++
				index++
			}
			continue
		}
		currentStack = 3
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, index)
		currentStack += 2
		state.pushStack(class, objectArray)
		state.pushStack(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_INT})
		stack, es := m.build(class, code, v, context, state)
		if len(es) > 0 {
			backfillExit(es, code.CodeLength)
			state.pushStack(class, v.Value)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1) // bool value
		}
		state.popStack(2)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.Value.IsPointer() == false {
			typeConverter.putPrimitiveInObject(class, code, v.Value)
		}
		code.Codes[code.CodeLength] = cg.OP_aastore
		code.CodeLength++
		index++
	}
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      java_print_stream_class,
		Method:     "printf",
		Descriptor: "(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/PrintStream;",
	},
		code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_pop
	code.CodeLength += 4
	return
}
