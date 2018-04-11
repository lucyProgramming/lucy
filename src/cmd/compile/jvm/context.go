package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	method            *cg.MethodHighLevel
	StackMapDelta     int
	function          *ast.Function
	currentSoureFile  string
	currentLineNUmber int
	Defers            []*ast.Defer
	Locals            []*cg.StackMap_verification_type_info
	Stacks            []*cg.StackMap_verification_type_info
}

func (c *Context) appendLimeNumberAndSourceFile(pos *ast.Pos, code *cg.AttributeCode, class *cg.ClassHighLevel) {
	if pos == nil {
		return
	}
	if pos.Filename != c.currentSoureFile {
		if class.SourceFiles == nil {
			class.SourceFiles = make(map[string]struct{})
		}
		class.SourceFiles[pos.Filename] = struct{}{}
		c.currentSoureFile = pos.Filename
		c.currentLineNUmber = pos.StartLine
		code.MKLineNumber(pos.StartLine)
		return
	}
	if c.currentLineNUmber != pos.StartLine {
		code.MKLineNumber(pos.StartLine)
		c.currentLineNUmber = pos.StartLine
	}
}

func (context *Context) MakeStackMap(last *StackMapStateLocalsNumber, offset int) cg.StackMap {
	var delta uint16
	if context.StackMapDelta == 0 {
		delta = uint16(offset)
	} else {
		delta = uint16(offset - context.StackMapDelta - 1)
	}
	defer func() {
		context.StackMapDelta = offset // rewrite
	}()
	if len(context.Locals) == last.Locals && len(context.Stacks) == 0 { // same frame or same frame extended
		if delta <= 63 {
			return &cg.StackMap_same_frame{FrameType: byte(delta)}
		} else {
			return &cg.StackMap_same_frame_extended{FrameType: 251, Delta: delta}
		}
	}
	if len(context.Locals) == last.Locals && len(context.Stacks) == 1 { // 1 stack or 1 stack extended
		if delta <= 64 {
			return &cg.StackMap_same_locals_1_stack_item_frame{
				FrameType: byte(delta),
				Stack:     context.Stacks[0],
			}
		} else {
			return &cg.StackMap_same_locals_1_stack_item_frame_extended{
				FrameType: 247,
				Delta:     delta,
				Stack:     context.Stacks[0],
			}
		}
	}
	if len(context.Locals) > last.Locals && len(context.Stacks) == 0 { // append frame
		num := len(context.Locals) - last.Locals
		if num <= 4 {
			appendFrame := &cg.StackMap_append_frame{}
			appendFrame.FrameType = byte(num + 251)
			appendFrame.Delta = delta
			appendFrame.Locals = context.Locals[last.Locals:][:] // make copy
			return appendFrame
		}

	}
	// full frame
	fullFrame := &cg.StackMap_full_frame{}
	fullFrame.FrameType = 255
	fullFrame.Locals = context.Locals[:] // make copy
	fullFrame.Stacks = context.Stacks[:] // make copy
	return fullFrame
}
