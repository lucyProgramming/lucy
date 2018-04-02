package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

const (
	ATYPE_T_BOOLEAN byte = 4
	ATYPE_T_CHAR    byte = 5
	ATYPE_T_FLOAT   byte = 6
	ATYPE_T_DOUBLE  byte = 7
	ATYPE_T_BYTE    byte = 8
	ATYPE_T_SHORT   byte = 9
	ATYPE_T_INT     byte = 10
	ATYPE_T_LONG    byte = 11
)

type ArrayMeta struct {
	classname, constructorFuncDescriptor  string
	getDescriptor, setDescriptor          string
	elementsFieldDescriptor               string
	sliceDescriptor                       string
	appendDescriptor, appendAllDescriptor string
	getJavaArrayDescriptor                string
}

func init() {
	ArrayMetas[ast.VARIABLE_TYPE_BOOL] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayBool",
		constructorFuncDescriptor: "([Z)V",
		getDescriptor:             "(I)Z",
		setDescriptor:             "(IZ)V",
		elementsFieldDescriptor:   "[Z",
		sliceDescriptor:           "(II)Llucy/deps/ArrayBool;",
		appendDescriptor:          "(Z)Llucy/deps/ArrayBool;",
		appendAllDescriptor:       "([Z)Llucy/deps/ArrayBool;",
		getJavaArrayDescriptor:    "()[Z",
	}
	ArrayMetas[ast.VARIABLE_TYPE_BYTE] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayByte",
		constructorFuncDescriptor: "([B)V",
		getDescriptor:             "(I)B",
		setDescriptor:             "(IB)V",
		elementsFieldDescriptor:   "[B",
		sliceDescriptor:           "(II)Llucy/deps/ArrayByte;",
		appendDescriptor:          "(B)Llucy/deps/ArrayByte;",
		appendAllDescriptor:       "([B)Llucy/deps/ArrayByte;",
		getJavaArrayDescriptor:    "()[B",
	}
	ArrayMetas[ast.VARIABLE_TYPE_SHORT] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayShort",
		constructorFuncDescriptor: "([S)V",
		getDescriptor:             "(I)S",
		setDescriptor:             "(IS)V",
		elementsFieldDescriptor:   "[S",
		sliceDescriptor:           "(II)Llucy/deps/ArrayShort;",
		appendDescriptor:          "(S)Llucy/deps/ArrayShort;",
		appendAllDescriptor:       "([S)Llucy/deps/ArrayShort;",
		getJavaArrayDescriptor:    "()[S",
	}
	ArrayMetas[ast.VARIABLE_TYPE_INT] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayInt",
		constructorFuncDescriptor: "([I)V",
		getDescriptor:             "(I)I",
		setDescriptor:             "(II)V",
		elementsFieldDescriptor:   "[I",
		sliceDescriptor:           "(II)Llucy/deps/ArrayInt;",
		appendDescriptor:          "(I)Llucy/deps/ArrayInt;",
		appendAllDescriptor:       "([I)Llucy/deps/ArrayInt;",
		getJavaArrayDescriptor:    "()[I",
	}
	ArrayMetas[ast.VARIABLE_TYPE_LONG] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayLong",
		constructorFuncDescriptor: "([J)V",
		getDescriptor:             "(I)J",
		setDescriptor:             "(IJ)V",
		elementsFieldDescriptor:   "[J",
		sliceDescriptor:           "(II)Llucy/deps/ArrayLong;",
		appendDescriptor:          "(J)Llucy/deps/ArrayLong;",
		appendAllDescriptor:       "([J)Llucy/deps/ArrayLong;",
		getJavaArrayDescriptor:    "()[J",
	}
	ArrayMetas[ast.VARIABLE_TYPE_FLOAT] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayFloat",
		constructorFuncDescriptor: "([F)V",
		getDescriptor:             "(I)F",
		setDescriptor:             "(IF)V",
		elementsFieldDescriptor:   "[F",
		sliceDescriptor:           "(II)Llucy/deps/ArrayFloat;",
		appendDescriptor:          "(F)Llucy/deps/ArrayFloat;",
		appendAllDescriptor:       "([F)Llucy/deps/ArrayFloat;",
		getJavaArrayDescriptor:    "()[F",
	}
	ArrayMetas[ast.VARIABLE_TYPE_DOUBLE] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayDouble",
		constructorFuncDescriptor: "([D)V",
		getDescriptor:             "(I)D",
		setDescriptor:             "(ID)V",
		elementsFieldDescriptor:   "[D",
		sliceDescriptor:           "(II)Llucy/deps/ArrayDouble;",
		appendDescriptor:          "(D)Llucy/deps/ArrayDouble;",
		appendAllDescriptor:       "([D)Llucy/deps/ArrayDouble;",
		getJavaArrayDescriptor:    "()[D",
	}
	ArrayMetas[ast.VARIABLE_TYPE_STRING] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayString",
		constructorFuncDescriptor: "([Ljava/lang/String;)V",
		getDescriptor:             "(I)Ljava/lang/String;",
		setDescriptor:             "(ILjava/lang/String;)V",
		elementsFieldDescriptor:   "[Ljava/lang/String;",
		sliceDescriptor:           "(II)Llucy/deps/ArrayString;",
		appendDescriptor:          "(Ljava/lang/String;)Llucy/deps/ArrayString;",
		appendAllDescriptor:       "([Ljava/lang/String;)Llucy/deps/ArrayString;",
		getJavaArrayDescriptor:    "()[Llucy/deps/ArrayString;",
	}
	ArrayMetas[ast.VARIABLE_TYPE_OBJECT] = &ArrayMeta{
		classname:                 "lucy/deps/ArrayObject",
		constructorFuncDescriptor: "([Ljava/lang/Object;)V",
		getDescriptor:             "(I)Ljava/lang/Object;",
		setDescriptor:             "(ILjava/lang/Object;)V",
		elementsFieldDescriptor:   "[Ljava/lang/Object;",
		sliceDescriptor:           "(II)Llucy/deps/ArrayObject;",
		appendDescriptor:          "(Ljava/lang/Object;)Llucy/deps/ArrayObject;",
		appendAllDescriptor:       "([Ljava/lang/Object;)Llucy/deps/ArrayObject;",
		getJavaArrayDescriptor:    "()[Llucy/deps/ArrayObject;",
	}
	ArrayMetas[ast.VARIABLE_TYPE_ARRAY] = ArrayMetas[ast.VARIABLE_TYPE_OBJECT]
	ArrayMetas[ast.VARIABLE_TYPE_MAP] = ArrayMetas[ast.VARIABLE_TYPE_OBJECT]
	for _, v := range ArrayMetas {
		ArrayMetasMap[v.classname] = v
	}
}
