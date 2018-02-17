package cg

import (
	"fmt"
	"math"
	"strings"
)

type ClassHighLevel struct {
	SourceFiles            map[string]struct{} // one class file can be compile form multi
	Name                   string
	IsClosureFunctionClass bool
	MainClass              *ClassHighLevel
	InnerClasss            []*ClassHighLevel
	AccessFlags            uint16
	IntConsts              map[int32][][]byte
	LongConsts             map[int64][][]byte
	FloatConsts            map[float32][][]byte
	DoubleConsts           map[float64][][]byte
	Classes                map[string][][]byte
	StringConsts           map[string][][]byte
	FieldRefs              map[CONSTANT_Fieldref_info_high_level][][]byte
	MethodRefs             map[CONSTANT_Methodref_info_high_level][][]byte
	NameAndTypes           map[CONSTANT_NameAndType_info_high_level][][]byte
	SuperClass             string
	Interfaces             []string
	Fields                 map[string]*FiledHighLevel
	Methods                map[string][]*MethodHighLevel
}

/*
	new a method name,mksure it does exists before
*/
func (c *ClassHighLevel) NewFunctionName(prefix string) string {
	for i := 0; i < math.MaxInt16; i++ {
		name := prefix + fmt.Sprintf("%d", i)
		if _, ok := c.Methods[name]; ok == false {
			return name
		}
	}
	panic(1)
}
func (c *ClassHighLevel) InsertStringConst(s string, location []byte) {
	if c.StringConsts == nil {
		c.StringConsts = make(map[string][][]byte)
	}
	if _, ok := c.StringConsts[s]; ok {
		c.StringConsts[s] = append(c.StringConsts[s], location)
	} else {
		c.StringConsts[s] = [][]byte{location}
	}
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

type CONSTANT_NameAndType_info_high_level struct {
	Name string
	Type string
}

func (c *ClassHighLevel) InsertNameAndTypeConst(nameAndType CONSTANT_NameAndType_info_high_level, location []byte) {
	if c.NameAndTypes == nil {
		c.NameAndTypes = make(map[CONSTANT_NameAndType_info_high_level][][]byte)
	}
	if _, ok := c.NameAndTypes[nameAndType]; ok {
		c.NameAndTypes[nameAndType] = append(c.NameAndTypes[nameAndType], location)
	} else {
		c.NameAndTypes[nameAndType] = [][]byte{location}
	}
}

type CONSTANT_Methodref_info_high_level struct {
	Class      string
	Name       string
	Descriptor string
}

func (c *ClassHighLevel) InsertMethodRefConst(mr CONSTANT_Methodref_info_high_level, location []byte) {
	if c.MethodRefs == nil {
		c.MethodRefs = make(map[CONSTANT_Methodref_info_high_level][][]byte)
	}
	if _, ok := c.MethodRefs[mr]; ok {
		c.MethodRefs[mr] = append(c.MethodRefs[mr], location)
	} else {
		c.MethodRefs[mr] = [][]byte{location}
	}
}

type CONSTANT_Fieldref_info_high_level struct {
	Class      string
	Name       string
	Descriptor string
}

func (c *ClassHighLevel) InsertFieldRefConst(fr CONSTANT_Fieldref_info_high_level, location []byte) {
	if c.FieldRefs == nil {
		c.FieldRefs = make(map[CONSTANT_Fieldref_info_high_level][][]byte)
	}
	if _, ok := c.FieldRefs[fr]; ok {
		c.FieldRefs[fr] = append(c.FieldRefs[fr], location)
	} else {
		c.FieldRefs[fr] = [][]byte{location}
	}
}
func (c *ClassHighLevel) InsertClassConst(classname string, location []byte) {
	if c.Classes == nil {
		c.Classes = make(map[string][][]byte)
	}
	if _, ok := c.Classes[classname]; ok {
		c.Classes[classname] = append(c.Classes[classname], location)
	} else {
		c.Classes[classname] = [][]byte{location}
	}
}
func (c *ClassHighLevel) InsertIntConst(i int32, location []byte) {
	if c.IntConsts == nil {
		c.IntConsts = make(map[int32][][]byte)
	}
	if _, ok := c.IntConsts[i]; ok {
		c.IntConsts[i] = append(c.IntConsts[i], location)
	} else {
		c.IntConsts[i] = [][]byte{location}
	}
}
func (c *ClassHighLevel) InsertLongConst(i int64, location []byte) {
	if c.LongConsts == nil {
		c.LongConsts = make(map[int64][][]byte)
	}
	if _, ok := c.LongConsts[i]; ok {
		c.LongConsts[i] = append(c.LongConsts[i], location)
	} else {
		c.LongConsts[i] = [][]byte{location}
	}
}

func (c *ClassHighLevel) InsertFloatConst(f float32, location []byte) {
	if c.FloatConsts == nil {
		c.FloatConsts = make(map[float32][][]byte)
	}
	if _, ok := c.FloatConsts[f]; ok {
		c.FloatConsts[f] = append(c.FloatConsts[f], location)
	} else {
		c.FloatConsts[f] = [][]byte{location}
	}
}

func (c *ClassHighLevel) InsertDoubleConst(d float64, location []byte) {
	if c.DoubleConsts == nil {
		c.DoubleConsts = make(map[float64][][]byte)
	}
	if _, ok := c.DoubleConsts[d]; ok {
		c.DoubleConsts[d] = append(c.DoubleConsts[d], location)
	} else {
		c.DoubleConsts[d] = [][]byte{location}
	}
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
