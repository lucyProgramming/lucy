package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildRelations(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if bin.Left.Value.IsNumber() { // in this case ,right must be a number type
		maxStack, _ = m.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		stack, _ := m.build(class, code, bin.Right, context, state)
		if t := jvmSize(bin.Left.Value) + stack; t > maxStack {
			maxStack = t
		}
		switch bin.Left.Value.Typ {
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
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength] = cg.OP_ifgt
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifle
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		} else if e.Typ == ast.EXPRESSION_TYPE_LT || e.Typ == ast.EXPRESSION_TYPE_GE { // < and >=
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength] = cg.OP_iflt
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifge
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		} else {
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength] = cg.OP_ifeq
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifne
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		}
		return
	}
	if bin.Left.Value.Typ == ast.VARIABLE_TYPE_BOOL ||
		bin.Right.Value.Typ == ast.VARIABLE_TYPE_BOOL { // bool type
		var es []*cg.Exit
		maxStack, es = m.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		if len(es) > 0 {
			context.MakeStackMap(code, state, code.CodeLength)
			backfillExit(es, code.CodeLength)
		}
		stack, es := m.build(class, code, bin.Right, context, state)
		state.pushStack(class, bin.Right.Value)
		if len(es) > 0 {
			context.MakeStackMap(code, state, code.CodeLength)
			backfillExit(es, code.CodeLength)
		}
		if t := jvmSize(bin.Left.Value) + stack; t > maxStack {
			maxStack = t
		}
		state.popStack(2) // 2 bool value
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_if_icmpeq
		} else {
			code.Codes[code.CodeLength] = cg.OP_if_icmpne
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
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
		maxStack, _ = m.build(class, code, notNullExpression, context, state)
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
		maxStack, _ = m.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		stack, _ := m.build(class, code, bin.Right, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      java_string_class,
			Method:     "compareTo",
			Descriptor: "(Ljava/lang/String;)I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		state.popStack(1) // pop left string
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		if e.Typ == ast.EXPRESSION_TYPE_GT || e.Typ == ast.EXPRESSION_TYPE_LE { // > and <=
			if e.Typ == ast.EXPRESSION_TYPE_GT {
				code.Codes[code.CodeLength] = cg.OP_ifgt
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifle
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		} else if e.Typ == ast.EXPRESSION_TYPE_LT || e.Typ == ast.EXPRESSION_TYPE_GE { // < and >=
			if e.Typ == ast.EXPRESSION_TYPE_LT {
				code.Codes[code.CodeLength] = cg.OP_iflt
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifge
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		} else {
			if e.Typ == ast.EXPRESSION_TYPE_EQ {
				code.Codes[code.CodeLength] = cg.OP_ifeq
			} else {
				code.Codes[code.CodeLength] = cg.OP_ifne
			}
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_0
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_1
			code.CodeLength += 8
		}
		return
	}

	if bin.Left.Value.IsPointer() && bin.Right.Value.IsPointer() { //
		stack, _ := m.build(class, code, bin.Left, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		state.pushStack(class, bin.Left.Value)
		stack, _ = m.build(class, code, bin.Right, context, state)
		if t := stack + 1; t > maxStack {
			maxStack = t
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
		if stack > maxStack {
			maxStack = stack
		}
		state.pushStack(class, bin.Left.Value)
		stack, _ = m.build(class, code, bin.Right, context, state)
		if t := stack + jvmSize(bin.Left.Value); t > maxStack {
			maxStack = t
		}
		state.popStack(1) //
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.VariableType{Typ: ast.VARIABLE_TYPE_BOOL})
		context.MakeStackMap(code, state, code.CodeLength+8) //result on stack
		if e.Typ == ast.EXPRESSION_TYPE_EQ {
			code.Codes[code.CodeLength] = cg.OP_if_icmpeq
		} else {
			code.Codes[code.CodeLength] = cg.OP_if_icmpne
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}
	return
}
