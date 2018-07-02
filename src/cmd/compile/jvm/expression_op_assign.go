package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	s += "456";
*/
func (buildExpression *BuildExpression) buildStrPlusAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	bin := e.Data.(*ast.ExpressionBinary)
	maxStack, remainStack, op, _, leftValueKind := buildExpression.getLeftValue(class, code, bin.Left, context, state)
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst("java/lang/StringBuilder", code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if t := remainStack + 2; t > maxStack {
		maxStack = t
	}
	state.pushStack(class, state.newObjectVariableType(javaStringBuilderClass))
	currentStack := remainStack + 1 //
	stack, _ := buildExpression.build(class, code, bin.Left, context, state)
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	//append origin string
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     `append`,
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	stack, _ = buildExpression.build(class, code, bin.Right, context, state)
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	//append right
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     `append`,
		Descriptor: "(Ljava/lang/String;)Ljava/lang/StringBuilder;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	// tostring
	code.Codes[code.CodeLength] = cg.OP_invokevirtual
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      javaStringBuilderClass,
		Method:     `toString`,
		Descriptor: "()Ljava/lang/String;",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if e.IsStatementExpression == false {
		currentStack += buildExpression.controlStack2FitAssign(code, leftValueKind, bin.Left.ExpressionValue)
	}
	//copy op
	copyOPs(code, op...)
	return

}
func (buildExpression *BuildExpression) buildOpAssign(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	bin := e.Data.(*ast.ExpressionBinary)
	if bin.Left.ExpressionValue.Type == ast.VariableTypeString {
		return buildExpression.buildStrPlusAssign(class, code, e, context, state)
	}
	maxStack, remainStack, op, _, leftValueKind := buildExpression.getLeftValue(class, code, bin.Left, context, state)
	//left value must can be used as right value,
	stack, _ := buildExpression.build(class, code, bin.Left, context, state) // load it`s value
	if t := stack + remainStack; t > maxStack {
		maxStack = t
	}
	state.pushStack(class, e.ExpressionValue)
	currentStack := jvmSlotSize(e.ExpressionValue) + remainStack // incase int -> long
	stack, _ = buildExpression.build(class, code, bin.Right, context, state)
	if t := currentStack + stack; t > maxStack {
		maxStack = t
	}
	switch bin.Left.ExpressionValue.Type {
	case ast.VariableTypeByte:
		if e.Type == ast.ExpressionTypePlusAssign {
			code.Codes[code.CodeLength] = cg.OP_iadd
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeMinusAssign {
			code.Codes[code.CodeLength] = cg.OP_isub
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeMulAssign {
			code.Codes[code.CodeLength] = cg.OP_imul
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeDivAssign {
			code.Codes[code.CodeLength] = cg.OP_idiv
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeModAssign {
			code.Codes[code.CodeLength] = cg.OP_irem
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeAndAssign {
			code.Codes[code.CodeLength] = cg.OP_iand
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeOrAssign {
			code.Codes[code.CodeLength] = cg.OP_ior
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeLshAssign {
			code.Codes[code.CodeLength] = cg.OP_ishl
			code.Codes[code.CodeLength+1] = cg.OP_i2b
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeRshAssign {
			code.Codes[code.CodeLength] = cg.OP_ishr
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeXorAssign {
			code.Codes[code.CodeLength] = cg.OP_ixor
			code.CodeLength++
		}
	case ast.VariableTypeShort:
		if e.Type == ast.ExpressionTypePlusAssign {
			code.Codes[code.CodeLength] = cg.OP_iadd
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeMinusAssign {
			code.Codes[code.CodeLength] = cg.OP_isub
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeMulAssign {
			code.Codes[code.CodeLength] = cg.OP_imul
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeDivAssign {
			code.Codes[code.CodeLength] = cg.OP_idiv
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeModAssign {
			code.Codes[code.CodeLength] = cg.OP_irem
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeAndAssign {
			code.Codes[code.CodeLength] = cg.OP_iand
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeOrAssign {
			code.Codes[code.CodeLength] = cg.OP_ior
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeLshAssign {
			code.Codes[code.CodeLength] = cg.OP_ishl
			code.Codes[code.CodeLength+1] = cg.OP_i2s
			code.CodeLength += 2
		} else if e.Type == ast.ExpressionTypeRshAssign {
			code.Codes[code.CodeLength] = cg.OP_ishr
			code.CodeLength++
		} else if e.Type == ast.ExpressionTypeXorAssign {
			code.Codes[code.CodeLength] = cg.OP_ixor
			code.CodeLength++
		}

	case ast.VariableTypeInt:
		if e.Type == ast.ExpressionTypePlusAssign {
			code.Codes[code.CodeLength] = cg.OP_iadd
		} else if e.Type == ast.ExpressionTypeMinusAssign {
			code.Codes[code.CodeLength] = cg.OP_isub
		} else if e.Type == ast.ExpressionTypeMulAssign {
			code.Codes[code.CodeLength] = cg.OP_imul
		} else if e.Type == ast.ExpressionTypeDivAssign {
			code.Codes[code.CodeLength] = cg.OP_idiv
		} else if e.Type == ast.ExpressionTypeModAssign {
			code.Codes[code.CodeLength] = cg.OP_irem
		} else if e.Type == ast.ExpressionTypeAndAssign {
			code.Codes[code.CodeLength] = cg.OP_iand
		} else if e.Type == ast.ExpressionTypeOrAssign {
			code.Codes[code.CodeLength] = cg.OP_ior
		} else if e.Type == ast.ExpressionTypeLshAssign {
			code.Codes[code.CodeLength] = cg.OP_ishl
		} else if e.Type == ast.ExpressionTypeRshAssign {
			code.Codes[code.CodeLength] = cg.OP_ishr
		} else if e.Type == ast.ExpressionTypeXorAssign {
			code.Codes[code.CodeLength] = cg.OP_ixor
		}
		code.CodeLength++
	case ast.VariableTypeLong:
		if e.Type == ast.ExpressionTypePlusAssign {
			code.Codes[code.CodeLength] = cg.OP_ladd
		} else if e.Type == ast.ExpressionTypeMinusAssign {
			code.Codes[code.CodeLength] = cg.OP_lsub
		} else if e.Type == ast.ExpressionTypeMulAssign {
			code.Codes[code.CodeLength] = cg.OP_lmul
		} else if e.Type == ast.ExpressionTypeDivAssign {
			code.Codes[code.CodeLength] = cg.OP_ldiv
		} else if e.Type == ast.ExpressionTypeModAssign {
			code.Codes[code.CodeLength] = cg.OP_lrem
		} else if e.Type == ast.ExpressionTypeAndAssign {
			code.Codes[code.CodeLength] = cg.OP_land
		} else if e.Type == ast.ExpressionTypeOrAssign {
			code.Codes[code.CodeLength] = cg.OP_lor
		} else if e.Type == ast.ExpressionTypeLshAssign {
			code.Codes[code.CodeLength] = cg.OP_lshl
		} else if e.Type == ast.ExpressionTypeRshAssign {
			code.Codes[code.CodeLength] = cg.OP_lshr
		} else if e.Type == ast.ExpressionTypeXorAssign {
			code.Codes[code.CodeLength] = cg.OP_lxor
		}
		code.CodeLength++
	case ast.VariableTypeFloat:
		if e.Type == ast.ExpressionTypePlusAssign {
			code.Codes[code.CodeLength] = cg.OP_ladd
		} else if e.Type == ast.ExpressionTypeMinusAssign {
			code.Codes[code.CodeLength] = cg.OP_lsub
		} else if e.Type == ast.ExpressionTypeMulAssign {
			code.Codes[code.CodeLength] = cg.OP_lmul
		} else if e.Type == ast.ExpressionTypeDivAssign {
			code.Codes[code.CodeLength] = cg.OP_ldiv
		} else if e.Type == ast.ExpressionTypeModAssign {
			code.Codes[code.CodeLength] = cg.OP_frem
		}
		code.CodeLength++
	case ast.VariableTypeDouble:
		if e.Type == ast.ExpressionTypePlusAssign {
			code.Codes[code.CodeLength] = cg.OP_dadd

		} else if e.Type == ast.ExpressionTypeMinusAssign {
			code.Codes[code.CodeLength] = cg.OP_dsub

		} else if e.Type == ast.ExpressionTypeMulAssign {
			code.Codes[code.CodeLength] = cg.OP_dmul

		} else if e.Type == ast.ExpressionTypeDivAssign {
			code.Codes[code.CodeLength] = cg.OP_ddiv

		} else if e.Type == ast.ExpressionTypeModAssign {
			code.Codes[code.CodeLength] = cg.OP_drem

		}
		code.CodeLength++
	}
	if e.IsStatementExpression == false {
		currentStack += buildExpression.controlStack2FitAssign(code, leftValueKind, bin.Left.ExpressionValue)
		if currentStack > maxStack {
			maxStack = currentStack
		}
	}
	//copy op
	copyOPs(code, op...)
	return
}
