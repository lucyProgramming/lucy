package cg

type AttributeLucyLucyMainClass struct {
}

func (a *AttributeLucyLucyMainClass) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtfConst(ATTRIBUTE_NAME_LUCY_LUCY_MAIN_CLASS)
	return ret
}
