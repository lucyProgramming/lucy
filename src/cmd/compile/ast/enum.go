package ast

import "fmt"

type EnumName struct {
	Enum  *Enum
	Name  string
	Pos   *Pos
	Value int32 // int32 is bigger enough
}

type Enum struct {
	AccessFlags uint16
	Name        string
	Pos         *Pos
	Enums       []*EnumName
	Init        *Expression //should be a int expression
	Used        bool
}

func (e *Enum) check() (err error) {
	if e.Init == nil {
		e.Init = &Expression{}
		e.Init.Typ = EXPRESSION_TYPE_INT
		e.Init.Data = int32(0)
		e.Init.Pos = e.Pos
	}
	is, err := e.Init.constFold()
	if err != nil || is == false || e.Init.Typ != EXPRESSION_TYPE_INT {
		if err == nil {
			err = fmt.Errorf("%s enum type must inited by integer_expression", errMsgPrefix(e.Pos))
		}
	}
	var initV int32 = 0
	if e.Init.Data != nil {
		if t, ok := e.Init.Data.(int32); ok {
			initV = t
		}
	}
	for k, v := range e.Enums {
		v.Value = int32(k) + initV
	}
	return nil
}
