package cg

type MethodHighLevel struct {
	IsConstruction                        bool
	Class                                 *ClassHighLevel
	Name                                  string
	Descriptor                            string
	AccessFlags                           uint16
	Code                                  *AttributeCode
	AttributeLucyMethodDescritor          *AttributeLucyMethodDescriptor
	AttributeLucyTriggerPackageInitMethod *AttributeLucyTriggerPackageInitMethod
	AttributeDefaultParameters            *AttributeDefaultParameters
	AttributeMethodParameters             *AttributeMethodParameters
	AttributeLucyReturnListNames          *AttributeMethodParameters
	AttributeCompilerAuto                 *AttributeCompilerAuto
}
