package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"path/filepath"
	"strings"
)

type Class struct {
	resolveFatherCalled               bool
	resolveInterfacesCalled           bool
	resolveFieldsAndMethodsTypeCalled bool
	NotImportedYet                    bool   // not imported
	Name                              string // binary name
	Pos                               *Position
	IsJava                            bool //class imported from CLASSPATH
	IsGlobal                          bool
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
}

func (c *Class) HaveStaticsCodes() bool {
	return len(c.StaticBlocks) > 0
}

func (c *Class) IsInterface() bool {
	return c.AccessFlags&cg.ACC_CLASS_INTERFACE != 0
}
func (c *Class) IsFinal() bool {
	return c.AccessFlags&cg.ACC_CLASS_FINAL != 0
}
func (c *Class) IsPublic() bool {
	return c.AccessFlags&cg.ACC_CLASS_PUBLIC != 0
}

func (c *Class) loadSelf() error {
	if c.NotImportedYet == false {
		return nil
	}
	cc, err := PackageBeenCompile.loadClass(c.Name)
	if err != nil {
		return err
	}
	*c = *cc
	return nil
}

func (c *Class) check(father *Block) []error {
	c.Block.inherit(father)
	c.Block.InheritedAttribute.Class = c
	errs := c.checkPhase1()
	es := c.checkPhase2()
	if esNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}

func (c *Class) checkPhase1() []error {
	errs := c.Block.checkConstants()
	err := c.resolveFather()
	if err != nil {
		errs = append(errs, err)
	} else {
		err = c.checkIfClassHierarchyErr()
		if err != nil {
			errs = append(errs, err)
		}
	}
	es := c.resolveFieldsAndMethodsType()
	if esNotEmpty(es) {
		errs = append(errs, es...)
	}
	es = c.resolveInterfaces()
	errs = append(errs, es...)
	es = c.suitableForInterfaces()
	errs = append(errs, es...)
	return errs
}

func (c *Class) checkPhase2() []error {
	errs := []error{}
	if c.Block.InheritedAttribute.ClassAndFunctionNames == "" {
		c.Block.InheritedAttribute.ClassAndFunctionNames = filepath.Base(c.Name)
	} else {
		c.Block.InheritedAttribute.ClassAndFunctionNames += "$" + filepath.Base(c.Name)
	}
	errs = append(errs, c.checkFields()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	c.mkClassInitMethod()
	for name, ms := range c.Methods {
		if c.Fields != nil && c.Fields[name] != nil {
			f := c.Fields[name]
			errMsg := fmt.Sprintf("%s class method named '%s' already declared as field,at:\n",
				errMsgPrefix(ms[0].Function.Pos),
			)
			errMsg += fmt.Sprintf("\t%s", errMsgPrefix(f.Pos))
			errs = append(errs, errors.New(errMsg))
			continue
		}
		if len(ms) > 1 {
			errMsg := fmt.Sprintf("%s class method named '%s' has declared %d times,which are:\n",
				errMsgPrefix(ms[0].Function.Pos),
				ms[0].Function.Name, len(ms))
			for _, v := range ms {
				errMsg += fmt.Sprintf("\t%s\n", errMsgPrefix(v.Function.Pos))
			}
			errs = append(errs, errors.New(errMsg))
		}
	}
	errs = append(errs, c.checkMethods()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	errs = append(errs, c.checkIfOverrideFinalMethod()...)
	return errs
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
	if c.Methods == nil {
		c.Methods = make(map[string][]*ClassMethod)
	}
	m := &ClassMethod{}
	m.isCompilerAuto = true
	m.Function = &Function{}
	m.Function.AccessFlags |= cg.ACC_METHOD_PUBLIC
	m.Function.Pos = c.Pos
	m.Function.Block.IsFunctionBlock = true
	m.Function.Name = SpecialMethodInit
	{
		e := &Expression{}
		e.Type = ExpressionTypeMethodCall
		e.Pos = c.Pos
		call := &ExpressionMethodCall{}
		call.Name = SUPER
		call.Expression = &Expression{
			Type: ExpressionTypeIdentifier,
			Data: &ExpressionIdentifier{
				Name: THIS,
			},
			Pos: c.Pos,
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

func (c *Class) checkIfClassHierarchyErr() error {
	m := make(map[string]struct{})
	arr := []string{}
	is := false
	class := c
	if err := c.loadSuperClass(); err != nil {
		return err
	}
	if c.SuperClass.IsFinal() {
		return fmt.Errorf("class name '%s' have super class  named '%s' that is final", c.Name, c.SuperClassName)
	}
	for class.Name != JavaRootClass {
		_, ok := m[class.Name]
		if ok {
			arr = append(arr, class.Name)
			is = true
			break
		}
		m[class.Name] = struct{}{}
		arr = append(arr, class.Name)
		err := class.loadSuperClass()
		if err != nil {
			return err
		}
		if class.SuperClass == nil {
			panic("class is nil")
		}
		class = class.SuperClass
	}
	if is == false {
		return nil
	}
	errMsg := fmt.Sprintf("%s class named '%s' detects a circularity in class hierarchy\n",
		errMsgPrefix(c.Pos), c.Name)
	tab := "\t"
	index := len(arr) - 1
	for index >= 0 {
		errMsg += tab + arr[index] + "\n"
		tab += " "
		index--
	}
	return fmt.Errorf(errMsg)
}

func (c *Class) checkIfOverrideFinalMethod() []error {
	err := c.loadSuperClass()
	if err != nil {
		return []error{err}
	}
	errs := []error{}
	for name, v := range c.Methods {
		if name == SpecialMethodInit {
			continue
		}
		if len(v) == 0 {
			continue
		}
		if len(c.SuperClass.Methods[name]) == 0 {
			// this class not found at super
			continue
		}
		m := v[0]
		for _, v := range c.SuperClass.Methods[name] {
			f1 := &Type{
				Type:         VariableTypeFunction,
				FunctionType: &m.Function.Type,
			}
			f2 := &Type{
				Type:         VariableTypeFunction,
				FunctionType: &v.Function.Type,
			}
			if f1.Equal(&errs, f2) == true {
				if v.IsFinal() {
					errs = append(errs, fmt.Errorf("%s override final method",
						errMsgPrefix(m.Function.Pos)))
				}
			}
		}
	}
	return errs
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
	f.Name = ClassInitMethod
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
			errs = append(errs, fmt.Errorf("%s super is special for access 'super'",
				errMsgPrefix(v.Pos)))
		}
		err = v.Type.resolve(&c.Block)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range c.Methods {
		for _, vv := range v {
			if vv.IsStatic() == false { // bind this
				if vv.Function.Block.Variables == nil {
					vv.Function.Block.Variables = make(map[string]*Variable)
				}
				vv.Function.Block.Variables[THIS] = &Variable{}
				vv.Function.Block.Variables[THIS].Name = THIS
				vv.Function.Block.Variables[THIS].Pos = vv.Function.Pos
				vv.Function.Block.Variables[THIS].Type = &Type{
					Type:  VariableTypeObject,
					Class: c,
				}
			}
			vv.Function.Block.inherit(&c.Block)
			vv.Function.Block.InheritedAttribute.Function = vv.Function
			vv.Function.Block.InheritedAttribute.ClassMethod = vv
			vv.Function.checkParametersAndReturns(&errs)
			if c.IsInterface() {
				for _, vvv := range vv.Function.Type.ParameterList {
					if vvv.Expression != nil {
						errs = append(errs, fmt.Errorf("%s interface method parameter '%s' cannot have default value",
							errMsgPrefix(vvv.Pos), vvv.Name))
					}
				}
				for _, vvv := range vv.Function.Type.ReturnList {
					if vvv.Expression != nil {
						errs = append(errs, fmt.Errorf("%s interface method return variable '%s' cannot have default value",
							errMsgPrefix(vvv.Pos), vvv.Name))
					}
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
			return fmt.Errorf("%s package name '%s' not imported", errMsgPrefix(c.Pos), t[0])
		}
		r, err := PackageBeenCompile.load(i.Import)
		if err != nil {
			return fmt.Errorf("%s %v", errMsgPrefix(c.Pos), err)
		}
		if p, ok := r.(*Package); ok && p != nil { // if package
			if _, ok := p.Block.NameExists(t[1]); ok == false {
				return fmt.Errorf("%s class not exists in package '%s' ", errMsgPrefix(c.Pos), t[1])
			}
			if p.Block.Classes == nil || p.Block.Classes[t[1]] == nil {
				return fmt.Errorf("%s class not exists in package '%s' ", errMsgPrefix(c.Pos), t[1])
			}
			c.SuperClass = p.Block.Classes[t[1]]
		} else if ss, ok := r.(*Class); ok && ss != nil { // must be class now
			t := ss
			c.SuperClassName = t.Name
			c.SuperClass = t
		} else {
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		}
	} else {
		variableType := Type{}
		variableType.Type = VariableTypeName // naming
		variableType.Name = c.SuperClassName
		variableType.Pos = c.Pos
		err := variableType.resolve(&c.Block, false)
		if err != nil {
			return err
		}
		if variableType.Type != VariableTypeObject {
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		}
		c.SuperClassName = variableType.Class.Name
		c.SuperClass = variableType.Class
	}
	err := c.loadSuperClass()
	if err != nil {
		return err
	}
	if c.IsInterface() {
		if c.SuperClass.Name == JavaRootClass {
			//nothing
		} else {
			return fmt.Errorf("%s interface`s super-class must be '%s'",
				errMsgPrefix(c.Pos), JavaRootClass)
		}
	}
	return nil
}
func (c *Class) resolveInterfaces() []error {
	if c.resolveInterfacesCalled {
		return []error{}
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
		if t.Type != VariableTypeObject {
			errs = append(errs, fmt.Errorf("%s '%s' is not a interface",
				errMsgPrefix(i.Pos), i.Name))
			continue
		}
		if t.Class.IsInterface() == false {
			errs = append(errs, fmt.Errorf("%s '%s' is not a interface",
				errMsgPrefix(i.Pos), i.Name))
			continue
		}
		c.Interfaces = append(c.Interfaces, t.Class)
	}
	return errs
}

func (c *Class) suitableForInterfaces() []error {
	errs := []error{}
	if c.IsInterface() {
		return errs
	}
	// c is class
	for _, i := range c.Interfaces {
		errs = append(errs, c.suitableForInterface(i, false)...)
	}
	return errs
}

func (c *Class) suitableForInterface(inter *Class, fromSub bool) []error {
	errs := []error{}
	for name, v := range inter.Methods {
		m := v[0]
		im, match := c.implementMethod(m, false, &errs, c.Pos)
		if match {
			if im.IsPublic() == false {
				err := fmt.Errorf("%s method '%s' is not public",
					errMsgPrefix(c.Pos), name)
				errs = append(errs, err)
			}
			if im.IsStatic() {
				err := fmt.Errorf("%s method '%s' is static",
					errMsgPrefix(c.Pos), name)
				errs = append(errs, err)
			}
			return errs
		} else {
			errs = append(errs, fmt.Errorf("%s missing implements method '%s' define on interface '%s'",
				errMsgPrefix(c.Pos), m.Function.readableMsg(), inter.Name))
		}
	}
	for _, vv := range inter.Interfaces {
		err := vv.loadSelf()
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		es := c.suitableForInterface(vv, true)
		if esNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	return errs
}

func (c *Class) implementMethod(m *ClassMethod, fromSub bool, errs *[]error, pos *Position) (*ClassMethod, bool) {
	if c.Methods == nil || len(c.Methods[m.Function.Name]) == 0 {
		//no same name method at current class
		if c.Name == JavaRootClass {
			return nil, false
		}
		err := c.loadSuperClass()
		if err != nil {
			*errs = append(*errs,
				fmt.Errorf("%s %v", errMsgPrefix(pos), err))
			return nil, false
		} else {
			return c.SuperClass.implementMethod(m, true, errs, pos)
		}
	}
	for _, v := range c.Methods[m.Function.Name] {
		if fromSub && v.IsPrivate() {
			return nil, false
		}
		if len(v.Function.Type.ParameterList) != len(m.Function.Type.ParameterList) {
			// parameter count not match
			continue
		}
		if len(v.Function.Type.ReturnList) != len(m.Function.Type.ReturnList) {
			// return list count not match
			continue
		}
		match := true
		for kk, p := range v.Function.Type.ParameterList {
			if p.Type.StrictEqual(m.Function.Type.ParameterList[kk].Type) == false {
				match = false
				break
			}
		}
		if match {
			for kk, p := range v.Function.Type.ReturnList {
				if p.Type.StrictEqual(m.Function.Type.ReturnList[kk].Type) == false {
					match = false
					break
				}
			}
		}
		if match {
			return v, true
		}
	}
	return nil, false
}

func (c *Class) haveSuper(superclassName string) (bool, error) {
	err := c.loadSelf()
	if err != nil {
		return false, err
	}
	if c.Name == superclassName {
		return true, nil
	}
	if c.Name == JavaRootClass {
		return false, nil
	}
	err = c.loadSuperClass()
	if err != nil {
		return false, err
	}
	return c.SuperClass.haveSuper(superclassName)
}

func (c *Class) implemented(inter string) (bool, error) {
	err := c.loadSelf()
	if err != nil {
		return false, err
	}
	for _, v := range c.Interfaces {
		if v.Name == inter {
			return true, nil
		}
		im, _ := v.implemented(inter)
		if im {
			return im, nil
		}
	}
	if c.Name == JavaRootClass {
		return false, nil
	}
	err = c.loadSuperClass()
	if err != nil {
		return false, err
	}
	return false, nil
}

func (c *Class) checkFields() []error {
	errs := []error{}
	ss := []*Statement{}
	for _, v := range c.Fields {
		if v.Expression != nil {
			t, es := v.Expression.checkSingleValueContextExpression(&c.Block)
			if esNotEmpty(es) {
				errs = append(errs, es...)
			}
			if v.Type.Equal(&errs, t) == false {
				errs = append(errs, fmt.Errorf("%s cannot assign '%s' as '%s' for default value",
					errMsgPrefix(v.Pos), t.TypeString(), v.Type.TypeString()))
				continue
			}
			//TODO:: should check or not ???
			if t.Type == VariableTypeNull {
				errs = append(errs, fmt.Errorf("%s pointer types default value is '%s' already",
					errMsgPrefix(v.Pos), t.TypeString()))
				continue
			}
			if v.IsStatic() && v.Expression.IsLiteral() {
				v.DefaultValue = v.Expression.Data
				continue
			}
			if v.IsStatic() == false {
				// nothing to do
				continue
			}
			bin := &ExpressionBinary{}
			bin.Right = &Expression{
				Type: ExpressionTypeList,
				Data: []*Expression{v.Expression},
			}
			{
				selection := &ExpressionSelection{}
				selection.Expression = &Expression{}
				selection.Expression.ExpressionValue = &Type{
					Type:  VariableTypeClass,
					Class: c,
				}
				selection.Name = v.Name
				selection.Field = v
				left := &Expression{
					Type: ExpressionTypeSelection,
					Data: selection,
				}
				left.ExpressionValue = v.Type
				bin.Left = &Expression{
					Type: ExpressionTypeList,
					Data: []*Expression{left},
				}
			}
			e := &Expression{
				Type: ExpressionTypeAssign,
				Data: bin,
				IsStatementExpression: true,
			}
			ss = append(ss, &Statement{
				Type:                      StatementTypeExpression,
				Expression:                e,
				isStaticFieldDefaultValue: true,
			})
		}
	}
	if len(ss) > 0 {
		b := &Block{}
		b.Statements = ss
		if c.StaticBlocks != nil {
			c.StaticBlocks = append([]*Block{b}, c.StaticBlocks...)
		} else {
			c.StaticBlocks = []*Block{b}
		}
	}
	return errs
}

func (c *Class) checkMethods() []error {
	errs := []error{}
	if c.IsInterface() {
		return errs
	}
	for name, methods := range c.Methods {
		for _, method := range methods {
			if c.IsInterface() {
				continue
			}
			isConstruction := (name == SpecialMethodInit)
			if isConstruction {
				if method.IsFirstStatementCallFatherConstruction() == false {
					errs = append(errs, fmt.Errorf("%s construction method should call father construction method first",
						errMsgPrefix(method.Function.Pos)))
				}
				if method.IsFinal() {
					errs = append(errs, fmt.Errorf("%s construction method cannot be final",
						errMsgPrefix(method.Function.Pos)))
				}
			}
			if isConstruction && method.Function.NoReturnValue() == false {
				errs = append(errs, fmt.Errorf("%s construction method expect no return values",
					errMsgPrefix(method.Function.Type.ParameterList[0].Pos)))
			}
			if c.IsInterface() == false {
				method.Function.Block.InheritedAttribute.IsConstructionMethod = isConstruction
				method.Function.checkBlock(&errs)
			}
		}
	}
	return errs
}

func (c *Class) loadSuperClass() error {
	if c.SuperClass != nil {
		return nil
	}
	if c.Name == JavaRootClass {
		return fmt.Errorf("root class already")
	}
	class, err := PackageBeenCompile.loadClass(c.SuperClassName)
	if err != nil {
		return err
	}
	c.SuperClass = class
	return nil
}
