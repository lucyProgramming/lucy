package cg

type ConstantInfoNameAndTypeHighLevel struct {
	Name       string
	Descriptor string
}

type ConstantInfoMethodrefHighLevel struct {
	Class      string
	Method     string
	Descriptor string
}

type ConstantInfoInterfaceMethodrefHighLevel struct {
	Class      string
	Method     string
	Descriptor string
}

type ConstantInfoFieldrefHighLevel struct {
	Class      string
	Field      string
	Descriptor string
}

type ConstantInfoMethodTypeHighLevel struct {
	Descriptor string
}
