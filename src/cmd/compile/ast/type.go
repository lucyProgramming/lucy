package ast

import (
	"fmt"
	"strings"
)

const (
	//primitive type
	_ = iota
	VARIABLE_TYPE_BOOL
	VARIABLE_TYPE_BYTE
	VARIABLE_TYPE_SHORT
	VARIABLE_TYPE_CHAR
	VARIABLE_TYPE_INT
	VARIABLE_TYPE_LONG
	VARIABLE_TYPE_FLOAT
	VARIABLE_TYPE_DOUBLE
	VARIABLE_TYPE_STRING
	VARIABLE_TYPE_OBJECT
	//function
	VARIABLE_TYPE_FUNCTION
	VARIABLE_TYPE_FUNCTION_TYPE
	//enum
	VARIABLE_TYPE_ENUM //enum
	//class
	VARIABLE_TYPE_CLASS //
	VARIABLE_TYPE_ARRAY //[]int
	VARIABLE_TYPE_NAME  //naming
	VARIABLE_TYPE_VOID
)

type VariableType struct {
	Pos             *Pos
	Typ             int
	Name            string // Lname.Rname
	CombinationType *VariableType
	FunctionType    *FunctionType
	Resource        *VariableTypeResource
}

func (v *VariableType) rightValueValid() bool {
	return v.Typ != VARIABLE_TYPE_VOID
}

type VariableTypeResource struct {
	Const    *Const
	Var      *VariableDefinition
	Class    *Class
	Enum     *Enum
	Function *Function
}

/*
	clone a type
*/
func (t *VariableType) Clone() *VariableType {
	ret := &VariableType{}
	ret.Typ = t.Typ // primitive copied,name should be copied too
	ret.Pos = &Pos{}
	*ret.Pos = *t.Pos
	if t.Resource != nil {
		ret.Resource = &VariableTypeResource{}
		*ret.Resource = *t.Resource
	}
	if ret.Typ == VARIABLE_TYPE_ARRAY {
		ret.CombinationType = t.CombinationType.Clone()
	}
	return ret
}

func (t *VariableType) assignAble() error {
	if t.Resource == nil {
		return fmt.Errorf("cannot been assign value")
	}
	if t.Resource.Const != nil {
		return fmt.Errorf("const '%v' cannot been assign value", t.Resource.Const.Name)
	}
	if t.Resource.Var == nil {
		return fmt.Errorf("cannot been assign value")
	}
	return nil
}

func (t *VariableType) markAsUsed() {
	if t.Resource == nil {
		return
	}

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
		*t = *d.(*Class).VariableType
		return nil
	case *Enum:
		*t = *d.(*Enum).VariableType
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
		return nil, fmt.Errorf("package %v not imported", accessname)
	}
	if _, ok = f.Imports[accessname]; !ok {
		return nil, fmt.Errorf("package %s is not imported", accessname)
	}
	p, err := block.loadPackage(f.Imports[accessname].Name)
	if err != nil {
		return nil, err
	}

	return p.Block.searchByName(strings.Trim(t.Name, accessname+"."))
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
	case VARIABLE_TYPE_VOID:
		return "void"
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
	return t.isInteger() || t.isFloat()

}

func (t *VariableType) isInteger() bool {
	return t.Typ == VARIABLE_TYPE_BYTE ||
		t.Typ == VARIABLE_TYPE_SHORT ||
		t.Typ == VARIABLE_TYPE_CHAR ||
		t.Typ == VARIABLE_TYPE_INT ||
		t.Typ == VARIABLE_TYPE_LONG
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
	return nil, fmt.Errorf("cannot assign %s to %s", e.OpName(), t.TypeString())
}

//assign some simple expression
func (t *VariableType) assignExpression(p *Package, e *Expression) (data interface{}, err error) {
	switch t.Typ {
	case VARIABLE_TYPE_BOOL:
		if e.Typ == EXPRESSION_TYPE_BOOL {
			data = e.Data.(bool)
			return
		}
	case VARIABLE_TYPE_ENUM:
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

//可读的类型信息
func (v *VariableType) TypeString_(ret *string) {
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		*ret = "bool"
	case VARIABLE_TYPE_BYTE:
		*ret = "byte"
	case VARIABLE_TYPE_SHORT:
		*ret = "short"
	case VARIABLE_TYPE_CHAR:
		*ret = "char"
	case VARIABLE_TYPE_INT:
		*ret = "int"
	case VARIABLE_TYPE_LONG:
		*ret = "long"
	case VARIABLE_TYPE_FLOAT:
		*ret = "float"
	case VARIABLE_TYPE_DOUBLE:
		*ret = "double"
	case VARIABLE_TYPE_FUNCTION:
		*ret = "function"
	case VARIABLE_TYPE_CLASS:
		*ret = v.Name
	case VARIABLE_TYPE_ENUM:
		*ret = v.Name + "(enum)"
	case VARIABLE_TYPE_ARRAY:
		*ret += "[]"
		v.CombinationType.TypeString_(ret)
	case VARIABLE_TYPE_VOID:
		*ret = "void"
	}
}

//可读的类型信息
func (v *VariableType) TypeString() string {
	t := ""
	v.TypeString_(&t)
	return t
}

func (v *VariableType) isPrimitive() bool {
	return v.isNumber() || v.Typ == VARIABLE_TYPE_STRING
}

func (t1 *VariableType) Equal(t2 *VariableType) bool {
	if t1.isPrimitive() && t2.isPrimitive() {
		return t1.Typ == t2.Typ
	}
	if t1.Typ != t2.Typ {
		return false
	}
	return t1.CombinationType.Equal(t2.CombinationType)
}
