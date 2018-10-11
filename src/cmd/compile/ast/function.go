package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Function struct {
	CallFatherConstructionExpression *Expression
	TemplateFunction                 *TemplateFunction
	parameterTypes                   map[string]*Type    //typed parameters
	Entrance                         *cg.MethodHighLevel // make call from
	buildInFunctionChecker           func(
		f *Function, e *ExpressionFunctionCall, block *Block, errs *[]error, args []*Type, pos *Pos) // used in build function
	AccessFlags                uint16
	Type                       FunctionType
	Closure                    Closure
	Name                       string // if name is nil string,means no name function
	Block                      Block
	Pos                        *Pos
	JvmDescriptor              string
	ClosureVariableOffSet      uint16 // for closure
	SourceCode                 []byte // source code for template function
	HasDefer                   bool
	HaveDefaultValue           bool
	DefaultValueStartAt        int
	IsGlobal                   bool
	IsBuildIn                  bool
	LoadedFromCorePackage      bool
	Used                       bool
	TemplateClonedFunction     bool
	IsClosureFunction          bool
	isGlobalVariableDefinition bool
	isPackageInitBlockFunction bool
	Comment                    string
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

func (f *Function) readableMsg() string {
	if f.Name == "" {
		return "fn " + f.Type.typeString()
	} else {
		return "fn " + f.Name + " " + f.Type.typeString()
	}
}

func (f *Function) makeName() {
	if f.Name == "" {
		if f.Block.InheritedAttribute.ClassAndFunctionNames == "" {
			f.Name = fmt.Sprintf("literal$%d", f.Pos.Line)
		} else {
			f.Name = fmt.Sprintf("%s$literal%d",
				f.Block.InheritedAttribute.ClassAndFunctionNames, f.Pos.Line)
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
	f.checkParametersAndReturns(&errs, true, false)
	if f.TemplateFunction == nil {
		f.checkBlock(&errs)
	}
	return errs
}

func (f *Function) clone() (ret *Function, es []error) {
	ret, es = ParseFunctionHandler(f.SourceCode, f.Pos)
	if len(es) > 0 {
		return ret, es
	}
	ret.TemplateClonedFunction = true
	return ret, es
}
func (f *Function) makeLastReturnStatement() {
	s := &StatementReturn{}
	f.Block.Statements = append(f.Block.Statements, &Statement{
		Type:            StatementTypeReturn,
		StatementReturn: s,
		Pos:             f.Block.EndPos,
	})
}
func (f *Function) IsGlobalMain() bool {
	return f.IsGlobal &&
		f.Name == MainFunctionName
}

func (f *Function) checkParametersAndReturns(errs *[]error, checkReturnVarExpression bool, isAbstract bool) {
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
			if isAbstract == false {
				err = f.Block.Insert(v.Name, v.Pos, v)
				if err != nil {
					*errs = append(*errs, err)
					continue
				}
			}
		}
		if v.Type.IsVariableArgs && v.DefaultValueExpression != nil {
			*errs = append(*errs, fmt.Errorf("%s vargs cannot have default value",
				errMsgPrefix(v.Type.Pos)))
		}
		if v.Type.IsVariableArgs {
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
		if v.DefaultValueExpression != nil {
			if f.HaveDefaultValue == false {
				f.DefaultValueStartAt = k
			}
			f.HaveDefaultValue = true
			t, es := v.DefaultValueExpression.checkSingleValueContextExpression(&f.Block)
			*errs = append(*errs, es...)
			if t != nil {
				if v.Type.assignAble(errs, t) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(v.DefaultValueExpression.Pos), t.TypeString(), v.Type.TypeString()))
					continue
				}
			}
			if v.DefaultValueExpression.IsLiteral() == false {
				*errs = append(*errs, fmt.Errorf("%s default value must be literal",
					errMsgPrefix(v.DefaultValueExpression.Pos)))
				continue
			}
			if v.DefaultValueExpression.Type == ExpressionTypeNull {
				*errs = append(*errs, fmt.Errorf("%s cannot use 'null' as default value",
					errMsgPrefix(v.DefaultValueExpression.Pos)))
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
				if isAbstract == false {
					err = f.Block.Insert(v.Name, v.Pos, v)
					if err != nil {
						*errs = append(*errs, err)
						continue
					}
				}
			}
			if f.TemplateFunction != nil {
				continue
			}
			if v.DefaultValueExpression == nil {
				v.DefaultValueExpression = v.Type.mkDefaultValueExpression()
				continue
			}
			if checkReturnVarExpression == false {
				// eval expression later
				continue
			}
			t, es := v.DefaultValueExpression.checkSingleValueContextExpression(&f.Block)
			if len(es) > 0 {
				*errs = append(*errs, es...)
				continue
			}
			if t != nil && v.Type.assignAble(errs, t) == false {
				err = fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(v.DefaultValueExpression.Pos),
					t.TypeString(), v.Type.TypeString())
				*errs = append(*errs, err)
			}
		}
	}
}
