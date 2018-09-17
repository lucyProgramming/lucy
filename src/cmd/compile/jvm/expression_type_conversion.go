package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildTypeConversion(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	conversion := e.Data.(*ast.ExpressionTypeConversion)
	currentStack := uint16(0)
	// []byte("aaaaaaaaaaaa")
	if conversion.Type.Type == ast.VariableTypeArray &&
		conversion.Type.Array.Type == ast.VariableTypeByte {
		currentStack = 2
		meta := ArrayMetas[ast.VariableTypeByte]
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		t := &cg.StackMapVerificationTypeInfo{}
		t.Verify = &cg.StackMapUninitializedVariableInfo{
			CodeOffset: uint16(code.CodeLength),
		}
		state.Stacks = append(state.Stacks, t, t)
		code.CodeLength += 4
	}
	// string
	if (conversion.Type.Type == ast.VariableTypeString && conversion.Expression.Value.Type == ast.VariableTypeArray &&
		conversion.Expression.Value.Array.Type == ast.VariableTypeByte) ||
		(conversion.Type.Type == ast.VariableTypeString && conversion.Expression.Value.Type == ast.VariableTypeJavaArray &&
			conversion.Expression.Value.Array.Type == ast.VariableTypeByte) {
		currentStack = 2
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst(javaStringClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		t := &cg.StackMapVerificationTypeInfo{}
		t.Verify = &cg.StackMapUninitializedVariableInfo{
			CodeOffset: uint16(code.CodeLength),
		}
		state.Stacks = append(state.Stacks, t, t)
		code.CodeLength += 4
	}
	stack := buildExpression.build(class, code, conversion.Expression, context, state)
	maxStack = currentStack + stack
	if e.Value.IsNumber() {
		buildExpression.numberTypeConverter(code, conversion.Expression.Value.Type, conversion.Type.Type)
		if t := jvmSlotSize(conversion.Type); t > maxStack {
			maxStack = t
		}
		return
	}
	// int(enum)
	if conversion.Type.Type == ast.VariableTypeInt &&
		conversion.Expression.Value.Type == ast.VariableTypeEnum {
		return
	}
	// enum(int)
	if conversion.Type.Type == ast.VariableTypeEnum &&
		conversion.Expression.Value.Type == ast.VariableTypeInt {
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, conversion.Type.Enum.DefaultValue)
		wrongExit := (&cg.Exit{}).Init(cg.OP_if_icmplt, code)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		loadInt32(class, code, conversion.Type.Enum.Enums[len(conversion.Type.Enum.Enums)-1].Value)
		wrongExit2 := (&cg.Exit{}).Init(cg.OP_if_icmpgt, code)
		okExit := (&cg.Exit{}).Init(cg.OP_goto, code)
		state.pushStack(class, conversion.Expression.Value)
		defer state.popStack(1)
		context.MakeStackMap(code, state, code.CodeLength)
		writeExits([]*cg.Exit{wrongExit, wrongExit2}, code.CodeLength)
		code.Codes[code.CodeLength] = cg.OP_pop
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst(javaExceptionClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst("int value not found in enum names",
			code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if 3 > maxStack {
			maxStack = 3
		}
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaExceptionClass,
			Method:     specialMethodInit,
			Descriptor: "(Ljava/lang/String;)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_athrow
		code.CodeLength++
		context.MakeStackMap(code, state, code.CodeLength)
		writeExits([]*cg.Exit{okExit}, code.CodeLength)
		return
	}

	// []byte("hello world")
	if conversion.Type.Type == ast.VariableTypeArray && conversion.Type.Array.Type == ast.VariableTypeByte &&
		conversion.Expression.Value.Type == ast.VariableTypeString {
		//stack top must be a string
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "getBytes",
			Descriptor: "()[B",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if 3 > maxStack { // arraybyteref arraybyteref byte[]
			maxStack = 3
		}
		meta := ArrayMetas[ast.VariableTypeByte]
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     specialMethodInit,
			Descriptor: meta.constructorFuncDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	// byte[]("hello world")
	if conversion.Type.Type == ast.VariableTypeJavaArray && conversion.Type.Array.Type == ast.VariableTypeByte &&
		conversion.Expression.Value.Type == ast.VariableTypeString {
		//stack top must be a string
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "getBytes",
			Descriptor: "()[B",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if 3 > maxStack { // arraybyteref arraybyteref byte[]
			maxStack = 3
		}
		return
	}
	//  string(['h','e'])
	if conversion.Type.Type == ast.VariableTypeString &&
		conversion.Expression.Value.Type == ast.VariableTypeArray &&
		conversion.Expression.Value.Array.Type == ast.VariableTypeByte {
		type autoVar struct {
			start  uint16
			length uint16
		}
		var a autoVar
		a.start = code.MaxLocals
		a.length = code.MaxLocals + 1
		state.appendLocals(class, &ast.Type{
			Type: ast.VariableTypeInt,
		})
		state.appendLocals(class, &ast.Type{
			Type: ast.VariableTypeInt,
		})
		code.MaxLocals += 2
		currentStack = 3
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		meta := ArrayMetas[ast.VariableTypeByte]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "start",
			Descriptor: "I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, a.start)...)
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     "size",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOPs(code, storeLocalVariableOps(ast.VariableTypeInt, a.length)...)
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      "elements",
			Descriptor: meta.elementsFieldDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, a.start)...)
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, a.length)...)
		code.Codes[code.CodeLength] = cg.OP_ldc_w
		class.InsertStringConst("utf-8", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if 6 > maxStack { // stack is ... stringRef stringRef byte[] start length "utf-8"
			maxStack = 6
		}
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     specialMethodInit,
			Descriptor: "([BIILjava/lang/String;)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	//  string(byte[])
	if conversion.Type.Type == ast.VariableTypeString &&
		conversion.Expression.Value.Type == ast.VariableTypeJavaArray &&
		conversion.Expression.Value.Array.Type == ast.VariableTypeByte {
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     specialMethodInit,
			Descriptor: "([B)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	if conversion.Type.Type == ast.VariableTypeString {
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(javaStringClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	// objects
	code.Codes[code.CodeLength] = cg.OP_checkcast
	code.CodeLength++
	insertTypeAssertClass(class, code, conversion.Type)
	return
}
