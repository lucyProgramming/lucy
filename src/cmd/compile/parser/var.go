package parser

import (
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
