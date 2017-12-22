package ast

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
	"strings"
)

type Class struct {
	Block   Block
	Access  uint16
	Pos     *Pos
	Name    string
	Fields  map[string]*ClassField
	Methods map[string][]*ClassMethod
	//SuperClassExpression *Expression // a or a.b
	SuperClassName string
	SuperClass     *Class
	Interfaces     []*Class
	Constructors   []*ClassMethod // can be nil
	Signature      *class_json.ClassSignature
	SouceFile      string
	Used           bool
	VariableType   VariableType
	ClosureVars    map[string]*VariableDefinition // closure variable
}

func (c *Class) check(father *Block) []error {
	errs := make([]error, 0)
	c.Block.check(father)
	errs = append(errs, c.checkFields()...)
	errs = append(errs, c.checkConstructionFunctions()...)
	errs = append(errs, c.checkMethods()...)
	return errs
}

func (c *Class) isInterface() bool {
	return c.Access&cg.ACC_CLASS_INTERFACE != 0
}

func (c *Class) implementInterfaceOf(father *Class) bool {
	if father.Access&cg.ACC_CLASS_INTERFACE == 0 {
		panic("not a interface")
	}
	for _, v := range c.Interfaces {
		if v.Name == father.Name {
			return true
		}
	}
	return false
}

func (c *Class) instanceOf(father *Class) bool {
	if father.Access&cg.ACC_CLASS_INTERFACE != 0 {
		return c.implementInterfaceOf(father)
	}
	return false
}

func (c *Class) mkVariableType() {
	c.VariableType.Typ = VARIABLE_TYPE_CLASS
	c.VariableType.Class = c
}

func (c *Class) checkConstructionFunctions() []error {
	errs := []error{}
	c.checkReloadFunctions(c.Constructors, &errs)
	return errs
}

func (c *Class) checkReloadFunctions(ms []*ClassMethod, errs *[]error) {
	m := make(map[string][]*ClassMethod)
	for _, v := range ms {
		if v.Func.AccessFlags&cg.ACC_METHOD_STATIC == 0 { // bind this
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
		}
		es := v.Func.check(&c.Block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
		}
		v.Func.MkDescriptor()
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
	err := c.VariableType.resolve(&c.Block)
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

func (c *Class) accessField(name string) (f *ClassField, accessable bool, err error) {
	if c.Fields[name] == nil {
		err = fmt.Errorf("field %s not found", name)
		return
	}
	f = c.Fields[name]
	accessable = (f.AccessFlags & cg.ACC_FIELD_PUBLIC) != 0
	return
}
func (c *Class) accessMethod(name string, args []*VariableType) (f *ClassMethod, errs []error) {
	errs = []error{}
	ms, ok := c.Methods[name]
	if !ok {
		return
	}
	s := mkSignatureByVariableTypes(args)
	s = "(" + s + ")"
	for _, m := range ms {
		if s == m.Func.Descriptor {
			f = m
			return
		}
	}
	cc := c
	for cc.Name != JAVA_ADAMS_CLASS {
		if cc.SuperClass == nil { // super class is not loaded
			es := cc.loadSuperClass()
			if errsNotEmpty(es) {
				errs = append(errs, es...)
				break
			}
		}
	}
	// not found,trying to access father`s method

	return
}

func (c *Class) loadSuperClass() []error {
	errs := []error{}
	if c.SuperClassName == "" {
		c.SuperClassName = "java/lang/Object"
	}
	pname := c.SuperClassName[0:strings.LastIndex(c.SuperClassName, "/")]
	c.Block.loadPackage(pname)
	return errs
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

type ClassMethod struct {
	Func *Function
}

type ClassField struct {
	VariableDefinition
}
