package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"math"
)

type Block struct {
	Exits []*cg.Exit // for switch template
	/*
		should analyse at ast stage
	*/
	NotExecuteToLastStatement bool
	Defers                    []*StatementDefer
	Fn                        *Function
	IsFunctionBlock           bool // function block
	IsClassBlock              bool // class block
	Class                     *Class
	IsForBlock                bool // for top block
	IsSwitchBlock             bool // switch statement list block
	IsWhenBlock               bool // template swtich statement list block
	Pos                       *Pos
	EndPos                    *Pos
	Outer                     *Block
	InheritedAttribute        InheritedAttribute
	Statements                []*Statement
	Constants                 map[string]*Constant
	Functions                 map[string]*Function
	Classes                   map[string]*Class
	Enums                     map[string]*Enum
	EnumNames                 map[string]*EnumName
	Labels                    map[string]*StatementLabel
	TypeAliases               map[string]*Type
	Variables                 map[string]*Variable
	checkConstantsCalled      bool
}

func (this *Block) NameExists(name string) (interface{}, bool) {
	if this.Functions != nil {
		if t, ok := this.Functions[name]; ok {
			return t, true
		}
	}
	if this.Variables != nil {
		if t, ok := this.Variables[name]; ok {
			return t, true
		}
	}
	if this.Constants != nil {
		if t, ok := this.Constants[name]; ok {
			return t, true
		}
	}
	if this.EnumNames != nil {
		if t, ok := this.EnumNames[name]; ok {
			return t, true
		}
	}
	if this.Classes != nil {
		if t, ok := this.Classes[name]; ok {
			return t, true
		}
	}
	if this.Enums != nil {
		if t, ok := this.Enums[name]; ok {
			return t, true
		}
	}
	if this.TypeAliases != nil {
		if t, ok := this.TypeAliases[name]; ok {
			return t, true
		}
	}
	if this.Labels != nil { // should be useless
		if t, ok := this.Labels[name]; ok {
			return t, true
		}
	}
	return nil, false
}

/*
	search label
*/
func (this *Block) searchLabel(name string) *StatementLabel {
	outer := this
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
func (this *Block) searchType(name string) interface{} {
	bb := this
	for bb != nil {
		if bb.Classes != nil {
			if t, ok := bb.Classes[name]; ok {
				t.Used = true
				return t
			}
		}
		if bb.Enums != nil {
			if t, ok := bb.Enums[name]; ok {
				t.Used = true
				return t
			}
		}
		if bb.TypeAliases != nil {
			if t, ok := bb.TypeAliases[name]; ok {
				return t
			}
		}
		if bb.IsFunctionBlock && bb.Fn != nil {
			if bb.Fn.parameterTypes != nil {
				if t := bb.Fn.parameterTypes[name]; t != nil {
					return t
				}
			}
		}
		bb = bb.Outer
	}
	return nil
}

func (this *Block) identifierIsWhat(d interface{}) string {
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
	case *Package:
		return "package" // impossible, no big deal
	default:
		return "new item , call author"
	}
}

/*
	search identifier
*/
func (this *Block) searchIdentifier(from *Pos, name string, isCaptureVar *bool) (interface{}, error) {
	if this.Functions != nil {
		if t, ok := this.Functions[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if this.Variables != nil {
		if t, ok := this.Variables[name]; ok {
			return t, nil
		}
	}
	if this.Constants != nil {
		if t, ok := this.Constants[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if this.EnumNames != nil {
		if t, ok := this.EnumNames[name]; ok {
			t.Enum.Used = true
			return t, nil
		}
	}
	if this.Enums != nil {
		if t, ok := this.Enums[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if this.Classes != nil {
		if t, ok := this.Classes[name]; ok {
			t.Used = true
			return t, nil
		}
	}
	if this.TypeAliases != nil {
		if t, ok := this.TypeAliases[name]; ok {
			return t, nil
		}
	}
	if this.IsFunctionBlock && this.Fn != nil {
		if this.Fn.parameterTypes != nil {
			if t := this.Fn.parameterTypes[name]; t != nil {
				return t, nil
			}
		}
	}
	// search closure
	if this.InheritedAttribute.Function != nil {
		v := this.InheritedAttribute.Function.Closure.Search(name)
		if v != nil {
			return v, nil
		}
	}
	if this.IsFunctionBlock &&
		len(this.InheritedAttribute.Function.parameterTypes) > 0 {
		return searchBuildIns(name), nil
	}
	if this.IsFunctionBlock &&
		name == ThisPointerName {
		return nil, nil
	}
	if this.Outer == nil {
		return searchBuildIns(name), nil
	}
	t, err := this.Outer.searchIdentifier(from, name, isCaptureVar) // search by outer block
	if err != nil {
		return t, err
	}
	if t != nil { //
		switch t.(type) {
		case *Variable:
			v := t.(*Variable)
			if v.IsGlobal == false { // not a global variable
				if this.IsFunctionBlock &&
					this.InheritedAttribute.Function.IsGlobal == false {
					this.InheritedAttribute.Function.Closure.InsertVar(from, v)
					if isCaptureVar != nil {
						*isCaptureVar = true
					}
				}
				//cannot search variable from class body
				if this.IsClassBlock {
					return nil, fmt.Errorf("%s trying to access variable '%s' from class",
						from.ErrMsgPrefix(), name)
				}
			}
		case *Function:
			f := t.(*Function)
			if f.IsGlobal == false {
				if this.IsClassBlock {
					this.Class.closure.InsertFunction(from, f)
				}
				if this.IsFunctionBlock {
					this.Fn.Closure.InsertFunction(from, f)
				}
			}
		}
	}
	return t, nil
}

func (this *Block) inherit(father *Block) {
	if this.Outer != nil {
		return
	}
	if this == father {
		panic("inherit from self")
	}
	this.InheritedAttribute = father.InheritedAttribute
	this.Outer = father
	if this.IsFunctionBlock || this.IsClassBlock {
		this.InheritedAttribute.ForBreak = nil
		this.InheritedAttribute.ForContinue = nil
		this.InheritedAttribute.StatementOffset = 0
		this.InheritedAttribute.IsConstructionMethod = false
		this.InheritedAttribute.ClassMethod = nil
		this.InheritedAttribute.Defer = nil
	}
}

func (this *Block) checkUnUsed() (es []error) {
	if common.CompileFlags.DisableCheckUnUse {
		return nil
	}
	es = []error{}
	for _, v := range this.Constants {
		if v.Used ||
			v.IsGlobal {
			continue
		}
		es = append(es, fmt.Errorf("%s const '%s' has declared,but not used",
			v.Pos.ErrMsgPrefix(), v.Name))
	}
	for _, v := range this.Enums {
		if v.Used ||
			v.IsGlobal {
			continue
		}
		es = append(es, fmt.Errorf("%s enum '%s' has declared,but not used",
			v.Pos.ErrMsgPrefix(), v.Name))
	}
	for _, v := range this.Classes {
		if v.Used ||
			v.IsGlobal {
			continue
		}
		es = append(es, fmt.Errorf("%s class '%s' has declared,but not used",
			v.Pos.ErrMsgPrefix(), v.Name))
	}
	for _, v := range this.Functions {
		if v.Used ||
			v.IsGlobal {
			continue
		}
		es = append(es, fmt.Errorf("%s function '%s' has declared,but not used",
			v.Pos.ErrMsgPrefix(), v.Name))
	}
	for _, v := range this.Labels {
		if v.Used {
			continue
		}
		es = append(es, fmt.Errorf("%s enum '%s' has declared,but not used",
			v.Pos.ErrMsgPrefix(), v.Name))
	}
	for _, v := range this.Variables {
		if v.Used ||
			v.IsGlobal ||
			v.IsFunctionParameter ||
			v.Name == ThisPointerName ||
			v.IsReturn {
			continue
		}
		es = append(es, fmt.Errorf("%s variable '%s' has declared,but not used",
			v.Pos.ErrMsgPrefix(), v.Name))
	}
	return es
}

func (this *Block) check() []error {
	errs := []error{}
	for k, s := range this.Statements {
		if s.isStaticFieldDefaultValue {
			// no need to check
			// compile auto statement , checked before
			continue
		}
		this.InheritedAttribute.StatementOffset = k
		errs = append(errs, s.check(this)...)
		if PackageBeenCompile.shouldStop(errs) {
			return errs
		}
	}
	errs = append(errs, this.checkUnUsed()...)
	return errs
}

func (this *Block) checkConstants() []error {
	if this.checkConstantsCalled {
		return []error{}
	}
	this.checkConstantsCalled = true
	errs := make([]error, 0)
	for _, c := range this.Constants {
		if err := this.nameIsValid(c.Name, c.Pos); err != nil {
			errs = append(errs, err)
			delete(this.Constants, c.Name)
			continue
		}
		err := checkConst(this, c)
		if err != nil {
			errs = append(errs, err)
		}
		if err != nil && c.Type == nil {
			delete(this.Constants, c.Name)
		}
	}
	return errs
}

func (this *Block) checkNameExist(name string, pos *Pos) error {
	if this.Variables == nil {
		this.Variables = make(map[string]*Variable)
	}
	if v, ok := this.Variables[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as variable,first declared at:\n",
			pos.ErrMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", v.Pos.ErrMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if this.Classes == nil {
		this.Classes = make(map[string]*Class)
	}
	if c, ok := this.Classes[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as class,first declared at:\n",
			pos.ErrMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", c.Pos.ErrMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if this.Functions == nil {
		this.Functions = make(map[string]*Function)
	}
	if f, ok := this.Functions[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as function,first declared at:\n",
			pos.ErrMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", f.Pos.ErrMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if this.Constants == nil {
		this.Constants = make(map[string]*Constant)
	}
	if c, ok := this.Constants[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as const,first declared at:\n",
			pos.ErrMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", c.Pos.ErrMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if this.EnumNames == nil {
		this.EnumNames = make(map[string]*EnumName)
	}
	if en, ok := this.EnumNames[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:\n",
			pos.ErrMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", en.Pos.ErrMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if this.TypeAliases == nil {
		this.TypeAliases = make(map[string]*Type)
	}
	if t, ok := this.TypeAliases[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:\n",
			pos.ErrMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", t.Pos.ErrMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	if this.Enums == nil {
		this.Enums = make(map[string]*Enum)
	}
	if e, ok := this.Enums[name]; ok {
		errMsg := fmt.Sprintf("%s name %s already declared as enum,first declared at:\n",
			pos.ErrMsgPrefix(), name)
		errMsg += fmt.Sprintf("\t%s", e.Pos.ErrMsgPrefix())
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (this *Block) nameIsValid(name string, pos *Pos) error {
	if name == "" {
		return fmt.Errorf(`%s "" is not a valid name`, pos.ErrMsgPrefix())
	}
	if name == ThisPointerName {
		return fmt.Errorf("%s '%s' already been taken", pos.ErrMsgPrefix(), ThisPointerName)
	}
	if name == "_" {
		return fmt.Errorf("%s '%s' is not a valid name", pos.ErrMsgPrefix(), name)
	}
	if isMagicIdentifier(name) {
		return fmt.Errorf("%s '%s' is not a magic identifier", pos.ErrMsgPrefix(), name)
	}
	if searchBuildIns(name) != nil {
		return fmt.Errorf("%s '%s' is buildin", pos.ErrMsgPrefix(), name)
	}
	return nil
}

func (this *Block) Insert(name string, pos *Pos, d interface{}) error {
	if err := this.nameIsValid(name, pos); err != nil {
		return err
	}
	// handle label
	if label, ok := d.(*StatementLabel); ok && label != nil {
		if this.Labels == nil {
			this.Labels = make(map[string]*StatementLabel)
		}
		if l, ok := this.Labels[name]; ok {
			errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
				pos.ErrMsgPrefix(), name)
			errMsg += fmt.Sprintf("\t%s", l.Statement.Pos.ErrMsgPrefix())
			return fmt.Errorf(errMsg)
		}
		this.Labels[name] = label
		return nil
	}
	err := this.checkNameExist(name, pos)
	if err != nil {
		return err
	}
	switch d.(type) {
	case *Class:
		this.Classes[name] = d.(*Class)
	case *Function:
		t := d.(*Function)
		if buildInFunctionsMap[t.Name] != nil {
			return fmt.Errorf("%s function named '%s' is buildin",
				pos.ErrMsgPrefix(), name)
		}
		this.Functions[name] = t
	case *Constant:
		this.Constants[name] = d.(*Constant)
	case *Variable:
		t := d.(*Variable)
		t.LocalValOffset = math.MaxUint16 // overflow
		this.Variables[name] = t
	case *Enum:
		e := d.(*Enum)
		this.Enums[name] = e
		for _, v := range e.Enums {
			err := this.Insert(v.Name, v.Pos, v)
			if err != nil {
				return err
			}
		}
	case *EnumName:
		this.EnumNames[name] = d.(*EnumName)
	case *Type:
		this.TypeAliases[name] = d.(*Type)
	}
	return nil
}
