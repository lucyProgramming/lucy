package jvm

import (
	"encoding/binary"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

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

//func mkPath(path, name string) string {
//	if path == "" {
//		return name
//	}
//	return path + "$" + name
//}
