package cg

type AttributeCompilerAuto struct {
}

func (a *AttributeCompilerAuto) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtf8Const(ATTRIBUTE_NAME_LUCY_COMPILTER_AUTO)
	return ret
}
