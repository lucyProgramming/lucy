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
	LocalVars          []string
}

func (b *Block) isTop() bool {
	return b.Outter == nil
}
func (b *Block) searchByName(name string) (interface{}, error) {
	if t, ok := b.Funcs[name]; ok {
		return t, nil
	}
	if t, ok := b.Classes[name]; ok {
		return t, nil
	}
	if t, ok := b.Vars[name]; ok {
		return t, nil
	}
	if t, ok := b.Consts[name]; ok {
		return t, nil
	}
	if t, ok := b.Enums[name]; ok {
		return t, nil
	}
	if t, ok := b.EnumNames[name]; ok {
		return t, nil
	}
	if b.InheritedAttribute.function.Typ.ClosureVars != nil && b.InheritedAttribute.function.Typ.ClosureVars[name] != nil {
		return b.InheritedAttribute.function.Typ.ClosureVars[name], nil
	}
	if b.Outter != nil {
		t, err := b.Outter.searchByName(name)
		if err == nil && b.isFuntionTopBlock() { //found and in function top block
			if _, ok := t.(*VariableDefinition); ok {
				b.InheritedAttribute.function.Typ.ClosureVars[name] = t.(*VariableDefinition)
			}
		}
	}
	return nil, fmt.Errorf("%s not found", name)
}

func (b *Block) inherite(father *Block) {
	b.InheritedAttribute.p = father.InheritedAttribute.p
	b.InheritedAttribute.istop = father.InheritedAttribute.istop
	b.InheritedAttribute.infor = father.InheritedAttribute.infor
	b.InheritedAttribute.function = father.InheritedAttribute.function
	b.Outter = father
}

//func (b *Block) searchFunction(name string) *Function {
//	bb := b
//	for bb != nil {
//		if i, ok := bb.Funcs[name]; ok {
//			return i[0]
//		}
//		bb = bb.Outter
//	}
//	return nil
//}

type InheritedAttribute struct {
	istop    bool // if it is a top block
	infor    bool // if this statement is in for or not
	function *Function
	returns  ReturnList
	p        *Package
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
	t, err := b.checkExpression(e)
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

func (b *Block) checkExpression_(e *Expression) (t []*VariableType, errs []error) {
	return e.check(b)
}

func (b *Block) checkExpression(e *Expression) (t *VariableType, errs []error) {
	ts, es := b.checkExpression_(e)
	if ts != nil && len(ts) > 1 {
		if es == nil {
			es = make([]error, 0)
		}
		es = append(es, fmt.Errorf("%s multi returns in single value context", errMsgPrefix(e.Pos)))
	}
	if len(ts) > 0 {
		t = ts[0]
	}
	errs = es
	return
}

func (b *Block) checkVar(v *VariableDefinition) []error {
	if v.Expression == nil && v.Typ == nil {
		panic(1)
	}
	var err error
	var expressionVariableType *VariableType
	if v.Expression != nil {
		var es []error
		expressionVariableType, es = b.checkExpression(v.Expression)
		if err != nil {
			return es
		}
	}
	if v.Typ != nil { //means variable typed by assignment
		err = v.Typ.resolve(b)
		if err != nil {
			if err != nil {
				return []error{fmt.Errorf("%s err", errMsgPrefix(v.Pos))}
			}
		}
		if !v.Typ.typeCompatible(expressionVariableType) {
			return []error{fmt.Errorf("%s variable %s defined wrong,cannot assign %s to %s", errMsgPrefix(v.Pos), v.Typ.TypeString(), expressionVariableType.TypeString())}
		}
		return nil
	} else {
		v.Typ = expressionVariableType
	}
	return nil
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

func (b *Block) checkFunctions() []error {
	errs := []error{}
	for _, v := range b.Funcs {
		//function has the sames
		for _, vv := range v {
			errs = append(errs, vv.check(b)...)
		}
		//redeclare errors

	}
	return errs
}

func (b *Block) insert(name string, pos *Pos, d interface{}) error {
	if b.Vars == nil {
		b.Vars = make(map[string]*VariableDefinition)
	}
	if b.Vars[name] == nil {
		return fmt.Errorf("%s name %s already declared as variable", name)
	}
	if b.Classes != nil {
		b.Classes = make(map[string]*Class)
	}
	if b.Classes[name] != nil {
		return fmt.Errorf("%s name %s already declared as class", errMsgPrefix(pos), name)
	}
	if b.Funcs == nil {
		b.Funcs = make(map[string][]*Function)
	}
	if b.Funcs[name] != nil {
		return fmt.Errorf("%s name %s already declared as function", errMsgPrefix(pos), name)
	}
	if b.Consts == nil {
		b.Consts = make(map[string]*Const)
	}
	if b.Consts[name] != nil {
		return fmt.Errorf("%s name %s already declared as const", errMsgPrefix(pos), name)
	}
	if b.Enums == nil {
		b.Enums = make(map[string]*Enum)
	}
	if b.Enums[name] != nil {
		return fmt.Errorf("%s name %s already declared as enum", errMsgPrefix(pos), name)
	}
	if b.EnumNames == nil {
		b.EnumNames = make(map[string]*EnumName)
	}
	if b.EnumNames[name] != nil {
		return fmt.Errorf("%s name %s already declared as enumName", errMsgPrefix(pos), name)
	}
	switch d.(type) {
	case *Class:
		b.Classes[name] = d.(*Class)
	case *Function:
		b.Funcs[name] = append(b.Funcs[name], d.(*Function))
	case *Const:
		b.Consts[name] = d.(*Const)
	case *VariableDefinition:
		b.Vars[name] = d.(*VariableDefinition)
		b.LocalVars = append(b.LocalVars, name)
	case *Enum:
		e := d.(*Enum)
		b.Enums[name] = e
		for _, v := range e.Names {
			err := b.insert(v.Name, v.Pos, v)
			if err != nil {
				return err
			}
		}
	case *EnumName:
		b.EnumNames[name] = d.(*EnumName)
	default:
		panic("????????")
	}
	return nil
}

func (b *Block) loadPackage(name string) (*Package, error) {
	return b.InheritedAttribute.p.loadPackage(name)
}

func (b *Block) isFuntionTopBlock() bool {
	if b.InheritedAttribute.function == nil { // not in a function
		return false
	}
	return b.InheritedAttribute.function != b.Outter.InheritedAttribute.function
}

//func (b *Block) checkExpression(Expression *e) (*VariableType, error) {
//	return nil, nil
//}

//func (b *Block) getPackageName(name )

//func (b *Block) checkVars() []error {
//	errs := make([]error, 0)
//	var es []error
//	for _, v := range b.Vars {
//		es = b.checkVar(v)
//		if errsNotEmpty(es) {
//			errs = append(errs, es...)
//		}
//	}
//	return errs
//}
