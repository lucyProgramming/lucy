package cg

import (
	"encoding/binary"
	"io"
)

func (c *Class) OutPut(dest io.Writer) error {
	c.destination = dest
	//magic number
	_, err := dest.Write([]byte{0xca, 0xfe, 0xba, 0xbe})
	if err != nil {
		return err
	}
	// minorversion
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(c.MinorVersion))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	// major version
	binary.BigEndian.PutUint16(bs2, uint16(c.MajorVersion))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	//const pool
	binary.BigEndian.PutUint16(bs2, c.constPoolUint16Length())
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.ConstPool {
		if v == nil {
			continue
		}
		_, err = dest.Write([]byte{byte(v.Tag)})
		if err != nil {
			return err
		}
		_, err = dest.Write(v.Info)
		if err != nil {
			return err
		}
	}
	//access flag
	binary.BigEndian.PutUint16(bs2, uint16(c.AccessFlag))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(bs2, c.ThisClass)
	//this class
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	//super Class
	binary.BigEndian.PutUint16(bs2, c.SuperClass)
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	// interface
	binary.BigEndian.PutUint16(bs2, uint16(len(c.Interfaces)))
	_, err = dest.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.Interfaces {
		binary.BigEndian.PutUint16(bs2, uint16(v))
		_, err = dest.Write(bs2)
		if err != nil {
			return err
		}
	}
	err = c.writeFields()
	if err != nil {
		return err
	}
	//methods
	err = c.writeMethods()
	if err != nil {
		return err
	}
	// attribute

	return c.writeAttributeInfo(c.Attributes)
}

func (c *Class) writeMethods() error {
	var err error
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(c.Methods)))
	_, err = c.destination.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.Methods {
		binary.BigEndian.PutUint16(bs2, uint16(v.AccessFlags))
		_, err = c.destination.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.NameIndex)
		_, err = c.destination.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.DescriptorIndex)
		_, err = c.destination.Write(bs2)
		if err != nil {
			return err
		}

		err = c.writeAttributeInfo(v.Attributes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Class) writeFields() error {
	var err error
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(c.Fields)))
	_, err = c.destination.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range c.Fields {
		binary.BigEndian.PutUint16(bs2, uint16(v.AccessFlags))
		_, err = c.destination.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.NameIndex)
		_, err = c.destination.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.DescriptorIndex)
		_, err = c.destination.Write(bs2)
		if err != nil {
			return err
		}

		err = c.writeAttributeInfo(v.Attributes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Class) writeAttributeInfo(as []*AttributeInfo) error {
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(as)))
	_, err := c.destination.Write(bs2)
	if err != nil {
		return err
	}
	bs4 := make([]byte, 4)
	for _, v := range as {
		binary.BigEndian.PutUint16(bs2, v.NameIndex)
		_, err = c.destination.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint32(bs4, v.attributeLength)
		_, err = c.destination.Write(bs4)
		if err != nil {
			return err
		}
		_, err = c.destination.Write(v.Info)
		if err != nil {
			return err
		}
	}
	return nil
}
