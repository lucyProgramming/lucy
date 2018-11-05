package cg

import (
	"encoding/binary"
	"io"
)

func (this *Class) OutPut(writer io.Writer) error {
	this.writer = writer
	//magic number
	_, err := writer.Write([]byte{0xca, 0xfe, 0xba, 0xbe})
	if err != nil {
		return err
	}
	// minor version
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(this.MinorVersion))
	_, err = writer.Write(bs2)
	if err != nil {
		return err
	}
	// major version
	binary.BigEndian.PutUint16(bs2, uint16(this.MajorVersion))
	_, err = writer.Write(bs2)
	if err != nil {
		return err
	}
	//const pool
	binary.BigEndian.PutUint16(bs2, this.constPoolUint16Length())
	_, err = writer.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range this.ConstPool {
		if v == nil {
			continue
		}
		_, err = writer.Write([]byte{byte(v.Tag)})
		if err != nil {
			return err
		}
		_, err = writer.Write(v.Info)
		if err != nil {
			return err
		}
	}
	//access flag
	binary.BigEndian.PutUint16(bs2, uint16(this.AccessFlag))
	_, err = writer.Write(bs2)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(bs2, this.ThisClass)
	//this class
	_, err = writer.Write(bs2)
	if err != nil {
		return err
	}
	//super Class
	binary.BigEndian.PutUint16(bs2, this.SuperClass)
	_, err = writer.Write(bs2)
	if err != nil {
		return err
	}
	// interface
	binary.BigEndian.PutUint16(bs2, uint16(len(this.Interfaces)))
	_, err = writer.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range this.Interfaces {
		binary.BigEndian.PutUint16(bs2, uint16(v))
		_, err = writer.Write(bs2)
		if err != nil {
			return err
		}
	}
	err = this.writeFields()
	if err != nil {
		return err
	}
	//methods
	err = this.writeMethods()
	if err != nil {
		return err
	}
	// attribute

	return this.writeAttributeInfo(this.Attributes)
}

func (this *Class) writeMethods() error {
	var err error
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(this.Methods)))
	_, err = this.writer.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range this.Methods {
		binary.BigEndian.PutUint16(bs2, uint16(v.AccessFlags))
		_, err = this.writer.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.NameIndex)
		_, err = this.writer.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.DescriptorIndex)
		_, err = this.writer.Write(bs2)
		if err != nil {
			return err
		}

		err = this.writeAttributeInfo(v.Attributes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Class) writeFields() error {
	var err error
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(this.Fields)))
	_, err = this.writer.Write(bs2)
	if err != nil {
		return err
	}
	for _, v := range this.Fields {
		binary.BigEndian.PutUint16(bs2, uint16(v.AccessFlags))
		_, err = this.writer.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.NameIndex)
		_, err = this.writer.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(bs2, v.DescriptorIndex)
		_, err = this.writer.Write(bs2)
		if err != nil {
			return err
		}

		err = this.writeAttributeInfo(v.Attributes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Class) writeAttributeInfo(as []*AttributeInfo) error {
	bs2 := make([]byte, 2)
	binary.BigEndian.PutUint16(bs2, uint16(len(as)))
	_, err := this.writer.Write(bs2)
	if err != nil {
		return err
	}
	bs4 := make([]byte, 4)
	for _, v := range as {
		binary.BigEndian.PutUint16(bs2, v.NameIndex)
		_, err = this.writer.Write(bs2)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint32(bs4, v.attributeLength)
		_, err = this.writer.Write(bs4)
		if err != nil {
			return err
		}
		_, err = this.writer.Write(v.Info)
		if err != nil {
			return err
		}
	}
	return nil
}
