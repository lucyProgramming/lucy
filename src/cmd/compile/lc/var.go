package lc

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/yacc"
)

var (
	Tops = make([]*ast.Node, 0)
	ast.Nodes = &Tops
)

