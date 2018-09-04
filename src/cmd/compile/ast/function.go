package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Function struct {
	CallFatherConstructionExpression *Expression
	TemplateFunction                 *TemplateFunction
	parameterTypes                   map[string]*Type     //typed parameters
	ClassMethod                      *cg.MethodHighLevel  // make call from
	buildInFunctionChecker           buildFunctionChecker // used in build function
	AccessFlags                      uint16
	Type                             FunctionType
	Closure                          Closure
	Name                             string // if name is nil string,means no name function
	Block                            Block
	Pos                              *Pos
	JvmDescriptor                    string
	ExpressionCount                  int
	ClosureVariableOffSet            uint16 // for closure
	SourceCodes                      []byte // source code for template function
	HasDefer                         bool
	HaveDefaultValue                 bool
	DefaultValueStartAt              int
	IsGlobal                         bool
	IsBuildIn                        bool
	LoadedFromLucyLangPackage        bool
	Used                             bool
	TemplateClonedFunction           bool
	IsClosureFunction                bool
	isGlobalVariableDefinition       bool
	isPackageInitBlockFunction       bool
	AccessByName                     int
}

func (f *Function) IsPublic() bool {
	return f.AccessFlags&cg.ACC_METHOD_PUBLIC != 0
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

type buildFunctionChecker func(f *Function, e *ExpressionFunctionCall, block *Block,
	errs *[]error, args []*Type, pos *Pos)

type AutoVariableForReturnBecauseOfDefer struct {
	Offset uint16
}

func (f *Function) VoidReturn() bool {
	return f.Type.VoidReturn()
}

func (f *Function) readableMsg(name ...string) string {
	var s string
	if len(name) > 0 {
		s = "fn " + name[0] + "("
	} else {
		s = "fn " + f.Name + "("
	}
	for k, v := range f.Type.ParameterList {
		s += v.Name + " "
		s += v.Type.TypeString()
		if v.Expression != nil {
			s += " = " + v.Expression.Description
		}
		if k != len(f.Type.ParameterList)-1 {
			s += ","
		}
	}
	if f.Type.VArgs != nil {
		if len(f.Type.ParameterList) > 0 {
			s += ","
		}
		s += f.Type.VArgs.Name + " "
		s += f.Type.VArgs.Type.TypeString()
	}
	s += ")"
	if len(f.Type.ReturnList) > 0 && f.VoidReturn() == false {
		s += "->( "
		for k, v := range f.Type.ReturnList {
			s += v.Name + " "
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
	name := f.Name
	if name == SpecialMethodInit {
		name = "init"
	}
	if f.Block.InheritedAttribute.ClassAndFunctionNames == "" {
		f.Block.InheritedAttribute.ClassAndFunctionNames = name
	} else {
		f.Block.InheritedAttribute.ClassAndFunctionNames += "$" + name
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
	f.Block.InheritedAttribute.Function = f
	f.checkParametersAndReturns(&errs)
	if f.TemplateFunction == nil {
		f.checkBlock(&errs)
	}
	return errs
}

func (f *Function) clone() (ret *Function, es []error) {
	ret, es = ParseFunctionHandler(f.SourceCodes, f.Pos)
	if len(es) > 0 {
		return ret, es
	}
	ret.TemplateClonedFunction = true
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
func (f *Function) IsGlobalMain() bool {
	return f.IsGlobal && f.Name == MainFunctionName
}
func (f *Function) checkParametersAndReturns(errs *[]error) {
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
		if v.Type.IsVArgs && v.Expression != nil {
			*errs = append(*errs, fmt.Errorf("%s vargs cannot have default value",
				errMsgPrefix(v.Type.Pos)))
		}
		if v.Type.IsVArgs {
			if k != len(f.Type.ParameterList)-1 {
				*errs = append(*errs, fmt.Errorf("%s only last parameter can be use as vargs",
					errMsgPrefix(v.Type.Pos)))
			} else {
				f.Type.ParameterList = f.Type.ParameterList[0:k]
				f.Type.VArgs = v
			}
			continue
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

			*errs = append(*errs, es...)

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
	if f.Type.VoidReturn() == false {
		//handler return
		for _, v := range f.Type.ReturnList {
			v.IsReturn = true
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
			if len(es) > 0 {
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
}
