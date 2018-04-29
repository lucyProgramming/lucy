package cg

type AttributeClosureFunctionClass struct {
}

func (a *AttributeClosureFunctionClass) ToAttributeInfo(class *Class) *AttributeInfo {
	ret := &AttributeInfo{}
	ret.NameIndex = class.insertUtf8Const(ATTRIBUTE_NAME_LUCY_CLOSURE_FUNCTION_CLASS)
	return ret
}
