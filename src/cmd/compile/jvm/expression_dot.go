package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

func (m *MakeExpression) buildDot(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	dot := e.Data.(*ast.ExpressionDot)
	if dot.Expression.VariableType.Typ == ast.VARIABLE_TYPE_PACKAGE {
		if dot.PackageVariableDefinition != nil {
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      dot.Expression.VariableType.Package.Name + "/main",
				Field:      dot.PackageVariableDefinition.Name,
				Descriptor: dot.PackageVariableDefinition.Descriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}
	if dot.Name == ast.SUPER_FIELD_NAME {
		if dot.Expression.VariableType.Typ == ast.VARIABLE_TYPE_OBJECT {
			maxstack, _ = m.build(class, code, dot.Expression, context, nil)
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst(e.VariableType.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}
	if dot.Expression.VariableType.Typ == ast.VARIABLE_TYPE_CLASS {
		maxstack = e.VariableType.JvmSlotSize()
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      dot.Expression.VariableType.Class.Name,
			Field:      dot.Name,
			Descriptor: Descriptor.typeDescriptor(e.VariableType),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	maxstack, _ = m.build(class, code, dot.Expression, context, nil)
	if t := e.VariableType.JvmSlotSize(); t > maxstack {
		maxstack = t
	}
	code.Codes[code.CodeLength] = cg.OP_getfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      dot.Expression.VariableType.Class.Name,
		Field:      dot.Name,
		Descriptor: Descriptor.typeDescriptor(e.VariableType),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}
