package ast

import (
	"fmt"
)

func (e *Expression) checkIdentiferExpression(block *Block) (t *VariableType, err error) {
	identifer := e.Data.(*ExpressionIdentifer)
	d := block.SearchByName(identifer.Name)
	if d == nil {
		i := PackageBeenCompile.getImport(e.Pos.Filename, identifer.Name)
		if i != nil {
			d, err = PackageBeenCompile.load(i.Resource)
			if err != nil {
				return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
			}
		}
	}
	if d == nil {
		return nil, fmt.Errorf("%s %s not found", errMsgPrefix(e.Pos), identifer.Name)
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		if f.IsGlobal && f.IsBuildin == false {
			i, should := shouldAccessFromImports(identifer.Name, e.Pos, f.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Resource)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &VariableType{}
				tt.Pos = e.Pos
				tt.Typ = VARIABLE_TYPE_PACKAGE
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Typ = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Typ = VARIABLE_TYPE_OBJECT
				}
				return tt, nil
			}
		}
		f.Used = true
		tt := &VariableType{}
		tt.Typ = VARIABLE_TYPE_FUNCTION
		tt.Pos = e.Pos
		tt.Function = f
		return tt, nil
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		if t.IsGlobal {
			i, should := shouldAccessFromImports(identifer.Name, e.Pos, t.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Resource)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &VariableType{}
				tt.Pos = e.Pos
				tt.Typ = VARIABLE_TYPE_PACKAGE
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Typ = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Typ = VARIABLE_TYPE_OBJECT
				}
				return tt, nil
			}
		}
		t.Used = true
		tt := t.Typ.Clone()
		tt.Pos = e.Pos
		identifer.Var = t
		return tt, nil
	case *Const:
		t := d.(*Const)
		if t.IsGlobal {
			i, should := shouldAccessFromImports(identifer.Name, e.Pos, t.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Resource)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &VariableType{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Typ = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Typ = VARIABLE_TYPE_OBJECT
				}
				return tt, nil
			}
		}
		t.Used = true
		e.fromConst(t)
		tt := t.Typ.Clone()
		tt.Pos = e.Pos
		return tt, nil
	case *Class:
		c := d.(*Class)
		if c.IsGlobal {
			i, should := shouldAccessFromImports(identifer.Name, e.Pos, c.Pos)
			if should {
				p, err := PackageBeenCompile.load(i.Resource)
				if err != nil {
					return nil, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err)
				}
				tt := &VariableType{}
				tt.Pos = e.Pos
				if pp, ok := p.(*Package); ok {
					tt.Package = pp
					tt.Typ = VARIABLE_TYPE_PACKAGE
				} else {
					tt.Class = p.(*Class)
					tt.Typ = VARIABLE_TYPE_OBJECT
				}
				return tt, nil
			}
		}
		t := &VariableType{}
		t.Typ = VARIABLE_TYPE_CLASS
		e.Pos = e.Pos
		t.Class = c
		return t, nil
	case (*Package):
		t := &VariableType{}
		t.Pos = e.Pos
		t.Typ = VARIABLE_TYPE_PACKAGE
		t.Package = d.(*Package)
		return t, nil
	default:
		return nil, fmt.Errorf("%s identifier named '%s' is not a expression", errMsgPrefix(e.Pos), identifer.Name)
	}
	return nil, nil
}

func (e *Expression) isThisIdentifierExpression() bool {
	if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
		return false
	}
	return e.Data.(*ExpressionIdentifer).Name == THIS
}
