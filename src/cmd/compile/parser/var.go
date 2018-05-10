package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func init() {
	ast.ParseFunctionHandler = ParseFunction
}

var (
	untils_lc = map[int]bool{
		lex.TOKEN_LC: true,
	}
	untils_gt = map[int]bool{
		lex.TOKEN_GT: true,
	}
	untils_rc = map[int]bool{
		lex.TOKEN_RC: true,
	}
	untils_semicolon = map[int]bool{
		lex.TOKEN_SEMICOLON: true,
	}
	untils_rc_semicolon = map[int]bool{
		lex.TOKEN_RC:        true,
		lex.TOKEN_SEMICOLON: true,
	}
)

func ParseFunction(bs []byte, pos *ast.Pos) (*ast.Function, []error) {
	p := &Parser{}
	p.filename = pos.Filename
	p.nerr = 10
	tops := []*ast.Node{}
	p.tops = &tops
	p.bs = bs
	es := p.Parse(pos.StartLine, pos.StartColumn)
	if es != nil && len(es) > 0 {
		return nil, es
	}
	if len(tops) == 0 {
		return nil, []error{fmt.Errorf("function not found")}
	}
	f, ok := tops[0].Data.(*ast.Function)
	if ok == false {
		return nil, []error{fmt.Errorf("not a function")}
	}
	return f, nil
}
