package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

/*
	make a default construction
*/
func mkClassDefaultContruction(class *cg.ClassHighLevel) {
	method := &cg.MethodHighLevel{}
	method.Name = specail_method_init
	method.Descriptor = "()V"
	method.AccessFlags |= cg.ACC_METHOD_PRIVATE
	method.Code.Codes = make([]byte, 5)
	method.Code.Codes[0] = cg.OP_aload_0
	method.Code.Codes[1] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      class.SuperClass,
		Name:       specail_method_init,
		Descriptor: "()V",
	}, method.Code.Codes[2:4])
	method.Code.Codes[4] = cg.OP_return
	method.Code.MaxStack = 1
	method.Code.MaxLocals = 1
	method.Code.CodeLength = uint16(len(method.Code.Codes))
	class.AppendMethod(method)
}

func backPatchEs(es []*cg.JumpBackPatch, to uint16) {
	for _, e := range es {
		offset := int16(int(to) - int(e.CurrentCodeLength))
		e.Bs[0] = byte(offset >> 8)
		e.Bs[1] = byte(offset)
	}
}

func jumpto(op byte, code *cg.AttributeCode, to uint16) {
	code.Codes[code.CodeLength] = cg.OP_goto
	b := &cg.JumpBackPatch{}
	b.CurrentCodeLength = code.CodeLength
	b.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	backPatchEs([]*cg.JumpBackPatch{b}, to)
	code.CodeLength += 3
}
