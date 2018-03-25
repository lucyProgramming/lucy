package main

import (
	"flag"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lc"
)

func main() {
	flag.BoolVar(&lc.CompileFlags.OnlyImport, "io", false, "only parse import package")
	flag.StringVar(&lc.CompileFlags.PackageName, "pn", "", "package name")
	flag.Parse()
	lc.Main(flag.Args())
}
