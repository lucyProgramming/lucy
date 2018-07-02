package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TypeConverterAndPrimitivePacker struct {
}

func (TypeConverterAndPrimitivePacker) unPackPrimitives(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.Type) {
	switch t.Type {
	case ast.VariableTypeBool:
		c := "java/lang/Boolean"
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(c, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      c,
			Method:     "booleanValue",
			Descriptor: "()Z",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeByte:
		c := "java/lang/Byte"
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(c, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      c,
			Method:     "byteValue",
			Descriptor: "()B",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeShort:
		c := "java/lang/Short"
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(c, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      c,
			Method:     "shortValue",
			Descriptor: "()S",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst("java/lang/Integer", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaIntegerClass,
			Method:     "intValue",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeLong:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(javaLongClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaLongClass,
			Method:     "longValue",
			Descriptor: "()J",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeFloat:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(javaFloatClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaFloatClass,
			Method:     "floatValue",
			Descriptor: "()F",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeDouble:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(javaDoubleClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaDoubleClass,
			Method:     "doubleValue",
			Descriptor: "()D",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
}

func (c *TypeConverterAndPrimitivePacker) packPrimitives(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.Type) {
	copyOPs(code, c.packPrimitivesBytes(class, t)...)
}

func (TypeConverterAndPrimitivePacker) packPrimitivesBytes(class *cg.ClassHighLevel, t *ast.Type) (bs []byte) {
	bs = make([]byte, 3)
	bs[0] = cg.OP_invokestatic
	switch t.Type {
	case ast.VariableTypeBool:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Boolean",
			Method:     "valueOf",
			Descriptor: "(Z)Ljava/lang/Boolean;",
		}, bs[1:3])
	case ast.VariableTypeByte:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Byte",
			Method:     "valueOf",
			Descriptor: "(B)Ljava/lang/Byte;",
		}, bs[1:3])
	case ast.VariableTypeShort:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Short",
			Method:     "valueOf",
			Descriptor: "(S)Ljava/lang/Short;",
		}, bs[1:3])
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaIntegerClass,
			Method:     "valueOf",
			Descriptor: "(I)Ljava/lang/Integer;",
		}, bs[1:3])
	case ast.VariableTypeFloat:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaFloatClass,
			Method:     "valueOf",
			Descriptor: "(F)Ljava/lang/Float;",
		}, bs[1:3])
	case ast.VariableTypeDouble:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaDoubleClass,
			Method:     "valueOf",
			Descriptor: "(D)Ljava/lang/Double;",
		}, bs[1:3])
	case ast.VariableTypeLong:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaLongClass,
			Method:     "valueOf",
			Descriptor: "(J)Ljava/lang/Long;",
		}, bs[1:3])
	}
	return
}

func (TypeConverterAndPrimitivePacker) castPointerTypeToRealType(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.Type) {
	if t.IsPointer() == false {
		panic("...")
	}
	code.Codes[code.CodeLength] = cg.OP_checkcast
	switch t.Type {
	case ast.VariableTypeString:
		class.InsertClassConst(javaStringClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeObject:
		if t.Class.Name != ast.JavaRootClass {
			class.InsertClassConst(t.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	case ast.VariableTypeArray:
		meta := ArrayMetas[t.Array.Type]
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeMap:
		class.InsertClassConst(javaMapClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeJavaArray:
		class.InsertClassConst(JvmDescriptor.typeDescriptor(t), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeFunction:
		class.InsertClassConst(javaMethodHandleClass, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	default:
		panic("1")
	}
}
