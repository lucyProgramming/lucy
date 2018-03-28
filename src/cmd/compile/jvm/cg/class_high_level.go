package cg

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

type ClassHighLevel struct {
	Class                  Class
	SourceFiles            map[string]struct{} // one class file can be compile form multi
	Name                   string
	IsClosureFunctionClass bool
	MainClass              *ClassHighLevel
	InnerClasss            []*ClassHighLevel
	AccessFlags            uint16
	SuperClass             string
	Interfaces             []string
	Fields                 map[string]*FiledHighLevel
	Methods                map[string][]*MethodHighLevel
}

type CONSTANT_NameAndType_info_high_level struct {
	Name       string
	Descriptor string
}

type CONSTANT_Methodref_info_high_level struct {
	Class      string
	Name       string
	Descriptor string
}

type CONSTANT_InterfaceMethodref_info_high_level struct {
	Class      string
	Name       string
	Descriptor string
}

func (c *ClassHighLevel) InsertMethodRefConst(mr CONSTANT_Methodref_info_high_level, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertMethodrefConst(mr))
}

type CONSTANT_Fieldref_info_high_level struct {
	Class      string
	Name       string
	Descriptor string
}

/*
	new a method name,mksure it does exists before
*/
func (c *ClassHighLevel) NewFunctionName(prefix string) string {
	if c.Methods == nil && c.Methods[prefix] == nil {
		return prefix
	}
	for i := 0; i < math.MaxInt16; i++ {
		name := prefix + fmt.Sprintf("%d", i)
		if _, ok := c.Methods[name]; ok == false {
			return name
		}
	}
	panic("names over flow")
}
func (c *ClassHighLevel) InsertStringConst(s string, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertStringConst(s))
}

func (c *ClassHighLevel) AppendMethod(ms ...*MethodHighLevel) {
	if c.Methods == nil {
		c.Methods = make(map[string][]*MethodHighLevel)
	}
	for _, v := range ms {
		if v.Name == "" {
			panic(1)
		}
		if _, ok := c.Methods[v.Name]; ok {
			c.Methods[v.Name] = append(c.Methods[v.Name], v)
		} else {
			c.Methods[v.Name] = []*MethodHighLevel{v}
		}
	}
}

func (c *ClassHighLevel) InsertInterfaceMethodrefConst(fr CONSTANT_InterfaceMethodref_info_high_level, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertInterfaceMethodrefConst(fr))
}
func (c *ClassHighLevel) InsertFieldRefConst(fr CONSTANT_Fieldref_info_high_level, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertFieldRefConst(fr))
}
func (c *ClassHighLevel) InsertClassConst(classname string, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertClassConst(classname))
}
func (c *ClassHighLevel) InsertIntConst(i int32, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertIntConst(i))
}
func (c *ClassHighLevel) InsertLongConst(i int64, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertLongConst(i))
}

func (c *ClassHighLevel) InsertFloatConst(f float32, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertFloatConst(f))
}

func (c *ClassHighLevel) InsertDoubleConst(d float64, location []byte) {
	binary.BigEndian.PutUint16(location, c.Class.InsertDoubleConst(d))
}
func (c *ClassHighLevel) getSourceFile() string {
	s := ""
	for f, _ := range c.SourceFiles {
		s += f + " "
	}
	if s != "" {
		s = strings.TrimRight(s, " ")
	}
	return s
}
