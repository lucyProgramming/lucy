package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

var (
	ArrayMetas                = map[int]*ArrayMeta{}
	ArrayMetasMap             = make(map[string]*ArrayMeta)
	primitiveObjectConverter  PrimitiveObjectConverter
	Descriptor                Descript
	LucyMethodSignatureParser LucyMethodSignatureParse
	LucyFieldSignatureParser  LucyFieldSignatureParse
)

const (
	java_arrylist_class                     = "java/util/ArrayList"
	special_method_init                     = "<init>"
	java_string_class                       = "java/lang/String"
	java_hashmap_class                      = "java/util/HashMap"
	java_integer_class                      = "java/lang/Integer"
	java_float_class                        = "java/lang/Float"
	java_double_class                       = "java/lang/Double"
	java_long_class                         = "java/lang/Long"
	java_index_out_of_range_exception_class = "java/lang/ArrayIndexOutOfBoundsException"
)

func init() {
	ast.JvmSlotSizeHandler = func(v *ast.VariableType) uint16 {
		if v.Typ == ast.VARIABLE_TYPE_DOUBLE || ast.VARIABLE_TYPE_LONG == v.Typ {
			return 2
		} else {
			return 1
		}
	}
}
