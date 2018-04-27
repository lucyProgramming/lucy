package ast

type BuildinPrintMeta struct {
	FirstParameterIsStream bool
}

type BuildinPrintfMeta struct {
	BuildinPrintMeta
	ArgsLength int
}

type BuildinSprintMeta struct {
	ArgsLength int
}
