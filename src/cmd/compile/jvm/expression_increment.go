package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildSelfIncrement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	ee := e.Data.(*ast.Expression)
	// identifer  and not captured and type`s int
	if t, ok := ee.Data.(*ast.ExpressionIdentifer); ee.Typ == ast.EXPRESSION_TYPE_IDENTIFIER &&
		ok &&
		t.Var.BeenCaptured == false &&
		t.Var.Typ.Typ == ast.VARIABLE_TYPE_INT {
		if t.Var.LocalValOffset > 255 { // early check
			panic("over 255")
		}
		if e.IsStatementExpression == false { // I still need it`s value
			if e.Typ == ast.EXPRESSION_TYPE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_DECREMENT {
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, t.Var.LocalValOffset)...) // load to stack top
				maxstack = 1
			}
		}
		if e.Typ == ast.EXPRESSION_TYPE_PRE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_INCREMENT {
			code.Codes[code.CodeLength] = cg.OP_iinc
			code.Codes[code.CodeLength+1] = byte(t.Var.LocalValOffset)
			code.Codes[code.CodeLength+2] = 1
			code.CodeLength += 3
		} else { // --
			code.Codes[code.CodeLength] = cg.OP_iinc
			code.Codes[code.CodeLength+1] = byte(t.Var.LocalValOffset)
			code.Codes[code.CodeLength+2] = 255 // -1
			code.CodeLength += 3
		}
		if e.IsStatementExpression == false { // I still need it`s value
			if e.Typ == ast.EXPRESSION_TYPE_PRE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_PRE_DECREMENT { // decrement
				copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_INT, t.Var.LocalValOffset)...) // load to stack top
				maxstack = 1
			}
		}
		return
	}
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length)
	}()
	maxstack, remainStack, op, _, classname, name, descriptor := m.getLeftValue(class, code, ee, context, state)
	/*
		left value must can be used as right value
	*/
	stack, _ := m.build(class, code, ee, context, state) // load it`s value
	if t := stack + remainStack; t > maxstack {
		maxstack = t
	}
	currentStack := jvmSize(ee.Value) + remainStack
	if e.IsStatementExpression == false {
		if e.Typ == ast.EXPRESSION_TYPE_INCREMENT || e.Typ == ast.EXPRESSION_TYPE_DECREMENT {
			currentStack += m.controlStack2FitAssign(code, op, classname, e.Value)
			if currentStack > maxstack {
				maxstack = currentStack
			}
		}
	}
	switch e.Value.Typ {
	case ast.VARIABLE_TYPE_BYTE:
		if e.IsSelfIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.Codes[code.CodeLength+2] = cg.OP_i2b
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_SHORT:
		if e.IsSelfIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.Codes[code.CodeLength+2] = cg.OP_i2s
		code.CodeLength += 3
	case ast.VARIABLE_TYPE_INT:
		if e.IsSelfIncrement() {
			code.Codes[code.CodeLength] = cg.OP_iconst_1
		} else {
			code.Codes[code.CodeLength] = cg.OP_iconst_m1
		}
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength+1] = cg.OP_iadd
		code.CodeLength += 2
	case ast.VARIABLE_TYPE_LONG:
		if e.IsSelfIncrement() {
			code.Codes[code.CodeLength] = cg.OP_lconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc2_w
			class.InsertLongConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if t := currentStack + 2; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength] = cg.OP_ladd
		code.CodeLength++
	case ast.VARIABLE_TYPE_FLOAT:
		if e.IsSelfIncrement() {
			code.Codes[code.CodeLength] = cg.OP_fconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc_w
			class.InsertFloatConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if t := currentStack + 1; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength] = cg.OP_fadd
		code.CodeLength++
	case ast.VARIABLE_TYPE_DOUBLE:
		if e.IsSelfIncrement() {
			code.Codes[code.CodeLength] = cg.OP_dconst_1
			code.CodeLength++
		} else {
			code.Codes[code.CodeLength] = cg.OP_ldc2_w
			class.InsertDoubleConst(-1, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if t := currentStack + 2; t > maxstack { // last op will change stack
			maxstack = t
		}
		code.Codes[code.CodeLength] = cg.OP_dadd
		code.CodeLength++
	}
	if e.IsStatementExpression == false {
		if e.Typ == ast.EXPRESSION_TYPE_PRE_INCREMENT ||
			e.Typ == ast.EXPRESSION_TYPE_PRE_DECREMENT {
			currentStack += m.controlStack2FitAssign(code, op, classname, e.Value)
			if currentStack > maxstack {
				maxstack = currentStack
			}
		}
	}
	if classname == java_hashmap_class && e.Value.IsPointer() == false { // map detination
		typeConverter.putPrimitiveInObject(class, code, e.Value)
	}
	//copy op
	copyOPLeftValue(class, code, op, classname, name, descriptor)
	return
}
