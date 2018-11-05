package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StackMapState struct {
	Locals []*cg.StackMapVerificationTypeInfo
	Stacks []*cg.StackMapVerificationTypeInfo
}

func (this *StackMapState) appendLocals(class *cg.ClassHighLevel, v *ast.Type) {
	this.Locals = append(this.Locals,
		this.newStackMapVerificationTypeInfo(class, v))
}

func (this *StackMapState) addTop(absent *StackMapState) {
	if this == absent {
		return
	}
	length := len(absent.Locals) - len(this.Locals)
	if length == 0 {
		return
	}
	oldLength := len(this.Locals)
	verify := &cg.StackMapVerificationTypeInfo{}
	verify.Verify = &cg.StackMapTopVariableInfo{}
	for i := 0; i < length; i++ {
		tt := absent.Locals[i+oldLength].Verify
		_, isDouble := tt.(*cg.StackMapDoubleVariableInfo)
		_, isLong := tt.(*cg.StackMapLongVariableInfo)
		if isDouble || isLong {
			this.Locals = append(this.Locals, verify, verify)
		} else {
			this.Locals = append(this.Locals, verify)
		}
	}
}

func (this *StackMapState) newObjectVariableType(name string) *ast.Type {
	ret := &ast.Type{}
	ret.Type = ast.VariableTypeObject
	ret.Class = &ast.Class{}
	ret.Class.Name = name
	return ret
}

func (this *StackMapState) popStack(pop int) {
	if pop == 0 {
		return
	}
	if pop < 0 {
		panic("negative pop")
	}
	if len(this.Stacks) == 0 {
		panic("already 0")
	}
	this.Stacks = this.Stacks[:len(this.Stacks)-pop]
}
func (this *StackMapState) pushStack(class *cg.ClassHighLevel, v *ast.Type) {
	if this == nil {
		panic("s is nil")
	}
	this.Stacks = append(this.Stacks, this.newStackMapVerificationTypeInfo(class, v))
}
func (this *StackMapState) initFromLast(last *StackMapState) *StackMapState {
	this.Locals = make([]*cg.StackMapVerificationTypeInfo, len(last.Locals))
	copy(this.Locals, last.Locals)
	return this
}

func (this *StackMapState) newStackMapVerificationTypeInfo(class *cg.ClassHighLevel,
	t *ast.Type) (ret *cg.StackMapVerificationTypeInfo) {
	ret = &cg.StackMapVerificationTypeInfo{}
	switch t.Type {
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
		ret.Verify = &cg.StackMapIntegerVariableInfo{}
	case ast.VariableTypeLong:
		ret.Verify = &cg.StackMapLongVariableInfo{}
	case ast.VariableTypeFloat:
		ret.Verify = &cg.StackMapFloatVariableInfo{}
	case ast.VariableTypeDouble:
		ret.Verify = &cg.StackMapDoubleVariableInfo{}
	case ast.VariableTypeNull:
		ret.Verify = &cg.StackMapNullVariableInfo{}
	case ast.VariableTypeString:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index:    class.Class.InsertClassConst(javaStringClass),
			Readable: javaStringClass,
		}
	case ast.VariableTypeObject:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index:    class.Class.InsertClassConst(t.Class.Name),
			Readable: t.Class.Name,
		}
	case ast.VariableTypeFunction:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index:    class.Class.InsertClassConst(javaMethodHandleClass),
			Readable: javaMethodHandleClass,
		}
	case ast.VariableTypeMap:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index:    class.Class.InsertClassConst(mapClass),
			Readable: mapClass,
		}
	case ast.VariableTypeArray:
		meta := ArrayMetas[t.Array.Type]
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index:    class.Class.InsertClassConst(meta.className),
			Readable: meta.className,
		}
	case ast.VariableTypeJavaArray:
		d := Descriptor.typeDescriptor(t)
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index:    class.Class.InsertClassConst(d),
			Readable: d,
		}
	default:
		panic(1)
	}
	return ret
}
