package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
)

type Context struct {
	Vars map[string]*ast.VariableDefinition
}
