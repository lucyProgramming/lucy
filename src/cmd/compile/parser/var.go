package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

var (
	untils_lc = map[int]bool{
		lex.TOKEN_LC: true,
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

func (p *Parser) unexpectedErr() {
	p.errs = append(p.errs, p.mkUnexpectedEofErr())
}
func (p *Parser) mkUnexpectedEofErr() error {
	return fmt.Errorf("%s unexpected EOF", p.errorMsgPrefix())
}
