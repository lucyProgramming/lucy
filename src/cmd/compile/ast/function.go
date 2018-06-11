package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Function struct {
	BlockChecked                   bool // template function may be multi time check
	TemplateFunction               *TemplateFunction
	TypeParameters                 map[string]*VariableType //typed parameters
	ParameterAndRetrunListOK       bool
	ClassMethod                    *cg.MethodHighLevel // make call from
	ConstructionMethodCalledByUser bool
	HaveDefaultValue               bool
	DefaultValueStartAt            int
	IsClosureFunction              bool
	isGlobalVariableDefinition     bool
	isPackageBlockFunction         bool
	buildinFunctionChecker         buildFunctionChecker // used in build function
	IsGlobal                       bool
	IsBuildin                      bool
	Used                           bool
	AccessFlags                    uint16 // public private or protected
	Typ                            FunctionType
	Closure                        Closure
	Name                           string // if name is nil string,means no name function
	Block                          Block
	Pos                            *Pos
	Descriptor                     string
	AutoVarForException            *AutoVarForException
	AutoVarForReturnBecauseOfDefer *AutoVarForReturnBecauseOfDefer
	AutoVarForMultiReturn          *AutoVarForMultiReturn
	VarOffSet                      uint16 // for closure
	SourceCode                     []byte // source code for T
	//	OutterFunction                 *Function
}

type CallChecker func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*VariableType, pos *Pos)

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
	if f.BlockChecked {
		return
	}
	f.mkLastRetrunStatement()
	*errs = append(*errs, f.Block.checkStatements()...)
	f.BlockChecked = true
}

func (f *Function) check(b *Block) []error {
	errs := make([]error, 0)
	f.Block.inherit(b)
	f.checkParaMeterAndRetuns(&errs)
	f.Block.InheritedAttribute.Function = f
	if f.TemplateFunction == nil {
		f.checkBlock(&errs)
	}
	return errs
}

func (f *Function) clone(block *Block) (ret *Function, es []error) {
	ret, es = ParseFunctionHandler(f.SourceCode, f.Pos)
	if errsNotEmpty(es) {
		return ret, es
	}
	ret.Block.inherit(block)
	ret.checkParaMeterAndRetuns(&es)
	ret.Block.InheritedAttribute.Function = ret
	return ret, es
}
func (f *Function) mkLastRetrunStatement() {
	if len(f.Block.Statements) == 0 ||
		(f.Block.Statements[len(f.Block.Statements)-1].Typ != STATEMENT_TYPE_RETURN) {
		s := &StatementReturn{}
		f.Block.Statements = append(f.Block.Statements, &Statement{
			Typ:             STATEMENT_TYPE_RETURN,
			StatementReturn: s,
			Pos:             f.Block.EndPos,
		})
	}
}

func (f *Function) checkParaMeterAndRetuns(errs *[]error) {
	errsLength := len(*errs)
	defer func() {
		if len(*errs) == errsLength {
			f.ParameterAndRetrunListOK = true
		}
	}()
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
			f.Typ.ParameterList[0].LocalValOffset = 1
			f.Typ.ParameterList[0].IsFunctionParameter = true
		}

		return
	}
	var err error
	for k, v := range f.Typ.ParameterList {
		v.IsFunctionParameter = true
		if v.Typ.Typ != VARIABLE_TYPE_T {
			err = v.Typ.resolve(&f.Block)
			if err != nil {
				*errs = append(*errs, err)
			}
		}
		err = f.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, err)
			continue
		}
		if v.Typ.haveT() && f.TemplateFunction == nil {
			if f.HaveDefaultValue {
				*errs = append(*errs, fmt.Errorf("%s cannot have typed parameter after default value",
					errMsgPrefix(f.Pos)))
			}
			f.TemplateFunction = &TemplateFunction{}
		}
		if f.HaveDefaultValue && v.Expression == nil {
			*errs = append(*errs, fmt.Errorf("%s expect default value", errMsgPrefix(v.Pos)))
			continue
		}
		if v.Expression != nil {
			if v.Typ.Typ == VARIABLE_TYPE_T {
				*errs = append(*errs, fmt.Errorf("%s typ is tempalate,cannot have default value",
					errMsgPrefix(v.Pos)))
				continue
			}
			if f.HaveDefaultValue == false {
				f.DefaultValueStartAt = k
			}
			f.HaveDefaultValue = true
			t, es := v.Expression.checkSingleValueContextExpression(&f.Block)
			if errsNotEmpty(es) {
				*errs = append(*errs, es...)
			}
			if t != nil {
				if v.Typ.Equal(errs, t) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(v.Expression.Pos), t.TypeString(), v.Typ.TypeString()))
					continue
				}

			}
			if v.Expression.IsLiteral() == false {
				*errs = append(*errs, fmt.Errorf("%s default value must be literal",
					errMsgPrefix(v.Expression.Pos)))
				continue
			}
			if v.Expression.Typ == EXPRESSION_TYPE_NULL {
				*errs = append(*errs, fmt.Errorf("%s cannot use 'null' as default value",
					errMsgPrefix(v.Expression.Pos)))
			}
		}
	}

	//handler return
	for _, v := range f.Typ.ReturnList {
		v.IsFunctionRetrunVar = true
		if v.Typ.Typ != VARIABLE_TYPE_T {
			err = v.Typ.resolve(&f.Block)
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
		}
		if v.Typ.haveT() && f.TemplateFunction == nil {
			f.TemplateFunction = &TemplateFunction{}
		}
		err = f.Block.insert(v.Name, v.Pos, v)
		if err != nil {
			*errs = append(*errs, fmt.Errorf("%s err:%v", errMsgPrefix(v.Pos), err))
		}
		if v.Typ.Typ == VARIABLE_TYPE_T && v.Expression != nil {
			*errs = append(*errs, fmt.Errorf("%s typ is tempalate,cannot have default value",
				errMsgPrefix(v.Pos)))
			continue
		}

		if v.Expression == nil {
			v.Expression = v.Typ.mkDefaultValueExpression()
			continue
		}
		t, es := v.Expression.checkSingleValueContextExpression(&f.Block)
		if errsNotEmpty(es) {
			*errs = append(*errs, es...)
			continue
		}
		if t != nil && t.Equal(errs, v.Typ) == false {
			err = fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(v.Expression.Pos),
				t.TypeString(), v.Typ.TypeString())
			*errs = append(*errs, err)
		}
	}
}
