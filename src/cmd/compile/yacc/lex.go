package yacc

import (
	"github.com/756445638/lucy/src/cmd/compile/lex"
	"log"
)

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type lucyLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (x *lucyLex) Lex(yylval *lucySymType) int {

}

// The parser calls this method on a parse error.
func (x *lucyLex) Error(s string) {
	log.Printf("parse error: %s", s)
}
