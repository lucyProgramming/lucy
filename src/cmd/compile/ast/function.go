package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Function struct {
	TemplateFunction                    *TemplateFunction
	parameterTypes                      map[string]*Type    //typed parameters
	ClassMethod                         *cg.MethodHighLevel // make call from
	HaveDefaultValue                    bool
	DefaultValueStartAt                 int
	IsClosureFunction                   bool
	isGlobalVariableDefinition          bool
	isPackageBlockFunction              bool
	buildInFunctionChecker              buildFunctionChecker // used in build function
	IsGlobal                            bool
	IsBuildIn                           bool
	Used                                bool
	AccessFlags                         uint16
	Type                                FunctionType
	Closure                             Closure
	Name                                string // if name is nil string,means no name function
	Block                               Block
	Pos                                 *Position
	Descriptor                          string
	AutoVariableForException            *AutoVariableForException
	AutoVariableForReturnBecauseOfDefer *AutoVariableForReturnBecauseOfDefer
	AutoVariableForMultiReturn          *AutoVariableForMultiReturn
	ClosureVariableOffSet               uint16 // for closure
	SourceCodes                         []byte // source code for template function
}

func (f *Function) NameLiteralFunction() string {
	if f.Name != "" {
		return f.Name
	}
	var t string
	if f.Name != "" {
		t = f.Block.InheritedAttribute.ClassAndFunctionNames + f.Name
	}
	return t
}

type CallChecker func(f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Position)

type buildFunctionChecker CallChecker

type AutoVariableForReturnBecauseOfDefer struct {
	Offset uint16
}

func (f *Function) MkAutoVarForReturnBecauseOfDefer() {
	if f.AutoVariableForReturnBecauseOfDefer != nil {
		return
	}
	f.AutoVariableForReturnBecauseOfDefer = &AutoVariableForReturnBecauseOfDefer{}
}

func (f *Function) NoReturnValue() bool {
	return f.Type.NoReturnValue()
}

type AutoVariableForException struct {
	Offset uint16
}

func (f *Function) mkAutoVarForException() {
	if f.AutoVariableForException != nil {
		return
	}
	f.AutoVariableForException = &AutoVariableForException{}
}

func (f *Function) mkAutoVarForMultiReturn() {
	if f.AutoVariableForMultiReturn != nil {
		return
	}
	f.AutoVariableForMultiReturn = &AutoVariableForMultiReturn{}
}

type AutoVariableForMultiReturn struct {
	Offset uint16
}

func (f *Function) readableMsg(name ...string) string {
	var s string
	if len(name) > 0 {
		s = "fn " + name[0] + "("
	} else {
		s = "fn " + f.Name + "("
	}
	for k, v := range f.Type.ParameterList {
		s += " " + v.Name + " "
		s += v.Type.TypeString()
		if v.Expression != nil {
			s += " = " + v.Expression.OpName()
		}
		s += " "
		if k != len(f.Type.ParameterList)-1 {
			s += ","
		}
	}
	s += ")"
	if len(f.Type.ReturnList) > 0 && f.NoReturnValue() == false {
		s += "->( "
		for k, v := range f.Type.ReturnList {
			s += " " + v.Name + " "
			s += v.Type.TypeString()
			if k != len(f.Type.ReturnList)-1 {
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
func (f *Function) badParameterMsg(name string, args []*Type) string {
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

func (f *Function) makeName() {
	if f.Name == "" {
		if f.Block.InheritedAttribute.ClassAndFunctionNames == "" {
			f.Name = fmt.Sprintf("literal$%d", f.Pos.StartLine)
		} else {
			f.Name = fmt.Sprintf("%s$literal%d",
				f.Block.InheritedAttribute.ClassAndFunctionNames, f.Pos.StartLine)
		}
	}
	if f.Block.InheritedAttribute.ClassAndFunctionNames == "" {
		f.Block.InheritedAttribute.ClassAndFunctionNames = f.Name
	} else {
		f.Block.InheritedAttribute.ClassAndFunctionNames += "$" + f.Name
	}
}

func (f *Function) checkBlock(errs *[]error) {
	f.makeName()
	f.makeLastReturnStatement()
	*errs = append(*errs, f.Block.checkStatements()...)
}

func (f *Function) check(b *Block) []error {
	errs := make([]error, 0)
	f.Block.inherit(b)
	f.checkParametersAndReturns(&errs)
	f.Block.InheritedAttribute.Function = f
	if f.TemplateFunction == nil {
		f.checkBlock(&errs)
	}
	return errs
}

func (f *Function) clone() (ret *Function, es []error) {
	ret, es = ParseFunctionHandler(f.SourceCodes, f.Pos)
	if esNotEmpty(es) {
		return ret, es
	}
	return ret, es
}
func (f *Function) makeLastReturnStatement() {
	if len(f.Block.Statements) == 0 ||
		(f.Block.Statements[len(f.Block.Statements)-1].Type != StatementTypeReturn) {
		s := &StatementReturn{}
		f.Block.Statements = append(f.Block.Statements, &Statement{
			Type:            StatementTypeReturn,
			StatementReturn: s,
			Pos:             f.Block.EndPos,
		})
	}
}

func (f *Function) checkParametersAndReturns(errs *[]error) {
	if f.Name == MainFunctionName {
		errMain := func() {
			*errs = append(*errs, fmt.Errorf("%s function '%s' expect declared as 'main(args []string)'",
				errMsgPrefix(f.Pos), MainFunctionName))
		}
		if len(f.Type.ParameterList) != 1 {
			errMain()
		} else { //
			if f.Type.ParameterList[0].Type.Type == VariableTypeArray &&
				f.Type.ParameterList[0].Type.Array.Type == VariableTypeString {
				err := f.Block.Insert(f.Type.ParameterList[0].Name, f.Type.ParameterList[0].Pos, f.Type.ParameterList[0])
				if err != nil {
					*errs = append(*errs, err)
				}
			} else {
				errMain()
			}
			//f.Type.ParameterList[0].LocalValOffset = 1
			f.Type.ParameterList[0].IsFunctionParameter = true
		}
		if len(f.Type.ReturnList) > 0 {
			errMain()
		}
		return
	}
	var err error
	for k, v := range f.Type.ParameterList {
		v.IsFunctionParameter = true
		if len(v.Type.getParameterType()) > 0 {
			if f.TemplateFunction == nil {
				f.TemplateFunction = &TemplateFunction{}
			}
		} else {
			err = v.Type.resolve(&f.Block)
			if err != nil {
				*errs = append(*errs, err)
			}
			err = f.Block.Insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
		}
		if f.TemplateFunction != nil {
			continue
		}
		if v.Expression != nil {
			if f.HaveDefaultValue == false {
				f.DefaultValueStartAt = k
			}
			f.HaveDefaultValue = true
			t, es := v.Expression.checkSingleValueContextExpression(&f.Block)
			if esNotEmpty(es) {
				*errs = append(*errs, es...)
			}
			if t != nil {
				if v.Type.Equal(errs, t) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(v.Expression.Pos), t.TypeString(), v.Type.TypeString()))
					continue
				}

			}
			if v.Expression.IsLiteral() == false {
				*errs = append(*errs, fmt.Errorf("%s default value must be literal",
					errMsgPrefix(v.Expression.Pos)))
				continue
			}
			if v.Expression.Type == ExpressionTypeNull {
				*errs = append(*errs, fmt.Errorf("%s cannot use 'null' as default value",
					errMsgPrefix(v.Expression.Pos)))
			}
		}
	}

	//handler return
	for _, v := range f.Type.ReturnList {
		v.IsFunctionReturnVariable = true
		if len(v.Type.getParameterType()) > 0 {
			if f.TemplateFunction == nil {
				f.TemplateFunction = &TemplateFunction{}
			}
		} else {
			err = v.Type.resolve(&f.Block)
			if err != nil {
				*errs = append(*errs, err)
			}
			err = f.Block.Insert(v.Name, v.Pos, v)
			if err != nil {
				*errs = append(*errs, err)
				continue
			}
		}
		if f.TemplateFunction != nil {
			continue
		}
		if v.Expression == nil {
			v.Expression = v.Type.mkDefaultValueExpression()
			continue
		}
		t, es := v.Expression.checkSingleValueContextExpression(&f.Block)
		if esNotEmpty(es) {
			*errs = append(*errs, es...)
			continue
		}
		if t != nil && t.Equal(errs, v.Type) == false {
			err = fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(v.Expression.Pos),
				t.TypeString(), v.Type.TypeString())
			*errs = append(*errs, err)
		}
	}
}
