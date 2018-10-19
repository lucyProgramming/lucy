package ast

import (
	"fmt"
	"strings"
)

type VariableTypeKind int

const (
	_ VariableTypeKind = iota
	//primitive types
	VariableTypeBool
	VariableTypeByte
	VariableTypeShort
	VariableTypeChar
	VariableTypeInt
	VariableTypeLong
	VariableTypeFloat
	VariableTypeDouble
	//enum
	VariableTypeEnum
	//ref types
	VariableTypeString
	VariableTypeObject
	VariableTypeMap
	VariableTypeArray
	VariableTypeJavaArray
	VariableTypeFunction
	VariableTypeClass
	VariableTypeName
	VariableTypeTemplate
	VariableTypeVoid

	VariableTypePackage
	VariableTypeNull
	VariableTypeGlobal
	VariableTypeMagicFunction
	VariableTypeDynamicSelector //
)

type Type struct {
	Type           VariableTypeKind
	IsBuildIn      bool // build in type alias
	IsVariableArgs bool // int ...
	Resolved       bool
	Pos            *Pos
	Name           string
	Array          *Type
	Class          *Class
	Enum           *Enum
	EnumName       *EnumName // indicate a const
	Function       *Function
	FunctionType   *FunctionType
	Map            *Map
	Package        *Package
	Alias          *TypeAlias
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
	case VariableTypeChar:
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
	default:
		e.Type = ExpressionTypeNull
	}
	return e
}

func (typ *Type) rightValueValid() error {
	if typ.Type == VariableTypeBool ||
		typ.IsNumber() ||
		typ.Type == VariableTypeString ||
		typ.Type == VariableTypeObject ||
		typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeMap ||
		typ.Type == VariableTypeNull ||
		typ.Type == VariableTypeJavaArray ||
		typ.Type == VariableTypeEnum ||
		typ.Type == VariableTypeFunction {
		return nil
	}
	switch typ.Type {
	case VariableTypePackage:
		return fmt.Errorf("%s use package '%s' without selector",
			typ.Pos.ErrMsgPrefix(), typ.Package.Name)
	case VariableTypeClass:
		return fmt.Errorf("%s use class '%s' without selector",
			typ.Pos.ErrMsgPrefix(), typ.Class.Name)
	case VariableTypeMagicFunction:
		return fmt.Errorf("%s use '%s' without selector",
			typ.Pos.ErrMsgPrefix(), magicIdentifierFunction)
	default:
		return fmt.Errorf("%s '%s' is not right value valid",
			typ.Pos.ErrMsgPrefix(), typ.TypeString())
	}
}

/*
	have type or not
*/
func (typ *Type) isTyped() error {
	if err := typ.rightValueValid(); err != nil {
		return err
	}
	/*
		null is only untyped right value
	*/
	if typ.Type == VariableTypeNull {
		return fmt.Errorf("%s '%s' is not typed",
			typ.Pos.ErrMsgPrefix(), typ.TypeString())
	}
	return nil
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
	if typ.Type == VariableTypeFunction {
		ret.FunctionType = typ.FunctionType.Clone()
	}
	return ret
}

func (typ *Type) resolve(block *Block) error {
	if typ.Resolved {
		return nil
	}
	typ.Resolved = true // single threading
	switch typ.Type {
	case VariableTypeTemplate:
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
	case VariableTypeName:
		return typ.resolveName(block)
	case VariableTypeGlobal:
		d, exists := PackageBeenCompile.Block.NameExists(typ.Name)
		if exists == false {
			return fmt.Errorf("%s '%s' not found",
				errMsgPrefix(typ.Pos), typ.Name)
		}
		return typ.makeTypeFrom(d)
	case VariableTypeArray, VariableTypeJavaArray:
		return typ.Array.resolve(block)
	case VariableTypeMap:
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
	case VariableTypeFunction:
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

func (typ *Type) resolveName(block *Block) error {
	var err error
	var d interface{}
	if strings.Contains(typ.Name, ".") == false {
		var loadFromImport bool
		d = block.searchType(typ.Name)
		if d != nil {
			switch d.(type) {
			case *Class:
				if t, ok := d.(*Class); ok &&
					t.IsBuildIn == false {
					_, loadFromImport = shouldAccessFromImports(typ.Name, typ.Pos, t.Pos)
				}
			case *Type:
				if t, ok := d.(*Type); ok &&
					t.IsBuildIn == false {
					_, loadFromImport = shouldAccessFromImports(typ.Name, typ.Pos, t.Pos)
				}
			case *Enum:
				if t, ok := d.(*Enum); ok &&
					t.IsBuildIn == false {
					_, loadFromImport = shouldAccessFromImports(typ.Name, typ.Pos, t.Pos)
				}
			}
		} else {
			loadFromImport = true
		}
		if loadFromImport {
			d, err = typ.getNameFromImport()
			if err != nil {
				return err
			}
		}
	} else { // a.b  in type situation,must be package name
		d, err = typ.getNameFromImport()
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
	return nil
}

func (typ *Type) getNameFromImport() (d interface{}, err error) {
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
	if pp, ok := p.(*Package); ok &&
		pp != nil {
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
		if dd.LoadFromOutSide && dd.IsPublic() == false {
			PackageBeenCompile.Errors = append(PackageBeenCompile.Errors,
				fmt.Errorf("%s class '%s' is not public",
					errMsgPrefix(typ.Pos), dd.Name))
		}
		typ.Type = VariableTypeObject
		typ.Class = dd
		return nil
	case *Type:
		pos := typ.Pos
		alias := typ.Alias
		resolved := typ.Resolved
		*typ = *d.(*Type)
		typ.Pos = pos
		typ.Alias = alias
		typ.Resolved = resolved
		return nil
	case *Enum:
		dd := d.(*Enum)
		if dd.LoadFromOutSide && dd.IsPublic() == false {
			PackageBeenCompile.Errors = append(PackageBeenCompile.Errors,
				fmt.Errorf("%s enum '%s' is not public",
					errMsgPrefix(typ.Pos), dd.Name))
		}
		typ.Type = VariableTypeEnum
		typ.Enum = dd
		return nil
	}
	return fmt.Errorf("%s name '%s' is not a type",
		errMsgPrefix(typ.Pos), typ.Name)
}

func (typ *Type) IsNumber() bool {
	return typ.isInteger() ||
		typ.isFloat()
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

func (typ *Type) isInteger() bool {
	return typ.Type == VariableTypeByte ||
		typ.Type == VariableTypeShort ||
		typ.Type == VariableTypeInt ||
		typ.Type == VariableTypeLong ||
		typ.Type == VariableTypeChar
}

/*
	float or double
*/
func (typ *Type) isFloat() bool {
	return typ.Type == VariableTypeFloat ||
		typ.Type == VariableTypeDouble
}

func (typ *Type) IsPrimitive() bool {
	return typ.IsNumber() ||
		typ.Type == VariableTypeString ||
		typ.Type == VariableTypeBool
}

func (typ *Type) typeString(ret *string) {
	if typ.Alias != nil {
		*ret += typ.Alias.Name
		return
	}
	switch typ.Type {
	case VariableTypeBool:
		*ret += "bool"
	case VariableTypeByte:
		*ret += "byte"
	case VariableTypeChar:
		*ret += "char"
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
		*ret += "object@" + typ.Class.Name
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
		if typ.IsVariableArgs {
			*ret += typ.Array.TypeString() + "..."
		} else {
			*ret += typ.Array.TypeString() + "[]"
		}
	case VariableTypeFunction:
		*ret += "fn " + typ.FunctionType.typeString()
	case VariableTypeEnum:
		*ret += "enum(" + typ.Enum.Name + ")"
	case VariableTypeClass:
		*ret += fmt.Sprintf("class@%s", typ.Class.Name)
	case VariableTypeName:
		*ret += typ.Name // resolve wrong, but typeString is ok to return
	case VariableTypeTemplate:
		*ret += typ.Name
	case VariableTypeDynamicSelector:
		*ret += "dynamicSelector@" + typ.Class.Name
	case VariableTypeVoid:
		*ret += "void"
	case VariableTypePackage:
		*ret += "package@" + typ.Package.Name
	case VariableTypeNull:
		*ret += "null"
	case VariableTypeGlobal:
		*ret += typ.Name
	case VariableTypeMagicFunction:
		*ret = magicIdentifierFunction
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
func (typ *Type) getParameterType(ft *FunctionType) []string {
	if typ.Type == VariableTypeName &&
		ft.haveTemplateName(typ.Name) {
		typ.Type = VariableTypeTemplate // convert to type
	}
	if typ.Type == VariableTypeTemplate {
		return []string{typ.Name}
	}
	if typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeJavaArray {
		return typ.Array.getParameterType(ft)
	}
	if typ.Type == VariableTypeMap {
		ret := []string{}
		if t := typ.Map.K.getParameterType(ft); t != nil {
			ret = append(ret, t...)
		}
		if t := typ.Map.V.getParameterType(ft); t != nil {
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
}

/*
	if there is error,this function will crash
*/
func (typ *Type) bindWithParameterTypes(ft *FunctionType, parameterTypes map[string]*Type) error {
	if typ.Type == VariableTypeTemplate {
		t, ok := parameterTypes[typ.Name]
		if ok == false {
			panic(fmt.Sprintf("typed parameter '%s' not found", typ.Name))
		}
		*typ = *t.Clone() // real bind
		return nil
	}
	if typ.Type == VariableTypeArray || typ.Type == VariableTypeJavaArray {
		return typ.Array.bindWithParameterTypes(ft, parameterTypes)
	}
	if typ.Type == VariableTypeMap {
		if len(typ.Map.K.getParameterType(ft)) > 0 {
			err := typ.Map.K.bindWithParameterTypes(ft, parameterTypes)
			if err != nil {
				return err
			}
		}
		if len(typ.Map.V.getParameterType(ft)) > 0 {
			return typ.Map.V.bindWithParameterTypes(ft, parameterTypes)
		}
	}
	panic("not T")
}

/*

 */
func (typ *Type) canBeBindWithType(ft *FunctionType, mkParameterTypes map[string]*Type, bind *Type) error {
	if err := bind.rightValueValid(); err != nil {
		return err
	}
	if bind.Type == VariableTypeNull {
		return fmt.Errorf("'%s' is un typed", bind.TypeString())
	}
	if typ.Type == VariableTypeTemplate {
		mkParameterTypes[typ.Name] = bind
		return nil
	}
	if typ.Type == VariableTypeArray && bind.Type == VariableTypeArray {
		return typ.Array.canBeBindWithType(ft, mkParameterTypes, bind.Array)
	}
	if typ.Type == VariableTypeJavaArray && bind.Type == VariableTypeJavaArray {
		return typ.Array.canBeBindWithType(ft, mkParameterTypes, bind.Array)
	}
	if typ.Type == VariableTypeMap && bind.Type == VariableTypeMap {
		if len(typ.Map.K.getParameterType(ft)) > 0 {
			err := typ.Map.K.canBeBindWithType(ft, mkParameterTypes, bind.Map.K)
			if err != nil {
				return err
			}
		}
		if len(typ.Map.V.getParameterType(ft)) > 0 {
			return typ.Map.V.canBeBindWithType(ft, mkParameterTypes, bind.Map.V)
		}
	}
	return fmt.Errorf("cannot bind '%s' to '%s'", bind.TypeString(), typ.TypeString())
}

func (leftValue *Type) assignAble(errs *[]error, rightValue *Type) bool {
	if leftValue == rightValue { // equal
		return true
	}
	if leftValue.IsPrimitive() &&
		rightValue.IsPrimitive() {
		return leftValue.Type == rightValue.Type
	}
	if leftValue.IsPointer() && rightValue.Type == VariableTypeNull {
		return true
	}
	if leftValue.Type == VariableTypeObject &&
		leftValue.Class.Name == JavaRootClass &&
		rightValue.IsPointer() {
		return true
	}
	if leftValue.Type == VariableTypeArray &&
		rightValue.Type == VariableTypeArray {
		return leftValue.Array.assignAble(errs, rightValue.Array)
	}
	if leftValue.Type == VariableTypeJavaArray &&
		rightValue.Type == VariableTypeJavaArray {
		if leftValue.IsVariableArgs != rightValue.IsVariableArgs {
			return false
		}
		return leftValue.Array.assignAble(errs, rightValue.Array)
	}

	if leftValue.Type == VariableTypeEnum && rightValue.Type == VariableTypeEnum {
		return leftValue.Enum.Name == rightValue.Enum.Name // same enum
	}
	if leftValue.Type == VariableTypeMap && rightValue.Type == VariableTypeMap {
		return leftValue.Map.K.assignAble(errs, rightValue.Map.K) &&
			leftValue.Map.V.assignAble(errs, rightValue.Map.V)
	}
	if leftValue.Type == VariableTypeFunction &&
		rightValue.Type == VariableTypeFunction {
		return leftValue.FunctionType.equal(rightValue.FunctionType)
	}
	if leftValue.Type == VariableTypeObject && rightValue.Type == VariableTypeObject { // object
		if err := leftValue.Class.loadSelf(leftValue.Pos); err != nil {
			*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(rightValue.Pos), err))
			return false
		}
		if err := rightValue.Class.loadSelf(rightValue.Pos); err != nil {
			*errs = append(*errs, err)
			return false
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
	return false
}

func (typ *Type) Equal(compareTo *Type) bool {
	if typ.Type != compareTo.Type {
		return false //early check
	}
	if typ.IsPrimitive() {
		return true //
	}
	if typ.Type == VariableTypeArray ||
		typ.Type == VariableTypeJavaArray {
		if typ.Type == VariableTypeJavaArray {
			if typ.Type == VariableTypeJavaArray &&
				typ.IsVariableArgs != compareTo.IsVariableArgs {
				return false
			}
		}
		return typ.Array.Equal(compareTo.Array)
	}
	if typ.Type == VariableTypeMap {
		if false == typ.Map.K.Equal(compareTo.Map.K) {
			return false
		}
		return typ.Map.V.Equal(compareTo.Map.V)
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
