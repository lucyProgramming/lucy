package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

const (
	ArrayTypeBoolean byte = 4
	ArrayTypeChar    byte = 5
	ArrayTypeFloat   byte = 6
	ArrayTypeDouble  byte = 7
	ArrayTypeByte    byte = 8
	ArrayTypeShort   byte = 9
	ArrayTypeInt     byte = 10
	ArrayTypeLong    byte = 11
)

type ArrayMeta struct {
	className                 string
	constructorFuncDescriptor string
	sliceDescriptor           string
	appendDescriptor          string
	appendAllDescriptor       string
	elementsFieldDescriptor   string
	setMethodDescription      string
	getMethodDescription      string
}

func init() {
	ArrayMetas[ast.VariableTypeBool] = &ArrayMeta{
		className:                 "lucy/deps/ArrayBool",
		constructorFuncDescriptor: "([Z)V",
		elementsFieldDescriptor:   "[Z",
		sliceDescriptor:           "(II)Llucy/deps/ArrayBool;",
		appendDescriptor:          "(Z)Llucy/deps/ArrayBool;",
		appendAllDescriptor:       "(Llucy/deps/ArrayBool;)Llucy/deps/ArrayBool;",
		setMethodDescription:      "(IZ)V",
		getMethodDescription:      "(I)Z",
	}
	ArrayMetas[ast.VariableTypeByte] = &ArrayMeta{
		className:                 "lucy/deps/ArrayByte",
		constructorFuncDescriptor: "([B)V",
		elementsFieldDescriptor:   "[B",
		sliceDescriptor:           "(II)Llucy/deps/ArrayByte;",
		appendDescriptor:          "(B)Llucy/deps/ArrayByte;",
		appendAllDescriptor:       "(Llucy/deps/ArrayByte;)Llucy/deps/ArrayByte;",
		setMethodDescription:      "(IB)V",
		getMethodDescription:      "(I)B",
	}
	ArrayMetas[ast.VariableTypeShort] = &ArrayMeta{
		className:                 "lucy/deps/ArrayShort",
		constructorFuncDescriptor: "([S)V",
		elementsFieldDescriptor:   "[S",
		sliceDescriptor:           "(II)Llucy/deps/ArrayShort;",
		appendDescriptor:          "(S)Llucy/deps/ArrayShort;",
		appendAllDescriptor:       "(Llucy/deps/ArrayShort;)Llucy/deps/ArrayShort;",
		setMethodDescription:      "(IS)V",
		getMethodDescription:      "(I)S",
	}
	ArrayMetas[ast.VariableTypeInt] = &ArrayMeta{
		className:                 "lucy/deps/ArrayInt",
		constructorFuncDescriptor: "([I)V",
		elementsFieldDescriptor:   "[I",
		sliceDescriptor:           "(II)Llucy/deps/ArrayInt;",
		appendDescriptor:          "(I)Llucy/deps/ArrayInt;",
		appendAllDescriptor:       "(Llucy/deps/ArrayInt;)Llucy/deps/ArrayInt;",
		setMethodDescription:      "(II)V",
		getMethodDescription:      "(I)I",
	}
	ArrayMetas[ast.VariableTypeLong] = &ArrayMeta{
		className:                 "lucy/deps/ArrayLong",
		constructorFuncDescriptor: "([J)V",
		elementsFieldDescriptor:   "[J",
		sliceDescriptor:           "(II)Llucy/deps/ArrayLong;",
		appendDescriptor:          "(J)Llucy/deps/ArrayLong;",
		appendAllDescriptor:       "(Llucy/deps/ArrayLong;)Llucy/deps/ArrayLong;",
		setMethodDescription:      "(IJ)V",
		getMethodDescription:      "(I)J",
	}
	ArrayMetas[ast.VariableTypeFloat] = &ArrayMeta{
		className:                 "lucy/deps/ArrayFloat",
		constructorFuncDescriptor: "([F)V",
		elementsFieldDescriptor:   "[F",
		sliceDescriptor:           "(II)Llucy/deps/ArrayFloat;",
		appendDescriptor:          "(F)Llucy/deps/ArrayFloat;",
		appendAllDescriptor:       "(Llucy/deps/ArrayFloat;)Llucy/deps/ArrayFloat;",
		setMethodDescription:      "(IF)V",
		getMethodDescription:      "(I)F",
	}
	ArrayMetas[ast.VariableTypeDouble] = &ArrayMeta{
		className:                 "lucy/deps/ArrayDouble",
		constructorFuncDescriptor: "([D)V",
		elementsFieldDescriptor:   "[D",
		sliceDescriptor:           "(II)Llucy/deps/ArrayDouble;",
		appendDescriptor:          "(D)Llucy/deps/ArrayDouble;",
		appendAllDescriptor:       "(Llucy/deps/ArrayDouble;)Llucy/deps/ArrayDouble;",
		setMethodDescription:      "(ID)V",
		getMethodDescription:      "(I)D",
	}
	ArrayMetas[ast.VariableTypeString] = &ArrayMeta{
		className:                 "lucy/deps/ArrayString",
		constructorFuncDescriptor: "([Ljava/lang/String;)V",
		elementsFieldDescriptor:   "[Ljava/lang/String;",
		sliceDescriptor:           "(II)Llucy/deps/ArrayString;",
		appendDescriptor:          "(Ljava/lang/String;)Llucy/deps/ArrayString;",
		appendAllDescriptor:       "(Llucy/deps/ArrayString;)Llucy/deps/ArrayString;",
		setMethodDescription:      "(ILlucy/deps/ArrayString;)V",
		getMethodDescription:      "(ILlucy/deps/ArrayString;)",
	}
	ArrayMetas[ast.VariableTypeObject] = &ArrayMeta{
		className:                 "lucy/deps/ArrayObject",
		constructorFuncDescriptor: "([Ljava/lang/Object;)V",
		elementsFieldDescriptor:   "[Ljava/lang/Object;",
		sliceDescriptor:           "(II)Llucy/deps/ArrayObject;",
		appendDescriptor:          "(Ljava/lang/Object;)Llucy/deps/ArrayObject;",
		appendAllDescriptor:       "(Llucy/deps/ArrayObject;)Llucy/deps/ArrayObject;",
		setMethodDescription:      "(ILjava/lang/Object;)V",
		getMethodDescription:      "(ILjava/lang/Object;)",
	}
	ArrayMetas[ast.VariableTypeArray] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeMap] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeJavaArray] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeFunction] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeEnum] = ArrayMetas[ast.VariableTypeInt]

}
