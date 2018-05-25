package ast

import (
	"fmt"
)

func (c *Class) accessField(name string, fromSub bool) (f *ClassField, err error) {
	err = c.loadSelf()
	if err != nil {
		return
	}
	notFoundErr := fmt.Errorf("field '%s' not found", name)
	if c.Fields != nil && nil != c.Fields[name] {
		if fromSub && c.Fields[name].IsPrivate() {
			// private field
			return nil, notFoundErr
		} else {
			return c.Fields[name], nil
		}
	}
	if c.Name == JAVA_ROOT_CLASS { // root class
		return nil, notFoundErr
	}
	err = c.loadSuperClass()
	if err != nil {
		return
	}
	return c.SuperClass.accessField(name, true)
}
