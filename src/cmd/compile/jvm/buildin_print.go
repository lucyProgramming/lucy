package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	function print
*/
func (m *MakeExpression) mkBuildinPrint(class *cg.ClassHighLevel, code *cg.AttributeCode, call *ast.ExpressionFunctionCall, context *Context, state *StackMapState) (maxstack uint16) {
	code.Codes[code.CodeLength] = cg.OP_getstatic
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      "java/lang/System",
		Field:      "out",
		Descriptor: "Ljava/io/PrintStream;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 1
	if len(call.Args) == 0 {
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/io/PrintStream",
			Method:     "println",
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	{
		t := &ast.VariableType{}
		t.Typ = ast.VARIABLE_TYPE_OBJECT
		t.Class = &ast.Class{}
		t.Class.Name = "java/io/PrintStream"
		state.pushStack(class, t)
	}
	defer func() {
		// print have no return value,stack is empty
		state.Stacks = []*cg.StackMap_verification_type_info{}
	}()
	if len(call.Args) == 1 && call.Args[0].HaveOnlyOneValue() {
		stack, es := m.build(class, code, call.Args[0], context, state)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			state.pushStack(class, call.Args[0].Value)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)
		}
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
		switch call.Args[0].GetTheOnlyOneVariableType().Typ {
		case ast.VARIABLE_TYPE_BOOL:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Z)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_ENUM:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(I)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(J)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(F)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(D)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_STRING:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Ljava/lang/String;)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_JAVA_ARRAY:
			fallthrough
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY:
			fallthrough
		case ast.VARIABLE_TYPE_MAP:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/io/PrintStream",
				Method:     "println",
				Descriptor: "(Ljava/lang/Object;)V",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}

	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Method:     special_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 3
	currentStack := uint16(2)
	app := func(isLast bool) {
		//
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      "java/lang/StringBuilder",
			Method:     "append",
			Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if isLast == false {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertStringConst(" ", code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/StringBuilder",
				Method:     "append",
				Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	}
	{
		t := &ast.VariableType{}
		t.Typ = ast.VARIABLE_TYPE_OBJECT
		t.Class = &ast.Class{}
		t.Class.Name = "java/lang/StringBuilder"
		state.pushStack(class, t)
	}

	for k, v := range call.Args {
		var variableType *ast.VariableType
		if v.MayHaveMultiValue() && len(v.Values) > 1 {
			stack, _ := m.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			arrayListPacker.buildStoreArrayListAutoVar(code, context)
			for kk, tt := range v.Values {
				stack = arrayListPacker.unPack(class, code, kk, tt, context)
				if t := stack + currentStack; t > maxstack {
					maxstack = t
				}
				m.stackTop2String(class, code, tt, context, state)
				if tt.IsPointer() && tt.Typ != ast.VARIABLE_TYPE_STRING {
					if t := 2 + currentStack; t > maxstack {
						maxstack = t
					}
				}
				app(k == len(call.Args)-1 && kk == len(v.Values)-1) // last and last
			}
			continue
		}
		variableType = v.Value
		if v.MayHaveMultiValue() {
			variableType = v.Values[0]
		}
		stack, es := m.build(class, code, v, context, state)
		if len(es) > 0 {
			backPatchEs(es, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
			state.pushStack(class, variableType)
			state.popStack(1)
		}
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		m.stackTop2String(class, code, variableType, context, state)
		app(k == len(call.Args)-1)
	}
	// tostring
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/lang/StringBuilder",
		Method:     "toString",
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// call println
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      "java/io/PrintStream",
		Method:     "println",
		Descriptor: "(Ljava/lang/String;)V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
