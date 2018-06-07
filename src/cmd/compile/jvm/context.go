package jvm

import (
	//	"runtime/debug"

	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	class                   *ast.Class
	lastStackMapState       *StackMapState
	lastStackMapStateLocals []*cg.StackMap_verification_type_info
	lastStackMapStateStacks []*cg.StackMap_verification_type_info
	LastStackMapOffset      int
	NotFirstStackMap        bool
	function                *ast.Function
	currentSoureFile        string
	currentLineNUmber       int
	Defer                   *ast.Defer
}

func (context *Context) MakeStackMap(code *cg.AttributeCode, state *StackMapState, offset int) {
	//fmt.Println(offset)
	if context.LastStackMapOffset == offset && context.NotFirstStackMap {
		if state.isSame(context.lastStackMapStateLocals, context.lastStackMapStateStacks) {
			return // no need
		} else {
			panic(fmt.Sprintf("missing checking same offset:%d", offset))
		}
	}
	var delta uint16
	if context.NotFirstStackMap == false {
		delta = uint16(offset)
	} else {
		delta = uint16(offset - context.LastStackMapOffset - 1)
	}
	defer func() {
		context.LastStackMapOffset = offset // rewrite
		context.lastStackMapStateLocals = make([]*cg.StackMap_verification_type_info, len(state.Locals))
		copy(context.lastStackMapStateLocals, state.Locals)
		context.lastStackMapStateStacks = make([]*cg.StackMap_verification_type_info, len(state.Stacks))
		copy(context.lastStackMapStateStacks, state.Stacks)
		context.NotFirstStackMap = true
		context.lastStackMapState = state
	}()
	if state == context.lastStackMapState {
		if len(state.Locals) == len(context.lastStackMapStateLocals) && len(state.Stacks) == 0 { // same frame or same frame extended
			if delta <= 63 {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMap_same_frame{FrameType: byte(delta)})

			} else {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMap_same_frame_extended{FrameType: 251, Delta: delta})
			}
			return
		}
		if len(context.lastStackMapStateLocals) == len(state.Locals) && len(state.Stacks) == 1 { // 1 stack or 1 stack extended
			if delta <= 64 {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMap_same_locals_1_stack_item_frame{
						FrameType: byte(delta + 64),
						Stack:     state.Stacks[0],
					})
			} else {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMap_same_locals_1_stack_item_frame_extended{
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
				appendFrame := &cg.StackMap_append_frame{}
				appendFrame.FrameType = byte(num + 251)
				appendFrame.Delta = delta
				appendFrame.Locals = make([]*cg.StackMap_verification_type_info, num)
				copy(appendFrame.Locals, state.Locals[len(state.Locals)-num:])
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, appendFrame)
				return
			}
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
	code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, fullFrame)
	return
}

func (context *Context) appendLimeNumberAndSourceFile(pos *ast.Pos,
	code *cg.AttributeCode, class *cg.ClassHighLevel) {
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
