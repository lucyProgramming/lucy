package ast

import (
	"fmt"
	"strings"
)

type Block struct {
	IsFunctionTopBlock bool
	IsClassBlock       bool
	Pos                *Pos
	p                  *Package
	Consts             map[string]*Const
	Funcs              map[string]*Function
	Classes            map[string]*Class
	Enums              map[string]*Enum
	EnumNames          map[string]*EnumName
	Outter             *Block //for closure,capture variables
	InheritedAttribute InheritedAttribute
	Statements         []*Statement
	Vars               map[string]*VariableDefinition
}

func (b *Block) shouldStop(errs []error) bool {
	return (len(b.InheritedAttribute.p.Errors) + len(errs)) >= b.InheritedAttribute.p.NErros2Stop
}

func (b *Block) searchByName(name string) interface{} {
	if b.Funcs != nil {
		if t, ok := b.Funcs[name]; ok {
			return t
		}
	}
	if b.Classes != nil {
		if t, ok := b.Classes[name]; ok {
			return t
		}
	}
	if b.Vars != nil {
		if t, ok := b.Vars[name]; ok {
			return t
		}
	}
	if b.Consts != nil {
		if t, ok := b.Consts[name]; ok {
			return t
		}
	}
	if b.Enums != nil {
		if t, ok := b.Enums[name]; ok {
			return t
		}
	}
	if b.EnumNames != nil {
		if t, ok := b.EnumNames[name]; ok {
			return t
		}
	}
	if b.InheritedAttribute.function != nil {
		v := b.InheritedAttribute.function.ClosureVars.Search(name)
		if v != nil {
			return v
		}
	}
	if b.Outter == nil {
		return nil
	}
	t := b.Outter.searchByName(name)
	if t != nil { //
		if v, ok := t.(*VariableDefinition); ok && v.IsGlobal == false {
			if b.InheritedAttribute.function != nil && b.IsFunctionTopBlock &&
				b.InheritedAttribute.function.IsGlobal == false {
				b.InheritedAttribute.function.ClosureVars.Insert(v)
			}
			//cannot search variable from class body
			if b.InheritedAttribute.class != nil && b.IsClassBlock {
				return nil //
			}
		}
	}
	return t
}

func (b *Block) inherite(father *Block) {
	b.InheritedAttribute = father.InheritedAttribute
	b.Outter = father
}

type InheritedAttribute struct {
	StatementFor                 *StatementFor // if this statement is in for or not
	StatementSwitch              *StatementSwitch
	mostCloseForOrSwitchForBreak interface{}
	function                     *Function
	class                        *Class
	p                            *Package
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
	if b.shouldStop(errs) {
		return errs
	}
	errs = append(errs, b.checkFunctions()...)
	if father.shouldStop(errs) {
		return errs
	}
	errs = append(errs, b.checkClass()...)
	if father.shouldStop(errs) {
		return errs
	}
	for _, v := range b.Vars {
		errs = append(errs, b.checkVar(v)...)
		if father.shouldStop(errs) {
			return errs
		}
	}
	for _, s := range b.Statements {
		errs = append(errs, s.check(b)...)
		if father.shouldStop(errs) {
			return errs
		}
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

func (b *Block) checkClass() []error {
	errs := []error{}
	for _, v := range b.Classes {
		errs = append(errs, v.check(b)...)
	}
	return errs
}

func (b *Block) checkConst() []error {
	errs := make([]error, 0)
	for _, v := range b.Consts {
		if v.Expression == nil {
			panic("should not happen")
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
		v.Expression.Typ = t
		v.Expression.Data = value
		ts, _ := v.Expression.check(b)
		v.Value = value
		if v.Typ != nil {
			err = v.Typ.resolve(b)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if !v.Typ.Equal(ts[0]) {
				errs = append(errs, fmt.Errorf("%s cannot assign %s %s", v.Typ.TypeString(), ts[0].TypeString()))
				continue
			}
		}
	}
	return errs
}

func (b *Block) checkFunctions() []error {
	errs := []error{}
	for _, v := range b.Funcs {
		if v.Isbuildin {
			continue
		}
		errs = append(errs, v.check(b)...)

	}
	return errs
}

func (b *Block) insert(name string, pos *Pos, d interface{}) error {
	fmt.Println("EEEEEEEEEEEEEEEE", name, d)
	if name == "" {
		panic("null name")
	}
	if name == "__main__" { // special name
		return fmt.Errorf("%s '__main__' already been token", errMsgPrefix(pos))
	}
	if name == THIS {
		return fmt.Errorf("%s 'this' already been token", errMsgPrefix(pos))
	}
	if name == "_" {
		panic("_")
	}
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
	if _, ok := buildinFunctionsMap[name]; ok {
		return fmt.Errorf("%s function named '%s' is buildin", errMsgPrefix(pos), name)
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
		t := d.(*Function)
		if buildinFunctionsMap[t.Name] != nil {
			return fmt.Errorf("%s function named '%s' is buildin", errMsgPrefix(pos), name)
		}
		b.Funcs[name] = t
	case *Const:
		b.Consts[name] = d.(*Const)
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		t.LocalValOffset = b.InheritedAttribute.function.Varoffset
		b.InheritedAttribute.function.Varoffset += t.NameWithType.Typ.JvmSlotSize()
		t.mkTypRight()
		b.Vars[name] = t
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
func (b *Block) loadClass(name string) (*Class, error) {
	pname := name[0:strings.LastIndex(name, "/")]
	cname := name[strings.LastIndex(name, "/")+1:]
	p, err := b.loadPackage(pname)
	if err != nil {
		return nil, err
	}
	c := p.Block.Classes[cname]
	if c == nil {
		err = fmt.Errorf("class %s not found", cname)
	}
	return c, err
}
