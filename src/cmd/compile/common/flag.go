package common

type Flags struct {
	OnlyImport  bool
	PackageName string
	JvmVersion  int
}

var (
	CompileFlags Flags
)
