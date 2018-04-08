package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Function struct {
	HaveDefaultValue               bool
	DefaultValueStartAt            int
	IsClosureFunction              bool
	ClassMethod                    *cg.MethodHighLevel
	VarOffSetForClosure            uint16
	isGlobalVariableDefinition     bool
	isPackageBlockFunction         bool
	callchecker                    CallChecker // used in build function
	IsGlobal                       bool
	IsBuildin                      bool
	Used                           bool
	AccessFlags                    uint16 // public private or protected
	Typ                            *FunctionType
	ClosureVars                    Closure
	Name                           string // if name is nil string,means no name function
	Block                          *Block
	Pos                            *Pos
	Descriptor                     string
	VarOffset                      uint16
	AutoVarForException            *AutoVarForException
	AutoVarForReturnBecauseOfDefer *AutoVarForReturnBecauseOfDefer
	AutoVarForMultiReturn          *AutoVarForMultiReturn
	OffsetDestinations             []*uint16
}

func (f *Function) MkAutoVarForReturnBecauseOfDefer() {
	if f.NoReturnValue() {
		return
	}
	if f.AutoVarForReturnBecauseOfDefer != nil {
		return
	}
	t := &AutoVarForReturnBecauseOfDefer{}
	f.AutoVarForReturnBecauseOfDefer = t
	t.ExceptionIsNotNilWhenEnter = f.VarOffset
	f.VarOffset++
	f.OffsetDestinations = append(f.OffsetDestinations, &t.ExceptionIsNotNilWhenEnter)
	if len(f.Typ.ReturnList) > 1 {
		t.MultiValueOffset = f.VarOffset
		f.OffsetDestinations = append(f.OffsetDestinations, &t.MultiValueOffset)
		f.VarOffset++
		t.IfReachBotton = f.VarOffset
		f.VarOffset++
		f.OffsetDestinations = append(f.OffsetDestinations, &t.IfReachBotton)
	}
}

type AutoVarForReturnBecauseOfDefer struct {
	/*
		flag is 1 means there is exception,but handled
	*/
	ExceptionIsNotNilWhenEnter uint16
	/*
		for multi return value
	*/
	MultiValueOffset uint16
	IfReachBotton    uint16
}

func (f *Function) NoReturnValue() bool {
	return len(f.Typ.ReturnList) == 0 || f.Typ.ReturnList[0].Typ.Typ == VARIABLE_TYPE_BOOL
}

type AutoVarForException struct {
	Offset uint16
}

func (f *Function) mkAutoVarForMultiReturn() {
	if f.AutoVarForMultiReturn != nil {
		return
	}
	t := &AutoVarForMultiReturn{}
	t.Offset = f.VarOffset
	f.AutoVarForMultiReturn = t
	f.OffsetDestinations = append(f.OffsetDestinations, &t.Offset)
	f.VarOffset++
}

type AutoVarForMultiReturn struct {
	Offset uint16
}

func (f *Function) readableMsg() string {
	s := "fn " + f.Name + " ("
	for k, v := range f.Typ.ParameterList {
		s += v.Name + " " + v.Typ.TypeString()
		if k != len(f.Typ.ParameterList)-1 {
			s += ","
		}
	}
	s += ") "
	if len(f.Typ.ReturnList) > 0 {
		s += "-> ( "
		for k, v := range f.Typ.ReturnList {
			s += v.Name + " " + v.Typ.TypeString()
			if k != len(f.Typ.ReturnList)-1 {
				s += ","
			}
		}
		s += " )"
	}
	return s
}

func (f *Function) checkBlock(errs *[]error) {
	if f.Typ != nil {
		f.mkLastRetrunStatement()
		*errs = append(*errs, f.Block.check()...)
	}
}

func (f *Function) check(b *Block) []error {
	errs := make([]error, 0)
	f.Block.inherite(b)
	f.Block.InheritedAttribute.Function = f
	f.checkParaMeterAndRetuns(&errs)
	f.checkBlock(&errs)
	return errs
}

func (f *Function) mkLastRetrunStatement() {
	if len(f.Block.Statements) == 0 ||
		(f.Block.Statements[len(f.Block.Statements)-1].Typ != STATEMENT_TYPE_RETURN && f.Block.Statements[len(f.Block.Statements)-1].Typ != STATEMENT_TYPE_SKIP) {
		s := &StatementReturn{}
		f.Block.Statements = append(f.Block.Statements, &Statement{Typ: STATEMENT_TYPE_RETURN, StatementReturn: s})
	}
}

func (f *Function) checkParaMeterAndRetuns(errs *[]error) {
	if f.Name == MAIN_FUNCTION_NAME {
		errFunc := func() {
			*errs = append(*errs, fmt.Errorf("%s function %s expect declared as 'main(args []string)'", errMsgPrefix(f.Pos), MAIN_FUNCTION_NAME))
		}
		if len(f.Typ.ParameterList) != 1 {
			errFunc()
		} else { //
			if f.Typ.ParameterList[0].Typ.Typ == VARIABLE_TYPE_ARRAY &&
				f.Typ.ParameterList[0].Typ.ArrayType.Typ == VARIABLE_TYPE_STRING {
				err := f.Block.insert(f.Typ.ParameterList[0].Name, f.Typ.ParameterList[0].Pos, f.Typ.ParameterList[0])
				if err != nil {
					*errs = append(*errs, err)
				}
			} else {
				errFunc()
			}
		}
		return
	}
	if f.Typ != nil {
		var err error
		for k, v := range f.Typ.ParameterList {
			v.IsFunctionParameter = true
			err = v.Typ.resolve(f.Block)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s %s", errMsgPrefix(v.Pos), err.Error()))
			}
			err = f.Block.insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
			if f.HaveDefaultValue && v.Expression == nil {
				*errs = append(*errs, fmt.Errorf("%s expect default value", errMsgPrefix(v.Pos)))
				continue
			}
			if v.Expression != nil {
				if f.HaveDefaultValue == false {
					f.DefaultValueStartAt = k
				}
				f.HaveDefaultValue = true
				ts, es := v.Expression.check(f.Block)
				if errsNotEmpty(es) {
					*errs = append(*errs, es...)
				}
				t, err := v.Expression.mustBeOneValueContext(ts)
				if err != nil {
					*errs = append(*errs, err)
				}
				if t != nil {
					if v.Typ.Equal(t) == false {
						*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
							errMsgPrefix(v.Expression.Pos), t.TypeString(), v.Typ.TypeString()))
						continue
					}
				}
				if v.Expression.IsLiteral() == false {
					*errs = append(*errs, fmt.Errorf("%s default value must be literal",
						errMsgPrefix(v.Expression.Pos), t.TypeString(), v.Typ.TypeString()))
				}
			}
		}
		//handler return
		for _, v := range f.Typ.ReturnList {
			err = v.Typ.resolve(f.Block)
			if err != nil {
				*errs = append(*errs, err)
			}
			err = f.Block.insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s err:%v", errMsgPrefix(v.Pos), err))
			}
			if v.Expression == nil {
				v.Expression = v.Typ.mkDefaultValueExpression()
			} else {
				ts, es := v.Expression.check(f.Block)
				if errsNotEmpty(es) {
					*errs = append(*errs, es...)
				}
				t, err := v.Expression.mustBeOneValueContext(ts)
				if err != nil {
					*errs = append(*errs, fmt.Errorf("%s err:%v", errMsgPrefix(v.Pos), err))
				}
				if t.TypeCompatible(v.Typ) == false {
					err = fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(v.Expression.Pos), t.TypeString(), v.Typ.TypeString())
					*errs = append(*errs, err)
				}
			}
		}
	}
}

type FunctionType struct {
	ParameterList ParameterList
	ReturnList    ReturnList
}

type ParameterList []*VariableDefinition
type ReturnList []*VariableDefinition

func (r ReturnList) retTypes(pos *Pos) []*VariableType {
	if r == nil || len(r) == 0 {
		t := &VariableType{}
		t.Typ = VARIABLE_TYPE_VOID // means no return;
		t.Pos = pos
		return []*VariableType{t}
	}
	ret := make([]*VariableType, len(r))
	for k, v := range r {
		ret[k] = v.Typ.Clone()
		ret[k].Pos = pos
	}
	return ret
}
