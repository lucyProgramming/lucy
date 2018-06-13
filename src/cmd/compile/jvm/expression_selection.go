package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

func (m *MakeExpression) buildSelection(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxStack uint16) {
	dot := e.Data.(*ast.ExpressionSelection)
	if dot.Expression.Value.Typ == ast.VARIABLE_TYPE_PACKAGE {
		maxStack = jvmSize(e.Value)
		if dot.PackageVariable != nil {
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      dot.Expression.Value.Package.Name + "/main",
				Field:      dot.PackageVariable.Name,
				Descriptor: dot.PackageVariable.Descriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if dot.EnumName != nil {
			loadInt32(class, code, dot.EnumName.Value)
		}
		return
	}
	// check cast to super class
	if dot.Name == ast.SUPER_FIELD_NAME {
		if dot.Expression.Value.Typ == ast.VARIABLE_TYPE_OBJECT {
			maxStack, _ = m.build(class, code, dot.Expression, context, nil)
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst(e.Value.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}
	if dot.Expression.Value.Typ == ast.VARIABLE_TYPE_CLASS {
		maxStack = jvmSize(e.Value)
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      dot.Expression.Value.Class.Name,
			Field:      dot.Name,
			Descriptor: Descriptor.typeDescriptor(e.Value),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	// object
	maxStack, _ = m.build(class, code, dot.Expression, context, state)
	if t := jvmSize(e.Value); t > maxStack {
		maxStack = t
	}
	code.Codes[code.CodeLength] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      dot.Expression.Value.Class.Name,
		Field:      dot.Name,
		Descriptor: Descriptor.typeDescriptor(e.Value),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
