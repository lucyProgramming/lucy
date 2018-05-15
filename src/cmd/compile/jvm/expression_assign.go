package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildExpressionAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	stackLength := len(state.Stacks)
	defer func() {
		state.popStack(len(state.Stacks) - stackLength)
	}()
	bin := e.Data.(*ast.ExpressionBinary)
	left := bin.Left.Data.([]*ast.Expression)[0]
	right := bin.Right.Data.([]*ast.Expression)[0]
	maxstack, remainStack, op, target, classname, name, descriptor := m.getLeftValue(class, code, left, context, state)
	stack, es := m.build(class, code, right, context, state)
	if len(es) > 0 {
		state.pushStack(class, right.Value)
		context.MakeStackMap(code, state, code.CodeLength)
		backPatchEs(es, code.CodeLength)
	}
	if t := remainStack + stack; t > maxstack {
		maxstack = t
	}
	currentStack := remainStack + jvmSize(target)
	currentStack += m.controlStack2FitAssign(code, op, classname, target)
	if currentStack > maxstack {
		maxstack = currentStack
	}
	if classname == java_hashmap_class && e.Value.IsPointer() == false {
		typeConverter.putPrimitiveInObject(class, code, target)
	}
	copyOPLeftValue(class, code, op, classname, name, descriptor)
	return
}

// a,b,c = 122,fdfd2232,"hello";
func (m *MakeExpression) buildAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context, state *StackMapState) (maxstack uint16) {
	if e.IsStatementExpression == false {
		return m.buildExpressionAssign(class, code, e, context, state)
	}
	bin := e.Data.(*ast.ExpressionBinary)
	lefts := bin.Left.Data.([]*ast.Expression)
	currentStack := uint16(0)
	targets := make([]*ast.VariableType, len(lefts))
	ops := make([][]byte, len(lefts))
	classnames := make([]string, len(lefts))
	names := make([]string, len(lefts))
	descriptors := make([]string, len(lefts))
	remainstacks := make([]uint16, len(lefts))
	noDestinations := make([]bool, len(lefts))
	stackHeights := make([]int, len(lefts))
	// put left value one the stack
	index := len(lefts) - 1
	for index >= 0 { //
		if lefts[index].Typ == ast.EXPRESSION_TYPE_IDENTIFIER &&
			lefts[index].Data.(*ast.ExpressionIdentifer).Name == ast.NO_NAME_IDENTIFIER {
			noDestinations[index] = true
		} else {
			h := len(state.Stacks)
			stack, remainstack, op, target, classname, name, descriptor := m.getLeftValue(class, code, lefts[index], context, state)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			targets[index] = target
			ops[index] = op
			classnames[index] = classname
			names[index] = name
			descriptors[index] = descriptor
			remainstacks[index] = remainstack
			stackHeights[index] = len(state.Stacks) - h
			currentStack += remainstack
		}
		index--
	}
	//
	rights := bin.Right.Data.([]*ast.Expression)
	slice := func() {
		currentStack -= remainstacks[0]
		if noDestinations[0] == false {
			copyOPLeftValue(class, code, ops[0], classnames[0], names[0], descriptors[0])
		}
		//let`s slice
		targets = targets[1:]
		ops = ops[1:]
		classnames = classnames[1:]
		names = names[1:]
		descriptors = descriptors[1:]
		remainstacks = remainstacks[1:]
		noDestinations = noDestinations[1:]
		state.popStack(stackHeights[0])
		stackHeights = stackHeights[1:]
	}
	for _, v := range rights {
		if v.MayHaveMultiValue() && len(v.Values) > 1 {
			stack, _ := m.build(class, code, v, context, state)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			// stack top is arrayListObject object
			if t := 1 + currentStack; t > maxstack {
				maxstack = t
			}
			arrayListPacker.storeArrayListAutoVar(code, context) // store it into local
			for k, v := range v.Values {                         // unpack
				needPutInObject := (classnames[0] == java_hashmap_class && targets[0].IsPointer() == false)
				stack = arrayListPacker.unPack(class, code, k, v, context)
				if t := stack + currentStack; t > maxstack {
					maxstack = t
				}
				if noDestinations[0] == false {
					if t := currentStack + jvmSize(targets[0]); t > maxstack { // incase int convert to double or long
						maxstack = t
					}
					if needPutInObject { // convert to primitive
						typeConverter.putPrimitiveInObject(class, code, targets[0])
					}
				} else { // pop fron stack
					if jvmSize(v) == 1 {
						code.Codes[code.CodeLength] = cg.OP_pop
					} else {
						code.Codes[code.CodeLength] = cg.OP_pop2
					}
					code.CodeLength++
				}
				slice()
			}
			continue
		}
		variableType := v.Value
		if v.MayHaveMultiValue() {
			variableType = v.Values[0]
		}
		needPutInObject := (classnames[0] == java_hashmap_class && targets[0].IsPointer() == false)
		stack, es := m.build(class, code, v, context, state)
		if len(es) > 0 {
			state.pushStack(class, variableType)
			context.MakeStackMap(code, state, code.CodeLength)
			state.popStack(1)                // must be interger
			backPatchEs(es, code.CodeLength) // true or false need to backpatch
		}
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if noDestinations[0] == false {
			if t := currentStack + jvmSize(targets[0]); t > maxstack { // incase int convert to double or long
				maxstack = t
			}
			if needPutInObject { // convert to primitive
				typeConverter.putPrimitiveInObject(class, code, targets[0])
			}
		} else { // pop fron stack
			if jvmSize(variableType) == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop2
			}
			code.CodeLength++
		}
		slice()
	}
	return
}

func (m *MakeExpression) controlStack2FitAssign(code *cg.AttributeCode, op []byte, classname string,
	stackTopType *ast.VariableType) (increment uint16) {
	// no object after value,just dup top
	if op[0] == cg.OP_istore || // 将栈顶 int 型数值存入指定局部变量。
		op[0] == cg.OP_lstore || //将栈顶 long 型数值存入指定局部变量。
		op[0] == cg.OP_fstore || //将栈顶 float 型数值存入指定局部变量。
		op[0] == cg.OP_dstore || //将栈顶 double 型数值存入指定局部变量。
		op[0] == cg.OP_astore || // 将栈顶引用型数值存入指定局部变量。
		op[0] == cg.OP_istore_0 || //将栈顶 int 型数值存入第一个局部变量。
		op[0] == cg.OP_istore_1 || // 将栈顶 int 型数值存入第二个局部变量。
		op[0] == cg.OP_istore_2 || //将栈顶 int 型数值存入第三个局部变量。
		op[0] == cg.OP_istore_3 || // 将栈顶 int 型数值存入第四个局部变量。
		op[0] == cg.OP_lstore_0 || //将栈顶 long 型数值存入第一个局部变量。
		op[0] == cg.OP_lstore_1 || // 将栈顶 long 型数值存入第二个局部变量。
		op[0] == cg.OP_lstore_2 || //将栈顶 long 型数值存入第三个局部变量。
		op[0] == cg.OP_lstore_3 || // 将栈顶 long 型数值存入第四个局部变量。
		op[0] == cg.OP_fstore_0 || //将栈顶 float 型数值存入第一个局部变量。
		op[0] == cg.OP_fstore_1 || //将栈顶 float 型数值存入第二个局部变量。
		op[0] == cg.OP_fstore_2 || //将栈顶 float 型数值存入第三个局部变量。
		op[0] == cg.OP_fstore_3 || //将栈顶 float 型数值存入第四个局部变量。
		op[0] == cg.OP_dstore_0 || //将栈顶 double 型数值存入第一个局部变量。
		op[0] == cg.OP_dstore_1 || //将栈顶 double 型数值存入第二个局部变量。
		op[0] == cg.OP_dstore_2 || // 将栈顶 double 型数值存入第三个局部变量。
		op[0] == cg.OP_dstore_3 || //将栈顶 double 型数值存入第四个局部变量。
		op[0] == cg.OP_astore_0 || //将栈顶引用型数值存入第一个局部变量。
		op[0] == cg.OP_astore_1 || ///将栈顶引用型数值存入第二个局部变量。
		op[0] == cg.OP_astore_2 || //将栈顶引用型数值存入第三个局部变量
		op[0] == cg.OP_astore_3 ||
		op[0] == cg.OP_putstatic { //为指定的类的静态域赋值。
		if jvmSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup
		} else {
			code.Codes[code.CodeLength] = cg.OP_dup2
			increment = 2
		}
		code.CodeLength++
		return
	}
	if op[0] == cg.OP_putfield {
		if jvmSize(stackTopType) == 1 {
			increment = 1
			code.Codes[code.CodeLength] = cg.OP_dup_x1
		} else {
			increment = 2
			code.Codes[code.CodeLength] = cg.OP_dup2_x1
		}
		code.CodeLength++
		return
	}
	if op[0] == cg.OP_invokevirtual { // array or map
		if ArrayMetasMap[classname] != nil || classname == java_hashmap_class {
			// stack is arrayref index or mapref kref which are all category 1 type
			if jvmSize(stackTopType) == 1 {
				increment = 1
				code.Codes[code.CodeLength] = cg.OP_dup_x2
				code.CodeLength++
			} else {
				increment = 2
				code.Codes[code.CodeLength] = cg.OP_dup2_x2
				code.CodeLength++
			}
			return
		}
	}
	panic("other case")
	return
}
