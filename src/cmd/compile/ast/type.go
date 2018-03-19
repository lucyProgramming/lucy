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
	//ref types

	VARIABLE_TYPE_STRING
	VARIABLE_TYPE_OBJECT
	VARIABLE_TYPE_MAP
	VARIABLE_TYPE_ARRAY      //[]int
	VARIABLE_TYPE_JAVA_ARRAY // java array int[]
	VARIABLE_TYPE_FUNCTION
	//function type
	VARIABLE_TYPE_FUNCTION_TYPE
	//enum
	VARIABLE_TYPE_ENUM //enum
	//class
	VARIABLE_TYPE_CLASS //

	VARIABLE_TYPE_NAME //naming
	VARIABLE_TYPE_VOID
	VARIABLE_TYPE_NULL
)

type VariableType struct {
	Pos       *Pos
	Typ       int
	Name      string
	ArrayType *VariableType
	Const     *Const
	Var       *VariableDefinition
	Class     *Class
	Enum      *Enum
	Function  *Function
	Map       *Map
}

type Map struct {
	K *VariableType
	V *VariableType
}

func (v *VariableType) mkDefaultValueExpression() *Expression {
	var e Expression
	e.IsCompileAutoExpression = true
	e.Pos = v.Pos
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data = false
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_BYTE:
		e.Typ = EXPRESSION_TYPE_BYTE
		e.Data = byte(0)
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_SHORT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = int32(0)
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_INT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = int32(0)
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_LONG:
		e.Typ = EXPRESSION_TYPE_LONG
		e.Data = int64(0)
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_FLOAT:
		e.Typ = EXPRESSION_TYPE_FLOAT
		e.Data = float32(0)
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_DOUBLE:
		e.Typ = EXPRESSION_TYPE_DOUBLE
		e.Data = float64(0)
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_STRING:
		e.Typ = EXPRESSION_TYPE_STRING
		e.Data = ""
		e.VariableType = v.Clone()
	case VARIABLE_TYPE_OBJECT:
		fallthrough
	case VARIABLE_TYPE_MAP:
		fallthrough
	case VARIABLE_TYPE_ARRAY:
		e.Typ = EXPRESSION_TYPE_NULL
		e.VariableType = v.Clone()
	default:
		panic("....")
	}
	return &e
}

func (v *VariableType) JvmSlotSize() uint16 {
	if v.rightValueValid() == false {
		panic("right value invalid")
	}
	return JvmSlotSizeHandler(v)
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
		v.Typ == VARIABLE_TYPE_ARRAY ||
		v.Typ == VARIABLE_TYPE_MAP ||
		v.Typ == VARIABLE_TYPE_NULL
}

/*
	isTyped means can get type from this
*/
func (v *VariableType) isTyped() bool {
	return v.Typ == VARIABLE_TYPE_BOOL ||
		v.Typ == VARIABLE_TYPE_BYTE ||
		v.Typ == VARIABLE_TYPE_SHORT ||
		v.Typ == VARIABLE_TYPE_INT ||
		v.Typ == VARIABLE_TYPE_LONG ||
		v.Typ == VARIABLE_TYPE_FLOAT ||
		v.Typ == VARIABLE_TYPE_DOUBLE ||
		v.Typ == VARIABLE_TYPE_STRING ||
		v.Typ == VARIABLE_TYPE_OBJECT ||
		v.Typ == VARIABLE_TYPE_ARRAY ||
		v.Typ == VARIABLE_TYPE_MAP
}

/*
	clone a type
*/
func (t *VariableType) Clone() *VariableType {
	ret := &VariableType{}
	*ret = *t
	if ret.Typ == VARIABLE_TYPE_ARRAY {
		ret.ArrayType = t.ArrayType.Clone()
	}
	return ret
}

/*
	when typ is VARIABLE_TYPE_CLASS ,convert to object assoiate with class
*/
func (t *VariableType) actionNeedBeenDoneWhenDescribeVariable() {
	if t == nil { // in some case identifier is insert into block,but typ is still wrong
		return
	}
	if t.Typ == VARIABLE_TYPE_CLASS {
		t.Typ = VARIABLE_TYPE_OBJECT // correct
		return                       // no further need be done
	}
	if t.Typ == VARIABLE_TYPE_ARRAY {
		t.ArrayType.actionNeedBeenDoneWhenDescribeVariable() // recursive do action
		return
	}
}

func (t *VariableType) resolve(block *Block) error {
	if t.isPrimitive() {
		return nil
	}
	if t.Typ == VARIABLE_TYPE_NAME { //
		return t.resolveName(block)
	}
	if t.Typ == VARIABLE_TYPE_ARRAY {
		return t.ArrayType.resolve(block)
	}
	return nil
}

func (t *VariableType) resolveName(block *Block) error {
	var d interface{}
	if !strings.Contains(t.Name, ".") {
		d = block.searchByName(t.Name)
		if d == nil {
			return fmt.Errorf("%s not found", t.Name)
		}
	} else { // a.b  in type situation,must be package name
		//d, err = t.resolvePackageName(block)
		//if err != nil {
		//	return err
		//}
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

/*
	check const if valid
	number
*/
func (t *VariableType) IsNumber() bool {
	return t.IsInteger() || t.IsFloat()

}

func (t *VariableType) IsPointer() bool {
	return t.Typ == VARIABLE_TYPE_OBJECT ||
		t.Typ == VARIABLE_TYPE_ARRAY ||
		t.Typ == VARIABLE_TYPE_MAP ||
		t.Typ == VARIABLE_TYPE_STRING

}

func (t *VariableType) IsInteger() bool {
	return t.Typ == VARIABLE_TYPE_BYTE ||
		t.Typ == VARIABLE_TYPE_SHORT ||
		t.Typ == VARIABLE_TYPE_INT ||
		t.Typ == VARIABLE_TYPE_LONG
}

/*
	float or double
*/
func (t *VariableType) IsFloat() bool {
	return t.Typ == VARIABLE_TYPE_FLOAT ||
		t.Typ == VARIABLE_TYPE_DOUBLE
}

func (v *VariableType) isPrimitive() bool {
	return v.IsNumber() ||
		v.Typ == VARIABLE_TYPE_STRING ||
		v.Typ == VARIABLE_TYPE_BOOL
}

//可读的类型信息
func (v *VariableType) typeString(ret *string) {
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		*ret += "bool"
	case VARIABLE_TYPE_BYTE:
		*ret += "byte"
	case VARIABLE_TYPE_SHORT:
		*ret += "short"
	case VARIABLE_TYPE_INT:
		*ret += "int"
	case VARIABLE_TYPE_LONG:
		*ret += "long"
	case VARIABLE_TYPE_FLOAT:
		*ret += "float"
	case VARIABLE_TYPE_DOUBLE:
		*ret += "double"
	case VARIABLE_TYPE_FUNCTION:
		*ret += "function(" + v.Function.Name + ")"
	case VARIABLE_TYPE_CLASS:
		*ret += v.Name
	case VARIABLE_TYPE_ENUM:
		*ret += "enum(" + v.Name + ")"
	case VARIABLE_TYPE_ARRAY:
		*ret += "[]"
		v.ArrayType.typeString(ret)
	case VARIABLE_TYPE_VOID:
		*ret += "void"
	case VARIABLE_TYPE_STRING:
		*ret += "string"
	case VARIABLE_TYPE_OBJECT: // class name
		*ret += v.Class.Name
	case VARIABLE_TYPE_MAP:
		*ret += "map{"
		*ret += v.Map.K.TypeString()
		*ret += " -> "
		*ret += v.Map.V.TypeString()
		*ret += "}"
	case VARIABLE_TYPE_JAVA_ARRAY:
		*ret += v.ArrayType.TypeString() + "[]"
	default:
		panic(1)
	}
}

//可读的类型信息
func (v *VariableType) TypeString() string {
	t := ""
	v.typeString(&t)
	return t
}

func (t1 *VariableType) Equal(t2 *VariableType) bool {
	if t1 == t2 {
		return true
	}
	if t1.isPrimitive() || t2.isPrimitive() {
		return t1.Typ == t2.Typ
	}
	if t1.IsPointer() && t1.Typ == VARIABLE_TYPE_NULL {
		return true
	}
	return t1.ArrayType.Equal(t2.ArrayType)
}
