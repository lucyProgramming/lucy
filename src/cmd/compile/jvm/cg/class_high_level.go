package cg

import (
	"encoding/binary"
	"fmt"
	"math"
)

type ClassHighLevel struct {
	Class Class
	/*
		one class file can be compile form multi source file
	*/
	SourceFiles              map[string]struct{}
	Name                     string
	AccessFlags              uint16
	SuperClass               string
	Interfaces               []string
	Fields                   map[string]*FieldHighLevel
	Methods                  map[string][]*MethodHighLevel
	TriggerPackageInitMethod *MethodHighLevel
	TemplateFunctions        []*AttributeTemplateFunction
}

func (classHighLevel *ClassHighLevel) InsertMethodRefConst(mr CONSTANT_Methodref_info_high_level, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertMethodrefConst(mr))
}

/*
	new a method name,make sure it does exists before
*/
func (classHighLevel *ClassHighLevel) NewFunctionName(prefix string) string {
	if classHighLevel.Methods == nil || classHighLevel.Methods[prefix] == nil {
		return prefix
	}
	for i := 0; i < math.MaxInt16; i++ {
		name := fmt.Sprintf("%s$%d", prefix, i)
		if _, ok := classHighLevel.Methods[name]; ok == false {
			return name
		}
	}
	panic("names over flow") // this is not happening
}

func (classHighLevel *ClassHighLevel) NewFieldName(prefix string) string {
	if classHighLevel.Fields == nil || classHighLevel.Fields[prefix] == nil {
		return prefix
	}
	for i := 0; i < math.MaxInt16; i++ {
		name := fmt.Sprintf("%s_%d", prefix, i)
		if _, ok := classHighLevel.Fields[name]; ok == false {
			return name
		}
	}
	panic("names over flow") // this is not happening
}
func (classHighLevel *ClassHighLevel) InsertStringConst(s string, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertStringConst(s))
}

func (classHighLevel *ClassHighLevel) AppendMethod(ms ...*MethodHighLevel) {
	if classHighLevel.Methods == nil {
		classHighLevel.Methods = make(map[string][]*MethodHighLevel)
	}
	for _, v := range ms {
		if v.Name == "" {
			panic("null name")
		}
		if _, ok := classHighLevel.Methods[v.Name]; ok {
			classHighLevel.Methods[v.Name] = append(classHighLevel.Methods[v.Name], v)
		} else {
			classHighLevel.Methods[v.Name] = []*MethodHighLevel{v}
		}
	}
}

func (classHighLevel *ClassHighLevel) InsertInterfaceMethodrefConst(constant CONSTANT_InterfaceMethodref_info_high_level, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertInterfaceMethodrefConst(constant))
}

func (classHighLevel *ClassHighLevel) InsertMethodTypeConst(constant CONSTANT_MethodType_info_high_level, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertMethodTypeConst(constant))
}

func (classHighLevel *ClassHighLevel) InsertFieldRefConst(constant CONSTANT_Fieldref_info_high_level, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertFieldRefConst(constant))
}
func (classHighLevel *ClassHighLevel) InsertClassConst(className string, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertClassConst(className))
}
func (classHighLevel *ClassHighLevel) InsertIntConst(i int32, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertIntConst(i))
}

func (classHighLevel *ClassHighLevel) InsertLongConst(value int64, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertLongConst(value))
}

func (classHighLevel *ClassHighLevel) InsertFloatConst(value float32, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertFloatConst(value))
}

func (classHighLevel *ClassHighLevel) InsertDoubleConst(value float64, location []byte) {
	binary.BigEndian.PutUint16(location, classHighLevel.Class.InsertDoubleConst(value))
}

/*
	source files
*/
func (classHighLevel *ClassHighLevel) getSourceFile() string {
	if len(classHighLevel.SourceFiles) == 1 {
		for k, _ := range classHighLevel.SourceFiles {
			return k
		}
	}
	var s string
	if len(classHighLevel.SourceFiles) > 1 {
		s = "multi source compile into one class file,which are:\n"
	}
	prefix := ""
	if len(classHighLevel.SourceFiles) > 1 {
		prefix = "\t\t: "
	}

	for f, _ := range classHighLevel.SourceFiles {
		s += prefix + f
		s += "\n"
	}
	return s
}
