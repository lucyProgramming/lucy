package parser

import (
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

var (
	untils_rc = map[int]bool{
		lex.TOKEN_RC: true,
	}
	untils_semicolon = map[int]bool{
		lex.TOKEN_SEMICOLON: true,
	}
	untils_block_statement = map[int]bool{
		lex.TOKEN_RC:        true,
		lex.TOKEN_SEMICOLON: true,
	}
)
