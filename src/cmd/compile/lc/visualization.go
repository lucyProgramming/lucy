package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type Visualization interface {
	VisualNodes(filename string, nodes []*ast.TopNode)
	VisualPackage(p *ast.Package)
}
