package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClosureKind int

const (
	_ ClosureKind = iota
	ClosureClassInt
	ClosureClassLong
	ClosureClassFloat
	ClosureClassDouble
	ClosureClassString
	ClosureClassObject
)

type Closure struct {
	ClosureObjectMetas map[ClosureKind]*ClosureObjectMeta
}

type ClosureObjectMeta struct {
	className        string
	fieldName        string
	fieldDescription string
}

func init() {
	closure.ClosureObjectMetas = make(map[ClosureKind]*ClosureObjectMeta)
	closure.ClosureObjectMetas[ClosureClassInt] = &ClosureObjectMeta{
		className:        "lucy/deps/ClosureInt",
		fieldName:        "value",
		fieldDescription: "I",
	}
	closure.ClosureObjectMetas[ClosureClassLong] = &ClosureObjectMeta{
		className:        "lucy/deps/ClosureLong",
		fieldName:        "value",
		fieldDescription: "J",
	}
	closure.ClosureObjectMetas[ClosureClassFloat] = &ClosureObjectMeta{
		className:        "lucy/deps/ClosureFloat",
		fieldName:        "value",
		fieldDescription: "F",
	}
	closure.ClosureObjectMetas[ClosureClassDouble] = &ClosureObjectMeta{
		className:        "lucy/deps/ClosureDouble",
		fieldName:        "value",
		fieldDescription: "D",
	}
	closure.ClosureObjectMetas[ClosureClassString] = &ClosureObjectMeta{
		className:        "lucy/deps/ClosureString",
		fieldName:        "value",
		fieldDescription: "Ljava/lang/String;",
	}
	closure.ClosureObjectMetas[ClosureClassObject] = &ClosureObjectMeta{
		className:        "lucy/deps/ClosureObject",
		fieldName:        "value",
		fieldDescription: "Ljava/lang/Object;",
	}
}

func (this *Closure) getMeta(t ast.VariableTypeKind) (meta *ClosureObjectMeta) {
	switch t {
	case ast.VariableTypeBool:
		fallthrough
	case ast.VariableTypeByte:
		fallthrough
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeEnum:
		fallthrough
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		meta = this.ClosureObjectMetas[ClosureClassInt]
	case ast.VariableTypeLong:
		meta = this.ClosureObjectMetas[ClosureClassLong]
	case ast.VariableTypeFloat:
		meta = this.ClosureObjectMetas[ClosureClassFloat]
	case ast.VariableTypeDouble:
		meta = this.ClosureObjectMetas[ClosureClassDouble]
	case ast.VariableTypeString:
		meta = this.ClosureObjectMetas[ClosureClassString]
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeArray: //[]int
		fallthrough
	case ast.VariableTypeJavaArray: // java array int[]
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeMap:
		meta = this.ClosureObjectMetas[ClosureClassObject]
	}
	return
}

/*
	create a closure var, init and leave on stack
*/
func (this *Closure) createClosureVar(class *cg.ClassHighLevel,
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
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		meta = this.ClosureObjectMetas[ClosureClassInt]
	case ast.VariableTypeLong:
		meta = this.ClosureObjectMetas[ClosureClassLong]
	case ast.VariableTypeFloat:
		meta = this.ClosureObjectMetas[ClosureClassFloat]
	case ast.VariableTypeDouble:
		meta = this.ClosureObjectMetas[ClosureClassDouble]
	case ast.VariableTypeString:
		meta = this.ClosureObjectMetas[ClosureClassString]
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeArray: //[]int
		fallthrough
	case ast.VariableTypeJavaArray: // java array int[]
		fallthrough
	case ast.VariableTypeMap:
		meta = this.ClosureObjectMetas[ClosureClassObject]
	}
	code.Codes[code.CodeLength] = cg.OP_new
	class.InsertClassConst(meta.className, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.Codes[code.CodeLength+3] = cg.OP_dup
	code.CodeLength += 4
	code.Codes[code.CodeLength] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      meta.className,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	return
}

func (this *Closure) storeLocalClosureVar(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.Variable) {
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
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		meta = this.ClosureObjectMetas[ClosureClassInt]
	case ast.VariableTypeLong:
		meta = this.ClosureObjectMetas[ClosureClassLong]
	case ast.VariableTypeFloat:
		meta = this.ClosureObjectMetas[ClosureClassFloat]
	case ast.VariableTypeDouble:
		meta = this.ClosureObjectMetas[ClosureClassDouble]
	case ast.VariableTypeString:
		meta = this.ClosureObjectMetas[ClosureClassString]
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeMap:
		fallthrough
	case ast.VariableTypeArray:
		fallthrough
	case ast.VariableTypeJavaArray:
		meta = this.ClosureObjectMetas[ClosureClassObject]
	}
	code.Codes[code.CodeLength] = cg.OP_putfield
	class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
		Class:      meta.className,
		Field:      meta.fieldName,
		Descriptor: meta.fieldDescription,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
}

/*
	create a closure var on stack
*/
func (this *Closure) loadLocalClosureVar(class *cg.ClassHighLevel, code *cg.AttributeCode,
	v *ast.Variable) (maxStack uint16) {
	copyOPs(code, loadLocalVariableOps(ast.VariableTypeObject, v.LocalValOffset)...)
	this.unPack(class, code, v.Type)
	maxStack = jvmSlotSize(v.Type)
	return
}

/*
	closure object is on stack
*/
func (this *Closure) unPack(class *cg.ClassHighLevel, code *cg.AttributeCode, v *ast.Type) {
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
	case ast.VariableTypeChar:
		fallthrough
	case ast.VariableTypeInt:
		meta = this.ClosureObjectMetas[ClosureClassInt]
	case ast.VariableTypeLong:
		meta = this.ClosureObjectMetas[ClosureClassLong]
	case ast.VariableTypeFloat:
		meta = this.ClosureObjectMetas[ClosureClassFloat]
	case ast.VariableTypeDouble:
		meta = this.ClosureObjectMetas[ClosureClassDouble]
	case ast.VariableTypeString:
		meta = this.ClosureObjectMetas[ClosureClassString]
	case ast.VariableTypeFunction:
		fallthrough
	case ast.VariableTypeMap:
		fallthrough
	case ast.VariableTypeObject:
		fallthrough
	case ast.VariableTypeArray:
		fallthrough
	case ast.VariableTypeJavaArray:
		meta = this.ClosureObjectMetas[ClosureClassObject]
	}
	code.Codes[code.CodeLength] = cg.OP_getfield
	class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
		Class:      meta.className,
		Field:      meta.fieldName,
		Descriptor: meta.fieldDescription,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
	if v.IsPointer() && v.Type != ast.VariableTypeString {
		typeConverter.castPointer(class, code, v)
	}
}
