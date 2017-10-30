package ast

const (
	CLASS_FIELD_PUBLIC = iota
	CLASS_FIELD_PROTECTED
	CLASS_FIELD_PRIVATE
)

type Class struct {
	Fields      []*ClassField
	Methods     []*ClassMethod
	Father      *Class
	Constructor *Function // can be nil
}

type ClassMethod struct {
	ClassFieldProperty
	Name string
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
}
