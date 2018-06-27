package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildSelfIncrement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	ee := e.Data.(*ast.Expression)
	// identifier  and not captured and type`s int
	if t, ok := ee.Data.(*ast.ExpressionIdentifier); ee.Type == ast.ExpressionTypeIdentifier &&
		ok &&
		t.Variable.BeenCaptured == false &&
		t.Variable.Type.Type == ast.VariableTypeInt {
		if t.Variable.LocalValOffset > 255 { // early check
			panic("over 255")
		}
		if e.IsStatementExpression == false { //  need it`s value
			if e.Type == ast.ExpressionTypeIncrement || e.Type == ast.ExpressionTypeDecrement {
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, t.Variable.LocalValOffset)...) // load to stack top
				maxStack = 1
			}
		}
		if e.Type == ast.ExpressionTypePrefixIncrement || e.Type == ast.ExpressionTypeIncrement {
			code.Codes[code.CodeLength] = cg.OP_iinc
			code.Codes[code.CodeLength+1] = byte(t.Variable.LocalValOffset)
			code.Codes[code.CodeLength+2] = 1
			code.CodeLength += 3
		} else { // --
			code.Codes[code.CodeLength] = cg.OP_iinc
			code.Codes[code.CodeLength+1] = byte(t.Variable.LocalValOffset)
			code.Codes[code.CodeLength+2] = 255 // -1
			code.CodeLength += 3
		}
		if e.IsStatementExpression == false { // I still need it`s value
			if e.Type == ast.ExpressionTypePrefixIncrement || e.Type == ast.ExpressionTypePrefixDecrement { // decrement
				copyOPs(code, loadLocalVariableOps(ast.VariableTypeInt, t.Variable.LocalValOffset)...) // load to stack top
				maxStack = 1
			}
		}
		return
	}
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	maxStack, remainStack, op, _, className, name, descriptor := makeExpression.getLeftValue(class, code, ee, context, state)
	/*
		left value must can be used as right value
	*/
	stack, _ := makeExpression.build(class, code, ee, context, state) // load it`s value
	if t := stack + remainStack; t > maxStack {
		maxStack = t
	}
	currentStack := jvmSlotSize(ee.ExpressionValue) + remainStack
	if e.IsStatementExpression == false {
		if e.Type == ast.ExpressionTypeIncrement || e.Type == ast.ExpressionTypeDecrement {
			currentStack += makeExpression.controlStack2FitAssign(code, op, className, e.ExpressionValue)
			if currentStack > maxStack {
				maxStack = currentStack
			}
		}
	}
	switch e.ExpressionValue.Type {
	case ast.VariableTypeByte:
		if e.IsIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		if t := currentStack + 1; t > maxStack { // last op will change stack
			maxStack = t
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
		if t := currentStack + 1; t > maxStack { // last op will change stack
			maxStack = t
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
		if t := currentStack + 1; t > maxStack { // last op will change stack
			maxStack = t
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
		if t := currentStack + 2; t > maxStack { // last op will change stack
			maxStack = t
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
		if t := currentStack + 1; t > maxStack { // last op will change stack
			maxStack = t
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
		if t := currentStack + 2; t > maxStack { // last op will change stack
			maxStack = t
		}
		code.Codes[code.CodeLength] = cg.OP_dadd
		code.CodeLength++
	}
	if e.IsStatementExpression == false {
		if e.Type == ast.ExpressionTypePrefixIncrement ||
			e.Type == ast.ExpressionTypePrefixDecrement {
			currentStack += makeExpression.controlStack2FitAssign(code, op, className, e.ExpressionValue)
			if currentStack > maxStack {
				maxStack = currentStack
			}
		}
	}
	//copy op
	copyLeftValueOps(class, code, op, className, name, descriptor)
	return
}
