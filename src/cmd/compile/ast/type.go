package ast

import (
	"fmt"
	"strings"
)

const (
	//primitive type
	_ = iota
	//value types
	VARIABLE_TYPE_BOOL
	VARIABLE_TYPE_BYTE
	VARIABLE_TYPE_SHORT
	VARIABLE_TYPE_INT
	VARIABLE_TYPE_LONG
	VARIABLE_TYPE_FLOAT
	VARIABLE_TYPE_DOUBLE
	VARIABLE_TYPE_STRING
	VARIABLE_TYPE_OBJECT
	VARIABLE_TYPE_ARRAY_INSTANCE
	VARIABLE_TYPE_FUNCTION
	//function type
	VARIABLE_TYPE_FUNCTION_TYPE
	//enum
	VARIABLE_TYPE_ENUM //enum
	//class
	VARIABLE_TYPE_CLASS //
	VARIABLE_TYPE_ARRAY //[]int
	VARIABLE_TYPE_NAME  //naming
	VARIABLE_TYPE_VOID
	VARIABLE_TYPE_NULL
)

type VariableType struct {
	Pos             *Pos
	Typ             int
	Name            string
	CombinationType *VariableType
	FunctionType    *FunctionType
	Const           *Const
	Var             *VariableDefinition
	Class           *Class
	Enum            *Enum
	Function        *Function
}

func (v *VariableType) JvmSlotSize() uint16 {
	if v.rightValueValid() == false {
		panic("right value invalid")
	}
	if v.Typ == VARIABLE_TYPE_DOUBLE || VARIABLE_TYPE_LONG == v.Typ {
		return 2
	} else {
		return 1
	}
}
func (v *VariableType) rightValueValid() bool {
	return v.Typ == VARIABLE_TYPE_BOOL ||
		v.Typ == VARIABLE_TYPE_BYTE ||
		v.Typ == VARIABLE_TYPE_SHORT ||
		v.Typ == VARIABLE_TYPE_INT ||
		v.Typ == VARIABLE_TYPE_LONG ||
		v.Typ == VARIABLE_TYPE_FLOAT ||
		v.Typ == VARIABLE_TYPE_DOUBLE ||
		v.Typ == VARIABLE_TYPE_STRING ||
		v.Typ == VARIABLE_TYPE_OBJECT ||
		v.Typ == VARIABLE_TYPE_ARRAY_INSTANCE
}

/*
	clone a type
*/
func (t *VariableType) Clone() *VariableType {
	ret := &VariableType{}
	*ret = *t
	if ret.Typ == VARIABLE_TYPE_ARRAY {
		ret.CombinationType = t.CombinationType.Clone()
	}
	return ret
}

func (t *VariableType) resolve(block *Block) error {
	if t.isPrimitive() {
		return nil
	}
	if t.Typ == VARIABLE_TYPE_NAME { //
		return t.resolveName(block)
	}
	if t.Typ == VARIABLE_TYPE_ARRAY {
		return t.CombinationType.resolve(block)
	}
	return nil
}

func (t *VariableType) resolveName(block *Block) error {
	var d interface{}
	var err error
	if !strings.Contains(t.Name, ".") {
		d = block.searchByName(t.Name)
		if d == nil {
			return fmt.Errorf("%s not found", t.Name)
		}
	} else { // a.b  in type situation,must be package name
		d, err = t.resolvePackageName(block)
		if err != nil {
			return err
		}
	}
	switch d.(type) {
	case *VariableDefinition:
		return fmt.Errorf("%s name %s is a variable,not a type", errMsgPrefix(t.Pos), t.Name)
	case *Function:
		return fmt.Errorf("%s name %s is a function,not a type", errMsgPrefix(t.Pos), t.Name)
	case *Const:
		return fmt.Errorf("%s name %s is a const,not a type", errMsgPrefix(t.Pos), t.Name)
	case *Class:
		*t = d.(*Class).VariableType
		return nil
	case *Enum:
		*t = d.(*Enum).VariableType
		return nil
	default:
		return fmt.Errorf("name %s is not type")
	}
	return nil
}

func (t *VariableType) resolvePackageName(block *Block) (interface{}, error) {
	accessname := t.Name[0:strings.Index(t.Name, ".")] // package name
	var f *File
	var ok bool
	if f, ok = block.InheritedAttribute.p.Files[t.Pos.Filename]; !ok {
		return nil, fmt.Errorf("%s package %v not imported", errMsgPrefix(t.Pos), accessname)
	}
	if _, ok = f.Imports[accessname]; !ok {
		return nil, fmt.Errorf("%s package %s is not imported", errMsgPrefix(t.Pos), accessname)
	}
	p, err := block.loadPackage(f.Imports[accessname].Name)
	if err != nil {
		return nil, err
	}
	d := p.Block.searchByName(strings.Trim(t.Name, accessname+"."))
	if d == nil {
		err = fmt.Errorf("%s not found", t.Name)
	}
	return d, err
}

/*
	number convert rule
*/
func (t *VariableType) NumberTypeConvertRule(t2 *VariableType) int {
	if t.Typ == t2.Typ {
		return t.Typ
	}
	if t.Typ == VARIABLE_TYPE_DOUBLE || t2.Typ == VARIABLE_TYPE_DOUBLE {
		return VARIABLE_TYPE_DOUBLE
	}
	if t.Typ == VARIABLE_TYPE_FLOAT || t2.Typ == VARIABLE_TYPE_FLOAT {
		if t.Typ == VARIABLE_TYPE_LONG || t2.Typ == VARIABLE_TYPE_LONG {
			return VARIABLE_TYPE_DOUBLE
		} else {
			return VARIABLE_TYPE_FLOAT
		}
	}
	if t.Typ == VARIABLE_TYPE_LONG || t2.Typ == VARIABLE_TYPE_LONG {
		return VARIABLE_TYPE_LONG
	}
	if t.Typ == VARIABLE_TYPE_INT || t2.Typ == VARIABLE_TYPE_INT {
		return VARIABLE_TYPE_INT
	}
	if t.Typ == VARIABLE_TYPE_SHORT || t2.Typ == VARIABLE_TYPE_SHORT {
		return VARIABLE_TYPE_SHORT
	}
	return VARIABLE_TYPE_BYTE
}

func (t *VariableType) TypeCompatible(comp *VariableType, err ...*error) bool {
	if t.Equal(comp) {
		return true
	}
	if t.IsNumber() && comp.IsNumber() {
		return true
	}
	if t.Typ != VARIABLE_TYPE_OBJECT || comp.Typ != VARIABLE_TYPE_OBJECT {
		return false
	}
	return comp.Class.instanceOf(t.Class)
}

func (v *VariableType) Descriptor() string {
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		return "Z"
	case VARIABLE_TYPE_BYTE:
		return "B"
	case VARIABLE_TYPE_SHORT:
		return "S"
		//	case VARIABLE_TYPE_CHAR:
		//		return "C"
	case VARIABLE_TYPE_INT:
		return "I"
	case VARIABLE_TYPE_LONG:
		return "J"
	case VARIABLE_TYPE_FLOAT:
		return "F"
	case VARIABLE_TYPE_DOUBLE:
		return "D"
	case VARIABLE_TYPE_ARRAY:
		return "[" + v.CombinationType.Descriptor()
	case VARIABLE_TYPE_STRING:
		return "Ljava/lang/String;"
	case VARIABLE_TYPE_VOID:
		return "V"
	case VARIABLE_TYPE_OBJECT:
		return "L" + v.Class.Name + ";"
	}
	panic("unhandle type signature")
}

/*
	check const if valid
	number
*/
func (t *VariableType) IsNumber() bool {
	return t.IsInteger() || t.IsFloat()

}

func (t *VariableType) IsPointer() bool {
	return t.Typ == VARIABLE_TYPE_OBJECT || t.Typ == VARIABLE_TYPE_ARRAY_INSTANCE

}

func (t *VariableType) IsInteger() bool {
	return t.Typ == VARIABLE_TYPE_BYTE ||
		t.Typ == VARIABLE_TYPE_SHORT ||
		//		t.Typ == VARIABLE_TYPE_CHAR ||
		t.Typ == VARIABLE_TYPE_INT ||
		t.Typ == VARIABLE_TYPE_LONG
}

/*
	float or double
*/
func (t *VariableType) IsFloat() bool {
	return t.Typ == VARIABLE_TYPE_FLOAT || t.Typ == VARIABLE_TYPE_DOUBLE
}

func (v *VariableType) isPrimitive() bool {
	return v.IsNumber() || v.Typ == VARIABLE_TYPE_STRING ||
		v.Typ == VARIABLE_TYPE_BOOL
}

//可读的类型信息
func (v *VariableType) typeString_(ret *string) {
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		*ret = "bool"
	case VARIABLE_TYPE_BYTE:
		*ret = "byte"
	case VARIABLE_TYPE_SHORT:
		*ret = "short"
	case VARIABLE_TYPE_INT:
		*ret = "int"
	case VARIABLE_TYPE_LONG:
		*ret = "long"
	case VARIABLE_TYPE_FLOAT:
		*ret = "float"
	case VARIABLE_TYPE_DOUBLE:
		*ret = "double"
	case VARIABLE_TYPE_FUNCTION:
		*ret = "function(" + v.Function.Name + ")"
	case VARIABLE_TYPE_CLASS:
		*ret = v.Name
	case VARIABLE_TYPE_ENUM:
		*ret = "enum(" + v.Name + ")"
	case VARIABLE_TYPE_ARRAY:
		*ret += "[]"
		v.CombinationType.typeString_(ret)
	case VARIABLE_TYPE_VOID:
		*ret = "void"
	case VARIABLE_TYPE_STRING:
		*ret = "string"
	case VARIABLE_TYPE_OBJECT:
		*ret = "object"
	case VARIABLE_TYPE_ARRAY_INSTANCE:
		*ret = "array_instance"
	default:
		panic(1)
	}
}

//可读的类型信息
func (v *VariableType) TypeString() string {
	t := ""
	v.typeString_(&t)
	return t
}

func (t1 *VariableType) Equal(t2 *VariableType) bool {
	if t1.isPrimitive() || t2.isPrimitive() {
		return t1.Typ == t2.Typ
	}
	return t1.CombinationType.Equal(t2.CombinationType)
}
