package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildMapMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	maxStack = buildExpression.build(class, code, call.Expression, context, state)
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	hashMapVerifyType := state.newObjectVariableType(javaMapClass)
	state.pushStack(class, hashMapVerifyType)
	switch call.Name {
	case common.MapMethodKeyExist:
		variableType := call.Args[0].Value
		stack := buildExpression.build(class, code, call.Args[0], context, state)
		if t := 1 + stack; t > maxStack {
			maxStack = t
		}
		if variableType.IsPointer() == false {
			typeConverter.packPrimitives(class, code, variableType)
		}
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaMapClass,
			Method:     "containsKey",
			Descriptor: "(Ljava/lang/Object;)Z",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	case common.MapMethodRemove:
		currentStack := uint16(1)
		callRemove := func() {
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      javaMapClass,
				Method:     "remove",
				Descriptor: "(Ljava/lang/Object;)Ljava/lang/Object;",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.Codes[code.CodeLength+3] = cg.OP_pop
			code.CodeLength += 4
		}
		for k, v := range call.Args {
			currentStack = 1
			if v.HaveMultiValue() {
				stack := buildExpression.build(class, code, v, context, state)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				autoVar := newMultiValueAutoVar(class, code, state) // store to temp
				for kk, tt := range v.MultiValues {
					currentStack = 1
					if k != len(call.Args)-1 || kk != len(v.MultiValues)-1 {
						code.Codes[code.CodeLength] = cg.OP_dup
						code.CodeLength++
						currentStack++
						state.pushStack(class, hashMapVerifyType)
					}
					//load
					stack = autoVar.unPack(class, code, kk, tt)
					if t := stack + currentStack; t > maxStack {
						maxStack = t
					}
					//remove
					callRemove()
					if k != len(call.Args)-1 || kk != len(v.MultiValues)-1 {
						state.popStack(1)
					}
				}
				continue
			}
			variableType := v.Value
			if k != len(call.Args)-1 {
				code.Codes[code.CodeLength] = cg.OP_dup
				currentStack++
				if currentStack > maxStack {
					maxStack = currentStack
				}
				state.pushStack(class, hashMapVerifyType)
			}
			stack := buildExpression.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxStack {
				maxStack = t
			}
			if variableType.IsPointer() == false {
				typeConverter.packPrimitives(class, code, variableType)
			}
			//call remove
			callRemove()
			if k != len(call.Args)-1 {
				state.popStack(1)
			}
		}
	case common.MapMethodRemoveAll:
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaMapClass,
			Method:     "clear",
			Descriptor: "()V",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case common.MapMethodSize:
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      javaMapClass,
			Method:     "size",
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}
