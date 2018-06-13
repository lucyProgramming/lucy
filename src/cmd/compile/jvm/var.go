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
)

const (
	special_method_init                     = "<init>"
	java_root_object_array                  = "[Ljava/lang/Object;"
	java_string_class                       = "java/lang/String"
	java_root_class                         = "java/lang/Object"
	java_hashmap_class                      = "java/util/HashMap"
	java_integer_class                      = "java/lang/Integer"
	java_float_class                        = "java/lang/Float"
	java_double_class                       = "java/lang/Double"
	java_long_class                         = "java/lang/Long"
	java_index_out_of_range_exception_class = "java/lang/ArrayIndexOutOfBoundsException"
	java_string_builder_class               = "java/lang/StringBuilder"
	java_throwable_class                    = "java/lang/Throwable"
	java_print_stream_class                 = "java/io/PrintStream"
)
