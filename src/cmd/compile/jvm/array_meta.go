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
	classname, initFuncDescriptor string
	getDescriptor, setDescriptor  string
	elementsFieldDescriptor       string
}

var (
	ArrayMetas = map[int]*ArrayMeta{}
)

func init() {
	ArrayMetas[ast.VARIABLE_TYPE_BOOL] = &ArrayMeta{
		classname:               "lucy/lang/Arrayboolean",
		initFuncDescriptor:      "([ZI)V",
		getDescriptor:           "(I)Z",
		setDescriptor:           "(IZ)V",
		elementsFieldDescriptor: "[Z",
	}
	ArrayMetas[ast.VARIABLE_TYPE_BYTE] = &ArrayMeta{
		classname:               "lucy/lang/Arraybyte",
		initFuncDescriptor:      "([BI)V",
		getDescriptor:           "(I)B",
		setDescriptor:           "(IB)V",
		elementsFieldDescriptor: "[B",
	}
	ArrayMetas[ast.VARIABLE_TYPE_SHORT] = &ArrayMeta{
		classname:               "lucy/lang/Arrayshort",
		initFuncDescriptor:      "([SI)V",
		getDescriptor:           "(I)S",
		setDescriptor:           "(IS)V",
		elementsFieldDescriptor: "[S",
	}
	ArrayMetas[ast.VARIABLE_TYPE_INT] = &ArrayMeta{
		classname:               "lucy/lang/Arrayint",
		initFuncDescriptor:      "([II)V",
		getDescriptor:           "(I)I",
		setDescriptor:           "(II)V",
		elementsFieldDescriptor: "[I",
	}
	ArrayMetas[ast.VARIABLE_TYPE_LONG] = &ArrayMeta{
		classname:               "lucy/lang/Arraylong",
		initFuncDescriptor:      "([JI)V",
		getDescriptor:           "(I)J",
		setDescriptor:           "(IJ)V",
		elementsFieldDescriptor: "[J",
	}
	ArrayMetas[ast.VARIABLE_TYPE_FLOAT] = &ArrayMeta{
		classname:               "lucy/lang/Arrayfloat",
		initFuncDescriptor:      "([FI)V",
		getDescriptor:           "(I)F",
		setDescriptor:           "(IF)V",
		elementsFieldDescriptor: "[F",
	}
	ArrayMetas[ast.VARIABLE_TYPE_DOUBLE] = &ArrayMeta{
		classname:               "lucy/lang/Arraydouble",
		initFuncDescriptor:      "([DI)V",
		getDescriptor:           "(I)D",
		setDescriptor:           "(ID)V",
		elementsFieldDescriptor: "[D",
	}
	ArrayMetas[ast.VARIABLE_TYPE_STRING] = &ArrayMeta{
		classname:               "lucy/lang/ArrayString",
		initFuncDescriptor:      "([Ljava/lang/String;I)V",
		getDescriptor:           "(I)Ljava/lang/String;",
		setDescriptor:           "(ILjava/lang/String;)V",
		elementsFieldDescriptor: "[Ljava/lang/String;",
	}
	ArrayMetas[ast.VARIABLE_TYPE_OBJECT] = &ArrayMeta{
		classname:               "lucy/lang/ArrayObject",
		initFuncDescriptor:      "([Ljava/lang/Object;I)V",
		getDescriptor:           "(I)Ljava/lang/Object;",
		setDescriptor:           "(ILjava/lang/Object;)V",
		elementsFieldDescriptor: "[Ljava/lang/Object;",
	}
	ArrayMetas[ast.VARIABLE_TYPE_ARRAY_INSTANCE] = &ArrayMeta{
		classname:               "lucy/lang/ArrayObject",
		initFuncDescriptor:      "([Ljava/lang/Object;I)V",
		getDescriptor:           "(I)Ljava/lang/Object;",
		setDescriptor:           "(ILjava/lang/Object;)V",
		elementsFieldDescriptor: "[Ljava/lang/Object;",
	}
}
