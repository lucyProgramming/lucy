package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildSelfIncrement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	ee := e.Data.(*ast.Expression)
	// identifier  and not captured and type`s int
	if identifier, ok := ee.Data.(*ast.ExpressionIdentifier); ok && ee.Type == ast.ExpressionTypeIdentifier &&
		identifier.Variable.BeenCaptured == 0 &&
		identifier.Variable.Type.Type == ast.VariableTypeInt &&
		identifier.Variable.IsGlobal == false {
		if identifier.Variable.LocalValOffset > 255 { // early check
			panic("over 255")
		}
		if e.IsStatementExpression == false { //  need it`s value
			if e.Type == ast.ExpressionTypeIncrement || e.Type == ast.ExpressionTypeDecrement {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, identifier.Variable.LocalValOffset)...) // load to stack top
				maxStack = 1
			}
		}
		code.Codes[code.CodeLength] = cg.OP_iinc
		code.Codes[code.CodeLength+1] = byte(identifier.Variable.LocalValOffset)
		if e.Type == ast.ExpressionTypePrefixIncrement || e.Type == ast.ExpressionTypeIncrement {
			code.Codes[code.CodeLength+2] = 1
		} else { // --
			code.Codes[code.CodeLength+2] = 255 // -1
			code.CodeLength += 3
		}
		code.CodeLength += 3
		if e.IsStatementExpression == false { // I still need it`s value
			if e.Type == ast.ExpressionTypePrefixIncrement || e.Type == ast.ExpressionTypePrefixDecrement { // decrement
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, identifier.Variable.LocalValOffset)...) // load to stack top
				maxStack = 1
			}
		}
		return
	}
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	maxStack, remainStack, op, leftValueKind :=
		buildExpression.getLeftValue(class, code, ee, context, state)
	/*
		left value must can be used as right value
	*/
	stack := buildExpression.build(class, code, ee, context, state) // load it`s value
	if t := stack + remainStack; t > maxStack {
		maxStack = t
	}
	currentStack := jvmSlotSize(e.Value) + remainStack
	if e.IsStatementExpression == false {
		if e.Type == ast.ExpressionTypeIncrement || e.Type == ast.ExpressionTypeDecrement {
			currentStack += buildExpression.controlStack2FitAssign(code, leftValueKind, e.Value)
			if currentStack > maxStack {
				maxStack = currentStack
			}
		}
	}
	if t := currentStack + jvmSlotSize(e.Value); t > maxStack {
		//
		maxStack = t
	}
	switch e.Value.Type {
	case ast.VariableTypeByte:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.Codes[code.CodeLength+2] = cg.OP_i2b
		code.CodeLength += 3
	case ast.VariableTypeShort:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.Codes[code.CodeLength+2] = cg.OP_i2s
		code.CodeLength += 3
	case ast.VariableTypeInt:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.CodeLength += 2
	case ast.VariableTypeLong:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_lconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc2_w
			class.InsertLongConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		code.Codes[code.CodeLength] = cg.OP_ladd
		code.CodeLength++
	case ast.VariableTypeFloat:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_fconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertFloatConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		code.Codes[code.CodeLength] = cg.OP_fadd
		code.CodeLength++
	case ast.VariableTypeDouble:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_dconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc2_w
			class.InsertDoubleConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		code.Codes[code.CodeLength] = cg.OP_dadd
		code.CodeLength++
	}
	if e.IsStatementExpression == false {
		if e.Type == ast.ExpressionTypePrefixIncrement ||
			e.Type == ast.ExpressionTypePrefixDecrement {
			currentStack += buildExpression.controlStack2FitAssign(code, leftValueKind, e.Value)
			if currentStack > maxStack {
				maxStack = currentStack
			}
		}
	}
	//copy op
	copyOPs(code, op...)
	return
}
