package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeClass) buildFunctionParameterAndReturnList(class *cg.ClassHighLevel, code *cg.AttributeCode, ft *ast.FunctionType, context *Context) {
	for _, v := range ft.ReturnList {
		currentStack := uint16(0)
		if v.BeenCaptured { //create closure object
			maxstack := closure.createCloureVar(class, code, v)
			if maxstack > code.MaxStack {
				code.MaxStack = maxstack
			}
			// then load
			copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, v.LocalValOffset)...)
			currentStack = 1
		}
		stack, es := m.MakeExpression.build(class, code, v.Expression, context)
		backPatchEs(es, code.CodeLength)
		if t := currentStack + stack; t > code.MaxStack {
			code.MaxStack = t
		}
		if v.Typ.IsNumber() && v.Typ.Typ != v.Expression.VariableType.Typ {
			m.MakeExpression.numberTypeConverter(code, v.Expression.VariableType.Typ, v.Typ.Typ)
		}
		if t := currentStack + v.Typ.JvmSlotSize(); t > code.MaxStack {
			code.MaxStack = t
		}
		m.storeLocalVar(class, code, v)
	}
}

func (m *MakeClass) buildFunction(class *cg.ClassHighLevel, method *cg.MethodHighLevel, f *ast.Function) {
	context := &Context{}
	context.function = f
	method.Code.Codes = make([]byte, 65536)
	method.Code.CodeLength = 0
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
		method.Code.MaxLocals = f.VarOffset // could  new slot when compile
	}()

	// if function is main
	if f.Name == ast.MAIN_FUNCTION_NAME {
		code := &method.Code
		code.Codes[code.CodeLength] = cg.OP_new
		meta := ArrayMetas[ast.VARIABLE_TYPE_STRING]
		class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
		copyOP(code, loadSimpleVarOp(ast.VARIABLE_TYPE_STRING, 0)...)
		if 3 > code.MaxStack {
			code.MaxStack = 3
		}
		code.Codes[code.CodeLength] = cg.OP_invokespecial
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.classname,
			Method:     special_method_init,
			Descriptor: meta.constructorFuncDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		copyOP(code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, 0)...)
	}

	m.buildFunctionParameterAndReturnList(class, &method.Code, f.Typ, context)
	if f.AutoVarForReturnBecauseOfDefer != nil {
		method.Code.Codes[method.Code.CodeLength] = cg.OP_iconst_0
		method.Code.CodeLength++
		copyOP(&method.Code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, f.AutoVarForReturnBecauseOfDefer.ExceptionIsNotNilWhenEnter)...)
		if len(f.Typ.ReturnList) > 1 {
			method.Code.Codes[method.Code.CodeLength] = cg.OP_iconst_0
			method.Code.CodeLength++
			copyOP(&method.Code, storeSimpleVarOp(ast.VARIABLE_TYPE_OBJECT, f.AutoVarForReturnBecauseOfDefer.MultiValueOffset)...)
			method.Code.Codes[method.Code.CodeLength] = cg.OP_iconst_0
			method.Code.CodeLength++
			copyOP(&method.Code, storeSimpleVarOp(ast.VARIABLE_TYPE_INT, f.AutoVarForReturnBecauseOfDefer.IfReachBotton)...)
		}
	}
	m.buildBlock(class, &method.Code, f.Block, context)
	return
}

func (m *MakeClass) loadLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		return closure.loadLocalCloureVar(class, code, v)
	}
	maxstack = v.Typ.JvmSlotSize()
	switch v.Typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_iload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_iload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_iload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_iload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_iload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_LONG:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_lload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_lload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_lload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_lload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_lload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_fload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_fload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_fload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_fload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_fload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_dload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_dload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_dload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_dload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_dload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY: //[]int
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_aload_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_aload_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_aload_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_aload_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_aload
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	}
	return
}

func (m *MakeClass) storeLocalVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.VariableDefinition) (maxstack uint16) {
	if v.BeenCaptured {
		closure.storeLocalCloureVar(class, code, v)
		return
	}
	maxstack = v.Typ.JvmSlotSize()
	switch v.Typ.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_istore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_istore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_istore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_istore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_istore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_LONG:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_lstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_lstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_lstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_lstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_lstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_FLOAT:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_fstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_fstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_fstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_fstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_fstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_DOUBLE:
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_dstore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_dstore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_dstore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_dstore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_dstore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	case ast.VARIABLE_TYPE_STRING:
		fallthrough
	case ast.VARIABLE_TYPE_OBJECT:
		fallthrough
	case ast.VARIABLE_TYPE_MAP:
		fallthrough
	case ast.VARIABLE_TYPE_ARRAY: //[]int
		switch v.LocalValOffset {
		case 0:
			code.Codes[code.CodeLength] = cg.OP_astore_0
			code.CodeLength++
		case 1:
			code.Codes[code.CodeLength] = cg.OP_astore_1
			code.CodeLength++
		case 2:
			code.Codes[code.CodeLength] = cg.OP_astore_2
			code.CodeLength++
		case 3:
			code.Codes[code.CodeLength] = cg.OP_astore_3
			code.CodeLength++
		default:
			if v.LocalValOffset > 255 {
				panic("over 255")
			}
			code.Codes[code.CodeLength] = cg.OP_astore
			code.Codes[code.CodeLength+1] = byte(v.LocalValOffset)
			code.CodeLength += 2
		}
	}
	return
}
