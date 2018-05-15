package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

var (
	CompileFlags         Flags
	compiler             LucyCompile
	loader               RealNameLoader
	ParseFunctionHandler func(bs []byte, pos *ast.Pos) (*ast.Function, []error)
)

type Flags struct {
	OnlyImport  bool
	PackageName string
	JvmVersion  int
}

func init() {
	ast.NameLoader = &loader
	loader.caches = make(map[string]interface{})
	ParseFunctionHandler = ast.ParseFunctionHandler
}

const (
	mainClassName = "main.class"
)
