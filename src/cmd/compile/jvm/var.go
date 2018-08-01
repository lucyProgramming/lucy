package jvm

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"

var (
	ArrayMetas                 = map[ast.VariableTypeKind]*ArrayMeta{}
	typeConverter              TypeConverterAndPrimitivePacker
	Descriptor                 Description
	LucyMethodSignatureParser  LucyMethodSignature
	LucyFieldSignatureParser   LucyFieldSignature
	LucyTypeAliasParser        LucyTypeAlias
	FunctionDefaultValueParser FunctionDefaultValueParse
	closure                    Closure
)

type LeftValueKind int

const (
	_ LeftValueKind = iota
	LeftValueTypeLucyArray
	LeftValueTypeMap
	LeftValueTypeLocalVar
	LeftValueTypePutStatic
	LeftValueTypePutField
	LeftValueTypeArray
)

const (
	functionPointerInvokeMethod = "invoke"
	specialMethodInit           = "<init>"
	javaRootObjectArray         = "[Ljava/lang/Object;"
	javaStringClass             = "java/lang/String"
	javaMethodHandleClass       = "java/lang/invoke/MethodHandle"
	javaRootClass               = "java/lang/Object"
	javaMapClass                = "java/util/HashMap"
	javaIntegerClass            = "java/lang/Integer"
	javaFloatClass              = "java/lang/Float"
	javaDoubleClass             = "java/lang/Double"
	javaLongClass               = "java/lang/Long"
	javaStringBuilderClass      = "java/lang/StringBuilder"
	throwableClass              = "java/lang/Throwable"
	javaPrintStreamClass        = "java/io/PrintStream"
)
