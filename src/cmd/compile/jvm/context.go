package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	//"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type Context struct {
	Vars     map[string]*ast.VariableDefinition
	Receiver *ast.VariableType
	Function *ast.Function
	Class    *ast.Class
}
