package ast

type BuildinFunctionPrintMeta struct {
	Stream *Expression
}

type BuildinFunctionPrintfMeta struct {
	BuildinFunctionPrintMeta
	Format     *Expression
	ArgsLength int
}

type BuildinFunctionSprintfMeta struct {
	Format     *Expression
	ArgsLength int
}
