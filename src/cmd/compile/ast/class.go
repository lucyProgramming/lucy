package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"strings"
)

/*
	public abstract final class {}
*/
type Class struct {
	IsBuildIn                         bool
	Used                              bool
	resolveFatherCalled               bool
	resolveInterfacesCalled           bool
	resolveFieldsAndMethodsTypeCalled bool
	NotImportedYet                    bool   // not imported
	Name                              string // binary name
	Pos                               *Pos
	FinalPos                          *Pos // final pos
	IsJava                            bool //class imported from CLASSPATH
	Block                             Block
	AccessFlags                       uint16
	Fields                            map[string]*ClassField
	Methods                           map[string][]*ClassMethod
	SuperClassName                    string
	SuperClass                        *Class
	InterfaceNames                    []*NameWithPos
	Interfaces                        []*Class
	LoadFromOutSide                   bool
	StaticBlocks                      []*Block
	Comment                           string
}

func (c *Class) HaveStaticsCodes() bool {
	return len(c.StaticBlocks) > 0
}

func (c *Class) IsInterface() bool {
	return c.AccessFlags&cg.ACC_CLASS_INTERFACE != 0
}
func (c *Class) IsAbstract() bool {
	return c.AccessFlags&cg.ACC_CLASS_ABSTRACT != 0
}
func (c *Class) IsFinal() bool {
	return c.AccessFlags&cg.ACC_CLASS_FINAL != 0
}
func (c *Class) IsPublic() bool {
	return c.AccessFlags&cg.ACC_CLASS_PUBLIC != 0
}

func (c *Class) loadSelf(pos *Pos) error {
	if c.NotImportedYet == false {
		return nil
	}
	load, err := PackageBeenCompile.loadClass(c.Name)
	if err != nil {
		return fmt.Errorf("%s %v", errMsgPrefix(pos), err)
	}
	*c = *load
	return nil
}

func (c *Class) mkDefaultConstruction() {
	if c.IsInterface() {
		return
	}
	if c.Methods == nil {
		c.Methods = make(map[string][]*ClassMethod)
	}
	if len(c.Methods[SpecialMethodInit]) > 0 {
		return
	}
	m := &ClassMethod{}
	m.isCompilerAuto = true
	m.Function = &Function{}
	m.Function.Pos = c.Pos
	m.Function.Block.IsFunctionBlock = true
	m.Function.Name = SpecialMethodInit
	{
		e := &Expression{}
		e.Description = "methodCall"
		e.Type = ExpressionTypeMethodCall
		e.Pos = c.Pos
		call := &ExpressionMethodCall{}
		call.Name = SUPER
		call.Expression = &Expression{
			Type: ExpressionTypeIdentifier,
			Data: &ExpressionIdentifier{
				Name: THIS,
			},
			Pos:         c.Pos,
			Description: "methodCall",
		}
		e.Data = call
		m.Function.Block.Statements = make([]*Statement, 1)
		m.Function.Block.Statements[0] = &Statement{
			Type:       StatementTypeExpression,
			Expression: e,
		}
	}
	c.Methods[SpecialMethodInit] = []*ClassMethod{m}
}

func (c *Class) mkClassInitMethod() {
	if c.HaveStaticsCodes() == false {
		return // no need
	}
	method := &ClassMethod{}
	method.Function = &Function{}
	method.Function.Type.ParameterList = make(ParameterList, 0)
	method.Function.Type.ReturnList = make(ReturnList, 0)
	f := method.Function
	f.Pos = c.Pos
	f.Block.Statements = make([]*Statement, len(c.StaticBlocks))
	for k, _ := range f.Block.Statements {
		s := &Statement{}
		s.Type = StatementTypeBlock
		s.Block = c.StaticBlocks[k]
		f.Block.Statements[k] = s
	}
	f.makeLastReturnStatement()
	f.AccessFlags |= cg.ACC_METHOD_PUBLIC
	f.AccessFlags |= cg.ACC_METHOD_STATIC
	f.AccessFlags |= cg.ACC_METHOD_FINAL
	f.AccessFlags |= cg.ACC_METHOD_BRIDGE
	f.Name = classInitMethod
	f.Block.IsFunctionBlock = true
	if c.Methods == nil {
		c.Methods = make(map[string][]*ClassMethod)
	}
	f.Block.inherit(&c.Block)
	f.Block.InheritedAttribute.Function = f
	f.Block.InheritedAttribute.ClassMethod = method
	c.Methods[f.Name] = []*ClassMethod{method}
}

func (c *Class) resolveFieldsAndMethodsType() []error {
	if c.resolveFieldsAndMethodsTypeCalled {
		return []error{}
	}
	defer func() {
		c.resolveFieldsAndMethodsTypeCalled = true
	}()
	errs := []error{}
	var err error
	for _, v := range c.Fields {
		if v.Name == SUPER {
			errs = append(errs, fmt.Errorf("%s 'super' not allow for field name",
				errMsgPrefix(v.Pos)))
			continue
		}
		err = v.Type.resolve(&c.Block)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, ms := range c.Methods {
		for _, m := range ms {
			if m.IsAbstract() {
				for _, v := range m.Function.Type.ParameterList {
					if v.Expression != nil {
						errs = append(errs, fmt.Errorf("%s abstract method parameter '%s' cannot have default value '%s'",
							errMsgPrefix(v.Pos), v.Name, v.Expression.Description))
					}
				}
				for _, v := range m.Function.Type.ReturnList {
					if v.Expression != nil {
						errs = append(errs, fmt.Errorf("%s abstract method return variable '%s' cannot have default value '%s'",
							errMsgPrefix(v.Pos), v.Name, v.Expression.Description))
					}
				}
			}
			m.Function.Block.inherit(&c.Block)
			m.Function.Block.InheritedAttribute.Function = m.Function
			m.Function.Block.InheritedAttribute.ClassMethod = m
			m.Function.checkParametersAndReturns(&errs, false, m.IsAbstract())
			if m.IsStatic() == false { // bind this
				if m.Function.Block.Variables == nil {
					m.Function.Block.Variables = make(map[string]*Variable)
				}
				m.Function.Block.Variables[THIS] = &Variable{}
				m.Function.Block.Variables[THIS].Name = THIS
				m.Function.Block.Variables[THIS].Pos = m.Function.Pos
				m.Function.Block.Variables[THIS].Type = &Type{
					Type:  VariableTypeObject,
					Class: c,
				}
			}
		}
	}
	return errs
}

func (c *Class) resolveFather() error {
	if c.resolveFatherCalled {
		return nil
	}
	defer func() {
		if c.SuperClassName == "" {
			if c.IsInterface() == false {
				c.SuperClassName = LucyRootClass
			} else {
				c.SuperClassName = JavaRootClass
			}
		}
		c.resolveFatherCalled = true
	}()
	if c.SuperClassName == "" {
		return nil
	}
	if strings.Contains(c.SuperClassName, ".") {
		t := strings.Split(c.SuperClassName, ".")
		i := PackageBeenCompile.getImport(c.Pos.Filename, t[0])
		if i == nil {
			c.SuperClassName = ""
			return fmt.Errorf("%s package name '%s' not imported", errMsgPrefix(c.Pos), t[0])
		}
		i.Used = true
		r, err := PackageBeenCompile.load(i.Import)
		if err != nil {
			c.SuperClassName = ""
			return fmt.Errorf("%s %v", errMsgPrefix(c.Pos), err)
		}
		if p, ok := r.(*Package); ok && p != nil { // if package
			if _, ok := p.Block.NameExists(t[1]); ok == false {
				c.SuperClassName = ""
				return fmt.Errorf("%s class not exists in package '%s' ", errMsgPrefix(c.Pos), t[1])
			}
			if p.Block.Classes == nil || p.Block.Classes[t[1]] == nil {
				c.SuperClassName = ""
				return fmt.Errorf("%s class not exists in package '%s' ", errMsgPrefix(c.Pos), t[1])
			}
			c.SuperClass = p.Block.Classes[t[1]]
		} else if ss, ok := r.(*Class); ok && ss != nil { // must be class now
			t := ss
			c.SuperClassName = t.Name
			c.SuperClass = t
		} else {
			c.SuperClassName = ""
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		}
	} else {
		variableType := Type{}
		variableType.Type = VariableTypeName // naming
		variableType.Name = c.SuperClassName
		variableType.Pos = c.Pos
		err := variableType.resolve(&c.Block)
		if err != nil {
			c.SuperClassName = ""
			return err
		}
		if variableType.Type != VariableTypeObject {
			c.SuperClassName = ""
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		}
		c.SuperClassName = variableType.Class.Name
		c.SuperClass = variableType.Class
	}
	err := c.loadSuperClass(c.Pos)
	if err != nil {
		c.SuperClassName = ""
		return err
	}
	if c.IsInterface() {
		if c.SuperClass.Name == JavaRootClass {
			//nothing
		} else {
			c.SuperClassName = ""
			return fmt.Errorf("%s interface`s super-class must be '%s'",
				errMsgPrefix(c.Pos), JavaRootClass)
		}
	}
	return nil
}

func (c *Class) resolveInterfaces() []error {
	if c.resolveInterfacesCalled {
		return nil
	}
	defer func() {
		c.resolveInterfacesCalled = true
	}()
	errs := []error{}
	for _, i := range c.InterfaceNames {
		t := &Type{}
		t.Type = VariableTypeName
		t.Pos = i.Pos
		t.Name = i.Name
		err := t.resolve(&c.Block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if t.Type != VariableTypeObject || t.Class.IsInterface() == false {
			errs = append(errs, fmt.Errorf("%s '%s' is not a interface",
				errMsgPrefix(i.Pos), i.Name))
			continue
		}
		c.Interfaces = append(c.Interfaces, t.Class)
	}
	return errs
}

func (c *Class) implementMethod(pos *Pos, m *ClassMethod, fromSub bool, errs *[]error) *ClassMethod {
	if c.Methods != nil {
		for _, v := range c.Methods[m.Function.Name] {
			if v.IsAbstract() {
				continue
			}
			if fromSub && v.ableAccessFromSubClass() == false {
				return nil
			}
			if v.Function.Type.equal(&m.Function.Type) {
				return v
			}
		}
	}
	//no same name method at current class
	if c.Name == JavaRootClass {
		return nil
	}
	err := c.loadSuperClass(pos)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	} else {
		//trying find fathte`s implementation
		return c.SuperClass.implementMethod(pos, m, true, errs)
	}
	return nil
}

func (c *Class) haveSuperClass(pos *Pos, superclassName string) (bool, error) {
	err := c.loadSelf(pos)
	if err != nil {
		return false, err
	}
	if c.Name == superclassName {
		return true, nil
	}
	if c.Name == JavaRootClass {
		return false, nil
	}
	err = c.loadSuperClass(pos)
	if err != nil {
		return false, err
	}
	return c.SuperClass.haveSuperClass(pos, superclassName) // check father is implements
}

func (c *Class) implementedInterface(pos *Pos, inter string) (bool, error) {
	err := c.loadSelf(pos)
	if err != nil {
		return false, err
	}
	for _, v := range c.Interfaces {
		if v.Name == inter {
			return true, nil
		}
		im, _ := v.implementedInterface(pos, inter)
		if im {
			return im, nil
		}
	}
	if c.Name == JavaRootClass {
		return false, nil
	}
	err = c.loadSuperClass(pos)
	if err != nil {
		return false, err
	}
	return c.SuperClass.implementedInterface(pos, inter) // check father is implements
}

func (c *Class) loadSuperClass(pos *Pos) error {
	if c.SuperClass != nil {
		return nil
	}
	if c.Name == JavaRootClass {
		return fmt.Errorf("%s root class already", errMsgPrefix(pos))
	}
	class, err := PackageBeenCompile.loadClass(c.SuperClassName)
	if err != nil {
		return fmt.Errorf("%s %v", errMsgPrefix(pos), err)
	}
	c.SuperClass = class
	return nil
}

func (c *Class) accessConstructionFunction(pos *Pos, errs *[]error, newCase *ExpressionNew,
	callFatherCase *ExpressionMethodCall, callArgs []*Type) (ms []*ClassMethod, matched bool, err error) {
	err = c.loadSelf(pos)
	if err != nil {
		return nil, false, err
	}

	var args *CallArgs
	if newCase != nil {
		args = &newCase.Args
	} else {
		args = &callFatherCase.Args
	}
	for _, v := range c.Methods[SpecialMethodInit] {
		vArgs, err := v.Function.Type.fitArgs(pos, args, callArgs, v.Function)
		if err == nil {
			if newCase != nil {
				newCase.VArgs = vArgs
			} else {
				callFatherCase.VArgs = vArgs
			}
			return []*ClassMethod{v}, true, nil
		} else {
			ms = append(ms, v)
		}
	}
	return ms, false, nil
}

//func (c *Class) classAccessAble(pos *Pos) error {
//	if c.LoadFromOutSide == false {
//		return nil
//	}
//	err := c.loadSelf(pos)
//	if err != nil {
//		return err
//	}
//	if c.IsPublic() == false {
//		return fmt.Errorf("%s class '%s' is not public",
//			errMsgPrefix(pos), c.Name)
//	}
//	return nil
//}

/*
	ret is *ClassField or *ClassMethod
*/
func (c *Class) getFieldOrMethod(pos *Pos, name string, fromSub bool) (interface{}, error) {
	err := c.loadSelf(pos)
	if err != nil {
		return nil, err
	}
	notFoundErr := fmt.Errorf("%s field or method named '%s' not found", errMsgPrefix(pos), name)
	if c.Fields != nil && nil != c.Fields[name] {
		if fromSub && c.Fields[name].ableAccessFromSubClass() == false {
			// private field
			// break find
			return nil, notFoundErr
		} else {
			return c.Fields[name], nil
		}
	}
	if c.Methods != nil && nil != c.Methods[name] {
		m := c.Methods[name][0]
		if fromSub && m.ableAccessFromSubClass() == false {
			return nil, notFoundErr
		} else {
			return m, nil
		}
	}
	if c.Name == JavaRootClass { // root class
		return nil, notFoundErr
	}
	err = c.loadSuperClass(pos)
	if err != nil {
		return nil, err
	}
	return c.SuperClass.getFieldOrMethod(pos, name, true)
}
