package cg

const (
	CONSTANT_POOL_MAX_SIZE            = 65536
	CLASS_MAGIC                uint32 = 0xcafebabe
	ATTRIBUTE_NAME_SOURCE_FILE        = "SourceFile"
	ATTRIBUTE_NAME_CONST_VALUE        = "ConstantValue"
	ATTRIBUTE_NAME_STACK_MAP          = "StackMapTable"
	ATTRIBUTE_NAME_SIGNATURE          = "Signature"
	// lucy attribute
	ATTRIBUTE_NAME_LUCY_FIELD_DESCRIPTOR     = "LucyFieldDescriptor"
	ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR    = "LucyMethodDescriptor"
	ATTRIBUTE_NAME_LUCY_COMPILTER_AUTO       = "LucyCompilerAuto"
	ATTRIBUTE_NAME_LUCY_INNER_STATIC_METHOD  = "LucyInnerStaticMethod"
	ATTRIBUTE_NAME_LUCY_TRIGGER_PACKAGE_INIT = "LucyTriggerPackageInitMethod"
	ATTRIBUTE_NAME_LUCY_DEFAULT_PARAMETERS   = "LucyDefaultParameters"
	ATTRIBUTE_NAME_LUCY_TYPE_ALIAS           = "LucyTypeAlias"
	ATTRIBUTE_NAME_LUCY_ENUM                 = "LucyEnum"
	ATTRIBUTE_NAME_LUCY_CONST                = "LucyConst" // indicate a package const
	ATTRIBUTE_NAME_METHOD_PARAMETERS         = "MethodParameters"
	ATTRIBUTE_NAME_LUCY_RETURNLIST_NAMES     = "LucyReturnListName"
	ATTRIBUTE_NAME_LUCY_TEMPLATE_FUNCTION    = "LucyTemplateFunction"
)
