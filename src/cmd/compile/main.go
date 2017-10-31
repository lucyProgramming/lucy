package main

import (
	"github.com/756445638/lucy/src/cmd/compile/lc"
)

func main() {
	flag.Parse()
	lc.Compile(flag.Args())
}
