package cg

type AttributeCompilerAuto struct {
}

func (a *AttributeCompilerAuto) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyCompilerAuto)
	return ret
}
