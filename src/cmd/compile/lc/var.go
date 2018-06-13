package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

var (
	compiler             LucyCompile
	loader               FileLoader
	ParseFunctionHandler func(bs []byte, pos *ast.Pos) (*ast.Function, []error)
)

func init() {
	ast.ImportsLoader = &loader
	loader.caches = make(map[string]interface{})
	ParseFunctionHandler = ast.ParseFunctionHandler
}

const (
	mainClassName = "main.class"
)
