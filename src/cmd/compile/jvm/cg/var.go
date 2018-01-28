package cg

import (
	"encoding/binary"
)

const (
	CONST_POOL_MAX_SIZE = 65535 // const pool index begin at 1
)

func backPatchIndex(locations [][]byte, index uint16) {
	for _, v := range locations {
		binary.BigEndian.PutUint16(v, index)
	}

}
