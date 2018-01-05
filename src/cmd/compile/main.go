package main

import (
	"flag"
	"github.com/756445638/lucy/src/cmd/compile/lc"
)

func main() {
	flag.BoolVar(&lc.CompileFlags.OnlyImport, "io", false, "only parse import package")
	flag.Parse()
	lc.Main(flag.Args())
}
