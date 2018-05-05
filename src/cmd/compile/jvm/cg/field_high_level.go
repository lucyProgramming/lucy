package cg

type FieldHighLevel struct {
	Name                        string
	Descriptor                  string
	AccessFlags                 uint16
	AttributeConstantValue      *AttributeConstantValue
	AttributeLucyFieldDescritor *AttributeLucyFieldDescriptor
	AttributeLucyConst          *AttributeLucyConst
}
