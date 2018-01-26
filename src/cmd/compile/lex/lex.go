package lex

import (
	"fmt"
)

func New(bs []byte) *LucyLexer {
	lex := &LucyLexer{bs: bs}
	lex.end = len(bs) - 1
	if lex.end == -1 {
		lex.end = 0
	}
	return lex
}

type LucyLexer struct {
	bs                                 []byte
	lastline, lastcolumn, line, column int
	offset, end                        int
}

func (lex *LucyLexer) incrementLine() {
	lex.lastline = lex.line
	lex.line++
}

func (lex *LucyLexer) getchar() (c byte, eof bool) {
	if lex.offset == lex.end {
		eof = true
		return
	}
	lex.lastcolumn = lex.column
	lex.column++
	offset := lex.offset
	lex.offset++
	c = lex.bs[offset]
	return
}
func (lex *LucyLexer) isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
func (lex *LucyLexer) isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
func (lex *LucyLexer) isOctal(c byte) bool {
	return '0' <= c && c <= '7'
}
func (lex *LucyLexer) isHex(c byte) bool {
	return lex.isDigit(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}
func (lex *LucyLexer) hexbyte2Byte(c byte) byte {
	if 'a' <= c && c <= 'f' {
		return c - 'a'
	}
	if 'A' <= c && c <= 'F' {
		return c - 'A'
	}
	return c - '0'
}

func (lex *LucyLexer) ungetchar() {
	lex.offset--
	lex.line, lex.column = lex.lastline, lex.lastcolumn
}
func (lex *LucyLexer) lexNumber(c byte) (token *Token, eof bool, err error) {
	return
}
func (lex *LucyLexer) lexIdentifier(c byte) (token *Token, eof bool, err error) {
	token = &Token{}
	token.StartLine = lex.line
	token.StartLine = lex.column
	bs := []byte{c}
	c, eof = lex.getchar()
	for eof == false {
		if lex.isLetter(c) || c == '_' || lex.isDigit(c) {
			bs = append(bs, c)
			c, eof = lex.getchar()
		} else {
			break
		}
	}
	lex.ungetchar()
	token.EndLine = lex.line
	token.EndColumn = lex.column
	identifier := string(bs)
	if t, ok := keywordMap[identifier]; ok {
		token.Type = t
		return
	}
	token.Type = TOKEN_IDENTIFIER
	token.Data = identifier
	token.Desp = identifier
	token.EndLine = lex.line
	token.EndColumn = lex.column
	return
}

func (lex *LucyLexer) lexString() (token *Token, eof bool, err error) {
	token = &Token{}
	token.StartLine = lex.line
	token.StartLine = lex.column
	bs := []byte{}
	var c byte
	c, eof = lex.getchar()
	for c != '"' && c != '\n' && eof == false {
		if c != '\\' {
			bs = append(bs, c)
			c, eof = lex.getchar()
			continue
		}
		c, eof = lex.getchar()
		if eof {
			break
		}
		switch c {
		case 'a':
			bs = append(bs, '\a')
			c, eof = lex.getchar()
		case 'b':
			bs = append(bs, '\b')
			c, eof = lex.getchar()
		case 'f':
			bs = append(bs, '\f')
			c, eof = lex.getchar()
		case 'n':
			bs = append(bs, '\n')
			c, eof = lex.getchar()
		case 'r':
			bs = append(bs, '\r')
			c, eof = lex.getchar()
		case 't':
			bs = append(bs, '\t')
			c, eof = lex.getchar()
		case 'v':
			bs = append(bs, '\v')
			c, eof = lex.getchar()
		case '\\':
			bs = append(bs, '\\')
			c, eof = lex.getchar()
		case '\'':
			bs = append(bs, '\'')
			c, eof = lex.getchar()
		case '"':
			bs = append(bs, '"')
			c, eof = lex.getchar()
		case 'x':
			var c1, c2 byte
			c1, eof = lex.getchar() // skip 'x'
			if eof {
				err = fmt.Errorf("unexpect EOF")
				break
			}
			if !lex.isHex(c) {
				err = fmt.Errorf("unknown escape sequence")
				continue
			}
			b := lex.hexbyte2Byte(c1)
			c2, eof = lex.getchar()
			if lex.isHex(c2) {
				b2 := lex.hexbyte2Byte(c2)
				if t := b*16 + b2; t > 127 { // only support standard accii
					bs = append(bs, b)
					lex.ungetchar()
				} else {
					bs = append(bs, t)
					lex.ungetchar()
				}
			} else {
				bs = append(bs, b)
				lex.ungetchar()
			}
		case '0', '1', '2', '3', '4', '5', '7':
			b := byte(0)
			for i := 0; i < 3; i++ {
				bb := lex.hexbyte2Byte(c)
				b = b*16 + bb
				if b > 127 {
					break
				}
				c, eof = lex.getchar()
				if lex.isOctal(c) == false || eof {
					break
				}
			}
			bs = append(bs, b)
		}
	}
	if c == '\n' && eof == false {
		lex.incrementLine()
	}
	return
}
func (lex *LucyLexer) lexMultiLineComment() {
redo:
	c, eof := lex.getchar()
	if eof {
		return
	}
	if c == '\n' {
		lex.incrementLine()
	}
	for c != '*' && eof == false {
		c, eof = lex.getchar()
		if c == '\n' {
			lex.incrementLine()
		}
	}
	if c == '\n' {
		lex.incrementLine()
	}
	c, eof = lex.getchar()
	if c == '/' {
		return
	}
	if c == '\n' {
		lex.incrementLine()
	}
	goto redo
}

func (lex *LucyLexer) Next() (token *Token, eof bool, err error) {
redo:
	var c byte
	c, eof = lex.getchar()
	if eof {
		return
	}
	for c == ' ' || c == '\t' || c == '\r' {
		c, eof = lex.getchar()
	}
	if eof {
		return
	}
	if lex.isLetter(c) || c == '_' {
		return lex.lexIdentifier(c)
	}
	if lex.isDigit(c) {
		return lex.lexNumber(c)
	}
	token = &Token{}
	token.StartLine = lex.line
	token.StartColumn = lex.column
	switch c {
	case '(':
		token.Type = TOKEN_LP
	case ')':
		token.Type = TOKEN_RP
	case '{':
		token.Type = TOKEN_LC
	case '}':
		token.Type = TOKEN_RC
	case '[':
		token.Type = TOKEN_LB
	case ']':
		token.Type = TOKEN_RB
	case ';':
		token.Type = TOKEN_SEMICOLON
	case ',':
		token.Type = TOKEN_COMMA
	case '&':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_AND
			break
		}
		if c == '&' {
			token.Type = TOKEN_LOGICAL_AND
		} else {
			lex.ungetchar()
			token.Type = TOKEN_AND
		}
	case '|':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_OR
			break
		}
		if c == '|' {
			token.Type = TOKEN_LOGICAL_OR
		} else {
			lex.ungetchar()
			token.Type = TOKEN_OR
		}
	case '=':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_ASSIGN
			break
		}
		if c == '=' {
			token.Type = TOKEN_EQUAL
		} else {
			lex.ungetchar()
			token.Type = TOKEN_ASSIGN
		}
	case '!':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_NOT
			break
		}
		if c == '=' {
			token.Type = TOKEN_NE
		} else {
			lex.ungetchar()
			token.Type = TOKEN_NOT
		}
	case '>':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_GT
			break
		}
		if c == '=' {
			token.Type = TOKEN_GE
		} else if c == '>' {
			token.Type = TOKEN_RIGHT_SHIFT
		} else {
			lex.ungetchar()
			token.Type = TOKEN_GT
		}
	case '<':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_GT
			break
		}
		if c == '=' {
			token.Type = TOKEN_LE
		} else if c == '<' {
			token.Type = TOKEN_LEFT_SHIFT
		} else {
			lex.ungetchar()
			token.Type = TOKEN_LT
		}
	case '+':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_ADD
			break
		}
		if c == '+' {
			token.Type = TOKEN_INCREMENT
		} else if c == '=' {
			token.Type = TOKEN_ADD_ASSIGN
		} else {
			lex.ungetchar()
			token.Type = TOKEN_ADD
		}
	case '-':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_SUB
			break
		}
		if c == '-' {
			token.Type = TOKEN_DECREMENT
		} else if c == '=' {
			token.Type = TOKEN_SUB_ASSIGN
		} else if c == '>' {
			token.Type = TOKEN_ARROW
		} else {
			lex.ungetchar()
			token.Type = TOKEN_SUB
		}
	case '*':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_MUL
			break
		}
		if c == '=' {
			token.Type = TOKEN_MUL_ASSIGN
		} else {
			lex.ungetchar()
			token.Type = TOKEN_MUL
		}
	case '%':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_MOD
			break
		}
		if c == '=' {
			token.Type = TOKEN_MOD_ASSIGN
		} else {
			lex.ungetchar()
			token.Type = TOKEN_MOD
		}
	case '/':
		c, eof = lex.getchar()
		if eof == true {
			token.Type = TOKEN_DIV
			break
		}
		if c == '=' {
			token.Type = TOKEN_DIV_ASSIGN
		} else if c == '/' {
			for c != '\n' && eof == false {
				c, eof = lex.getchar()
			}
			lex.incrementLine()
			goto redo
		} else if c == '*' {
			lex.lexMultiLineComment()
			goto redo
		} else {
			lex.ungetchar()
			token.Type = TOKEN_DIV
		}
	case '\n':
		lex.incrementLine()
		token.Type = TOKEN_CRLF
	case '.':
		token.Type = TOKEN_DOT
	case '`':
		bs := []byte{}
		c, eof = lex.getchar()
		for c != '`' && eof == false {
			bs = append(bs, c)
			c, eof = lex.getchar()
		}
		token.Type = TOKEN_LITERAL_STRING
		token.Data = string(bs)
		token.Desp = string(bs)
	case '"':
		return lex.lexString()
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	return
}
