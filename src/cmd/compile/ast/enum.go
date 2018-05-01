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

func (e *Enum) check() error {
	if e.Init == nil {
		e.Init = &Expression{}
		e.Init.Typ = EXPRESSION_TYPE_INT
		e.Init.Data = int32(0)
		e.Pos = e.Pos
	}
	is, typ, value, err := e.Init.getConstValue()
	if err != nil || is == false || typ != EXPRESSION_TYPE_INT {
		return fmt.Errorf("%s enum type must inited by integer", errMsgPrefix(e.Pos))
	}
	for k, v := range e.Enums {
		v.Value = int32(k) + value.(int32)
	}
	return nil
}
