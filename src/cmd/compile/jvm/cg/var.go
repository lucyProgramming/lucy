package cg

import (
	"encoding/binary"
)

const (
	CONSTANT_POOL_MAX_SIZE = 65536
)

func backPatchIndex(locations [][]byte, index uint16) {
	for _, v := range locations {
		binary.BigEndian.PutUint16(v, index)
	}
}

type JumpBackPatch struct {
	CurrentCodeLength uint16
	Bs                []byte
}

func (j *JumpBackPatch) FromCode(op byte, code *AttributeCode) *JumpBackPatch {
	j.CurrentCodeLength = code.CodeLength
	code.Codes[code.CodeLength] = op
	j.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	return j
}
