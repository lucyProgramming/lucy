package ast

import (
	"errors"
	"fmt"
	"strings"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Class struct {
	FatherNameResolved bool
	NotImportedYet     bool // not imported
	Name               string
	Pos                *Position
	IsJava             bool //class found in CLASSPATH
	IsGlobal           bool
	Block              Block
	AccessFlags        uint16
	Fields             map[string]*ClassField
	Methods            map[string][]*ClassMethod
	SuperClassName     string
	SuperClass         *Class
	InterfaceNames     []*NameWithPos
	Interfaces         []*Class
	LoadFromOutSide    bool
	StaticBlocks       []*Block
}

func (c *Class) HaveStaticsCodes() bool {
	s := 0
	for _, v := range c.StaticBlocks {
		s += len(v.Statements)
	}
	return s > 0
}

func (c *Class) IsInterface() bool {
	return c.AccessFlags&cg.ACC_CLASS_INTERFACE != 0
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
	errs := c.checkPhase1(father)
	es := c.checkPhase2(father)
	if errorsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}

func (c *Class) mkDefaultConstruction() {
	if c.Methods == nil {
		c.Methods = make(map[string][]*ClassMethod)
	}
	if len(c.Methods[ConstructionMethodName]) > 0 {
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
	{
		e := &Expression{}
		e.Type = EXPRESSION_TYPE_METHOD_CALL
		e.Pos = c.Pos
		call := &ExpressionMethodCall{}
		call.Name = SUPER
		call.Expression = &Expression{
			Type: EXPRESSION_TYPE_IDENTIFIER,
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
	c.Methods[ConstructionMethodName] = []*ClassMethod{m}
}

func (c *Class) checkIfClassHierarchyCircularity() error {
	m := make(map[string]struct{})
	arr := []string{}
	is := false
	class := c
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
func (c *Class) checkPhase1(father *Block) []error {
	c.Block.inherit(father)
	errs := c.Block.checkConstants()
	c.Block.InheritedAttribute.Class = c
	es := c.resolveAllNames(father)
	if errorsNotEmpty(es) {
		errs = append(errs, es...)
	}
	err := c.resolveFather(father)
	if err != nil {
		errs = append(errs, err)
	} else {
		err = c.checkIfClassHierarchyCircularity()
		if err != nil {
			errs = append(errs, err)
		}
	}

	es = c.resolveInterfaces(father)
	errs = append(errs, es...)
	es = c.suitableForInterfaces()
	errs = append(errs, es...)
	return errs
}

func (c *Class) checkPhase2(father *Block) []error {
	errs := []error{}
	c.Block.InheritedAttribute.Class = c
	errs = append(errs, c.checkFields()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	c.mkClassInitMethod()
	for _, ms := range c.Methods {
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
	f.mkLastReturnStatement()
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

func (c *Class) resolveAllNames(b *Block) []error {
	errs := []error{}
	var err error
	for _, v := range c.Fields {
		if v.Name == SUPER {
			errs = append(errs, fmt.Errorf("%s super is special for access 'super'",
				errMsgPrefix(v.Pos)))
		}
		err = v.Type.resolve(b)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range c.Methods {
		for _, vv := range v {
			if vv.Function.AccessFlags&cg.ACC_METHOD_STATIC == 0 { // bind this
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
		}
	}
	return errs
}

func (c *Class) resolveFather(block *Block) error {
	if c.SuperClass != nil || c.FatherNameResolved {
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
		c.FatherNameResolved = true
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
		r, err := PackageBeenCompile.load(i.ImportName)
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
		err := variableType.resolve(block)
		if err != nil {
			return err
		}
		if variableType.Type != VariableTypeObject {
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		}
		c.SuperClassName = variableType.Class.Name
		c.SuperClass = variableType.Class
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
func (c *Class) resolveInterfaces(block *Block) []error {
	errs := []error{}
	for _, i := range c.InterfaceNames {
		t := &Type{}
		t.Type = VariableTypeName
		t.Pos = i.Pos
		t.Name = i.Name
		err := t.resolve(block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if t.Type != VariableTypeObject {
			errs = append(errs, fmt.Errorf("%s '%s' is not a class",
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
		if fromSub == false || m.IsPrivate() == false {
			continue
		}
		args := make([]*Type, len(m.Function.Type.ParameterList))
		for k, v := range m.Function.Type.ParameterList {
			args[k] = v.Type
		}
		_, match, _ := c.accessMethod(c.Pos, &errs, name, args, nil, false)
		if match == false {
			err := fmt.Errorf("%s class named '%s' does not implement '%s',missing method '%s'",
				errMsgPrefix(c.Pos), c.Name, inter.Name, m.Function.readableMsg())
			errs = append(errs, err)
		}
	}
	return errs
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
	}
	if c.Name == JavaRootClass {
		return false, nil
	}
	err = c.loadSuperClass()
	if err != nil {
		return false, err
	}
	return c.SuperClass.implemented(inter)
}

func (c *Class) checkFields() []error {
	errs := []error{}
	ss := []*Statement{}
	for _, v := range c.Fields {
		if v.Expression != nil {
			t, es := v.Expression.checkSingleValueContextExpression(&c.Block)
			if errorsNotEmpty(es) {
				errs = append(errs, es...)
			}
			if v.Type.Equal(&errs, t) == false {
				errs = append(errs, fmt.Errorf("%s cannot assign '%s' as '%s' for default value",
					errMsgPrefix(v.Pos), t.TypeString(), v.Type.TypeString()))
				continue
			}
			//TODO:: should check or not ???
			//if t.Type == VARIABLE_TYPE_NULL {
			//	errs = append(errs, fmt.Errorf("%s pointer types default value is '%s' already",
			//		errMsgPrefix(v.Pos), t.TypeString()))
			//	continue
			//}
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
				Type: EXPRESSION_TYPE_LIST,
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
					Type: EXPRESSION_TYPE_SELECTION,
					Data: selection,
				}
				left.ExpressionValue = v.Type
				bin.Left = &Expression{
					Type: EXPRESSION_TYPE_LIST,
					Data: []*Expression{left},
				}
			}
			e := &Expression{
				Type: EXPRESSION_TYPE_ASSIGN,
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
	for name, v := range c.Methods {
		for _, vv := range v {
			if c.IsInterface() {
				continue
			}
			isConstruction := (name == ConstructionMethodName)
			if isConstruction &&
				vv.IsFirstStatementCallFatherConstruction() == false {
				errs = append(errs, fmt.Errorf("%s construction method should call father construction method first",
					errMsgPrefix(vv.Function.Pos)))
			}
			if isConstruction && vv.Function.NoReturnValue() == false {
				errs = append(errs, fmt.Errorf("%s construction method expect no return values",
					errMsgPrefix(vv.Function.Type.ParameterList[0].Pos)))
			}
			if c.IsInterface() == false {
				vv.Function.Block.InheritedAttribute.IsConstructionMethod = isConstruction
				vv.Function.checkBlock(&errs)
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
