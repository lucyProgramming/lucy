package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementLable struct {
	Block            *Block
	StatementsOffset int
	Name             string
	Pos              *Pos
	BackPatches      []*cg.JumpBackPatch
}

func (s *StatementLable) Ready(from *Pos) error {
	ss := []*Statement{}
	for i := 0; i < s.StatementsOffset; i++ {
		if s.Block.Statements[i].isVariableDefinition() && s.Block.Statements[i].Checked == false {
			ss = append(ss, s.Block.Statements[i])
		}
	}
	if len(ss) == 0 {
		return nil
	}
	errmsg := fmt.Sprintf("%s cannot jump over variable definition:", errMsgPrefix(from))
	for _, v := range ss {
		errmsg += fmt.Sprintf("\t%s constains variable definition", errMsgPrefix(v.Pos))
	}
	return errors.New(errmsg)
}
