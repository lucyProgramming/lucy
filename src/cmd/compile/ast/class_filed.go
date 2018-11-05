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

func (this *Class) getField(pos *Pos, name string, fromSub bool) (*ClassField, error) {
	err := this.loadSelf(pos)
	if err != nil {
		return nil, err
	}
	notFoundErr := fmt.Errorf("%s field named '%s' not found",
		pos.ErrMsgPrefix(), name)
	if this.Fields != nil && nil != this.Fields[name] {
		if fromSub && this.Fields[name].ableAccessFromSubClass() == false {
			// private field
			return nil, notFoundErr
		} else {
			return this.Fields[name], nil
		}
	}
	if this.Name == JavaRootClass { // root class
		return nil, notFoundErr
	}
	err = this.loadSuperClass(pos)
	if err != nil {
		return nil, err
	}
	if this.SuperClass == nil {
		return nil, notFoundErr
	}
	return this.SuperClass.getField(pos, name, true)
}
