package ast

import (
	"fmt"
	"strings"
)

const (
	//primitive type
	VARIABLE_TYPE_BOOL = iota
	VARIABLE_TYPE_BYTE
	VARIABLE_TYPE_SHORT
	VARIABLE_TYPE_CHAR
	VARIABLE_TYPE_INT
	VARIABLE_TYPE_LONG
	VARIABLE_TYPE_FLOAT
	VARIABLE_TYPE_DOUBLE
	VARIABLE_TYPE_STRING
	//function
	VARIALBE_TYPE_FUNCTION
	//enum
	VARIALBE_TYPE_ENUM //enum
	//class
	VARIABLE_TYPE_CLASS //
	VARIABLE_TYPE_ARRAY // []int
	VARIABLE_TYPE_NAME  // naming
	VARIABLE_TYPE_VOID
)

type VariableType struct {
	Pos             *Pos
	Typ             int
	Name            string // Lname.Rname
	CombinationType *VariableType
	Class           *Class
	Enum            *Enum
	EnumName        *EnumName
	FunctionType    *FunctionType
}

/*
	clone a type
	only copy primitive type or type basic on name
*/
func (t *VariableType) CloneNameTypeOrPrimitive() *VariableType {
	ret := &VariableType{}
	*ret = *t // primitive copied,name should be copied too
	if ret.Typ == VARIABLE_TYPE_ARRAY {
		ret.CombinationType = t.CombinationType.CloneNameTypeOrPrimitive()
	}
	return nil
}

func (t *VariableType) resolve(block *Block) error {
	if t.Typ == VARIABLE_TYPE_NAME { //
		err := t.resolveName(block)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (t *VariableType) resolveName(block *Block) error {
	var d interface{}
	var err error
	if !strings.Contains(t.Name, ".") {
		d, err = block.searchByName(t.Name)
		if err != nil {
			return err
		}
	} else { // a.b  in type situation,must be package name
		d, err = t.resolvePackageName(block)
		if err != nil {
			return err
		}
	}
	switch d.(type) {
	case *VariableDefinition:
		return fmt.Errorf("name %s is a variable,not a type", t.Name)
	case []*Function:
		return fmt.Errorf("name %s is a function,not a type", t.Name)
	case *Const:
		return fmt.Errorf("name %s is a const,not a type", t.Name)
	case *Class:
		t.Typ = VARIABLE_TYPE_CLASS
		t.Class = d.(*Class)
	case *Enum:
		t.Typ = VARIALBE_TYPE_ENUM
		t.Enum = d.(*Enum)
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
		return nil, fmt.Errorf("package %v not imported", accessname)
	}
	pname, accessResouce, err := f.fullPackageName(accessname, t.Name)
	if err != nil {
		return nil, err
	}
	p, err := block.loadPackage(pname)
	if err != nil {
		return nil, err
	}
	return p.Block.searchByName(accessResouce)
}

/*
	mk jvm descriptor
*/
func (v *VariableType) Descriptor() string {
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		return "Z"
	case VARIABLE_TYPE_BYTE:
		return "B"
	case VARIABLE_TYPE_SHORT:
		return "S"
	case VARIABLE_TYPE_CHAR:
		return "C"
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
		return "Ljava/lang/String"
	}
	panic("unhandle type signature")
}

func (t *VariableType) typeCompatible(t2 *VariableType) bool {
	if t.Equal(t2) {
		return true
	}
	return false
}

/*
	check const if valid
	number
*/
func (t *VariableType) isNumber() bool {
	return t.Typ == VARIABLE_TYPE_BYTE ||
		t.Typ == VARIABLE_TYPE_SHORT ||
		t.Typ == VARIABLE_TYPE_CHAR ||
		t.Typ == VARIABLE_TYPE_INT ||
		t.Typ == VARIABLE_TYPE_LONG ||
		t.Typ == VARIABLE_TYPE_FLOAT ||
		t.Typ == VARIABLE_TYPE_DOUBLE
}

/*
	float or double
*/
func (t *VariableType) isFloat() bool {
	return t.Typ == VARIABLE_TYPE_FLOAT || t.Typ == VARIABLE_TYPE_DOUBLE
}

func (t *VariableType) constValueValid(e *Expression) (data interface{}, err error) {
	//number type
	if t.isNumber() {
		if !t.isNumber() {
			err = fmt.Errorf("expression is not number")
			return
		}
		if t.isFloat() {
			f := e.literalValue2Float64()
			if t.Typ == VARIABLE_TYPE_FLOAT {
				return float32(f), nil
			} else {
				return f, nil
			}
		} else {
			d := e.literalValue2Int64()
			switch t.Typ {
			case VARIABLE_TYPE_BYTE:
				return byte(d), nil
			case VARIABLE_TYPE_SHORT:
				return int16(d), nil
			case VARIABLE_TYPE_CHAR:
				return int16(d), nil
			case VARIABLE_TYPE_INT:
				return int32(d), nil
			case VARIABLE_TYPE_LONG:
				return d, nil
			case VARIABLE_TYPE_FLOAT:
				return float32(d), nil
			case VARIABLE_TYPE_DOUBLE:
				return float64(d), nil
			}
		}
	}
	if t.Typ == VARIABLE_TYPE_BOOL {
		return e.canBeCovert2Bool()
	}
	if t.Typ == VARIABLE_TYPE_STRING {
		if e.Typ != EXPRESSION_TYPE_STRING {
			return nil, fmt.Errorf("expression is not string but %s", e.OpName())
		}
		return e.Data.(string), nil
	}
	return nil, fmt.Errorf("cannot assign %s to %s", e.OpName(), t.TypeString)
}

//assign some simple expression
func (t *VariableType) assignExpression(p *Package, e *Expression) (data interface{}, err error) {
	switch t.Typ {
	case VARIABLE_TYPE_BOOL:
		if e.Typ == EXPRESSION_TYPE_BOOL {
			data = e.Data.(bool)
			return
		}
	case VARIALBE_TYPE_ENUM:
		if e.Typ == EXPRESSION_TYPE_IDENTIFIER {
			if _, ok := p.Block.EnumNames[e.Data.(string)]; ok {
				data = p.Block.EnumNames[e.Data.(string)]
			}
		}
	case VARIABLE_TYPE_BYTE:
		if e.Typ == EXPRESSION_TYPE_BYTE {
			data = e.Data.(byte)
			return
		}
	case VARIABLE_TYPE_INT:
		if e.Typ == EXPRESSION_TYPE_BYTE {
			data = int64(e.Data.(byte))
			return
		} else if e.Typ == EXPRESSION_TYPE_INT {
			data = e.Data.(int64)
			return
		}
	case VARIABLE_TYPE_FLOAT:
		if e.Typ == EXPRESSION_TYPE_BYTE {
			data = int64(e.Data.(byte))
			return
		} else if e.Typ == EXPRESSION_TYPE_INT {
			data = e.Data.(int64)
			return
		} else if EXPRESSION_TYPE_FLOAT == e.Typ {
			data = e.Data.(float64)
			return
		}
	case VARIABLE_TYPE_STRING:
		if e.Typ == EXPRESSION_TYPE_STRING {
			data = e.Data.(string)
			return
		}
	case VARIABLE_TYPE_CLASS:
		if e.Typ == EXPRESSION_TYPE_NULL {
			return nil, nil // null pointer
		}
	}
	err = fmt.Errorf("can`t covert type accroding to type")
	return
}

////把树型转化为可读字符串
//func (c *CombinationType) TypeString(ret *string) {
//	if c.Typ == COMBINATION_TYPE_ARRAY {
//		*ret += "[]"
//	}
//	c.Combination.TypeString(ret)
//}

//可读的类型信息
func (v *VariableType) TypeString_(ret *string) {
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		*ret = "bool"
	case VARIABLE_TYPE_BYTE:
		*ret = "byte"
	case VARIABLE_TYPE_INT:
		*ret = "int"
	case VARIABLE_TYPE_FLOAT:
		*ret = "float"
	case VARIALBE_TYPE_FUNCTION:
		*ret = "function"
	case VARIABLE_TYPE_CLASS:
		*ret = v.Name
	case VARIALBE_TYPE_ENUM:
		*ret = v.Name + "(enum)"
	case VARIABLE_TYPE_ARRAY:
		*ret += "[]"
		v.CombinationType.TypeString_(ret)
	}
}

//可读的类型信息
func (v *VariableType) TypeString() string {
	t := ""
	v.TypeString_(&t)
	return t
}

func (v *VariableType) Equal(e *VariableType) bool {
	if v.Typ != e.Typ {
		return false
	}
	if v.CombinationType.Typ != e.CombinationType.Typ {
		return false
	}
	return v.CombinationType.Equal(e.CombinationType)
}
