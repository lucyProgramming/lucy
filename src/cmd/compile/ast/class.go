package ast

import "github.com/756445638/lucy/src/cmd/compile/jvm/class_json"

//const (
//	ACCESS_PUBLIC = iota
//	ACCESS_PROTECTED
//	ACCESS_PRIVATE
//)

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

type ClassMethod struct {
	ClassFieldProperty
	Func      *Function
	Signature *class_json.MethodSignature
}

type ClassFieldProperty struct {
	IsStatic    bool   //static or not
	AccessFlags uint16 // public private or protected
}

type ClassField struct {
	ClassFieldProperty
	VariableDefinition
	Tag       string //for reflect
	Pos       *Pos
	Signature *class_json.FieldSignature
}
