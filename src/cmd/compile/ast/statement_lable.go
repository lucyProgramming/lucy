package ast

import (
	"errors"
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementLabel struct {
	CodeOffsetGenerated bool
	CodeOffset          int
	Block               *Block
	Name                string
	Exits               []*cg.Exit
	Statement           *Statement
}

func (s *StatementLabel) Ready(from *Pos) error {
	ss := []*Statement{}
	for _, v := range s.Block.Statements {
		if v.StatementLabel == s { // this is me
			break
		}
		if v.isVariableDefinition() && v.Checked == false {
			ss = append(ss, v)
		}
	}
	if len(ss) == 0 {
		return nil
	}
	errMsg := fmt.Sprintf("%s cannot jump over variable definition:\n", errMsgPrefix(from))
	for _, v := range ss {
		errMsg += fmt.Sprintf("\t%s constains variable definition\n", errMsgPrefix(v.Pos))
	}
	return errors.New(errMsg)
}
