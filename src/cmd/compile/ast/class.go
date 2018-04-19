package ast

import (
	"errors"
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
	SouceFile          string
	Used               bool
	LoadFromOutSide    bool
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
	//errs = append(errs, c.checkConstructionFunctions()...)
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
}

func (c *Class) checkFields() []error {
	errs := []error{}
	for _, v := range c.Fields {
		err := v.Typ.resolve(&c.Block)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
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
