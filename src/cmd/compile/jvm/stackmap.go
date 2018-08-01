package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StackMapState struct {
	Locals []*cg.StackMapVerificationTypeInfo
	Stacks []*cg.StackMapVerificationTypeInfo
}

func (stackMapState *StackMapState) appendLocals(class *cg.ClassHighLevel, v *ast.Type) {
	stackMapState.Locals = append(stackMapState.Locals,
		stackMapState.newStackMapVerificationTypeInfo(class, v))
}

func (stackMapState *StackMapState) addTop(absent *StackMapState) {
	if stackMapState == absent {
		return
	}
	length := len(absent.Locals) - len(stackMapState.Locals)
	if length == 0 {
		return
	}
	oldLength := len(stackMapState.Locals)
	verify := &cg.StackMapVerificationTypeInfo{}
	verify.Verify = &cg.StackMapTopVariableInfo{}
	for i := 0; i < length; i++ {
		tt := absent.Locals[i+oldLength].Verify
		_, isDouble := tt.(*cg.StackMapDoubleVariableInfo)
		_, isLong := tt.(*cg.StackMapLongVariableInfo)
		if isDouble || isLong {
			stackMapState.Locals = append(stackMapState.Locals, verify, verify)
		} else {
			stackMapState.Locals = append(stackMapState.Locals, verify)
		}
	}
}

func (stackMapState *StackMapState) newObjectVariableType(name string) *ast.Type {
	ret := &ast.Type{}
	ret.Type = ast.VariableTypeObject
	ret.Class = &ast.Class{}
	ret.Class.Name = name
	return ret
}

func (stackMapState *StackMapState) popStack(pop int) {
	if pop == 0 {
		return
	}
	if pop < 0 {
		panic("negative pop")
	}
	if len(stackMapState.Stacks) == 0 {
		panic("already 0")
	}
	stackMapState.Stacks = stackMapState.Stacks[:len(stackMapState.Stacks)-pop]
}
func (stackMapState *StackMapState) pushStack(class *cg.ClassHighLevel, v *ast.Type) {
	if stackMapState == nil {
		panic("s is nil")
	}
	stackMapState.Stacks = append(stackMapState.Stacks, stackMapState.newStackMapVerificationTypeInfo(class, v))
}
func (stackMapState *StackMapState) FromLast(last *StackMapState) *StackMapState {
	stackMapState.Locals = make([]*cg.StackMapVerificationTypeInfo, len(last.Locals))
	copy(stackMapState.Locals, last.Locals)
	return stackMapState
}

func (stackMapState *StackMapState) newStackMapVerificationTypeInfo(class *cg.ClassHighLevel,
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
			Index: class.Class.InsertClassConst(javaStringClass),
		}
	case ast.VariableTypeObject:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(t.Class.Name),
		}
	case ast.VariableTypeFunction:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(javaMethodHandleClass),
		}
	case ast.VariableTypeMap:
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(javaMapClass),
		}
	case ast.VariableTypeArray:
		meta := ArrayMetas[t.Array.Type]
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(meta.className),
		}
	case ast.VariableTypeJavaArray:
		d := Descriptor.typeDescriptor(t)
		ret.Verify = &cg.StackMapObjectVariableInfo{
			Index: class.Class.InsertClassConst(d),
		}
	default:
		panic(1)
	}
	return ret
}
