package lc

type ResourceKind int

const (
	_ ResourceKind = iota
	ResourceKindJavaClass
	ResourceKindJavaPackage
	ResourceKindLucyClass
	ResourceKindLucyPackage
)

type Resource struct {
	kind     ResourceKind
	realPath string
	name     string
}
