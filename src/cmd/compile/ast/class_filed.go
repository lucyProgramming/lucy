package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClassField struct {
	Class                  *Class
	DefaultValueExpression *Expression
	Name                   string
	Type                   *Type
	Pos                    *Pos
	Comment                string
	AccessFlags            uint16
	JvmDescriptor          string      // jvm
	DefaultValue           interface{} // value base on type
}

func (f *ClassField) IsStatic() bool {
	return (f.AccessFlags & cg.AccFieldStatic) != 0
}
func (f *ClassField) IsPublic() bool {
	return (f.AccessFlags & cg.AccFieldPublic) != 0
}
func (f *ClassField) IsProtected() bool {
	return (f.AccessFlags & cg.AccFieldProtected) != 0
}
func (f *ClassField) IsPrivate() bool {
	return (f.AccessFlags & cg.AccFieldPrivate) != 0
}
func (f *ClassField) ableAccessFromSubClass() bool {
	return f.IsPublic() ||
		f.IsProtected()
}

func (c *Class) getField(pos *Pos, name string, fromSub bool) (*ClassField, error) {
	err := c.loadSelf(pos)
	if err != nil {
		return nil, err
	}
	notFoundErr := fmt.Errorf("%s field named '%s' not found",
		pos.ErrMsgPrefix(), name)
	if c.Fields != nil && nil != c.Fields[name] {
		if fromSub && c.Fields[name].ableAccessFromSubClass() == false {
			// private field
			return nil, notFoundErr
		} else {
			return c.Fields[name], nil
		}
	}
	if c.Name == JavaRootClass { // root class
		return nil, notFoundErr
	}
	err = c.loadSuperClass(pos)
	if err != nil {
		return nil, err
	}
	if c.SuperClass == nil {
		return nil, notFoundErr
	}
	return c.SuperClass.getField(pos, name, true)
}
