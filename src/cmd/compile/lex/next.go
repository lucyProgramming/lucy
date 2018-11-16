package lex

import "fmt"

func (this *Lexer) Next() (token *Token, err error) {
	token = &Token{}
	c, eof := this.getChar()
	if eof {
		token.Type = TokenEof
		token.Description = "EOF"
		return
	}
	for c == ' ' ||
		c == '\t' ||
		c == '\r' { // skip empty
		c, eof = this.getChar()

	}

	if eof {
		token.Type = TokenEof
		token.Description = "EOF"
		return
	}
	token.StartLine = this.line
	token.StartColumn = this.column - 1
	if this.isLetter(c) || c == '_' || c == '$' { // start of a identifier
		return this.lexIdentifier(c)
	}
	if this.isDigit(c) {
		eof, err = this.lexNumber(token, c)
		return
	}
	token.Offset = this.offset
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
		c, eof = this.getChar()
		if c == '&' {
			token.Type = TokenLogicalAnd
			token.Description = "&&"
		} else if c == '=' {
			token.Type = TokenAndAssign
			token.Description = "&="
		} else {
			this.unGetChar()
			token.Type = TokenAnd
			token.Description = "&"
		}
	case '|':
		c, eof = this.getChar()
		if c == '|' {
			token.Type = TokenLogicalOr
			token.Description = "||"
		} else if c == '=' {
			token.Type = TokenOrAssign
			token.Description = "|="
		} else {
			this.unGetChar()
			token.Type = TokenOr
			token.Description = "|"
		}
	case '=':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenEqual
			token.Description = "=="
		} else {
			this.unGetChar()
			token.Type = TokenAssign
			token.Description = "="
		}
	case '!':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenNe
			token.Description = "!="
		} else {
			this.unGetChar()
			token.Type = TokenNot
			token.Description = "!"
		}
	case '>':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenGe
			token.Description = ">="
		} else if c == '>' {
			c, eof = this.getChar()
			if c == '=' {
				token.Type = TokenRshAssign
				token.Description = ">>="
			} else {
				this.unGetChar()
				token.Type = TokenRsh
				token.Description = ">>"
			}
		} else {
			this.unGetChar()
			token.Type = TokenGt
			token.Description = ">"
		}
	case '<':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenLe
			token.Description = "<="
		} else if c == '<' {
			c, eof = this.getChar()
			if c == '=' {
				token.Type = TokenLshAssign
				token.Description = "<<="
			} else {
				this.unGetChar()
				token.Type = TokenLsh
				token.Description = "<<"
			}
		} else {
			this.unGetChar()
			token.Type = TokenLt
			token.Description = "<"
		}
	case '^':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenXorAssign
			token.Description = "^="
		} else {
			this.unGetChar()
			token.Type = TokenXor
			token.Description = "^"
		}
	case '~':
		token.Type = TokenBitNot
		token.Description = "~"
	case '+':
		c, eof = this.getChar()
		if c == '+' {
			token.Type = TokenIncrement
			token.Description = "++"
		} else if c == '=' {
			token.Type = TokenAddAssign
			token.Description = "+="
		} else {
			this.unGetChar()
			token.Type = TokenAdd
			token.Description = "+"
		}
	case '-':
		c, eof = this.getChar()
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
			this.unGetChar()
			token.Type = TokenSub
			token.Description = "-"
		}
	case '*':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenMulAssign
			token.Description = "*="
		} else {
			this.unGetChar()
			token.Type = TokenMul
			token.Description = "*"
		}
	case '%':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenModAssign
			token.Description = "%="
		} else {
			this.unGetChar()
			token.Type = TokenMod
			token.Description = "%"
		}
	case '/':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenDivAssign
			token.Description = "/="
		} else if c == '/' {
			bs := []byte{}
			for c != '\n' && eof == false {
				c, eof = this.getChar()
				bs = append(bs, c)
			}
			token.Type = TokenComment
			token.Data = string(bs)
			token.Description = string(bs)
		} else if c == '*' {
			comment, err := this.lexMultiLineComment()
			if err != nil {
				return nil, err
			}
			token.Type = TokenMultiLineComment
			token.Data = comment
			token.Description = comment
		} else {
			this.unGetChar()
			token.Type = TokenDiv
			token.Description = "/"
		}
	case '\n':
		token.Type = TokenLf
		token.Description = "\\n"
	case '.':
		if this.lexVArgs() {
			token.Type = TokenVArgs
			token.Description = "..."
		} else {
			token.Type = TokenSelection
			token.Description = "."
		}
	case '`':
		bs := []byte{}
		c, eof = this.getChar()
		for c != '`' && eof == false {
			bs = append(bs, c)
			c, eof = this.getChar()
		}
		token.Type = TokenLiteralString
		token.Data = string(bs)
		token.Description = string(bs)
	case '"':
		return this.lexString('"')
	case '\'':
		isChar := this.isChar()
		token, err = this.lexString('\'')
		if err == nil {
			if t := []rune(token.Data.(string)); len(t) != 1 {
				err = fmt.Errorf("expect one char")
			} else { // correct token
				if isChar {
					token.Type = TokenLiteralChar
					token.Data = int64(t[0])
				} else {
					token.Type = TokenLiteralByte
					token.Data = int64(t[0])
				}
			}
		}
		return
	case ':':
		c, eof = this.getChar()
		if c == '=' {
			token.Type = TokenVarAssign
			token.Description = ":="
		} else if c == ':' {
			token.Type = TokenSelectConst
			token.Description = "::"
		} else {
			token.Type = TokenColon
			token.Description = ":"
			this.unGetChar()
		}
	default:
		err = fmt.Errorf("unkown beginning of token:%x", c)

		panic(err)
		return nil, err
	}
	token.EndLine = this.line
	token.EndColumn = this.column
	return
}
