package common

type Flags struct {
	OnlyImport        bool
	PackageName       string
	JvmMajorVersion   int
	JvmMinorVersion   int
	DisableCheckUnUse bool
}

var (
	CompileFlags Flags
)
