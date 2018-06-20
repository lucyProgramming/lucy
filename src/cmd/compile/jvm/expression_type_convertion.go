package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildTypeConversion(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	{
		length := len(state.Stacks)
		defer func() {
			state.popStack(len(state.Stacks) - length)
		}()
	}
	conversion := e.Data.(*ast.ExpressionTypeConversion)
	currentStack := uint16(0)
	// []byte("aaaaaaaaaaaa")
	if conversion.Type.Type == ast.VARIABLE_TYPE_ARRAY &&
		conversion.Type.ArrayType.Type == ast.VARIABLE_TYPE_BYTE {
		currentStack = 2
		meta := ArrayMetas[ast.VARIABLE_TYPE_BYTE]
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
	if (conversion.Type.Type == ast.VARIABLE_TYPE_STRING && conversion.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_ARRAY &&
		conversion.Expression.ExpressionValue.ArrayType.Type == ast.VARIABLE_TYPE_BYTE) ||
		(conversion.Type.Type == ast.VARIABLE_TYPE_STRING && conversion.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_JAVA_ARRAY &&
			conversion.Expression.ExpressionValue.ArrayType.Type == ast.VARIABLE_TYPE_BYTE) {
		currentStack = 2
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst(java_string_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		t := &cg.StackMapVerificationTypeInfo{}
		t.Verify = &cg.StackMapUninitializedVariableInfo{
			CodeOffset: uint16(code.CodeLength),
		}
		state.Stacks = append(state.Stacks, t, t)
		code.CodeLength += 4
	}
	stack, _ := makeExpression.build(class, code, conversion.Expression, context, state)
	maxStack = currentStack + stack
	if e.ExpressionValue.IsNumber() {
		makeExpression.numberTypeConverter(code, conversion.Expression.ExpressionValue.Type, conversion.Type.Type)
		if t := jvmSize(conversion.Type); t > maxStack {
			maxStack = t
		}
		return
	}
	//  []byte("hello world")
	if conversion.Type.Type == ast.VARIABLE_TYPE_ARRAY && conversion.Type.ArrayType.Type == ast.VARIABLE_TYPE_BYTE &&
		conversion.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_STRING {
		//stack top must be a string
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "getBytes",
			Descriptor: "()[B",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if 3 > maxStack { // arraybyteref arraybyteref byte[]
			maxStack = 3
		}
		meta := ArrayMetas[ast.VARIABLE_TYPE_BYTE]
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     special_method_init,
			Descriptor: meta.constructorFuncDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	// byte[]("hello world")
	if conversion.Type.Type == ast.VARIABLE_TYPE_JAVA_ARRAY && conversion.Type.ArrayType.Type == ast.VARIABLE_TYPE_BYTE &&
		conversion.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_STRING {
		//stack top must be a string
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
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
	if conversion.Type.Type == ast.VARIABLE_TYPE_STRING &&
		conversion.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_ARRAY &&
		conversion.Expression.ExpressionValue.ArrayType.Type == ast.VARIABLE_TYPE_BYTE {
		meta := ArrayMetas[ast.VARIABLE_TYPE_BYTE]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.className,
			Method:     "getJavaArray",
			Descriptor: "()[B",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     special_method_init,
			Descriptor: "([B)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	//  string(byte[])
	if conversion.Type.Type == ast.VARIABLE_TYPE_STRING &&
		conversion.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_JAVA_ARRAY &&
		conversion.Expression.ExpressionValue.ArrayType.Type == ast.VARIABLE_TYPE_BYTE {
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     special_method_init,
			Descriptor: "([B)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}

	if conversion.Type.Type == ast.VARIABLE_TYPE_STRING {
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(java_string_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}

	// objects
	code.Codes[code.CodeLength] = cg.OP_checkcast
	if conversion.Type.Type == ast.VARIABLE_TYPE_OBJECT {
		class.InsertClassConst(conversion.Type.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else if conversion.Type.Type == ast.VARIABLE_TYPE_ARRAY { // arrays
		meta := ArrayMetas[conversion.Type.ArrayType.Type]
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
	} else {
		class.InsertClassConst(Descriptor.typeDescriptor(conversion.Type), code.Codes[code.CodeLength+1:code.CodeLength+3])
	}
	code.CodeLength += 3
	return
}

func (makeExpression *MakeExpression) stackTop2Byte(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		// already is
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2b
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	}
}

func (makeExpression *MakeExpression) stackTop2Short(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		// already is
	case ast.VARIABLE_TYPE_SHORT:
		// already is
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2s
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	}
}

func (makeExpression *MakeExpression) stackTop2Int(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		// already is
	case ast.VARIABLE_TYPE_SHORT:
		// already is
	case ast.VARIABLE_TYPE_INT:
		// already is
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.CodeLength++
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.CodeLength++
	}
}

func (makeExpression *MakeExpression) stackTop2Float(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2f
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		// already is
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2f
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2f
		code.CodeLength++
	}
}

func (makeExpression *MakeExpression) stackTop2Long(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2l
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2l
		code.CodeLength++
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2l
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		// already is
	}
}

func (makeExpression *MakeExpression) stackTop2Double(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2d
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2d
		code.CodeLength++
	case ast.VARIABLE_TYPE_DOUBLE:
		// already is
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2d
		code.CodeLength++
	}
}

/*
	convert stack top to target
*/
func (makeExpression *MakeExpression) numberTypeConverter(code *cg.AttributeCode, typ int, target int) {
	if typ == target {
		return
	}
	switch target {
	case ast.VARIABLE_TYPE_BYTE:
		makeExpression.stackTop2Byte(code, typ)
	case ast.VARIABLE_TYPE_SHORT:
		makeExpression.stackTop2Short(code, typ)
	case ast.VARIABLE_TYPE_INT:
		makeExpression.stackTop2Int(code, typ)
	case ast.VARIABLE_TYPE_LONG:
		makeExpression.stackTop2Long(code, typ)
	case ast.VARIABLE_TYPE_FLOAT:
		makeExpression.stackTop2Float(code, typ)
	case ast.VARIABLE_TYPE_DOUBLE:
		makeExpression.stackTop2Double(code, typ)
	}
}

func (makeExpression *MakeExpression) stackTop2String(class *cg.ClassHighLevel, code *cg.AttributeCode,
	typ *ast.Type, context *Context, state *StackMapState) (maxstack uint16) {
	if typ.Type == ast.VARIABLE_TYPE_STRING {
		return
	}
	maxstack = jvmSize(typ)
	switch typ.Type {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "valueOf",
			Descriptor: "(Z)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "valueOf",
			Descriptor: "(I)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "valueOf",
			Descriptor: "(J)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "valueOf",
			Descriptor: "(F)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "valueOf",
			Descriptor: "(D)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY:
		fallthrough
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		if 2 > maxstack {
			maxstack = 2
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		{
			state.pushStack(class, typ)
			context.MakeStackMap(code, state, code.CodeLength+10)
			state.popStack(1)
			state.pushStack(class, &ast.Type{Type: ast.VARIABLE_TYPE_STRING})
			context.MakeStackMap(code, state, code.CodeLength+13)
			state.popStack(1)
		}
		code.Codes[code.CodeLength] = cg.OP_ifnonnull
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 10)
		code.Codes[code.CodeLength+3] = cg.OP_pop
		code.Codes[code.CodeLength+4] = cg.OP_ldc_w
		class.InsertStringConst("null", code.Codes[code.CodeLength+5:code.CodeLength+7])
		code.Codes[code.CodeLength+7] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+8:code.CodeLength+10], 6)
		code.Codes[code.CodeLength+10] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Object",
			Method:     "toString",
			Descriptor: "()Ljava/lang/String;",
		}, code.Codes[code.CodeLength+11:code.CodeLength+13])
		code.CodeLength += 13
	}

	return
}
