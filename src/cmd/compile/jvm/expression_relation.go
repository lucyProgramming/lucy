package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildRelations(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if bin.Left.Value.IsNumber() { // in this case ,right must be a number type
		maxstack = 4
		stack, _ := m.build(class, code, bin.Left, context, state)
		if stack > maxstack {
			maxstack = stack
		}
		target := bin.Left.Value.NumberTypeConvertRule(bin.Right.Value)
		state.pushStack(class, &ast.VariableType{Typ: target})
		stack, _ = m.build(class, code, bin.Right, context, state)
		if t := 2 + stack; t > maxstack {
			maxstack = t
		}
		switch target {
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			code.Codes[code.CodeLength] = cg.OP_isub
		case ast.VARIABLE_TYPE_LONG:
			code.Codes[code.CodeLength] = cg.OP_lcmp
		case ast.VARIABLE_TYPE_FLOAT:
			code.Codes[code.CodeLength] = cg.OP_fcmpl
		case ast.VARIABLE_TYPE_DOUBLE:
			code.Codes[code.CodeLength] = cg.OP_dcmpl
		}
		code.CodeLength++
		state.popStack(1)
		context.MakeStackMap(code, state, code.CodeLength+7)

		state.pushStack(class, &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		if e.Typ == ast.EXPRESSION_TYPE_GT || e.Typ == ast.EXPRESSION_TYPE_LE { // > and <=
			code.Codes[code.CodeLength] = cg.OP_ifgt
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		} else if e.Typ == ast.EXPRESSION_TYPE_LT || e.Typ == ast.EXPRESSION_TYPE_GE { // < and >=
			code.Codes[code.CodeLength] = cg.OP_iflt
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		} else {
			code.Codes[code.CodeLength] = cg.OP_ifeq
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		}
		return
	}
	if bin.Left.Value.Typ == ast.VARIABLE_TYPE_BOOL ||
		bin.Right.Value.Typ == ast.VARIABLE_TYPE_BOOL { // bool type
		var es []*cg.JumpBackPatch
		state.pushStack(class, bin.Left.Value)
		maxstack, es = m.build(class, code, bin.Left, context, state)
		if len(es) > 0 {
			context.MakeStackMap(code, state, code.CodeLength)
			backPatchEs(es, code.CodeLength)
		}
		stack, es := m.build(class, code, bin.Right, context, state)
		state.pushStack(class, bin.Left.Value)
		if len(es) > 0 {
			context.MakeStackMap(code, state, code.CodeLength)
			backPatchEs(es, code.CodeLength)
		}
		state.popStack(2) // 2 bool value
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)

		code.Codes[code.CodeLength] = cg.OP_if_icmpeq
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		} else {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_1
		}
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength+7] = cg.OP_iconst_0
		}
		code.CodeLength += 8
		return
	}
	if bin.Left.Value.Typ == ast.VARIABLE_TYPE_NULL ||
		bin.Right.Value.Typ == ast.VARIABLE_TYPE_NULL { // must not null-null
		var notNullExpression *ast.Expression
		if bin.Left.Value.Typ != ast.VARIABLE_TYPE_NULL {
			notNullExpression = bin.Left
		} else {
			notNullExpression = bin.Right
		}
		maxstack, _ = m.build(class, code, notNullExpression, context, state)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_ifnull
		} else { // ne
			code.Codes[code.CodeLength] = cg.OP_ifnonnull
		}
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}

	//string compare
	if bin.Left.Value.Typ == ast.VARIABLE_TYPE_STRING ||
		bin.Right.Value.Typ == ast.VARIABLE_TYPE_STRING {
		maxstack, _ = m.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		stack, _ := m.build(class, code, bin.Right, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "compareTo",
			Descriptor: "(Ljava/lang/String;)I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := 1 + stack; t > maxstack {
			maxstack = t
		}
		state.popStack(1) // pop left string
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		if e.Typ == ast.EXPRESSION_TYPE_GT || e.Typ == ast.EXPRESSION_TYPE_LE {
			code.Codes[code.CodeLength] = cg.OP_ifgt
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		} else if e.Typ == ast.EXPRESSION_TYPE_LT || e.Typ == ast.EXPRESSION_TYPE_GE {
			code.Codes[code.CodeLength] = cg.OP_iflt
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		} else { //   == or !=
			code.Codes[code.CodeLength] = cg.OP_ifeq
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			} else {
				code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			}
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			} else {
				code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			}
			code.CodeLength += 8
		}
		return
	}

	if bin.Left.Value.IsPointer() && bin.Right.Value.IsPointer() { //
		stack, _ := m.build(class, code, bin.Left, context, state)
		if stack > maxstack {
			maxstack = stack
		}
		state.pushStack(class, bin.Left.Value)
		stack, _ = m.build(class, code, bin.Right, context, state)
		if t := stack + 1; t > maxstack {
			maxstack = t
		}
		state.popStack(1) // pop bin left
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_if_acmpeq
		} else { // ne
			code.Codes[code.CodeLength] = cg.OP_if_acmpne
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}

	// enum
	if bin.Left.Value.Typ == ast.VARIABLE_TYPE_ENUM {
		stack, _ := m.build(class, code, bin.Left, context, state)
		if stack > maxstack {
			maxstack = stack
		}
		state.pushStack(class, bin.Left.Value)
		stack, _ = m.build(class, code, bin.Right, context, state)
		if t := stack + 1; t > maxstack {
			maxstack = t
		}
		state.popStack(1) //
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_BOOL})
		context.MakeStackMap(code, state, code.CodeLength+8) //result on stack
		code.Codes[code.CodeLength] = cg.OP_if_icmpeq
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		} else {
			code.Codes[code.CodeLength+3] = cg.OP_iconst_1
		}
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength+7] = cg.OP_iconst_0
		}
		code.CodeLength += 8
		return
	}
	return
}
