package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

const (
	_ = iota
	ClosureClassInt
	ClosureClassLong
	ClosureClassFloat
	ClosureClassDouble
	ClosureClassString
	ClosureClassObject
)

type Closure struct {
	ClosureObjectMetas map[int]*ClosureObjectMeta
}

var (
	closure Closure
)

type ClosureObjectMeta struct {
	className       string
	fieldName       string
	fieldDescriptor string
}

func init() {
	closure.ClosureObjectMetas = make(map[int]*ClosureObjectMeta)
	closure.ClosureObjectMetas[ClosureClassInt] = &ClosureObjectMeta{
		className:       "lucy/deps/ClosureInt",
		fieldName:       "value",
		fieldDescriptor: "I",
	}
	closure.ClosureObjectMetas[ClosureClassLong] = &ClosureObjectMeta{
		className:       "lucy/deps/ClosureLong",
		fieldName:       "value",
		fieldDescriptor: "J",
	}
	closure.ClosureObjectMetas[ClosureClassFloat] = &ClosureObjectMeta{
		className:       "lucy/deps/ClosureFloat",
		fieldName:       "value",
		fieldDescriptor: "F",
	}
	closure.ClosureObjectMetas[ClosureClassDouble] = &ClosureObjectMeta{
		className:       "lucy/deps/ClosureDouble",
		fieldName:       "value",
		fieldDescriptor: "D",
	}
	closure.ClosureObjectMetas[ClosureClassString] = &ClosureObjectMeta{
		className:       "lucy/deps/ClosureString",
		fieldName:       "value",
		fieldDescriptor: "Ljava/lang/String;",
	}
	closure.ClosureObjectMetas[ClosureClassObject] = &ClosureObjectMeta{
		className:       "lucy/deps/ClosureObject",
		fieldName:       "value",
		fieldDescriptor: "Ljava/lang/Object;",
	}
}

func (closure *Closure) getMeta(t int) (meta *ClosureObjectMeta) {
	switch t {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		meta = closure.ClosureObjectMetas[ClosureClassInt]
	case ast.VariableTypeLong:
		meta = closure.ClosureObjectMetas[ClosureClassLong]
	case ast.VariableTypeFloat:
		meta = closure.ClosureObjectMetas[ClosureClassFloat]
	case ast.VariableTypeDouble:
		meta = closure.ClosureObjectMetas[ClosureClassDouble]
	case ast.VariableTypeString:
		meta = closure.ClosureObjectMetas[ClosureClassString]
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeArray: //[]int
		fallthrough
	case ast.VariableTypeJavaArray: // java array int[]
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeMap:
		meta = closure.ClosureObjectMetas[ClosureClassObject]
	}
	return
}

/*
	create a closure var, inited and leave on stack
*/
func (closure *Closure) createClosureVar(class *cg.ClassHighLevel,
	code *cg.AttributeCode, v *ast.Type) (maxStack uint16) {
	maxStack = 2
	var meta *ClosureObjectMeta
	switch v.Type {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		meta = closure.ClosureObjectMetas[ClosureClassInt]
	case ast.VariableTypeLong:
		meta = closure.ClosureObjectMetas[ClosureClassLong]
	case ast.VariableTypeFloat:
		meta = closure.ClosureObjectMetas[ClosureClassFloat]
	case ast.VariableTypeDouble:
		meta = closure.ClosureObjectMetas[ClosureClassDouble]
	case ast.VariableTypeString:
		meta = closure.ClosureObjectMetas[ClosureClassString]
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeArray: //[]int
		fallthrough
	case ast.VariableTypeJavaArray: // java array int[]
		fallthrough
	case ast.VariableTypeMap:
		meta = closure.ClosureObjectMetas[ClosureClassObject]
	}
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      meta.className,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}

func (closure *Closure) storeLocalClosureVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.Variable) {
	var meta *ClosureObjectMeta
	switch v.Type.Type {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		meta = closure.ClosureObjectMetas[ClosureClassInt]
	case ast.VariableTypeLong:
		meta = closure.ClosureObjectMetas[ClosureClassLong]
	case ast.VariableTypeFloat:
		meta = closure.ClosureObjectMetas[ClosureClassFloat]
	case ast.VariableTypeDouble:
		meta = closure.ClosureObjectMetas[ClosureClassDouble]
	case ast.VariableTypeString:
		meta = closure.ClosureObjectMetas[ClosureClassString]
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeMap:
		fallthrough
	case ast.VariableTypeArray:
		fallthrough
	case ast.VariableTypeJavaArray:
		meta = closure.ClosureObjectMetas[ClosureClassObject]
	}
	code.Codes[code.CodeLength] = cg.OP_putfield
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      meta.className,
		Field:      meta.fieldName,
		Descriptor: meta.fieldDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
}

/*
	create a closure var on stack
*/
func (closure *Closure) loadLocalClosureVar(class *cg.ClassHighLevel, code *cg.AttributeCode,
	v *ast.Variable) (maxStack uint16) {
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, v.LocalValOffset)...)
	closure.unPack(class, code, v.Type)
	maxStack = jvmSlotSize(v.Type)
	return
}

/*
	closure object is on stack
*/
func (closure *Closure) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.Type) {
	switch v.Type {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeInt:
		meta := closure.ClosureObjectMetas[ClosureClassInt]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      meta.fieldName,
			Descriptor: meta.fieldDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3

	case ast.VariableTypeLong:
		meta := closure.ClosureObjectMetas[ClosureClassLong]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      meta.fieldName,
			Descriptor: meta.fieldDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3

	case ast.VariableTypeFloat:
		meta := closure.ClosureObjectMetas[ClosureClassFloat]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      meta.fieldName,
			Descriptor: meta.fieldDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3

	case ast.VariableTypeDouble:
		meta := closure.ClosureObjectMetas[ClosureClassDouble]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      meta.fieldName,
			Descriptor: meta.fieldDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3

	case ast.VariableTypeString:
		meta := closure.ClosureObjectMetas[ClosureClassString]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      meta.fieldName,
			Descriptor: meta.fieldDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeMap:
		fallthrough
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeArray:
		fallthrough
	case ast.VariableTypeJavaArray:
		meta := closure.ClosureObjectMetas[ClosureClassObject]
		code.Codes[code.CodeLength] = cg.OP_getfield
		class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
			Class:      meta.className,
			Field:      meta.fieldName,
			Descriptor: meta.fieldDescriptor,
		}, code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.CodeLength += 3
		typeConverter.castPointerTypeToRealType(class, code, v)

	}

}
