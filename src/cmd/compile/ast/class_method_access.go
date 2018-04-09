package ast

import "fmt"

/*
	access method lucy style
*/
func (c *Class) accessMethod(name string, args []*VariableType, callArgs *CallArgs, fromsub bool) (ms []*ClassMethod, matched bool, err error) {
	if c.IsJava {
		return c.accessMethodAsJava(name, args, false)
	}
	if len(c.Methods[name]) > 0 {
		m := c.Methods[name][0]
		if fromsub {
			if m.IsPrivate() {
				return nil, false, fmt.Errorf("method '%s' not found", name)
			}
		}
		if len(args) > len(m.Func.Typ.ParameterList) {
			return nil, false, fmt.Errorf("too many paramaters to call method '%s'", m.Func.Name)
		}
		if len(args) < len(m.Func.Typ.ParameterList) {
			if m.Func.HaveDefaultValue {
				if len(args) < m.Func.DefaultValueStartAt {
					return []*ClassMethod{m}, false, fmt.Errorf("too few paramaters to call method '%s'", m.Func.Name)
				}
				for i := len(args); i < len(m.Func.Typ.ParameterList); i++ {
					*callArgs = append(*callArgs, m.Func.Typ.ParameterList[i].Expression)
				}
			} else { // no default value
				if len(args) < len(m.Func.Typ.ParameterList) {
					return nil, false, fmt.Errorf("too few paramaters to call method '%s'", m.Func.Name)
				}
			}
		}
		for k, v := range m.Func.Typ.ParameterList {
			if k < len(args) {
				if !v.Typ.TypeCompatible(args[k]) {
					return nil, false, fmt.Errorf("type '%s' is not compatible with '%s'",
						v.Typ.TypeString(), args[k].TypeString())
				}
			}
		}
		return []*ClassMethod{m}, true, nil
	}
	err = c.loadSuperClass()
	if err != nil {
		return nil, false, err
	}
	return c.SuperClass.accessMethod(name, args, callArgs, true)
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
	ms_, matched, err := c.SuperClass.accessMethodAsJava(name, args, true)
	if err != nil {
		return ms, false, err
	}
	if matched { // perfect match in father
		return ms_, matched, nil
	}
	return append(ms, ms_...), false, nil
}

func (c *Class) matchContructionFunction(args []*VariableType, callArgs *CallArgs) (ms []*ClassMethod, matched bool, err error) {
	return c.accessMethod("<init>", args, callArgs, false)
}
