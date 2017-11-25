package cg

type AttributeInfo struct {
	attributeIndex  U2
	attributeLength U4
	info            []byte
}

type ConstantValue_attribute struct {
	AttributeInfo
	constvalueIndex U2
}

type AttributeSignature struct {
	AttributeInfo
	index U2
}

type AttributeSourceFile struct {
	AttributeInfo
	index U2
}

type AttributeLineNumber struct {
	AttributeInfo
	length      U2
	linenumbers []*AttributeLinePc
}

type AttributeLinePc struct {
	startPc    U2
	lineNumber U2
}

type AttributeBootstrapMethods struct {
	AttributeInfo
	numBootStrapMethods U2
}

type BootStrapMethod struct {
	bootStrapMethodRef         U2
	numBootStrapMethodArgument U2
	bootStrapMethodArguments   []U2
}
