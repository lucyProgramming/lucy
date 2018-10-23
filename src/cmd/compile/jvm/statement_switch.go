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
		switch t.Type {
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeEnum:
			fallthrough
		case ast.VariableTypeChar:
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
		default:
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
	for _, v := range s.PrefixExpressions {
		stack := buildPackage.BuildExpression.build(class, code, v, context, state)
		if stack > maxStack {
			maxStack = stack
		}
	}
	stack := buildPackage.BuildExpression.build(class, code, s.Condition, context, state)
	if stack > maxStack {
		maxStack = stack
	}
	//value is on stack
	var notMatch *cg.Exit
	size := jvmSlotSize(s.Condition.Value)
	currentStack := size
	state.pushStack(class, s.Condition.Value)
	for _, c := range s.StatementSwitchCases {
		if notMatch != nil {
			writeExits([]*cg.Exit{notMatch}, code.CodeLength)
			context.MakeStackMap(code, state, code.CodeLength)
		}
		matches := []*cg.Exit{}
		for _, ee := range c.Matches {
			if ee.Value.Type == ast.VariableTypeBool {
				currentStack = size
				stack := buildPackage.BuildExpression.build(class, code, ee, context, state)
				if t := currentStack + stack; t > maxStack {
					maxStack = t
				}
				matches = append(matches, (&cg.Exit{}).Init(cg.OP_ifne, code)) // comsume result on stack
			} else {
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
				state.popStack(1)
				compare(s.Condition.Value)
				matches = append(matches, (&cg.Exit{}).Init(cg.OP_ifeq, code)) // comsume result on stack
			}
		}
		// should be goto next,here is no match
		notMatch = (&cg.Exit{}).Init(cg.OP_goto, code)
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
			ss := (&StackMapState{}).initFromLast(state)
			buildPackage.buildBlock(class, code, c.Block, context, ss)
			state.addTop(ss)
		}
		if c.Block == nil || c.Block.NotExecuteToLastStatement == false {
			s.Exits = append(s.Exits,
				(&cg.Exit{}).Init(cg.OP_goto, code)) // matched,goto switch outside
		}
	}
	writeExits([]*cg.Exit{notMatch}, code.CodeLength)
	context.MakeStackMap(code, state, code.CodeLength)
	if size == 1 {
		code.Codes[code.CodeLength] = cg.OP_pop
	} else {
		code.Codes[code.CodeLength] = cg.OP_pop2
	}
	code.CodeLength++
	state.popStack(1)
	if s.Default != nil {
		ss := (&StackMapState{}).initFromLast(state)
		buildPackage.buildBlock(class, code, s.Default, context, ss)
		state.addTop(ss)
	}
	return
}
