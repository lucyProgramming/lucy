package cg

type AttributeLucyEnum struct {
}

func (a *AttributeLucyEnum) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtf8Const(ATTRIBUTE_NAME_LUCY_ENUM)
	return ret
}
