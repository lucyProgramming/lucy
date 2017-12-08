package classJson

type ClassJson struct {
	Magic        string    `json:"magic"`
	MinorVersion uint16    `json:"minorVersion"`
	MajorVersion uint16    `json:"majorVersion"`
	AccessFlags  uint16    `json:"access_flags"`
	ThisClass    string    `json:"this_class"`
	SuperClass   string    `json:"super_class"`
	Fields       []*Field  `json:"fields"`
	Methods      []*Method `json:"method"`
}

type Field struct {
	Name        string          `json:"name"`
	AccessFlags uint16          `json:"access_flags"`
	Descriptor  string          `json:"descriptor"`
	Signature   *FieldSignature `json:"signature"`
}

type FieldSignature FieldTypeSingture

type FieldTypeSingture struct {
	Kind       string              `json:"kind"`
	Identifier string              `json:"identifier"`
	ArrayType  *TypeSignature      `json:"array_type"`
	ClassType  *ClassTypeSignature `json:"class_type"`
}

type ClassTypeSignature struct {
	Name           string               `json:"name"`
	TypedArguments []*FieldTypeSingture `json:""`
}

type TypeSignature FieldTypeSingture

type Method struct {
	Name        string           `json:"name"`
	AccessFlags uint16           `json:"access_flags"`
	Typ         *MethodTyp       `json:"typ"`
	Signature   *MethodSignature `json:"signature"`
}
type MethodSignature struct {
	TypedParameters []*FormalTypeParameter `json:"typed_parameters"`
	Parameters      []*TypeSignature       `json:"parameters"`
	Returns         []*TypeSignature       `json:"returns"`
}

type FormalTypeParameter struct {
	Name       string            `json:"name"`
	Super      *ClassSignature   `json:"super"`
	Interfaces []*ClassSignature `json:"interfaces"`
}

type MethodTyp struct {
	Parameters []string
	Return     string
}

type MethodSignature struct {
}

type ClassSignature struct {
}
