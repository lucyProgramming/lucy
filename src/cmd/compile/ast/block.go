package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Block struct {
	Exits []*cg.Exit // for switch template
	/*
		should analyse at ast stage
	*/
	NotExecuteToLastStatement       bool
	Defers                          []*StatementDefer
	IsGlobalVariableDefinitionBlock bool
	IsFunctionBlock                 bool // function block
	IsClassBlock                    bool // class block
	IsForBlock                      bool // for top block
	IsSwitchBlock                   bool // switch statement list block
	IsSwitchTemplateBlock           bool // template swtich statement list block
	Pos                             *Pos
	EndPos                          *Pos
	Outer                           *Block
	InheritedAttribute              InheritedAttribute
	Statements                      []*Statement
	Constants                       map[string]*Constant
	Functions                       map[string]*Function
	Classes                         map[string]*Class
	Enums                           map[string]*Enum
	EnumNames                       map[string]*EnumName
	Labels                          map[string]*StatementLabel
	TypeAliases                     map[string]*Type
	Variables                       map[string]*Variable
	ClosureFunctions                map[string]*Function //in "Functions" too
	checkConstantsCalled            bool
}

func (b *Block) NameExists(name string) (interface{}, bool) {
	if b.Functions != nil {
		if t, ok := b.Functions[name]; ok {
			return t, true
		}
	}
	if b.Variables != nil {
		if t, ok := b.Variables[name]; ok {
			return t, true
		}
	}
	if b.Constants != nil {
		if t, ok := b.Constants[name]; ok {
			return t, true
		}
	}
	if b.EnumNames != nil {
		if t, ok := b.EnumNames[name]; ok {
			return t, true
		}
	}
	if b.Classes != nil {
		if t, ok := b.Classes[name]; ok {
			return t, true
		}
	}
	if b.Enums != nil {
		if t, ok := b.Enums[name]; ok {
			return t, true
		}
	}
	if b.TypeAliases != nil {
		if t, ok := b.TypeAliases[name]; ok {
			return t, true
		}
	}
	if b.Labels != nil { // should be useless
		if t, ok := b.Labels[name]; ok {
			return t, true
		}
	}
	return nil, false
}

/*
	search label
*/
func (b *Block) searchLabel(name string) *StatementLabel {
	outer := b
	for {
		if outer.Labels != nil {
			if l, ok := outer.Labels[name]; ok {
				l.Used = true
				return l
			}
		}
		if outer.IsFunctionBlock {
			return nil
		}
		outer = outer.Outer
	}
	return nil
}

/*
	search type
*/
func (b *Block) searchType(name string) interface{} {
	outer := b
	for outer != nil {
		if outer.Classes != nil {
			if t, ok := outer.Classes[name]; ok {
				t.Used = true
				return t
			}
			if t, ok := outer.Enums[name]; ok {
				t.Used = true
				return t
			}
			if t, ok := outer.TypeAliases[name]; ok {
				return t
			}
		}
		outer = outer.Outer
	}
	return nil
}

func (b *Block) identifierIsWhat(d interface{}) string {
	switch d.(type) {
	case *Function:
		return "function"
	case *Variable:
		return "variable"
	case *Constant:
		return "constant"
	case *EnumName:
		return "enum name"
	case *Enum:
		return "enum"
	case *Class:
		return "class"
	case *Type:
		return "type alias"
	default:
		return "new item,call author"
	}
}

/*
	search identifier
*/
func (b *Block) searchIdentifier(from *Pos, name string) (interface{}, error) {
	if b.Functions != nil {
		if t, ok := b.Functions[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if b.Variables != nil {
		if t, ok := b.Variables[name]; ok {
			return t, nil
		}
	}
	if b.Constants != nil {
		if t, ok := b.Constants[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if b.EnumNames != nil {
		if t, ok := b.EnumNames[name]; ok {
			t.Enum.Used = true
			return t, nil
		}
	}
	if b.Enums != nil {
		if t, ok := b.Enums[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if b.Classes != nil {
		if t, ok := b.Classes[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if b.TypeAliases != nil {
		if t, ok := b.TypeAliases[name]; ok {
			return t, nil
		}
	}
	// search closure
	if b.InheritedAttribute.Function != nil {
		v := b.InheritedAttribute.Function.Closure.Search(name)
		if v != nil {
			return v, nil
		}
	}
	if b.IsFunctionBlock &&
		len(b.InheritedAttribute.Function.parameterTypes) > 0 {
		return searchBuildIns(name), nil
	}
	if b.Outer == nil {
		return searchBuildIns(name), nil
	}
	t, err := b.Outer.searchIdentifier(from, name) // search by outer block
	if err != nil {
		return t, err
	}
	if t != nil { //
		switch t.(type) {
		case *Variable:
			v := t.(*Variable)
			if v.IsGlobal == false { // not a global variable
				if b.IsFunctionBlock &&
					b.InheritedAttribute.Function.IsGlobal == false {
					if v.Name == THIS {
						return nil, fmt.Errorf("%s capture '%s' not allow",
							from.errMsgPrefix(), name) // capture this not allow
					}
					b.InheritedAttribute.Function.Closure.InsertVar(v)
				}
				//cannot search variable from class body
				if b.InheritedAttribute.Class != nil && b.IsClassBlock {
					return nil, fmt.Errorf("%s trying to access variable '%s' from class",
						from.errMsgPrefix(), name)
				}
			}
		case *Function:
			f := t.(*Function)
			if f.IsGlobal == false {
				if b.IsFunctionBlock {
					b.InheritedAttribute.Function.Closure.InsertFunction(f)
				}
				if b.IsClassBlock && f.IsClosureFunction {
					return nil, fmt.Errorf("%s trying to access closure function '%s' from class",
						from.errMsgPrefix(), name)
				}
			}
		}
	}
	return t, nil
}

func (b *Block) inherit(father *Block) {
	if b.Outer != nil {
		return
	}
	if b == father {
		panic("inherit from self")
	}
	b.InheritedAttribute = father.InheritedAttribute
	b.Outer = father
}

func (b *Block) checkUnUsed() (es []error) {
	es = []error{}
	for _, v := range b.Variables {
		if v.Used ||
			v.IsGlobal ||
			v.IsFunctionParameter ||
			v.Name == THIS ||
			v.IsReturn {
			continue
		}
		es = append(es, fmt.Errorf("%s variable '%s' has declared,but not used",
			v.Pos.errMsgPrefix(), v.Name))
	}

	return es
}

func (b *Block) checkStatementsAndUnused() []error {
	errs := []error{}
	for k, s := range b.Statements {
		if s.isStaticFieldDefaultValue {
			// no need to check
			// compile auto statement , checked before
			continue
		}
		b.InheritedAttribute.StatementOffset = k
		errs = append(errs, s.check(b)...)
		if PackageBeenCompile.shouldStop(errs) {
			return errs
		}
	}
	errs = append(errs, b.checkUnUsed()...)
	return errs
}

func (b *Block) checkConstants() []error {
	if b.checkConstantsCalled {
		return []error{}
	}
	b.checkConstantsCalled = true
	errs := make([]error, 0)
	for _, c := range b.Constants {
		if err := b.nameIsValid(c.Name, c.Pos); err != nil {
			errs = append(errs, err)
			delete(b.Constants, c.Name)
			continue
		}
		err := checkConst(b, c)
		if err != nil {
			errs = append(errs, err)
		}
		if err != nil && c.Type == nil {
			delete(b.Constants, c.Name)
		}
	}
	return errs
}

func (b *Block) checkNameExist(name string, pos *Pos) error {
	if b.Variables == nil {
		b.Variables = make(map[string]*Variable)
	}
	if v, ok := b.Variables[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as variable,first declared at:\n",
			pos.errMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", v.Pos.errMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if b.Classes == nil {
		b.Classes = make(map[string]*Class)
	}
	if c, ok := b.Classes[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as class,first declared at:\n",
			pos.errMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", c.Pos.errMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if b.Functions == nil {
		b.Functions = make(map[string]*Function)
	}
	if f, ok := b.Functions[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as function,first declared at:\n",
			pos.errMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", f.Pos.errMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if b.Constants == nil {
		b.Constants = make(map[string]*Constant)
	}
	if c, ok := b.Constants[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as const,first declared at:\n",
			pos.errMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", c.Pos.errMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if b.EnumNames == nil {
		b.EnumNames = make(map[string]*EnumName)
	}
	if en, ok := b.EnumNames[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:\n",
			pos.errMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", en.Pos.errMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if b.TypeAliases == nil {
		b.TypeAliases = make(map[string]*Type)
	}
	if t, ok := b.TypeAliases[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:\n",
			pos.errMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", t.Pos.errMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if b.Enums == nil {
		b.Enums = make(map[string]*Enum)
	}
	if e, ok := b.Enums[name]; ok {
		errMsg := fmt.Sprintf("%s name %s already declared as enum,first declared at:\n",
			pos.errMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", e.Pos.errMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (b *Block) nameIsValid(name string, pos *Pos) error {
	if name == "" {
		return fmt.Errorf(`%s "" is not a valid name`, pos.errMsgPrefix())
	}
	if name == THIS {
		return fmt.Errorf("%s '%s' already been taken", pos.errMsgPrefix(), THIS)
	}
	if name == "_" {
		return fmt.Errorf("%s '%s' is not a valid name", pos.errMsgPrefix(), name)
	}
	if isMagicIdentifier(name) {
		return fmt.Errorf("%s '%s' is not a magic identifier", pos.errMsgPrefix(), name)
	}
	if searchBuildIns(name) != nil {
		return fmt.Errorf("%s '%s' is buildin", pos.errMsgPrefix(), name)
	}
	return nil
}

func (b *Block) Insert(name string, pos *Pos, d interface{}) error {
	if err := b.nameIsValid(name, pos); err != nil {
		return err
	}
	if v, ok := d.(*Variable); ok &&
		b.InheritedAttribute.Function.isGlobalVariableDefinition {
		b := PackageBeenCompile.Block
		err := b.checkNameExist(name, pos)
		if err != nil {
			return err
		}
		b.Variables[name] = v
		v.IsGlobal = true // it`s global
		return nil
	}
	// handle label
	if label, ok := d.(*StatementLabel); ok && label != nil {
		if b.Labels == nil {
			b.Labels = make(map[string]*StatementLabel)
		}
		if l, ok := b.Labels[name]; ok {
			errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
				pos.errMsgPrefix(), name)
			errMsg += fmt.Sprintf("\t%s", l.Statement.Pos.errMsgPrefix())
			return fmt.Errorf(errMsg)
		}
		b.Labels[name] = label
		return nil
	}
	err := b.checkNameExist(name, pos)
	if err != nil {
		return err
	}
	switch d.(type) {
	case *Class:
		b.Classes[name] = d.(*Class)
	case *Function:
		t := d.(*Function)
		if buildInFunctionsMap[t.Name] != nil {
			return fmt.Errorf("%s function named '%s' is buildin",
				pos.errMsgPrefix(), name)
		}
		b.Functions[name] = t
	case *Constant:
		b.Constants[name] = d.(*Constant)
	case *Variable:
		t := d.(*Variable)
		b.Variables[name] = t
	case *Enum:
		e := d.(*Enum)
		b.Enums[name] = e
		for _, v := range e.Enums {
			err := b.Insert(v.Name, v.Pos, v)
			if err != nil {
				return err
			}
		}
	case *EnumName:
		b.EnumNames[name] = d.(*EnumName)
	case *Type:
		b.TypeAliases[name] = d.(*Type)
	}
	return nil
}
