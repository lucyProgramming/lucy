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
				FrameType: byte(delta + 64),
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
	if len(context.Locals) < last.Locals && len(context.Stacks) == 0 { // append frame
		num := len(context.Locals) - last.Locals
		if num <= 4 {
			appendFrame := &cg.StackMap_append_frame{}
			appendFrame.FrameType = byte(num + 251)
			appendFrame.Delta = delta
			appendFrame.Locals = make([]*cg.StackMap_verification_type_info, len(context.Locals[last.Locals:]))
			for k, _ := range appendFrame.Locals {
				appendFrame.Locals[k] = &cg.StackMap_verification_type_info{}
				appendFrame.Locals[k].T = &cg.StackMap_Top_variable_info{}
			}
			return appendFrame
		}
	}
	// full frame
	fullFrame := &cg.StackMap_full_frame{}
	fullFrame.FrameType = 255
	fullFrame.Delta = delta
	fullFrame.Locals = make([]*cg.StackMap_verification_type_info, len(context.Locals))
	copy(fullFrame.Locals, context.Locals)
	fullFrame.Stacks = make([]*cg.StackMap_verification_type_info, len(context.Stacks))
	copy(fullFrame.Stacks, context.Stacks)
	return fullFrame
}
func (context *Context) newStackMapVerificationTypeInfo(class *cg.ClassHighLevel, t *ast.VariableType, classname ...string) (ret []*cg.StackMap_verification_type_info) {
	ret = []*cg.StackMap_verification_type_info{}
	switch t.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		fallthrough
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Integer_variable_info{}
	case ast.VARIABLE_TYPE_LONG:
		ret = make([]*cg.StackMap_verification_type_info, 2)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[1] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Long_variable_info{}
		ret[1].T = &cg.StackMap_Top_variable_info{}
	case ast.VARIABLE_TYPE_FLOAT:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Float_variable_info{}
	case ast.VARIABLE_TYPE_DOUBLE:
		ret = make([]*cg.StackMap_verification_type_info, 2)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[1] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Double_variable_info{}
		ret[1].T = &cg.StackMap_Top_variable_info{}
	case ast.VARIABLE_TYPE_NULL:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Null_variable_info{}
	case ast.VARIABLE_TYPE_STRING:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_string_class),
		}
	case ast.VARIABLE_TYPE_OBJECT:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(classname[0]),
		}
	case ast.VARIABLE_TYPE_MAP:
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(java_hashmap_class),
		}
	case ast.VARIABLE_TYPE_ARRAY:
		meta := ArrayMetas[t.ArrayType.Typ]
		ret = make([]*cg.StackMap_verification_type_info, 1)
		ret[0] = &cg.StackMap_verification_type_info{}
		ret[0].T = &cg.StackMap_Object_variable_info{
			Index: class.Class.InsertClassConst(meta.classname),
		}
	}
	return ret
}
