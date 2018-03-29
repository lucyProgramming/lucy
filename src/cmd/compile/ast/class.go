package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"strings"
)

type Class struct {
	Package            *Package
	Checked            bool
	Name               string
	NameWithOutPackage string
	IsGlobal           bool
	Block              Block
	Access             uint16
	Pos                *Pos
	Fields             map[string]*ClassField
	Methods            map[string][]*ClassMethod
	SuperClassName     string
	SuperClass         *Class
	Interfaces         []*Class
	Constructors       []*ClassMethod // can be nil
	SouceFile          string
	Used               bool
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
		v.Func.checkParaMeterAndRetuns(&errs)
	}
	for _, v := range c.Methods {
		for _, vv := range v {
			vv.Func.checkParaMeterAndRetuns(&errs)
		}
	}
	return errs
}

func (c *Class) check(father *Block) []error {
	if c.Checked {
		return nil
	}
	c.Block.inherite(father)
	c.Checked = true
	//super class name
	if c.SuperClassName == "" {
		c.SuperClassName = LUCY_ROOT_CLASS
	} else {
		if strings.Contains(c.SuperClassName, ".") {

		} else {
			t := father.SearchByName(c.SuperClassName)
			if t == nil {
				c.SuperClassName = LUCY_ROOT_CLASS

			} else {
				if _, ok := t.(*Class); ok == false {
					c.SuperClassName = LUCY_ROOT_CLASS
				}
			}
		}
	}
	c.Name = PackageBeenCompile.Name + "/" + c.Name // binary name
	c.Package = PackageBeenCompile
	errs := make([]error, 0)
	c.Block.check() // check innerclass mainly
	c.Block.InheritedAttribute.class = c
	errs = append(errs, c.checkFields()...)
	if father.shouldStop(errs) {
		return errs
	}
	errs = append(errs, c.checkConstructionFunctions()...)
	if father.shouldStop(errs) {
		return errs
	}
	errs = append(errs, c.checkMethods()...)
	if father.shouldStop(errs) {
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
	if father.shouldStop(errs) {
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
	if father.shouldStop(errs) {
		return errs
	}
	return errs
}

func (c *Class) loadInterfaces(errs *[]error) {

}

func (c *Class) isInterface() bool {
	return c.Access&cg.ACC_CLASS_INTERFACE != 0
}

func (c *Class) implementInterfaceOf(super *Class) bool {
	//	if super.Access&cg.ACC_CLASS_INTERFACE == 0 {
	//		panic("not a interface")
	//	}
	//	for _, v := range c.Interfaces {
	//		if v.Name == super.Name {
	//			return true
	//		}
	//	}
	return false
}

func (c *Class) instanceOf(super *Class) bool {
	if super.Access&cg.ACC_CLASS_INTERFACE != 0 {
		return c.implementInterfaceOf(super)
	}
	return false
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
			err := fmt.Errorf("%s %s redeclared", errMsgPrefix(vv.Func.Pos))
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
	if c.SuperClassName == "" {
		//class, err := c.Block.InheritedAttribute.p.load("lucy/lang", "Object")
		//if err != nil {
		//	return fmt.Errorf("%s load super failed err:%v", errMsgPrefix(c.Pos), err)
		//}
		//if _,ok := class.(*Class); ok == false {
		//	return fmt.Errorf("%s load super failed err:%v", errMsgPrefix(c.Pos), err)
		//}
		//c.SuperClass =
		//if c.SuperClass == nil {
		//	panic("........")
		//}
	} else {
		if false == strings.Contains(c.SuperClassName, "/") {
			d := c.Block.SearchByName(c.SuperClassName)
			if c, ok := d.(*Class); ok {
				c.SuperClass = c
			} else {
				return fmt.Errorf("%s '%s' is not a class", errMsgPrefix(c.Pos), c.SuperClassName)
			}
		}
		t := strings.Split(c.SuperClassName, "/")
		f, ok := PackageBeenCompile.Files[t[0]]
		if ok == false {
			return fmt.Errorf("%s package named '%s' not imported", errMsgPrefix(c.Pos), t[0])
		}
		pname, ok := f.Imports[t[0]]
		if ok == false {
			return fmt.Errorf("%s package named '%s' not imported", errMsgPrefix(c.Pos), t[0])
		}
		class, err := PackageBeenCompile.load(pname.Name, "Object")
		if err != nil {
			return fmt.Errorf("%s load super failed err:%v", errMsgPrefix(c.Pos), err)
		}
		if _, ok := class.(*Class); ok == false {
			return fmt.Errorf("%s %s is not a class", errMsgPrefix(c.Pos), err)
		}
		c.SuperClass = class.(*Class)
		return nil
	}
	return nil
}

func (c *Class) matchContructionFunction(args []*VariableType) (f *ClassMethod, accessable bool, err error) {
	if len(c.Constructors) == 0 && len(args) == 0 {
		return nil, true, nil
	}
	f, err = c.reloadMethod(c.Constructors, args)
	if (f.Func.AccessFlags & cg.ACC_METHOD_PUBLIC) != 0 {
		accessable = true
	}
	return
}

func (c *Class) reloadMethod(ms []*ClassMethod, args []*VariableType) (f *ClassMethod, err error) {
	return
}
