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
	VARIABLE_TYPE_NULL
	VARIABLE_TYPE_STRING
	VARIABLE_TYPE_OBJECT
	VARIABLE_TYPE_MAP
	VARIABLE_TYPE_ARRAY      //[]int
	VARIABLE_TYPE_JAVA_ARRAY // java array int[]
	VARIABLE_TYPE_FUNCTION

	//enum
	VARIABLE_TYPE_ENUM //enum
	//class
	VARIABLE_TYPE_CLASS //

	VARIABLE_TYPE_NAME //naming
	VARIABLE_TYPE_VOID

	VARIABLE_TYPE_PACKAGE //
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
	Function  *Function
	Map       *Map
	Package   *Package
}

type Map struct {
	K *VariableType
	V *VariableType
}

func (v *VariableType) mkDefaultValueExpression() *Expression {
	var e Expression
	e.IsCompileDefaultValueExpression = true
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
	}
	return &e
}

func (v *VariableType) JvmSlotSize() uint16 {
	if v.Typ == VARIABLE_TYPE_VOID {
		return 0
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
	return v.rightValueValid() && v.Typ != VARIABLE_TYPE_NULL
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

func (t *VariableType) resolve(block *Block) error {
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
		return nil, fmt.Errorf("%s class '%s' not found", errMsgPrefix(t.Pos), t.Name)
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
		if pp.Block.SearchByName(packageAndName[1]) == nil {
			return nil, fmt.Errorf("%s '%s' not found", errMsgPrefix(t.Pos), packageAndName[1])
		}
		if pp.Block.Types != nil && pp.Block.Types[packageAndName[1]] != nil {
			return pp.Block.Types[packageAndName[1]], nil
		}
		if pp.Block.Classes != nil && pp.Block.Classes[packageAndName[1]] != nil {
			return pp.Block.Classes[packageAndName[1]], nil
		}
		return nil, fmt.Errorf("%s '%s' is not a type", errMsgPrefix(t.Pos), packageAndName[1])
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
	}
	return fmt.Errorf("%s name '%s' is not a type", errMsgPrefix(t.Pos), t.Name)
}

func (t *VariableType) resolveName(block *Block) error {
	var err error
	var d interface{}
	if strings.Contains(t.Name, ".") == false {
		variableT := block.searchType(t.Name)
		var classT *Class
		useClassT := false
		if variableT != nil {
			classT := block.searchClass(t.Name)
			if classT != nil { // d2 is not nil
				if moreClose(t.Pos, classT.Pos, variableT.Pos) {
					useClassT = true
				}
			}
		} else {
			classT = block.searchClass(t.Name)
			useClassT = true
		}
		loadFromImport := (classT == nil && variableT == nil)
		if loadFromImport == false {
			if useClassT {
				if classT != nil {
					_, loadFromImport = shouldAccessFromImports(t.Name, t.Pos, classT.Pos)
				}
			} else {
				if variableT != nil {
					_, loadFromImport = shouldAccessFromImports(t.Name, t.Pos, variableT.Pos)
				}
			}

		}
		if loadFromImport {
			d, err = t.resolveNameFromImport()
			if err != nil {
				return err
			}
		} else {
			if useClassT {
				d = classT
			} else {
				d = variableT
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

func (v *VariableType) IsPrimitive() bool {
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
	case VARIABLE_TYPE_CLASS:
		*ret += fmt.Sprintf("class named %s", v.Class.Name)
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
		*ret += "object@" + v.Class.Name
	case VARIABLE_TYPE_MAP:
		*ret += "map{"
		*ret += v.Map.K.TypeString()
		*ret += " -> "
		*ret += v.Map.V.TypeString()
		*ret += "}"
	case VARIABLE_TYPE_JAVA_ARRAY:
		*ret += v.ArrayType.TypeString() + "[]"
	case VARIABLE_TYPE_PACKAGE:
		*ret = v.Package.Name
	default:
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
	if t.IsInteger() && t2.IsInteger() {
		return true
	}
	if t.IsFloat() && t2.IsFloat() {
		return true
	}
	return t.Equal(t2)
}

/*
	t2 can be cast to t1
*/
func (t1 *VariableType) Equal(t2 *VariableType) bool {
	if t1 == t2 { // this is not happening
		return true
	}
	if t1.IsPrimitive() || t2.IsPrimitive() {
		return t1.Typ == t2.Typ
	}
	if (t1.IsPointer() && t1.Typ != VARIABLE_TYPE_STRING) && t2.Typ == VARIABLE_TYPE_NULL {
		return true
	}
	if t1.Typ == VARIABLE_TYPE_ARRAY && t2.Typ == VARIABLE_TYPE_ARRAY {
		return t1.ArrayType.Equal(t2.ArrayType)
	}
	if t1.Typ == VARIABLE_TYPE_MAP && t2.Typ == VARIABLE_TYPE_MAP {
		return t1.Map.K.Equal(t1.Map.K) && t1.Map.V.Equal(t1.Map.V)
	}
	if t1.Typ == VARIABLE_TYPE_OBJECT && t2.Typ == VARIABLE_TYPE_OBJECT { // object
		if t1.Class.isInterface() {
			i, _ := t2.Class.implemented(t1.Class.Name)
			return i
		} else { // class
			has, _ := t2.Class.haveSuper(t1.Class.Name)
			return has
		}
	}
	return false
}
