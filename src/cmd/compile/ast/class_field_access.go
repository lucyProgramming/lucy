package ast

func (c *Class) accessField(name string, fromSub bool) (f *ClassField, err error) {
	if c.Fields != nil && nil != c.Fields[name] {
		if fromSub { // private field
			if c.Fields[name].IsPrivate() {
				return nil, nil
			}
		} else {
			return c.Fields[name], nil
		}
	}
	err = c.loadSuperClass()
	if err != nil {
		return
	}
	return c.SuperClass.accessField(name, true)
}
