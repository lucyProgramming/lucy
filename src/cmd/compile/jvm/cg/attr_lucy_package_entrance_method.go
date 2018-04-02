package cg

type AttributeLucyPackageEntranceMethod struct {
}

func (a *AttributeLucyPackageEntranceMethod) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtfConst(ATTRIBUTE_NAME_LUCY_PACKAGE_ENTRANCE_METHOD)
	return ret
}
