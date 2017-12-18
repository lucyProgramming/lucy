package ast

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
)

type Class struct {
	Access               uint16 // public private or protected
	Pos                  *Pos
	Name                 string
	Fields               map[string]*ClassField
	Methods              map[string][]*ClassMethod
	Consts               map[string]*Const
	SuperClassExpression *Expression // a or a.b
	SuperClassName       string
	SuperClass           *Class
	Constructors         []*ClassMethod // can be nil
	Signature            *class_json.ClassSignature
	SouceFile            string
	Used                 bool
	VariableType         VariableType
}

func (c *Class) mkVariableType() {
	c.VariableType.Typ = VARIABLE_TYPE_CLASS
	c.VariableType.Class = c
}

func (c *Class) check() []error {
	errs := make([]error, 0)
	errs = append(errs, c.checkFields()...)
	errs = append(errs, c.checkFields()...)
	return errs
}

func (c *Class) checkFields() []error {
	errs := []error{}
	//	for _, v := range c.Fields {
	//	}
	return errs
}

func (c *Class) checkMethods() []error {
	errs := []error{}
	for _, v := range c.Methods {
		for _, vv := range v {
			errs = append(errs, vv.Func.check(nil)...)
		}
	}
	return errs
}

func (c *Class) accessField(name string) (f *ClassField, accessable bool, err error) {
	if c.Fields[name] == nil {
		err = fmt.Errorf("field %s not found")
		return
	}
	f = c.Fields[name]
	accessable = (f.AccessFlags | cg.ACC_FIELD_PUBLIC) != 0
	return
}

type ClassMethod struct {
	Func *Function
}

type ClassField struct {
	VariableDefinition
}
