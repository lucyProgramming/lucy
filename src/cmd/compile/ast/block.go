package ast

import (
	"fmt"
)

type Block struct {
	Pos                *Pos
	InheritedAttribute InheritedAttribute
	Outter             *Block //for closure,capture variables
	SymbolicTable      SymbolicTable
	Statements         []*Statement
	p                  *Package
}

func (b *Block) searchSymbolicItem(name string) *SymbolicItem {
	bb := b
	for bb != nil {
		if x := bb.SymbolicTable.ItemsMap[name]; x != nil {
			return x
		}
		bb = bb.Outter
	}
	return nil
}

func (b *Block) searchSymbolicItemAlsoGlobalVar(name string) *SymbolicItem {
	x := b.searchSymbolicItem(name)
	if x != nil {
		return x
	}
	//try global variable
	t := b.p.Vars[name]
	if t != nil {
		return nil
	}
	return nil
}

func (b *Block) inherite(father *Block) {
	b.InheritedAttribute.p = father.InheritedAttribute.p
	b.InheritedAttribute.istop = father.InheritedAttribute.istop
	b.InheritedAttribute.infor = father.InheritedAttribute.infor
	b.InheritedAttribute.infunction = father.InheritedAttribute.infunction
	b.Outter = father
}

func (b *Block) searchFunction(e *Expression) *Function {
	//bb := b
	//for bb != nil {
	//	if i, ok := bb.SymbolicTable.itemsMap[name]; ok {
	//		if i.Typ.Typ == VARIALBE_TYPE_FUNCTION {
	//			return
	//		}
	//	}
	//	bb = bb.Outter
	//}
	if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
		return nil
	}
	return b.p.Funcs[e.Data.(string)][0]
}

type InheritedAttribute struct {
	istop      bool // if it is a top block
	infor      bool // if this statement is in for or not
	infunction bool // if this in a function situation
	returns    ReturnList
	p          *Package
}

const (
	ITEM_TYPE_CONST = iota
	ITEM_TYPE_VAR
)

type SymbolicItem struct {
	Typ   int
	Var   *VariableDefinition
	Const *Const
}

type SymbolicTable struct {
	ItemsMap map[string]*SymbolicItem
}

func (s *SymbolicTable) Insert(name string, pos *Pos, d interface{}) error {
	if name == "" {
		panic("name is null string")
	}
	if s.ItemsMap == nil {
		s.ItemsMap = make(map[string]*SymbolicItem)
	}
	switch d.(type) {
	case *VariableDefinition:
		if _, ok := s.ItemsMap[name]; ok {
			return fmt.Errorf("%s varaible %s already declared", errMsgPrefix(pos), name)
		}
		s.ItemsMap[name] = &SymbolicItem{
			Typ: ITEM_TYPE_CONST,
			Var: d.(*VariableDefinition),
		}
	case *Const:
		if _, ok := s.ItemsMap[name]; ok {
			return fmt.Errorf("%s const %s already declared", errMsgPrefix(pos), name)
		}
		s.ItemsMap[name] = &SymbolicItem{
			Typ:   ITEM_TYPE_CONST,
			Const: d.(*Const),
		}
	default:
		panic(d.(*VariableDefinition)) // == panic(false) ,runtime panic definitely
	}
	return nil
}

type NameWithType struct {
	Name string
	Typ  *VariableType
}

//check out if expression is bool,must fold const before call this function
func (b *Block) isBoolValue(e *Expression) (bool, []error) {
	if e.Typ == EXPRESSION_TYPE_BOOL { //bool literal
		return true, nil
	}
	t, errs := b.getTypeFromExpression(e)
	if errs != nil && len(errs) > 0 {
		return false, errs
	}
	return t.Typ == VARIABLE_TYPE_BOOL, nil
}

func (b *Block) check(p *Package) []error {
	b.InheritedAttribute.p = p
	errs := []error{}
	for _, s := range b.Statements {
		errs = append(errs, s.check(b)...)
	}
	return errs
}

func (b *Block) getTypeFromExpression(e *Expression) (t *VariableType, errs []error) {
	errs = []error{}
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		t = &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
		}
	case EXPRESSION_TYPE_BYTE:
		t = &VariableType{
			Typ: VARIABLE_TYPE_BYTE,
		}
	case EXPRESSION_TYPE_INT:
		t = &VariableType{
			Typ: VARIABLE_TYPE_INT,
		}
	case EXPRESSION_TYPE_FLOAT:
		t = &VariableType{
			Typ: VARIABLE_TYPE_FLOAT,
		}
	case EXPRESSION_TYPE_STRING:
		t = &VariableType{
			Typ: VARIABLE_TYPE_STRING,
		}
	default:
		panic("unhandled type inference")
	}
	return
}
