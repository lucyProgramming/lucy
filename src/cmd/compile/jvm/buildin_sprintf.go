package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) mkBuildInSprintf(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	// format,must be string
	call := e.Data.(*ast.ExpressionFunctionCall)
	meta := call.BuildInFunctionMeta.(*ast.BuildInFunctionSprintfMeta)
	maxStack, _ = makeExpression.build(class, code, meta.Format, context, state)
	state.pushStack(class, state.newObjectVariableType(javaStringClass))
	loadInt32(class, code, int32(meta.ArgsLength))
	code.Codes[code.CodeLength] = cg.OP_anewarray
	class.InsertClassConst("java/lang/Object", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	currentStack := uint16(2)
	if currentStack > maxStack {
		maxStack = currentStack
	}

	objectArray := &ast.Type{}
	objectArray.Type = ast.VariableTypeJavaArray
	objectArray.ArrayType = state.newObjectVariableType(javaRootClass)
	state.pushStack(class, objectArray)
	index := int32(0)
	for _, v := range call.Args {
		if v.MayHaveMultiValue() && len(v.ExpressionMultiValues) > 1 {
			currentStack = 2
			stack, _ := makeExpression.build(class, code, v, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			// store in temp var
			multiValuePacker.storeMultiValueAutoVar(code, context)
			for kk, _ := range v.ExpressionMultiValues {
				currentStack = 2
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
		currentStack = 2
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, index)
		currentStack += 2
		state.pushStack(class, objectArray)
		state.pushStack(class, &ast.Type{Type: ast.VariableTypeInt})
		stack, es := makeExpression.build(class, code, v, context, state)
		if len(es) > 0 {
			fillOffsetForExits(es, code.CodeLength)
			state.pushStack(class, v.ExpressionValue)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1) // bool value
		}
		state.popStack(2)
		if t := currentStack + stack; t > maxStack {
			maxStack = t
		}
		if v.ExpressionValue.IsPointer() == false {
			typeConverter.packPrimitives(class, code, v.ExpressionValue)
		}
		code.Codes[code.CodeLength] = cg.OP_aastore
		code.CodeLength++
		index++
	}
	code.Codes[code.CodeLength] = cg.OP_invokestatic
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringClass,
		Method:     "format",
		Descriptor: "(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;",
	},
		code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if e.IsStatementExpression {
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
	}

	return
}
