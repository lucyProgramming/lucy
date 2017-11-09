package main

import (
	"flag"

	"github.com/756445638/lucy/src/cmd/compile/lc"
)

func main() {
	flag.Parse()
	lc.Main(flag.Args())
}
