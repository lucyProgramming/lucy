package cg

type FieldHighLevel struct {
	Name                        string
	Descriptor                  string
	AccessFlags                 uint16
	ConstantValue               *AttributeConstantValue
	AttributeLucyFieldDescritor *AttributeLucyFieldDescriptor
	AttributeLucyConst          *AttributeLucyConst
}
