package cg

type AttributeLucyEnum struct {
}

func (this *AttributeLucyEnum) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyEnum)
	return ret
}
