package ast

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
)

type Class struct {
	Block                Block
	Access               uint16
	Pos                  *Pos
	Name                 string
	Fields               map[string]*ClassField
	Methods              map[string][]*ClassMethod
	SuperClassExpression *Expression // a or a.b
	SuperClassName       string
	SuperClass           *Class
	Constructors         []*ClassMethod // can be nil
	Signature            *class_json.ClassSignature
	SouceFile            string
	Used                 bool
	VariableType         VariableType
	ClosureVars          map[string]*VariableDefinition // closure variable
}

func (c *Class) mkVariableType() {
	c.VariableType.Typ = VARIABLE_TYPE_CLASS
	c.VariableType.Class = c
}

func (c *Class) check(father *Block) []error {
	errs := make([]error, 0)
	c.Block.check(father)
	errs = append(errs, c.checkFields()...)
	errs = append(errs, c.checkMethods()...)
	return errs
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
			errs = append(errs, vv.Func.check(&c.Block)...)
		}
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
func (c *Class) accessMethod(name string, args []*VariableType) (f *ClassMethod, accessable bool, err error) {

	return
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
