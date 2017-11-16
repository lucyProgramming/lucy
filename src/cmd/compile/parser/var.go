package parser

import (
	"github.com/756445638/lucy/src/cmd/compile/lex"
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
