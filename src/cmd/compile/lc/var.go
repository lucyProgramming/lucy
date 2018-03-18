package lc

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
)

var (
	CompileFlags Flags
	compiler     LucyCompile
)

type Flags struct {
	OnlyImport  bool
	PackageName string
}

func init() {
	ast.NameLoader = &RealNameLoader{}
}
