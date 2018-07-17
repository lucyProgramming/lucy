package lc

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

var (
	compiler LucyCompile
	loader   FileLoader
)

func init() {
	ast.ImportsLoader = &loader
	loader.caches = make(map[string]interface{})

}

const (
	mainClassName = "main.class"
)
