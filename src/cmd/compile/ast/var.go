package ast

import (
	"regexp"
)

var (
	small_float          = 0.0001
	negative_small_float = -small_float
	Nodes                *[]*Node //
	packageAliasReg      *regexp.Regexp
)

type NameWithPos struct {
	Name string
	Pos  *Pos
}

func init() {
	var err error
	packageAliasReg, err = regexp.Compile(`^[a-zA-Z][[a-zA-Z1-9\_]+$`)
	if err != nil {
		panic(err)
	}
}
