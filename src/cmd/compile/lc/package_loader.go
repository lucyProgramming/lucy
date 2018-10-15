package lc

type ResourceKind int

const (
	_ ResourceKind = iota
	resourceKindJavaClass
	resourceKindJavaPackage
	resourceKindLucyClass
	resourceKindLucyPackage
)

type Resource struct {
	kind     ResourceKind
	realPath string
	name     string
}
