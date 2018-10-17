package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildExpression *BuildExpression) buildCapturedIdentifier(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context) (maxStack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifier)
	//fmt.Println(identifier.Name, identifier.Variable.BeenCapturedAsLeftValue, context.function.Closure.ClosureVariableExist(identifier.Variable))
	maxStack = jvmSlotSize(identifier.Variable.Type)
	if context.function.Closure.ClosureVariableExist(identifier.Variable) {
		if identifier.Variable.BeenCapturedAsLeftValue > 0 {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			meta := closure.getMeta(identifier.Variable.Type.Type)
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      identifier.Variable.Name,
				Descriptor: "L" + meta.className + ";",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
			closure.unPack(class, code, identifier.Variable.Type)
		} else {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
			code.Codes[code.CodeLength] = cg.OP_getfield
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      class.Name,
				Field:      identifier.Variable.Name,
				Descriptor: Descriptor.typeDescriptor(identifier.Variable.Type),
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
	} else {
		if identifier.Variable.BeenCapturedAsLeftValue > 0 {
			copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, identifier.Variable.LocalValOffset)...)
			closure.unPack(class, code, identifier.Variable.Type)
		} else {
			copyOPs(code, loadLocalVariableOps(e.Value.Type, identifier.Variable.LocalValOffset)...)
		}
	}

	return
}

func (buildExpression *BuildExpression) buildIdentifier(
	class *cg.ClassHighLevel,
	code *cg.AttributeCode,
	e *ast.Expression,
	context *Context) (maxStack uint16) {
	if e.Value.Type == ast.VariableTypeClass {
		panic("this is not happening")
	}
	identifier := e.Data.(*ast.ExpressionIdentifier)

	switch {
	case e.Value.Type == ast.VariableTypeEnum && identifier.EnumName != nil:
		loadInt32(class, code, identifier.EnumName.Value)
		maxStack = 1
		return
	case identifier.Function != nil:
		return buildExpression.BuildPackage.packFunction2MethodHandle(class, code, identifier.Function, context)
	case identifier.Variable.IsGlobal:
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      buildExpression.BuildPackage.mainClass.Name,
			Field:      identifier.Name,
			Descriptor: Descriptor.typeDescriptor(identifier.Variable.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxStack = jvmSlotSize(identifier.Variable.Type)
		return
	case identifier.Variable.BeenCapturedAsLeftValue+identifier.Variable.BeenCapturedAsRightValue > 0:
		return buildExpression.buildCapturedIdentifier(class, code, e, context)
	default:
		maxStack = jvmSlotSize(e.Value)
		copyOPs(code, loadLocalVariableOps(e.Value.Type, identifier.Variable.LocalValOffset)...)
		return
	}
}
