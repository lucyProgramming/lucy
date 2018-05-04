package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildJavaArrayMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	call := e.Data.(*ast.ExpressionMethodCall)
	maxstack, _ = m.build(class, code, call.Expression, context, state)
	switch call.Name {
	case common.ARRAY_METHOD_SIZE:
		code.Codes[code.CodeLength] = cg.OP_arraylength
		code.CodeLength++
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}

func (m *MakeExpression) buildArrayMethodCall(class *cg.ClassHighLevel, code *cg.AttributeCode,
	e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	length := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - length) // ref type
	}()
	call := e.Data.(*ast.ExpressionMethodCall)
	maxstack, _ = m.build(class, code, call.Expression, context, state)
	state.pushStack(class, call.Expression.Value)
	switch call.Name {
	case common.ARRAY_METHOD_CAP,
		common.ARRAY_METHOD_SIZE,
		common.ARRAY_METHOD_START,
		common.ARRAY_METHOD_END:
		meta := ArrayMetas[call.Expression.Value.ArrayType.Typ]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.classname,
			Method:     call.Name,
			Descriptor: "()I",
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	case common.ARRAY_METHOD_APPEND:
		meta := ArrayMetas[call.Expression.Value.ArrayType.Typ]
		appendName := "append"
		appendDescriptor := meta.appendDescriptor
		for _, v := range call.Args {
			currentStack := uint16(1)
			if v.MayHaveMultiValue() && len(v.Values) > 0 {
				stack, _ := m.build(class, code, v, context, nil)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				arrayListPacker.storeArrayListAutoVar(code, context)
				for kk, t := range v.Values {
					currentStack = 1
					if t := arrayListPacker.unPack(class, code, kk, t, context) + currentStack; t > maxstack {
						maxstack = t
					}
					if t := currentStack + jvmSize(t); t > maxstack {
						maxstack = t
					}
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      meta.classname,
						Method:     appendName,
						Descriptor: meta.appendDescriptor,
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				}
				continue
			}
			stack, es := m.build(class, code, v, context, state)
			if len(es) > 0 {
				backPatchEs(es, code.CodeLength)
				state.pushStack(class, v.Value)
				context.MakeStackMap(code, state, code.CodeLength)
				state.popStack(1) // must be a logical expression
			}
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.classname,
				Method:     appendName,
				Descriptor: appendDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	case common.ARRAY_METHOD_APPEND_ALL:
		meta := ArrayMetas[call.Expression.Value.ArrayType.Typ]
		for _, v := range call.Args {
			currentStack := uint16(1)
			appendName := "append"
			appendDescriptor := meta.appendAllDescriptor
			if v.MayHaveMultiValue() && len(v.Values) > 0 {
				stack, _ := m.build(class, code, v, context, state)
				if t := currentStack + stack; t > maxstack {
					maxstack = t
				}
				arrayListPacker.storeArrayListAutoVar(code, context)
				for kk, _ := range v.Values {
					currentStack := uint16(1)
					arrayListPacker.buildLoadArrayListAutoVar(code, context)
					loadInt32(class, code, int32(kk))
					if t := currentStack + 2; t > maxstack {
						maxstack = t
					}
					//
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      java_arrylist_class,
						Method:     "get",
						Descriptor: "(I)Ljava/lang/Object;",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
					//cast to real object
					code.Codes[code.CodeLength] = cg.OP_checkcast
					class.InsertClassConst(meta.classname, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
						Class:      meta.classname,
						Method:     appendName,
						Descriptor: appendDescriptor,
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				}
				continue
			}
			stack, _ := m.build(class, code, v, context, state)
			if t := stack + currentStack; t > maxstack {
				maxstack = t
			}
			//get elements field
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      meta.classname,
				Method:     appendName,
				Descriptor: appendDescriptor,
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		}
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	case common.ARRAY_METHOD_TO_JAVA:
		meta := ArrayMetas[call.Expression.Value.ArrayType.Typ]
		code.Codes[code.CodeLength] = cg.OP_invokevirtual
		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      meta.classname,
			Method:     "getJavaArray",
			Descriptor: meta.getJavaArrayDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		if e.IsStatementExpression {
			code.Codes[code.CodeLength] = cg.OP_pop
			code.CodeLength++
		}
	}
	return
}
