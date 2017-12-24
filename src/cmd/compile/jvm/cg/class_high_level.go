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
	Descriptor string
	FieldInfo
}
type MethodHighLevel struct {
	Descriptor string
	MethodInfo
}
