package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type Enum struct {
	IsGlobal        bool
	IsBuildIn       bool
	AccessFlags     uint16
	Name            string
	Pos             *Pos
	Enums           []*EnumName
	Init            *Expression //should be a int expression
	Used            bool
	DefaultValue    int32
	LoadFromOutSide bool
	FirstValueIndex int
	Comment         string
}

func (e *Enum) isPublic() bool {
	return e.AccessFlags&cg.AccClassPublic != 0
}

type EnumName struct {
	Enum    *Enum
	Name    string
	Pos     *Pos
	Value   int32 // int32 is bigger enough
	Comment string
	NoNeed  *Expression
}

func (e *Enum) check() (errs []error) {
	var initV int32 = 0
	errs = []error{}
	if e.Init != nil {
		if is, err := e.Init.constantFold(); err != nil {
			errs = append(errs, err)
		} else {
			if is == false {
				err := fmt.Errorf("%s enum type must inited by 'int' literal",
					e.Pos.ErrMsgPrefix())
				errs = append(errs, err)
			} else {
				initV = int32(e.Init.getLongValue())
			}
		}
	}
	e.DefaultValue = initV
	for k, v := range e.Enums {
		if v.NoNeed != nil {
			errs = append(errs, fmt.Errorf("%s enum only expect 1 init value",
				v.Pos.ErrMsgPrefix()))
		}
		if k < e.FirstValueIndex {
			v.Value = initV - int32(e.FirstValueIndex-k)
		} else {
			v.Value = initV + int32(k-e.FirstValueIndex)
		}
	}
	return errs
}
