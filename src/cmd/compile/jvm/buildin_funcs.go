package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *BuildExpression) mkBuildInFunctionCall(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context,
	state *StackMapState) (maxStack uint16) {
	call := e.Data.(*ast.ExpressionFunctionCall)
	if call.Function.LoadedFromCorePackage {
		maxStack = this.buildCallArgs(class, code, call.Args, call.VArgs, context, state)
		code.Codes[code.CodeLength] = cg.OP_invokestatic
		class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      call.Function.Entrance.Class.Name,
			Method:     call.Function.Name,
			Descriptor: call.Function.Entrance.Descriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.IsStatementExpression {
			if call.Function.Type.VoidReturn() == false {
				if len(call.Function.Type.ReturnList) > 1 {
					code.Codes[code.CodeLength] = cg.OP_pop
					code.CodeLength++
				} else {
					if jvmSlotSize(e.Value) == 1 {
						code.Codes[code.CodeLength] = cg.OP_pop
						code.CodeLength++
					} else {
						code.Codes[code.CodeLength] = cg.OP_pop2
						code.CodeLength++
					}
				}
			}
		}

		return
	}
	switch call.Function.Name {
	case common.BuildInFunctionPrint:
		return this.mkBuildInPrint(class, code, e, context, state)
	case common.BuildInFunctionPanic:
		return this.mkBuildInPanic(class, code, e, context, state)
	case common.BuildInFunctionCatch:
		return this.mkBuildInCatch(class, code, e, context)
	case common.BuildInFunctionMonitorEnter, common.BuildInFunctionMonitorExit:
		maxStack = this.build(class, code, call.Args[0], context, state)
		if call.Function.Name == common.BuildInFunctionMonitorEnter {
			code.Codes[code.CodeLength] = cg.OP_monitorenter
		} else { // monitor enter on exit
			code.Codes[code.CodeLength] = cg.OP_monitorexit
		}
		code.CodeLength++
	case common.BuildInFunctionPrintf:
		return this.mkBuildInPrintf(class, code, e, context, state)
	case common.BuildInFunctionSprintf:
		return this.mkBuildInSprintf(class, code, e, context, state)
	case common.BuildInFunctionLen:
		return this.mkBuildInLen(class, code, e, context, state)
	case common.BuildInFunctionBlockHole:
		return this.mkBuildInBlackHole(class, code, e, context, state)
	case common.BuildInFunctionAssert:
		return this.mkBuildInAssert(class, code, e, context, state)
	default:
		panic("unknown  buildIn function:" + call.Function.Name)
	}
	return
}
