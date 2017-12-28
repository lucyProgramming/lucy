package cg

type ClassHighLevel struct {
	AccessFlags uint16
	Name        string
	SuperClass  string
	Interfaces  []string
	Fields      map[string]*FiledHighLevel
	Methods     map[string][]*MethodHighLevel
}

type FiledHighLevel struct {
	Name       string
	Descriptor string
	FieldInfo
}
type MethodHighLevel struct {
	Name       string
	Descriptor string
	MethodInfo
	Code AttributeCode
}
