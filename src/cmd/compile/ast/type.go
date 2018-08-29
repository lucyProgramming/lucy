package ast

import (
	"fmt"
	"strings"
)

type VariableTypeKind int

const (
	_ VariableTypeKind = iota
	//value types
	VariableTypeBool
	VariableTypeByte
	VariableTypeShort
	VariableTypeInt
	VariableTypeLong
	VariableTypeFloat
	VariableTypeDouble
	//ref types
	VariableTypeString
	VariableTypeObject
	VariableTypeMap
	VariableTypeArray
	VariableTypeJavaArray
	VariableTypeFunction
	VariableTypeEnum
	VariableTypeClass
	VariableTypeName
	VariableTypeTemplate
	VariableTypeVoid
	VariableTypeTypeAlias
	VariableTypePackage
	VariableTypeNull
	VariableTypeSelectGlobal
	VariableTypeMagicFunction
)

type Type struct {
	Type         VariableTypeKind
	IsBuildIn    bool // build in type alias
	IsVArgs      bool // int ...
	Resolved     bool
	Pos          *Pos
	Name         string
	Array        *Type
	Class        *Class
	Enum         *Enum
	EnumName     *EnumName // is a const
	Function     *Function
	FunctionType *FunctionType
	Map          *Map
	Package      *Package
	Alias        string
	AliasType    *Type
}

type Map struct {
	K *Type
	V *Type
}

func (typ *Type) validForTypeAssertOrConversion() bool {
	if typ.IsPointer() == false {
		return false
	}
	// object or string
	if typ.Type == VariableTypeObject || typ.Type == VariableTypeString {
		return true
	}
	if typ.Type == VariableTypeArray && typ.Array.IsPrimitive() {
		return true
	}
	if typ.Type == VariableTypeJavaArray {
		if typ.Array.IsPointer() {
			return typ.Array.validForTypeAssertOrConversion()
		} else {
			return true
		}
	}
	return false
}

func (typ *Type) mkDefaultValueExpression() *Expression {
	e := &Expression{}
	e.Description = "compilerAuto"
	e.IsCompileAuto = true
	e.Pos = typ.Pos
	e.Value = typ.Clone()
	switch typ.Type {
	case VariableTypeBool:
		e.Type = ExpressionTypeBool
		e.Data = false
	case VariableTypeByte:
		e.Type = ExpressionTypeByte
		e.Data = byte(0)
	case VariableTypeShort:
		e.Type = ExpressionTypeInt
		e.Data = int32(0)
	case VariableTypeInt:
		e.Type = ExpressionTypeInt
		e.Data = int32(0)
	case VariableTypeLong:
		e.Type = ExpressionTypeLong
		e.Data = int64(0)
	case VariableTypeFloat:
		e.Type = ExpressionTypeFloat
		e.Data = float32(0)
	case VariableTypeDouble:
		e.Type = ExpressionTypeDouble
		e.Data = float64(0)
	case VariableTypeEnum:
		e.Type = ExpressionTypeInt
		e.Data = typ.Enum.DefaultValue
	case VariableTypeFunction:
		fallthrough
	case VariableTypeString:
		fallthrough
	case VariableTypeObject:
		fallthrough
	case VariableTypeJavaArray:
		fallthrough
	case VariableTypeMap:
		fallthrough
	case VariableTypeArray:
		e.Type = ExpressionTypeNull
	}
	return e
}

func (typ *Type) RightValueValid() bool {
	return typ.Type == VariableTypeBool ||
		typ.IsNumber() ||
		typ.Type == VariableTypeString ||
		typ.Type == VariableTypeObject ||
		typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeMap ||
		typ.Type == VariableTypeNull ||
		typ.Type == VariableTypeJavaArray ||
		typ.Type == VariableTypeEnum ||
		typ.Type == VariableTypeFunction
}

/*
	have type or not
*/
func (typ *Type) isTyped() bool {
	//null is only untyped right value
	return typ.RightValueValid() &&
		typ.Type != VariableTypeNull
}

/*
	deep clone
*/
func (typ *Type) Clone() *Type {
	ret := &Type{}
	*ret = *typ
	if ret.Type == VariableTypeArray ||
		ret.Type == VariableTypeJavaArray {
		ret.Array = typ.Array.Clone()
	}
	if ret.Type == VariableTypeMap {
		ret.Map = &Map{}
		ret.Map.K = typ.Map.K.Clone()
		ret.Map.V = typ.Map.V.Clone()
	}
	return ret
}

func (typ *Type) resolve(block *Block) error {
	if typ == nil {
		return nil
	}
	if typ.Resolved {
		return nil
	}
	defer func() {
		typ.Resolved = true
	}()
	if typ.Type == VariableTypeTemplate {
		if block.InheritedAttribute.Function.parameterTypes == nil {
			return fmt.Errorf("%s parameter type '%s' not in a template function",
				errMsgPrefix(typ.Pos), typ.Name)
		}
		if block.InheritedAttribute.Function.parameterTypes[typ.Name] == nil {
			return fmt.Errorf("%s parameter type '%s' not found",
				errMsgPrefix(typ.Pos), typ.Name)
		}
		pos := typ.Pos // keep pos
		*typ = *block.InheritedAttribute.Function.parameterTypes[typ.Name]
		typ.Pos = pos // keep pos
		return nil
	}
	if typ.Type == VariableTypeName { //
		return typ.resolveName(block)
	}
	if typ.Type == VariableTypeSelectGlobal {
		d, exists := PackageBeenCompile.Block.NameExists(typ.Name)
		if exists == false {
			return fmt.Errorf("%s '%s' not found",
				errMsgPrefix(typ.Pos), typ.Name)
		}
		switch d.(type) {
		case *Class:
			typ.Type = VariableTypeObject
			typ.Class = d.(*Class)
		case *Enum:
			typ.Type = VariableTypeEnum
			typ.Enum = d.(*Enum)
		case *Type:
			pos := typ.Pos
			*typ = *d.(*Type)
			typ.Pos = pos
		default:
			return fmt.Errorf("%s '%s' is not a type",
				errMsgPrefix(typ.Pos), typ.Name)
		}
	}
	if typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeJavaArray {
		return typ.Array.resolve(block)
	}
	if typ.Type == VariableTypeMap {
		var err error
		if typ.Map.K != nil {
			err = typ.Map.K.resolve(block)
			if err != nil {
				return err
			}
		}
		if typ.Map.V != nil {
			return typ.Map.V.resolve(block)
		}
	}
	if typ.Type == VariableTypeFunction {
		for _, v := range typ.FunctionType.ParameterList {
			if err := v.Type.resolve(block); err != nil {
				return err
			}
		}
		for _, v := range typ.FunctionType.ReturnList {
			if err := v.Type.resolve(block); err != nil {
				return err
			}
		}
	}
	return nil
}

func (typ *Type) resolveNameFromImport() (d interface{}, err error) {
	if strings.Contains(typ.Name, ".") == false {
		i := PackageBeenCompile.getImport(typ.Pos.Filename, typ.Name)
		if i != nil {
			i.Used = true
			return PackageBeenCompile.load(i.Import)
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
	i.Used = true
	p, err := PackageBeenCompile.load(i.Import)
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

func (typ *Type) makeTypeFrom(d interface{}, loadFromImport bool) error {
	switch d.(type) {
	case *Class:
		dd := d.(*Class)
		if loadFromImport && dd.IsPublic() == false {
			PackageBeenCompile.Errors = append(PackageBeenCompile.Errors, fmt.Errorf("%s class '%s' is not public",
				errMsgPrefix(typ.Pos), dd.Name))
		}
		if typ != nil {
			typ.Type = VariableTypeObject
			typ.Class = dd
			return nil
		}
	case *Type:
		dd := d.(*Type)
		tt := dd.Clone()
		tt.Pos = typ.Pos
		*typ = *tt
		return nil
	case *Enum:
		dd := d.(*Enum)
		if loadFromImport && dd.IsPublic() == false {
			PackageBeenCompile.Errors = append(PackageBeenCompile.Errors, fmt.Errorf("%s enum '%s' is not public",
				errMsgPrefix(typ.Pos), dd.Name))
		}
		typ.Type = VariableTypeEnum
		typ.Enum = dd
		return nil
	}
	return fmt.Errorf("%s name '%s' is not a type",
		errMsgPrefix(typ.Pos), typ.Name)
}

func (typ *Type) resolveName(block *Block) error {
	var err error
	var d interface{}
	var loadFromImport bool
	if strings.Contains(typ.Name, ".") == false {
		d = block.searchType(typ.Name)
		if d != nil {
			if loadFromImport == false { // d is not nil
				switch d.(type) {
				case *Class:
					if t := d.(*Class); t != nil && t.IsBuildIn {
						loadFromImport = false
					} else {
						_, loadFromImport = shouldAccessFromImports(typ.Name, typ.Pos, t.Pos)
					}
				case *Type:
					if t := d.(*Type); t != nil && t.IsBuildIn {
						loadFromImport = false
					} else {
						_, loadFromImport = shouldAccessFromImports(typ.Name, typ.Pos, t.Pos)
					}
				case *Enum:
					if t := d.(*Enum); t != nil && t.IsBuildIn {
						loadFromImport = true
					} else {
						_, loadFromImport = shouldAccessFromImports(typ.Name, typ.Pos, t.Pos)
					}
				}
			}
		} else {
			loadFromImport = true
		}
		if loadFromImport {
			d, err = typ.resolveNameFromImport()
			if err != nil {
				return err
			}
		}
	} else { // a.b  in type situation,must be package name
		loadFromImport = true
		d, err = typ.resolveNameFromImport()
		if err != nil {
			return err
		}
	}
	if d == nil {
		return fmt.Errorf("%s type named '%s' not found", errMsgPrefix(typ.Pos), typ.Name)
	}
	err = typ.makeTypeFrom(d, loadFromImport)
	if err != nil {
		return err
	}
	return nil
}

func (typ *Type) IsNumber() bool {
	return typ.IsInteger() || typ.IsFloat()
}

func (typ *Type) IsPointer() bool {
	return typ.Type == VariableTypeObject ||
		typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeJavaArray ||
		typ.Type == VariableTypeMap ||
		typ.Type == VariableTypeString ||
		typ.Type == VariableTypeNull ||
		typ.Type == VariableTypeFunction

}

func (typ *Type) IsInteger() bool {
	return typ.Type == VariableTypeByte ||
		typ.Type == VariableTypeShort ||
		typ.Type == VariableTypeInt ||
		typ.Type == VariableTypeLong
}

/*
	float or double
*/
func (typ *Type) IsFloat() bool {
	return typ.Type == VariableTypeFloat ||
		typ.Type == VariableTypeDouble
}

func (typ *Type) IsPrimitive() bool {
	return typ.IsNumber() ||
		typ.Type == VariableTypeString ||
		typ.Type == VariableTypeBool
}

//可读的类型信息
func (typ *Type) typeString(ret *string) {
	if typ.Alias != "" {
		*ret += typ.Alias
		return
	}
	switch typ.Type {
	case VariableTypeBool:
		*ret += "bool"
	case VariableTypeByte:
		*ret += "byte"
	case VariableTypeShort:
		*ret += "short"
	case VariableTypeInt:
		*ret += "int"
	case VariableTypeLong:
		*ret += "long"
	case VariableTypeFloat:
		*ret += "float"
	case VariableTypeDouble:
		*ret += "double"
	case VariableTypeString:
		*ret += "string"
	case VariableTypeObject: // class name
		*ret += "object@(" + typ.Class.Name + ")"
	case VariableTypeMap:
		*ret += "map{"
		*ret += typ.Map.K.TypeString()
		*ret += " -> "
		*ret += typ.Map.V.TypeString()
		*ret += "}"
	case VariableTypeArray:
		*ret += "[]"
		typ.Array.typeString(ret)
	case VariableTypeJavaArray:
		if typ.IsVArgs {
			*ret += typ.Array.TypeString() + "..."
		} else {
			*ret += typ.Array.TypeString() + "[]"
		}
	case VariableTypeFunction:
		s := "fn ("
		for k, v := range typ.FunctionType.ParameterList {
			if v.Name != "" {
				s += v.Name + " "
			}
			s += v.Type.TypeString()
			if k != len(typ.FunctionType.ParameterList)-1 {
				s += " , "
			}
		}
		if typ.FunctionType.VArgs != nil {
			if len(typ.FunctionType.ParameterList) > 0 {
				s += ","
			}
			s += typ.FunctionType.VArgs.Name + " "
			s += typ.FunctionType.VArgs.Type.TypeString()
		}
		s += ")"
		if len(typ.FunctionType.ReturnList) > 0 {
			s += " -> ("
			for k, v := range typ.FunctionType.ReturnList {
				if v.Name != "" {
					s += v.Name + " "
				}
				s += v.Type.TypeString()
				if k != len(typ.FunctionType.ReturnList)-1 {
					s += ","
				}
			}
			s += ")"
		}
		*ret += s
	case VariableTypeEnum:
		*ret += "enum(" + typ.Enum.Name + ")"
	case VariableTypeClass:
		*ret += fmt.Sprintf("class(%s)", typ.Class.Name)
	case VariableTypeName:
		*ret += typ.Name // resolve wrong, but typeString is ok to return
	case VariableTypeTemplate:
		*ret += typ.Name

	case VariableTypeVoid:
		*ret += "void"
	case VariableTypeTypeAlias:
		*ret += "type_alias@" + typ.AliasType.TypeString()
	case VariableTypePackage:
		*ret += "package@" + typ.Package.Name
	case VariableTypeNull:
		*ret += "null"
	case VariableTypeSelectGlobal:
		*ret += typ.Name
	case VariableTypeMagicFunction:
		*ret = MagicIdentifierFunction
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
func (typ *Type) getParameterType() []string {
	if typ.Type == VariableTypeTemplate {
		return []string{typ.Name}
	}
	if typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeJavaArray {
		return typ.Array.getParameterType()
	}
	if typ.Type == VariableTypeMap {
		ret := []string{}
		if t := typ.Map.K.getParameterType(); t != nil {
			ret = append(ret, t...)
		}
		if t := typ.Map.V.getParameterType(); t != nil {
			ret = append(ret, t...)
		}
		return ret
	}
	return nil
}

func (typ *Type) canBeBindWithParameterTypes(parameterTypes map[string]*Type) error {
	if typ.Type == VariableTypeTemplate {
		_, ok := parameterTypes[typ.Name]
		if ok == false {
			return fmt.Errorf("typed parameter '%s' not found", typ.Name)
		}
		return nil
	}
	if typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeJavaArray {
		return typ.Array.canBeBindWithParameterTypes(parameterTypes)
	}
	if typ.Type == VariableTypeMap {
		err := typ.Map.K.canBeBindWithParameterTypes(parameterTypes)
		if err != nil {
			return err
		}
		return typ.Map.V.canBeBindWithParameterTypes(parameterTypes)
	}
	return nil
	//return fmt.Errorf("not T") // looks impossible
}

/*
	if there is error,this function will crash
*/
func (typ *Type) bindWithParameterTypes(parameterTypes map[string]*Type) error {
	if typ.Type == VariableTypeTemplate {
		t, ok := parameterTypes[typ.Name]
		if ok == false {
			panic(fmt.Sprintf("typed parameter '%s' not found", typ.Name))
		}
		*typ = *t.Clone() // real bind
		return nil
	}
	if typ.Type == VariableTypeArray || typ.Type == VariableTypeJavaArray {
		return typ.Array.bindWithParameterTypes(parameterTypes)
	}
	if typ.Type == VariableTypeMap {
		if len(typ.Map.K.getParameterType()) > 0 {
			err := typ.Map.K.bindWithParameterTypes(parameterTypes)
			if err != nil {
				return err
			}
		}
		if len(typ.Map.V.getParameterType()) > 0 {
			return typ.Map.V.bindWithParameterTypes(parameterTypes)
		}
	}
	panic("not T")
}

/*

 */
func (typ *Type) canBeBindWithType(mkParameterTypes map[string]*Type, bind *Type) error {
	if bind.RightValueValid() == false {
		return fmt.Errorf("'%s' is not right value valid", bind.TypeString())
	}
	if bind.Type == VariableTypeNull {
		return fmt.Errorf("'%s' is un typed", bind.TypeString())
	}
	if typ.Type == VariableTypeTemplate {
		mkParameterTypes[typ.Name] = bind
		return nil
	}
	if typ.Type == VariableTypeArray && bind.Type == VariableTypeArray {
		return typ.Array.canBeBindWithType(mkParameterTypes, bind.Array)
	}
	if typ.Type == VariableTypeJavaArray && bind.Type == VariableTypeJavaArray {
		return typ.Array.canBeBindWithType(mkParameterTypes, bind.Array)
	}
	if typ.Type == VariableTypeMap && bind.Type == VariableTypeMap {
		if len(typ.Map.K.getParameterType()) > 0 {
			err := typ.Map.K.canBeBindWithType(mkParameterTypes, bind.Map.K)
			if err != nil {
				return err
			}
		}
		if len(typ.Map.V.getParameterType()) > 0 {
			return typ.Map.V.canBeBindWithType(mkParameterTypes, bind.Map.V)
		}
	}
	return fmt.Errorf("cannot bind '%s' to '%s'", bind.TypeString(), typ.TypeString())
}

func (leftValue *Type) Equal(errs *[]error, rightValue *Type) bool {
	if leftValue == rightValue { // equal
		return true
	}
	if leftValue.IsPrimitive() && rightValue.IsPrimitive() {
		return leftValue.Type == rightValue.Type
	}
	if leftValue.IsPointer() && rightValue.Type == VariableTypeNull {
		return true
	}
	if leftValue.Type == VariableTypeObject && leftValue.Class.Name == JavaRootClass &&
		rightValue.IsPointer() {
		return true
	}
	if leftValue.Type == VariableTypeArray && rightValue.Type == VariableTypeArray {
		return leftValue.Array.Equal(errs, rightValue.Array)
	}
	if leftValue.Type == VariableTypeJavaArray && rightValue.Type == VariableTypeJavaArray {
		if leftValue.IsVArgs != rightValue.IsVArgs {
			return false
		}
		return leftValue.Array.Equal(errs, rightValue.Array)
	}

	if leftValue.Type == VariableTypeEnum && rightValue.Type == VariableTypeEnum {
		return leftValue.Enum.Name == rightValue.Enum.Name
	}
	if leftValue.Type == VariableTypeMap && rightValue.Type == VariableTypeMap {
		return leftValue.Map.K.Equal(errs, rightValue.Map.K) && leftValue.Map.V.Equal(errs, rightValue.Map.V)
	}
	if leftValue.Type == VariableTypeObject && rightValue.Type == VariableTypeObject { // object
		if leftValue.Class.NotImportedYet {
			if err := leftValue.Class.loadSelf(leftValue.Pos); err != nil {
				*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(rightValue.Pos), err))
				return false
			}
		}
		if rightValue.Class.NotImportedYet {
			if err := rightValue.Class.loadSelf(rightValue.Pos); err != nil {
				*errs = append(*errs, err)
				return false
			}
		}
		if leftValue.Class.IsInterface() {
			i, err := rightValue.Class.implementedInterface(leftValue.Pos, leftValue.Class.Name)
			if err != nil {
				*errs = append(*errs, err)
			}
			return i
		} else { // class
			has, err := rightValue.Class.haveSuperClass(rightValue.Pos, leftValue.Class.Name)
			if err != nil {
				*errs = append(*errs, err)
			}
			return has
		}
	}
	if leftValue.Type == VariableTypeFunction && rightValue.Type == VariableTypeFunction {
		return leftValue.FunctionType.equal(rightValue.FunctionType)
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
	if typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeJavaArray {
		if typ.Type == VariableTypeJavaArray {
			if typ.IsVArgs != compareTo.IsVArgs {
				return false
			}
		}
		return typ.Array.StrictEqual(compareTo.Array)
	}
	if typ.Type == VariableTypeMap {
		if false == typ.Map.K.StrictEqual(compareTo.Map.K) {
			return false
		}
		return typ.Map.V.StrictEqual(compareTo.Map.V)
	}
	if typ.Type == VariableTypeEnum {
		return typ.Enum.Name == compareTo.Enum.Name
	}
	if typ.Type == VariableTypeObject {
		return typ.Class.Name == compareTo.Class.Name
	}
	if typ.Type == VariableTypeVoid {
		return true
	}
	if typ.Type == VariableTypeFunction {
		return typ.FunctionType.equal(compareTo.FunctionType)
	}
	return false
}
