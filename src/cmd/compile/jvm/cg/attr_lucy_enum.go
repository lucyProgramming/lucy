package cg

type AttributeLucyEnum struct {
}

func (a *AttributeLucyEnum) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyEnum)
	return ret
}
