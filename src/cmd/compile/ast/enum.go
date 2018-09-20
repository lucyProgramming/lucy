package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Enum struct {
	IsBuildIn       bool
	AccessFlags     uint16
	Name            string
	Pos             *Pos
	Enums           []*EnumName
	Init            *Expression //should be a int expression
	Used            bool
	DefaultValue    int32
	LoadFromOutSide bool
	Comment         string
}

func (e *Enum) IsPublic() bool {
	return e.AccessFlags&cg.ACC_CLASS_PUBLIC != 0
}

type EnumName struct {
	Enum       *Enum
	Name       string
	Pos        *Pos
	Value      int32 // int32 is bigger enough
	Comment    string
	Expression *Expression
}

func (e *Enum) check() (err error) {
	var initV int32 = 0
	if e.Init != nil {
		is, err := e.Init.constantFold()
		if err != nil {
			return err
		}
		if is == false ||
			e.Init.Type != ExpressionTypeInt {
			if err == nil {
				err = fmt.Errorf("%s enum type must inited by 'int' literal",
					errMsgPrefix(e.Pos))
				return err
			}
		}
		initV = e.Init.Data.(int32)
	}
	e.DefaultValue = initV
	for k, v := range e.Enums {
		if v.Expression != nil && err == nil {
			err = fmt.Errorf("%s enum only expect 1 init value",
				errMsgPrefix(v.Pos))
		}
		v.Value = int32(k) + initV
	}
	return err
}
