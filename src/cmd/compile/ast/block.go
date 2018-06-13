package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
)

type Block struct {
	DeadEnding                 bool // will execute to last statement
	Defers                     []*Defer
	isGlobalVariableDefinition bool
	IsFunctionTopBlock         bool
	IsClassBlock               bool
	Pos                        *Pos
	EndPos                     *Pos
	Outer                      *Block
	InheritedAttribute         InheritedAttribute
	Statements                 []*Statement
	Constants                  map[string]*Constant
	Functions                  map[string]*Function
	Classes                    map[string]*Class
	Enums                      map[string]*Enum
	EnumNames                  map[string]*EnumName
	Labels                     map[string]*StatementLabel
	Types                      map[string]*VariableType
	Variables                  map[string]*VariableDefinition
	ClosureFunctions           map[string]*Function //in "Functions" too
}

func (b *Block) HaveVariableDefinition() bool {
	if b.ClosureFunctions == nil && b.Variables == nil {
		return false
	}
	return len(b.ClosureFunctions) > 0 || len(b.Variables) > 0
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
	if b.Types != nil {
		if t, ok := b.Types[name]; ok {
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
	bb := b
	for bb != nil {
		if bb.Labels != nil {
			if l, ok := bb.Labels[name]; ok {
				return l
			}
		}
		if bb.IsFunctionTopBlock {
			return nil
		}
		bb = bb.Outer
	}
	return nil
}

/*
	search type
*/
func (b *Block) searchType(name string) interface{} {
	bb := b
	for bb != nil {
		if bb.Classes != nil {
			if t, ok := bb.Classes[name]; ok {
				return t
			}
			if t, ok := bb.Enums[name]; ok {
				return t
			}
			if t, ok := bb.Types[name]; ok {
				return t
			}
		}
		bb = bb.Outer
	}
	return nil
}

/*
	search rightValue
*/
func (b *Block) SearchByName(name string) (interface{}, error) {
	if b.Functions != nil {
		if t, ok := b.Functions[name]; ok {
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
			return t, nil
		}
	}
	if b.EnumNames != nil {
		if t, ok := b.EnumNames[name]; ok {
			return t, nil
		}
	}
	if b.Classes != nil {
		if t, ok := b.Classes[name]; ok {
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
	if b.IsFunctionTopBlock &&
		len(b.InheritedAttribute.Function.TypeParameters) > 0 { // this is a template function
		return searchBuildIns(name), nil
	}
	if b.Outer == nil {
		return searchBuildIns(name), nil
	}
	t, err := b.Outer.SearchByName(name) // search by outter block
	if err != nil {
		return t, err
	}
	if t != nil { //
		if _, ok := t.(*VariableDefinition); ok && b.IsFunctionTopBlock &&
			len(b.InheritedAttribute.Function.TypeParameters) > 0 { // template function
			return nil, nil
		}
		if v, ok := t.(*VariableDefinition); ok && v.IsGlobal == false { // not a global variable
			if b.IsFunctionTopBlock &&
				b.InheritedAttribute.Function.IsGlobal == false { // 	b.InheritedAttribute.Function.IsGlobal == false  no need to check
				if v.Name == THIS {
					return nil, nil // capture this not allow
				}
				b.InheritedAttribute.Function.Closure.InsertVar(v)
			}
			//cannot search variable from class body
			if b.InheritedAttribute.Class != nil && b.IsClassBlock {
				return nil, nil //
			}
		}
		// if it is a function
		if f, ok := t.(*Function); ok && f.IsGlobal == false {
			if b.IsFunctionTopBlock {
				b.InheritedAttribute.Function.Closure.InsertFunction(f)
			}
			if b.IsClassBlock && f.IsClosureFunction {
				return nil, fmt.Errorf("trying to access closure function '%s' from class", name)
			}
		}
	}
	return t, nil
}

func (b *Block) inherit(father *Block) {
	if b == father {
		panic("inherit  from self")
	}
	b.InheritedAttribute = father.InheritedAttribute
	b.Outer = father

}

func (b *Block) checkUnUsedVariable() (es []error) {
	es = []error{}
	for _, v := range b.Variables {
		if v.Used ||
			v.IsGlobal ||
			v.IsFunctionParameter ||
			v.IsFunctionReturnVar ||
			v.Name == THIS {
			continue
		}
		es = append(es, fmt.Errorf("%s variable '%s' has declared,but not used",
			errMsgPrefix(v.Pos), v.Name))
	}
	return es
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
	if common.CompileFlags.DisableCheckUnUsedVariable == false {
		errs = append(errs, b.checkUnUsedVariable()...)
	}
	return errs
}

func (b *Block) checkConst() []error {
	errs := make([]error, 0)
	for _, c := range b.Constants {
		if c.Name == NO_NAME_IDENTIFIER {
			err := fmt.Errorf("%s '%s' is not a valid name",
				errMsgPrefix(c.Pos), c.Name)
			errs = append(errs, err)
			delete(b.Constants, c.Name)
			continue
		}
		err := checkConst(b, c, &errs)
		if err != nil && c.Typ == nil {
			delete(b.Constants, c.Name)
		}
	}
	return errs
}

func (b *Block) Insert(name string, pos *Pos, d interface{}) error {
	return b.insert(name, pos, d)
}
func (b *Block) insert(name string, pos *Pos, d interface{}) error {
	//fmt.Println(name, pos)
	// global var insert into block
	if v, ok := d.(*VariableDefinition); ok && b.InheritedAttribute.Function.isGlobalVariableDefinition {
		b := PackageBeenCompile.Block
		if vv, ok := b.Variables[name]; ok {
			errMsg := fmt.Sprintf("%s name '%s' already declared as variable,first declared at:\n",
				errMsgPrefix(pos), name)
			errMsg += fmt.Sprintf("\t%s", errMsgPrefix(vv.Pos))
			return fmt.Errorf(errMsg)
		}
		b.Variables[name] = v
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
	if b.Variables == nil {
		b.Variables = make(map[string]*VariableDefinition)
	}
	if v, ok := b.Variables[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as variable,first declared at:\n",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(v.Pos))
		return fmt.Errorf(errMsg)
	}
	if b.Classes == nil {
		b.Classes = make(map[string]*Class)
	}
	if c, ok := b.Classes[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as class,first declared at:",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(c.Pos))
		return fmt.Errorf(errMsg)
	}
	if b.Functions == nil {
		b.Functions = make(map[string]*Function)
	}
	if f, ok := b.Functions[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as function,first declared at:",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(f.Pos))
		return fmt.Errorf(errMsg)
	}
	if b.Constants == nil {
		b.Constants = make(map[string]*Constant)
	}
	if c, ok := b.Constants[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as const,first declared at:",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(c.Pos))
		return fmt.Errorf(errMsg)
	}
	if b.Enums == nil {
		b.Enums = make(map[string]*Enum)
	}
	if e, ok := b.Enums[name]; ok {
		errMsg := fmt.Sprintf("%s name %s already declared as enum,first declared at:",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(e.Pos))
		return fmt.Errorf(errMsg)
	}
	if b.EnumNames == nil {
		b.EnumNames = make(map[string]*EnumName)
	}
	if en, ok := b.EnumNames[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(en.Pos))
		return fmt.Errorf(errMsg)
	}
	if b.Labels == nil {
		b.Labels = make(map[string]*StatementLabel)
	}
	if l, ok := b.Labels[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(l.Statement.Pos))
		return fmt.Errorf(errMsg)
	}
	if b.Types == nil {
		b.Types = make(map[string]*VariableType)
	}
	if t, ok := b.Types[name]; ok {
		errMsg := fmt.Sprintf("%s name '%s' already declared as enumName,first declared at:",
			errMsgPrefix(pos), name)
		errMsg += fmt.Sprintf("\t%s", errMsgPrefix(t.Pos))
		return fmt.Errorf(errMsg)
	}
	// name exists in buildin, not allow
	if lucyBuildInPackage != nil {
		if _, exits := lucyBuildInPackage.Block.NameExists(name); exits {
			return fmt.Errorf("%s '%s' is buildin", errMsgPrefix(pos), name)
		}
	}
	switch d.(type) {
	case *Class:
		b.Classes[name] = d.(*Class)
	case *Function:
		t := d.(*Function)
		if buildInFunctionsMap[t.Name] != nil {
			return fmt.Errorf("%s function named '%s' is buildin",
				errMsgPrefix(pos), name)
		}
		b.Functions[name] = t
	case *Constant:
		b.Constants[name] = d.(*Constant)
	case *VariableDefinition:
		t := d.(*VariableDefinition)
		b.Variables[name] = t
	case *Enum:
		e := d.(*Enum)
		b.Enums[name] = e
		for _, v := range e.Enums {
			err := b.insert(v.Name, v.Pos, v)
			if err != nil {
				return err
			}
		}
	case *EnumName:
		b.EnumNames[name] = d.(*EnumName)
	case *StatementLabel:
		b.Labels[name] = d.(*StatementLabel)
	case *VariableType:
		b.Types[name] = d.(*VariableType)
	default:
		panic(d) // panic d
	}
	return nil
}
