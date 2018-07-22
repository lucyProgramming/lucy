package jvm

var (
	ArrayMetas                 = map[int]*ArrayMeta{}
	typeConverter              TypeConverterAndPrimitivePacker
	Descriptor                 Description
	LucyMethodSignatureParser  LucyMethodSignature
	LucyFieldSignatureParser   LucyFieldSignature
	LucyTypeAliasParser        LucyTypeAlias
	FunctionDefaultValueParser FunctionDefaultValueParse
	multiValuePacker           MultiValuePacker
	closure                    Closure
)

const (
	_ = iota
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
