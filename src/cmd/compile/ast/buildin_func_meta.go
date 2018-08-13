package ast

type BuildInFunctionPrintfMeta struct {
	Format *Expression
	Length int
}

type BuildInFunctionSprintfMeta struct {
	Format     *Expression
	ArgsLength int
}
