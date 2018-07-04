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

/*
	access method lucy style
*/
func (c *Class) accessMethod(from *Position, errs *[]error, name string, args []*Type,
	callArgs *CallArgs, fromSub bool) (ms []*ClassMethod, matched bool, err error) {
	err = c.loadSelf()
	if err != nil {
		return nil, false, err
	}
	if c.IsJava {
		return c.accessMethodAsJava(from, errs, name, args, false)
	}
	if len(c.Methods[name]) > 0 {
		for _, m := range c.Methods[name] {
			if fromSub {
				if m.IsPrivate() { // break the looking
					return nil, false, fmt.Errorf("method '%s' not found", name)
				}
			}
			if len(args) > len(m.Function.Type.ParameterList) {
				errMsg := fmt.Sprintf("too many paramaters to call function '%s':\n", name)
				errMsg += fmt.Sprintf("\thave %s\n", m.Function.badParameterMsg(name, args))
				errMsg += fmt.Sprintf("\twant %s\n", m.Function.readableMsg())
				return nil, false, fmt.Errorf(errMsg)
			}
			if len(args) < len(m.Function.Type.ParameterList) {
				if m.Function.HaveDefaultValue && len(args) >= m.Function.DefaultValueStartAt && callArgs != nil {
					for i := len(args); i < len(m.Function.Type.ParameterList); i++ {
						*callArgs = append(*callArgs, m.Function.Type.ParameterList[i].Expression)
					}
				} else { // no default value
					errMsg := fmt.Sprintf("too few paramaters to call function '%s'\n", name)
					errMsg += fmt.Sprintf("\thave %s\n", m.Function.badParameterMsg(m.Function.Name, args))
					errMsg += fmt.Sprintf("\twant %s\n", m.Function.readableMsg())
					return nil, false, fmt.Errorf(errMsg)
				}
			} else {
				convertLiteralExpressionsToNeeds(*callArgs, m.Function.Type.getParameterTypes(), args)
			}
			for k, v := range m.Function.Type.ParameterList {
				if k < len(args) {
					if args[k] != nil && !v.Type.Equal(errs, args[k]) {
						errMsg := fmt.Sprintf("cannot use '%s' as '%s'\n", args[k].TypeString(), v.Type.TypeString())
						errMsg += fmt.Sprintf("\thave %s\n", m.Function.badParameterMsg(m.Function.Name, args))
						errMsg += fmt.Sprintf("\twant %s\n", m.Function.readableMsg())
						return nil, false, fmt.Errorf(errMsg)
					}
				}
			}
			return []*ClassMethod{m}, true, nil
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
	return c.SuperClass.accessMethod(from, errs, name, args, callArgs, true)
}

/*
	access method java style
*/
func (c *Class) accessMethodAsJava(from *Position, errs *[]error, name string,
	args []*Type, fromSub bool) (ms []*ClassMethod, matched bool, err error) {
	for _, v := range c.Methods[name] {
		if len(v.Function.Type.ParameterList) != len(args) {
			if fromSub == false || v.IsPublic() || v.IsProtected() {
				ms = append(ms, v)
			}
			continue
		}
		noError := true
		for kk, vv := range v.Function.Type.ParameterList {
			if args[kk] != nil && vv.Type.Equal(errs, args[kk]) == false {
				noError = false
				ms = append(ms, v)
				break
			}
		}
		if noError {
			return []*ClassMethod{v}, true, nil
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
	ms_, matched, err := c.SuperClass.accessMethodAsJava(from, errs, name, args, true)
	if err != nil {
		return ms, false, err
	}
	if matched { // perfect match in father
		return ms_, matched, nil
	}
	return append(ms, ms_...), false, nil // methods have the same name
}

func (c *Class) matchConstructionFunction(from *Position, errs *[]error, args []*Type,
	callArgs *CallArgs) (ms []*ClassMethod, matched bool, err error) {
	return c.accessMethod(from, errs, SpecialMethodInit, args, callArgs, false)
}
