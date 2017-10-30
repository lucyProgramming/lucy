package ast

const (
	ACCESS_PUBLIC = iota
	ACCESS_PROTECTED
	ACCESS_PRIVATE
)

type Class struct {
	Fields      []*ClassField
	Methods     []*ClassMethod
	Father      *Class
	Constructor *Function // can be nil
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
}
