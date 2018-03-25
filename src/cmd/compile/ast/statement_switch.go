package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementSwitch struct {
	Pos                 *Pos
	BackPatchs          []*cg.JumpBackPatch
	Condition           *Expression //switch
	StatmentSwitchCases []*StatmentSwitchCase
	Default             *Block
}

type StatmentSwitchCase struct {
	Matches []*Expression
	Block   *Block
}

func (s *StatementSwitch) check(b *Block) []error {
	errs := []error{}
	conditionType, es := b.checkExpression(s.Condition)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if conditionType == nil {
		return errs
	}
	if conditionType.Typ == VARIABLE_TYPE_BOOL {
		errs = append(errs, fmt.Errorf("%s bool not allow for switch", errMsgPrefix(conditionType.Pos)))
		return errs
	}
	if len(s.StatmentSwitchCases) == 0 {
		errs = append(errs, fmt.Errorf("%s switch statement has no cases", errMsgPrefix(s.Pos)))
	}
	byteMap := make(map[byte]*Pos)
	shortMap := make(map[int16]*Pos)
	int32Map := make(map[int32]*Pos)
	int64Map := make(map[int64]*Pos)
	type floatExist struct {
		value float32
		pos   *Pos
	}
	floatMap := []floatExist{}
	type doubleExist struct {
		value float64
		pos   *Pos
	}
	doubleMap := []doubleExist{}
	stringMap := make(map[string]*Pos)
	for _, v := range s.StatmentSwitchCases {
		for _, e := range v.Matches {
			t, es := b.checkExpression(e)
			if errsNotEmpty(es) {
				errs = append(errs, es...)
			}
			if t == nil {
				continue
			}
			var byteValue byte
			var shortValue int16
			var int32Vavlue int32
			var int64Value int64
			var floatValue float32
			var doubleValue float64
			var stringValue string
			if conditionType.IsPrimitive() {
				if e.IsLiteral() {
					switch e.Typ {
					case EXPRESSION_TYPE_BYTE:
						byteValue = e.Data.(byte)
					case EXPRESSION_TYPE_INT:
						int32Vavlue = e.Data.(int32)
					case EXPRESSION_TYPE_LONG:
						int64Value = e.Data.(int64)
					case EXPRESSION_TYPE_FLOAT:
						floatValue = e.Data.(float32)
					case EXPRESSION_TYPE_DOUBLE:
						doubleValue = e.Data.(float64)
					case EXPRESSION_TYPE_STRING:
						stringValue = e.Data.(string)
					}
				} else {
					errs = append(errs, fmt.Errorf("%s expression is not a literal value", errMsgPrefix(e.Pos)))
					continue
				}
			}
			if conditionType.Equal(t) == false {
				errs = append(errs, fmt.Errorf("%s cannot use '%s' as '%s'", errMsgPrefix(e.Pos), t.TypeString(), conditionType.TypeString()))
				continue
			}
			errMsg := func(first *Pos) string {
				errmsg := fmt.Sprintf("%s duplicate case ,first declared at:\n", errMsgPrefix(e.Pos))
				errmsg += fmt.Sprintf("%s", errMsgPrefix(first))
				return errmsg
			}

			if conditionType.IsPrimitive() {
				switch conditionType.Typ {
				case VARIABLE_TYPE_BYTE:
					if first, ok := byteMap[byteValue]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
					} else {
						byteMap[byteValue] = e.Pos
					}
				case VARIABLE_TYPE_SHORT:
					if first, ok := shortMap[shortValue]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
					} else {
						shortMap[shortValue] = e.Pos
					}
				case VARIABLE_TYPE_INT:
					if first, ok := int32Map[int32Vavlue]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
					} else {
						int32Map[int32Vavlue] = e.Pos
					}
				case VARIABLE_TYPE_LONG:
					if first, ok := int64Map[int64Value]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
					} else {
						int64Map[int64Value] = e.Pos
					}
				case VARIABLE_TYPE_FLOAT:
					var first *Pos
					found := false
					for _, v := range floatMap {
						if v.value == floatValue {
							first = v.pos
							found = true
							break
						}
					}
					if found {
						errs = append(errs, fmt.Errorf(errMsg(first)))
					} else {
						floatMap = append(floatMap, floatExist{value: floatValue, pos: e.Pos})
					}
				case VARIABLE_TYPE_DOUBLE:
					var first *Pos
					found := false
					for _, v := range doubleMap {
						if v.value == doubleValue {
							found = true
							first = v.pos
							break
						}
					}
					if found {
						errs = append(errs, fmt.Errorf(errMsg(first)))
					} else {
						doubleMap = append(doubleMap, doubleExist{value: doubleValue, pos: e.Pos})
					}
				case VARIABLE_TYPE_STRING:
					if first, ok := stringMap[stringValue]; ok {
						errs = append(errs, fmt.Errorf(errMsg(first)))
					} else {
						stringMap[stringValue] = e.Pos
					}
				}
			}
		}
		if v.Block != nil {
			v.Block.inherite(b)
			errs = append(errs, v.Block.check()...)
		}
	}
	if s.Default != nil {
		s.Default.inherite(b)
		errs = append(errs, s.Default.check()...)
	}
	return errs
}
