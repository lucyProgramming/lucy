package ast

import (
	"fmt"
	"strings"
)

const (
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
	VARIABLE_TYPE_ARRAY
	VARIABLE_TYPE_JAVA_ARRAY
	VARIABLE_TYPE_FUNCTION
	VARIABLE_TYPE_ENUM

	VARIABLE_TYPE_CLASS

	VARIABLE_TYPE_NAME
	VARIABLE_TYPE_VOID

	VARIABLE_TYPE_PACKAGE
	VARIABLE_TYPE_NULL
)

type VariableType struct {
	Resolved  bool
	Pos       *Pos
	Typ       int
	Name      string
	ArrayType *VariableType
	Const     *Const
	Var       *VariableDefinition
	Class     *Class
	Enum      *Enum
	EnumName  *EnumName
	Function  *Function
	Map       *Map
	Package   *Package
	Alias     string
}

type Map struct {
	K *VariableType
	V *VariableType
}

func (v *VariableType) mkDefaultValueExpression() *Expression {
	var e Expression
	e.IsCompileAuto = true
	e.Pos = v.Pos
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data = false
		e.Value = v.Clone()
	case VARIABLE_TYPE_BYTE:
		e.Typ = EXPRESSION_TYPE_BYTE
		e.Data = byte(0)
		e.Value = v.Clone()
	case VARIABLE_TYPE_SHORT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = int32(0)
		e.Value = v.Clone()
	case VARIABLE_TYPE_INT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = int32(0)
		e.Value = v.Clone()
	case VARIABLE_TYPE_LONG:
		e.Typ = EXPRESSION_TYPE_LONG
		e.Data = int64(0)
		e.Value = v.Clone()
	case VARIABLE_TYPE_FLOAT:
		e.Typ = EXPRESSION_TYPE_FLOAT
		e.Data = float32(0)
		e.Value = v.Clone()
	case VARIABLE_TYPE_DOUBLE:
		e.Typ = EXPRESSION_TYPE_DOUBLE
		e.Data = float64(0)
		e.Value = v.Clone()
	case VARIABLE_TYPE_STRING:
		e.Typ = EXPRESSION_TYPE_STRING
		e.Data = ""
		e.Value = v.Clone()
	case VARIABLE_TYPE_OBJECT:
		fallthrough
	case VARIABLE_TYPE_MAP:
		fallthrough
	case VARIABLE_TYPE_ARRAY:
		e.Typ = EXPRESSION_TYPE_NULL
		e.Value = v.Clone()
	case VARIABLE_TYPE_ENUM:
		e.Typ = EXPRESSION_TYPE_INT
		e.Value = v.Clone()
		e.Data = v.Enum.Enums[0].Value
	}
	return &e
}

func (v *VariableType) RightValueValid() bool {
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
		v.Typ == VARIABLE_TYPE_NULL ||
		v.Typ == VARIABLE_TYPE_JAVA_ARRAY ||
		v.Typ == VARIABLE_TYPE_ENUM
}

/*
	isTyped means can get type from this
*/
func (v *VariableType) isTyped() bool {
	return v.RightValueValid() && v.Typ != VARIABLE_TYPE_NULL
}

/*
	shallow clone
*/
func (t *VariableType) Clone() *VariableType {
	ret := &VariableType{}
	*ret = *t
	return ret
}

func (t *VariableType) resolve(block *Block) error {
	if t == nil {
		return nil
	}
	if t.Resolved {
		return nil
	}
	t.Resolved = true
	if t.Typ == VARIABLE_TYPE_NAME { //
		return t.resolveName(block)
	}
	if t.Typ == VARIABLE_TYPE_ARRAY {
		return t.ArrayType.resolve(block)
	}
	if t.Typ == VARIABLE_TYPE_MAP {
		var err error
		if t.Map.K != nil {
			err = t.Map.K.resolve(block)
			if err != nil {
				return err
			}
		}
		if t.Map.V != nil {
			return t.Map.V.resolve(block)
		}
	}
	return nil
}

func (t *VariableType) resolveNameFromImport() (d interface{}, err error) {
	if strings.Contains(t.Name, ".") == false {
		i := PackageBeenCompile.getImport(t.Pos.Filename, t.Name)
		if i != nil {
			return PackageBeenCompile.load(i.Resource)
		}
		return nil, fmt.Errorf("%s '%s' not found", errMsgPrefix(t.Pos), t.Name)
	}
	packageAndName := strings.Split(t.Name, ".")
	i := PackageBeenCompile.getImport(t.Pos.Filename, packageAndName[0])
	if nil == i {
		return nil, fmt.Errorf("%s package '%s' not imported", errMsgPrefix(t.Pos), packageAndName[0])
	}
	p, err := PackageBeenCompile.load(i.Resource)
	if err != nil {
		return nil, fmt.Errorf("%s %v", errMsgPrefix(t.Pos), err)
	}
	if pp, ok := p.(*Package); ok == false && pp != nil {
		return nil, fmt.Errorf("%s '%s' is not a package", errMsgPrefix(t.Pos), packageAndName[0])
	} else {
		d = pp.Block.SearchByName(packageAndName[1])
		if d == nil {
			err = fmt.Errorf("%s '%s' not found", errMsgPrefix(t.Pos))
		}
		return d, err
	}

}

func (t *VariableType) mkTypeFromInterface(d interface{}) error {
	switch d.(type) {
	case *Class:
		dd := d.(*Class)
		if t != nil {
			t.Typ = VARIABLE_TYPE_OBJECT
			t.Class = dd
			return nil
		}
	case *VariableType:
		dd := d.(*VariableType)
		if dd != nil {
			tt := dd.Clone()
			tt.Pos = t.Pos
			*t = *tt
			return nil
		}
	case *Enum:
		dd := d.(*Enum)
		if dd != nil {
			t.Typ = VARIABLE_TYPE_ENUM
			t.Enum = dd
			return nil
		}
	}
	return fmt.Errorf("%s name '%s' is not a type", errMsgPrefix(t.Pos), t.Name)
}

func (t *VariableType) resolveName(block *Block) error {
	var err error
	var d interface{}
	if strings.Contains(t.Name, ".") == false {
		d = block.SearchByName(t.Name)
		loadFromImport := (d == nil)
		if loadFromImport == false { // d is not nil
			switch d.(type) {
			case *Class:
				if t := d.(*Class); t == nil {
					loadFromImport = true
				} else {
					_, loadFromImport = shouldAccessFromImports(t.Name, t.Pos, t.Pos)
				}
			case *VariableType:
				if t := d.(*VariableType); t == nil {
					loadFromImport = true
				} else {
					_, loadFromImport = shouldAccessFromImports(t.Name, t.Pos, t.Pos)
				}
			case *Enum:
				if t := d.(*Enum); t == nil {
					loadFromImport = true
				} else {
					_, loadFromImport = shouldAccessFromImports(t.Name, t.Pos, t.Pos)
				}

			}
		}
		if loadFromImport {
			d, err = t.resolveNameFromImport()
			if err != nil {
				return err
			}
		}
	} else { // a.b  in type situation,must be package name
		d, err = t.resolveNameFromImport()
		if err != nil {
			return err
		}
	}
	if d == nil {
		return fmt.Errorf("%s type named '%s' not found", errMsgPrefix(t.Pos), t.Name)
	}

	return t.mkTypeFromInterface(d)
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

func (t *VariableType) IsNumber() bool {
	return t.IsInteger() || t.IsFloat()

}

func (t *VariableType) IsPointer() bool {
	return t.Typ == VARIABLE_TYPE_OBJECT ||
		t.Typ == VARIABLE_TYPE_ARRAY ||
		t.Typ == VARIABLE_TYPE_MAP ||
		t.Typ == VARIABLE_TYPE_STRING ||
		t.Typ == VARIABLE_TYPE_JAVA_ARRAY
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

func (v *VariableType) IsPrimitive() bool {
	return v.IsNumber() ||
		v.Typ == VARIABLE_TYPE_STRING ||
		v.Typ == VARIABLE_TYPE_BOOL
}

//可读的类型信息
func (v *VariableType) typeString(ret *string) {
	if v.Alias != "" {
		*ret += v.Alias
		return
	}
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
	case VARIABLE_TYPE_CLASS:
		*ret += fmt.Sprintf("class(%s)", v.Class.Name)
	case VARIABLE_TYPE_ENUM:
		*ret += "enum(" + v.Enum.Name + ")"
	case VARIABLE_TYPE_ARRAY:
		*ret += "[]"
		v.ArrayType.typeString(ret)
	case VARIABLE_TYPE_VOID:
		*ret += "void"
	case VARIABLE_TYPE_STRING:
		*ret += "string"
	case VARIABLE_TYPE_OBJECT: // class name
		*ret += "object@(" + v.Class.Name + ")"
	case VARIABLE_TYPE_MAP:
		*ret += "map{"
		*ret += v.Map.K.TypeString()
		*ret += " -> "
		*ret += v.Map.V.TypeString()
		*ret += "}"
	case VARIABLE_TYPE_JAVA_ARRAY:
		*ret += v.ArrayType.TypeString() + "[]"
	case VARIABLE_TYPE_PACKAGE:
		*ret += v.Package.Name
	case VARIABLE_TYPE_NULL:
		*ret += "null"
	case VARIABLE_TYPE_NAME:
		*ret += v.Name // resove wrong, but typeString is ok to return
	case VARIABLE_TYPE_FUNCTION:
		*ret += v.Function.readableMsg()
	default:
		fmt.Println(VARIABLE_TYPE_JAVA_ARRAY)
		panic(v.Typ)
	}
}

//可读的类型信息
func (v *VariableType) TypeString() string {
	t := ""
	v.typeString(&t)
	return t
}

func (t *VariableType) TypeCompatible(t2 *VariableType) bool {
	// if t.IsInteger() && t2.IsInteger() {
	// 	return true
	// }
	// if t.IsFloat() && t2.IsFloat() {
	// 	return true
	// }
	return t.Equal(t2)
}

/*
	t2 can be cast to t1
*/
func (v *VariableType) Equal(assignMent *VariableType, subPart ...bool) bool {
	if v == assignMent {
		return true
	}

	if v.IsPrimitive() && assignMent.IsPrimitive() {
		return v.Typ == assignMent.Typ
	}
	if v.IsPointer() && assignMent.Typ == VARIABLE_TYPE_NULL {
		return true
	}
	if v.Typ == VARIABLE_TYPE_ARRAY && assignMent.Typ == VARIABLE_TYPE_ARRAY {
		return v.ArrayType.Equal(assignMent.ArrayType, true)
	}
	if v.Typ == VARIABLE_TYPE_JAVA_ARRAY && assignMent.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return v.ArrayType.Equal(assignMent.ArrayType, true)
	}

	if v.Typ == VARIABLE_TYPE_ENUM && assignMent.Typ == VARIABLE_TYPE_ENUM {
		return v.Enum.Name == assignMent.Enum.Name
	}
	if v.Typ == VARIABLE_TYPE_MAP && assignMent.Typ == VARIABLE_TYPE_MAP {
		return v.Map.K.Equal(assignMent.Map.K, true) && v.Map.V.Equal(assignMent.Map.V, true)
	}
	if v.Typ == VARIABLE_TYPE_OBJECT && assignMent.Typ == VARIABLE_TYPE_OBJECT { // object
		if len(subPart) > 0 {
			return v.Class.Name == assignMent.Class.Name
		} else {
			if v.Class.IsInterface() {
				i, _ := assignMent.Class.implemented(v.Class.Name)
				return i
			} else { // class
				has, _ := assignMent.Class.haveSuper(v.Class.Name)
				return has
			}
		}
	}
	return false
}
