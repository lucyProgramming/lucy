package cg

type MethodHighLevel struct {
	Class                          *ClassHighLevel
	Name                           string
	Descriptor                     string
	AccessFlags                    uint16
	Code                           AttributeCode
	AttributeLucyInnerStaticMethod *AttributeLucyInnerStaticMethod
}
