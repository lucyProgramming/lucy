package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementSwitch struct {
	Pos                  *Pos
	Condition            *Expression //switch
	StatementSwitchCases []*StatementSwitchCase
	Default              *Block
	Exits                []*cg.Exit
}

type StatementSwitchCase struct {
	Matches []*Expression
	Block   *Block
}

func (s *StatementSwitch) check(block *Block) []error {
	errs := []error{}
	if s.Condition == nil { // must be a error at parse stage
		return errs
	}
	conditionType, es := s.Condition.checkSingleValueContextExpression(block)
	errs = append(errs, es...)
	if conditionType == nil {
		return errs
	}
	if conditionType.isTyped() == false {
		errs = append(errs, fmt.Errorf("%s condtion is not typed",
			errMsgPrefix(conditionType.Pos)))
		return errs
	}
	if conditionType.Type == VariableTypeBool {
		errs = append(errs, fmt.Errorf("%s bool expression not allow for switch",
			errMsgPrefix(conditionType.Pos)))
		return errs
	}
	if len(s.StatementSwitchCases) == 0 {
		errs = append(errs, fmt.Errorf("%s switch statement has no cases",
			errMsgPrefix(s.Pos)))
		return errs
	}
	byteMap := make(map[byte]*Pos)
	shortMap := make(map[int32]*Pos)
	int32Map := make(map[int32]*Pos)
	charMap := make(map[int32]*Pos)
	int64Map := make(map[int64]*Pos)
	floatMap := make(map[float32]*Pos)
	doubleMap := make(map[float64]*Pos)
	stringMap := make(map[string]*Pos)
	enumNamesMap := make(map[string]*Pos)
	enumPackageName := ""
	var byteValue byte
	var shortValue int32
	var intValue int32
	var charValue int32
	var longValue int64
	var floatValue float32
	var doubleValue float64
	var stringValue string
	var enumName string
	for _, v := range s.StatementSwitchCases {
		for _, e := range v.Matches {
			valueValid := false
			t, es := e.checkSingleValueContextExpression(block)
			errs = append(errs, es...)
			if t == nil {
				continue
			}
			if conditionType.assignAble(&errs, t) == false {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix(e.Pos), t.TypeString(), conditionType.TypeString()))
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
				if e.IsLiteral() {
					switch e.Type {
					case ExpressionTypeByte:
						byteValue = e.Data.(byte)
					case ExpressionTypeShort:
						shortValue = e.Data.(int32)
					case ExpressionTypeChar:
						charValue = e.Data.(int32)
					case ExpressionTypeInt:
						intValue = e.Data.(int32)
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
					errs = append(errs, fmt.Errorf("%s expression is not a literal value",
						errMsgPrefix(e.Pos)))
					continue
				}
			}
			if e.canBeUsedAsCondition() == false {
				errs = append(errs, fmt.Errorf("%s expression cannot use as condition",
					errMsgPrefix(e.Pos)))
				continue
			}
			if valueValid {
				errMsg := func(first *Pos, which string) error {
					errMsg := fmt.Sprintf("%s  '%s' duplicate case,first declared at:\n",
						errMsgPrefix(e.Pos), which)
					errMsg += fmt.Sprintf("\t%s", errMsgPrefix(first))
					return errors.New(errMsg)
				}
				switch conditionType.Type {
				case VariableTypeByte:
					if first, ok := byteMap[byteValue]; ok {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", byteValue)))
						continue // no check body
					} else {
						byteMap[byteValue] = e.Pos
					}
				case VariableTypeShort:
					if first, ok := shortMap[shortValue]; ok {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", shortValue)))
						continue // no check body
					} else {
						shortMap[shortValue] = e.Pos
					}
				case VariableTypeChar:
					if first, ok := charMap[charValue]; ok {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", charValue)))
						continue // no check body
					} else {
						charMap[charValue] = e.Pos
					}
				case VariableTypeInt:
					if first, ok := int32Map[intValue]; ok {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", intValue)))
						continue // no check body
					} else {
						int32Map[intValue] = e.Pos
					}
				case VariableTypeLong:
					if first, ok := int64Map[longValue]; ok {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", longValue)))
						continue // no check body
					} else {
						int64Map[longValue] = e.Pos
					}
				case VariableTypeFloat:
					if first, found := floatMap[floatValue]; found {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", floatValue)))
						continue // no check body
					} else {
						floatMap[floatValue] = e.Pos
					}
				case VariableTypeDouble:
					if first, found := doubleMap[doubleValue]; found {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", doubleValue)))
						continue // no check body
					} else {
						doubleMap[doubleValue] = e.Pos
					}
				case VariableTypeString:
					if first, ok := stringMap[stringValue]; ok {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", stringValue)))
						continue // no check body
					} else {
						stringMap[stringValue] = e.Pos
					}
				case VariableTypeEnum:
					if first, ok := enumNamesMap[enumName]; ok {
						errs = append(errs, errMsg(first, fmt.Sprintf("%v", enumName)))
						continue // no check body
					} else {
						enumNamesMap[enumName] = e.Pos
					}
				}
			}
		}
		if v.Block != nil {
			v.Block.inherit(block)
			v.Block.InheritedAttribute.ForBreak = s
			errs = append(errs, v.Block.checkStatements()...)
		}
	}
	if s.Default != nil {
		s.Default.inherit(block)
		s.Default.InheritedAttribute.ForBreak = s
		errs = append(errs, s.Default.checkStatements()...)
	}
	if conditionType.Type == VariableTypeEnum &&
		len(enumNamesMap) < len(conditionType.Enum.Enums) &&
		s.Default == nil {
		//some enum are missing, not allow
		errMsg := fmt.Sprintf("%s switch for enum '%s' is not complete\n",
			errMsgPrefix(s.Pos), conditionType.Enum.Name)
		errMsg += "\tyou can use 'default:' or give missing enums,which are:\n"
		for _, v := range conditionType.Enum.Enums {
			_, ok := enumNamesMap[v.Name]
			if ok {
				//handled
				continue
			}
			if enumPackageName == "" {
				errMsg += fmt.Sprintf("\t\tcase %s:\n", v.Name)
			} else {
				errMsg += fmt.Sprintf("\t\tcase %s.%s:\n", enumPackageName, v.Name)
			}
		}
		errs = append(errs, errors.New(errMsg))
	}
	return errs
}
