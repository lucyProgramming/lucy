package cg

type AttributeLucyConst struct {
}

func (a *AttributeLucyConst) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(ATTRIBUTE_NAME_LUCY_CONST)
	return ret
}
