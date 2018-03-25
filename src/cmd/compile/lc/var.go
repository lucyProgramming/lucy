package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
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
