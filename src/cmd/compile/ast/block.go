package ast

type Block struct {
	Pos                Pos
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
	b.file = father.file
	b.Outter = father
}
func (b *Block) searchFunction(name string) *Function {
	//bb := b
	//for bb != nil {
	//	if i, ok := bb.SymbolicTable.itemsMap[name]; ok {
	//		if i.Typ.Typ == VARIALBE_TYPE_FUNCTION {
	//			return
	//		}
	//	}
	//	bb = bb.Outter
	//}
	return b.p.Funcs[name]
}

type InheritedAttribute struct {
	istop bool // if it is a top block
	infor bool // if this statement is in for or not
	p     *Package
}

type SymbolicTable struct {
	items    []*SymbolicItem
	itemsMap map[string]*SymbolicItem // easy to access by name
}

type SymbolicItem struct {
	Name string
	Typ  *VariableType
}

func (b *Block) check(p *Package) []error {
	b.p = p
	if _, ok := b.p.Files[b.Pos.Filename]; !ok {
		panic("block has no files")
	}
	b.file = b.p.Files[b.Pos.Filename]
	errs := []error{}
	for _, s := range b.Statements {
		errs = append(errs, s.check(b)...)
	}
	return errs
}
