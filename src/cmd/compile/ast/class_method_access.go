package ast

import "fmt"

/*
	access method lucy style
*/
func (c *Class) accessMethod(name string, args []*VariableType, fromsub bool) (ms []*ClassMethod, matched bool, err error) {
	if c.IsJava {
		return c.accessMethodAsJava(name, args, false)
	}
	if len(c.Methods[name]) > 0 {
		m := c.Methods[name][0]
		if len(args) > len(m.Func.Typ.ParameterList) {
			return nil, false, fmt.Errorf("too mary argument to call")
		}
		if len(args) < len(m.Func.Typ.ParameterList) {
			return nil, false, fmt.Errorf("too few argument to call")
		}
		for k, v := range m.Func.Typ.ParameterList {
			if v.Typ.TypeCompatible(args[k]) == false {
				return nil, false, fmt.Errorf("cannot use '%s' as '%s'", args[k].TypeString(), v.Typ.TypeString())
			}
		}
		return []*ClassMethod{m}, true, nil
	}
	if c.SuperClass == nil {
		c.loadSuperClass()
	}
	return c.SuperClass.accessMethod(name, args, false)
}

/*
	access method java style
*/
func (c *Class) accessMethodAsJava(name string, args []*VariableType, fromsub bool) (ms []*ClassMethod, matched bool, err error) {
	for _, v := range c.Methods[name] {
		if len(v.Func.Typ.ParameterList) != len(args) {
			if fromsub == false || v.IsPublic() || v.IsProtected() {
				ms = append(ms, v)
			}
			continue
		}
		noError := true
		for kk, vv := range v.Func.Typ.ParameterList {
			if vv.Typ.Equal(args[kk]) == false {
				noError = false
				break
			}
		}
		if noError {
			return []*ClassMethod{v}, true, nil
		}
	}
	if c.Name == JAVA_ROOT_CLASS {
		return ms, false, nil
	}
	err = c.loadSuperClass()
	if err != nil {
		return nil, false, err
	}
	ms_, matched, err := c.SuperClass.accessMethod(name, args, true)
	if err != nil {
		return ms, false, err
	}
	if matched { // perfect match in father
		return ms_, matched, nil
	}
	return append(ms, ms_...), false, nil
}
