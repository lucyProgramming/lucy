package ast

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
)

type Class struct {
	IsGlobal bool
	Block    Block
	Access   uint16
	Pos      *Pos
	Name     string
	Fields   map[string]*ClassField
	Methods  map[string][]*ClassMethod
	//SuperClassExpression *Expression // a or a.b
	SuperClassName string
	SuperClass     *Class
	Interfaces     []*Class
	Constructors   []*ClassMethod // can be nil
	Signature      *class_json.ClassSignature
	SouceFile      string
	Used           bool
	VariableType   VariableType
}

func (c *Class) check(father *Block) []error {
	errs := make([]error, 0)
	err := c.loadSuperClass()
	if err != nil {
		errs = append(errs, err)
	}
	c.loadInterfaces(&errs)
	c.Block.inherite((father))
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
		errs = append(errs, fmt.Errorf("class named '%s' has %d contructor,declare at:", c.Name, len(c.Constructors)))
		for _, v := range c.Constructors {
			errs = append(errs, fmt.Errorf("%s contructor method", errMsgPrefix(v.Func.Pos)))
		}
	}
	if father.shouldStop(errs) {
		return errs
	}
	for _, ms := range c.Methods {
		if len(ms) > 1 {
			errs = append(errs, fmt.Errorf("class named '%s' has %d contructor,declare at:", c.Name, len(c.Constructors)))
			for _, v := range ms {
				errs = append(errs, fmt.Errorf("%s contructor method", errMsgPrefix(v.Func.Pos)))
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
	if super.Access&cg.ACC_CLASS_INTERFACE == 0 {
		panic("not a interface")
	}
	for _, v := range c.Interfaces {
		if v.Name == super.Name {
			return true
		}
	}
	return false
}

func (c *Class) instanceOf(super *Class) bool {
	if super.Access&cg.ACC_CLASS_INTERFACE != 0 {
		return c.implementInterfaceOf(super)
	}
	return false
}

func (c *Class) mkVariableType() {
	c.VariableType.Typ = VARIABLE_TYPE_CLASS
	c.VariableType.Class = c
}

func (c *Class) checkConstructionFunctions() []error {
	errs := []error{}
	if c.Constructors == nil || len(c.Constructors) == 0 {
		c.Constructors = []*ClassMethod{c.mkDefaultConstructionMethod()}
	}
	c.checkReloadFunctions(c.Constructors, &errs)
	return errs
}
func (c *Class) mkDefaultConstructionMethod() *ClassMethod {
	ret := &ClassMethod{}
	ret.Func = &Function{}
	ret.Func.AccessFlags = cg.ACC_METHOD_PUBLIC
	ret.Func.AccessFlags = cg.ACC_METHOD_FINAL
	ret.Func.Typ = &FunctionType{}
	ret.Func.Descriptor = "()V"
	ret.Func.Block = &Block{}
	ret.Func.Block.Statements = append(ret.Func.Block.Statements, &Statement{
		Typ:             STATEMENT_TYPE_RETURN,
		StatementReturn: &StatementReturn{},
	})
	return ret
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
			v.Func.Varoffset = 1 //this function
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

func (c *Class) accessField(name string) (f *ClassField, err error) {
	if c.Fields[name] == nil {
		err = fmt.Errorf("field %s not found", name)
		return
	}
	f = c.Fields[name]
	return
}
func (c *Class) accessMethod(name string, pos *Pos, args []*VariableType) (f *ClassMethod, errs []error) {
	//	errs = []error{}
	//	cc := c
	//	matchString := mkSignatureByVariableTypes(args)
	//	matchString = "(" + matchString + ")"
	//	for cc.SuperClassName != "" { // java/lang/Object has no super class
	//		ms, ok := c.Methods[name]
	//		if ok {
	//			for _, m := range ms {
	//				if matchString == m.Func.Descriptor {
	//					f = m
	//					return
	//				}
	//			}
	//		}
	//		if cc.SuperClass == nil { // super class is not loaded
	//			err := cc.loadSuperClass()
	//			if err != nil {
	//				errs = append(errs, fmt.Errorf("%s %s", errMsgPrefix(pos), err.Error()))
	//				break
	//			} else {
	//				cc = cc.SuperClass
	//			}
	//		}
	//	}
	//	// not found,trying to access father`s method
	return
}

func (c *Class) loadSuperClass() error {
	if c.Name == LUCY_ROOT_CLASS { // root class
		c.SuperClassName = ""
		c.SuperClass = nil
		return nil
	}
	if c.SuperClassName == "" {
		c.SuperClassName = LUCY_ROOT_CLASS
	}
	t := &VariableType{Typ: VARIABLE_TYPE_NAME, Name: c.SuperClassName}
	err := t.resolve(&c.Block)
	if err != nil {
		return err
	}
	if t.Typ != VARIABLE_TYPE_CLASS {
		return fmt.Errorf("superclass is not class,but %s", t.TypeString())
	}
	if t.Class.isInterface() {
		return fmt.Errorf("superclass is interface", t.TypeString())
	}
	c.SuperClass = t.Class
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
