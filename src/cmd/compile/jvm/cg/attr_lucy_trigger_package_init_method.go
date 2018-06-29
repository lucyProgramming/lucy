package cg

type AttributeLucyTriggerPackageInitMethod struct {
}

func (a *AttributeLucyTriggerPackageInitMethod) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.InsertUtf8Const(AttributeNameLucyTriggerPackageInit)
	return ret
}
