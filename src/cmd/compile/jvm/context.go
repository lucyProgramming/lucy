package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	class                   *ast.Class
	function                *ast.Function
	exceptionVarOffset      uint16
	multiValueVarOffset     uint16
	currentSourceFile       string
	currentLineNumber       int
	Defer                   *ast.StatementDefer
	lastStackMapState       *StackMapState
	lastStackMapStateLocals []*cg.StackMapVerificationTypeInfo
	lastStackMapStateStacks []*cg.StackMapVerificationTypeInfo
	lastStackMapOffset      int
	stackMapOffsets         []int
}

func (this *Context) MakeStackMap(code *cg.AttributeCode, state *StackMapState, offset int) {
	if this.lastStackMapOffset == offset {
		code.AttributeStackMap.StackMaps =
			code.AttributeStackMap.StackMaps[0 : len(code.AttributeStackMap.StackMaps)-1]
		this.stackMapOffsets = this.stackMapOffsets[0 : len(this.stackMapOffsets)-1]
		this.lastStackMapState = nil
		if len(this.stackMapOffsets) > 0 {
			this.lastStackMapOffset = this.stackMapOffsets[len(this.stackMapOffsets)-1]
		} else {
			this.lastStackMapOffset = -1
		}
	}
	var delta uint16
	if this.lastStackMapOffset == -1 {
		/*
			first one
		*/
		delta = uint16(offset)
	} else {
		delta = uint16(offset - this.lastStackMapOffset - 1)
	}
	defer func() {
		this.lastStackMapOffset = offset // rewrite
		this.lastStackMapState = state
		this.lastStackMapStateLocals = make([]*cg.StackMapVerificationTypeInfo, len(state.Locals))
		copy(this.lastStackMapStateLocals, state.Locals)
		this.lastStackMapStateStacks = make([]*cg.StackMapVerificationTypeInfo, len(state.Stacks))
		copy(this.lastStackMapStateStacks, state.Stacks)
		this.stackMapOffsets = append(this.stackMapOffsets, offset)
	}()
	if state == this.lastStackMapState { // same state
		if len(state.Locals) == len(this.lastStackMapStateLocals) && len(state.Stacks) == 0 {
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
		if len(this.lastStackMapStateLocals) == len(state.Locals) && len(state.Stacks) == 1 { // 1 stack or 1 stack extended
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
		if len(this.lastStackMapStateLocals) < len(state.Locals) && len(state.Stacks) == 0 { // append frame
			num := len(state.Locals) - len(this.lastStackMapStateLocals)
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

func (this *Context) appendLimeNumberAndSourceFile(
	pos *ast.Pos,
	code *cg.AttributeCode,
	class *cg.ClassHighLevel) {
	if pos == nil {
		return
	}
	if pos.Filename != this.currentSourceFile {
		class.InsertSourceFile(pos.Filename)
		this.currentSourceFile = pos.Filename
		this.currentLineNumber = pos.Line
		code.MKLineNumber(pos.Line)
		return
	}
	if this.currentLineNumber != pos.Line {
		code.MKLineNumber(pos.Line)
		this.currentLineNumber = pos.Line
	}
}
