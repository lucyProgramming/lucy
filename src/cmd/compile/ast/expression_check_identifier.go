package ast

import (
	"fmt"
	"time"
)

func (e *Expression) checkIdentifierExpression(block *Block) (*Type, error) {
	identifier := e.Data.(*ExpressionIdentifier)
	if identifier.Name == NoNameIdentifier {
		return nil, fmt.Errorf("%s '%s' is not a valid name",
			errMsgPrefix(e.Pos), NoNameIdentifier)
	}
	switch identifier.Name {
	case MagicIdentifierFile:
		e.Type = ExpressionTypeString
		e.Data = e.Pos.Filename
		ts, _ := e.check(block)
		return ts[0], nil
	case MagicIdentifierLine:
		e.Type = ExpressionTypeInt
		e.Data = int32(e.Pos.StartLine)
		ts, _ := e.check(block)
		return ts[0], nil
	case MagicIdentifierTime:
		e.Type = ExpressionTypeLong
		e.Data = int64(time.Now().UnixNano())
		ts, _ := e.check(block)
		return ts[0], nil
	case MagicIdentifierClass:
		if block.InheritedAttribute.Class == nil {
			return nil,
				fmt.Errorf("%s '%s' must in class scope", errMsgPrefix(e.Pos), identifier.Name)
		}
		t := &Type{}
		t.Type = VariableTypeClass
		t.Pos = e.Pos
		t.Class = block.InheritedAttribute.Class
		return t, nil
	}
	fromImport := false
	d, err := block.searchIdentifier(identifier.Name)
	if err != nil {
		return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
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
	//fmt.Println(identifier.Name, e.Pos)
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		if fromImport == false && f.IsBuildIn == false && f.IsGlobal && f.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, f.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				result := &Type{}
				result.Pos = e.Pos
				result.Type = VariableTypePackage
				if pp, ok := p.(*Package); ok {
					result.Package = pp
					result.Type = VariableTypePackage
				} else {
					result.Class = p.(*Class)
					result.Type = VariableTypeObject
				}
				return result, nil
			}
		}
		f.Used = true
		result := &Type{}
		result.Type = VariableTypeFunction
		result.FunctionType = &f.Type
		result.Pos = e.Pos
		result.Function = f
		identifier.Function = f
		return result, nil
	case *Variable:
		t := d.(*Variable)
		if fromImport == false && t.IsGlobal && t.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, t.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				result := &Type{}
				result.Pos = e.Pos
				result.Type = VariableTypePackage
				if pp, ok := p.(*Package); ok {
					result.Package = pp
					result.Type = VariableTypePackage
				} else {
					result.Class = p.(*Class)
					result.Type = VariableTypeObject
				}
				return result, nil
			}
		}
		t.Used = true
		result := t.Type.Clone()
		result.Pos = e.Pos
		identifier.Variable = t
		return result, nil
	case *Constant:
		t := d.(*Constant)
		if fromImport == false && t.IsGlobal && t.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, t.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
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
		}
		t.Used = true
		e.fromConst(t)
		result := t.Type.Clone()
		result.Pos = e.Pos
		return result, nil
	case *Class:
		c := d.(*Class)
		if fromImport == false && c.IsGlobal && c.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, c.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
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
		}
		result := &Type{}
		result.Type = VariableTypeClass
		result.Pos = e.Pos
		result.Class = c
		return result, nil
	case *EnumName:
		enumName := d.(*EnumName)
		if fromImport == false && enumName.Enum.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, enumName.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
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
		}
		if enumName != nil {
			result := &Type{}
			result.Pos = enumName.Pos
			result.Type = VariableTypeEnum
			result.EnumName = enumName
			result.Enum = enumName.Enum
			identifier.EnumName = enumName
			return result, nil
		}
	case *Type:
		typ := d.(*Type)
		if fromImport == false && typ.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, typ.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
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
		}
		result := &Type{}
		result.Pos = e.Pos
		result.Type = VariableTypeTypeAlias
		result.AliasType = typ
		return result, nil
	case *Package:
		result := &Type{}
		result.Pos = e.Pos
		result.Type = VariableTypePackage
		result.Package = d.(*Package)
		return result, nil
	}
	return nil, fmt.Errorf("%s identifier named '%s' is not a expression",
		errMsgPrefix(e.Pos), identifier.Name)
}
