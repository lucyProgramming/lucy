package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) getCaptureIdentiferLeftValue(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack, remainStack uint16, op []byte, target *ast.VariableType, classname, fieldname, fieldDescriptor string) {
	return
}

func (m *MakeExpression) getMapLeftValue(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack, remainStack uint16, op []byte, target *ast.VariableType, classname, name, descriptor string) {
	index := e.Data.(*ast.ExpressionIndex)
	maxstack, _ = m.build(class, code, index.Expression, context)
	stack, _ := m.build(class, code, index.Index, context)
	if t := 1 + stack; t > maxstack {
		maxstack = t
	}
	PrimitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, index.Index.VariableType)
	remainStack = 2
	op = []byte{cg.OP_invokevirtual, cg.OP_pop}
	target = index.Expression.VariableType.Map.V
	classname = java_hashmap_class
	name = "put"
	descriptor = "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"
	return
}

func (m *MakeExpression) getLeftValue(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack, remainStack uint16, op []byte, target *ast.VariableType, classname, name, descriptor string) {
	switch e.Typ {
	case ast.EXPRESSION_TYPE_IDENTIFIER:
		identifier := e.Data.(*ast.ExpressionIdentifer)
		if identifier.Var.IsGlobal {
			op = []byte{cg.OP_putstatic}
			target = identifier.Var.Typ
			classname = context.mainclass.Name
			name = identifier.Name
			descriptor = Descriptor.typeDescriptor(identifier.Var.Typ)
			return
		}
		if identifier.Var.BeenCaptured {
			return m.getCaptureIdentiferLeftValue(class, code, e, context)
		}
		if identifier.Name == ast.NO_NAME_IDENTIFIER {
			return //
		}
		switch identifier.Var.Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			fallthrough
		case ast.VARIABLE_TYPE_BYTE:
			fallthrough
		case ast.VARIABLE_TYPE_SHORT:
			fallthrough
		case ast.VARIABLE_TYPE_INT:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_istore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_istore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_istore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_istore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_istore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local int var offset > 255")
			}
		case ast.VARIABLE_TYPE_FLOAT:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_fstore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_fstore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_fstore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_fstore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_fstore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		case ast.VARIABLE_TYPE_DOUBLE:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_dstore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_dstore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_dstore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_dstore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_dstore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		case ast.VARIABLE_TYPE_LONG:
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_lstore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_lstore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_lstore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_lstore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_lstore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		default: // must be a object type
			if identifier.Var.LocalValOffset == 0 {
				op = []byte{cg.OP_astore_0}
			} else if identifier.Var.LocalValOffset == 1 {
				op = []byte{cg.OP_astore_1}
			} else if identifier.Var.LocalValOffset == 2 {
				op = []byte{cg.OP_astore_2}
			} else if identifier.Var.LocalValOffset == 3 {
				op = []byte{cg.OP_astore_3}
			} else if identifier.Var.LocalValOffset <= 255 {
				op = []byte{cg.OP_astore, byte(identifier.Var.LocalValOffset)}
			} else {
				panic("local float var out of range")
			}
		}
		target = identifier.Var.Typ
	case ast.EXPRESSION_TYPE_INDEX:
		index := e.Data.(*ast.ExpressionIndex)
		if index.Expression.VariableType.Typ == ast.VARIABLE_TYPE_ARRAY {
			maxstack, _ = m.build(class, code, index.Expression, context)
			stack, _ := m.build(class, code, index.Index, context)
			if t := stack + 1; t > maxstack {
				maxstack = t
			}
			meta := ArrayMetas[e.VariableType.Typ]
			classname = meta.classname
			name = "set"
			descriptor = meta.setDescriptor
			target = e.VariableType
			remainStack = 2 // [objectref ,index]
			op = []byte{cg.OP_invokevirtual}
		} else { // map
			return m.getMapLeftValue(class, code, e, context)
		}

	case ast.EXPRESSION_TYPE_DOT:
		index := e.Data.(*ast.ExpressionIndex)
		if index.Expression.Typ == ast.VARIABLE_TYPE_CLASS {
			op = []byte{cg.OP_getstatic}
			classname = index.Expression.VariableType.Class.Name
			name = index.Name
			descriptor = index.Field.Descriptor
		} else {
			maxstack, _ = m.build(class, code, index.Expression, context)
			classname = index.Expression.VariableType.Class.Name
			name = index.Name
			descriptor = index.Field.Descriptor
		}
	default:
		panic("unkown type ")
	}
	return
}
