package lc

type resourceKind int

const (
	_ resourceKind = iota
	resourceKindJavaClass
	resourceKindJavaPackage
	resourceKindLucyClass
	resourceKindLucyPackage
)

type Resource struct {
	kind     resourceKind
	realPath string
	name     string
}
