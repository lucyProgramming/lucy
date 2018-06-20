package cg

type AttributeRuntimeVisibleAnnotation struct {
	Annotations []*Annotation
}

func (a *AttributeRuntimeVisibleAnnotation) ToAttributeInfo(class *Class) *AttributeInfo {
	if a == nil || len(a.Annotations) == 0 {
		return nil
	}
	ret := &AttributeInfo{}

	return ret
}
