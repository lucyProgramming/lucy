package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) stackTop2Byte(code *cg.AttributeCode, on ast.VariableTypeKind) {
	switch on {
	case ast.VariableTypeByte:
		// already is
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2b
		code.CodeLength++
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2b
		code.CodeLength += 2
	}
}

func (buildExpression *BuildExpression) stackTop2Short(code *cg.AttributeCode, on ast.VariableTypeKind) {
	switch on {
	case ast.VariableTypeByte:
		// already is
	case ast.VariableTypeShort:
		// already is
	case ast.VariableTypeChar:
		code.Codes[code.CodeLength] = cg.OP_i2s
		code.CodeLength++
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2s
		code.CodeLength++
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.Codes[code.CodeLength+1] = cg.OP_i2s
		code.CodeLength += 2

	}
}

func (buildExpression *BuildExpression) stackTop2Char(code *cg.AttributeCode, on ast.VariableTypeKind) {
	switch on {
	case ast.VariableTypeByte:
		// already is
	case ast.VariableTypeShort:
		// already is
	case ast.VariableTypeChar:
		// already is
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2c
		code.CodeLength++
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_i2c
		code.CodeLength++
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_i2c
		code.CodeLength++
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.CodeLength++
		code.Codes[code.CodeLength] = cg.OP_i2c
		code.CodeLength++

	}
}
func (buildExpression *BuildExpression) stackTop2Int(code *cg.AttributeCode, on ast.VariableTypeKind) {
	switch on {
	case ast.VariableTypeByte:
		// already is
	case ast.VariableTypeShort:
		// already is
	case ast.VariableTypeChar:
		// already is
	case ast.VariableTypeInt:
		// already is
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.CodeLength++
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.CodeLength++
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.CodeLength++

	}
}

func (buildExpression *BuildExpression) stackTop2Float(code *cg.AttributeCode, on ast.VariableTypeKind) {
	switch on {
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2f
		code.CodeLength++
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2f
		code.CodeLength++
	case ast.VariableTypeFloat:
		// already is
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2f
		code.CodeLength++

	}
}

func (buildExpression *BuildExpression) stackTop2Long(code *cg.AttributeCode, on ast.VariableTypeKind) {
	switch on {
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2l
		code.CodeLength++
	case ast.VariableTypeLong:
		// already is
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2l
		code.CodeLength++
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_d2l
		code.CodeLength++

	}
}

func (buildExpression *BuildExpression) stackTop2Double(code *cg.AttributeCode, on ast.VariableTypeKind) {
	switch on {
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_i2d
		code.CodeLength++
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_l2d
		code.CodeLength++
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_f2d
		code.CodeLength++
	case ast.VariableTypeDouble:
		// already is
	}
}

/*
	convert stack top to target
*/
func (buildExpression *BuildExpression) numberTypeConverter(code *cg.AttributeCode,
	on ast.VariableTypeKind, target ast.VariableTypeKind) {
	if on == target {
		return
	}
	switch target {
	case ast.VariableTypeByte:
		buildExpression.stackTop2Byte(code, on)
	case ast.VariableTypeShort:
		buildExpression.stackTop2Short(code, on)
	case ast.VariableTypeChar:
		buildExpression.stackTop2Char(code, on)
	case ast.VariableTypeInt:
		buildExpression.stackTop2Int(code, on)
	case ast.VariableTypeLong:
		buildExpression.stackTop2Long(code, on)
	case ast.VariableTypeFloat:
		buildExpression.stackTop2Float(code, on)
	case ast.VariableTypeDouble:
		buildExpression.stackTop2Double(code, on)
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
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(Z)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeChar:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(C)Ljava/lang/String;",
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
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(I)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(J)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(F)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      javaStringClass,
			Method:     "valueOf",
			Descriptor: "(D)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	default:
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
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      "java/lang/Object",
			Method:     "toString",
			Descriptor: "()Ljava/lang/String;",
		}, code.Codes[code.CodeLength+11:code.CodeLength+13])
		code.CodeLength += 13
	}
	return
}
