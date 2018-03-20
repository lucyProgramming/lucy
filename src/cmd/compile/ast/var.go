package ast

import (
	"regexp"
)

type LoadName interface {
	LoadName(p *Package, packageName, name string) (interface{}, error)
}

var (
	MAIN_FUNCTION_NAME  = "main"
	THIS                = "this"
	NO_NAME_IDENTIFIER  = "_"
	Nodes               *[]*Node
	packageAliasReg     *regexp.Regexp
	NameLoader          LoadName
	LUCY_ROOT_CLASS     = "lucy/lang/Object"
	JAVA_ROOT_CLASS     = "java/lang/Object"
	buildinFunctionsMap = make(map[string]*buildFunction)
	JvmSlotSizeHandler  func(v *VariableType) uint16 // implements by outside
	PackageBeenCompile  *Package
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
