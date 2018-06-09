package common

type Flags struct {
	OnlyImport                 bool
	PackageName                string
	JvmVersion                 int
	DisableCheckUnUsedVariable bool //
}

var (
	CompileFlags Flags
)
