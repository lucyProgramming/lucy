package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (makeExpression *MakeExpression) buildCapturedIdentifier(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context) (maxStack uint16) {
	identifier := e.Data.(*ast.ExpressionIdentifier)
	captured := context.function.Closure.ClosureVariableExist(identifier.Variable)
	if captured == false {
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, identifier.Variable.LocalValOffset)...)
	} else {
		copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, 0)...)
		meta := closure.getMeta(identifier.Variable.Type.Type)
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      class.Name,
			Field:      identifier.Variable.Name,
			Descriptor: "L" + meta.className + ";",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	}
	if 1 > maxStack {
		maxStack = 1
	}
	closure.unPack(class, code, identifier.Variable.Type)
	if t := jvmSlotSize(identifier.Variable.Type); t > maxStack {
		maxStack = t
	}
	return
}

func (makeExpression *MakeExpression) buildIdentifier(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression,
	context *Context) (maxStack uint16) {
	if e.ExpressionValue.Type == ast.VariableTypeClass {
		return
	}
	identifier := e.Data.(*ast.ExpressionIdentifier)
	if e.ExpressionValue.Type == ast.VariableTypeEnum && identifier.EnumName != nil { // not a var
		loadInt32(class, code, identifier.EnumName.Value)
		maxStack = 1
		return
	}
	if identifier.Function != nil {
		return makeExpression.MakeClass.packFunction2MethodHandle(class, code, identifier.Function, context)
	}

	if identifier.Variable.IsGlobal { //fetch global var
		code.Codes[code.CodeLength] = cg.OP_getstatic
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      makeExpression.MakeClass.mainClass.Name,
			Field:      identifier.Name,
			Descriptor: JvmDescriptor.typeDescriptor(identifier.Variable.Type),
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		maxStack = jvmSlotSize(identifier.Variable.Type)
		return
	}
	if identifier.Variable.BeenCaptured {
		return makeExpression.buildCapturedIdentifier(class, code, e, context)
	}
	switch identifier.Variable.Type.Type {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		if identifier.Variable.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_iload_0
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_iload_1
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_iload_2
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_iload_3
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_iload
			code.Codes[code.CodeLength+1] = byte(identifier.Variable.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxStack = 1
	case ast.VariableTypeFloat:
		if identifier.Variable.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_fload_0
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_fload_1
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_fload_2
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_fload_3
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_fload
			code.Codes[code.CodeLength+1] = byte(identifier.Variable.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxStack = 1
	case ast.VariableTypeDouble:
		if identifier.Variable.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_dload_0
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_dload_1
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_dload_2
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_dload_3
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_dload
			code.Codes[code.CodeLength+1] = byte(identifier.Variable.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxStack = 2
	case ast.VariableTypeLong:
		if identifier.Variable.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_lload_0
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_lload_1
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_lload_2
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_lload_3
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_lload
			code.Codes[code.CodeLength+1] = byte(identifier.Variable.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxStack = 2
	default: // object types
		if identifier.Variable.LocalValOffset == 0 {
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 1 {
			code.Codes[code.CodeLength] = cg.OP_aload_1
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 2 {
			code.Codes[code.CodeLength] = cg.OP_aload_2
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset == 3 {
			code.Codes[code.CodeLength] = cg.OP_aload_3
			code.CodeLength++
		} else if identifier.Variable.LocalValOffset < 255 {
			code.Codes[code.CodeLength] = cg.OP_aload
			code.Codes[code.CodeLength+1] = byte(identifier.Variable.LocalValOffset)
			code.CodeLength += 2
		} else {
			panic("over 255")
		}
		maxStack = 1
	}
	return
}
