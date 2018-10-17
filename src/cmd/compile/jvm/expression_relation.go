package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildRelations(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	bin := e.Data.(*ast.ExpressionBinary)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	if bin.Left.Value.IsNumber() ||
		bin.Left.Value.Type == ast.VariableTypeEnum { // in this case ,right must be a number type
		maxStack = buildExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		stack := buildExpression.build(class, code, bin.Right, context, state)
		if t := jvmSlotSize(bin.Left.Value) + stack; t > maxStack {
			maxStack = t
		}
		if bin.Left.Value.Type == ast.VariableTypeByte ||
			bin.Left.Value.Type == ast.VariableTypeShort ||
			bin.Left.Value.Type == ast.VariableTypeChar ||
			bin.Left.Value.Type == ast.VariableTypeInt ||
			bin.Left.Value.Type == ast.VariableTypeEnum {
			switch e.Type {
			case ast.ExpressionTypeGt:
				code.Codes[code.CodeLength] = cg.OP_if_icmpgt
			case ast.ExpressionTypeLe:
				code.Codes[code.CodeLength] = cg.OP_if_icmple
			case ast.ExpressionTypeLt:
				code.Codes[code.CodeLength] = cg.OP_if_icmplt
			case ast.ExpressionTypeGe:
				code.Codes[code.CodeLength] = cg.OP_if_icmpge
			case ast.ExpressionTypeEq:
				code.Codes[code.CodeLength] = cg.OP_if_icmpeq
			case ast.ExpressionTypeNe:
				code.Codes[code.CodeLength] = cg.OP_if_icmpne
			}
			code.CodeLength++
		} else {
			switch bin.Left.Value.Type {
			case ast.VariableTypeLong:
				code.Codes[code.CodeLength] = cg.OP_lcmp
			case ast.VariableTypeFloat:
				code.Codes[code.CodeLength] = cg.OP_fcmpl
			case ast.VariableTypeDouble:
				code.Codes[code.CodeLength] = cg.OP_dcmpl
			}
			code.CodeLength++
			switch e.Type {
			case ast.ExpressionTypeGt:
				code.Codes[code.CodeLength] = cg.OP_ifgt
			case ast.ExpressionTypeLe:
				code.Codes[code.CodeLength] = cg.OP_ifle
			case ast.ExpressionTypeLt:
				code.Codes[code.CodeLength] = cg.OP_iflt
			case ast.ExpressionTypeGe:
				code.Codes[code.CodeLength] = cg.OP_ifge
			case ast.ExpressionTypeEq:
				code.Codes[code.CodeLength] = cg.OP_ifeq
			case ast.ExpressionTypeNe:
				code.Codes[code.CodeLength] = cg.OP_ifne
			}
			code.CodeLength++
		}
		state.popStack(1)
		context.MakeStackMap(code, state, code.CodeLength+6)
		state.pushStack(class, &ast.Type{
			Type: ast.VariableTypeBool,
		})
		context.MakeStackMap(code, state, code.CodeLength+7)
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength:code.CodeLength+2], 7)
		code.Codes[code.CodeLength+2] = cg.OP_iconst_0
		code.Codes[code.CodeLength+3] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+4:code.CodeLength+6], 4)
		code.Codes[code.CodeLength+6] = cg.OP_iconst_1
		code.CodeLength += 7
		return
	}
	if bin.Left.Value.Type == ast.VariableTypeBool ||
		bin.Right.Value.Type == ast.VariableTypeBool { // bool type
		var dependOnOne *ast.Expression
		if bin.Left.IsBoolLiteral(true) && e.Type == ast.ExpressionTypeEq {
			/*
				a == true <==> a
			*/
			dependOnOne = bin.Right
		} else if bin.Left.IsBoolLiteral(false) && e.Type == ast.ExpressionTypeNe {
			dependOnOne = bin.Right
		} else if bin.Right.IsBoolLiteral(true) && e.Type == ast.ExpressionTypeEq {
			dependOnOne = bin.Left
		} else if bin.Right.IsBoolLiteral(false) && e.Type == ast.ExpressionTypeNe {
			dependOnOne = bin.Left
		}
		if dependOnOne != nil {
			return buildExpression.build(class, code, dependOnOne, context, state)
		} else {
			maxStack = buildExpression.build(class, code, bin.Left, context, state)
			state.pushStack(class, bin.Left.Value)
			stack := buildExpression.build(class, code, bin.Right, context, state)
			state.pushStack(class, bin.Right.Value)
			if t := jvmSlotSize(bin.Left.Value) + stack; t > maxStack {
				maxStack = t
			}
			state.popStack(2) // 2 bool value
			context.MakeStackMap(code, state, code.CodeLength+7)
			state.pushStack(class, &ast.Type{
				Type: ast.VariableTypeBool,
			})
			context.MakeStackMap(code, state, code.CodeLength+8)
			if e.Type == ast.ExpressionTypeEq {
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
	}
	if bin.Left.Value.Type == ast.VariableTypeNull ||
		bin.Right.Value.Type == ast.VariableTypeNull { // must not null-null
		var notNullExpression *ast.Expression
		if bin.Left.Value.Type != ast.VariableTypeNull {
			notNullExpression = bin.Left
		} else {
			notNullExpression = bin.Right
		}
		maxStack = buildExpression.build(class, code, notNullExpression, context, state)
		if e.Type == ast.ExpressionTypeEq {
			code.Codes[code.CodeLength] = cg.OP_ifnull
		} else { // ne
			code.Codes[code.CodeLength] = cg.OP_ifnonnull
		}
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.Type{
			Type: ast.VariableTypeBool,
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
	if bin.Left.Value.Type == ast.VariableTypeString {
		maxStack = buildExpression.build(class, code, bin.Left, context, state)
		state.pushStack(class, bin.Left.Value)
		stack := buildExpression.build(class, code, bin.Right, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaStringClass,
			Method:     "compareTo",
			Descriptor: "(Ljava/lang/String;)I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		state.popStack(1) // pop left string
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.Type{
			Type: ast.VariableTypeBool,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		switch e.Type {
		case ast.ExpressionTypeGt:
			code.Codes[code.CodeLength] = cg.OP_ifgt
		case ast.ExpressionTypeLe:
			code.Codes[code.CodeLength] = cg.OP_ifle
		case ast.ExpressionTypeLt:
			code.Codes[code.CodeLength] = cg.OP_iflt
		case ast.ExpressionTypeGe:
			code.Codes[code.CodeLength] = cg.OP_ifge
		case ast.ExpressionTypeEq:
			code.Codes[code.CodeLength] = cg.OP_ifeq
		case ast.ExpressionTypeNe:
			code.Codes[code.CodeLength] = cg.OP_ifne
		}
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
		code.Codes[code.CodeLength+3] = cg.OP_iconst_0
		code.Codes[code.CodeLength+4] = cg.OP_goto
		binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
		code.Codes[code.CodeLength+7] = cg.OP_iconst_1
		code.CodeLength += 8
		return
	}

	if bin.Left.Value.IsPointer() && bin.Right.Value.IsPointer() { //
		stack := buildExpression.build(class, code, bin.Left, context, state)
		if stack > maxStack {
			maxStack = stack
		}
		state.pushStack(class, bin.Left.Value)
		stack = buildExpression.build(class, code, bin.Right, context, state)
		if t := stack + 1; t > maxStack {
			maxStack = t
		}
		state.popStack(1) // pop bin left
		context.MakeStackMap(code, state, code.CodeLength+7)
		state.pushStack(class, &ast.Type{
			Type: ast.VariableTypeBool,
		})
		context.MakeStackMap(code, state, code.CodeLength+8)
		if e.Type == ast.ExpressionTypeEq {
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
	return
}
