package ast

func (c *Class) accessField(name string, fromSub bool) (f *ClassField, err error) {
	if c.Fields != nil && c.Fields[name] != nil {
		f = c.Fields[name]
		if fromSub {
			if f.IsPublic() == false && f.IsProtected() == false {
				return nil, nil
			}
		}
		return
	}
	// not found in current class
	if c.SuperClassName == "" {
		c.SuperClassName = LUCY_ROOT_CLASS
	}
	err = c.loadSuperClass()
	if err != nil {
		return
	}
	return c.SuperClass.accessField(name, true)
}
