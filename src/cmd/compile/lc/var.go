package lc

var (
	CompileFlags Flags
	compiler     LucyCompile
)

type Flags struct {
	OnlyImport bool
}
