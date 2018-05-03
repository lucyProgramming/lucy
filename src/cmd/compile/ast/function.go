package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Function struct {
	ClassMethod                    *cg.MethodHighLevel // make call from
	ConstructionMethodCalledByUser bool
	HaveDefaultValue               bool
	DefaultValueStartAt            int
	IsClosureFunction              bool
	isGlobalVariableDefinition     bool
	isPackageBlockFunction         bool
	buildChecker                   buildFunctionChecker // used in build function
	IsGlobal                       bool
	IsBuildin                      bool
	Used                           bool
	AccessFlags                    uint16 // public private or protected
	Typ                            FunctionType
	ClosureVars                    Closure
	Name                           string // if name is nil string,means no name function
	Block                          *Block
	Pos                            *Pos
	Descriptor                     string
	AutoVarForException            *AutoVarForException
	AutoVarForReturnBecauseOfDefer *AutoVarForReturnBecauseOfDefer
	AutoVarForMultiReturn          *AutoVarForMultiReturn
	VarOffSet                      uint16 // for closure
}
type CallChecker func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error,
	args []*VariableType, pos *Pos)

type buildFunctionChecker CallChecker

type AutoVarForReturnBecauseOfDefer struct {
	ForArrayList uint16
}

func (f *Function) MkAutoVarForReturnBecauseOfDefer() {
	if f.AutoVarForReturnBecauseOfDefer != nil {
		return
	}
	f.AutoVarForReturnBecauseOfDefer = &AutoVarForReturnBecauseOfDefer{}
}

func (f *Function) NoReturnValue() bool {
	return len(f.Typ.ReturnList) == 0 ||
		f.Typ.ReturnList[0].Typ.Typ == VARIABLE_TYPE_VOID
}

type AutoVarForException struct {
	Offset uint16
}

func (f *Function) mkAutoVarForException() {
	if f.AutoVarForException != nil {
		return
	}
	f.AutoVarForException = &AutoVarForException{}
}

func (f *Function) mkAutoVarForMultiReturn() {
	if f.AutoVarForMultiReturn != nil {
		return
	}
	f.AutoVarForMultiReturn = &AutoVarForMultiReturn{}
}

type AutoVarForMultiReturn struct {
	Offset uint16
}

func (f *Function) readableMsg(name ...string) string {
	var s string
	if len(name) > 0 {
		s = "fn " + name[0] + "("
	} else {
		s = "fn " + f.Name + "("
	}
	for k, v := range f.Typ.ParameterList {
		s += " " + v.Name + " "
		s += v.Typ.TypeString()
		if v.Expression != nil {
			s += " = " + v.Expression.OpName()
		}
		s += " "
		if k != len(f.Typ.ParameterList)-1 {
			s += ","
		}
	}
	s += ")"
	if len(f.Typ.ReturnList) > 0 && f.NoReturnValue() == false {
		s += "->( "
		for k, v := range f.Typ.ReturnList {
			s += " " + v.Name + " "
			s += v.Typ.TypeString()
			if k != len(f.Typ.ReturnList)-1 {
				s += ","
			}
		}
		s += " )"
	}
	return s
}

/*
	no need return list
*/
func (f *Function) badParameterMsg(name string, args []*VariableType) string {
	s := "fn " + name + "("
	for k, v := range args {
		s += " " + v.TypeString() + " "
		if k != len(args)-1 {
			s += ","
		}
	}
	s += ")"
	return s
}

func (f *Function) checkBlock(errs *[]error) {
	f.mkLastRetrunStatement()
	*errs = append(*errs, f.Block.check()...)
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
		(f.Block.Statements[len(f.Block.Statements)-1].Typ != STATEMENT_TYPE_RETURN &&
			f.Block.Statements[len(f.Block.Statements)-1].Typ != STATEMENT_TYPE_SKIP) {
		s := &StatementReturn{}
		f.Block.Statements = append(f.Block.Statements, &Statement{
			Typ:             STATEMENT_TYPE_RETURN,
			StatementReturn: s,
			Pos:             f.Block.EndPos,
		})
	}
}

func (f *Function) checkParaMeterAndRetuns(errs *[]error) {
	if f.Name == MAIN_FUNCTION_NAME {
		errFunc := func() {
			*errs = append(*errs, fmt.Errorf("%s function '%s' expect declared as 'main(args []string)'",
				errMsgPrefix(f.Pos), MAIN_FUNCTION_NAME))
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
		f.Typ.ParameterList[0].LocalValOffset = 1
		//f.VarOffset = 2 //
		return
	}
	var err error
	for k, v := range f.Typ.ParameterList {
		//v.IsFunctionParameter = true
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
					errMsgPrefix(v.Expression.Pos)))
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
				continue
			}
			t, err := v.Expression.mustBeOneValueContext(ts)
			if err != nil {
				*errs = append(*errs, fmt.Errorf("%s err:%v", errMsgPrefix(v.Pos), err))
				continue
			}
			if t.TypeCompatible(v.Typ) == false {
				err = fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(v.Expression.Pos),
					t.TypeString(), v.Typ.TypeString())
				*errs = append(*errs, err)
			}
		}
	}

}
