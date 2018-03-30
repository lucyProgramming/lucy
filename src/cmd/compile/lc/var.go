package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"path/filepath"
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

func classShortName(binaryName string) (shortName string) {
	shortName = filepath.Base(binaryName)
	return
}
