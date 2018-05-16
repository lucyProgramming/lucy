package common

type Flags struct {
	OnlyImport                 bool
	PackageName                string
	JvmVersion                 int
	DisAbleCheckUnUsedVariable bool
}

var (
	CompileFlags Flags
)
