package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClassMethod struct {
	Func                 *Function
	IsConstructionMethod bool
}

func (m *ClassMethod) IsPublic() bool {
	return (m.Func.AccessFlags & cg.ACC_METHOD_PUBLIC) != 0
}

func (m *ClassMethod) IsStatic() bool {
	return (m.Func.AccessFlags & cg.ACC_METHOD_PUBLIC) != 0
}

func (m *ClassMethod) IsPrivate() bool {
	return (m.Func.AccessFlags & cg.ACC_METHOD_PRIVATE) != 0
}
