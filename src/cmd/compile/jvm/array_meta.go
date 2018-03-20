package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
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
		classname:                 "lucy/lang/Arrayboolean",
		constructorFuncDescriptor: "([ZI)V",
		getDescriptor:             "(I)Z",
		setDescriptor:             "(IZ)V",
		elementsFieldDescriptor:   "[Z",
		sliceDescriptor:           "(II)Llucy/lang/Arrayboolean;",
		appendDescriptor:          "(Z)Llucy/lang/Arrayboolean;",
		appendAllDescriptor:       "([Z)Llucy/lang/Arrayboolean;",
		getJavaArrayDescriptor:    "()[Z",
	}
	ArrayMetas[ast.VARIABLE_TYPE_BYTE] = &ArrayMeta{
		classname:                 "lucy/lang/Arraybyte",
		constructorFuncDescriptor: "([BI)V",
		getDescriptor:             "(I)B",
		setDescriptor:             "(IB)V",
		elementsFieldDescriptor:   "[B",
		sliceDescriptor:           "(II)Llucy/lang/Arraybyte;",
		appendDescriptor:          "(B)Llucy/lang/Arraybyte;",
		appendAllDescriptor:       "([B)Llucy/lang/Arraybyte;",
		getJavaArrayDescriptor:    "()[B",
	}
	ArrayMetas[ast.VARIABLE_TYPE_SHORT] = &ArrayMeta{
		classname:                 "lucy/lang/Arrayshort",
		constructorFuncDescriptor: "([SI)V",
		getDescriptor:             "(I)S",
		setDescriptor:             "(IS)V",
		elementsFieldDescriptor:   "[S",
		sliceDescriptor:           "(II)Llucy/lang/Arrayshort;",
		appendDescriptor:          "(S)Llucy/lang/Arrayshort;",
		appendAllDescriptor:       "([S)Llucy/lang/Arrayshort;",
		getJavaArrayDescriptor:    "()[S",
	}
	ArrayMetas[ast.VARIABLE_TYPE_INT] = &ArrayMeta{
		classname:                 "lucy/lang/Arrayint",
		constructorFuncDescriptor: "([II)V",
		getDescriptor:             "(I)I",
		setDescriptor:             "(II)V",
		elementsFieldDescriptor:   "[I",
		sliceDescriptor:           "(II)Llucy/lang/Arrayint;",
		appendDescriptor:          "(I)Llucy/lang/Arrayint;",
		appendAllDescriptor:       "([I)Llucy/lang/Arrayint;",
		getJavaArrayDescriptor:    "()[I",
	}
	ArrayMetas[ast.VARIABLE_TYPE_LONG] = &ArrayMeta{
		classname:                 "lucy/lang/Arraylong",
		constructorFuncDescriptor: "([JI)V",
		getDescriptor:             "(I)J",
		setDescriptor:             "(IJ)V",
		elementsFieldDescriptor:   "[J",
		sliceDescriptor:           "(II)Llucy/lang/Arraylong;",
		appendDescriptor:          "(J)Llucy/lang/Arraylong;",
		appendAllDescriptor:       "([J)Llucy/lang/Arraylong;",
		getJavaArrayDescriptor:    "()[J",
	}
	ArrayMetas[ast.VARIABLE_TYPE_FLOAT] = &ArrayMeta{
		classname:                 "lucy/lang/Arrayfloat",
		constructorFuncDescriptor: "([FI)V",
		getDescriptor:             "(I)F",
		setDescriptor:             "(IF)V",
		elementsFieldDescriptor:   "[F",
		sliceDescriptor:           "(II)Llucy/lang/Arrayfloat;",
		appendDescriptor:          "(F)Llucy/lang/Arrayfloat;",
		appendAllDescriptor:       "([F)Llucy/lang/Arrayfloat;",
		getJavaArrayDescriptor:    "()[F",
	}
	ArrayMetas[ast.VARIABLE_TYPE_DOUBLE] = &ArrayMeta{
		classname:                 "lucy/lang/Arraydouble",
		constructorFuncDescriptor: "([DI)V",
		getDescriptor:             "(I)D",
		setDescriptor:             "(ID)V",
		elementsFieldDescriptor:   "[D",
		sliceDescriptor:           "(II)Llucy/lang/Arraydouble;",
		appendDescriptor:          "(D)Llucy/lang/Arraydouble;",
		appendAllDescriptor:       "([D)Llucy/lang/Arraydouble;",
		getJavaArrayDescriptor:    "()[D",
	}
	ArrayMetas[ast.VARIABLE_TYPE_STRING] = &ArrayMeta{
		classname:                 "lucy/lang/ArrayString",
		constructorFuncDescriptor: "([Ljava/lang/String;I)V",
		getDescriptor:             "(I)Ljava/lang/String;",
		setDescriptor:             "(ILjava/lang/String;)V",
		elementsFieldDescriptor:   "[Ljava/lang/String;",
		sliceDescriptor:           "(II)Llucy/lang/ArrayString;",
		appendDescriptor:          "(Ljava/lang/String;)Llucy/lang/ArrayString;",
		appendAllDescriptor:       "([Ljava/lang/String;)Llucy/lang/ArrayString;",
		getJavaArrayDescriptor:    "()[Llucy/lang/ArrayString;",
	}
	ArrayMetas[ast.VARIABLE_TYPE_OBJECT] = &ArrayMeta{
		classname:                 "lucy/lang/ArrayObject",
		constructorFuncDescriptor: "([Ljava/lang/Object;I)V",
		getDescriptor:             "(I)Ljava/lang/Object;",
		setDescriptor:             "(ILjava/lang/Object;)V",
		elementsFieldDescriptor:   "[Ljava/lang/Object;",
		sliceDescriptor:           "(II)Llucy/lang/ArrayObject;",
		appendDescriptor:          "(Ljava/lang/Object;)Llucy/lang/ArrayObject;",
		appendAllDescriptor:       "([Ljava/lang/Object;)Llucy/lang/ArrayObject;",
		getJavaArrayDescriptor:    "()[Llucy/lang/ArrayObject;",
	}
	ArrayMetas[ast.VARIABLE_TYPE_ARRAY] = ArrayMetas[ast.VARIABLE_TYPE_OBJECT]
	ArrayMetas[ast.VARIABLE_TYPE_MAP] = ArrayMetas[ast.VARIABLE_TYPE_OBJECT]
	for _, v := range ArrayMetas {
		ArrayMetasMap[v.classname] = v
	}
}
