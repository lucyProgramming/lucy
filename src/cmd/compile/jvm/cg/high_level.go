package cg

type CONSTANT_NameAndType_info_high_level struct {
	Name       string
	Descriptor string
}

type CONSTANT_Methodref_info_high_level struct {
	Class      string
	Method     string
	Descriptor string
}

type CONSTANT_InterfaceMethodref_info_high_level struct {
	Class      string
	Method     string
	Descriptor string
}

type CONSTANT_Fieldref_info_high_level struct {
	Class      string
	Field      string
	Descriptor string
}
