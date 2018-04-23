package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"path/filepath"
)

var (
	CompileFlags Flags
	compiler     LucyCompile
	loader       RealNameLoader
)

type Flags struct {
	OnlyImport  bool
	PackageName string
}

func init() {
	ast.NameLoader = &loader
}

func classShortName(binaryName string) (shortName string) {
	shortName = filepath.Base(binaryName)
	return
}

const (
	mainClassName = "main.class"
)
