package ast

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementSwitch struct {
	Pos                  *Position
	Condition            *Expression //switch
	StatementSwitchCases []*StatementSwitchCase
	Default              *Block
	Exits                []*cg.Exit
}

type StatementSwitchCase struct {
	Matches []*Expression
	Block   *Block
}

func (s *StatementSwitch) check(b *Block) []error {
	errs := []error{}
	conditionType, es := s.Condition.checkSingleValueContextExpression(b)
	if errorsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if conditionType == nil {
		return errs
	}
	if conditionType.Type == VARIABLE_TYPE_BOOL {
		errs = append(errs, fmt.Errorf("%s bool expression not allow for switch",
			errMsgPrefix(conditionType.Pos)))
		return errs
	}
	if len(s.StatementSwitchCases) == 0 {
		errs = append(errs, fmt.Errorf("%s switch statement has no cases",
			errMsgPrefix(s.Pos)))
	}

	byteMap := make(map[byte]*Position)
	shortMap := make(map[int32]*Position)
	int32Map := make(map[int32]*Position)
	int64Map := make(map[int64]*Position)
	floatMap := make(map[float32]*Position)
	doubleMap := make(map[float64]*Position)
	stringMap := make(map[string]*Position)
	enumNamesMap := make(map[string]*Position)
	enumPackageName := ""
	for _, v := range s.StatementSwitchCases {
		for _, e := range v.Matches {
			var byteValue byte
			var shortValue int32
			var int32Value int32
			var int64Value int64
			var floatValue float32
			var doubleValue float64
			var stringValue string
			var enumName string
			valueValid := false
			valueFromExpression := func() {
				switch e.Type {
				case EXPRESSION_TYPE_BYTE:
					byteValue = e.Data.(byte)
				case EXPRESSION_TYPE_SHORT:
					shortValue = e.Data.(int32)
				case EXPRESSION_TYPE_INT:
					int32Value = e.Data.(int32)
				case EXPRESSION_TYPE_LONG:
					int64Value = e.Data.(int64)
				case EXPRESSION_TYPE_FLOAT:
					floatValue = e.Data.(float32)
				case EXPRESSION_TYPE_DOUBLE:
					doubleValue = e.Data.(float64)
				case EXPRESSION_TYPE_STRING:
					stringValue = e.Data.(string)
				}
			}
			t, es := e.checkSingleValueContextExpression(b)
			if errorsNotEmpty(es) {
				errs = append(errs, es...)
			}
			if t == nil {
				continue
			}
			if conditionType.Equal(&errs, t) == false {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'",
					errMsgPrefix(e.Pos), t.TypeString(), conditionType.TypeString()))
				continue
			}
			if conditionType.Type == VARIABLE_TYPE_ENUM {
				if t.EnumName == nil {
					errs = append(errs, fmt.Errorf("%s enum value is not literal",
						errMsgPrefix(e.Pos)))
					continue
				} else {
					if e.ExpressionValue.Type == VARIABLE_TYPE_PACKAGE &&
						enumPackageName == "" {
						enumPackageName = e.ExpressionValue.Package.Name
					}
					enumName = t.EnumName.Name
					valueValid = true
				}
			}
			if conditionType.IsPrimitive() {
				if e.IsLiteral() {
					valueFromExpression()
					valueValid = true
				} else {
					errs = append(errs, fmt.Errorf("%s expression is not a literal value", errMsgPrefix(e.Pos)))
					continue
				}
			}
			if e.canBeUsedAsCondition() == false {
				errs = append(errs, fmt.Errorf("%s expression cannot use as condition",
					errMsgPrefix(e.Pos)))
			}
			if valueValid {
				errMsg := func(first *Position) string {
					errMsg := fmt.Sprintf("%s duplicate case ,first declared at:\n", errMsgPrefix(e.Pos))
					errMsg += fmt.Sprintf("\t%s", errMsgPrefix(first))
					return errMsg
				}
				switch conditionType.Type {
				case VARIABLE_TYPE_BYTE:
					if first, ok := byteMap[byteValue]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						byteMap[byteValue] = e.Pos
					}
				case VARIABLE_TYPE_SHORT:
					if first, ok := shortMap[shortValue]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						shortMap[shortValue] = e.Pos
					}
				case VARIABLE_TYPE_INT:
					if first, ok := int32Map[int32Value]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						int32Map[int32Value] = e.Pos
					}
				case VARIABLE_TYPE_LONG:
					if first, ok := int64Map[int64Value]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						int64Map[int64Value] = e.Pos
					}
				case VARIABLE_TYPE_FLOAT:
					if first, found := floatMap[floatValue]; found {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						floatMap[floatValue] = e.Pos
					}
				case VARIABLE_TYPE_DOUBLE:
					if first, found := doubleMap[doubleValue]; found {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						doubleMap[doubleValue] = e.Pos
					}
				case VARIABLE_TYPE_STRING:
					if first, ok := stringMap[stringValue]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						stringMap[stringValue] = e.Pos
					}
				case VARIABLE_TYPE_ENUM:
					if first, ok := enumNamesMap[enumName]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
						continue // no check body
					} else {
						enumNamesMap[enumName] = e.Pos
					}
				}
			}
		}
		if v.Block != nil {
			v.Block.inherit(b)
			v.Block.InheritedAttribute.StatementSwitch = s
			v.Block.InheritedAttribute.ForBreak = s
			errs = append(errs, v.Block.checkStatements()...)
		}
	}
	if s.Default != nil {
		s.Default.inherit(b)
		s.Default.InheritedAttribute.StatementSwitch = s
		s.Default.InheritedAttribute.ForBreak = s
		errs = append(errs, s.Default.checkStatements()...)
	}
	if conditionType.Type == VARIABLE_TYPE_ENUM &&
		len(enumNamesMap) < len(conditionType.Enum.Enums) &&
		s.Default == nil {
		//some enum are missing, not allow
		errMsg := fmt.Sprintf("%s switch for enum '%s' is not complete\n",
			errMsgPrefix(s.Pos), conditionType.Enum.Name)
		errMsg += "\tyou can use 'default:' or give missing enums,which are:\n"
		for _, v := range conditionType.Enum.Enums {
			_, ok := enumNamesMap[v.Name]
			if ok {
				continue
			}
			if enumPackageName == "" {
				errMsg += fmt.Sprintf("\t\tcase %s:\n", v.Name)
			} else {
				errMsg += fmt.Sprintf("\t\tcase %s.%s:\n", enumPackageName, v.Name)
			}

		}
		errs = append(errs, fmt.Errorf(errMsg))
	}
	return errs
}
