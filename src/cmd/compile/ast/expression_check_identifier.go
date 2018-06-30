package ast

import (
	"fmt"
)

func (e *Expression) checkIdentifierExpression(block *Block) (t *Type, err error) {
	identifier := e.Data.(*ExpressionIdentifier)
	if identifier.Name == NoNameIdentifier {
		return nil, fmt.Errorf("%s '%s' is not a valid name",
			errMsgPrefix(e.Pos), NoNameIdentifier)
	}
	fromImport := false
	d, err := block.searchIdentifier(identifier.Name)
	if err != nil {
		return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
	}
	if d == nil {
		i := PackageBeenCompile.getImport(e.Pos.Filename, identifier.Name)
		if i != nil {
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
		if fromImport == false && f.IsGlobal && f.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, f.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				tt.Type = VariableTypePackage
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VariableTypePackage
				} else {
					tt.Class = p.(*Class)
					tt.Type = VariableTypeObject
				}
				return tt, nil
			}
		}
		f.Used = true
		tt := &Type{}
		tt.Type = VariableTypeFunction
		tt.FunctionType = &f.Type
		tt.Pos = e.Pos
		tt.Function = f
		identifier.Function = f
		return tt, nil
	case *Variable:
		t := d.(*Variable)
		if fromImport == false && t.IsGlobal { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, t.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				tt.Type = VariableTypePackage
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VariableTypePackage
				} else {
					tt.Class = p.(*Class)
					tt.Type = VariableTypeObject
				}
				return tt, nil
			}
		}
		t.Used = true
		tt := t.Type.Clone()
		tt.Pos = e.Pos
		identifier.Variable = t
		return tt, nil
	case *Constant:
		t := d.(*Constant)
		if fromImport == false && t.IsGlobal { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, t.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VariableTypePackage
				} else {
					tt.Class = p.(*Class)
					tt.Type = VariableTypeObject
				}
				return tt, nil
			}
		}
		t.Used = true
		e.fromConst(t)
		tt := t.Type.Clone()
		tt.Pos = e.Pos
		return tt, nil
	case *Class:
		c := d.(*Class)
		if fromImport == false && c.IsGlobal { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, c.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VariableTypePackage
				} else {
					tt.Class = p.(*Class)
					tt.Type = VariableTypeObject
				}
				return tt, nil
			}
		}
		t := &Type{}
		t.Type = VariableTypeClass
		t.Pos = e.Pos
		t.Class = c
		return t, nil
	case *EnumName:
		enumName := d.(*EnumName)
		if fromImport == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, enumName.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VariableTypePackage
				} else {
					tt.Class = p.(*Class)
					tt.Type = VariableTypeObject
				}
				return tt, nil
			}
		}
		if enumName != nil {
			t := &Type{}
			t.Pos = enumName.Pos
			t.Type = VariableTypeEnum
			t.EnumName = enumName
			t.Enum = enumName.Enum
			identifier.EnumName = enumName
			return t, nil
		}
	case *Type:
		typ := d.(*Type)
		if fromImport == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, typ.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Import)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VariableTypePackage
				} else {
					tt.Class = p.(*Class)
					tt.Type = VariableTypeObject
				}
				return tt, nil
			}
		}
		t := &Type{}
		t.Pos = e.Pos
		t.Type = VariableTypeAlias
		t.AliasType = typ
		return t, nil
	case *Package:
		t := &Type{}
		t.Pos = e.Pos
		t.Type = VariableTypePackage
		t.Package = d.(*Package)
		return t, nil

	}
	return nil, fmt.Errorf("%s identifier named '%s' is not a expression",
		errMsgPrefix(e.Pos), identifier.Name)
}
