package lc

const (
	_ = iota
	RESOURCE_KIND_JAVA_CLASS
	RESOURCE_KIND_JAVA_PACKAGE
	RESOURCE_KIND_LUCY_CLASS
	RESOURCE_KIND_LUCY_PACKAGE
)

type Resource struct {
	kind     int
	realPath string
	name     string
}
