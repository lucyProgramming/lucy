package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildTypeConvertion(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	convertion := e.Data.(*ast.ExpressionTypeConvertion)
	currentStack := uint16(0)
	// []byte("aaaaaaaaaaaa")
	if convertion.Typ.Typ == ast.VARIABLE_TYPE_ARRAY && convertion.Typ.ArrayType.Typ == ast.VARIABLE_TYPE_BYTE {
		currentStack = 2
		meta := ArrayMetas[ast.VARIABLE_TYPE_BYTE]
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
	}
	if convertion.Typ.Typ == ast.VARIABLE_TYPE_STRING {
		currentStack = 2
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst(java_string_class, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
	}
	stack, es := m.build(class, code, convertion.Expression, context)
	backPatchEs(es, code.CodeLength)
	maxstack = currentStack + stack
	if convertion.Typ.IsNumber() {
		m.numberTypeConverter(code, convertion.Expression.VariableType.Typ, convertion.Typ.Typ)
		return
	}
	//  []byte("hello world")
	if convertion.Typ.Typ == ast.VARIABLE_TYPE_ARRAY && convertion.Typ.ArrayType.Typ == ast.VARIABLE_TYPE_BYTE {
		//stack top must be a string
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Name:       "getBytes",
			Descriptor: "()[B",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_dup
		code.CodeLength++
		if 4 > maxstack { // arraybyteref arraybyteref byte[] byte[]
			maxstack = 4
		}
		code.Codes[code.CodeLength] = cg.OP_arraylength
		code.CodeLength++
		meta := ArrayMetas[ast.VARIABLE_TYPE_BYTE]
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.classname,
			Name:       special_method_init,
			Descriptor: meta.constructorFuncDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	//  string(['h','e'])
	if convertion.Typ.Typ == ast.VARIABLE_TYPE_STRING {
		meta := ArrayMetas[ast.VARIABLE_TYPE_BYTE]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.classname,
			Name:       "getJavaArray",
			Descriptor: meta.getJavaArrayDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Name:       special_method_init,
			Descriptor: "([B)V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	return
}

func (m *MakeExpression) stackTop2Byte(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
	case ast.VARIABLE_TYPE_SHORT:
		code.Codes[code.CodeLength] = cg.OP_i2b
		code.CodeLength++
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
		code.Codes[code.CodeLength+1] = cg.OP_l2i
		code.CodeLength += 2
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Short(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
	case ast.VARIABLE_TYPE_SHORT:
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
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Int(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
	case ast.VARIABLE_TYPE_SHORT:
	case ast.VARIABLE_TYPE_INT:
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_f2i
		code.CodeLength++
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2i
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2i
		code.CodeLength++
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Float(code *cg.AttributeCode, typ int) {
	switch typ {
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_i2f
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_d2f
		code.CodeLength++
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2f
		code.CodeLength++
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Long(code *cg.AttributeCode, typ int) {
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
	default:
		panic("~~~~~~~~~~~~")
	}
}

func (m *MakeExpression) stackTop2Double(code *cg.AttributeCode, typ int) {
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
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_l2d
		code.CodeLength++

	default:
		panic("~~~~~~~~~~~~")
	}
}

/*
	convert stack top to target
*/
func (m *MakeExpression) numberTypeConverter(code *cg.AttributeCode, typ int, target int) {
	if typ == target {
		return
	}
	switch target {
	case ast.VARIABLE_TYPE_BYTE:
		m.stackTop2Byte(code, typ)
	case ast.VARIABLE_TYPE_SHORT:
		m.stackTop2Short(code, typ)
	case ast.VARIABLE_TYPE_INT:
		m.stackTop2Int(code, typ)
	case ast.VARIABLE_TYPE_LONG:
		m.stackTop2Long(code, typ)
	case ast.VARIABLE_TYPE_FLOAT:
		m.stackTop2Float(code, typ)
	case ast.VARIABLE_TYPE_DOUBLE:
		m.stackTop2Double(code, typ)
	default:
		panic(1)
	}
}

func (m *MakeExpression) stackTop2String(class *cg.ClassHighLevel, code *cg.AttributeCode, typ *ast.VariableType) {
	if typ.Typ == ast.VARIABLE_TYPE_STRING {
		return
	}
	switch typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Name:       "valueOf",
			Descriptor: "(Z)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Name:       "valueOf",
			Descriptor: "(I)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Name:       "valueOf",
			Descriptor: "(J)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Name:       "valueOf",
			Descriptor: "(F)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Name:       "valueOf",
			Descriptor: "(D)Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/Object",
			Name:       "toString",
			Descriptor: "()Ljava/lang/String;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	default:
		panic(1111111111)
	}
}
