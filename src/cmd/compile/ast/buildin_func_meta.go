package ast

type BuildInFunctionPrintfMeta struct {
	Format     *Expression
	ArgsLength int
}

type BuildInFunctionSprintfMeta struct {
	Format     *Expression
	ArgsLength int
}
