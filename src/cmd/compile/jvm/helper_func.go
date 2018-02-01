package jvm

import (
	"encoding/binary"

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
	class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
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
func appendBackPatch(p *[][]byte, b []byte) {
	if *p == nil {
		*p = [][]byte{b}
	} else {
		*p = append(*p, b)
	}
}

/*
	backpatch exits
*/
func backPatchEs(es [][]byte, code *cg.AttributeCode) {
	for _, v := range es {
		binary.BigEndian.PutUint16(v, code.CodeLength)
	}
}
