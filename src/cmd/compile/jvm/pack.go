package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TypeConverterAndPrimitivePacker struct {
}

func (TypeConverterAndPrimitivePacker) unPackPrimitives(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.Type) {
	switch t.Type {
	case ast.VARIABLE_TYPE_BOOL:
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
	case ast.VARIABLE_TYPE_BYTE:
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
	case ast.VARIABLE_TYPE_SHORT:
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
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst("java/lang/Integer", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_integer_class,
			Method:     "intValue",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(java_long_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_long_class,
			Method:     "longValue",
			Descriptor: "()J",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3

	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(java_float_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_float_class,
			Method:     "floatValue",
			Descriptor: "()F",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_checkcast
		class.InsertClassConst(java_double_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_double_class,
			Method:     "doubleValue",
			Descriptor: "()D",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
}

func (c *TypeConverterAndPrimitivePacker) packPrimitives(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.Type) {
	copyOP(code, c.packPrimitivesBytes(class, t)...)
}

func (TypeConverterAndPrimitivePacker) packPrimitivesBytes(class *cg.ClassHighLevel, t *ast.Type) (bs []byte) {
	bs = make([]byte, 3)
	bs[0] = cg.OP_invokestatic
	switch t.Type {
	case ast.VARIABLE_TYPE_BOOL:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Boolean",
			Method:     "valueOf",
			Descriptor: "(Z)Ljava/lang/Boolean;",
		}, bs[1:3])
	case ast.VARIABLE_TYPE_BYTE:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Byte",
			Method:     "valueOf",
			Descriptor: "(B)Ljava/lang/Byte;",
		}, bs[1:3])
	case ast.VARIABLE_TYPE_SHORT:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Short",
			Method:     "valueOf",
			Descriptor: "(S)Ljava/lang/Short;",
		}, bs[1:3])
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_integer_class,
			Method:     "valueOf",
			Descriptor: "(I)Ljava/lang/Integer;",
		}, bs[1:3])
	case ast.VARIABLE_TYPE_FLOAT:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_float_class,
			Method:     "valueOf",
			Descriptor: "(F)Ljava/lang/Float;",
		}, bs[1:3])
	case ast.VARIABLE_TYPE_DOUBLE:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_double_class,
			Method:     "valueOf",
			Descriptor: "(D)Ljava/lang/Double;",
		}, bs[1:3])
	case ast.VARIABLE_TYPE_LONG:
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_long_class,
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
	case ast.VARIABLE_TYPE_STRING:
		class.InsertClassConst(java_string_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
		if t.Class.Name != ast.JAVA_ROOT_CLASS {
			class.InsertClassConst(t.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[t.ArrayType.Type]
		class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_MAP:
		class.InsertClassConst(java_hashmap_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_JAVA_ARRAY:
		class.InsertClassConst(Descriptor.typeDescriptor(t), code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
}
