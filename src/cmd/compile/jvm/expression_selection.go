package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

func (makeExpression *MakeExpression) buildSelection(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	selection := e.Data.(*ast.ExpressionSelection)
	if selection.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_PACKAGE {
		maxStack = jvmSize(e.ExpressionValue)
		if selection.PackageVariable != nil {
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      selection.Expression.ExpressionValue.Package.Name + "/main",
				Field:      selection.PackageVariable.Name,
				Descriptor: selection.PackageVariable.JvmDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if selection.PackageEnumName != nil {
			loadInt(class, code, selection.PackageEnumName.Value)
		}
		return
	}
	// check cast to super class
	if selection.Name == ast.SUPER_FIELD_NAME {
		if selection.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_OBJECT {
			maxStack, _ = makeExpression.build(class, code, selection.Expression, context, nil)
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst(e.ExpressionValue.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}
	if selection.Expression.ExpressionValue.Type == ast.VARIABLE_TYPE_CLASS {
		maxStack = jvmSize(e.ExpressionValue)
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      selection.Expression.ExpressionValue.Class.Name,
			Field:      selection.Name,
			Descriptor: Descriptor.typeDescriptor(e.ExpressionValue),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	// object
	maxStack, _ = makeExpression.build(class, code, selection.Expression, context, state)
	if t := jvmSize(e.ExpressionValue); t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      selection.Expression.ExpressionValue.Class.Name,
		Field:      selection.Name,
		Descriptor: Descriptor.typeDescriptor(e.ExpressionValue),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
