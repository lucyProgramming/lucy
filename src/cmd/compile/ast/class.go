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
	if c.Methods != nil && len(c.Methods[CONSTRUCTION_METHOD_NAME]) > 0 {
		return
	}
	if c.Methods == nil {
		c.Methods = make(map[string][]*ClassMethod)
	}
	m := &ClassMethod{}
	//m.IsConstructionMethod = true
	m.Func = &Function{}
	m.Func.AccessFlags |= cg.ACC_METHOD_PUBLIC
	m.Func.Pos = c.Pos
	m.Func.Block.IsFunctionBlock = true
	c.Methods[CONSTRUCTION_METHOD_NAME] = []*ClassMethod{m}
}

func (c *Class) checkIfClassHierarchyCircularity() error {
	m := make(map[string]struct{})
	arr := []string{}
	is := false
	class := c
	for class.Name != JAVA_ROOT_CLASS {
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
	for _, ms := range c.Methods {
		if len(ms) > 1 {
			errMsg := fmt.Sprintf("%s class method named '%s' has declared %d times,which are:\n",
				errMsgPrefix(ms[0].Func.Pos),
				ms[0].Func.Name, len(ms))
			for _, v := range ms {
				errMsg += fmt.Sprintf("\t%s\n", errMsgPrefix(v.Func.Pos))
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

func (c *Class) resolveAllNames(b *Block) []error {
	errs := []error{}
	var err error
	for _, v := range c.Fields {
		if v.Name == SUPER_FIELD_NAME {
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
			vv.Func.Block.inherit(&c.Block)
			vv.Func.Block.InheritedAttribute.Function = vv.Func
			vv.Func.checkParametersAndReturns(&errs)
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
				c.SuperClassName = LUCY_ROOT_CLASS
			} else {
				c.SuperClassName = JAVA_ROOT_CLASS
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
		variableType.Type = VARIABLE_TYPE_NAME // naming
		variableType.Name = c.SuperClassName
		variableType.Pos = c.Pos
		err := variableType.resolve(block)
		if err != nil {
			return err
		}
		if variableType.Type != VARIABLE_TYPE_OBJECT {
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		}
		c.SuperClassName = variableType.Class.Name
		c.SuperClass = variableType.Class
	}
	if c.IsInterface() {
		if c.SuperClass.Name == JAVA_ROOT_CLASS {
			//nothing
		} else {
			return fmt.Errorf("%s interface`s super-class must be '%s'",
				errMsgPrefix(c.Pos), JAVA_ROOT_CLASS)
		}
	}
	return nil
}
func (c *Class) resolveInterfaces(block *Block) []error {
	errs := []error{}
	for _, i := range c.InterfaceNames {
		t := &Type{}
		t.Type = VARIABLE_TYPE_NAME
		t.Pos = i.Pos
		t.Name = i.Name
		err := t.resolve(block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if t.Type != VARIABLE_TYPE_OBJECT {
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
		args := make([]*Type, len(m.Func.Type.ParameterList))
		for k, v := range m.Func.Type.ParameterList {
			args[k] = v.Type
		}
		_, match, _ := c.accessMethod(c.Pos, &errs, name, args, nil, false)
		if match == false {
			err := fmt.Errorf("%s class named '%s' does not implement '%s',missing method '%s'",
				errMsgPrefix(c.Pos), c.Name, inter.Name, m.Func.readableMsg())
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
	if c.Name == JAVA_ROOT_CLASS {
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
	if c.Name == JAVA_ROOT_CLASS {
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
	for _, v := range c.Fields {
		if v.Expression != nil {
			if v.Expression.IsLiteral() == false {
				errs = append(errs, fmt.Errorf("%s field default value must be literal",
					errMsgPrefix(v.Pos)))
				continue
			}
			ts, _ := v.Expression.check(&c.Block)
			if v.Type.Equal(&errs, ts[0]) == false {
				errs = append(errs, fmt.Errorf("%s cannot assign '%s' as '%s' for default value",
					errMsgPrefix(v.Pos), ts[0].TypeString(), v.Type.TypeString()))
				continue
			}
			v.DefaultValue = v.Expression.Data // copy default value
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
			if vv.Func.AccessFlags&cg.ACC_METHOD_STATIC == 0 { // bind this
				if vv.Func.Block.Variables == nil {
					vv.Func.Block.Variables = make(map[string]*Variable)
				}
				vv.Func.Block.Variables[THIS] = &Variable{}
				vv.Func.Block.Variables[THIS].Name = THIS
				vv.Func.Block.Variables[THIS].Pos = vv.Func.Pos
				vv.Func.Block.Variables[THIS].Type = &Type{
					Type:  VARIABLE_TYPE_OBJECT,
					Class: c,
				}
			}
			isConstruction := (name == CONSTRUCTION_METHOD_NAME)
			if isConstruction && vv.Func.NoReturnValue() == false {
				errs = append(errs, fmt.Errorf("%s construction method expect no return values",
					errMsgPrefix(vv.Func.Type.ParameterList[0].Pos)))
			}
			if c.IsInterface() == false {
				vv.Func.Block.InheritedAttribute.IsConstructionMethod = isConstruction
				vv.Func.checkBlock(&errs)
			}
		}
	}
	return errs
}

func (c *Class) loadSuperClass() error {
	if c.SuperClass != nil {
		return nil
	}
	if c.Name == JAVA_ROOT_CLASS {
		return fmt.Errorf("root class already")
	}
	class, err := PackageBeenCompile.loadClass(c.SuperClassName)
	if err != nil {
		return err
	}
	c.SuperClass = class
	return nil
}
