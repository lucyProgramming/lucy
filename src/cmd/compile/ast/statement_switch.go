package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementSwitch struct {
	PrefixExpressions    []*Expression
	initExpressionBlock  Block
	Pos                  *Pos
	EndPos               *Pos
	Condition            *Expression //switch
	StatementSwitchCases []*StatementSwitchCase
	Default              *Block
	Exits                []*cg.Exit
}

type StatementSwitchCase struct {
	Matches []*Expression
	Block   *Block
}

func (this *StatementSwitch) check(block *Block) []error {
	errs := []error{}
	if this.Condition == nil { // must be a error at parse stage
		return errs
	}
	this.initExpressionBlock.inherit(block)
	for _, v := range this.PrefixExpressions {
		v.IsStatementExpression = true
		_, es := v.check(&this.initExpressionBlock)
		errs = append(errs, es...)
		if err := v.canBeUsedAsStatement(); err != nil {
			errs = append(errs, err)
		}
	}
	if this.Condition == nil {
		return errs
	}
	conditionType, es := this.Condition.checkSingleValueContextExpression(&this.initExpressionBlock)
	errs = append(errs, es...)
	if conditionType == nil {
		return errs
	}
	if err := conditionType.isTyped(); err != nil {
		errs = append(errs, err)
		return errs
	}
	if conditionType.Type == VariableTypeBool {
		errs = append(errs, fmt.Errorf("%s bool expression not allow for switch",
			conditionType.Pos.ErrMsgPrefix()))
		return errs
	}
	if len(this.StatementSwitchCases) == 0 {
		errs = append(errs, fmt.Errorf("%s switch statement has no cases",
			this.EndPos.ErrMsgPrefix()))
		return errs
	}
	longMap := make(map[int64]*Pos)
	floatMap := make(map[float32]*Pos)
	doubleMap := make(map[float64]*Pos)
	stringMap := make(map[string]*Pos)
	enumNamesMap := make(map[string]*Pos)
	enumPackageName := ""
	var longValue int64
	var floatValue float32
	var doubleValue float64
	var stringValue string
	var enumName string

	for _, v := range this.StatementSwitchCases {
		for _, e := range v.Matches {
			valueValid := false
			t, es := e.checkSingleValueContextExpression(&this.initExpressionBlock)
			errs = append(errs, es...)
			if t == nil {
				continue
			}
			if t.Type == VariableTypeBool { // bool condition
				continue
			}
			if conditionType.assignAble(&errs, t) == false {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					e.Pos.ErrMsgPrefix(), t.TypeString(), conditionType.TypeString()))
				continue
			}
			if conditionType.Type == VariableTypeEnum {
				if t.EnumName == nil {
					errs = append(errs, fmt.Errorf("%s enum value is not literal",
						errMsgPrefix(e.Pos)))
					continue
				} else {
					if e.Value.Type == VariableTypePackage &&
						enumPackageName == "" {
						enumPackageName = e.Value.Package.Name
					}
					enumName = t.EnumName.Name
					valueValid = true
				}
			}
			if conditionType.IsPrimitive() {
				if e.isLiteral() {
					switch e.Type {
					case ExpressionTypeByte:
						fallthrough
					case ExpressionTypeShort:
						fallthrough
					case ExpressionTypeChar:
						fallthrough
					case ExpressionTypeInt:
						fallthrough
					case ExpressionTypeLong:
						longValue = e.Data.(int64)
					case ExpressionTypeFloat:
						floatValue = e.Data.(float32)
					case ExpressionTypeDouble:
						doubleValue = e.Data.(float64)
					case ExpressionTypeString:
						stringValue = e.Data.(string)
					}
					valueValid = true
				} else {
					errs = append(errs, fmt.Errorf("%s expression is not literal",
						errMsgPrefix(e.Pos)))
					continue
				}
			}
			if err := e.canBeUsedAsCondition(); err != nil {
				errs = append(errs, err)
				continue
			}
			if valueValid {
				errMsg := func(first *Pos, which interface{}) error {
					errMsg := fmt.Sprintf("%s  '%v' duplicate case,first declared at:\n",
						e.Pos.ErrMsgPrefix(), which)
					errMsg += fmt.Sprintf("\t%s", first.ErrMsgPrefix())
					return errors.New(errMsg)
				}
				switch conditionType.Type {
				case VariableTypeByte:
					fallthrough
				case VariableTypeShort:
					fallthrough
				case VariableTypeChar:
					fallthrough
				case VariableTypeInt:
					fallthrough
				case VariableTypeLong:
					if first, ok := longMap[longValue]; ok {
						errs = append(errs, errMsg(first, longValue))
						continue // no check body
					} else {
						longMap[longValue] = e.Pos
					}
				case VariableTypeFloat:
					if first, found := floatMap[floatValue]; found {
						errs = append(errs, errMsg(first, floatValue))
						continue // no check body
					} else {
						floatMap[floatValue] = e.Pos
					}
				case VariableTypeDouble:
					if first, found := doubleMap[doubleValue]; found {
						errs = append(errs, errMsg(first, doubleValue))
						continue // no check body
					} else {
						doubleMap[doubleValue] = e.Pos
					}
				case VariableTypeString:
					if first, ok := stringMap[stringValue]; ok {
						errs = append(errs, errMsg(first, stringValue))
						continue // no check body
					} else {
						stringMap[stringValue] = e.Pos
					}
				case VariableTypeEnum:
					if first, ok := enumNamesMap[enumName]; ok {
						errs = append(errs, errMsg(first, enumName))
						continue // no check body
					} else {
						enumNamesMap[enumName] = e.Pos
					}
				}
			}
		}
		if v.Block != nil {
			v.Block.inherit(&this.initExpressionBlock)
			v.Block.InheritedAttribute.ForBreak = this
			errs = append(errs, v.Block.check()...)
		}
	}
	if this.Default != nil {
		this.Default.inherit(&this.initExpressionBlock)
		this.Default.InheritedAttribute.ForBreak = this
		errs = append(errs, this.Default.check()...)
	}
	//if conditionType.Type == VariableTypeEnum &&
	//	len(enumNamesMap) < len(conditionType.Enum.Enums) &&
	//	this.Default == nil &&
	//	containsBool == false {
	//	//some enum are missing, not allow
	//	errMsg := fmt.Sprintf("%s switch for enum '%s' is not complete\n",
	//		this.EndPos.ErrMsgPrefix(), conditionType.Enum.Name)
	//	errMsg += "\tyou can use 'default:' or give missing enums,which are:\n"
	//	for _, v := range conditionType.Enum.Enums {
	//		_, ok := enumNamesMap[v.Name]
	//		if ok {
	//			//handled
	//			continue
	//		}
	//		if enumPackageName == "" {
	//			errMsg += fmt.Sprintf("\t\tcase %s:\n", v.Name)
	//		} else {
	//			errMsg += fmt.Sprintf("\t\tcase %s.%s:\n", enumPackageName, v.Name)
	//		}
	//	}
	//	errs = append(errs, errors.New(errMsg))
	//}
	return errs
}
