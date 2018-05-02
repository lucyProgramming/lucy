package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	lastStackMapState  *StackMapState
	LastStackMapOffset int
	function           *ast.Function
	currentSoureFile   string
	currentLineNUmber  int
	Defer              *ast.Defer
}

func (context *Context) MakeStackMap(code *cg.AttributeCode, state *StackMapState, offset int) {
	//fmt.Println("offset:", offset)
	if context.LastStackMapOffset == offset {
		return
	}
	var delta uint16
	if context.LastStackMapOffset == 0 {
		delta = uint16(offset)
	} else {
		delta = uint16(offset - context.LastStackMapOffset - 1)
	}
	defer func() {
		context.LastStackMapOffset = offset // rewrite
		context.lastStackMapState = state
		state.LastLocals = make([]*cg.StackMap_verification_type_info, len(state.Locals))
		copy(state.LastLocals, state.Locals)
	}()
	if context.lastStackMapState != nil && context.lastStackMapState == state {
		if context.LastStackMapOffset == 0 { // first one
			if len(state.Stacks) == 0 && len(state.Locals)-len(state.LastLocals) <= 3 {
				// append frame
				num := len(state.Stacks) - len(state.LastLocals)
				appendFrame := &cg.StackMap_append_frame{}
				appendFrame.FrameType = byte(len(state.Locals) + 251)
				appendFrame.Delta = delta
				appendFrame.Locals = make([]*cg.StackMap_verification_type_info, num)
				for i := len(state.LastLocals); i < num; i++ {
					appendFrame.Locals[i-len(state.LastLocals)] = state.Locals[i]
				}
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps, appendFrame)
				return
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

		if len(state.Locals) == len(state.LastLocals) && len(state.Stacks) == 0 { // same frame or same frame extended
			if delta <= 63 {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMap_same_frame{FrameType: byte(delta)})

			} else {
				code.AttributeStackMap.StackMaps = append(code.AttributeStackMap.StackMaps,
					&cg.StackMap_same_frame_extended{FrameType: 251, Delta: delta})
			}
			return
		}
		if len(state.LastLocals) == len(state.Locals) && len(state.Stacks) == 1 { // 1 stack or 1 stack extended
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
		if len(state.LastLocals) < len(state.Locals) && len(state.Stacks) == 0 { // append frame
			num := len(state.Locals) - len(state.LastLocals)
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
