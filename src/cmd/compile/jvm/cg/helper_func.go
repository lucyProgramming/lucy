package cg

import (
	"encoding/binary"
)

func backPatchIndex(locations [][]byte, index uint16) {
	for _, v := range locations {
		binary.BigEndian.PutUint16(v, index)
	}
}
