package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

func (m *MakeExpression) buildDot(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	dot := e.Data.(*ast.ExpressionDot)
	if dot.Expression.Value.Typ == ast.VARIABLE_TYPE_PACKAGE {
		if dot.PackageVariable != nil {
			code.Codes[code.CodeLength] = cg.OP_getstatic
			class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
				Class:      dot.Expression.Value.Package.Name + "/main",
				Field:      dot.PackageVariable.Name,
				Descriptor: dot.PackageVariable.Descriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}
	// check cast to super class
	if dot.Name == ast.SUPER_FIELD_NAME {
		if dot.Expression.Value.Typ == ast.VARIABLE_TYPE_OBJECT {
			maxstack, _ = m.build(class, code, dot.Expression, context, nil)
			code.Codes[code.CodeLength] = cg.OP_checkcast
			class.InsertClassConst(e.Value.Class.Name, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		return
	}
	if dot.Expression.Value.Typ == ast.VARIABLE_TYPE_CLASS {
		maxstack = jvmSize(e.Value)
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      dot.Expression.Value.Class.Name,
			Field:      dot.Name,
			Descriptor: Descriptor.typeDescriptor(e.Value),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		return
	}
	maxstack, _ = m.build(class, code, dot.Expression, context, state)
	if t := jvmSize(e.Value); t > maxstack {
		maxstack = t
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
