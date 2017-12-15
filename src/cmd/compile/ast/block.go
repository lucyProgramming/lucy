package ast

import (
	"fmt"
)

type Block struct {
	Pos                *Pos
	Vars               map[string]*VariableDefinition
	Consts             map[string]*Const
	Funcs              map[string]*Function
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
	fmt.Println("!!!!!!!!!!!!!!", name, b.Funcs, b.Outter)
	if b.Funcs != nil {
		if t, ok := b.Funcs[name]; ok {
			return t, nil
		}
	}
	if b.Classes != nil {
		if t, ok := b.Classes[name]; ok {
			return t, nil
		}
	}
	if b.Vars != nil {
		if t, ok := b.Vars[name]; ok {
			return t, nil
		}
	}
	if b.Consts != nil {
		if t, ok := b.Consts[name]; ok {
			return t, nil
		}
	}
	if b.Enums != nil {
		if t, ok := b.Enums[name]; ok {
			return t, nil
		}
	}
	if b.EnumNames != nil {
		if t, ok := b.EnumNames[name]; ok {
			return t, nil
		}
	}
	if b.InheritedAttribute.function != nil &&
		b.InheritedAttribute.function.Typ.ClosureVars != nil &&
		b.InheritedAttribute.function.Typ.ClosureVars[name] != nil {
		return b.InheritedAttribute.function.Typ.ClosureVars[name], nil
	}
	if b.Outter == nil {
		return nil, fmt.Errorf("%s not found", name)
	}
	t, err := b.Outter.searchByName(name)
	if err == nil && b.isFuntionTopBlock() && b.Outter.Outter != nil { //found and in function top block and b.Outter is not top block
		if _, ok := t.(*VariableDefinition); ok {
			b.InheritedAttribute.function.Typ.ClosureVars[name] = t.(*VariableDefinition)
		}
	}
	return t, err
}

func (b *Block) inherite(father *Block) {
	b.InheritedAttribute.p = father.InheritedAttribute.p
	b.InheritedAttribute.istop = father.InheritedAttribute.istop
	b.InheritedAttribute.infor = father.InheritedAttribute.infor
	b.InheritedAttribute.function = father.InheritedAttribute.function
	b.Outter = father
}

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

func (b *Block) check(father *Block) []error {
	if father != nil {
		b.inherite(father)
	}
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
	errs = []error{}
	ts, es := b.checkExpression_(e)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if ts != nil && len(ts) > 1 {
		errs = append(errs, fmt.Errorf("%s multi values in single value context", errMsgPrefix(e.Pos)))
	}
	if len(ts) > 0 {
		t = ts[0]
	}
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
		if expressionVariableType != nil && !v.Typ.typeCompatible(expressionVariableType) {
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
			v.Value = d
		}
	}
	return errs
}

func (b *Block) checkFunctions() []error {
	errs := []error{}
	for _, v := range b.Funcs {
		errs = append(errs, v.check(b)...)
	}
	return errs
}

func (b *Block) insert(name string, pos *Pos, d interface{}) error {
	if name == "__main__" { // special name
		return fmt.Errorf("%s '__main__' already been token", errMsgPrefix(pos))
	}
	fmt.Println("***************", name, d)
	if b.Vars == nil {
		b.Vars = make(map[string]*VariableDefinition)
	}
	if b.Vars[name] != nil {
		return fmt.Errorf("%s name '%s' already declared as variable", errMsgPrefix(pos), name)
	}
	if b.Classes == nil {
		b.Classes = make(map[string]*Class)
	}
	if b.Classes[name] != nil {
		return fmt.Errorf("%s name '%s' already declared as class", errMsgPrefix(pos), name)
	}
	if b.Funcs == nil {
		b.Funcs = make(map[string]*Function)
	}
	if b.Funcs[name] != nil {
		return fmt.Errorf("%s name '%s' already declared as function", errMsgPrefix(pos), name)
	}
	if b.Consts == nil {
		b.Consts = make(map[string]*Const)
	}
	if b.Consts[name] != nil {
		return fmt.Errorf("%s name '%s' already declared as const", errMsgPrefix(pos), name)
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
		return fmt.Errorf("%s name '%s' already declared as enumName", errMsgPrefix(pos), name)
	}

	switch d.(type) {
	case *Class:
		b.Classes[name] = d.(*Class)
	case *Function:
		b.Funcs[name] = d.(*Function)
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
