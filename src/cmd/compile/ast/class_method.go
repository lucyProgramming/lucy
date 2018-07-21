package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClassMethod struct {
	isCompilerAuto  bool // compile auto method
	Function        *Function
	LoadFromOutSide bool
}

func (m *ClassMethod) IsPublic() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_PUBLIC) != 0
}
func (m *ClassMethod) IsProtected() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_PROTECTED) != 0
}

func (m *ClassMethod) IsStatic() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_STATIC) != 0
}

func (m *ClassMethod) IsPrivate() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_PRIVATE) != 0
}
func (m *ClassMethod) IsFinal() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_FINAL) != 0
}

func (m *ClassMethod) IsFirstStatementCallFatherConstruction() bool {
	if len(m.Function.Block.Statements) == 0 {
		return false
	}
	s := m.Function.Block.Statements[0]
	if s.Type != StatementTypeExpression {
		return false
	}
	e := s.Expression
	if e.Type != ExpressionTypeMethodCall {
		return false
	}
	call := s.Expression.Data.(*ExpressionMethodCall)
	if call.Expression.isThis() == false || call.Name != SUPER {
		return false
	}
	return true
}

func (c *Class) accessInterfaceMethod(from *Position, errs *[]error, name string, call *ExpressionMethodCall, callArgTypes []*Type,
	fromSub bool) (ms []*ClassMethod, matched bool, err error) {
	err = c.loadSelf()
	if err != nil {
		return nil, false, err
	}
	if len(c.Methods[name]) > 0 {
		for _, m := range c.Methods[name] {
			if fromSub {
				if m.IsPrivate() { // break the looking
					return nil, false, fmt.Errorf("method '%s' not found", name)
				}
			}
			var fit bool
			fit, call.VArgs, _ = m.Function.Type.fitCallArgs(from, &call.Args, callArgTypes, nil)
			if fit {
				return []*ClassMethod{m}, true, nil
			} else {
				break
			}
		}
	}
	for _, v := range c.Interfaces {
		err := v.loadSelf()
		if err != nil {
			return nil, false, fmt.Errorf("%s %v", errMsgPrefix(from), err)
		}
		ms, matched, err := v.accessInterfaceMethod(from, errs, name, call, callArgTypes, true)
		if matched {
			return ms, matched, err
		}
	}
	return nil, false, fmt.Errorf("%s method '%s' not found", errMsgPrefix(from), name)
}

/*
	access method lucy style
*/
func (c *Class) accessMethod(from *Position, errs *[]error, name string, call *ExpressionMethodCall,
	callArgTypes []*Type, fromSub bool, fieldMethodHandler **ClassField) (ms []*ClassMethod, matched bool, err error) {
	err = c.loadSelf()
	if err != nil {
		return nil, false, err
	}
	if c.IsJava {
		return c.accessMethodAsJava(from, errs, name, call, callArgTypes, false)
	}
	//TODO:: can be accessed or not ???
	if f := c.Fields[name]; f != nil && f.Type.Type == VariableTypeFunction {
		if fromSub && f.IsPrivate() {
			//cannot access this field
		} else {
			var fit bool
			fit, call.VArgs, _ = c.Fields[name].Type.FunctionType.fitCallArgs(from, &call.Args, callArgTypes, nil)
			if fit {
				*fieldMethodHandler = f
			}
		}
	}
	if len(c.Methods[name]) > 0 {
		for _, m := range c.Methods[name] {
			if fromSub {
				if m.IsPrivate() { // break the looking
					return nil, false, fmt.Errorf("method '%s' not found", name)
				}
			}
			var fit bool
			fit, call.VArgs, _ =
				m.Function.Type.fitCallArgs(from, &call.Args, callArgTypes, m.Function)
			if fit {
				return []*ClassMethod{m}, true, nil
			}
		}
	}
	// don`t try father, when is is construction method
	if name == SpecialMethodInit {
		return nil, false, nil
	}
	err = c.loadSuperClass()
	if err != nil {
		return nil, false, err
	}
	return c.SuperClass.accessMethod(from, errs, name, call, callArgTypes, true, fieldMethodHandler)
}

/*
	access method java style
*/
func (c *Class) accessMethodAsJava(from *Position, errs *[]error, name string, call *ExpressionMethodCall,
	callArgTypes []*Type, fromSub bool) (ms []*ClassMethod, matched bool, err error) {
	for _, m := range c.Methods[name] {
		var fit bool
		fit, call.VArgs, _ =
			m.Function.Type.fitCallArgs(from, &call.Args, callArgTypes, m.Function)
		if fit {
			return []*ClassMethod{m}, true, nil
		}
	}
	// don`t try father, when is is construction method
	if name == SpecialMethodInit {
		return ms, false, nil
	}
	if c.Name == JavaRootClass {
		return ms, false, nil
	}
	err = c.loadSuperClass()
	if err != nil {
		return nil, false, err
	}
	ms_, matched, err := c.SuperClass.accessMethodAsJava(from, errs, name, call, callArgTypes, true)
	if err != nil {
		return ms, false, err
	}
	if matched { // perfect match in father
		return ms_, matched, nil
	}
	return append(ms, ms_...), false, nil // methods have the same name
}
