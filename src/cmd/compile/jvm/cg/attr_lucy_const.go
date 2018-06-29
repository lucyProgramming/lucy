package cg

type AttributeLucyConst struct {
}

func (a *AttributeLucyConst) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyConst)
	return ret
}
