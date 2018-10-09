package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClassMethod struct {
	IsCompilerAuto  bool
	Function        *Function
	LoadFromOutSide bool
}

func (m *ClassMethod) narrowDownAccessRange(implementation *ClassMethod) bool {
	if m.IsPublic() {
		return !implementation.IsPublic()
	}
	if m.IsProtected() {
		return implementation.IsPrivate() ||
			implementation.isAccessFlagDefault()
	}
	if m.isAccessFlagDefault() {
		return implementation.IsPrivate()
	}
	return false
}

func (m *ClassMethod) accessString() string {
	if m.IsPublic() {
		return "public"
	}
	if m.IsProtected() {
		return "protected"
	}
	if m.IsPrivate() {
		return "private"
	}
	return `default ""` //TODO:: ???
}

func (m *ClassMethod) isAccessFlagDefault() bool {
	return m.Function.AccessFlags&
		(cg.ACC_METHOD_PUBLIC|
			cg.ACC_METHOD_PRIVATE|
			cg.ACC_METHOD_PROTECTED) == 0
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
func (m *ClassMethod) IsAbstract() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_ABSTRACT) != 0
}

func (m *ClassMethod) ableAccessFromSubClass() bool {
	return m.IsPublic() ||
		m.IsProtected()
}

func (m *ClassMethod) checkModifierOk() []error {
	errs := []error{}
	if m.IsAbstract() && m.IsFinal() {
		errs = append(errs, fmt.Errorf("%s abstract method cannot be final",
			errMsgPrefix(m.Function.Pos)))
	}
	if m.IsAbstract() && m.IsPrivate() {
		errs = append(errs, fmt.Errorf("%s abstract method cannot be private",
			errMsgPrefix(m.Function.Pos)))
	}
	if m.IsAbstract() && m.Function.Name == SpecialMethodInit {
		errs = append(errs, fmt.Errorf("%s construction method cannot be abstract",
			errMsgPrefix(m.Function.Pos)))
	}
	return errs
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
	if call.Expression.IsIdentifier(THIS) == false || call.Name != SUPER {
		return false
	}
	return true
}
func (c *Class) accessInterfaceObjectMethod(pos *Pos, errs *[]error, name string, call *ExpressionMethodCall, callArgTypes []*Type,
	fromSub bool) (ms []*ClassMethod, matched bool, err error) {
	ms, matched, err = c.accessInterfaceMethod(pos, errs, name, call, callArgTypes, fromSub)
	if err != nil {
		return nil, false, err
	}
	if matched {
		return ms, matched, err
	}
	err = c.loadSuperClass(pos)
	if err != nil {
		return nil, false, err
	}
	return c.SuperClass.accessMethod(pos, errs, call, callArgTypes, fromSub, nil)
}

func (c *Class) accessInterfaceMethod(pos *Pos, errs *[]error, name string, call *ExpressionMethodCall, callArgTypes []*Type,
	fromSub bool) (ms []*ClassMethod, matched bool, err error) {
	err = c.loadSelf(pos)
	if err != nil {
		return nil, false, err
	}
	if nil != c.Methods {
		for _, m := range c.Methods[name] {
			if fromSub && m.ableAccessFromSubClass() == false {
				continue
			}
			call.VArgs, err = m.Function.Type.fitArgs(pos, &call.Args, callArgTypes, nil)
			if err == nil {
				return []*ClassMethod{m}, true, nil
			} else {
				return nil, false, err
			}
		}
	}
	for _, v := range c.Interfaces {
		err := v.loadSelf(pos)
		if err != nil {
			return nil, false, err
		}
		ms2, matched2, err2 := v.accessInterfaceMethod(pos, errs, name, call, callArgTypes, true)
		if err2 != nil {
			return nil, false, err2
		}
		if matched {
			return ms2, matched2, nil
		}
	}
	return nil, false, nil // no found , no error
}

/*
	access method lucy style
*/
func (c *Class) accessMethod(pos *Pos, errs *[]error, call *ExpressionMethodCall,
	callArgTypes []*Type, fromSub bool, fieldMethodHandler **ClassField) (ms []*ClassMethod, matched bool, err error) {
	err = c.loadSelf(pos)
	if err != nil {
		return nil, false, err
	}

	if c.IsJava {
		return c.accessMethodAsJava(pos, errs, call, callArgTypes, false)
	}
	//TODO:: can be accessed or not ???
	if f := c.Fields[call.Name]; f != nil &&
		f.Type.Type == VariableTypeFunction &&
		fieldMethodHandler != nil {
		if fromSub && f.ableAccessFromSubClass() == false {
			//cannot access this field
		} else {
			call.VArgs, err = c.Fields[call.Name].Type.FunctionType.fitArgs(pos, &call.Args,
				callArgTypes, nil)
			if err == nil {
				*fieldMethodHandler = f
				matched = true
				return
			} else {
				return nil, false, err
			}
		}
	}
	if len(c.Methods[call.Name]) > 0 {
		for _, m := range c.Methods[call.Name] {
			if fromSub && m.ableAccessFromSubClass() == false {
				return nil, false, fmt.Errorf("%s method '%s' not found",
					errMsgPrefix(pos), call.Name)
			}
			call.VArgs, err = m.Function.Type.fitArgs(pos, &call.Args,
				callArgTypes, m.Function)
			if err == nil {
				return []*ClassMethod{m}, true, nil
			} else {
				return []*ClassMethod{m}, false, err
			}
		}
	}
	err = c.loadSuperClass(pos)
	if err != nil {
		return ms, false, err
	}
	return c.SuperClass.accessMethod(pos, errs, call,
		callArgTypes, true, fieldMethodHandler)
}

/*
	access method java style
*/
func (c *Class) accessMethodAsJava(pos *Pos, errs *[]error, call *ExpressionMethodCall,
	callArgTypes []*Type, fromSub bool) (ms []*ClassMethod, matched bool, err error) {
	if c.Methods != nil {
		for _, m := range c.Methods[call.Name] {
			if fromSub == true && m.ableAccessFromSubClass() == false {
				//cannot access from sub
				continue
			}
			call.VArgs, err = m.Function.Type.fitArgs(pos, &call.Args, callArgTypes, m.Function)
			if err == nil {
				return []*ClassMethod{m}, true, nil
			}
		}
	}
	if c.Name == JavaRootClass {
		return ms, false, nil
	}
	err = c.loadSuperClass(pos)
	if err != nil {
		return nil, false, err
	}
	ms_, matched, err := c.SuperClass.accessMethodAsJava(pos, errs, call, callArgTypes, true)
	if err != nil {
		return ms, false, err
	}
	if matched { // perfect match in father
		return ms_, matched, nil
	}
	return append(ms, ms_...), false, nil // methods have the same name
}

func (m *ClassMethod) implementationMethodIsOk(pos *Pos, implementation *ClassMethod) error {
	if implementation.Function.Pos != nil {
		pos = implementation.Function.Pos
	}
	if implementation.IsStatic() {
		return fmt.Errorf("%s method '%s' is static", errMsgPrefix(pos), m.Function.Name)
	}
	if m.narrowDownAccessRange(implementation) {
		return fmt.Errorf("%s implementation of method '%s' should not narrow down access range, '%s' -> '%s'",
			errMsgPrefix(pos), m.Function.Name, m.accessString(), implementation.accessString())
	}
	return nil
}
