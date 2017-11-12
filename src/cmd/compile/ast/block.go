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
	file               *File
}

func (b *Block) searchSymbolicItem(name string) *SymbolicItem {
	bb := b
	for bb != nil {
		if x := bb.SymbolicTable.itemsMap[name]; x != nil {
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
		return &t.SymbolicItem
	}
	return nil
}

func (b *Block) inherite(father *Block) {
	b.InheritedAttribute.p = father.InheritedAttribute.p
	b.InheritedAttribute.istop = father.InheritedAttribute.istop
	b.InheritedAttribute.infor = father.InheritedAttribute.infor
	b.InheritedAttribute.infunction = father.InheritedAttribute.infunction

	b.file = father.file
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
	return b.p.Funcs[e.Data.(string)]
}

type InheritedAttribute struct {
	istop      bool // if it is a top block
	infor      bool // if this statement is in for or not
	infunction bool // if this in a function situation
	returns    ReturnList
	p          *Package
}

type SymbolicTable struct {
	//items    []*SymbolicItem
	itemsMap map[string]*SymbolicItem // easy to access by name
}

func (s *SymbolicTable) Insert(name string, item *SymbolicItem) error {
	if s.itemsMap == nil {
		s.itemsMap = make(map[string]*SymbolicItem)
	}
	if _, ok := s.itemsMap[name]; ok {
		return fmt.Errorf("symbolic %s already declared", name)
	}
	s.itemsMap[name] = item
	return nil
}

type SymbolicItem struct {
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
func (b *Block) check(p ...*Package) []error {
	if len(p) > 0 {
		b.p = p[0]
		if _, ok := b.p.Files[b.Pos.Filename]; !ok {
			panic("block has no files")
		}
		b.file = b.p.Files[b.Pos.Filename]
	}
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
