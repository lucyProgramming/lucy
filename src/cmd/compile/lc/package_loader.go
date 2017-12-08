package lc

import "github.com/756445638/lucy/src/cmd/compile/ast"

type PackageLoader struct {
}

func (*PackageLoader) LoadPackage(name string) (*ast.Package, error) {
	return nil, nil
}

func init() {
	ast.PackageLoad = &PackageLoader{}
}
