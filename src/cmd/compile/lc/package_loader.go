package lc

const (
	_ = iota
	ResourceKindJavaClass
	ResourceKindJavaPackage
	ResourceKindLucyClass
	ResourceKindLucyPackage
)

type Resource struct {
	kind     int
	realPath string
	name     string
}
