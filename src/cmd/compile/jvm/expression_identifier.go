package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildCapturedIdentifier(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context) (maxStack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifier)
	if identifier.Variable.BeenCapturedAndModifiedInCaptureFunction {
		if context.function.Closure.ClosureVariableExist(identifier.Variable) {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			meta := closure.getMeta(identifier.Variable.Type.Type)
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      identifier.Variable.Name,
				Descriptor: "L" + meta.className + ";",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, identifier.Variable.LocalValOffset)...)
		}
		if 1 > maxStack {
			maxStack = 1
		}
		closure.unPack(class, code, identifier.Variable.Type)
		if t := jvmSlotSize(identifier.Variable.Type); t > maxStack {
			maxStack = t
		}
	} else {
		if identifier.Variable.JvmDescriptor == "" {
			identifier.Variable.JvmDescriptor = Descriptor.typeDescriptor(identifier.Variable.Type)
		}
		if context.function.Closure.ClosureVariableExist(identifier.Variable) {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      identifier.Variable.Name,
				Descriptor: identifier.Variable.JvmDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		} else {
			copyOPs(code, loadLocalVariableOps(identifier.Variable.Type.Type, identifier.Variable.LocalValOffset)...)
		}
		if t := jvmSlotSize(e.Value); t > maxStack {
			maxStack = t
		}
	}
	return
}

func (buildExpression *BuildExpression) buildIdentifier(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context) (maxStack uint16) {
	if e.Value.Type == ast.VariableTypeClass {
		panic("this is not happening")
	}
	identifier := e.Data.(*ast.ExpressionIdentifier)
	if e.Value.Type == ast.VariableTypeEnum && identifier.EnumName != nil {
		loadInt32(class, code, identifier.EnumName.Value)
		maxStack = 1
		return
	}
	if identifier.Function != nil {
		return buildExpression.BuildPackage.packFunction2MethodHandle(class, code, identifier.Function, context)
	}
	if identifier.Variable.IsGlobal { //fetch global var
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      buildExpression.BuildPackage.mainClass.Name,
			Field:      identifier.Name,
			Descriptor: Descriptor.typeDescriptor(identifier.Variable.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxStack = jvmSlotSize(identifier.Variable.Type)
		return
	}
	if identifier.Variable.BeenCaptured {
		return buildExpression.buildCapturedIdentifier(class, code, e, context)
	}
	maxStack = jvmSlotSize(e.Value)
	copyOPs(code, loadLocalVariableOps(e.Value.Type, identifier.Variable.LocalValOffset)...)
	return
}
