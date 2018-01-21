package ast

import (
	"regexp"
)

type PackageLoader interface {
	LoadPackage(name string) (*Package, error)
}

var (
	THIS                    = "this"
	SIGNAL_UNERDERLINE_NAME = "_"
	Nodes                   *[]*Node //
	packageAliasReg         *regexp.Regexp
	PackageLoad             PackageLoader
	ROOT_CLASS              = "lucy/lang/LucyObject"
	buildinFunctionsMap     = make(map[string]*buildFunction)
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
