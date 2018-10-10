package ast

import (
	"fmt"
	"time"
)

func (e *Expression) checkIdentifierExpression(block *Block) (*Type, error) {
	identifier := e.Data.(*ExpressionIdentifier)
	if identifier.Name == NoNameIdentifier {
		//_ is not valid
		return nil, fmt.Errorf("%s '%s' is not a valid name",
			errMsgPrefix(e.Pos), NoNameIdentifier)
	}
	//handle magic identifier
	switch identifier.Name {
	case MagicIdentifierFile:
		e.Type = ExpressionTypeString
		e.Data = e.Pos.Filename
		result, _ := e.checkSingleValueContextExpression(block)
		return result, nil
	case MagicIdentifierLine:
		e.Type = ExpressionTypeInt
		e.Data = int32(e.Pos.Line)
		result, _ := e.checkSingleValueContextExpression(block)
		return result, nil
	case MagicIdentifierTime:
		e.Type = ExpressionTypeLong
		e.Data = int64(time.Now().UnixNano())
		result, _ := e.checkSingleValueContextExpression(block)
		return result, nil
	case MagicIdentifierClass:
		if block.InheritedAttribute.Class == nil {
			return nil,
				fmt.Errorf("%s '%s' must in class scope", errMsgPrefix(e.Pos), identifier.Name)
		}
		result := &Type{}
		result.Type = VariableTypeClass
		result.Pos = e.Pos
		result.Class = block.InheritedAttribute.Class
		return result, nil
	case MagicIdentifierFunction:
		if block.InheritedAttribute.Function.isGlobalVariableDefinition ||
			block.InheritedAttribute.Function.isPackageInitBlockFunction {
			return nil,
				fmt.Errorf("%s '%s' must in function scope", errMsgPrefix(e.Pos), identifier.Name)
		}
		result := &Type{}
		result.Type = VariableTypeMagicFunction
		result.Pos = e.Pos
		result.Function = block.InheritedAttribute.Function
		return result, nil
	}
	fromImport := false
	d, err := block.searchIdentifier(e.Pos, identifier.Name)
	if err != nil {
		return nil, err
	}
	if d == nil {
		i := PackageBeenCompile.getImport(e.Pos.Filename, identifier.Name)
		if i != nil {
			i.Used = true
			fromImport = true
			d, err = PackageBeenCompile.load(i.Import)
			if err != nil {
				return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
			}
		}
	}
	if d == nil {
		return nil, fmt.Errorf("%s '%s' not found", errMsgPrefix(e.Pos), identifier.Name)
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		if fromImport == false &&
			f.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, f.Pos)
			if should {
				return e.checkIdentifierThroughImports(i)
			}
		}
		if f.IsGlobalMain() {
			return nil, fmt.Errorf("%s fucntion is global main", errMsgPrefix(e.Pos))
		}
		if f.IsBuildIn {
			return nil, fmt.Errorf("%s fucntion '%s' is buildin",
				errMsgPrefix(e.Pos), f.Name)
		}
		if f.TemplateFunction != nil {
			return nil, fmt.Errorf("%s fucntion '%s' a template function",
				errMsgPrefix(e.Pos), f.Name)
		}
		f.Used = true
		result := &Type{}
		result.Type = VariableTypeFunction
		result.FunctionType = &f.Type
		result.Pos = e.Pos
		identifier.Function = f
		return result, nil
	case *Variable:
		t := d.(*Variable)
		if fromImport == false &&
			t.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, t.Pos)
			if should {
				return e.checkIdentifierThroughImports(i)
			}
		}
		t.Used = true
		result := t.Type.Clone()
		result.Pos = e.Pos
		identifier.Variable = t
		return result, nil
	case *Constant:
		t := d.(*Constant)
		if fromImport == false && t.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, t.Pos)
			if should {
				return e.checkIdentifierThroughImports(i)
			}
		}
		t.Used = true
		e.fromConst(t)
		result := t.Type.Clone()
		result.Pos = e.Pos
		return result, nil
	case *Class:
		c := d.(*Class)
		if fromImport == false && c.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, c.Pos)
			if should {
				return e.checkIdentifierThroughImports(i)
			}
		}
		result := &Type{}
		result.Type = VariableTypeClass
		result.Pos = e.Pos
		result.Class = c
		return result, nil
	case *EnumName:
		enumName := d.(*EnumName)
		if fromImport == false &&
			enumName.Enum.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, enumName.Pos)
			if should {
				return e.checkIdentifierThroughImports(i)
			}
		}
		result := &Type{}
		result.Pos = enumName.Pos
		result.Type = VariableTypeEnum
		result.EnumName = enumName
		result.Enum = enumName.Enum
		identifier.EnumName = enumName
		return result, nil
	case *Package:
		// must load from import
		result := &Type{}
		result.Pos = e.Pos
		result.Type = VariableTypePackage
		result.Package = d.(*Package)
		return result, nil
	}
	return nil, fmt.Errorf("%s identifier '%s' is not a expression , but '%s'",
		errMsgPrefix(e.Pos), identifier.Name, block.identifierIsWhat(d))
}

func (e *Expression) checkIdentifierThroughImports(it *Import) (*Type, error) {
	p, err := PackageBeenCompile.load(it.Import)
	if err != nil {
		return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
	}
	result := &Type{}
	result.Pos = e.Pos
	if pp, ok := p.(*Package); ok {
		result.Package = pp
		result.Type = VariableTypePackage
	} else {
		result.Class = p.(*Class)
		result.Type = VariableTypeObject
	}
	return result, nil
}
