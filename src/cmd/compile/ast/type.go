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
	//VARIABLE_TYPE_FUNCTION_POINTER
	VARIABLE_TYPE_ENUM
	VARIABLE_TYPE_CLASS

	VARIABLE_TYPE_NAME
	VARIABLE_TYPE_T
	VARIABLE_TYPE_VOID

	VARIABLE_TYPE_PACKAGE
	VARIABLE_TYPE_NULL
)

type Type struct {
	TNames       []string
	Resolved     bool
	Pos          *Position
	Type         int
	Name         string
	ArrayType    *Type
	Class        *Class
	Enum         *Enum
	EnumName     *EnumName
	Function     *Function
	FunctionType *FunctionType
	Map          *Map
	Package      *Package
	Alias        string
}

func (typ *Type) validForTypeAssertOrConversion() bool {
	if typ.IsPointer() == false {
		return false
	}
	if typ.Type == VARIABLE_TYPE_ARRAY && typ.ArrayType.IsPrimitive() {
		return true
	}
	if typ.Type == VARIABLE_TYPE_OBJECT || typ.Type == VARIABLE_TYPE_STRING {
		return true
	}
	if typ.Type == VARIABLE_TYPE_JAVA_ARRAY {
		if typ.ArrayType.IsPointer() {
			return typ.ArrayType.validForTypeAssertOrConversion()
		} else {
			return true
		}
	}

	return false
}

type Map struct {
	Key   *Type
	Value *Type
}

func (typ *Type) mkDefaultValueExpression() *Expression {
	var e Expression
	e.IsCompileAuto = true
	e.Pos = typ.Pos
	e.ExpressionValue = typ.Clone()
	switch typ.Type {
	case VARIABLE_TYPE_BOOL:
		e.Type = EXPRESSION_TYPE_BOOL
		e.Data = false
	case VARIABLE_TYPE_BYTE:
		e.Type = EXPRESSION_TYPE_BYTE
		e.Data = byte(0)
	case VARIABLE_TYPE_SHORT:
		e.Type = EXPRESSION_TYPE_INT
		e.Data = int32(0)
	case VARIABLE_TYPE_INT:
		e.Type = EXPRESSION_TYPE_INT
		e.Data = int32(0)
	case VARIABLE_TYPE_LONG:
		e.Type = EXPRESSION_TYPE_LONG
		e.Data = int64(0)
	case VARIABLE_TYPE_FLOAT:
		e.Type = EXPRESSION_TYPE_FLOAT
		e.Data = float32(0)
	case VARIABLE_TYPE_DOUBLE:
		e.Type = EXPRESSION_TYPE_DOUBLE
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
		e.Type = EXPRESSION_TYPE_NULL
	case VARIABLE_TYPE_ENUM:
		e.Type = EXPRESSION_TYPE_INT
		e.Data = typ.Enum.Enums[0].Value
	}
	return &e
}

func (typ *Type) RightValueValid() bool {
	return typ.Type == VARIABLE_TYPE_BOOL ||
		typ.IsNumber() ||
		typ.Type == VARIABLE_TYPE_STRING ||
		typ.Type == VARIABLE_TYPE_OBJECT ||
		typ.Type == VARIABLE_TYPE_ARRAY ||
		typ.Type == VARIABLE_TYPE_MAP ||
		typ.Type == VARIABLE_TYPE_NULL ||
		typ.Type == VARIABLE_TYPE_JAVA_ARRAY ||
		typ.Type == VARIABLE_TYPE_ENUM ||
		typ.Type == VARIABLE_TYPE_FUNCTION
}

/*
	isTyped means can get type from this
*/
func (typ *Type) isTyped() bool {
	return typ.RightValueValid() && typ.Type != VARIABLE_TYPE_NULL
}

/*
	deep clone
*/
func (typ *Type) Clone() *Type {
	ret := &Type{}
	*ret = *typ
	if ret.Type == VARIABLE_TYPE_ARRAY ||
		ret.Type == VARIABLE_TYPE_JAVA_ARRAY {
		ret.ArrayType = typ.ArrayType.Clone()
	}
	if ret.Type == VARIABLE_TYPE_MAP {
		ret.Map = &Map{}
		ret.Map.Key = typ.Map.Key.Clone()
		ret.Map.Value = typ.Map.Value.Clone()
	}
	return ret
}

func (typ *Type) resolve(block *Block, subPart ...bool) error {
	if typ == nil {
		return nil
	}
	if typ.Resolved {
		return nil
	}
	typ.Resolved = true
	if typ.Type == VARIABLE_TYPE_T {
		if block.InheritedAttribute.Function.parameterTypes == nil ||
			block.InheritedAttribute.Function.parameterTypes[typ.Name] == nil {
			return fmt.Errorf("%s parameter type '%s' not found",
				errMsgPrefix(typ.Pos), typ.Name)
		}
		pos := typ.Pos
		*typ = *block.InheritedAttribute.Function.parameterTypes[typ.Name]
		typ.Pos = pos // keep pos
		return nil
	}
	if typ.Type == VARIABLE_TYPE_NAME { //
		return typ.resolveName(block, len(subPart) > 0)
	}
	if typ.Type == VARIABLE_TYPE_ARRAY ||
		typ.Type == VARIABLE_TYPE_JAVA_ARRAY {
		return typ.ArrayType.resolve(block, true)
	}
	if typ.Type == VARIABLE_TYPE_MAP {
		var err error
		if typ.Map.Key != nil {
			err = typ.Map.Key.resolve(block, true)
			if err != nil {
				return err
			}
		}
		if typ.Map.Value != nil {
			return typ.Map.Value.resolve(block, true)
		}
	}
	return nil
}

func (typ *Type) resolveNameFromImport() (d interface{}, err error) {
	if strings.Contains(typ.Name, ".") == false {
		i := PackageBeenCompile.getImport(typ.Pos.Filename, typ.Name)
		if i != nil {
			return PackageBeenCompile.load(i.ImportName)
		}
		return nil, fmt.Errorf("%s type named '%s' not found",
			errMsgPrefix(typ.Pos), typ.Name)
	}
	packageAndName := strings.Split(typ.Name, ".")
	i := PackageBeenCompile.getImport(typ.Pos.Filename, packageAndName[0])
	if nil == i {
		return nil, fmt.Errorf("%s package '%s' not imported",
			errMsgPrefix(typ.Pos), packageAndName[0])
	}
	p, err := PackageBeenCompile.load(i.ImportName)
	if err != nil {
		return nil, fmt.Errorf("%s %v",
			errMsgPrefix(typ.Pos), err)
	}
	if pp, ok := p.(*Package); ok && pp != nil {
		var exists bool
		d, exists = pp.Block.NameExists(packageAndName[1])
		if exists == false {
			err = fmt.Errorf("%s '%s' not found",
				errMsgPrefix(typ.Pos), packageAndName[1])
		}
		return d, err
	} else {
		return nil, fmt.Errorf("%s '%s' is not a package",
			errMsgPrefix(typ.Pos), packageAndName[0])
	}

}

func (typ *Type) makeTypeFrom(d interface{}) error {
	switch d.(type) {
	case *Class:
		dd := d.(*Class)
		if typ != nil {
			typ.Type = VARIABLE_TYPE_OBJECT
			typ.Class = dd
			return nil
		}
	case *Type:
		dd := d.(*Type)
		if dd != nil {
			tt := dd.Clone()
			tt.Pos = typ.Pos
			*typ = *tt
			return nil
		}
	case *Enum:
		dd := d.(*Enum)
		if dd != nil {
			typ.Type = VARIABLE_TYPE_ENUM
			typ.Enum = dd
			return nil
		}
	}
	return fmt.Errorf("%s name '%s' is not a type",
		errMsgPrefix(typ.Pos), typ.Name)
}

func (typ *Type) resolveName(block *Block, subPart bool) error {
	var err error
	var d interface{}
	if strings.Contains(typ.Name, ".") == false {
		d = block.searchType(typ.Name)
		loadFromImport := (d == nil)
		if loadFromImport == false { // d is not nil
			switch d.(type) {
			case *Class:
				if t := d.(*Class); t == nil {
					loadFromImport = true
				} else {
					_, loadFromImport = shouldAccessFromImports(t.Name, t.Pos, t.Pos)
				}
			case *Type:
				if t := d.(*Type); t == nil {
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
			d, err = typ.resolveNameFromImport()
			if err != nil {
				return err
			}
		}
	} else { // a.b  in type situation,must be package name
		d, err = typ.resolveNameFromImport()
		if err != nil {
			return err
		}
	}
	if d == nil {
		return fmt.Errorf("%s type named '%s' not found", errMsgPrefix(typ.Pos), typ.Name)
	}
	err = typ.makeTypeFrom(d)
	if err != nil {
		return err
	}
	if typ.Type == VARIABLE_TYPE_ENUM && subPart {
		if typ.Enum.Enums[0].Value != 0 {
			return fmt.Errorf("%s enum named '%s' as subPart of a type,first enum value named by '%s' must have value '0'",
				errMsgPrefix(typ.Pos), typ.Enum.Name, typ.Enum.Enums[0].Name)
		}
	}
	return nil
}

func (typ *Type) IsNumber() bool {
	return typ.IsInteger() || typ.IsFloat()
}

func (typ *Type) IsPointer() bool {
	return typ.Type == VARIABLE_TYPE_OBJECT ||
		typ.Type == VARIABLE_TYPE_ARRAY ||
		typ.Type == VARIABLE_TYPE_JAVA_ARRAY ||
		typ.Type == VARIABLE_TYPE_MAP ||
		typ.Type == VARIABLE_TYPE_STRING ||
		typ.Type == VARIABLE_TYPE_NULL

}

func (typ *Type) IsInteger() bool {
	return typ.Type == VARIABLE_TYPE_BYTE ||
		typ.Type == VARIABLE_TYPE_SHORT ||
		typ.Type == VARIABLE_TYPE_INT ||
		typ.Type == VARIABLE_TYPE_LONG
}

/*
	float or double
*/
func (typ *Type) IsFloat() bool {
	return typ.Type == VARIABLE_TYPE_FLOAT ||
		typ.Type == VARIABLE_TYPE_DOUBLE
}

func (typ *Type) IsPrimitive() bool {
	return typ.IsNumber() ||
		typ.Type == VARIABLE_TYPE_STRING ||
		typ.Type == VARIABLE_TYPE_BOOL
}

//可读的类型信息
func (typ *Type) typeString(ret *string) {
	if typ.Alias != "" {
		*ret += typ.Alias
		return
	}
	switch typ.Type {
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
		*ret += fmt.Sprintf("class(%s)", typ.Class.Name)
	case VARIABLE_TYPE_ENUM:
		*ret += "enum(" + typ.Enum.Name + ")"
	case VARIABLE_TYPE_ARRAY:
		*ret += "[]"
		typ.ArrayType.typeString(ret)
	case VARIABLE_TYPE_VOID:
		*ret += "void"
	case VARIABLE_TYPE_STRING:
		*ret += "string"
	case VARIABLE_TYPE_OBJECT: // class name
		*ret += "object@(" + typ.Class.Name + ")"
	case VARIABLE_TYPE_MAP:
		*ret += "map{"
		*ret += typ.Map.Key.TypeString()
		*ret += " -> "
		*ret += typ.Map.Value.TypeString()
		*ret += "}"
	case VARIABLE_TYPE_JAVA_ARRAY:
		*ret += typ.ArrayType.TypeString() + "[]"
	case VARIABLE_TYPE_PACKAGE:
		*ret += typ.Package.Name
	case VARIABLE_TYPE_NULL:
		*ret += "null"
	case VARIABLE_TYPE_NAME:
		*ret += typ.Name // resolve wrong, but typeString is ok to return
	case VARIABLE_TYPE_FUNCTION:
		*ret += typ.Function.readableMsg()
	case VARIABLE_TYPE_T:
		*ret += typ.Name
	default:
		panic(typ.Type)
	}
}

//可读的类型信息
func (typ *Type) TypeString() string {
	t := ""
	typ.typeString(&t)
	return t
}
func (typ *Type) haveParameterType() (ret []string) {
	ret = []string{}
	defer func() {
		typ.TNames = ret
	}()
	if typ.TNames != nil {
		return typ.TNames
	}
	if typ.Type == VARIABLE_TYPE_T {
		ret = []string{typ.Name}
		return
	}
	if typ.Type == VARIABLE_TYPE_ARRAY || typ.Type == VARIABLE_TYPE_JAVA_ARRAY {
		ret = typ.ArrayType.haveParameterType()
		return
	}
	if typ.Type == VARIABLE_TYPE_MAP {
		ret = []string{}
		if t := typ.Map.Key.haveParameterType(); t != nil {
			ret = append(ret, t...)
		}
		if t := typ.Map.Value.haveParameterType(); t != nil {
			ret = append(ret, t...)
		}
		return
	}
	return nil
}

func (typ *Type) canBeBindWithParameterTypes(parameterTypes map[string]*Type) error {
	if typ.Type == VARIABLE_TYPE_T {
		_, ok := parameterTypes[typ.Name]
		if ok == false {
			return fmt.Errorf("typed parameter '%s' not found", typ.Name)
		}
		return nil
	}
	if typ.Type == VARIABLE_TYPE_ARRAY || typ.Type == VARIABLE_TYPE_JAVA_ARRAY {
		return typ.ArrayType.canBeBindWithParameterTypes(parameterTypes)
	}
	if typ.Type == VARIABLE_TYPE_MAP {
		err := typ.Map.Key.canBeBindWithParameterTypes(parameterTypes)
		if err != nil {
			return err
		}
		return typ.Map.Value.canBeBindWithParameterTypes(parameterTypes)
	}
	return fmt.Errorf("not T") // looks impossible
}

/*
	if there is error,this function will crash
*/
func (typ *Type) bindWithParameterTypes(parameterTypes map[string]*Type) error {
	if typ.Type == VARIABLE_TYPE_T {
		t, ok := parameterTypes[typ.Name]
		if ok == false {
			panic(fmt.Sprintf("typed parameter '%s' not found", typ.Name))
		}
		*typ = *t.Clone() // real bind
		return nil
	}
	if typ.Type == VARIABLE_TYPE_ARRAY || typ.Type == VARIABLE_TYPE_JAVA_ARRAY {
		return typ.ArrayType.bindWithParameterTypes(parameterTypes)
	}
	if typ.Type == VARIABLE_TYPE_MAP {
		err := typ.Map.Key.bindWithParameterTypes(parameterTypes)
		if err != nil {
			return err
		}
		return typ.Map.Value.bindWithParameterTypes(parameterTypes)
	}
	panic("not T")
}

/*

 */
func (typ *Type) canBeBindWithType(mkParameterTypes map[string]*Type, bind *Type) error {
	if bind.RightValueValid() == false {
		return fmt.Errorf("'%s' is not right value valid", bind.TypeString())
	}
	if bind.Type == VARIABLE_TYPE_NULL {
		return fmt.Errorf("'%s' is un typed", bind.TypeString())
	}
	if typ.Type == VARIABLE_TYPE_T {
		mkParameterTypes[typ.Name] = bind
		return nil
	}
	if typ.Type == VARIABLE_TYPE_ARRAY && bind.Type == VARIABLE_TYPE_ARRAY {
		return typ.ArrayType.canBeBindWithType(mkParameterTypes, bind.ArrayType)
	}
	if typ.Type == VARIABLE_TYPE_JAVA_ARRAY && bind.Type == VARIABLE_TYPE_JAVA_ARRAY {
		return typ.ArrayType.canBeBindWithType(mkParameterTypes, bind.ArrayType)
	}
	if typ.Type == VARIABLE_TYPE_MAP && bind.Type == VARIABLE_TYPE_MAP {
		err := typ.Map.Key.canBeBindWithType(mkParameterTypes, bind.Map.Key)
		if err != nil {
			return err
		}
		return typ.Map.Value.canBeBindWithType(mkParameterTypes, bind.Map.Value)
	}
	return fmt.Errorf("cannot bind '%s' to '%s'", bind.TypeString(), typ.TypeString())
}

func (typ *Type) Equal(errs *[]error, compareTo *Type) bool {
	if typ == compareTo { // equal
		return true
	}
	if typ.IsPrimitive() && compareTo.IsPrimitive() {
		return typ.Type == compareTo.Type
	}
	if typ.IsPointer() && compareTo.Type == VARIABLE_TYPE_NULL {
		return true
	}
	if typ.Type == VARIABLE_TYPE_OBJECT && typ.Class.Name == JAVA_ROOT_CLASS &&
		compareTo.IsPointer() {
		return true
	}
	if typ.Type == VARIABLE_TYPE_ARRAY && compareTo.Type == VARIABLE_TYPE_ARRAY {
		return typ.ArrayType.Equal(errs, compareTo.ArrayType)
	}
	if typ.Type == VARIABLE_TYPE_JAVA_ARRAY && compareTo.Type == VARIABLE_TYPE_JAVA_ARRAY {
		return typ.ArrayType.Equal(errs, compareTo.ArrayType)
	}

	if typ.Type == VARIABLE_TYPE_ENUM && compareTo.Type == VARIABLE_TYPE_ENUM {
		return typ.Enum.Name == compareTo.Enum.Name
	}
	if typ.Type == VARIABLE_TYPE_MAP && compareTo.Type == VARIABLE_TYPE_MAP {
		return typ.Map.Key.Equal(errs, compareTo.Map.Key) && typ.Map.Value.Equal(errs, compareTo.Map.Value)
	}
	if typ.Type == VARIABLE_TYPE_OBJECT && compareTo.Type == VARIABLE_TYPE_OBJECT { // object
		if typ.Class.NotImportedYet {
			if err := typ.Class.loadSelf(); err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(compareTo.Pos), err))
				return false
			}
		}
		if compareTo.Class.NotImportedYet {
			if err := compareTo.Class.loadSelf(); err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(compareTo.Pos), err))
				return false
			}
		}
		if typ.Class.IsInterface() {
			i, err := compareTo.Class.implemented(typ.Class.Name)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(compareTo.Pos), err))
			}
			return i
		} else { // class
			has, err := compareTo.Class.haveSuper(typ.Class.Name)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(compareTo.Pos), err))
			}
			return has
		}
	}
	if typ.Type == VARIABLE_TYPE_FUNCTION && compareTo.Type == VARIABLE_TYPE_FUNCTION {
		var compareToFunctionType *FunctionType
		if compareTo.FunctionType != nil {
			compareToFunctionType = compareTo.FunctionType
		} else {
			compareToFunctionType = &compareTo.Function.Type
		}
		if len(typ.FunctionType.ParameterList) != len(compareToFunctionType.ParameterList) {
			return false
		}
		if len(typ.FunctionType.ReturnList) != len(compareToFunctionType.ReturnList) {
			return false
		}
		for k, v := range typ.FunctionType.ParameterList {
			//TODO :: force to equal or not ???
			if v.Name != compareToFunctionType.ParameterList[k].Name {
				return false
			}
			if false == v.Type.StrictEqual(compareToFunctionType.ParameterList[k].Type) {
				return false
			}
		}
		for k, v := range typ.FunctionType.ReturnList {
			//TODO ::  force to equal or not ???
			if v.Name != compareToFunctionType.ReturnList[k].Name {
				return false
			}
			if false == v.Type.StrictEqual(compareToFunctionType.ReturnList[k].Type) {
				return false
			}
		}
		return true
	}
	return false
}

func (typ *Type) StrictEqual(compareTo *Type) bool {
	if typ.Type != compareTo.Type {
		return false
	}
	if typ.IsPrimitive() {
		return typ.Type == compareTo.Type
	}
	if typ.Type == VARIABLE_TYPE_ARRAY || typ.Type == VARIABLE_TYPE_JAVA_ARRAY {
		return typ.ArrayType.StrictEqual(compareTo.ArrayType)
	}
	if typ.Type == VARIABLE_TYPE_MAP {
		if false == typ.Map.Key.StrictEqual(compareTo.Map.Key) {
			return false
		}
		return typ.Map.Value.StrictEqual(compareTo.Map.Value)
	}
	if typ.Type == VARIABLE_TYPE_ENUM {
		return typ.Enum.Name == compareTo.Enum.Name
	}
	if typ.Type == VARIABLE_TYPE_OBJECT {
		return typ.Class.Name == compareTo.Class.Name
	}
	return false
}
