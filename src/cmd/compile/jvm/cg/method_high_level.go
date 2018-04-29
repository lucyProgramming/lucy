package cg

type MethodHighLevel struct {
	IsConstruction                        bool
	Class                                 *ClassHighLevel
	Name                                  string
	Descriptor                            string
	AccessFlags                           uint16
	Code                                  *AttributeCode
	AttributeLucyInnerStaticMethod        *AttributeLucyInnerStaticMethod
	AttributeLucyMethodDescritor          *AttributeLucyMethodDescriptor
	AttributeLucyTriggerPackageInitMethod *AttributeLucyTriggerPackageInitMethod
	AttributeDefaultParameters            *AttributeDefaultParameters
}
