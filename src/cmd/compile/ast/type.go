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
	haveTCalled bool
	TNames      []string
	Resolved    bool
	Pos         *Pos
	Typ         int
	Name        string
	ArrayType   *VariableType
	Class       *Class
	Enum        *Enum
	EnumName    *EnumName
	Function    *Function
	Map         *Map
	Package     *Package
	Alias       string
}

func (variableType *VariableType) validForTypeAssert() bool {
	if variableType.IsPointer() == false {
		return false
	}
	if variableType.Typ == VARIABLE_TYPE_ARRAY && variableType.ArrayType.IsPrimitive() {
		return true
	}
	if variableType.Typ == VARIABLE_TYPE_OBJECT || variableType.Typ == VARIABLE_TYPE_STRING {
		return true
	}
	if variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		if variableType.ArrayType.IsPointer() {
			return variableType.ArrayType.validForTypeAssert()
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

func (variableType *VariableType) mkDefaultValueExpression() *Expression {
	var e Expression
	e.IsCompileAuto = true
	e.Pos = variableType.Pos
	e.Value = variableType.Clone()
	switch variableType.Typ {
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
		e.Data = variableType.Enum.Enums[0].Value
	}
	return &e
}

func (variableType *VariableType) RightValueValid() bool {
	return variableType.Typ == VARIABLE_TYPE_BOOL ||
		variableType.IsNumber() ||
		variableType.Typ == VARIABLE_TYPE_STRING ||
		variableType.Typ == VARIABLE_TYPE_OBJECT ||
		variableType.Typ == VARIABLE_TYPE_ARRAY ||
		variableType.Typ == VARIABLE_TYPE_MAP ||
		variableType.Typ == VARIABLE_TYPE_NULL ||
		variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY ||
		variableType.Typ == VARIABLE_TYPE_ENUM
}

/*
	isTyped means can get type from this
*/
func (variableType *VariableType) isTyped() bool {
	return variableType.RightValueValid() && variableType.Typ != VARIABLE_TYPE_NULL
}

/*
	deep clone
*/
func (variableType *VariableType) Clone() *VariableType {
	ret := &VariableType{}
	*ret = *variableType
	if ret.Typ == VARIABLE_TYPE_ARRAY ||
		ret.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		ret.ArrayType = variableType.ArrayType.Clone()
	}
	if ret.Typ == VARIABLE_TYPE_MAP {
		ret.Map = &Map{}
		ret.Map.K = variableType.Map.K.Clone()
		ret.Map.V = variableType.Map.V.Clone()
	}
	return ret
}

func (variableType *VariableType) resolve(block *Block, subPart ...bool) error {
	if variableType == nil {
		return nil
	}
	if variableType.Resolved {
		return nil
	}
	variableType.Resolved = true
	if variableType.Typ == VARIABLE_TYPE_T {
		if block.InheritedAttribute.Function.TypeParameters == nil ||
			block.InheritedAttribute.Function.TypeParameters[variableType.Name] == nil {
			return fmt.Errorf("%s typed parameter '%s' not found",
				errMsgPrefix(variableType.Pos), variableType.Name)
		}
		pos := variableType.Pos
		*variableType = *block.InheritedAttribute.Function.TypeParameters[variableType.Name]
		variableType.Pos = pos // keep pos
		return nil
	}
	if variableType.Typ == VARIABLE_TYPE_NAME { //
		return variableType.resolveName(block, len(subPart) > 0)
	}
	if variableType.Typ == VARIABLE_TYPE_ARRAY ||
		variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return variableType.ArrayType.resolve(block, true)
	}
	if variableType.Typ == VARIABLE_TYPE_MAP {
		var err error
		if variableType.Map.K != nil {
			err = variableType.Map.K.resolve(block, true)
			if err != nil {
				return err
			}
		}
		if variableType.Map.V != nil {
			return variableType.Map.V.resolve(block, true)
		}
	}
	return nil
}

func (variableType *VariableType) resolveNameFromImport() (d interface{}, err error) {
	if strings.Contains(variableType.Name, ".") == false {
		i := PackageBeenCompile.getImport(variableType.Pos.Filename, variableType.Name)
		if i != nil {
			return PackageBeenCompile.load(i.Resource)
		}
		return nil, fmt.Errorf("%s type named '%s' not found", errMsgPrefix(variableType.Pos), variableType.Name)
	}
	packageAndName := strings.Split(variableType.Name, ".")
	i := PackageBeenCompile.getImport(variableType.Pos.Filename, packageAndName[0])
	if nil == i {
		return nil, fmt.Errorf("%s package '%s' not imported", errMsgPrefix(variableType.Pos), packageAndName[0])
	}
	p, err := PackageBeenCompile.load(i.Resource)
	if err != nil {
		return nil, fmt.Errorf("%s %v", errMsgPrefix(variableType.Pos), err)
	}
	if pp, ok := p.(*Package); ok && pp != nil {
		var exists bool
		d, exists = pp.Block.NameExists(packageAndName[1])
		if exists == false {
			err = fmt.Errorf("%s '%s' not found", errMsgPrefix(variableType.Pos), packageAndName[1])
		}
		return d, err
	} else {
		return nil, fmt.Errorf("%s '%s' is not a package", errMsgPrefix(variableType.Pos), packageAndName[0])
	}

}

func (variableType *VariableType) mkTypeFrom(d interface{}) error {
	switch d.(type) {
	case *Class:
		dd := d.(*Class)
		if variableType != nil {
			variableType.Typ = VARIABLE_TYPE_OBJECT
			variableType.Class = dd
			return nil
		}
	case *VariableType:
		dd := d.(*VariableType)
		if dd != nil {
			tt := dd.Clone()
			tt.Pos = variableType.Pos
			*variableType = *tt
			return nil
		}
	case *Enum:
		dd := d.(*Enum)
		if dd != nil {
			variableType.Typ = VARIABLE_TYPE_ENUM
			variableType.Enum = dd
			return nil
		}
	}
	return fmt.Errorf("%s name '%s' is not a type", errMsgPrefix(variableType.Pos), variableType.Name)
}

func (variableType *VariableType) resolveName(block *Block, subPart bool) error {
	var err error
	var d interface{}
	if strings.Contains(variableType.Name, ".") == false {
		d = block.searchType(variableType.Name)
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
			d, err = variableType.resolveNameFromImport()
			if err != nil {
				return err
			}
		}
	} else { // a.b  in type situation,must be package name
		d, err = variableType.resolveNameFromImport()
		if err != nil {
			return err
		}
	}
	if d == nil {
		return fmt.Errorf("%s type named '%s' not found", errMsgPrefix(variableType.Pos), variableType.Name)
	}
	err = variableType.mkTypeFrom(d)
	if err != nil {
		return err
	}
	if variableType.Typ == VARIABLE_TYPE_ENUM && subPart {
		if variableType.Enum.Enums[0].Value != 0 {
			return fmt.Errorf("%s enum named '%s' as subPart of a type,first enum value named by '%s' must have value '0'",
				errMsgPrefix(variableType.Pos), variableType.Enum.Name, variableType.Enum.Enums[0].Name)
		}
	}
	return nil
}

func (variableType *VariableType) IsNumber() bool {
	return variableType.IsInteger() || variableType.IsFloat()
}

func (variableType *VariableType) IsPointer() bool {
	return variableType.Typ == VARIABLE_TYPE_OBJECT ||
		variableType.Typ == VARIABLE_TYPE_ARRAY ||
		variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY ||
		variableType.Typ == VARIABLE_TYPE_MAP ||
		variableType.Typ == VARIABLE_TYPE_STRING

}

func (variableType *VariableType) IsInteger() bool {
	return variableType.Typ == VARIABLE_TYPE_BYTE ||
		variableType.Typ == VARIABLE_TYPE_SHORT ||
		variableType.Typ == VARIABLE_TYPE_INT ||
		variableType.Typ == VARIABLE_TYPE_LONG
}

/*
	float or double
*/
func (variableType *VariableType) IsFloat() bool {
	return variableType.Typ == VARIABLE_TYPE_FLOAT ||
		variableType.Typ == VARIABLE_TYPE_DOUBLE
}

func (variableType *VariableType) IsPrimitive() bool {
	return variableType.IsNumber() ||
		variableType.Typ == VARIABLE_TYPE_STRING ||
		variableType.Typ == VARIABLE_TYPE_BOOL
}

//可读的类型信息
func (variableType *VariableType) typeString(ret *string) {
	if variableType.Alias != "" {
		*ret += variableType.Alias
		return
	}
	switch variableType.Typ {
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
		*ret += fmt.Sprintf("class(%s)", variableType.Class.Name)
	case VARIABLE_TYPE_ENUM:
		*ret += "enum(" + variableType.Enum.Name + ")"
	case VARIABLE_TYPE_ARRAY:
		*ret += "[]"
		variableType.ArrayType.typeString(ret)
	case VARIABLE_TYPE_VOID:
		*ret += "void"
	case VARIABLE_TYPE_STRING:
		*ret += "string"
	case VARIABLE_TYPE_OBJECT: // class name
		*ret += "object@(" + variableType.Class.Name + ")"
	case VARIABLE_TYPE_MAP:
		*ret += "map{"
		*ret += variableType.Map.K.TypeString()
		*ret += " -> "
		*ret += variableType.Map.V.TypeString()
		*ret += "}"
	case VARIABLE_TYPE_JAVA_ARRAY:
		*ret += variableType.ArrayType.TypeString() + "[]"
	case VARIABLE_TYPE_PACKAGE:
		*ret += variableType.Package.Name
	case VARIABLE_TYPE_NULL:
		*ret += "null"
	case VARIABLE_TYPE_NAME:
		*ret += variableType.Name // resolve wrong, but typeString is ok to return
	case VARIABLE_TYPE_FUNCTION:
		*ret += variableType.Function.readableMsg()
	case VARIABLE_TYPE_T:
		*ret += variableType.Name
	default:
		panic(variableType.Typ)
	}
}

//可读的类型信息
func (variableType *VariableType) TypeString() string {
	t := ""
	variableType.typeString(&t)
	return t
}
func (variableType *VariableType) haveT() (ret []string) {
	defer func() {
		variableType.haveTCalled = true
		variableType.TNames = ret
	}()
	if variableType.Typ == VARIABLE_TYPE_T {
		ret = []string{variableType.Name}
		return
	}
	if variableType.Typ == VARIABLE_TYPE_ARRAY || variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		ret = variableType.ArrayType.haveT()
		return
	}
	if variableType.Typ == VARIABLE_TYPE_MAP {
		ret = []string{}
		if t := variableType.Map.K.haveT(); t != nil {
			ret = append(ret, t...)
		}
		if t := variableType.Map.V.haveT(); t != nil {
			ret = append(ret, t...)
		}
		return
	}
	return nil
}

func (variableType *VariableType) canBeBindWithTypedParameters(typedParaMeters map[string]*VariableType) error {
	if variableType.Typ == VARIABLE_TYPE_T {
		_, ok := typedParaMeters[variableType.Name]
		if ok == false {
			return fmt.Errorf("typed parameter '%s' not found", variableType.Name)
		}
		return nil
	}
	if variableType.Typ == VARIABLE_TYPE_ARRAY || variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return variableType.ArrayType.canBeBindWithTypedParameters(typedParaMeters)
	}
	if variableType.Typ == VARIABLE_TYPE_MAP {
		err := variableType.Map.K.canBeBindWithTypedParameters(typedParaMeters)
		if err != nil {
			return err
		}
		return variableType.Map.V.canBeBindWithTypedParameters(typedParaMeters)
	}
	return fmt.Errorf("not T") // looks impossible
}

/*
	if there is error,this function will crash
*/
func (variableType *VariableType) bindWithTypedParameters(typedParaMeters map[string]*VariableType) error {
	if variableType.Typ == VARIABLE_TYPE_T {
		t, ok := typedParaMeters[variableType.Name]
		if ok == false {
			panic(fmt.Sprintf("typed parameter '%s' not found", variableType.Name))
		}
		*variableType = *t.Clone() // real bind
		return nil
	}
	if variableType.Typ == VARIABLE_TYPE_ARRAY {
		return variableType.ArrayType.bindWithTypedParameters(typedParaMeters)
	}
	if variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return variableType.ArrayType.bindWithTypedParameters(typedParaMeters)
	}
	if variableType.Typ == VARIABLE_TYPE_MAP {
		err := variableType.Map.K.bindWithTypedParameters(typedParaMeters)
		if err != nil {
			return err
		}
		return variableType.Map.V.bindWithTypedParameters(typedParaMeters)
	}
	panic("not T")
}

/*

 */
func (variableType *VariableType) canBebindWithType(typedParaMeters map[string]*VariableType, t *VariableType) error {
	if t.RightValueValid() == false {
		return fmt.Errorf("'%s' is not right value valid", t.TypeString())
	}
	if t.Typ == VARIABLE_TYPE_NULL {
		return fmt.Errorf("'%s' is un typed", t.TypeString())
	}
	if variableType.Typ == VARIABLE_TYPE_T {
		typedParaMeters[variableType.Name] = t
		return nil
	}
	if variableType.Typ == VARIABLE_TYPE_ARRAY && t.Typ == VARIABLE_TYPE_ARRAY {
		return variableType.ArrayType.canBebindWithType(typedParaMeters, t.ArrayType)
	}
	if variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY && t.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return variableType.ArrayType.canBebindWithType(typedParaMeters, t.ArrayType)
	}
	if variableType.Typ == VARIABLE_TYPE_MAP && t.Typ == VARIABLE_TYPE_MAP {
		err := variableType.Map.K.canBebindWithType(typedParaMeters, t.Map.K)
		if err != nil {
			return err
		}
		return variableType.Map.V.canBebindWithType(typedParaMeters, t.Map.V)
	}
	return fmt.Errorf("cannot bind '%s' to '%s'", t.TypeString(), variableType.TypeString())
}

func (variableType *VariableType) Equal(errs *[]error, assignMent *VariableType) bool {
	if variableType == assignMent { // equal
		return true
	}
	if variableType.IsPrimitive() && assignMent.IsPrimitive() {
		return variableType.Typ == assignMent.Typ
	}
	if variableType.IsPointer() && assignMent.Typ == VARIABLE_TYPE_NULL {
		return true
	}
	if variableType.Typ == VARIABLE_TYPE_OBJECT && variableType.Class.Name == JAVA_ROOT_CLASS &&
		assignMent.IsPointer() {
		return true
	}
	if variableType.Typ == VARIABLE_TYPE_ARRAY && assignMent.Typ == VARIABLE_TYPE_ARRAY {
		return variableType.ArrayType.Equal(errs, assignMent.ArrayType)
	}
	if variableType.Typ == VARIABLE_TYPE_JAVA_ARRAY && assignMent.Typ == VARIABLE_TYPE_JAVA_ARRAY {
		return variableType.ArrayType.Equal(errs, assignMent.ArrayType)
	}

	if variableType.Typ == VARIABLE_TYPE_ENUM && assignMent.Typ == VARIABLE_TYPE_ENUM {
		return variableType.Enum.Name == assignMent.Enum.Name
	}
	if variableType.Typ == VARIABLE_TYPE_MAP && assignMent.Typ == VARIABLE_TYPE_MAP {
		return variableType.Map.K.Equal(errs, assignMent.Map.K) && variableType.Map.V.Equal(errs, assignMent.Map.V)
	}
	if variableType.Typ == VARIABLE_TYPE_OBJECT && assignMent.Typ == VARIABLE_TYPE_OBJECT { // object
		if variableType.Class.NotImportedYet {
			if err := variableType.Class.loadSelf(); err != nil {
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
		if variableType.Class.IsInterface() {
			i, err := assignMent.Class.implemented(variableType.Class.Name)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(assignMent.Pos), err))
			}
			return i
		} else { // class
			has, err := assignMent.Class.haveSuper(variableType.Class.Name)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(assignMent.Pos), err))
			}
			return has
		}
	}
	return false
}
