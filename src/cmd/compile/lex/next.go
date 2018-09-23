package lex

import "fmt"

func (lex *Lexer) Next() (token *Token, err error) {
	token = &Token{}
	var c byte
	c, eof := lex.getChar()
	if eof {
		token.Type = TokenEof
		token.Description = "EOF"
		return
	}
	for c == ' ' || c == '\t' || c == '\r' { // skip empty
		c, eof = lex.getChar()
	}
	token.StartLine = lex.line
	token.StartColumn = lex.column - 1
	if eof {
		token.Type = TokenEof
		token.Description = "EOF"
		return
	}
	if lex.isLetter(c) || c == '_' || c == '$' { // start of a identifier
		return lex.lexIdentifier(c)
	}
	if lex.isDigit(c) {
		eof, err = lex.lexNumber(token, c)
		return
	}
	token.Offset = lex.offset
	switch c {
	case '?':
		token.Type = TokenQuestion
		token.Description = "?"
	case '(':
		token.Type = TokenLp
		token.Description = "("
	case ')':
		token.Type = TokenRp
		token.Description = ")"
	case '{':
		token.Type = TokenLc
		token.Description = "{"
	case '}':
		token.Type = TokenRc
		token.Description = "}"
	case '[':
		token.Type = TokenLb
		token.Description = "["
	case ']':
		token.Type = TokenRb
		token.Description = "]"
	case ';':
		token.Type = TokenSemicolon
		token.Description = ";"
	case ',':
		token.Type = TokenComma
		token.Description = ","
	case '&':
		c, eof = lex.getChar()
		if c == '&' {
			token.Type = TokenLogicalAnd
			token.Description = "&&"
		} else if c == '=' {
			token.Type = TokenAndAssign
			token.Description = "&="
		} else {
			lex.unGetChar()
			token.Type = TokenAnd
			token.Description = "&"
		}
	case '|':
		c, eof = lex.getChar()
		if c == '|' {
			token.Type = TokenLogicalOr
			token.Description = "||"
		} else if c == '=' {
			token.Type = TokenOrAssign
			token.Description = "|="
		} else {
			lex.unGetChar()
			token.Type = TokenOr
			token.Description = "|"
		}
	case '=':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenEqual
			token.Description = "=="
		} else {
			lex.unGetChar()
			token.Type = TokenAssign
			token.Description = "="
		}
	case '!':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenNe
			token.Description = "!="
		} else {
			lex.unGetChar()
			token.Type = TokenNot
			token.Description = "!"
		}
	case '>':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenGe
			token.Description = ">="
		} else if c == '>' {
			c, eof = lex.getChar()
			if c == '=' {
				token.Type = TokenRshAssign
				token.Description = ">>="
			} else {
				lex.unGetChar()
				token.Type = TokenRsh
				token.Description = ">>"
			}
		} else {
			lex.unGetChar()
			token.Type = TokenGt
			token.Description = ">"
		}
	case '<':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenLe
			token.Description = "<="
		} else if c == '<' {
			c, eof = lex.getChar()
			if c == '=' {
				token.Type = TokenLshAssign
				token.Description = "<<="
			} else {
				lex.unGetChar()
				token.Type = TokenLsh
				token.Description = "<<"
			}
		} else {
			lex.unGetChar()
			token.Type = TokenLt
			token.Description = "<"
		}
	case '^':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenXorAssign
			token.Description = "^="
		} else {
			lex.unGetChar()
			token.Type = TokenXor
			token.Description = "^"
		}
	case '~':
		token.Type = TokenBitNot
		token.Description = "~"
	case '+':
		c, eof = lex.getChar()
		if c == '+' {
			token.Type = TokenIncrement
			token.Description = "++"
		} else if c == '=' {
			token.Type = TokenAddAssign
			token.Description = "+="
		} else {
			lex.unGetChar()
			token.Type = TokenAdd
			token.Description = "+"
		}
	case '-':
		c, eof = lex.getChar()
		if c == '-' {
			token.Type = TokenDecrement
			token.Description = "--"
		} else if c == '=' {
			token.Type = TokenSubAssign
			token.Description = "-="
		} else if c == '>' {
			token.Type = TokenArrow
			token.Description = "->"
		} else {
			lex.unGetChar()
			token.Type = TokenSub
			token.Description = "-"
		}
	case '*':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenMulAssign
			token.Description = "*="
		} else {
			lex.unGetChar()
			token.Type = TokenMul
			token.Description = "*"
		}
	case '%':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenModAssign
			token.Description = "%="
		} else {
			lex.unGetChar()
			token.Type = TokenMod
			token.Description = "%"
		}
	case '/':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenDivAssign
			token.Description = "/="
		} else if c == '/' {
			bs := []byte{}
			for c != '\n' && eof == false {
				c, eof = lex.getChar()
				bs = append(bs, c)
			}
			token.Type = TokenComment
			token.Data = string(bs)
			token.Description = string(bs)
		} else if c == '*' {
			comment, err := lex.lexMultiLineComment()
			if err != nil {
				return nil, err
			}
			token.Type = TokenCommentMultiLine
			token.Data = comment
			token.Description = comment
		} else {
			lex.unGetChar()
			token.Type = TokenDiv
			token.Description = "/"
		}
	case '\n':
		token.Type = TokenLf
		token.Description = "\\n"
	case '.':
		if lex.lexVArgs() {
			token.Type = TokenVArgs
			token.Description = "..."
		} else {
			token.Type = TokenSelection
			token.Description = "."
		}
	case '`':
		bs := []byte{}
		c, eof = lex.getChar()
		for c != '`' && eof == false {
			bs = append(bs, c)
			c, eof = lex.getChar()
		}
		token.Type = TokenLiteralString
		token.Data = string(bs)
		token.Description = string(bs)
	case '"':
		return lex.lexString('"')
	case '\'':
		isChar := lex.isChar()
		token, err = lex.lexString('\'')
		if err == nil {
			if t := []rune(token.Data.(string)); len(t) != 1 {
				err = fmt.Errorf("expect one char")
			} else { // correct token
				if isChar {
					token.Type = TokenLiteralChar
					token.Data = int32(t[0])
				} else {
					token.Type = TokenLiteralByte
					token.Data = byte(t[0])
				}
			}
		}
		return
	case ':':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TokenVarAssign
			token.Description = ":="
		} else if c == ':' {
			token.Type = Token2Colon
			token.Description = "::"
		} else {
			token.Type = TokenColon
			token.Description = ":"
			lex.unGetChar()
		}
	default:
		err = fmt.Errorf("unkown beginning of token:%d", c)
		return
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	return
}
