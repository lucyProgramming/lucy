package ast

import (
	"fmt"
)

type LoadImport interface {
	LoadImport(importName string) (interface{}, error)
}

const (
	magicIdentifierFile     = "__FILE__"
	magicIdentifierLine     = "__LINE__"
	magicIdentifierTime     = "__TIME__"
	magicIdentifierClass    = "__CLASS__"
	magicIdentifierFunction = "__FUNCTION__"
	MainFunctionName        = "main"
	THIS                    = "this"
	UnderScore              = "_"
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
	return name == magicIdentifierFile ||
		name == magicIdentifierLine ||
		name == magicIdentifierTime ||
		name == magicIdentifierClass ||
		name == magicIdentifierFunction
}

var (
	ImportsLoader       LoadImport
	PackageBeenCompile  Package
	buildInFunctionsMap = make(map[string]*Function)
	lucyBuildInPackage  *Package
	// this function implemented by package parse , special for clone template function
	ParseFunctionHandler func(bs []byte, pos *Pos) (f *Function, es []error)
	javaStringClass      *Class
	LucyBytesType        *Type // []byte
	JavaBytesType        *Type // byte[]
)

func init() {
	LucyBytesType = &Type{
		Type: VariableTypeArray,
		Array: &Type{
			Type: VariableTypeByte,
		},
	}
	JavaBytesType = &Type{
		Type: VariableTypeJavaArray,
		Array: &Type{
			Type: VariableTypeByte,
		},
	}
}

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
			pos.ErrMsgPrefix(), JavaStringClass)
	}
}
