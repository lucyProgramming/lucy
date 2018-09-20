package cg

type FieldHighLevel struct {
	Name                         string
	Descriptor                   string
	AccessFlags                  uint16
	AttributeConstantValue       *AttributeConstantValue
	AttributeLucyFieldDescriptor *AttributeLucyFieldDescriptor
	AttributeLucyConst           *AttributeLucyConst
	AttributeLucyComment         *AttributeLucyComment
}
