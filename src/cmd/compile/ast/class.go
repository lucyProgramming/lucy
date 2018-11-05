package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	public abstract final class {}
	unlike function class not accept capture
*/
type Class struct {
	IsGlobal                          bool
	IsBuildIn                         bool
	Used                              bool
	resolveFatherCalled               bool
	loadSuperClassCalled              bool
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
	SuperClassName                    *NameWithPos
	SuperClass                        *Class
	InterfaceNames                    []*NameWithPos
	Interfaces                        []*Class
	LoadFromOutSide                   bool
	StaticBlocks                      []*Block
	Comment                           string
	closure                           Closure
}

func (this *Class) HaveStaticsCodes() bool {
	return len(this.StaticBlocks) > 0
}

func (this *Class) IsInterface() bool {
	return this.AccessFlags&cg.AccClassInterface != 0
}
func (this *Class) IsAbstract() bool {
	return this.AccessFlags&cg.AccClassAbstract != 0
}
func (this *Class) IsFinal() bool {
	return this.AccessFlags&cg.AccClassFinal != 0
}
func (this *Class) IsPublic() bool {
	return this.AccessFlags&cg.AccClassPublic != 0
}

func (this *Class) loadSelf(pos *Pos) error {
	if this.NotImportedYet == false { // current compile class
		return nil
	}
	this.NotImportedYet = false
	load, err := PackageBeenCompile.loadClass(this.Name)
	if err != nil {
		return fmt.Errorf("%s %v", pos.ErrMsgPrefix(), err)
	}
	*this = *load
	return nil
}

func (this *Class) mkDefaultConstruction() {
	if this.IsInterface() {
		return
	}
	if this.Methods == nil {
		this.Methods = make(map[string][]*ClassMethod)
	}
	if len(this.Methods[SpecialMethodInit]) > 0 {
		return
	}
	m := &ClassMethod{}
	m.IsCompilerAuto = true
	m.Function = &Function{}
	m.Function.AccessFlags |= cg.AccMethodPublic
	m.Function.Pos = this.Pos
	m.Function.Block.IsFunctionBlock = true
	m.Function.Block.Fn = m.Function
	m.Function.Name = SpecialMethodInit
	{
		e := &Expression{}
		e.Op = "methodCall"
		e.Type = ExpressionTypeMethodCall
		e.Pos = this.Pos
		call := &ExpressionMethodCall{}
		call.Name = SUPER
		call.Expression = &Expression{
			Type: ExpressionTypeIdentifier,
			Data: &ExpressionIdentifier{
				Name: ThisPointerName,
			},
			Pos: this.Pos,
			Op:  "methodCall",
		}
		e.Data = call
		m.Function.Block.Statements = make([]*Statement, 1)
		m.Function.Block.Statements[0] = &Statement{
			Type:       StatementTypeExpression,
			Expression: e,
		}
	}
	this.Methods[SpecialMethodInit] = []*ClassMethod{m}
}

func (this *Class) mkClassInitMethod() {
	if this.HaveStaticsCodes() == false {
		return // no need
	}
	method := &ClassMethod{}
	method.Function = &Function{}
	method.Function.Type.ParameterList = make(ParameterList, 0)
	method.Function.Type.ReturnList = make(ReturnList, 0)
	f := method.Function
	f.Pos = this.Pos
	f.Block.Statements = make([]*Statement, len(this.StaticBlocks))
	for k, _ := range f.Block.Statements {
		s := &Statement{}
		s.Type = StatementTypeBlock
		s.Block = this.StaticBlocks[k]
		f.Block.Statements[k] = s
	}
	f.makeLastReturnStatement()
	f.AccessFlags |= cg.AccMethodPublic
	f.AccessFlags |= cg.AccMethodStatic
	f.AccessFlags |= cg.AccMethodFinal
	f.AccessFlags |= cg.AccMethodBridge
	f.Name = classInitMethod
	f.Block.IsFunctionBlock = true
	f.Block.Fn = method.Function
	if this.Methods == nil {
		this.Methods = make(map[string][]*ClassMethod)
	}
	f.Block.inherit(&this.Block)
	f.Block.InheritedAttribute.Function = f
	f.Block.InheritedAttribute.ClassMethod = method
	this.Methods[f.Name] = []*ClassMethod{method}
}

func (this *Class) resolveFieldsAndMethodsType() []error {
	if this.resolveFieldsAndMethodsTypeCalled {
		return []error{}
	}
	this.resolveFieldsAndMethodsTypeCalled = true
	errs := []error{}
	var err error
	for _, v := range this.Fields {
		if v.Name == SUPER {
			errs = append(errs, fmt.Errorf("%s 'super' not allow for field name",
				errMsgPrefix(v.Pos)))
			continue
		}
		err = v.Type.resolve(&this.Block)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, ms := range this.Methods {
		for _, m := range ms {
			if m.IsAbstract() {
				for _, v := range m.Function.Type.ParameterList {
					if v.DefaultValueExpression != nil {
						errs = append(errs,
							fmt.Errorf("%s abstract method parameter '%s' cannot have default value '%s'",
								errMsgPrefix(v.Pos), v.Name, v.DefaultValueExpression.Op))
					}
				}
				for _, v := range m.Function.Type.ReturnList {
					if v.DefaultValueExpression != nil {
						errs = append(errs,
							fmt.Errorf("%s abstract method return variable '%s' cannot have default value '%s'",
								errMsgPrefix(v.Pos), v.Name, v.DefaultValueExpression.Op))
					}
				}
			}
			m.Function.Block.inherit(&this.Block)
			m.Function.Block.InheritedAttribute.Function = m.Function
			m.Function.Block.InheritedAttribute.ClassMethod = m
			m.Function.checkParametersAndReturns(&errs, false, m.IsAbstract())
			if len(m.Function.Type.TemplateNames) > 0 {
				errs = append(errs, fmt.Errorf("%s cannot use template for method",
					errMsgPrefix(m.Function.Pos)))
			}
			if m.IsStatic() == false { // bind this
				if m.Function.Block.Variables == nil {
					m.Function.Block.Variables = make(map[string]*Variable)
				}
				m.Function.Block.Variables[ThisPointerName] = &Variable{}
				m.Function.Block.Variables[ThisPointerName].Name = ThisPointerName
				m.Function.Block.Variables[ThisPointerName].Pos = m.Function.Pos
				m.Function.Block.Variables[ThisPointerName].Type = &Type{
					Type:  VariableTypeObject,
					Class: this,
				}
			}
		}
	}
	return errs
}

func (this *Class) resolveFather() error {
	if this.resolveFatherCalled {
		return nil
	}
	this.resolveFatherCalled = true
	if this.SuperClassName == nil {
		superClassName := ""
		if PackageBeenCompile.Name == common.CorePackage {
			superClassName = JavaRootClass
		} else {
			if this.IsInterface() {
				superClassName = JavaRootClass
			} else {
				superClassName = LucyRootClass
			}
		}
		this.SuperClass = &Class{
			Name:           superClassName,
			NotImportedYet: true,
		}
	} else {
		variableType := Type{}
		variableType.Type = VariableTypeName // naming
		variableType.Name = this.SuperClassName.Name
		variableType.Pos = this.SuperClassName.Pos
		err := variableType.resolve(&this.Block)
		if err != nil {
			return err
		}
		if variableType.Type != VariableTypeObject {
			err := fmt.Errorf("%s '%s' is not a class",
				this.Pos.ErrMsgPrefix(), this.SuperClassName.Name)
			return err
		}
		this.SuperClass = variableType.Class
		if this.IsInterface() {
			if this.SuperClass.Name != JavaRootClass {
				err := fmt.Errorf("%s interface`s super-class must be '%s'",
					errMsgPrefix(this.SuperClassName.Pos), JavaRootClass)

				return err
			}
		}
	}
	return nil
}

func (this *Class) resolveInterfaces() []error {
	if this.resolveInterfacesCalled {
		return nil
	}
	this.resolveInterfacesCalled = true
	errs := []error{}
	for _, i := range this.InterfaceNames {
		t := &Type{}
		t.Type = VariableTypeName
		t.Pos = i.Pos
		t.Name = i.Name
		err := t.resolve(&this.Block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if t.Type != VariableTypeObject {
			errs = append(errs, fmt.Errorf("%s '%s' is not a object , but '%s'",
				errMsgPrefix(i.Pos), i.Name, t.TypeString()))
			continue
		}
		if t.Class.IsInterface() == false {
			errs = append(errs, fmt.Errorf("%s '%s' is not a interface",
				errMsgPrefix(i.Pos), i.Name))
			continue
		}
		this.Interfaces = append(this.Interfaces, t.Class)
	}
	return errs
}

func (this *Class) implementMethod(
	pos *Pos,
	m *ClassMethod,
	nameMatched **ClassMethod,
	fromSub bool,
	errs *[]error) *ClassMethod {
	if this.Methods != nil {
		for _, v := range this.Methods[m.Function.Name] {
			if v.IsAbstract() {
				continue
			}
			if fromSub && v.ableAccessFromSubClass() == false {
				return nil
			}
			if *nameMatched == nil {
				*nameMatched = v
			}
			if v.Function.Type.equal(&m.Function.Type) {
				return v
			}
		}
	}
	//no same name method at current class
	if this.Name == JavaRootClass {
		return nil
	}
	err := this.loadSuperClass(pos)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}

	if this.SuperClass == nil {
		return nil
	}

	return this.SuperClass.implementMethod(pos, m, nameMatched, true, errs)
}

func (this *Class) haveSuperClass(pos *Pos, superclassName string) (bool, error) {
	err := this.loadSelf(pos)
	if err != nil {
		return false, err
	}
	if this.Name == superclassName {
		return true, nil
	}
	if this.Name == JavaRootClass {
		return false, nil
	}

	err = this.loadSuperClass(pos)
	if err != nil {
		return false, err
	}

	if this.SuperClass == nil {
		return false, nil
	}

	return this.SuperClass.haveSuperClass(pos, superclassName) // check father is implements
}

func (this *Class) implementedInterface(pos *Pos, inter string) (bool, error) {
	err := this.loadSelf(pos)
	if err != nil {
		return false, err
	}
	for _, v := range this.Interfaces {
		if v.Name == inter {
			return true, nil
		}
		im, _ := v.implementedInterface(pos, inter)
		if im {
			return im, nil
		}
	}
	if this.Name == JavaRootClass {
		return false, nil
	}
	err = this.loadSuperClass(pos)
	if err != nil {
		return false, err
	}
	if this.SuperClass == nil {
		return false, nil
	}
	return this.SuperClass.implementedInterface(pos, inter) // check father is implements
}

func (this *Class) loadSuperClass(pos *Pos) error {
	if this.SuperClass != nil {
		return this.SuperClass.loadSelf(pos)
	}
	if this.resolveFatherCalled ||
		this.loadSuperClassCalled {
		return nil
	}
	this.loadSuperClassCalled = true
	if this.Name == JavaRootClass {
		err := fmt.Errorf("%s root class already", errMsgPrefix(pos))
		return err
	}
	if this.SuperClassName == nil {
		this.SuperClassName = &NameWithPos{
			Name: JavaRootClass,
			Pos:  this.Pos,
		}
	}
	class, err := PackageBeenCompile.loadClass(this.SuperClassName.Name)
	if err != nil {
		err := fmt.Errorf("%s %v", pos.ErrMsgPrefix(), err)
		return err
	}
	this.SuperClass = class
	return nil
}

func (this *Class) accessConstructionMethod(
	pos *Pos,
	errs *[]error,
	newCase *ExpressionNew,
	callFatherCase *ExpressionMethodCall,
	callArgs []*Type) (ms []*ClassMethod, matched bool, err error) {
	err = this.loadSelf(pos)
	if err != nil {
		return nil, false, err
	}
	var args *CallArgs
	if newCase != nil {
		args = &newCase.Args
	} else {
		args = &callFatherCase.Args
	}
	for _, v := range this.Methods[SpecialMethodInit] {
		vArgs, err := v.Function.Type.fitArgs(pos, args, callArgs, v.Function)
		if err == nil {
			if newCase != nil {
				newCase.VArgs = vArgs
			} else {
				callFatherCase.VArgs = vArgs
			}
			return []*ClassMethod{v}, true, nil
		} else {
			if this.IsJava {
				ms = append(ms, v)
			} else {
				return nil, false, err
			}
		}
	}
	return ms, false, nil
}

/*
	ret is *ClassField or *ClassMethod
*/
func (this *Class) getFieldOrMethod(
	pos *Pos,
	name string,
	fromSub bool) (interface{}, error) {
	err := this.loadSelf(pos)
	if err != nil {
		return nil, err
	}
	notFoundErr := fmt.Errorf("%s field or method named '%s' not found", pos.ErrMsgPrefix(), name)
	if this.Fields != nil && nil != this.Fields[name] {
		if fromSub && this.Fields[name].ableAccessFromSubClass() == false {
			// private field
			// break find
			return nil, notFoundErr
		} else {
			return this.Fields[name], nil
		}
	}
	if this.Methods != nil && nil != this.Methods[name] {
		m := this.Methods[name][0]
		if fromSub && m.ableAccessFromSubClass() == false {
			return nil, notFoundErr
		} else {
			return m, nil
		}
	}
	if this.Name == JavaRootClass { // root class
		return nil, notFoundErr
	}
	err = this.loadSuperClass(pos)
	if err != nil {
		return nil, err
	}
	if this.SuperClass == nil {
		return nil, notFoundErr
	}
	return this.SuperClass.getFieldOrMethod(pos, name, true)
}

func (this *Class) constructionMethodAccessAble(pos *Pos, method *ClassMethod) error {
	if this.LoadFromOutSide {
		if this.IsPublic() == false {
			return fmt.Errorf("%s class '%s' is not public",
				pos.ErrMsgPrefix(), this.Name)
		}
		if method.IsPublic() == false {
			return fmt.Errorf("%s method '%s' is not public",
				pos.ErrMsgPrefix(), method.Function.Name)
		}
	} else {
		if method.IsPrivate() {
			return fmt.Errorf("%s method '%s' is private",
				pos.ErrMsgPrefix(), method.Function.Name)
		}
	}
	return nil
}
