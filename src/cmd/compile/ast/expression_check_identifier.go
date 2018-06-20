package ast

import (
	"fmt"
)

func (e *Expression) checkIdentifierExpression(block *Block) (t *Type, err error) {
	identifier := e.Data.(*ExpressionIdentifier)
	if identifier.Name == NO_NAME_IDENTIFIER {
		return nil, fmt.Errorf("%s '%s' is not a valid name",
			errMsgPrefix(e.Pos), NO_NAME_IDENTIFIER)
	}
	fromImport := false
	d, err := block.searchByName(identifier.Name)
	if err != nil {
		return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
	}
	if d == nil {
		i := PackageBeenCompile.getImport(e.Pos.Filename, identifier.Name)
		if i != nil {
			fromImport = true
			d, err = PackageBeenCompile.load(i.ImportName)
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
				p, err := PackageBeenCompile.load(i.ImportName)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				tt.Type = VARIABLE_TYPE_PACKAGE
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Type = VARIABLE_TYPE_OBJECT
				}
				return tt, nil
			}
		}
		f.Used = true
		tt := &Type{}
		tt.Type = VARIABLE_TYPE_FUNCTION
		tt.Pos = e.Pos
		tt.Function = f
		return tt, nil
	case *Variable:
		t := d.(*Variable)
		if fromImport == false && t.IsGlobal { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, t.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.ImportName)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				tt.Type = VARIABLE_TYPE_PACKAGE
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Type = VARIABLE_TYPE_OBJECT
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
				p, err := PackageBeenCompile.load(i.ImportName)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Type = VARIABLE_TYPE_OBJECT
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
				p, err := PackageBeenCompile.load(i.ImportName)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Type = VARIABLE_TYPE_OBJECT
				}
				return tt, nil
			}
		}
		t := &Type{}
		t.Type = VARIABLE_TYPE_CLASS
		t.Pos = e.Pos
		t.Class = c
		return t, nil
	case *EnumName:
		e := d.(*EnumName)
		if fromImport == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, e.Pos, e.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.ImportName)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &Type{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Type = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Type = VARIABLE_TYPE_OBJECT
				}
				return tt, nil
			}
		}
		if e != nil {
			t := &Type{}
			t.Pos = e.Pos
			t.Type = VARIABLE_TYPE_ENUM
			t.EnumName = e
			t.Enum = e.Enum
			identifier.EnumName = e
			return t, nil
		}
	case *Package:
		t := &Type{}
		t.Pos = e.Pos
		t.Type = VARIABLE_TYPE_PACKAGE
		t.Package = d.(*Package)
		return t, nil
	}
	return nil, fmt.Errorf("%s identifier named '%s' is not a expression",
		errMsgPrefix(e.Pos), identifier.Name)
}
