package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	function    *ast.Function
	OutterClass *cg.ClassHighLevel
}
