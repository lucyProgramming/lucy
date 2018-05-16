package lex

func New(bs []byte, startLine, startColumn int) *LucyLexer {
	lex := &LucyLexer{bs: bs}
	lex.end = len(bs)
	lex.line = startLine
	lex.column = startColumn
	lex.lastline = startLine
	lex.lastcolumn = startColumn
	return lex
}
