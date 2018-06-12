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
	VARIABLE_TYPE_T
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
	Class     *Class
	Enum      *Enum
	EnumName  *EnumName
	Function  *Function
	Map       *Map
	Package   *Package
	Alias     string
}

func (v *VariableType) validForTypeAssert() bool {
	if v.IsPointer() == false {
		return false
	}
	if v.Typ == VARIABLE_TYPE_ARRAY && v.ArrayType.IsPrimitive() {
		return true
	}
	if v.Typ == VARIABLE_TYPE_OBJECT || v.Typ == VARIABLE_TYPE_STRING {
		return true
	}
	if v.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		if v.ArrayType.IsPointer() {
			return v.ArrayType.validForTypeAssert()
		} else {
			return true
		}
	}

	return false
}

type Map struct {
	K *VariableType
	V *VariableType
}

func (v *VariableType) mkDefaultValueExpression() *Expression {
	var e Expression
	e.IsCompileAuto = true
	e.Pos = v.Pos
	e.Value = v.Clone()
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		e.Typ = EXPRESSION_TYPE_BOOL
		e.Data = false
	case VARIABLE_TYPE_BYTE:
		e.Typ = EXPRESSION_TYPE_BYTE
		e.Data = byte(0)
	case VARIABLE_TYPE_SHORT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = int32(0)
	case VARIABLE_TYPE_INT:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = int32(0)
	case VARIABLE_TYPE_LONG:
		e.Typ = EXPRESSION_TYPE_LONG
		e.Data = int64(0)
	case VARIABLE_TYPE_FLOAT:
		e.Typ = EXPRESSION_TYPE_FLOAT
		e.Data = float32(0)
	case VARIABLE_TYPE_DOUBLE:
		e.Typ = EXPRESSION_TYPE_DOUBLE
		e.Data = float64(0)
	case VARIABLE_TYPE_STRING:
		fallthrough
	case VARIABLE_TYPE_OBJECT:
		fallthrough
	case VARIABLE_TYPE_JAVA_ARRAY:
		fallthrough
	case VARIABLE_TYPE_MAP:
		fallthrough
	case VARIABLE_TYPE_ARRAY:
		e.Typ = EXPRESSION_TYPE_NULL
	case VARIABLE_TYPE_ENUM:
		e.Typ = EXPRESSION_TYPE_INT
		e.Data = v.Enum.Enums[0].Value
	}
	return &e
}

func (v *VariableType) RightValueValid() bool {
	return v.Typ == VARIABLE_TYPE_BOOL ||
		v.IsNumber() ||
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

func (t *VariableType) resolve(block *Block, subPart ...bool) error {
	if t == nil {
		return nil
	}
	if t.Resolved {
		return nil
	}
	t.Resolved = true
	if t.Typ == VARIABLE_TYPE_T {
		if block.InheritedAttribute.Function.TypeParameters == nil ||
			block.InheritedAttribute.Function.TypeParameters[t.Name] == nil {
			return fmt.Errorf("%s typed parameter '%s' not found",
				errMsgPrefix(t.Pos), t.Name)
		}
		pos := t.Pos
		*t = *block.InheritedAttribute.Function.TypeParameters[t.Name]
		t.Pos = pos // keep pos
		return nil
	}
	if t.Typ == VARIABLE_TYPE_NAME { //
		return t.resolveName(block, len(subPart) > 0)
	}
	if t.Typ == VARIABLE_TYPE_ARRAY ||
		t.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return t.ArrayType.resolve(block, true)
	}
	if t.Typ == VARIABLE_TYPE_MAP {
		var err error
		if t.Map.K != nil {
			err = t.Map.K.resolve(block, true)
			if err != nil {
				return err
			}
		}
		if t.Map.V != nil {
			return t.Map.V.resolve(block, true)
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
		return nil, fmt.Errorf("%s type named '%s' not found", errMsgPrefix(t.Pos), t.Name)
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
	if pp, ok := p.(*Package); ok && pp != nil {
		var exists bool
		d, exists = pp.Block.NameExists(packageAndName[1])
		if exists == false {
			err = fmt.Errorf("%s '%s' not found", errMsgPrefix(t.Pos), packageAndName[1])
		}
		return d, err
	} else {
		return nil, fmt.Errorf("%s '%s' is not a package", errMsgPrefix(t.Pos), packageAndName[0])
	}

}

func (t *VariableType) mkTypeFrom(d interface{}) error {
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

func (t *VariableType) resolveName(block *Block, subPart bool) error {
	var err error
	var d interface{}

	if strings.Contains(t.Name, ".") == false {
		d, _ = block.SearchByName(t.Name)
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
	err = t.mkTypeFrom(d)
	if err != nil {
		return err
	}
	if t.Typ == VARIABLE_TYPE_ENUM && subPart {
		if t.Enum.Enums[0].Value != 0 {
			return fmt.Errorf("%s enum named '%s' as subPart of a type,first enum value named by '%s' must have value '0'",
				errMsgPrefix(t.Pos), t.Enum.Name, t.Enum.Enums[0].Name)
		}
	}
	return nil
}

func (t *VariableType) IsNumber() bool {
	return t.IsInteger() || t.IsFloat()
}

func (t *VariableType) IsPointer() bool {
	return t.Typ == VARIABLE_TYPE_OBJECT ||
		t.Typ == VARIABLE_TYPE_ARRAY ||
		t.Typ == VARIABLE_TYPE_JAVA_ARRAY ||
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
	case VARIABLE_TYPE_T:
		*ret += v.Name
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
func (v *VariableType) haveT() bool {
	if v.Typ == VARIABLE_TYPE_T {
		return true
	}
	if v.Typ == VARIABLE_TYPE_ARRAY || v.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return v.ArrayType.haveT()
	}
	if v.Typ == VARIABLE_TYPE_MAP {
		return v.Map.K.haveT() || v.Map.V.haveT()
	}
	return false
}

func (v *VariableType) canBeBindWith(t *VariableType) error {
	if v.Typ == VARIABLE_TYPE_T {
		return nil
	}
	if v.Typ == VARIABLE_TYPE_ARRAY && t.Typ == VARIABLE_TYPE_ARRAY {
		return v.ArrayType.canBeBindWith(t.ArrayType)
	}
	if v.Typ == VARIABLE_TYPE_JAVA_ARRAY && t.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return v.ArrayType.canBeBindWith(t.ArrayType)
	}
	if v.Typ == VARIABLE_TYPE_MAP && t.Typ == VARIABLE_TYPE_MAP {
		err := v.Map.K.canBeBindWith(t.Map.K)
		if err != nil {
			return err
		}
		return v.Map.V.canBeBindWith(t.Map.V)
	}
	return fmt.Errorf("cannot bind '%s' to '%s'", t.TypeString(), v.TypeString())
}

func (v *VariableType) Equal(errs *[]error, assignMent *VariableType) bool {
	if v == assignMent { // equal
		return true
	}
	if v.IsPrimitive() && assignMent.IsPrimitive() {
		return v.Typ == assignMent.Typ
	}
	if v.IsPointer() && assignMent.Typ == VARIABLE_TYPE_NULL {
		return true
	}
	if v.Typ == VARIABLE_TYPE_OBJECT && v.Class.Name == JAVA_ROOT_CLASS &&
		assignMent.IsPointer() {
		return true
	}
	if v.Typ == VARIABLE_TYPE_ARRAY && assignMent.Typ == VARIABLE_TYPE_ARRAY {
		return v.ArrayType.Equal(errs, assignMent.ArrayType)
	}
	if v.Typ == VARIABLE_TYPE_JAVA_ARRAY && assignMent.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return v.ArrayType.Equal(errs, assignMent.ArrayType)
	}

	if v.Typ == VARIABLE_TYPE_ENUM && assignMent.Typ == VARIABLE_TYPE_ENUM {
		return v.Enum.Name == assignMent.Enum.Name
	}
	if v.Typ == VARIABLE_TYPE_MAP && assignMent.Typ == VARIABLE_TYPE_MAP {
		return v.Map.K.Equal(errs, assignMent.Map.K) && v.Map.V.Equal(errs, assignMent.Map.V)
	}
	if v.Typ == VARIABLE_TYPE_OBJECT && assignMent.Typ == VARIABLE_TYPE_OBJECT { // object
		if v.Class.NotImportedYet {
			if err := v.Class.loadSelf(); err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(assignMent.Pos), err))
				return false
			}
		}
		if assignMent.Class.NotImportedYet {
			if err := assignMent.Class.loadSelf(); err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(assignMent.Pos), err))
				return false
			}
		}
		if v.Class.IsInterface() {
			i, err := assignMent.Class.implemented(v.Class.Name)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(assignMent.Pos), err))
			}
			return i
		} else { // class
			has, err := assignMent.Class.haveSuper(v.Class.Name)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(assignMent.Pos), err))
			}
			return has
		}
	}
	return false
}
