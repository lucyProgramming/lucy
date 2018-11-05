package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type buildInFunctionChecker func(
	f *Function,
	e *ExpressionFunctionCall,
	block *Block,
	errs *[]error,
	args []*Type,
	pos *Pos) // used in build function

type Function struct {
	CallFatherConstructionExpression *Expression
	TemplateFunction                 *TemplateFunction
	parameterTypes                   map[string]*Type    //typed parameters
	Entrance                         *cg.MethodHighLevel // make_node_objects call from
	buildInFunctionChecker           buildInFunctionChecker
	AccessFlags                      uint16
	Type                             FunctionType
	Closure                          Closure
	Name                             string // if name is nil string,means no name function
	Block                            Block
	Function                         *Function
	Pos                              *Pos
	JvmDescriptor                    string
	ClosureVariableOffSet            uint16 // for closure
	SourceCode                       []byte // source code for template function
	HasDefer                         bool
	HaveDefaultValue                 bool
	DefaultValueStartAt              int
	IsGlobal                         bool
	IsBuildIn                        bool
	LoadedFromCorePackage            bool
	Used                             bool
	TemplateClonedFunction           bool
	isPackageInitBlockFunction       bool
	Comment                          string
	IsClosureFunction                bool
}

func (this *Function) IsPublic() bool {
	return this.AccessFlags&cg.AccMethodPublic != 0
}

func (this *Function) NameLiteralFunction() string {
	if this.Name != "" {
		return this.Name
	}
	var t string
	if this.Name != "" {
		t = this.Block.InheritedAttribute.ClassAndFunctionNames + this.Name
	}
	return t
}

func (this *Function) readableMsg() string {
	if this.Name == "" {
		return "fn " + this.Type.TypeString()
	} else {
		return "fn " + this.Name + " " + this.Type.TypeString()
	}
}

func (this *Function) makeName() {
	if this.Name == "" {
		if this.Block.InheritedAttribute.ClassAndFunctionNames == "" {
			this.Name = fmt.Sprintf("literal$%d", this.Pos.Line)
		} else {
			this.Name = fmt.Sprintf("%s$literal%d",
				this.Block.InheritedAttribute.ClassAndFunctionNames, this.Pos.Line)
		}
	}
	name := this.Name
	if name == SpecialMethodInit {
		name = "init"
	}
	if this.Block.InheritedAttribute.ClassAndFunctionNames == "" {
		this.Block.InheritedAttribute.ClassAndFunctionNames = name
	} else {
		this.Block.InheritedAttribute.ClassAndFunctionNames += "$" + name
	}
}

func (this *Function) checkBlock(errs *[]error) {
	this.makeName()
	this.makeLastReturnStatement()
	*errs = append(*errs, this.Block.check()...)
}

func (this *Function) check(b *Block) []error {
	errs := make([]error, 0)
	this.Block.inherit(b)
	this.Block.InheritedAttribute.Function = this
	this.checkParametersAndReturns(&errs, true, false)
	if this.TemplateFunction == nil {
		this.checkBlock(&errs)
	}
	return errs
}

func (this *Function) clone() (ret *Function, es []error) {
	ret, es = ParseFunctionHandler(this.SourceCode, this.Pos)
	if len(es) > 0 {
		return ret, es
	}
	ret.TemplateClonedFunction = true
	return ret, es
}
func (this *Function) makeLastReturnStatement() {
	s := &StatementReturn{}
	this.Block.Statements = append(this.Block.Statements, &Statement{
		Type:            StatementTypeReturn,
		StatementReturn: s,
		Pos:             this.Block.EndPos,
	})
}

func (this *Function) IsGlobalMain() bool {
	return this.IsGlobal &&
		this.Name == MainFunctionName
}

func (this *Function) checkParametersAndReturns(
	errs *[]error,
	checkReturnVarExpression bool,
	isAbstract bool) {
	var err error
	for k, v := range this.Type.ParameterList {
		v.IsFunctionParameter = true
		if len(v.Type.getParameterType(&this.Type)) > 0 {
			if this.TemplateFunction == nil {
				this.TemplateFunction = &TemplateFunction{}
			}
		} else {
			err = v.Type.resolve(&this.Block)
			if err != nil {
				*errs = append(*errs, err)
			}
			if isAbstract == false {
				err = this.Block.Insert(v.Name, v.Pos, v)
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
			if k != len(this.Type.ParameterList)-1 {
				*errs = append(*errs, fmt.Errorf("%s only last parameter can be use as vargs",
					errMsgPrefix(v.Type.Pos)))
			} else {
				this.Type.ParameterList = this.Type.ParameterList[0:k]
				this.Type.VArgs = v
			}
			continue
		}
		if this.TemplateFunction != nil {
			continue
		}
		if v.DefaultValueExpression != nil {
			if this.HaveDefaultValue == false {
				this.DefaultValueStartAt = k
			}
			this.HaveDefaultValue = true
			t, es := v.DefaultValueExpression.checkSingleValueContextExpression(&this.Block)
			*errs = append(*errs, es...)
			if t != nil {
				if v.Type.assignAble(errs, t) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot use '%s' as '%s'",
						errMsgPrefix(v.DefaultValueExpression.Pos), t.TypeString(), v.Type.TypeString()))
					continue
				}
			}
			if v.DefaultValueExpression.isLiteral() == false {
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
	if this.Type.VoidReturn() == false {
		//handler return
		for _, v := range this.Type.ReturnList {
			v.IsReturn = true
			if len(v.Type.getParameterType(&this.Type)) > 0 {
				if this.TemplateFunction == nil {
					this.TemplateFunction = &TemplateFunction{}
				}
			} else {
				err = v.Type.resolve(&this.Block)
				if err != nil {
					*errs = append(*errs, err)
				}
				if isAbstract == false {
					err = this.Block.Insert(v.Name, v.Pos, v)
					if err != nil {
						*errs = append(*errs, err)
						continue
					}
				}
			}
			if this.TemplateFunction != nil {
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
			t, es := v.DefaultValueExpression.checkSingleValueContextExpression(&this.Block)
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

func (this *Function) checkReturnVarExpression() []error {
	if this.Type.VoidReturn() {
		return nil
	}
	var errs []error
	for _, v := range this.Type.ReturnList {
		t, es := v.DefaultValueExpression.checkSingleValueContextExpression(&this.Block)
		if len(es) > 0 {
			errs = append(errs, es...)
			continue
		}
		if t != nil && v.Type.assignAble(&errs, t) == false {
			err := fmt.Errorf("%s cannot assign '%s' to '%s'", errMsgPrefix(v.DefaultValueExpression.Pos),
				t.TypeString(), v.Type.TypeString())
			errs = append(errs, err)
		}
	}
	return errs
}
