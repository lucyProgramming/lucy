package ast

import (
	"errors"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type StatementLabel struct {
	Used                bool
	CodeOffsetGenerated bool
	CodeOffset          int
	Block               *Block
	Name                string
	Exits               []*cg.Exit
	Statement           *Statement
	Pos                 *Pos
}

/*
	defer {
		xxx:
	}
	defer block could be compile multi times,
	should reset the label
*/
func (l *StatementLabel) Reset() {
	l.CodeOffsetGenerated = false
	l.CodeOffset = -1
	l.Exits = []*cg.Exit{}
}

// check this label is read to goto
func (l *StatementLabel) Ready(from *Pos) error {
	ss := []*Statement{}
	for _, v := range l.Block.Statements {
		if v.StatementLabel == l { // this is me
			break
		}
		if v.isVariableDefinition() && v.Checked == false {
			ss = append(ss, v)
		}
	}
	if len(ss) == 0 {
		return nil
	}
	/*
		if false {
			if true {
				// jump over variable definition not allow
				goto some ;
			}
		}
		a := false ;
		some:
	*/
	errMsg := fmt.Sprintf("%s cannot jump over variable definition:\n", from.ErrMsgPrefix())
	for _, v := range ss {
		errMsg += fmt.Sprintf("\t%s constains variable definition\n", v.Pos.ErrMsgPrefix())
	}
	return errors.New(errMsg)
}
