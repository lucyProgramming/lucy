package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClassField struct {
	Variable
	Class           *Class
	LoadFromOutSide bool
	DefaultValue    interface{} // value base on type
}

func (f *ClassField) IsStatic() bool {
	return (f.AccessFlags & cg.ACC_FIELD_STATIC) != 0
}

func (f *ClassField) IsPublic() bool {
	return (f.AccessFlags & cg.ACC_FIELD_PUBLIC) != 0
}
func (f *ClassField) IsProtected() bool {
	return (f.AccessFlags & cg.ACC_FIELD_PROTECTED) != 0
}
func (f *ClassField) IsPrivate() bool {
	return (f.AccessFlags & cg.ACC_FIELD_PRIVATE) != 0
}

func (c *Class) accessField(name string, fromSub bool) (f *ClassField, err error) {
	err = c.loadSelf()
	if err != nil {
		return
	}
	notFoundErr := fmt.Errorf("field or method named '%s' not found", name)
	if c.Fields != nil && nil != c.Fields[name] {
		if fromSub && c.Fields[name].IsPrivate() {
			// private field
			return nil, notFoundErr
		} else {
			return c.Fields[name], nil
		}
	}
	if c.Name == JavaRootClass { // root class
		return nil, notFoundErr
	}
	err = c.loadSuperClass()
	if err != nil {
		return
	}
	return c.SuperClass.accessField(name, true)
}
