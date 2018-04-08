package cg

const (
	CONSTANT_POOL_MAX_SIZE                = 65536
	CLASS_MAGIC                    uint32 = 0xcafebabe
	ATTRIBUTE_NAME_SOURCE_FILE            = "SourceFile"
	ATTRIBUTE_NAME_CONST_VALUE            = "ConstantValue"
	ATTRIBUTE_NAME_LUCY_TYPE_ALIAS        = "LucyTypeAlias"
	ATTRIBUTE_NAME_SIGNATURE              = "Signature"
	// lucy attribute
	ATTRIBUTE_NAME_LUCY_FIELD_DESCRIPTOR       = "LucyFieldDescriptor"
	ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR      = "LucyMethodDescriptor"
	ATTRIBUTE_NAME_LUCY_CLOSURE_FUNCTION_CLASS = "LucyClosureFunctionClass"
	ATTRIBUTE_NAME_LUCY_INNER_STATIC_METHOD    = "LucyInnerStaticMethod"
	ATTRIBUTE_NAME_LUCY_TRIGGER_PACKAGE_INIT   = "LucyTriggerPackageInitMethod"
	ATTRIBUTE_NAME_LUCY_DEFAULT_PARAMETERS     = "LucyDefaultParameters"
)
