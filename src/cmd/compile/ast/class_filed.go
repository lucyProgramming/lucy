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

/*
	ret is *ClassField or *ClassMethod
*/
func (c *Class) getFieldOrMethod(from *Pos, name string, fromSub bool) (interface{}, error) {
	err := c.loadSelf()
	if err != nil {
		return nil, err
	}
	notFoundErr := fmt.Errorf("%s field or method named '%s' not found", errMsgPrefix(from), name)
	if c.Fields != nil && nil != c.Fields[name] {
		if fromSub && c.Fields[name].IsPrivate() {
			// private field
			return nil, notFoundErr
		} else {
			return c.Fields[name], nil
		}
	}
	if c.Methods != nil && nil != c.Methods[name] {
		m := c.Methods[name][0]
		if fromSub && m.IsPrivate() {
			return nil, notFoundErr
		} else {
			return m, nil
		}
	}
	if c.Name == JavaRootClass { // root class
		return nil, notFoundErr
	}
	err = c.loadSuperClass()
	if err != nil {
		return nil, err
	}
	return c.SuperClass.getFieldOrMethod(from, name, true)
}

func (c *Class) accessField(from *Pos, name string, fromSub bool) (*ClassField, error) {
	err := c.loadSelf()
	if err != nil {
		return nil, err
	}
	notFoundErr := fmt.Errorf("field named '%s' not found", name)
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
		return nil, err
	}
	return c.SuperClass.accessField(from, name, true)
}
