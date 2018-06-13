package ast

type BuildInFunctionPrintMeta struct {
	Stream *Expression
}

type BuildInFunctionPrintfMeta struct {
	BuildInFunctionPrintMeta
	Format     *Expression
	ArgsLength int
}

type BuildInFunctionSprintfMeta struct {
	Format     *Expression
	ArgsLength int
}
