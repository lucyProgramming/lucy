package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildReturnStatement(class *cg.ClassHighLevel, code *cg.AttributeCode, s *ast.StatementReturn, context *Context) (maxstack uint16) {
	if len(s.Function.Typ.ReturnList) == 0 {
		code.Codes[code.CodeLength] = cg.OP_return
		code.CodeLength++
		return
	}
	if len(s.Function.Typ.ReturnList) == 1 {
		if len(s.Expressions) != 1 {
			panic("this is not happening")
		}
		var es []*cg.JumpBackPatch
		maxstack, es = m.MakeExpression.build(class, code, s.Expressions[0], context)
		backPatchEs(es, code.CodeLength)
		switch s.Function.Typ.ReturnList[0].Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_ireturn
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lreturn
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_freturn
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dreturn
		case ast.VARIABLE_TYPE_STRING:
			fallthrough
		case ast.VARIABLE_TYPE_OBJECT:
			fallthrough
		case ast.VARIABLE_TYPE_ARRAY_INSTANCE:
			code.Codes[code.CodeLength] = cg.OP_areturn
		default:
			panic("......a")
		}
		code.CodeLength++
		return
	}
	//new a array list
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(arrylistclassname, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup // dup on stack
	code.CodeLength += 4
	//call init
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      arrylistclassname,
		Name:       specail_method_init,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	maxstack = 2 // max stack is 2
	currentStack := uint16(1)
	for _, v := range s.Expressions {
		code.Codes[code.CodeLength] = cg.OP_dup // dup array list
		code.CodeLength++
		currentStack++
		if currentStack > maxstack {
			maxstack = maxstack
		}
		if (v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL) && len(v.VariableTypes) > 0 {
			if currentStack > maxstack {
				maxstack = maxstack
			} // make the call
			stack, _ := m.MakeExpression.build(class, code, v, context)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/util/ArrayList",
				Name:       "addAll",
				Descriptor: "(Ljava/util/Collection;)Z",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = cg.OP_pop
			code.CodeLength += 4
			continue
		}
		stack, es := m.MakeExpression.build(class, code, v, context)
		backPatchEs(es, code.CodeLength)
		if t := stack + currentStack; t > maxstack {
			maxstack = t
		}
		//convert to object
		switch v.VariableType.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Integer",
				Name:       "valueOf",
				Descriptor: "(I)Ljava/lang/Integer;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Float",
				Name:       "valueOf",
				Descriptor: "(F)Ljava/lang/Float;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Double",
				Name:       "valueOf",
				Descriptor: "(D)Ljava/lang/Double;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_invokestatic
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      "java/lang/Long",
				Name:       "valueOf",
				Descriptor: "(J)Ljava/lang/Long;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		// append
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      arrylistclassname,
			Name:       "add",
			Descriptor: "(Ljava/lang/Object;)Z",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_pop
		code.CodeLength += 4
		currentStack = 1
	}
	code.Codes[code.CodeLength] = cg.OP_areturn
	code.CodeLength++
	return
}
