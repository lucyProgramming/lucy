package cg

type AttributeCode struct {
	AttributeInfo
	maxStack             U2
	maxLocals            U2
	codeLength           U4
	codes                []byte
	exceptionTableLength U2
	exceptions           []*ExceptionTable
	attributeCounts      U2
	attributes           []*AttributeInfo
}

type ExceptionTable struct {
	startPc   U2
	endpc     U2
	handlerPc U2
	catchType U2
}
