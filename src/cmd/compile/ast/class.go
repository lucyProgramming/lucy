package ast

const (
	ACCESS_PUBLIC = iota
	ACCESS_PROTECTED
	ACCESS_PRIVATE
)

type Class struct {
	Pos         *Pos
	Name        string
	Fields      map[string]*ClassField
	Methods     map[string]*ClassMethod
	Father      *Expression // a or a.b
	Constructor *Function   // can be nil
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
		errs = append(errs, v.Func.check(nil)...)
	}
	return errs
}

type ClassMethod struct {
	ClassFieldProperty
	Func Function
}

type ClassFieldProperty struct {
	IsStatic bool //static or not
	AccessProperty
}

type ClassField struct {
	ClassFieldProperty
	Name string
	Typ  VariableType
	Init *Expression // init value
	Tag  string      //for reflect
	Pos  *Pos
}
