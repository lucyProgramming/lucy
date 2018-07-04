package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	class                   *ast.Class
	function                *ast.Function
	currentSourceFile       string
	currentLineNUmber       int
	Defer                   *ast.StatementDefer
	lastStackMapState       *StackMapState
	lastStackMapStateLocals []*cg.StackMapVerificationTypeInfo
	lastStackMapStateStacks []*cg.StackMapVerificationTypeInfo
	lastStackMapOffset      int
	stackMapOffsets         []int
}

func (context *Context) MakeStackMap(code *cg.AttributeCode, state *StackMapState, offset int) {
	if context.lastStackMapOffset == offset {
		code.AttributeStackMap.StackMaps = code.AttributeStackMap.StackMaps[0 : len(code.AttributeStackMap.StackMaps)-1]
		context.stackMapOffsets = context.stackMapOffsets[0 : len(context.stackMapOffsets)-1]
		context.lastStackMapState = nil
		if len(context.stackMapOffsets) > 0 {
			context.lastStackMapOffset = context.stackMapOffsets[len(context.stackMapOffsets)-1]
		} else {
			context.lastStackMapOffset = 0
		}
	}
	var delta uint16
	if context.lastStackMapOffset == 0 {
		/*
			first one
		*/
		delta = uint16(offset)
	} else {
		delta = uint16(offset - context.lastStackMapOffset - 1)
	}
	defer func() {
		context.lastStackMapOffset = offset // rewrite
		context.lastStackMapState = state
		context.lastStackMapStateLocals = make([]*cg.StackMapVerificationTypeInfo, len(state.Locals))
		copy(context.lastStackMapStateLocals, state.Locals)
		context.lastStackMapStateStacks = make([]*cg.StackMapVerificationTypeInfo, len(state.Stacks))
		copy(context.lastStackMapStateStacks, state.Stacks)
		context.stackMapOffsets = append(context.stackMapOffsets, offset)
	}()
	if state == context.lastStackMapState { // same state
		if len(state.Locals) == len(context.lastStackMapStateLocals) && len(state.Stacks) == 0 {
			/*
				same frame or same frame extended
			*/
			if delta <= 63 {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMapSameFrame{FrameType: byte(delta)})
			} else {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMapSameFrameExtended{FrameType: 251, Delta: delta})
			}
			return
		}
		if len(context.lastStackMapStateLocals) == len(state.Locals) && len(state.Stacks) == 1 { // 1 stack or 1 stack extended
			if delta <= 64 {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMapSameLocals1StackItemFrame{
						FrameType: byte(delta + 64),
						Stack:     state.Stacks[0],
					})
			} else {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMapSameLocals1StackItemFrameExtended{
						FrameType: 247,
						Delta:     delta,
						Stack:     state.Stacks[0],
					})
			}
			return
		}
		if len(context.lastStackMapStateLocals) < len(state.Locals) && len(state.Stacks) == 0 { // append frame
			num := len(state.Locals) - len(context.lastStackMapStateLocals)
			if num <= 3 {
				appendFrame := &cg.StackMapAppendFrame{}
				appendFrame.FrameType = byte(num + 251)
				appendFrame.Delta = delta
				appendFrame.Locals = make([]*cg.StackMapVerificationTypeInfo, num)
				copy(appendFrame.Locals, state.Locals[len(state.Locals)-num:])
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, appendFrame)
				return
			}
		}
	}
	// full frame
	fullFrame := &cg.StackMapFullFrame{}
	fullFrame.FrameType = 255
	fullFrame.Delta = delta
	fullFrame.Locals = make([]*cg.StackMapVerificationTypeInfo, len(state.Locals))
	copy(fullFrame.Locals, state.Locals)
	fullFrame.Stacks = make([]*cg.StackMapVerificationTypeInfo, len(state.Stacks))
	copy(fullFrame.Stacks, state.Stacks)
	code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, fullFrame)
	return
}

func (context *Context) appendLimeNumberAndSourceFile(pos *ast.Position,
	code *cg.AttributeCode, class *cg.ClassHighLevel) {
	if pos == nil {
		return
	}
	if pos.Filename != context.currentSourceFile {
		if class.SourceFiles == nil {
			class.SourceFiles = make(map[string]struct{})
		}
		class.SourceFiles[pos.Filename] = struct{}{}
		context.currentSourceFile = pos.Filename
		context.currentLineNUmber = pos.StartLine
		code.MKLineNumber(pos.StartLine)
		return
	}
	if context.currentLineNUmber != pos.StartLine {
		code.MKLineNumber(pos.StartLine)
		context.currentLineNUmber = pos.StartLine
	}
}
