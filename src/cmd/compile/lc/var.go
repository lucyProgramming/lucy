package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

var (
	compiler             LucyCompile
	loader               RealNameLoader
	ParseFunctionHandler func(bs []byte, pos *ast.Pos) (*ast.Function, []error)
)

func init() {
	ast.NameLoader = &loader
	loader.caches = make(map[string]interface{})
	ParseFunctionHandler = ast.ParseFunctionHandler
}

const (
	mainClassName = "main.class"
)
