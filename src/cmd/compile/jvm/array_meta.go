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
		appendDescriptor:          "(Z)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayBool;)V",
		setMethodDescription:      "(IZ)V",
		getMethodDescription:      "(I)Z",
	}
	ArrayMetas[ast.VariableTypeByte] = &ArrayMeta{
		className:                 "lucy/deps/ArrayByte",
		constructorFuncDescriptor: "([B)V",
		elementsFieldDescriptor:   "[B",
		sliceDescriptor:           "(II)Llucy/deps/ArrayByte;",
		appendDescriptor:          "(B)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayByte;)V",
		setMethodDescription:      "(IB)V",
		getMethodDescription:      "(I)B",
	}
	ArrayMetas[ast.VariableTypeShort] = &ArrayMeta{
		className:                 "lucy/deps/ArrayShort",
		constructorFuncDescriptor: "([S)V",
		elementsFieldDescriptor:   "[S",
		sliceDescriptor:           "(II)Llucy/deps/ArrayShort;",
		appendDescriptor:          "(S)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayShort;)V",
		setMethodDescription:      "(IS)V",
		getMethodDescription:      "(I)S",
	}
	ArrayMetas[ast.VariableTypeChar] = &ArrayMeta{
		className:                 "lucy/deps/CharInt",
		constructorFuncDescriptor: "([I)V",
		elementsFieldDescriptor:   "[I",
		sliceDescriptor:           "(II)Llucy/deps/CharInt;",
		appendDescriptor:          "(I)V",
		appendAllDescriptor:       "(Llucy/deps/CharInt;)V",
		setMethodDescription:      "(IC)V",
		getMethodDescription:      "(I)C",
	}
	ArrayMetas[ast.VariableTypeInt] = &ArrayMeta{
		className:                 "lucy/deps/ArrayInt",
		constructorFuncDescriptor: "([I)V",
		elementsFieldDescriptor:   "[I",
		sliceDescriptor:           "(II)Llucy/deps/ArrayInt;",
		appendDescriptor:          "(I)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayInt;)V",
		setMethodDescription:      "(II)V",
		getMethodDescription:      "(I)I",
	}
	ArrayMetas[ast.VariableTypeLong] = &ArrayMeta{
		className:                 "lucy/deps/ArrayLong",
		constructorFuncDescriptor: "([J)V",
		elementsFieldDescriptor:   "[J",
		sliceDescriptor:           "(II)Llucy/deps/ArrayLong;",
		appendDescriptor:          "(J)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayLong;)V",
		setMethodDescription:      "(IJ)V",
		getMethodDescription:      "(I)J",
	}
	ArrayMetas[ast.VariableTypeFloat] = &ArrayMeta{
		className:                 "lucy/deps/ArrayFloat",
		constructorFuncDescriptor: "([F)V",
		elementsFieldDescriptor:   "[F",
		sliceDescriptor:           "(II)Llucy/deps/ArrayFloat;",
		appendDescriptor:          "(F)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayFloat;)V",
		setMethodDescription:      "(IF)V",
		getMethodDescription:      "(I)F",
	}
	ArrayMetas[ast.VariableTypeDouble] = &ArrayMeta{
		className:                 "lucy/deps/ArrayDouble",
		constructorFuncDescriptor: "([D)V",
		elementsFieldDescriptor:   "[D",
		sliceDescriptor:           "(II)Llucy/deps/ArrayDouble;",
		appendDescriptor:          "(D)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayDouble;)V",
		setMethodDescription:      "(ID)V",
		getMethodDescription:      "(I)D",
	}
	ArrayMetas[ast.VariableTypeString] = &ArrayMeta{
		className:                 "lucy/deps/ArrayString",
		constructorFuncDescriptor: "([Ljava/lang/String;)V",
		elementsFieldDescriptor:   "[Ljava/lang/String;",
		sliceDescriptor:           "(II)Llucy/deps/ArrayString;",
		appendDescriptor:          "(Ljava/lang/String;)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayString;)V",
		setMethodDescription:      "(ILjava/lang/String;)V",
		getMethodDescription:      "(I)Ljava/lang/String;",
	}
	ArrayMetas[ast.VariableTypeObject] = &ArrayMeta{
		className:                 "lucy/deps/ArrayObject",
		constructorFuncDescriptor: "([Ljava/lang/Object;)V",
		elementsFieldDescriptor:   "[Ljava/lang/Object;",
		sliceDescriptor:           "(II)Llucy/deps/ArrayObject;",
		appendDescriptor:          "(Ljava/lang/Object;)V",
		appendAllDescriptor:       "(Llucy/deps/ArrayObject;)V",
		setMethodDescription:      "(ILjava/lang/Object;)V",
		getMethodDescription:      "(I)Ljava/lang/Object;",
	}
	ArrayMetas[ast.VariableTypeArray] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeMap] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeJavaArray] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeFunction] = ArrayMetas[ast.VariableTypeObject]
	ArrayMetas[ast.VariableTypeEnum] = ArrayMetas[ast.VariableTypeInt]

}
