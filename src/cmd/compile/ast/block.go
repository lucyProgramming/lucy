package ast

import (
	"fmt"
)

type Block struct {
	DeadEnding                 bool
	Defers                     []*Defer
	isGlobalVariableDefinition bool
	IsFunctionTopBlock         bool
	IsClassBlock               bool
	Pos                        *Pos
	EndPos                     *Pos
	Consts                     map[string]*Const
	Funcs                      map[string]*Function
	Classes                    map[string]*Class
	Enums                      map[string]*Enum
	EnumNames                  map[string]*EnumName
	Lables                     map[string]*StatementLable
	Types                      map[string]*VariableType
	Outter                     *Block //for closure,capture variables
	InheritedAttribute         InheritedAttribute
	Statements                 []*Statement
	Vars                       map[string]*VariableDefinition
	ClosureFuncs               map[string]*Function //in "Funcs" too
}

func (b *Block) HaveVariableDefinition() bool {
	if b.ClosureFuncs == nil && b.Vars == nil {
		return false
	}
	return len(b.ClosureFuncs) > 0 || len(b.Vars) > 0
}

func (b *Block) NameExists(name string) bool {
	if b.Funcs != nil {
		if _, ok := b.Funcs[name]; ok {
			return true
		}
	}
	if b.Classes != nil {
		if _, ok := b.Classes[name]; ok {
			return true
		}
	}
	if b.Vars != nil {
		if _, ok := b.Vars[name]; ok {
			return true
		}
	}
	if b.Lables != nil {
		if _, ok := b.Lables[name]; ok {
			return true
		}
	}
	if b.Consts != nil {
		if _, ok := b.Consts[name]; ok {
			return true

		}
	}
	if b.Enums != nil {
		if _, ok := b.Enums[name]; ok {
			return true

		}
	}
	if b.EnumNames != nil {
		if _, ok := b.EnumNames[name]; ok {
			return true

		}
	}
	if b.Types != nil {
		if _, ok := b.Types[name]; ok {
			return true
		}
	}
	return false
}

func (b *Block) searchLable(name string) *StatementLable {
	for b != nil {
		if b.Lables != nil {
			if l, ok := b.Lables[name]; ok {
				return l
			}
		}
		b = b.Outter
	}
	return nil
}

/*
	search anything
*/
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
	if b.Types != nil {
		if t, ok := b.Types[name]; ok {
			return t
		}
	}
	// search closure
	if b.InheritedAttribute.Function != nil {
		v := b.InheritedAttribute.Function.ClosureVars.Search(name)
		if v != nil {
			return v
		}
	}
	if b.Outter == nil {
		return searchBuildIns(name)
	}
	t := b.Outter.SearchByName(name) // search by outter block
	if t != nil {                    //
		if v, ok := t.(*VariableDefinition); ok && v.IsGlobal == false { // not a global variable
			if b.IsFunctionTopBlock &&
				b.InheritedAttribute.Function.IsGlobal == false {

				if v.Name == THIS {
					return nil // capture this not allow
				}
				b.InheritedAttribute.Function.ClosureVars.InsertVar(v)
			}
			//cannot search variable from class body
			if b.InheritedAttribute.class != nil && b.IsClassBlock {
				return nil //
			}
		}
		if l, ok := t.(*StatementLable); ok {
			if b.IsFunctionTopBlock { // search lable from outside out not allow
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

func (b *Block) checkStatements() []error {
	errs := []error{}
	for k, s := range b.Statements {
		b.InheritedAttribute.StatementOffset = k
		errs = append(errs, s.check(b)...)
		if PackageBeenCompile.shouldStop(errs) {
			return errs
		}
	}
	return errs
}

func (b *Block) checkExpression(e *Expression, singleValueContext bool) (t *VariableType, errs []error) {
	errs = []error{}
	ts, es := e.check(b)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}

	if ts != nil && len(ts) > 1 && singleValueContext {
		errs = append(errs, fmt.Errorf("%s multi values in single value context",
			errMsgPrefix(e.Pos)))
	}
	if len(ts) > 0 {
		t = ts[0]
	}
	return
}

func (b *Block) checkConst() []error {
	errs := make([]error, 0)
	for _, c := range b.Consts {
		if c.Name == NO_NAME_IDENTIFIER {
			err := fmt.Errorf("%s '%s' is not a valid name",
				errMsgPrefix(c.Pos), c.Name)
			errs = append(errs, err)
			delete(b.Consts, c.Name)
			continue
		}
		err := checkConst(b, c)
		if err != nil && c.Typ == nil {
			errs = append(errs, err)
			delete(b.Consts, c.Name)
		}
	}
	return errs
}

func (b *Block) Insert(name string, pos *Pos, d interface{}) error {
	return b.insert(name, pos, d)
}
func (b *Block) insert(name string, pos *Pos, d interface{}) error {
	if v, ok := d.(*VariableDefinition); ok && b.InheritedAttribute.Function.isGlobalVariableDefinition { // global var insert into block
		b := PackageBeenCompile.Block
		if vv, ok := b.Vars[name]; ok {
			errmsg := fmt.Sprintf("%s name '%s' already declared as variable,first declared at:\n",
				errMsgPrefix(pos), name)
			errmsg += fmt.Sprintf("\t%s", errMsgPrefix(vv.Pos))
			return fmt.Errorf(errmsg)
		}
		b.Vars[name] = v
		v.IsGlobal = true // it`s global
		return nil
	}
	if name == "" {
		return fmt.Errorf("%s name is null string", errMsgPrefix(pos))
	}
	if name == THIS {
		return fmt.Errorf("%s '%s' already been taken", errMsgPrefix(pos), THIS)
	}
	if name == "_" {
		return fmt.Errorf("%s '%s' is not a valid name", errMsgPrefix(pos), name)
	}

	if b.Vars == nil {
		b.Vars = make(map[string]*VariableDefinition)
	}
	if v, ok := b.Vars[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as variable,first declared at:\n",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(v.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Classes == nil {
		b.Classes = make(map[string]*Class)
	}
	if c, ok := b.Classes[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as class,first declared at:",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(c.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Funcs == nil {
		b.Funcs = make(map[string]*Function)
	}
	if f, ok := b.Funcs[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as function,first declared at:",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(f.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Consts == nil {
		b.Consts = make(map[string]*Const)
	}
	if c, ok := b.Consts[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as const,first declared at:",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(c.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Enums == nil {
		b.Enums = make(map[string]*Enum)
	}
	if e, ok := b.Enums[name]; ok {
		errmsg := fmt.Sprintf("%s name %s already declared as enum,first declared at:",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(e.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.EnumNames == nil {
		b.EnumNames = make(map[string]*EnumName)
	}
	if en, ok := b.EnumNames[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(en.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Lables == nil {
		b.Lables = make(map[string]*StatementLable)
	}
	if l, ok := b.Lables[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(l.Statement.Pos))
		return fmt.Errorf(errmsg)
	}
	if b.Types == nil {
		b.Types = make(map[string]*VariableType)
	}
	if t, ok := b.Types[name]; ok {
		errmsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
			errMsgPrefix(pos), name)
		errmsg += fmt.Sprintf("\t%s", errMsgPrefix(t.Pos))
		return fmt.Errorf(errmsg)
	}
	switch d.(type) {
	case *Class:
		if t := searchBuildIns(name); t != nil {
			return fmt.Errorf("%s '%s' is buildin", errMsgPrefix(pos), name)
		}
		b.Classes[name] = d.(*Class)
	case *Function:
		if t := searchBuildIns(name); t != nil {
			return fmt.Errorf("%s '%s' is buildin", errMsgPrefix(pos), name)
		}
		t := d.(*Function)
		if buildinFunctionsMap[t.Name] != nil {
			return fmt.Errorf("%s function named '%s' is buildin",
				errMsgPrefix(pos), name)
		}
		b.Funcs[name] = t
	case *Const:
		b.Consts[name] = d.(*Const)
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		b.Vars[name] = t
	case *Enum:
		if t := searchBuildIns(name); t != nil {
			return fmt.Errorf("%s '%s' is buildin", errMsgPrefix(pos), name)
		}
		e := d.(*Enum)
		b.Enums[name] = e
		for _, v := range e.Enums {
			err := b.insert(v.Name, v.Pos, v)
			if err != nil {
				return err
			}
		}
	case *EnumName:
		if t := searchBuildIns(name); t != nil {
			return fmt.Errorf("%s '%s' is buildin", errMsgPrefix(pos), name)
		}
		b.EnumNames[name] = d.(*EnumName)
	case *StatementLable:
		b.Lables[name] = d.(*StatementLable)
	case *VariableType:
		if t := searchBuildIns(name); t != nil {
			return fmt.Errorf("%s '%s' is buildin", errMsgPrefix(pos), name)
		}
		b.Types[name] = d.(*VariableType)
	default:
		panic("????????")
	}
	return nil
}
