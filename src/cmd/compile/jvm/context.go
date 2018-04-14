package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	method            *cg.MethodHighLevel
	function          *ast.Function
	currentSoureFile  string
	currentLineNUmber int
	Defers            []*ast.Defer
	StackMapDelta     int
	block             *ast.Block
}

func (context *Context) appendLimeNumberAndSourceFile(pos *ast.Pos, code *cg.AttributeCode, class *cg.ClassHighLevel) {
	if pos == nil {
		return
	}
	if pos.Filename != context.currentSoureFile {
		if class.SourceFiles == nil {
			class.SourceFiles = make(map[string]struct{})
		}
		class.SourceFiles[pos.Filename] = struct{}{}
		context.currentSoureFile = pos.Filename
		context.currentLineNUmber = pos.StartLine
		code.MKLineNumber(pos.StartLine)
		return
	}
	if context.currentLineNUmber != pos.StartLine {
		code.MKLineNumber(pos.StartLine)
		context.currentLineNUmber = pos.StartLine
	}
}

func (context *Context) MakeStackMap(state *StackMapState, offset int) cg.StackMap {
	var delta uint16
	if context.StackMapDelta == 0 {
		delta = uint16(offset)
	} else {
		delta = uint16(offset - context.StackMapDelta - 1)
	}
	defer func() {
		context.StackMapDelta = offset // rewrite
		context.block.LastStackMapLocalVars = make([]*cg.StackMap_verification_type_info, len(state.Locals))
		copy(context.block.LastStackMapLocalVars, state.Locals)
	}()

	if context.StackMapDelta == 0 { // first one
		if len(state.Stacks) == 0 && len(state.Stacks)-len(context.block.LastStackMapLocalVars) <= 3 {
			// append frame
			num := len(state.Stacks) - len(context.block.LastStackMapLocalVars)
			appendFrame := &cg.StackMap_append_frame{}
			appendFrame.FrameType = byte(len(state.Locals) + 251)
			appendFrame.Delta = delta
			appendFrame.Locals = make([]*cg.StackMap_verification_type_info, num)
			for i := len(context.block.LastStackMapLocalVars); i < num; i++ {
				appendFrame.Locals[i-len(context.block.LastStackMapLocalVars)] = state.Locals[i]
			}
			return appendFrame
		}
		// full frame
		fullFrame := &cg.StackMap_full_frame{}
		fullFrame.FrameType = 255
		fullFrame.Delta = delta
		fullFrame.Locals = make([]*cg.StackMap_verification_type_info, len(state.Locals))
		copy(fullFrame.Locals, state.Locals)
		fullFrame.Stacks = make([]*cg.StackMap_verification_type_info, len(state.Stacks))
		copy(fullFrame.Stacks, state.Stacks)
		return fullFrame
	}

	if len(state.Locals) == len(context.block.LastStackMapLocalVars) && len(state.Stacks) == 0 { // same frame or same frame extended
		if delta <= 63 {
			return &cg.StackMap_same_frame{FrameType: byte(delta)}
		} else {
			panic(delta)
			return &cg.StackMap_same_frame_extended{FrameType: 251, Delta: delta}
		}
	}

	if len(context.block.LastStackMapLocalVars) == len(state.Locals) && len(state.Stacks) == 1 { // 1 stack or 1 stack extended
		if delta <= 64 {
			return &cg.StackMap_same_locals_1_stack_item_frame{
				FrameType: byte(delta + 64),
				Stack:     state.Stacks[0],
			}
		} else {
			return &cg.StackMap_same_locals_1_stack_item_frame_extended{
				FrameType: 247,
				Delta:     delta,
				Stack:     state.Stacks[0],
			}
		}
	}
	if len(context.block.LastStackMapLocalVars) < len(state.Locals) && len(state.Stacks) == 0 { // append frame
		num := len(state.Locals) - len(context.block.LastStackMapLocalVars)
		if num <= 3 {
			appendFrame := &cg.StackMap_append_frame{}
			appendFrame.FrameType = byte(num + 251)
			appendFrame.Delta = delta
			appendFrame.Locals = make([]*cg.StackMap_verification_type_info, num)
			copy(appendFrame.Locals, state.Locals[len(state.Locals)-num:])
			return appendFrame
		}
	}
	// full frame
	fullFrame := &cg.StackMap_full_frame{}
	fullFrame.FrameType = 255
	fullFrame.Delta = delta
	fullFrame.Locals = make([]*cg.StackMap_verification_type_info, len(state.Locals))
	copy(fullFrame.Locals, state.Locals)
	fullFrame.Stacks = make([]*cg.StackMap_verification_type_info, len(state.Stacks))
	copy(fullFrame.Stacks, state.Stacks)
	return fullFrame
}
