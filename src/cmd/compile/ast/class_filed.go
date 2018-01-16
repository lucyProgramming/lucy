package ast

import (
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type ClassField struct {
	VariableDefinition
}

func (f *ClassField) isStatic() bool {
	return (f.AccessFlags & cg.ACC_FIELD_STATIC) != 0
}

func (f *ClassField) isPublic() bool {
	return (f.AccessFlags & cg.ACC_FIELD_STATIC) != 0
}
