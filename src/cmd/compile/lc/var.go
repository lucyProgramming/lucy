package lc

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
)

var (
	Tops         = make([]*ast.Node, 0)
	CompileFlags Flags
)

type Flags struct {
	OnlyImport bool
}
