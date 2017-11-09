package yacc

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
)

var (
	current_pos ast.Pos
)

func packageDefination(s *ast.Expression) {
	*ast.Nodes = append(*ast.Nodes, &ast.PackageNameDeclare{
		Name: s.Data.(string),
		Pos:  s.Pos,
	})
}

func importDefination(pname *ast.Expression, alias *ast.Expression) {
	*ast.Nodes = append(*ast.Nodes, &ast.Imports{
		Name:  pname.Data.(string),
		Alias: alias.Data.(string),
		Pos:   pname.Pos,
	})
}
