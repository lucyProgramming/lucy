package lex

func New(bs []byte, startLine, startColumn int) *Lexer {
	lex := &Lexer{bs: bs}
	lex.end = len(bs)
	lex.line = startLine
	lex.column = startColumn
	lex.lastLine = startLine
	lex.lastColumn = startColumn
	return lex
}
