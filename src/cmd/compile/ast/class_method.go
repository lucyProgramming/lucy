package ast

import (
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type ClassMethod struct {
	Func *Function
}

func (m *ClassMethod) isPublic() bool {
	return (m.Func.AccessFlags & cg.ACC_METHOD_PUBLIC) != 0
}
