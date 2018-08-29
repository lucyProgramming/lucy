package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Enum struct {
	IsBuildIn    bool
	AccessFlags  uint16
	Name         string
	Pos          *Pos
	Enums        []*EnumName
	Init         *Expression //should be a int expression
	Used         bool
	DefaultValue int32
}

func (e *Enum) IsPublic() bool {
	return e.AccessFlags&cg.ACC_CLASS_PUBLIC != 0
}

type EnumName struct {
	Enum  *Enum
	Name  string
	Pos   *Pos
	Value int32 // int32 is bigger enough
}

func (e *Enum) check() (err error) {
	if e.Init == nil {
		e.Init = &Expression{}
		e.Init.Type = ExpressionTypeInt
		e.Init.Data = int32(0)
		e.Init.Pos = e.Pos
	}
	is, err := e.Init.constantFold()
	if err != nil || is == false || e.Init.Type != ExpressionTypeInt {
		if err == nil {
			err = fmt.Errorf("%s enum type must inited by integer_expression",
				errMsgPrefix(e.Pos))
		}
	}
	var initV int32 = 0
	if e.Init.Data != nil {
		if t, ok := e.Init.Data.(int32); ok {
			initV = t
		}
	}
	e.DefaultValue = initV
	for k, v := range e.Enums {
		v.Value = int32(k) + initV
	}
	return nil
}
