package jvm

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (buildPackage *BuildPackage) buildSwitchStatement(class *cg.ClassHighLevel, code *cg.AttributeCode,
	s *ast.StatementSwitch, context *Context, state *StackMapState) (maxStack uint16) {
	// if equal,leave 0 on stack
	compare := func(t *ast.Type) {
		state.popStack(2)
		switch t.Type {
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeInt:
			code.Codes[code.CodeLength] = cg.OP_isub
			code.CodeLength++
		case ast.VariableTypeLong:
			code.Codes[code.CodeLength] = cg.OP_lcmp
			code.CodeLength++
		case ast.VariableTypeFloat:
			code.Codes[code.CodeLength] = cg.OP_fcmpg
			code.CodeLength++
		case ast.VariableTypeDouble:
			code.Codes[code.CodeLength] = cg.OP_dcmpg
			code.CodeLength++
		case ast.VariableTypeString:
			code.Codes[code.CodeLength] = cg.OP_invokevirtual
			class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
				Class:      javaStringClass,
				Method:     "compareTo",
				Descriptor: "(Ljava/lang/String;)I",
			}, code.Codes[code.CodeLength+1:code.CodeLength+3])
			code.CodeLength += 3
		case ast.VariableTypeFunction:
			fallthrough
		case ast.VariableTypeObject:
			fallthrough
		case ast.VariableTypeMap:
			fallthrough
		case ast.VariableTypeArray:
			context.MakeStackMap(code, state, code.CodeLength+7)
			state.pushStack(class, &ast.Type{
				Type: ast.VariableTypeBool,
			})
			context.MakeStackMap(code, state, code.CodeLength+8)
			state.popStack(1)
			code.Codes[code.CodeLength] = cg.OP_if_acmpeq
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+1:code.CodeLength+3], 7)
			code.Codes[code.CodeLength+3] = cg.OP_iconst_1
			code.Codes[code.CodeLength+4] = cg.OP_goto
			binary.BigEndian.PutUint16(code.Codes[code.CodeLength+5:code.CodeLength+7], 4)
			code.Codes[code.CodeLength+7] = cg.OP_iconst_0
			code.CodeLength += 8
		}
	}
	maxStack = buildPackage.BuildExpression.build(class, code, s.Condition, context, state)
	//value is on stack
	var exit *cg.Exit
	size := jvmSlotSize(s.Condition.Value)
	currentStack := size
	state.pushStack(class, s.Condition.Value)
	for _, c := range s.StatementSwitchCases {
		if exit != nil {
			writeExits([]*cg.Exit{exit}, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
		matches := []*cg.Exit{}
		for _, ee := range c.Matches {
			if ee.MayHaveMultiValue() && len(ee.MultiValues) > 1 {
				stack := buildPackage.BuildExpression.build(class, code, ee, context, state)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				multiValuePacker.storeMultiValueAutoVar(code, context)
				for kkk, ttt := range ee.MultiValues {
					currentStack = size
					if size == 1 {
						code.Codes[code.CodeLength] = cg.OP_dup
					} else {
						code.Codes[code.CodeLength] = cg.OP_dup2
					}
					code.CodeLength++
					state.pushStack(class, s.Condition.Value)
					currentStack += size
					if currentStack > maxStack {
						maxStack = currentStack
					}
					stack = multiValuePacker.unPack(class, code, kkk, ttt, context)
					if t := stack + currentStack; t > maxStack {
						maxStack = t
					}
					state.pushStack(class, s.Condition.Value)
					compare(s.Condition.Value)
					// consume result on stack
					matches = append(matches, (&cg.Exit{}).Init(cg.OP_ifeq, code))
				}
				continue
			}
			currentStack = size
			// mk stack ready
			if size == 1 {
				code.Codes[code.CodeLength] = cg.OP_dup
			} else {
				code.Codes[code.CodeLength] = cg.OP_dup2
			}
			code.CodeLength++
			currentStack += size
			state.pushStack(class, s.Condition.Value)
			stack := buildPackage.BuildExpression.build(class, code, ee, context, state)
			if t := currentStack + stack; t > maxStack {
				maxStack = t
			}
			state.pushStack(class, s.Condition.Value)
			compare(s.Condition.Value)
			matches = append(matches, (&cg.Exit{}).Init(cg.OP_ifeq, code)) // comsume result on stack
		}
		// should be goto next,here is no match
		exit = (&cg.Exit{}).Init(cg.OP_goto, code)
		// if match goto here
		writeExits(matches, code.CodeLength)
		//before block,pop off stack
		context.MakeStackMap(code, state, code.CodeLength)
		if size == 1 {
			code.Codes[code.CodeLength] = cg.OP_pop
		} else {
			code.Codes[code.CodeLength] = cg.OP_pop2
		}
		code.CodeLength++
		//block is here
		if c.Block != nil {
			ss := (&StackMapState{}).FromLast(state)
			buildPackage.buildBlock(class, code, c.Block, context, ss)
			state.addTop(ss)
		}
		if c.Block == nil || c.Block.WillNotExecuteToEnd == false {
			s.Exits = append(s.Exits,
				(&cg.Exit{}).Init(cg.OP_goto, code)) // matched,goto switch outside
		}
	}
	writeExits([]*cg.Exit{exit}, code.CodeLength)
	context.MakeStackMap(code, state, code.CodeLength)
	if size == 1 {
		code.Codes[code.CodeLength] = cg.OP_pop
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop2
	}
	code.CodeLength++
	state.popStack(1)
	// build default
	if s.Default != nil {
		var ss *StackMapState
		if s.Default.HaveVariableDefinition() {
			ss = (&StackMapState{}).FromLast(state)
		} else {
			ss = state
		}
		buildPackage.buildBlock(class, code, s.Default, context, ss)
		state.addTop(ss)
	}
	return
}
