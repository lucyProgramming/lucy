package ast

import (
	"fmt"
	"regexp"
)

type LoadName interface {
	LoadName(resouceName string) (interface{}, error)
}

const (
	MAIN_FUNCTION_NAME       = "main"
	THIS                     = "this"
	NO_NAME_IDENTIFIER       = "_"
	LUCY_ROOT_CLASS          = "lucy/lang/Object"
	JAVA_ROOT_CLASS          = "java/lang/Object"
	DEFAULT_EXCEPTION_CLASS  = "java/lang/Exception"
	JAVA_THROWABLE_CLASS     = "java/lang/Throwable"
	JAVA_STRING_CLASS        = "java/lang/String"
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
	javaStringClass        *Class
)

func loadJavaStringClass(pos *Pos) error {
	if javaStringClass != nil {
		return nil
	}
	c, err := NameLoader.LoadName(JAVA_STRING_CLASS)
	if err != nil {
		return err
	}
	if cc, ok := c.(*Class); ok && cc != nil {
		javaStringClass = cc
		return nil
	} else {
		return fmt.Errorf("%s '%s' is not class", errMsgPrefix(pos), JAVA_STRING_CLASS)
	}
}

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
