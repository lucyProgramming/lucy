package ast

import (
	"regexp"
)

type LoadName interface {
	LoadName(resouceName string) (*Package, interface{}, error)
}

const (
	MAIN_FUNCTION_NAME       = "main"
	THIS                     = "this"
	NO_NAME_IDENTIFIER       = "_"
	LUCY_ROOT_CLASS          = "lucy/deps/Object"
	JAVA_ROOT_CLASS          = "java/lang/Object"
	DEFAULT_EXCEPTION_CLASS  = "java/lang/Exception"
	JAVA_THROWABLE_CLASS     = "java/lang/Throwable"
	SUPER_FIELD_NAME         = "super"
	CONSTRUCTION_METHOD_NAME = "<init>"
)

var (
	Nodes                  *[]*Node
	packageAliasReg        *regexp.Regexp
	NameLoader             LoadName
	PackageBeenCompile     Package
	buildinFunctionsMap    = make(map[string]*Function)
	lucyLangBuildinPackage *Package
	ParseFunctionHandler   func(bs []byte, pos *Pos) (*Function, []error)
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
