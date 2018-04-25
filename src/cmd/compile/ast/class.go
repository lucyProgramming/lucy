package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"strings"
)

type Class struct {
	NotImportedYet     bool // not imported
	Pos                *Pos
	IsJava             bool // compiled from java source file
	Name               string
	NameWithOutPackage string
	IsGlobal           bool
	Block              Block
	AccessFlags        uint16
	Fields             map[string]*ClassField
	Methods            map[string][]*ClassMethod
	SuperClassName     string
	SuperClass         *Class
	InterfaceNames     []*NameWithPos
	Interfaces         []*Class
	SouceFile          string
	Used               bool
	LoadFromOutSide    bool
}

func (c *Class) loadSelf() error {
	if c.NotImportedYet == false {
		return nil
	}
	cc, err := PackageBeenCompile.load(c.Name)
	if err != nil {
		return err
	}
	*c = *(cc.(*Class))
	return nil
}
func (c *Class) check(father *Block) []error {
	errs := c.checkPhase1(father)
	es := c.checkPhase2(father)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}

func (c *Class) checkPhase1(father *Block) []error {
	c.Block.inherite(father)
	c.Block.InheritedAttribute.class = c
	errs := c.resolveName(father)
	err := c.resolveFather(father)
	if err != nil {
		errs = append(errs, err)
	}
	es := c.resolveInterfaces(father)
	errs = append(errs, es...)
	es = c.suitableForInterfaces()
	errs = append(errs, es...)
	return errs
}

func (c *Class) checkPhase2(father *Block) []error {
	errs := []error{}
	c.Block.check() // check innerclass mainly
	c.Block.InheritedAttribute.class = c
	errs = append(errs, c.checkFields()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	for _, ms := range c.Methods {
		if len(ms) > 1 {
			errmsg := fmt.Sprintf("%s class method named '%s' has declared %d times,which are:\n",
				errMsgPrefix(ms[0].Func.Pos),
				ms[0].Func.Name, len(ms))
			for _, v := range ms {
				errmsg += fmt.Sprintf("\t%s\n", errMsgPrefix(v.Func.Pos))
			}
			errs = append(errs, errors.New(errmsg))
		}
	}
	errs = append(errs, c.checkMethods()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	return errs
}

func (c *Class) resolveName(b *Block) []error {
	errs := []error{}
	var err error
	for _, v := range c.Fields {
		if v.Name == SUPER_FIELD_NAME {
			errs = append(errs, fmt.Errorf("%s super is special for access 'super'",
				errMsgPrefix(v.Pos)))
		}
		err = v.Typ.resolve(b)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range c.Methods {
		for _, vv := range v {
			vv.Func.checkParaMeterAndRetuns(&errs)
		}
	}
	return errs
}

func (c *Class) resolveFather(block *Block) error {
	if c.SuperClass != nil {
		return nil
	}
	defer func() {
		if c.SuperClassName == "" {
			c.SuperClassName = LUCY_ROOT_CLASS
		}
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
		superClass, err := PackageBeenCompile.load(i.Resource + "/" + t[1])
		if err != nil {
			return fmt.Errorf("%s %v", errMsgPrefix(c.Pos), err)
		}
		if ss, ok := superClass.(*Class); ok == false || ss == nil {
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		} else {
			t := ss
			c.SuperClassName = t.Name
			c.SuperClass = t
		}
	} else {
		variableType := VariableType{}
		variableType.Typ = VARIABLE_TYPE_NAME // naming
		variableType.Name = c.SuperClassName
		variableType.Pos = c.Pos
		err := variableType.resolve(block)
		if err != nil {
			return err
		}
		if variableType.Typ != VARIABLE_TYPE_OBJECT {
			return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		}
		c.SuperClassName = variableType.Class.Name
		c.SuperClass = variableType.Class
	}
	return nil
}
func (c *Class) resolveInterfaces(block *Block) []error {
	errs := []error{}
	for _, i := range c.InterfaceNames {
		t := &VariableType{}
		t.Typ = VARIABLE_TYPE_NAME
		t.Pos = i.Pos
		t.Name = i.Name
		err := t.resolve(block)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if t.Typ != VARIABLE_TYPE_OBJECT {
			errs = append(errs, fmt.Errorf("%s '%s' is not a class", errMsgPrefix(i.Pos), i.Name))
			continue
		}
		c.Interfaces = append(c.Interfaces, t.Class)
	}
	return errs
}

func (c *Class) suitableForInterfaces() []error {
	errs := []error{}
	for _, i := range c.Interfaces {
		errs = append(errs, c.suitableForInterface(i)...)
	}
	return errs
}
func (c *Class) suitableForInterface(inter *Class) []error {
	errs := []error{}
	for name, v := range inter.Methods {
		m := v[0]
		args := make([]*VariableType, len(m.Func.Typ.ParameterList))
		for k, v := range m.Func.Typ.ParameterList {
			args[k] = v.Typ
		}
		_, match, _ := c.accessMethod(name, args, nil, false)
		if match == false {
			err := fmt.Errorf("%s class named '%s' does not implement '%s',missing method '%s'",
				errMsgPrefix(c.Pos), c.Name, inter.Name, m.Func.readableMsg())
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return errs
}

func (c *Class) IsInterface() bool {
	return c.AccessFlags&cg.ACC_CLASS_INTERFACE != 0
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
	return errs
}

func (c *Class) checkMethods() []error {
	errs := []error{}
	if c.IsInterface() {
		return errs
	}
	for name, v := range c.Methods {
		for _, vv := range v {
			if vv.Func.AccessFlags&cg.ACC_METHOD_STATIC == 0 { // bind this
				if vv.Func.Block.Vars == nil {
					vv.Func.Block.Vars = make(map[string]*VariableDefinition)
				}
				vv.Func.Block.Vars[THIS] = &VariableDefinition{}
				vv.Func.Block.Vars[THIS].Name = THIS
				vv.Func.Block.Vars[THIS].Pos = vv.Func.Pos
				vv.Func.Block.Vars[THIS].Typ = &VariableType{
					Typ:   VARIABLE_TYPE_OBJECT,
					Class: c,
				}
			}
			isConstruction := (name == CONSTRUCTION_METHOD_NAME)
			c.Block.InheritedAttribute.IsConstruction = isConstruction
			if isConstruction && vv.Func.NoReturnValue() == false {
				errs = append(errs, fmt.Errorf("%s construction method expect no return",
					errMsgPrefix(vv.Func.Typ.ParameterList[0].Pos)))
			}
			vv.Func.Block.inherite(&c.Block)
			vv.Func.Block.InheritedAttribute.Function = vv.Func
			es := vv.Func.Block.check()
			if errsNotEmpty(es) {
				errs = append(errs, es...)
			}
		}
	}
	return errs
}

func (c *Class) loadSuperClass() error {
	if c.SuperClass != nil {
		return nil
	}
	d, err := PackageBeenCompile.load(c.SuperClassName)
	if err != nil {
		return err
	}
	if class, ok := d.(*Class); ok && ok && class != nil {
		c.SuperClass = class
		return nil
	} else {
		return fmt.Errorf("'%s' is not a class", c.SuperClassName)
	}
}
