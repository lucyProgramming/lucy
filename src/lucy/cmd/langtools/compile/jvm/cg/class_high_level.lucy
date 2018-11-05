package cg

import (
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"
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

func (this *ClassHighLevel) InsertSourceFile(filename string) {
	if this.SourceFiles == nil {
		this.SourceFiles = make(map[string]struct{})
	}
	this.SourceFiles[filename] = struct{}{}
}

func (this *ClassHighLevel) InsertMethodRefConst(mr ConstantInfoMethodrefHighLevel,
	location []byte) {
	binary.BigEndian.PutUint16(location, this.Class.InsertMethodrefConst(mr))
}

func (this *ClassHighLevel) InsertMethodCall(code *AttributeCode, op byte,
	className, method, descriptor string) {
	code.Codes[code.CodeLength] = op
	this.InsertMethodRefConst(ConstantInfoMethodrefHighLevel{
		Class:      className,
		Method:     method,
		Descriptor: descriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
}

/*
	new a method name,make_node_objects sure it does exists before
*/
func (this *ClassHighLevel) NewMethodName(prefix string) string {
	if this.Methods == nil ||
		this.Methods[prefix] == nil {
		return prefix
	}
	for i := 0; i < math.MaxInt16; i++ {
		var name string
		if i == 0 {
			name = prefix
		} else {
			name = fmt.Sprintf("%s$%d", prefix, i)
		}
		if _, ok := this.Methods[name]; ok == false {
			return name
		}
	}
	panic("names over flow") // this is not happening
}

func (this *ClassHighLevel) InsertStringConst(s string, location []byte) {
	binary.BigEndian.PutUint16(location, this.Class.InsertStringConst(s))
}

func (this *ClassHighLevel) AppendMethod(ms ...*MethodHighLevel) {
	if this.Methods == nil {
		this.Methods = make(map[string][]*MethodHighLevel)
	}
	for _, v := range ms {
		if v.Name == "" {
			panic("null name")
		}
		if _, ok := this.Methods[v.Name]; ok {
			this.Methods[v.Name] = append(this.Methods[v.Name], v)
		} else {
			this.Methods[v.Name] = []*MethodHighLevel{v}
		}
	}
}

func (this *ClassHighLevel) InsertInterfaceMethodrefConst(
	constant ConstantInfoInterfaceMethodrefHighLevel,
	location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertInterfaceMethodrefConst(constant))
}

func (this *ClassHighLevel) InsertMethodTypeConst(constant ConstantInfoMethodTypeHighLevel,
	location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertMethodTypeConst(constant))
}

func (this *ClassHighLevel) InsertFieldRefConst(constant ConstantInfoFieldrefHighLevel,
	location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertFieldRefConst(constant))
}

func (this *ClassHighLevel) InsertClassConst(className string, location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertClassConst(className))
}
func (this *ClassHighLevel) InsertIntConst(i int32, location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertIntConst(i))
}

func (this *ClassHighLevel) InsertLongConst(value int64, location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertLongConst(value))
}

func (this *ClassHighLevel) InsertFloatConst(value float32, location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertFloatConst(value))
}

func (this *ClassHighLevel) InsertDoubleConst(value float64, location []byte) {
	binary.BigEndian.PutUint16(location,
		this.Class.InsertDoubleConst(value))
}

/*
	source files
*/
func (this *ClassHighLevel) getSourceFile() string {
	if len(this.SourceFiles) == 0 {
		return ""
	}
	if len(this.SourceFiles) == 1 {
		for k, _ := range this.SourceFiles {
			return k
		}
	}
	var s string
	for k, _ := range this.SourceFiles {
		s = filepath.Dir(k)
		break
	}
	s += "\\{"
	i := 0
	for f, _ := range this.SourceFiles {
		s += filepath.Base(f)
		if i != len(this.SourceFiles)-1 {
			s += ","
		}
		i++
	}
	s += "}"
	return s
}
