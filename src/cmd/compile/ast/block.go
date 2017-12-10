package ast

import (
	"fmt"
)

type Block struct {
	Pos                *Pos
	Vars               map[string]*VariableDefinition
	Consts             map[string]*Const
	Funcs              map[string][]*Function
	Classes            map[string]*Class
	Enums              map[string]*Enum
	EnumNames          map[string]*EnumName
	Outter             *Block //for closure,capture variables
	InheritedAttribute InheritedAttribute
	Statements         []*Statement
	p                  *Package
	LocalVars          []interface{}
}

func (b *Block) isTop() bool {
	return b.Outter == nil
}
func (b *Block) searchByName(name string) interface{} {
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
	//if e.Typ != EXPRESSION_TYPE_IDENTIFIER {
	//	return nil
	//}
	//return b.p.Funcs[e.Data.(string)][0]
	return nil
}

type InheritedAttribute struct {
	istop      bool // if it is a top block
	infor      bool // if this statement is in for or not
	infunction bool // if this in a function situation,return can be availale or not
	returns    ReturnList
	p          *Package
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
	t, err := b.getTypeFromExpression(e)
	if err != nil {
		return false, err
	}
	return t.Typ == VARIABLE_TYPE_BOOL, nil
}

func (b *Block) check(p *Package) []error {
	b.InheritedAttribute.p = p
	errs := []error{}
	errs = append(errs, b.checkConst()...)
	errs = append(errs, b.checkFunctions()...)
	for _, v := range b.Vars {
		errs = append(errs, b.checkVar(v)...)
	}
	errs = append(errs, b.checkClass()...)
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

func (b *Block) checkVar(v *VariableDefinition) []error {
	if v.Expression == nil && v.Typ == nil {
		panic(1)
	}
	var err error
	if v.Expression != nil {
		err = v.Expression.constFold() //fold const error
		if err != nil {
			return []error{fmt.Errorf("%s variable %s defined wrong,err:%v", errMsgPrefix(v.Pos), v.Name, err)}
		}
	}
	if v.Typ != nil { //means variable typed by assignment
		match := v.Typ.matchExpression(b, v.Expression)
		if !match {
			return []error{fmt.Errorf("%s variable %s dose not matched by %s ", errMsgPrefix(v.Pos), v.Name, v.Expression.typeName())}
		}
		return nil
	} else {
		var es []error
		v.Typ, es = b.getTypeFromExpression(v.Expression)
		return es
	}
}

func (p *Block) checkClass() []error {
	errs := []error{}
	for _, v := range p.Classes {
		errs = append(errs, v.check()...)
	}
	return errs
}

func (b *Block) checkConst() []error {
	errs := make([]error, 0)
	for _, v := range b.Consts {
		if v.Expression == nil && v.Typ == nil {
			errs = append(errs, fmt.Errorf("%s const %v has no initiation value", errMsgPrefix(v.Pos), v.Name))
			continue
		}
		is, t, value, err := v.Expression.getConstValue()
		if err != nil {
			errs = append(errs, fmt.Errorf("%s const %v cannot be defined by intiation value", errMsgPrefix(v.Pos), err))
			continue
		}
		if is == false {
			errs = append(errs, fmt.Errorf("%s const %s is not a const value", errMsgPrefix(v.Pos), v.Name))
			continue
		}
		//rewrite
		v.Expression = &Expression{}
		v.Expression.Typ = t
		v.Expression.Data = value
		if v.Typ != nil && v.Expression != nil {
			d, err := v.Typ.constValueValid(v.Expression)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s const %v has worng initiation value", errMsgPrefix(v.Pos), v.Name))
				continue
			}
			v.Data = d
		}
	}
	return errs
}

func (b *Block) checkGlobalVariables() []error {
	errs := make([]error, 0)
	var es []error
	for _, v := range b.Vars {
		es = b.checkVar(v)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
		}
	}
	return errs
}

func (b *Block) checkFunctions() []error {
	errs := []error{}
	for _, v := range b.Funcs {
		//function has the sames
		for _, vv := range v {
			errs = append(errs, vv.check(b)...)
		}
	}
	return errs
}

func (b *Block) insert(name string, pos *Pos, d interface{}) error {
	if b.Classes != nil {
		b.Classes = make(map[string]*Class)
	}
	if b.Funcs == nil {
		b.Funcs = make(map[string][]*Function)
	}
	if b.Consts == nil {
		b.Consts = make(map[string]*Const)
	}
	if b.Enums == nil {
		b.Enums = make(map[string]*Enum)
	}
	if b.EnumNames == nil {
		b.EnumNames = make(map[string]*EnumName)
	}
	return nil
}

//const (
//	ITEM_TYPE_CONST = iota
//	ITEM_TYPE_VAR
//)

//type SymbolicItem struct {
//	Typ   int
//	Var   *VariableDefinition
//	Const *Const
//}
//type SymbolicTable struct {
//	ItemsMap map[string]*SymbolicItem
//}
//
//func (s *SymbolicTable) Insert(name string, pos *Pos, d interface{}) error {
//	if name == "" {
//		panic("name is null string")
//	}
//	if s.ItemsMap == nil {
//		s.ItemsMap = make(map[string]*SymbolicItem)
//	}
//	switch d.(type) {
//	case *VariableDefinition:
//		if _, ok := s.ItemsMap[name]; ok {
//			return fmt.Errorf("%s varaible %s already declared", errMsgPrefix(pos), name)
//		}
//		s.ItemsMap[name] = &SymbolicItem{
//			Typ: ITEM_TYPE_CONST,
//			Var: d.(*VariableDefinition),
//		}
//	case *Const:
//		if _, ok := s.ItemsMap[name]; ok {
//			return fmt.Errorf("%s const %s already declared", errMsgPrefix(pos), name)
//		}
//		s.ItemsMap[name] = &SymbolicItem{
//			Typ:   ITEM_TYPE_CONST,
//			Const: d.(*Const),
//		}
//	default:
//		panic(d.(*VariableDefinition)) // == panic(false) ,runtime panic definitely
//	}
//	return nil
//}
