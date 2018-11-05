package cg

type AttributeLucyConst struct {
}

func (this *AttributeLucyConst) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyConst)
	return ret
}
