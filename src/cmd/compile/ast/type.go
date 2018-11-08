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

func (this *Type) validForTypeAssertOrConversion() bool {
	if this.IsPointer() == false {
		return false
	}
	return true
	//// object or string
	//if this.Type == VariableTypeObject || this.Type == VariableTypeString {
	//	return true
	//}
	//if this.Type == VariableTypeArray && this.Array.IsPrimitive() {
	//	return true
	//}
	//if this.Type == VariableTypeJavaArray {
	//	if this.Array.IsPointer() {
	//		return this.Array.validForTypeAssertOrConversion()
	//	} else {
	//		return true
	//	}
	//}
	//return false
}

func (this *Type) mkDefaultValueExpression() *Expression {
	e := &Expression{}
	e.Op = "defaultValueByCompiler"
	e.IsCompileAuto = true
	e.Pos = this.Pos
	e.Value = this.Clone()
	switch this.Type {
	case VariableTypeBool:
		e.Type = ExpressionTypeBool
		e.Data = false
	case VariableTypeByte:
		e.Type = ExpressionTypeByte
		e.Data = int64(0)
	case VariableTypeShort:
		e.Type = ExpressionTypeShort
		e.Data = int64(0)
	case VariableTypeChar:
		e.Type = ExpressionTypeChar
		e.Data = int64(0)
	case VariableTypeInt:
		e.Type = ExpressionTypeInt
		e.Data = int64(0)
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
		e.Data = this.Enum.DefaultValue
	default:
		e.Type = ExpressionTypeNull
	}
	return e
}

func (this *Type) rightValueValid() error {
	if this.Type == VariableTypeBool ||
		this.IsNumber() ||
		this.Type == VariableTypeString ||
		this.Type == VariableTypeObject ||
		this.Type == VariableTypeArray ||
		this.Type == VariableTypeMap ||
		this.Type == VariableTypeNull ||
		this.Type == VariableTypeJavaArray ||
		this.Type == VariableTypeEnum ||
		this.Type == VariableTypeFunction {
		return nil
	}
	switch this.Type {
	case VariableTypePackage:
		return fmt.Errorf("%s use package '%s' without selector",
			this.Pos.ErrMsgPrefix(), this.Package.Name)
	case VariableTypeClass:
		return fmt.Errorf("%s use class '%s' without selector",
			this.Pos.ErrMsgPrefix(), this.Class.Name)
	case VariableTypeMagicFunction:
		return fmt.Errorf("%s use '%s' without selector",
			this.Pos.ErrMsgPrefix(), magicIdentifierFunction)
	default:
		return fmt.Errorf("%s '%s' is not right value valid",
			this.Pos.ErrMsgPrefix(), this.TypeString())
	}
}

/*
	have type or not
*/
func (this *Type) isTyped() error {
	if err := this.rightValueValid(); err != nil {
		return err
	}
	/*
		null is only untyped right value
	*/
	if this.Type == VariableTypeNull {
		return fmt.Errorf("%s '%s' is not typed",
			this.Pos.ErrMsgPrefix(), this.TypeString())
	}
	return nil
}

/*
	deep clone
*/
func (this *Type) Clone() *Type {
	ret := &Type{}
	*ret = *this
	if ret.Type == VariableTypeArray ||
		ret.Type == VariableTypeJavaArray {
		ret.Array = this.Array.Clone()
	}
	if ret.Type == VariableTypeMap {
		ret.Map = &Map{}
		ret.Map.K = this.Map.K.Clone()
		ret.Map.V = this.Map.V.Clone()
	}
	if this.Type == VariableTypeFunction {
		ret.FunctionType = this.FunctionType.Clone()
	}
	return ret
}

func (this *Type) resolve(block *Block) error {
	if this.Resolved {
		return nil
	}
	this.Resolved = true // single threading
	switch this.Type {
	case VariableTypeTemplate:
		if block.InheritedAttribute.Function.parameterTypes == nil {
			return fmt.Errorf("%s parameter type '%s' not in a template function",
				errMsgPrefix(this.Pos), this.Name)
		}
		if block.InheritedAttribute.Function.parameterTypes[this.Name] == nil {
			return fmt.Errorf("%s parameter type '%s' not found",
				errMsgPrefix(this.Pos), this.Name)
		}
		pos := this.Pos // keep pos
		*this = *block.InheritedAttribute.Function.parameterTypes[this.Name]
		this.Pos = pos // keep pos
		return nil
	case VariableTypeName:
		return this.resolveName(block)
	case VariableTypeGlobal:
		d, exists := PackageBeenCompile.Block.NameExists(this.Name)
		if exists == false {
			return fmt.Errorf("%s '%s' not found",
				errMsgPrefix(this.Pos), this.Name)
		}
		return this.makeTypeFrom(d)
	case VariableTypeArray, VariableTypeJavaArray:
		return this.Array.resolve(block)
	case VariableTypeMap:
		var err error
		if this.Map.K != nil {
			err = this.Map.K.resolve(block)
			if err != nil {
				return err
			}
		}
		if this.Map.V != nil {
			return this.Map.V.resolve(block)
		}
	case VariableTypeFunction:
		for _, v := range this.FunctionType.ParameterList {
			if err := v.Type.resolve(block); err != nil {
				return err
			}
		}
		for _, v := range this.FunctionType.ReturnList {
			if err := v.Type.resolve(block); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Type) resolveName(block *Block) error {
	var err error
	var d interface{}
	if strings.Contains(this.Name, ".") == false {
		var loadFromImport bool
		d = block.searchType(this.Name)
		if d != nil {
			switch d.(type) {
			case *Class:
				if t, ok := d.(*Class); ok &&
					t.IsBuildIn == false {
					_, loadFromImport = shouldAccessFromImports(this.Name, this.Pos, t.Pos)
				}
			case *Type:
				if t, ok := d.(*Type); ok &&
					t.IsBuildIn == false {
					_, loadFromImport = shouldAccessFromImports(this.Name, this.Pos, t.Pos)
				}
			case *Enum:
				if t, ok := d.(*Enum); ok &&
					t.IsBuildIn == false {
					_, loadFromImport = shouldAccessFromImports(this.Name, this.Pos, t.Pos)
				}
			}
		} else {
			loadFromImport = true
		}
		if loadFromImport {
			d, err = this.getNameFromImport()
			if err != nil {
				return err
			}
		}
	} else { // a.b  in type situation,must be package name
		d, err = this.getNameFromImport()
		if err != nil {
			return err
		}
	}
	if d == nil {
		return fmt.Errorf("%s type named '%s' not found", errMsgPrefix(this.Pos), this.Name)
	}
	err = this.makeTypeFrom(d)
	if err != nil {
		return err
	}
	return nil
}

func (this *Type) getNameFromImport() (d interface{}, err error) {
	if strings.Contains(this.Name, ".") == false {
		i := PackageBeenCompile.getImport(this.Pos.Filename, this.Name)
		if i != nil {
			i.Used = true
			return PackageBeenCompile.load(i.Import)
		}
		return nil, fmt.Errorf("%s type named '%s' not found",
			errMsgPrefix(this.Pos), this.Name)
	}
	packageAndName := strings.Split(this.Name, ".")
	i := PackageBeenCompile.getImport(this.Pos.Filename, packageAndName[0])
	if nil == i {
		return nil, fmt.Errorf("%s package '%s' not imported",
			errMsgPrefix(this.Pos), packageAndName[0])
	}
	i.Used = true
	p, err := PackageBeenCompile.load(i.Import)
	if err != nil {
		return nil, fmt.Errorf("%s %v",
			errMsgPrefix(this.Pos), err)
	}
	if pp, ok := p.(*Package); ok &&
		pp != nil {
		var exists bool
		d, exists = pp.Block.NameExists(packageAndName[1])
		if exists == false {
			err = fmt.Errorf("%s '%s' not found",
				errMsgPrefix(this.Pos), packageAndName[1])
		}
		return d, err
	} else {
		return nil, fmt.Errorf("%s '%s' is not a package",
			errMsgPrefix(this.Pos), packageAndName[0])
	}
}

func (this *Type) makeTypeFrom(d interface{}) error {
	switch d.(type) {
	case *Class:
		dd := d.(*Class)
		if dd.LoadFromOutSide && dd.IsPublic() == false {
			PackageBeenCompile.errors = append(PackageBeenCompile.errors,
				fmt.Errorf("%s class '%s' is not public",
					errMsgPrefix(this.Pos), dd.Name))
		}
		this.Type = VariableTypeObject
		this.Class = dd
		return nil
	case *Type:
		pos := this.Pos
		alias := this.Alias
		resolved := this.Resolved
		*this = *d.(*Type)
		this.Pos = pos
		this.Alias = alias
		this.Resolved = resolved
		return nil
	case *Enum:
		dd := d.(*Enum)
		if dd.LoadFromOutSide && dd.isPublic() == false {
			PackageBeenCompile.errors = append(PackageBeenCompile.errors,
				fmt.Errorf("%s enum '%s' is not public",
					errMsgPrefix(this.Pos), dd.Name))
		}
		this.Type = VariableTypeEnum
		this.Enum = dd
		return nil
	}
	return fmt.Errorf("%s name '%s' is not a type",
		errMsgPrefix(this.Pos), this.Name)
}

func (this *Type) IsNumber() bool {
	return this.isInteger() ||
		this.isFloat()
}

func (this *Type) IsPointer() bool {
	return this.Type == VariableTypeObject ||
		this.Type == VariableTypeArray ||
		this.Type == VariableTypeJavaArray ||
		this.Type == VariableTypeMap ||
		this.Type == VariableTypeString ||
		this.Type == VariableTypeNull ||
		this.Type == VariableTypeFunction
}

func (this *Type) isInteger() bool {
	return this.Type == VariableTypeByte ||
		this.Type == VariableTypeShort ||
		this.Type == VariableTypeInt ||
		this.Type == VariableTypeLong ||
		this.Type == VariableTypeChar
}

/*
	float or double
*/
func (this *Type) isFloat() bool {
	return this.Type == VariableTypeFloat ||
		this.Type == VariableTypeDouble
}

func (this *Type) IsPrimitive() bool {
	return this.IsNumber() ||
		this.Type == VariableTypeString ||
		this.Type == VariableTypeBool
}

func (this *Type) typeString(ret *string) {
	if this.Alias != nil {
		*ret += this.Alias.Name
		return
	}
	switch this.Type {
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
		*ret += "object@" + this.Class.Name
	case VariableTypeMap:
		*ret += "map{"
		*ret += this.Map.K.TypeString()
		*ret += " -> "
		*ret += this.Map.V.TypeString()
		*ret += "}"
	case VariableTypeArray:
		*ret += "[]"
		this.Array.typeString(ret)
	case VariableTypeJavaArray:
		if this.IsVariableArgs {
			*ret += this.Array.TypeString() + "..."
		} else {
			*ret += this.Array.TypeString() + "[]"
		}
	case VariableTypeFunction:
		*ret += "fn " + this.FunctionType.TypeString()
	case VariableTypeEnum:
		*ret += "enum(" + this.Enum.Name + ")"
	case VariableTypeClass:
		*ret += fmt.Sprintf("class@%s", this.Class.Name)
	case VariableTypeName:
		*ret += this.Name // resolve wrong, but TypeString is ok to return
	case VariableTypeTemplate:
		*ret += this.Name
	case VariableTypeDynamicSelector:
		*ret += "dynamicSelector@" + this.Class.Name
	case VariableTypeVoid:
		*ret += "void"
	case VariableTypePackage:
		*ret += "package@" + this.Package.Name
	case VariableTypeNull:
		*ret += "null"
	case VariableTypeGlobal:
		*ret += this.Name
	case VariableTypeMagicFunction:
		*ret = magicIdentifierFunction
	default:
		panic(this.Type)
	}
}

//可读的类型信息
func (this *Type) TypeString() string {
	t := ""
	this.typeString(&t)
	return t
}
func (this *Type) getParameterType(ft *FunctionType) []string {
	if this.Type == VariableTypeName &&
		ft.haveTemplateName(this.Name) {
		this.Type = VariableTypeTemplate // convert to type
	}
	if this.Type == VariableTypeTemplate {
		return []string{this.Name}
	}
	if this.Type == VariableTypeArray ||
		this.Type == VariableTypeJavaArray {
		return this.Array.getParameterType(ft)
	}
	if this.Type == VariableTypeMap {
		ret := []string{}
		ret = append(ret, this.Map.K.getParameterType(ft)...)
		ret = append(ret, this.Map.V.getParameterType(ft)...)
		return ret
	}
	return nil
}

func (this *Type) canBeBindWithParameterTypes(parameterTypes map[string]*Type) error {
	if this.Type == VariableTypeTemplate {
		_, ok := parameterTypes[this.Name]
		if ok == false {
			return fmt.Errorf("typed parameter '%s' not found", this.Name)
		}
		return nil
	}
	if this.Type == VariableTypeArray ||
		this.Type == VariableTypeJavaArray {
		return this.Array.canBeBindWithParameterTypes(parameterTypes)
	}
	if this.Type == VariableTypeMap {
		err := this.Map.K.canBeBindWithParameterTypes(parameterTypes)
		if err != nil {
			return err
		}
		return this.Map.V.canBeBindWithParameterTypes(parameterTypes)
	}
	return nil
}

/*
	if there is error,this function will crash
*/
func (this *Type) bindWithParameterTypes(ft *FunctionType, parameterTypes map[string]*Type) error {
	if this.Type == VariableTypeTemplate {
		t, ok := parameterTypes[this.Name]
		if ok == false {
			panic(fmt.Sprintf("typed parameter '%s' not found", this.Name))
		}
		*this = *t.Clone() // real bind
		return nil
	}
	if this.Type == VariableTypeArray || this.Type == VariableTypeJavaArray {
		return this.Array.bindWithParameterTypes(ft, parameterTypes)
	}
	if this.Type == VariableTypeMap {
		if len(this.Map.K.getParameterType(ft)) > 0 {
			err := this.Map.K.bindWithParameterTypes(ft, parameterTypes)
			if err != nil {
				return err
			}
		}
		if len(this.Map.V.getParameterType(ft)) > 0 {
			return this.Map.V.bindWithParameterTypes(ft, parameterTypes)
		}
	}
	panic("not T")
}

/*

 */
func (this *Type) canBeBindWithType(ft *FunctionType, mkParameterTypes map[string]*Type, bind *Type) error {
	if err := bind.rightValueValid(); err != nil {
		return err
	}
	if bind.Type == VariableTypeNull {
		return fmt.Errorf("'%s' is un typed", bind.TypeString())
	}
	if this.Type == VariableTypeTemplate {
		mkParameterTypes[this.Name] = bind
		return nil
	}
	if this.Type == VariableTypeArray && bind.Type == VariableTypeArray {
		return this.Array.canBeBindWithType(ft, mkParameterTypes, bind.Array)
	}
	if this.Type == VariableTypeJavaArray && bind.Type == VariableTypeJavaArray {
		return this.Array.canBeBindWithType(ft, mkParameterTypes, bind.Array)
	}
	if this.Type == VariableTypeMap && bind.Type == VariableTypeMap {
		if len(this.Map.K.getParameterType(ft)) > 0 {
			err := this.Map.K.canBeBindWithType(ft, mkParameterTypes, bind.Map.K)
			if err != nil {
				return err
			}
		}
		if len(this.Map.V.getParameterType(ft)) > 0 {
			return this.Map.V.canBeBindWithType(ft, mkParameterTypes, bind.Map.V)
		}
	}
	return fmt.Errorf("cannot bind '%s' to '%s'", bind.TypeString(), this.TypeString())
}

func (this *Type) assignAble(errs *[]error, rightValue *Type) bool {
	leftValue := this
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

func (this *Type) Equal(compareTo *Type) bool {
	leftValue := this
	if leftValue.Type != compareTo.Type {
		return false //early check
	}
	if leftValue.IsPrimitive() {
		return true //
	}
	if leftValue.Type == VariableTypeVoid {
		return true
	}
	if leftValue.Type == VariableTypeArray ||
		leftValue.Type == VariableTypeJavaArray {
		if leftValue.Type == VariableTypeJavaArray &&
			leftValue.IsVariableArgs != compareTo.IsVariableArgs {
			return false
		}
		return leftValue.Array.Equal(compareTo.Array)
	}
	if leftValue.Type == VariableTypeMap {
		return leftValue.Map.K.Equal(compareTo.Map.K) &&
			leftValue.Map.V.Equal(compareTo.Map.V)
	}
	if leftValue.Type == VariableTypeEnum {
		return leftValue.Enum.Name == compareTo.Enum.Name
	}
	if leftValue.Type == VariableTypeObject {
		return leftValue.Class.Name == compareTo.Class.Name
	}
	if leftValue.Type == VariableTypeFunction {
		return leftValue.FunctionType.equal(compareTo.FunctionType)
	}
	return false
}
