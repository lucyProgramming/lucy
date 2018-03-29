package ast

import "fmt"

/*
	access method lucy style
*/
func (c *Class) accessMethod(name string, args []*VariableType) (ms []*ClassMethod, err error) {
	if c.Package.Kind == PACKAGE_KIND_JAVA {
		return c.accessMethodAsJava(name, args)
	}
	if len(c.Methods[name]) > 0 {
		m := c.Methods[name][0]
		if len(args) > len(m.Func.Typ.ParameterList) {
			return nil, fmt.Errorf("too mary argument to call")
		}
		if len(args) < len(m.Func.Typ.ParameterList) {
			return nil, fmt.Errorf("too mary argument to call")
		}
		for k, v := range m.Func.Typ.ParameterList {
			if v.Typ.TypeCompatible(args[k]) == false {
				return nil, fmt.Errorf("cannot use '%s' as '%s'", args[k].TypeString(), v.Typ.TypeString())
			}
		}
		return []*ClassMethod{m}, nil
	}
	if c.SuperClass == nil {
		c.loadSuperClass()
	}
	return c.SuperClass.accessMethod(name, args)
}

/*
	access method java style
*/
func (c *Class) accessMethodAsJava(name string, args []*VariableType) (ms []*ClassMethod, err error) {
	for _, v := range c.Methods[name] {
		if len(v.Func.Typ.ParameterList) == len(args) {
			ms = append(ms, v)
			continue
		}
		noError := true
		if len(v.Func.Typ.ParameterList) == len(args) {
			for kk, vv := range v.Func.Typ.ParameterList {
				if vv.Typ.Equal(args[kk]) == false {
					noError = false
					break
				}
			}
		}
		if noError {
			return []*ClassMethod{v}, nil
		}
	}
	if c.Name == JAVA_ROOT_CLASS {
		return ms, nil
	}
	// here is no match
	if c.SuperClass == nil {
		err = c.loadSuperClass()
		if err != nil {
			return nil, err
		}
	}
	ms_, err := c.SuperClass.accessMethod(name, args)
	if err != nil {
		return ms, err
	}
	if len(ms_) == 1 { // perfect match in father
		return ms_, nil
	}
	return append(ms, ms_...), nil
}
