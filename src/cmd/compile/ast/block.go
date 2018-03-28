package ast

import (
	"fmt"
)

type Block struct {
	Defers                     []*Defer
	isGlobalVariableDefinition bool
	IsFunctionTopBlock         bool
	IsClassBlock               bool
	Pos                        *Pos
	Consts                     map[string]*Const
	Funcs                      map[string]*Function
	Classes                    map[string]*Class
	Enums                      map[string]*Enum
	EnumNames                  map[string]*EnumName
	Lables                     map[string]*StatementLable
	Outter                     *Block //for closure,capture variables
	InheritedAttribute         InheritedAttribute
	Statements                 []*Statement
	Vars                       map[string]*VariableDefinition
}

type Defer struct {
	StartPc int
	Block   Block
}

func (b *Block) shouldStop(errs []error) bool {
	return (len(PackageBeenCompile.Errors) + len(errs)) >= PackageBeenCompile.NErros2Stop
}

func (b *Block) SearchByName(name string) interface{} {
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
	if b.Lables != nil {
		if l, ok := b.Lables[name]; ok {
			return l
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
	if b.InheritedAttribute.Function != nil {
		v := b.InheritedAttribute.Function.ClosureVars.Search(name)
		if v != nil {
			return v
		}
	}
	if b.Outter == nil {
		return nil
	}
	t := b.Outter.SearchByName(name) // search by outter block
	if t != nil {                    //
		if v, ok := t.(*VariableDefinition); ok && v.IsGlobal == false { // not a global variable
			if b.InheritedAttribute.Function != nil &&
				b.IsFunctionTopBlock &&
				b.InheritedAttribute.Function.IsGlobal == false {
				b.InheritedAttribute.Function.ClosureVars.Insert(b.InheritedAttribute.Function, v)
			}
			//cannot search variable from class body
			if b.InheritedAttribute.class != nil && b.IsClassBlock {
				return nil //
			}
		}
		if l, ok := t.(*StatementLable); ok {
			if b.IsFunctionTopBlock { // search lable from outside out allow
				return nil
			} else {
				return l
			}
		}
		// if it is a function
		if f, ok := t.(*Function); ok && f.IsGlobal == false {
			if b.IsFunctionTopBlock {
				b.InheritedAttribute.Function.ClosureVars.InsertFunction(f)
			}
		}
	}
	return t
}

func (b *Block) inherite(father *Block) {
	if b != father {
		b.InheritedAttribute = father.InheritedAttribute
		b.Outter = father
	}
}

type InheritedAttribute struct {
	StatementFor                 *StatementFor // if this statement is in for or not
	StatementSwitch              *StatementSwitch
	mostCloseForOrSwitchForBreak interface{}
	Function                     *Function
	//OutterFunction               *Function
	class  *Class
	Defer  *Defer
	Defers []*Defer
}

type NameWithType struct {
	Name string
	Typ  *VariableType
}

func (b *Block) check() []error {
	errs := []error{}
	errs = append(errs, b.checkConst()...)
	if b.shouldStop(errs) {
		return errs
	}
	errs = append(errs, b.checkFunctions()...)
	if b.shouldStop(errs) {
		return errs
	}
	errs = append(errs, b.checkClass()...)
	if b.shouldStop(errs) {
		return errs
	}
	for _, v := range b.Vars {
		errs = append(errs, b.checkVar(v)...)
		if b.shouldStop(errs) {
			return errs
		}
	}
	for _, s := range b.Statements {
		errs = append(errs, s.check(b)...)
		if b.shouldStop(errs) {
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
		e.VariableType = t
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
		if expressionVariableType != nil && !v.Typ.TypeCompatible(expressionVariableType) {
			return []error{fmt.Errorf("%s variable %s defined wrong,cannot assign '%s' to '%s'", errMsgPrefix(v.Pos), v.Typ.TypeString(), expressionVariableType.TypeString())}
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
			errs = append(errs, fmt.Errorf("%s const %v has no initiation value", errMsgPrefix(v.Pos), v.Name))
			continue
		}
		is, t, value, err := v.Expression.getConstValue()
		if err != nil {
			errs = append(errs, fmt.Errorf("%s const '%v' defined wrong", errMsgPrefix(v.Pos), v.Name, err))
			continue
		}
		if is == false {
			errs = append(errs, fmt.Errorf("%s const %s is not a const value", errMsgPrefix(v.Pos), v.Name))
			continue
		}
		v.Data = value
		v.Expression.Typ = t
		v.Expression.Data = value
		ts, _ := v.Expression.check(b)
		if v.Typ == nil {
			v.Typ = ts[0]
		} else {
			err = v.Typ.resolve(b)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if !v.Typ.Equal(ts[0]) {
				errs = append(errs, fmt.Errorf("%s cannot assign %s %s", errMsgPrefix(v.Pos), v.Typ.TypeString(), ts[0].TypeString()))
				continue
			}
		}
	}
	return errs
}

func (b *Block) checkFunctions() []error {
	errs := []error{}
	for _, v := range b.Funcs {
		if v.IsBuildin {
			continue
		}
		errs = append(errs, v.check(b)...)
	}
	return errs
}

func (b *Block) Insert(name string, pos *Pos, d interface{}) error {
	return b.insert(name, pos, d)
}
func (b *Block) insert(name string, pos *Pos, d interface{}) error {
	fmt.Println("insert to block:", name, pos, d)
	if v, ok := d.(*VariableDefinition); ok && b.InheritedAttribute.Function.isGlobalVariableDefinition { // global var insert into block
		b := PackageBeenCompile.Block
		if _, ok := b.Vars[name]; ok {
			errmsg := fmt.Sprintf("%s name '%s' already declared as variable,last declared at:\n", errMsgPrefix(pos), name)
			errmsg += fmt.Sprintf("%s", errMsgPrefix(v.Pos))
			return fmt.Errorf(errmsg)
		}
		b.Vars[name] = v
		v.Typ.actionNeedBeenDoneWhenDescribeVariable()
		v.IsGlobal = true // it`s global
		return nil
	}
	if name == "" {
		panic("null name")
	}
	if name == THIS {
		return fmt.Errorf("%s '%s' already been taken", errMsgPrefix(pos), THIS)
	}
	if name == "_" {
		panic("_")
	}
	if b.Vars == nil {
		b.Vars = make(map[string]*VariableDefinition)
	}
	if v, ok := b.Vars[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as variable,last declared at:\n", errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("%s", errMsgPrefix(v.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Classes == nil {
		b.Classes = make(map[string]*Class)
	}
	if c, ok := b.Classes[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as class,last declared at:", errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("%s", errMsgPrefix(c.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Funcs == nil {
		b.Funcs = make(map[string]*Function)
	}
	if f, ok := b.Funcs[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as function,last declared at:", errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("%s", errMsgPrefix(f.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Consts == nil {
		b.Consts = make(map[string]*Const)
	}
	if c, ok := b.Consts[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as const,last declared at:", errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("%s", errMsgPrefix(c.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Enums == nil {
		b.Enums = make(map[string]*Enum)
	}
	if e, ok := b.Enums[name]; ok {
		errmsg := fmt.Sprintf("%s name %s already declared as enum,last declared at:", errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("%s", errMsgPrefix(e.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.EnumNames == nil {
		b.EnumNames = make(map[string]*EnumName)
	}
	if en, ok := b.EnumNames[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as enumName,last declared at:", errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("%s", errMsgPrefix(en.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Lables == nil {
		b.Lables = make(map[string]*StatementLable)
	}
	if l, ok := b.Lables[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as enumName,last declared at:", errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("%s", errMsgPrefix(l.Pos))
		return fmt.Errorf(errmsg)
	}
	switch d.(type) {
	case *Class:
		b.Classes[name] = d.(*Class)
	case *Function:
		t := d.(*Function)
		t.MkVariableType()
		if buildinFunctionsMap[t.Name] != nil {
			return fmt.Errorf("%s function named '%s' is buildin", errMsgPrefix(pos), name)
		}
		if name == MAIN_FUNCTION_NAME && b.Outter != nil {
			return fmt.Errorf("%s '%s' is not available", errMsgPrefix(pos), MAIN_FUNCTION_NAME)
		}
		b.Funcs[name] = t
	case *Const:
		b.Consts[name] = d.(*Const)
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		t.Typ.actionNeedBeenDoneWhenDescribeVariable()
		t.LocalValOffset = b.InheritedAttribute.Function.VarOffset
		b.InheritedAttribute.Function.VarOffset += t.NameWithType.Typ.JvmSlotSize()
		b.InheritedAttribute.Function.OffsetDestinations = append(b.InheritedAttribute.Function.OffsetDestinations, &t.LocalValOffset)
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
	case *StatementLable:
		b.Lables[name] = d.(*StatementLable)
	default:
		panic("????????")
	}
	return nil
}
