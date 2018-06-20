package cg

type AttributeCompilerAuto struct {
}

func (a *AttributeCompilerAuto) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(ATTRIBUTE_NAME_LUCY_COMPILER_AUTO)
	return ret
}
