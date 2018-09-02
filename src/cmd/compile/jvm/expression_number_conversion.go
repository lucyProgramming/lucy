package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) stackTop2Byte(code *cg.AttributeCode, typ ast.VariableTypeKind) {
	switch typ {
	case ast.VariableTypeByte:
		// already is
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2b
		code.CodeLength++
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	}
}

func (buildExpression *BuildExpression) stackTop2Short(code *cg.AttributeCode, typ ast.VariableTypeKind) {
	switch typ {
	case ast.VariableTypeByte:
		// already is
	case ast.VariableTypeShort:
		// already is
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2s
		code.CodeLength++
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	}
}

func (buildExpression *BuildExpression) stackTop2Int(code *cg.AttributeCode, typ ast.VariableTypeKind) {
	switch typ {
	case ast.VariableTypeByte:
		// already is
	case ast.VariableTypeShort:
		// already is
	case ast.VariableTypeInt:
		// already is
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.CodeLength++
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.CodeLength++
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.CodeLength++
	}
}

func (buildExpression *BuildExpression) stackTop2Float(code *cg.AttributeCode, typ ast.VariableTypeKind) {
	switch typ {
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2f
		code.CodeLength++
	case ast.VariableTypeFloat:
		// already is
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2f
		code.CodeLength++
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2f
		code.CodeLength++
	}
}

func (buildExpression *BuildExpression) stackTop2Long(code *cg.AttributeCode, typ ast.VariableTypeKind) {
	switch typ {
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2l
		code.CodeLength++
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2l
		code.CodeLength++
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2l
		code.CodeLength++
	case ast.VariableTypeLong:
		// already is
	}
}

func (buildExpression *BuildExpression) stackTop2Double(code *cg.AttributeCode, typ ast.VariableTypeKind) {
	switch typ {
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2d
		code.CodeLength++
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2d
		code.CodeLength++
	case ast.VariableTypeDouble:
		// already is
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2d
		code.CodeLength++
	}
}

/*
	convert stack top to target
*/
func (buildExpression *BuildExpression) numberTypeConverter(code *cg.AttributeCode, typ ast.VariableTypeKind, target ast.VariableTypeKind) {
	if typ == target {
		return
	}
	{
		t1 := &ast.Type{
			Type: typ,
		}
		t2 := &ast.Type{
			Type: target,
		}
		if t1.IsNumber() == false ||
			t2.IsNumber() == false {
			panic("not a number type")
		}
	}
	switch target {
	case ast.VariableTypeByte:
		buildExpression.stackTop2Byte(code, typ)
	case ast.VariableTypeShort:
		buildExpression.stackTop2Short(code, typ)
	case ast.VariableTypeInt:
		buildExpression.stackTop2Int(code, typ)
	case ast.VariableTypeLong:
		buildExpression.stackTop2Long(code, typ)
	case ast.VariableTypeFloat:
		buildExpression.stackTop2Float(code, typ)
	case ast.VariableTypeDouble:
		buildExpression.stackTop2Double(code, typ)
	}
}

func (buildExpression *BuildExpression) stackTop2String(class *cg.ClassHighLevel, code *cg.AttributeCode,
	typ *ast.Type, context *Context, state *StackMapState) (maxStack uint16) {
	if typ.Type == ast.VariableTypeString {
		return
	}
	maxStack = jvmSlotSize(typ) * 2
	switch typ.Type {
	case ast.VariableTypeBool:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(Z)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(I)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(J)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(F)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(D)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeArray:
		fallthrough
	case ast.VariableTypeJavaArray:
		fallthrough
	case ast.VariableTypeMap:
		if 2 > maxStack {
			maxStack = 2
		}
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		{
			state.pushStack(class, typ)
			context.MakeStackMap(code, state, code.CodeLength+10)
			state.popStack(1)
			state.pushStack(class, &ast.Type{Type: ast.VariableTypeString})
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
