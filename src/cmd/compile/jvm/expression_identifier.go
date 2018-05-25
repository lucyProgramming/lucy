package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildCapturedIdentifer(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifer)
	captured := context.function.Closure.ClosureVariableExist(identifier.Var)
	if captured == false {
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, identifier.Var.LocalValOffset)...)
	} else {
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, 0)...)
		meta := closure.getMeta(identifier.Var.Typ.Typ)
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      class.Name,
			Field:      identifier.Var.Name,
			Descriptor: "L" + meta.className + ";",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	if 1 > maxstack {
		maxstack = 1
	}
	if t := jvmSize(identifier.Var.Typ); t > maxstack {
		maxstack = t
	}

	return
}

func (m *MakeExpression) buildIdentifer(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context) (maxstack uint16) {
	if e.Value.Typ == ast.VARIABLE_TYPE_CLASS {
		return
	}
	identifier := e.Data.(*ast.ExpressionIdentifer)
	if e.Value.Typ == ast.VARIABLE_TYPE_ENUM && identifier.EnumName != nil { // not a var
		loadInt32(class, code, identifier.EnumName.Value)
		maxstack = 1
		return
	}

	if identifier.Var.IsGlobal { //fetch global var
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      m.MakeClass.mainclass.Name,
			Field:      identifier.Name,
			Descriptor: Descriptor.typeDescriptor(identifier.Var.Typ),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxstack = jvmSize(identifier.Var.Typ)
		return
	}
	if identifier.Var.BeenCaptured {
		return m.buildCapturedIdentifer(class, code, e, context)
	}
	switch identifier.Var.Typ.Typ {

	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_ENUM:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_iload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_iload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_iload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_iload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_iload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxstack = 1
	case ast.VARIABLE_TYPE_FLOAT:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_fload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_fload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_fload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_fload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_fload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxstack = 1
	case ast.VARIABLE_TYPE_DOUBLE:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_dload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_dload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_dload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_dload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_dload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxstack = 2
	case ast.VARIABLE_TYPE_LONG:
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_lload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_lload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_lload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_lload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_lload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxstack = 2
	default: // object types
		if identifier.Var.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_aload_1
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_aload_2
			code.CodeLength++
		} else if identifier.Var.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_aload_3
			code.CodeLength++
		} else if identifier.Var.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_aload
			code.Codes[code.CodeLength+1] = byte(identifier.Var.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxstack = 1
	}
	return
}
