package lc

import (
	"encoding/binary"
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type ClassDecoder struct {
	bs []byte
}

func (c *ClassDecoder) parseConstPool() error {
	length := binary.BigEndian.Uint16(c.bs)
	c.bs = c.bs[2:]
	for i := 0; i < int(length)-1; i++ {

	}
	return nil
}

func (c *ClassDecoder) decode(bs []byte) (*cg.Class, error) {
	c.bs = bs
	if binary.BigEndian.Uint32(bs) != cg.CLASS_MAGIC {
		return nil, fmt.Errorf("magic number is not right")
	}
	c.bs = c.bs[4:]
	ret := &cg.Class{}
	ret.ConstPool = []*cg.ConstPool{nil} // pool start 1
	if err := c.parseConstPool(); err != nil {
		return ret, err
	}
	return ret, nil
}
