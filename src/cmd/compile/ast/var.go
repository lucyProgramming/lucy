package ast

import (
	"fmt"
)

type LoadImport interface {
	LoadImport(importName string) (interface{}, error)
}

const (
	MagicIdentifierFile     = "__FILE__"
	MagicIdentifierLine     = "__LINE__"
	MagicIdentifierTime     = "__TIME__"
	MagicIdentifierClass    = "__CLASS__"
	MagicIdentifierFunction = "__FUNCTION__"
	MainFunctionName        = "main"
	THIS                    = "this"
	NoNameIdentifier        = "_"
	LucyRootClass           = "lucy/lang/Lucy"
	JavaRootClass           = "java/lang/Object"
	DefaultExceptionClass   = "java/lang/Exception"
	JavaThrowableClass      = "java/lang/Throwable"
	JavaStringClass         = "java/lang/String"
	SUPER                   = "super"
	SpecialMethodInit       = "<init>"
	classInitMethod         = "<clinit>"
)

func isMagicIdentifier(name string) bool {
	return name == MagicIdentifierFile ||
		name == MagicIdentifierLine ||
		name == MagicIdentifierTime ||
		name == MagicIdentifierClass ||
		name == MagicIdentifierFunction
}

var (
	ImportsLoader        LoadImport
	PackageBeenCompile   Package
	buildInFunctionsMap  = make(map[string]*Function)
	lucyBuildInPackage   *Package
	ParseFunctionHandler func(bs []byte, pos *Pos) (f *Function, es []error)
	javaStringClass      *Class
)

func loadJavaStringClass(pos *Pos) error {
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
