package ast

import (
	"fmt"
	"regexp"
)

type LoadImport interface {
	LoadImport(importName string) (interface{}, error)
}

const (
	MainFunctionName      = "main"
	THIS                  = "this"
	NoNameIdentifier      = "_"
	LucyRootClass         = "lucy/lang/Lucy"
	JavaRootClass         = "java/lang/Object"
	DefaultExceptionClass = "java/lang/Exception"
	JavaThrowableClass    = "java/lang/Throwable"
	JavaStringClass       = "java/lang/String"
	SUPER                 = "super"
	SpecialMethodInit     = "<init>"
	ClassInitMethod       = "<clinit>"
)

var (
	packageAccessNameReg *regexp.Regexp
	ImportsLoader        LoadImport
	PackageBeenCompile   Package
	buildInFunctionsMap  = make(map[string]*Function)
	lucyBuildInPackage   *Package
	ParseFunctionHandler func(bs []byte, pos *Position) (*Function, []error)
	javaStringClass      *Class
)

func init() {
	var err error
	packageAccessNameReg, err = regexp.Compile(`^[a-zA-Z][[a-zA-Z1-9\_]+$`)
	if err != nil {
		panic(err)
	}
}

func loadJavaStringClass(pos *Position) error {
	if javaStringClass != nil {
		return nil
	}
	c, err := ImportsLoader.LoadImport(JavaStringClass)
	if err != nil {
		return err
	}
	if cc, ok := c.(*Class); ok && cc != nil {
		javaStringClass = cc
		return nil
	} else {
		return fmt.Errorf("%s '%s' is not class",
			errMsgPrefix(pos), JavaStringClass)
	}
}
