package cg

type AttributeRuntimeVisibleAnnotation struct {
	Annotations []*Annotation
}

type Annotation struct {
	Type              string
	ElementValuePairs []*ElementValuePair
}
type ElementValuePair struct {
	ElementName string
	Value       ElementValue
}

type ElementValue struct {
	Tag byte
}