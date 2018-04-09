package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"strings"
)

type Class struct {
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
	InterfacesName     []*NameWithPos
	Interfaces         []*Class
	Constructors       []*ClassMethod // can be nil
	SouceFile          string
	Used               bool
	LoadFromOutSide    bool
}

func (c *Class) resolveName(b *Block) []error {
	errs := []error{}
	var err error
	for _, v := range c.Fields {
		err = v.Typ.resolve(b)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range c.Constructors {
		for _, vv := range v.Func.Typ.ParameterList {
			err := vv.Typ.resolve(b)
			if err != nil {
				errs = append(errs, err)
			}
		}
		for _, vv := range v.Func.Typ.ReturnList {
			err := vv.Typ.resolve(b)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, v := range c.Methods {
		for _, vv := range v {
			for _, vvv := range vv.Func.Typ.ParameterList {
				err := vvv.Typ.resolve(b)
				if err != nil {
					errs = append(errs, err)
				}
			}
			for _, vvv := range vv.Func.Typ.ReturnList {
				err := vvv.Typ.resolve(b)
				if err != nil {
					errs = append(errs, err)
				}
			}
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
		if _, ok := superClass.(*Class); ok {
			return fmt.Errorf("%s   '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
		} else {
			t := superClass.(*Class)
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

func (c *Class) checkPhase1(father *Block) []error {
	c.Block.inherite(father)
	c.Block.InheritedAttribute.class = c
	errs := c.resolveName(father)
	err := c.resolveFather(father)
	if err != nil {
		errs = append(errs, err)
	}
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
	errs = append(errs, c.checkConstructionFunctions()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	errs = append(errs, c.checkMethods()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	if len(c.Constructors) > 1 {
		errs = append(errs, fmt.Errorf("%s class named '%s' has %d(more than 1) contructor,declare at:",
			errMsgPrefix(c.Pos),
			c.Name, len(c.Constructors)))
		for _, v := range c.Constructors {
			errs = append(errs, fmt.Errorf("\t %s contructor method...", errMsgPrefix(v.Func.Pos)))
		}
	}
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	for _, ms := range c.Methods {
		if len(ms) > 1 {
			errs = append(errs, fmt.Errorf("%s class named '%s' has %d contructor,declare at:",
				errMsgPrefix(ms[0].Func.Pos),
				c.Name, len(ms)))
			for _, v := range ms {
				errs = append(errs, fmt.Errorf("\t%s contructor method", errMsgPrefix(v.Func.Pos)))
			}
		}
	}
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	return errs
}
func (c *Class) check(father *Block) []error {
	errs := c.checkPhase1(father)
	es := c.checkPhase2(father)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}

func (c *Class) isInterface() bool {
	return c.AccessFlags&cg.ACC_CLASS_INTERFACE != 0
}

func (c *Class) haveSuper(superclassName string) (bool, error) {
	if c.Name == superclassName {
		return true, nil
	}
	err := c.loadSuperClass()
	if err != nil {
		return false, err
	}
	return c.SuperClass.haveSuper(superclassName)
}

func (c *Class) implemented(inter string) (bool, error) {
	for _, v := range c.InterfacesName {
		if v.Name == inter {
			return true, nil
		}
	}
	err := c.loadSuperClass()
	if err != nil {
		return false, err
	}
	return c.SuperClass.implemented(inter)
}

func (c *Class) checkConstructionFunctions() []error {
	errs := []error{}
	for _, v := range c.Constructors {
		v.IsConstructionMethod = true
		if v.IsStatic() {
			errs = append(errs, fmt.Errorf("%s construction method must not be static", errMsgPrefix(v.Func.Pos)))
		}
	}
	c.checkReloadFunctions(c.Constructors, &errs)
	return errs
}

func (c *Class) checkReloadFunctions(ms []*ClassMethod, errs *[]error) {
	m := make(map[string][]*ClassMethod)
	for _, v := range ms {
		if v.Func.AccessFlags&cg.ACC_METHOD_STATIC == 0 || v.IsConstructionMethod { // bind this
			if v.Func.Block.Vars == nil {
				v.Func.Block.Vars = make(map[string]*VariableDefinition)
			}
			v.Func.Block.Vars[THIS] = &VariableDefinition{}
			v.Func.Block.Vars[THIS].Name = THIS
			v.Func.Block.Vars[THIS].Pos = v.Func.Pos
			v.Func.Block.Vars[THIS].Typ = &VariableType{
				Typ:   VARIABLE_TYPE_OBJECT,
				Class: c,
			}
			v.Func.VarOffset = 1 //this function
		}
		es := v.Func.check(&c.Block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		if m[v.Func.Descriptor] == nil {
			m[v.Func.Descriptor] = []*ClassMethod{v}
		} else {
			m[v.Func.Descriptor] = append(m[v.Func.Descriptor], v)
		}
	}
	for _, v := range m {
		if len(v) == 1 {
			continue
		}
		for _, vv := range v {
			err := fmt.Errorf("%s %s redeclared", errMsgPrefix(vv.Func.Pos), vv.Func.Name)
			*errs = append(*errs, err)
		}
	}
}

func (c *Class) checkFields() []error {
	errs := []error{}
	for _, v := range c.Fields {
		c.checkField(v, &errs)
	}
	return errs
}

func (c *Class) checkField(f *ClassField, errs *[]error) {
	err := f.Typ.resolve(&c.Block)
	if err != nil {
		*errs = append(*errs, err)
	}
}

func (c *Class) checkMethods() []error {
	errs := []error{}
	for _, v := range c.Methods {
		c.checkReloadFunctions(v, &errs)
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
